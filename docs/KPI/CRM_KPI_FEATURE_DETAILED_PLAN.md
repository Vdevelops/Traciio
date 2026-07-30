# CRM KPI Feature — Detailed Implementation Plan (v2)

> Dokumen ini adalah versi diperdalam dari plan awal. Tujuannya: setiap keputusan desain punya formula/pseudocode/schema eksplisit, sehingga bisa langsung dipakai sebagai input ke AI coding assistant (Cursor/Claude Code) tanpa ambiguitas.

---

## 0. Ringkasan Eksekutif

Fitur ini memperluas modul **leaderboard/KPI** yang sudah ada di repo, bukan membuat jalur baru. Output berlapis dua:

1. **Scorecard mentah** — angka metrik individual, dapat diverifikasi langsung ke data source.
2. **Weighted composite score + evaluasi** — lapisan penilaian (grade) dan tren, dibangun di atas scorecard, dengan diagnostic notes yang actionable.

Dua konsumen utama: **Sales Rep** (scope: diri sendiri) dan **Sales Manager** (scope: tim dalam brick yang di-manage + agregasi per brick/area).

---

## 1. KPI Taxonomy — Formula Eksplisit

Setiap metrik didefinisikan dengan formula, source table, filter wajib, dan role visibility. Ini adalah **source of truth** — AI assistant tidak boleh mengarang formula lain di luar tabel ini.

| # | Metric | Formula | Source Table(s) | Filter Wajib | Visible untuk |
|---|--------|---------|------------------|--------------|----------------|
| 1 | Total Deals Closed | `COUNT(deals) WHERE status = 'closed_won'` | `deals` | `assigned_to = scope_user_id`, `closed_at BETWEEN start_date AND end_date` | Rep, Manager (per rep & agregat) |
| 2 | Total Revenue | `SUM(deals.value) WHERE status = 'closed_won'` | `deals` | sama seperti #1 | Rep, Manager |
| 3 | Deals Created (denominator conversion) | `COUNT(deals)` | `deals` | `assigned_to = scope_user_id`, `created_at BETWEEN start_date AND end_date` | Rep, Manager |
| 4 | Conversion Rate | `(Total Deals Closed / Deals Created) * 100` | derived | jika `Deals Created = 0` → return `null`, bukan `0` atau `divide by zero` | Rep, Manager |
| 5 | Average Deal Value | `Total Revenue / Total Deals Closed` | derived | jika `Total Deals Closed = 0` → `null` | Rep, Manager |
| 6 | Visit Completed | `COUNT(visit_reports) WHERE status IN ('approved','completed')` | `visit_reports` | `sales_rep_id = scope_user_id`, `visit_date BETWEEN start_date AND end_date` | Rep, Manager |
| 7 | Visit Planned | `COUNT(schedules)` atau target dari `monthly_targets.visit_target` | `schedules` / `monthly_targets` | periode sama, `assigned_to = scope_user_id` | Rep, Manager |
| 8 | Visit Compliance | `(Visit Completed / Visit Planned) * 100` | derived | jika `Visit Planned = 0` → `null` | Rep, Manager |
| 9 | Tasks Completed | `COUNT(tasks) WHERE status = 'completed'` | `tasks` | `assigned_to = scope_user_id`, `completed_at BETWEEN start_date AND end_date` | Rep, Manager |
| 10 | Overdue Task Rate | `(COUNT(tasks WHERE due_date < NOW() AND status != 'completed') / COUNT(tasks assigned in period)) * 100` | `tasks` | `assigned_to = scope_user_id` | Rep, Manager |
| 11 | Target Attainment (Revenue) | `(Total Revenue / monthly_targets.revenue_target) * 100` | `monthly_targets` | match by `user_id` + `month/year` dalam range | Rep, Manager |
| 12 | Target Attainment (Deals) | `(Total Deals Closed / monthly_targets.deal_target) * 100` | `monthly_targets` | sama seperti #11 | Rep, Manager |
| 13 | Pipeline Movement Score | Lihat §1.1 | `deals` (stage history) | `assigned_to = scope_user_id`, transisi dalam periode | Rep, Manager |
| 14 | Territory/Coverage Penetration | `(DISTINCT customers dengan interaksi dalam brick / total customers terdaftar di brick) * 100` | `area_mapping`, `deals`/`visit_reports` | `brick_id = scope_brick_id` | **Manager only** (brick-level) |

### 1.1 Pipeline Movement Score (definisi detail)

Karena "pergerakan pipeline" ambigu, definisikan sebagai **net forward progression**:

```
pipeline_movement_score =
  Σ (stage_weight[new_stage] - stage_weight[old_stage])
  untuk setiap stage transition dalam periode, per deal

stage_weight contoh (harus dikonfirmasi ke actual pipeline stages di repo):
  lead        = 1
  qualified   = 2
  proposal    = 3
  negotiation = 4
  closed_won  = 5
  closed_lost = 0 (atau exclude dari scoring, tapi tetap dicatat sbg lost)
```

> **Catatan wajib untuk implementer**: cek dulu apakah tabel `deals` punya history transisi stage (misal `deal_stage_history`), atau apakah stage transition harus diturunkan dari `updated_at` + `status` snapshot. Jangan asumsikan ada tabel history sebelum verifikasi ke schema aktual.

---

## 2. Weighted Composite Score

### 2.1 Struktur bobot (harus configurable, bukan hardcode)

```json
{
  "role": "sales_rep",
  "weights": {
    "conversion_rate": 0.25,
    "revenue_target_attainment": 0.25,
    "visit_compliance": 0.20,
    "overdue_task_rate": 0.10,
    "average_deal_value": 0.10,
    "pipeline_movement_score": 0.10
  }
}
```

```json
{
  "role": "sales_manager",
  "weights": {
    "team_target_attainment": 0.30,
    "team_conversion_rate": 0.20,
    "territory_coverage": 0.20,
    "team_visit_compliance": 0.15,
    "team_overdue_task_rate": 0.10,
    "brick_pipeline_movement": 0.05
  }
}
```

- Bobot disimpan di **config table** (misal `kpi_weight_config`) atau minimal named-constant di kode — **bukan angka lepas di tengah logic**. Ini menjawab pertanyaan "further consideration #1": ya, bobot **berbeda per role**, dan idealnya juga bisa di-override per tenant/company kalau sistem ini multi-tenant.
- Setiap metrik dinormalisasi ke skala 0–100 sebelum dikalikan bobot (lihat §2.2), supaya metrik dengan satuan berbeda (rupiah vs persen vs count) bisa dijumlahkan secara adil.

### 2.2 Normalisasi metrik

```
normalized_value = min(100, (raw_value / target_or_benchmark_value) * 100)

- Untuk metrik "rate"/"percentage" yang sudah 0-100 (conversion_rate, visit_compliance,
  target_attainment): pakai langsung, di-cap ke 100 (attainment >100% dianggap 100
  untuk composite score, tapi raw value tetap ditampilkan apa adanya di scorecard).
- Untuk overdue_task_rate: metrik ini "semakin rendah semakin baik", jadi:
  normalized_value = 100 - overdue_task_rate (di-clamp 0-100)
- Untuk average_deal_value: benchmark = rata-rata deal value tim/brick dalam periode
  yang sama (peer benchmark), bukan target absolut.
```

### 2.3 Formula composite

```
composite_score = Σ (normalized_value[metric] * weight[metric])
```

Contoh perhitungan (sales rep):

| Metric | Raw | Normalized | Weight | Contribution |
|---|---|---|---|---|
| conversion_rate | 18% | 72 (target 25%) | 0.25 | 18.0 |
| revenue_target_attainment | 110% | 100 (capped) | 0.25 | 25.0 |
| visit_compliance | 80% | 80 | 0.20 | 16.0 |
| overdue_task_rate | 10% | 90 (100-10) | 0.10 | 9.0 |
| average_deal_value | peer avg | 95 | 0.10 | 9.5 |
| pipeline_movement_score | — | 60 | 0.10 | 6.0 |
| **Composite Score** | | | | **83.5** |

---

## 3. Grading / Evaluation Layer

```
grade_bands = [
  { min: 85, max: 100, label: "Excellent",         color: "green"  },
  { min: 70, max: 84,  label: "Good",               color: "blue"   },
  { min: 55, max: 69,  label: "Needs Improvement",  color: "yellow" },
  { min: 0,  max: 54,  label: "Critical",           color: "red"    }
]
```

- Grade bands juga **configurable** (config table / constants file), jangan hardcode di handler.
- Setiap evaluasi wajib menyertakan **trend indicator** dibanding periode sebelumnya (periode sebelumnya = periode dengan panjang sama, langsung sebelum `start_date`):

```
trend = current_composite_score - previous_composite_score
trend_direction = "up" | "down" | "flat"  (flat jika |trend| < 1.0)
```

---

## 4. Aturan Atribusi Data (dengan Fallback Eksplisit)

```
FUNCTION resolve_brick(entity):
    IF entity.brick_id IS NOT NULL:
        RETURN entity.brick_id
    ELSE IF entity.assigned_to.default_brick_id IS NOT NULL:
        # fallback ke brick default milik user yang di-assign
        RETURN entity.assigned_to.default_brick_id
        MARK entity AS "brick_inferred = true"  # untuk diagnostic notes
    ELSE:
        RETURN "UNASSIGNED"
        MARK entity AS "brick_missing = true"
        EXCLUDE from brick-level aggregation
        INCLUDE in personal-level aggregation (tetap dihitung untuk KPI rep individu)
```

**Aturan tambahan:**
- Data dengan `brick_missing = true` **tidak boleh** memengaruhi coverage/territory penetration (karena itu metrik brick-level), tapi tetap masuk personal scorecard sales rep.
- Manager hanya melihat brick yang secara eksplisit ter-assign ke dia (lihat `brick.manager_id` atau tabel relasi `manager_brick`), bukan brick yang "kebetulan" punya deal dari anak buahnya.
- Jika satu sales rep terdaftar di lebih dari satu brick (multi-brick rep), agregasi personal tetap satu angka (semua brick digabung), tapi breakdown per brick tetap disediakan di response detail.
- Diagnostic notes wajib mencatat berapa banyak record yang kena fallback (`brick_inferred_count`, `brick_missing_count`) agar manager tahu data quality-nya.

---

## 5. Desain API

### 5.1 Endpoint

```
GET /api/v1/kpi/sales-rep
  Query params:
    userId          (optional, default = current authenticated user; hanya admin/manager
                     yang boleh query userId lain milik timnya)
    startDate       (required, ISO date)
    endDate         (required, ISO date)
    compareWithPrevious (bool, default true)

GET /api/v1/kpi/sales-manager
  Query params:
    managerId       (optional, default = current authenticated user)
    startDate       (required)
    endDate         (required)
    brickId         (optional — filter ke satu brick spesifik dalam scope manager)
    includeTeamBreakdown (bool, default true)
    compareWithPrevious (bool, default true)
```

### 5.2 Response Schema — Sales Rep

```json
{
  "scope": {
    "userId": "uuid",
    "role": "sales_rep",
    "startDate": "2026-06-01",
    "endDate": "2026-06-30"
  },
  "scorecard": {
    "totalDealsClosed": 12,
    "totalRevenue": 450000000,
    "dealsCreated": 40,
    "conversionRate": 30.0,
    "averageDealValue": 37500000,
    "visitCompleted": 18,
    "visitPlanned": 20,
    "visitCompliance": 90.0,
    "tasksCompleted": 25,
    "overdueTaskRate": 8.0,
    "revenueTargetAttainment": 112.5,
    "dealTargetAttainment": 100.0,
    "pipelineMovementScore": 34
  },
  "evaluation": {
    "compositeScore": 83.5,
    "grade": "Good",
    "trend": {
      "previousCompositeScore": 78.0,
      "delta": 5.5,
      "direction": "up"
    },
    "targetGap": {
      "revenue": { "target": 400000000, "actual": 450000000, "gapPercent": 12.5, "status": "above" },
      "deals":   { "target": 12, "actual": 12, "gapPercent": 0, "status": "met" }
    }
  },
  "diagnostics": [
    { "code": "LOW_VISIT_CADENCE", "severity": "info", "message": "Visit compliance 90%, mendekati target." }
  ],
  "meta": {
    "brickMissingCount": 0,
    "brickInferredCount": 1,
    "generatedAt": "2026-07-30T10:00:00Z"
  }
}
```

### 5.3 Response Schema — Sales Manager

```json
{
  "scope": {
    "managerId": "uuid",
    "role": "sales_manager",
    "startDate": "2026-06-01",
    "endDate": "2026-06-30",
    "bricks": ["brick-A", "brick-B"]
  },
  "teamSummary": {
    "totalRepsCount": 6,
    "totalDealsClosed": 55,
    "totalRevenue": 1800000000,
    "teamConversionRate": 27.4,
    "teamVisitCompliance": 81.2,
    "teamOverdueTaskRate": 14.0,
    "teamTargetAttainment": 96.0
  },
  "evaluation": {
    "compositeScore": 76.2,
    "grade": "Good",
    "trend": { "previousCompositeScore": 74.0, "delta": 2.2, "direction": "up" }
  },
  "teamBreakdown": [
    {
      "userId": "uuid-rep-1",
      "name": "Budi",
      "compositeScore": 88.0,
      "grade": "Excellent",
      "totalRevenue": 500000000,
      "conversionRate": 35.0,
      "rank": 1
    }
  ],
  "brickBreakdown": [
    {
      "brickId": "brick-A",
      "name": "Brick Semarang Timur",
      "coveragePenetration": 62.5,
      "totalRevenue": 900000000,
      "repsCount": 3,
      "compositeScore": 79.0
    }
  ],
  "topBottomPerformers": {
    "top": ["uuid-rep-1"],
    "bottom": ["uuid-rep-5"]
  },
  "diagnostics": [
    { "code": "COVERAGE_DECLINE", "severity": "warning", "brickId": "brick-B", "message": "Coverage turun 8% dibanding periode lalu." }
  ],
  "meta": {
    "brickMissingCount": 2,
    "brickInferredCount": 5,
    "generatedAt": "2026-07-30T10:00:00Z"
  }
}
```

> **Menjawab further consideration #2**: manager mendapat **keduanya** — `teamBreakdown` (leaderboard internal lengkap per rep, bukan cuma top/bottom) DAN `topBottomPerformers` (ringkasan cepat). `teamBreakdown` sudah cukup untuk UI menyusun leaderboard sendiri, jadi tidak perlu endpoint terpisah.

---

## 6. Security & Scope Enforcement

```
FUNCTION authorize_kpi_request(requester, requested_scope):
    IF requester.role == "sales_rep":
        IF requested_scope.userId != requester.id:
            RETURN 403 Forbidden
    IF requester.role == "sales_manager":
        IF requested_scope.type == "sales-rep":
            IF requested_scope.userId NOT IN get_team_members(requester.id):
                RETURN 403 Forbidden
        IF requested_scope.type == "sales-manager":
            IF requested_scope.managerId != requester.id:
                RETURN 403 Forbidden
            IF requested_scope.brickId NOT IN get_managed_bricks(requester.id):
                RETURN 403 Forbidden
    # admin/superadmin role: bebas akses semua scope (jika ada di sistem)
```

- Implementasikan sebagai **middleware/service-layer guard**, bukan dicek manual di tiap handler — reuse pola RBAC yang sudah ada di `dashboard_handler.go` kalau memang sudah ada pattern serupa (verifikasi dulu ke file aslinya).
- Jangan percaya `userId`/`managerId`/`brickId` dari query param tanpa validasi — ini rawan IDOR (Insecure Direct Object Reference).

---

## 7. Diagnostic / Actionable Insight Rules

> Menjawab further consideration #3: **ya**, KPI harus menyertakan rekomendasi tindakan otomatis berbasis rule table berikut (bukan hardcoded if-else berserakan — taruh di satu rule engine/table agar mudah ditambah).

| Code | Kondisi | Severity | Message Template |
|---|---|---|---|
| `LOW_CONVERSION` | `conversion_rate < 15%` AND `deals_created >= 10` | warning | "Conversion rate rendah ({rate}%) dari {count} deals dibuat." |
| `TARGET_UNDERPERFORM` | `revenue_target_attainment < 70%` | critical | "Pencapaian target revenue baru {attainment}%." |
| `LOW_VISIT_CADENCE` | `visit_compliance < 60%` | warning | "Visit compliance {compliance}%, di bawah standar minimum." |
| `HIGH_OVERDUE_TASK` | `overdue_task_rate > 25%` | critical | "{rate}% task melewati due date." |
| `COVERAGE_DECLINE` | `coverage_penetration` turun >5 poin vs periode lalu | warning | "Coverage brick {brickId} turun {delta}% dibanding periode sebelumnya." |
| `DATA_QUALITY_ISSUE` | `brickMissingCount > 0` | info | "{count} record tidak punya brick_id dan dikecualikan dari agregasi brick." |
| `STAGNANT_PIPELINE` | `pipeline_movement_score <= 0` AND `deals_created > 0` | warning | "Tidak ada progres stage berarti pada deal yang dibuat periode ini." |

---

## 8. Rencana Implementasi Backend (Go, Clean Architecture)

Mengikuti pola yang sudah ada di repo. **Jangan buat layer baru yang meniru dashboard** — KPI punya modul sendiri.

### Fase 1 — Fondasi & Scorecard Mentah
1. `apps/api/internal/domain/kpi/entity.go` — definisikan struct metrik mentah per §1.
2. `apps/api/internal/domain/kpi/dto.go` — DTO request/response sesuai schema §5.2 dan §5.3.
3. `apps/api/internal/repository/kpi/repository.go` — query terpisah per metrik (jangan satu query raksasa yang sulit di-test), gunakan pola yang sudah ada di `sales_overview/service.go` sebagai referensi.
4. `apps/api/internal/service/kpi/service.go` — orchestrate repository calls, terapkan attribution fallback (§4).
5. Verifikasi field asli di `deals`, `visit_reports`, `tasks`, `monthly_targets`, `bricks` sebelum menulis query — **jangan asumsikan nama kolom**.

### Fase 2 — Evaluation, Composite Score, Diagnostics
6. `apps/api/internal/service/kpi/scoring.go` — normalisasi (§2.2), composite score (§2.3), grading (§3).
7. `apps/api/internal/service/kpi/diagnostics.go` — rule engine dari tabel §7 (gunakan slice of rule struct + evaluator function, bukan if-else panjang).
8. Config: `apps/api/internal/config/kpi_weights.go` (atau table di DB) untuk weights (§2.1) dan grade bands (§3) — **no hardcode**.

### Fase 3 — API Layer & Security
9. `apps/api/internal/api/handlers/kpi_handler.go` — request parsing, validasi date range, panggil service.
10. `apps/api/internal/api/routes/kpi_routes.go` — daftarkan `/kpi/sales-rep` dan `/kpi/sales-manager` dengan role middleware (§6).
11. Terapkan authorize guard sebagai middleware terpisah agar reusable.

### Fase 4 — Dokumentasi & Sinkronisasi
12. Update `docs/TECHNICAL_FEATURESv2.md` dengan definisi final metrik & formula.
13. Update `docs/sprint/sprint2/SPRINT_PLANNING_DEV4.md` dengan scope KPI final.
14. Tambahkan request baru ke `docs/postman/CRM-Healthcare-API.postman_collection.json` dengan contoh response persis seperti §5.2/§5.3.

### Fase 5 — Testing
15. Unit test formula (khusus edge case: pembagian nol, data kosong, target null).
16. Unit test attribution fallback (brick null, multi-brick rep).
17. Integration test scope: rep akses rep lain (harus 403), manager akses brick di luar scope (harus 403).
18. Manual test dengan date range sempit (1 hari) dan lebar (1 tahun).

---

