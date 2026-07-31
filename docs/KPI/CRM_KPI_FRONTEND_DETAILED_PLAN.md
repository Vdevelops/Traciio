# CRM KPI Feature — Frontend Detailed Implementation Plan (v2)

> Dokumen ini adalah versi diperdalam dari plan FE awal, dan **align langsung** ke kontrak API di `CRM_KPI_FEATURE_DETAILED_PLAN.md` (§5.2 response Sales Rep, §5.3 response Sales Manager). Semua nama field TypeScript di dokumen ini mengikuti field JSON dari BE plan — jangan ubah casing atau struktur tanpa mengecek ulang ke dokumen BE.

---

## 0. Ringkasan & Prinsip Desain

- FE hanya menyasar **web app** (`apps/web`) — tidak ada `apps/mobile` di workspace saat ini.
- KPI mendapat **route & section sendiri** (`/kpi` atau `/evaluation`), bukan disisipkan ke dashboard summary.
- Reuse **pattern** (layout, table, card, drilldown) dari dashboard & sales-overview yang sudah ada — tapi **jangan reuse shape data lama**, karena shape dashboard tidak dirancang untuk composite score, grading, dan diagnostics.
- Prinsip utama: UI harus **menjelaskan** KPI (grade, trend, actionable insight), bukan cuma menampilkan angka mentah.

---

## 1. Information Architecture & Routing

```
apps/web/app/[locale]/kpi/
  page.client.tsx                     # entry point, role-based redirect/render
  components/
    kpi-role-router.tsx               # baca role dari JWT/session, render view sesuai role
    sales-rep-kpi-view.tsx            # view untuk role sales_rep
    sales-manager-kpi-view.tsx        # view untuk role sales_manager
  [userId]/
    page.client.tsx                   # drilldown detail per sales rep (dipakai manager)
  brick/[brickId]/
    page.client.tsx                   # drilldown detail per brick/area (dipakai manager)
```

**Aturan routing:**
- `scope` (userId, managerId, brickId) **tidak boleh** diambil bebas dari query param tanpa validasi terhadap sesi user yang login. Default selalu `self`; hanya manager yang boleh drilldown ke `userId` anak buahnya sendiri (validasi ini tetap terjadi di BE lewat guard §6 BE plan, tapi FE juga harus tidak menawarkan link ke scope yang tidak valid).
- `sales-rep-kpi-view.tsx` **tidak pernah** menerima `userId` sebagai prop dari luar kecuali dari drilldown manager — untuk rep yang login sebagai dirinya sendiri, `userId` selalu diambil dari session, tidak dari URL, supaya tidak ada state di mana rep bisa mengubah URL untuk melihat data rep lain (401/403 tetap ditangani BE, tapi FE sebaiknya juga tidak expose input bebas).

---

## 2. Kontrak Data Frontend (TypeScript types)

Tambahkan file baru — **jangan extend `index.d.ts` dashboard lama**, karena shape KPI berbeda secara fundamental (ada composite score, grade, diagnostics).

```
apps/web/src/features/kpi/types/index.d.ts
```

```typescript
// ==== Shared ====
export type KPIGrade = "Excellent" | "Good" | "Needs Improvement" | "Critical";
export type TrendDirection = "up" | "down" | "flat";
export type DiagnosticSeverity = "info" | "warning" | "critical";

export interface DiagnosticFlag {
  code: string;
  severity: DiagnosticSeverity;
  message: string;
  brickId?: string;
}

export interface TrendInfo {
  previousCompositeScore: number;
  delta: number;
  direction: TrendDirection;
}

export interface KPIMeta {
  brickMissingCount: number;
  brickInferredCount: number;
  generatedAt: string; // ISO datetime
}

// ==== Sales Rep ====
export interface SalesRepScorecard {
  totalDealsClosed: number;
  totalRevenue: number;
  dealsCreated: number;
  conversionRate: number | null;
  averageDealValue: number | null;
  visitCompleted: number;
  visitPlanned: number;
  visitCompliance: number | null;
  tasksCompleted: number;
  overdueTaskRate: number | null;
  revenueTargetAttainment: number | null;
  dealTargetAttainment: number | null;
  pipelineMovementScore: number;
}

export interface TargetGapItem {
  target: number;
  actual: number;
  gapPercent: number;
  status: "above" | "met" | "below";
}

export interface SalesRepEvaluation {
  compositeScore: number;
  grade: KPIGrade;
  trend: TrendInfo;
  targetGap: {
    revenue: TargetGapItem;
    deals: TargetGapItem;
  };
}

export interface SalesRepKPIResponse {
  scope: { userId: string; role: "sales_rep"; startDate: string; endDate: string };
  scorecard: SalesRepScorecard;
  evaluation: SalesRepEvaluation;
  diagnostics: DiagnosticFlag[];
  meta: KPIMeta;
}

// ==== Sales Manager ====
export interface TeamSummary {
  totalRepsCount: number;
  totalDealsClosed: number;
  totalRevenue: number;
  teamConversionRate: number | null;
  teamVisitCompliance: number | null;
  teamOverdueTaskRate: number | null;
  teamTargetAttainment: number | null;
}

export interface TeamBreakdownItem {
  userId: string;
  name: string;
  compositeScore: number;
  grade: KPIGrade;
  totalRevenue: number;
  conversionRate: number | null;
  rank: number;
}

export interface BrickBreakdownItem {
  brickId: string;
  name: string;
  coveragePenetration: number | null;
  totalRevenue: number;
  repsCount: number;
  compositeScore: number;
}

export interface SalesManagerKPIResponse {
  scope: { managerId: string; role: "sales_manager"; startDate: string; endDate: string; bricks: string[] };
  teamSummary: TeamSummary;
  evaluation: { compositeScore: number; grade: KPIGrade; trend: TrendInfo };
  teamBreakdown: TeamBreakdownItem[];
  brickBreakdown: BrickBreakdownItem[];
  topBottomPerformers: { top: string[]; bottom: string[] };
  diagnostics: DiagnosticFlag[];
  meta: KPIMeta;
}

// ==== Request params (dipakai hook & service) ====
export interface KPIDateRangeParams {
  startDate: string;
  endDate: string;
  compareWithPrevious?: boolean;
}

export interface SalesRepKPIParams extends KPIDateRangeParams {
  userId?: string; // hanya diisi manager saat drilldown
}

export interface SalesManagerKPIParams extends KPIDateRangeParams {
  managerId?: string;
  brickId?: string;
  includeTeamBreakdown?: boolean;
}
```

> **Catatan penting**: semua field numeric yang di BE bisa `null` (conversionRate, averageDealValue, visitCompliance, overdueTaskRate, dst) **wajib** bertipe `number | null` di FE, bukan `number`. Ini untuk memaksa setiap komponen render eksplisit menangani kasus "data tidak cukup untuk dihitung" — jangan render `null` sebagai `0` karena itu menyesatkan (0% conversion vs "belum ada data" adalah makna yang beda total).

---

## 3. Service Layer

**Jangan** taruh KPI request di `dashboardService.ts` atau `salesOverviewService.ts` yang sudah ada — buat service baru khusus supaya tidak mencampur concern, tapi tetap **reuse http client/axios instance** yang sama dengan service lain.

```
apps/web/src/features/kpi/services/kpiService.ts
```

```typescript
import { httpClient } from "@/lib/http-client"; // sesuaikan dengan instance yang dipakai dashboardService.ts
import type {
  SalesRepKPIResponse, SalesManagerKPIResponse,
  SalesRepKPIParams, SalesManagerKPIParams,
} from "../types";

export const kpiService = {
  getSalesRepKPI: (params: SalesRepKPIParams) =>
    httpClient.get<SalesRepKPIResponse>("/api/v1/kpi/sales-rep", { params }),

  getSalesManagerKPI: (params: SalesManagerKPIParams) =>
    httpClient.get<SalesManagerKPIResponse>("/api/v1/kpi/sales-manager", { params }),
};
```

> **Verifikasi wajib sebelum implementasi**: cek nama instance http client aktual yang dipakai `dashboardService.ts` (misal `apiClient`, `axiosInstance`, dll) — jangan asumsikan nama `httpClient` di atas, itu hanya placeholder.

---

## 4. React Query Hooks

```
apps/web/src/features/kpi/hooks/useSalesRepKPI.ts
apps/web/src/features/kpi/hooks/useSalesManagerKPI.ts
apps/web/src/features/kpi/hooks/useKPIDrilldown.ts   # untuk userId/brickId drilldown
```

```typescript
// useSalesRepKPI.ts
export const kpiQueryKeys = {
  salesRep: (params: SalesRepKPIParams) => ["kpi", "sales-rep", params] as const,
  salesManager: (params: SalesManagerKPIParams) => ["kpi", "sales-manager", params] as const,
};

export function useSalesRepKPI(params: SalesRepKPIParams) {
  return useQuery({
    queryKey: kpiQueryKeys.salesRep(params),
    queryFn: () => kpiService.getSalesRepKPI(params),
    enabled: Boolean(params.startDate && params.endDate),
    staleTime: 5 * 60 * 1000, // 5 menit — KPI tidak perlu realtime
  });
}
```

**Aturan query key**: query key **wajib** menyertakan seluruh params (termasuk `userId`/`brickId`/date range) sebagai bagian dari key — supaya React Query tidak menampilkan data cache dari scope lain saat user pindah filter/drilldown. Ini menjawab poin verification "filter tanggal, period preset, dan drilldown tidak merusak state React Query atau menyebabkan data campur antar scope" di plan awal.

---

## 5. Struktur Halaman — 3 Blok Utama

### 5.1 Sales Rep View (`sales-rep-kpi-view.tsx`)

```
Blok 1 — Ringkasan
  - CompositeScoreCard (grade badge, angka besar, warna sesuai grade)
  - TrendIndicator (panah naik/turun + delta vs periode lalu)
  - TargetGapSummary (revenue & deals: target vs actual vs gap%)

Blok 2 — Scorecard & Trend Detail
  - MetricGrid (semua field dari SalesRepScorecard, masing-masing sebagai MetricCard)
  - Setiap MetricCard: raw value + label + (opsional) mini sparkline jika BE nanti sediakan trend series
  - Handle null: tampilkan "Belum ada data" bukan "0" atau "-"

Blok 3 — Diagnostic & Actionable Insight
  - DiagnosticList: render dari `diagnostics[]`, grouped by severity (critical dulu, lalu warning, lalu info)
  - Setiap diagnostic item: icon sesuai severity + message (message sudah final dari BE, FE tidak perlu compose ulang teks)
```

Filter di atas semua blok: `DateRangePicker` + `PeriodPresetSelector` (This Month, Last Month, This Quarter, Custom).

### 5.2 Sales Manager View (`sales-manager-kpi-view.tsx`)

```
Blok 1 — Ringkasan Tim
  - CompositeScoreCard (team-level, sama komponen dengan rep tapi data dari `evaluation`)
  - TeamSummaryGrid (dari `teamSummary`)
  - TrendIndicator (sama pattern)

Blok 2 — Breakdown
  - Tab/Toggle: "Per Sales Rep" vs "Per Brick/Area"
  - Per Sales Rep -> TeamBreakdownTable (reuse pola SalesPerformanceList.tsx untuk table+ranking,
    tapi kolom disesuaikan: rank, name, compositeScore, grade badge, totalRevenue, conversionRate)
    -> klik row -> navigate ke /kpi/[userId] drilldown
  - Per Brick/Area -> BrickBreakdownTable (coveragePenetration, totalRevenue, repsCount, compositeScore)
    -> klik row -> navigate ke /kpi/brick/[brickId] drilldown
  - TopBottomPerformersPanel (reuse pola SalesPodium.tsx untuk top performer visual, tambah panel
    kecil terpisah untuk bottom performer dengan styling "perlu perhatian", bukan styling negatif/hukuman)

Blok 3 — Diagnostic & Actionable Insight
  - Sama seperti rep, tapi diagnostic bisa punya `brickId` -> render badge brick di item terkait
```

### 5.3 Drilldown Pages

- `/kpi/[userId]` — reuse struktur `sales-rep-kpi-view.tsx` apa adanya, tapi header menampilkan "Melihat KPI: {nama rep}" + tombol "Kembali ke Tim". Data fetch pakai `userId` dari path param (bukan dari session), hanya accessible kalau requester adalah manager dari rep tsb (BE yang final-validate, FE cukup sembunyikan link kalau bukan anak buah).
- `/kpi/brick/[brickId]` — halaman baru, fokus ke satu brick: `BrickBreakdownItem` detail + daftar rep di brick tsb + diagnostic yang brickId-nya cocok.

---

## 6. Reuse Pattern dari File Existing (referensi, bukan copy shape data)

| File existing | Yang di-reuse | Yang TIDAK di-reuse |
|---|---|---|
| `page.client.tsx` (dashboard) | pola routing role-based, shell layout | data fetching shape dashboard |
| `role-based-dashboard.tsx` | pola komposisi render-per-role | — |
| `sales-manager-dashboard.tsx` | pola layout manager (grid, section) | binding ke `dashboardService` |
| `sales-dashboard.tsx` | pola layout rep personal | binding ke `dashboardService` |
| `SalesPerformanceList.tsx` | struktur tabel ranking + sorting | field dan sort key (ganti ke composite score) |
| `SalesPodium.tsx` | visual top performer (podium style) | data source (ganti ke `teamBreakdown` top 3) |
| `sales-rep-detail-page-client.tsx` | pola drilldown page (header, back button, layout detail) | data hook (ganti ke `useSalesRepKPI` dengan `userId` dari path) |

---

## 7. State: Loading, Empty, Error, Partial Data

```
LoadingState:
  - Skeleton per blok (bukan satu skeleton besar full page) supaya blok yang sudah selesai
    fetch bisa langsung render duluan jika di masa depan dipecah jadi beberapa request.

EmptyState:
  - Trigger: scorecard semua 0 DAN dealsCreated = 0 DAN visitPlanned = 0 dalam range tsb.
  - Message: "Belum ada aktivitas tercatat pada periode ini." + suggestion ganti date range.
  - JANGAN tampilkan grade/composite score kalau kondisi ini true (grade dari 0 aktivitas
    menyesatkan) — tampilkan state "Tidak cukup data untuk evaluasi" alih-alih grade badge.

ErrorState:
  - 403 -> "Anda tidak memiliki akses ke data ini." + redirect ke KPI milik sendiri.
  - 500/network -> generic error card + retry button (React Query refetch).

PartialDataFallback:
  - Field individual null (conversionRate, visitCompliance, dst) -> render "N/A" atau
    "Belum ada data" di level MetricCard, JANGAN sembunyikan seluruh blok karena satu
    field null.
  - meta.brickMissingCount > 0 -> tampilkan info badge kecil "X data belum terpetakan
    ke brick" di dekat brickBreakdown (transparansi data quality ke manager).
```

---

## 8. Formatting & Copy Standard

Buat util terpusat, jangan format ad-hoc di tiap komponen:

```
apps/web/src/features/kpi/utils/formatters.ts

- formatCurrency(value: number | null): string       // "Rp 450.000.000" atau "Belum ada data"
- formatPercent(value: number | null, decimals=1)     // "30.0%" atau "N/A"
- formatCompositeScore(value: number): string          // "83.5"
- gradeColor(grade: KPIGrade): string                  // mapping ke design token warna
- trendIcon(direction: TrendDirection): IconComponent  // panah atas/bawah/datar
- diagnosticSeverityIcon(severity: DiagnosticSeverity)
```

**Copy guideline** (bukan cuma angka, harus menjelaskan):
- Composite score selalu didampingi grade label, jangan tampilkan angka telanjang.
- Trend selalu didampingi teks pembanding: "naik 5.5 poin dari bulan lalu", bukan cuma "+5.5".
- Target gap selalu dalam kalimat: "Revenue 12.5% di atas target" bukan cuma angka gap.

---

## 9. Navigasi & RBAC

```
FUNCTION renderKPIMenuItem(currentUser):
    IF currentUser.role IN ["sales_rep", "sales_manager"]:
        SHOW menu "Evaluasi KPI" di sidebar/nav
    ELSE:
        HIDE menu item

FUNCTION kpiRoleRouter(currentUser):
    IF currentUser.role == "sales_rep": RENDER SalesRepKPIView (userId = currentUser.id)
    IF currentUser.role == "sales_manager": RENDER SalesManagerKPIView (managerId = currentUser.id)
    ELSE: RENDER 403/NotAuthorized
```

- Role dibaca dari session/JWT yang sudah divalidasi di layer auth existing (cek pola yang dipakai di `role-based-dashboard.tsx`) — jangan buat mekanisme role-check baru yang terpisah.

---

## 10. Rencana Implementasi Frontend (bertahap)

### Fase 1 — Fondasi Data Layer
1. `features/kpi/types/index.d.ts` — semua type di §2.
2. `features/kpi/services/kpiService.ts` — verifikasi dulu http client instance aktual.
3. `features/kpi/hooks/useSalesRepKPI.ts`, `useSalesManagerKPI.ts` — query key wajib include semua params.
4. `features/kpi/utils/formatters.ts` — util formatting §8.

### Fase 2 — Sales Rep View
5. `app/[locale]/kpi/components/sales-rep-kpi-view.tsx` — 3 blok sesuai §5.1.
6. Komponen kecil: `CompositeScoreCard`, `TrendIndicator`, `TargetGapSummary`, `MetricCard`, `DiagnosticList`.
7. State handling: loading/empty/error/partial sesuai §7.

### Fase 3 — Sales Manager View
8. `app/[locale]/kpi/components/sales-manager-kpi-view.tsx` — 3 blok sesuai §5.2.
9. `TeamBreakdownTable` (reuse pola `SalesPerformanceList.tsx`), `BrickBreakdownTable`, `TopBottomPerformersPanel` (reuse `SalesPodium.tsx`).

### Fase 4 — Drilldown & Routing
10. `app/[locale]/kpi/page.client.tsx` + `kpi-role-router.tsx`.
11. `app/[locale]/kpi/[userId]/page.client.tsx`.
12. `app/[locale]/kpi/brick/[brickId]/page.client.tsx`.
13. Nav menu entry (§9), role guard.

### Fase 5 — Polish & Sinkronisasi Dokumen
14. Filter global: `DateRangePicker`, `PeriodPresetSelector` — pastikan konsisten dengan pattern filter yang sudah dipakai di sales-overview.
15. Update dokumentasi FE (struktur feature baru) dan pastikan kontrak response yang dipakai FE tetap sinkron dengan Postman collection BE.

### Fase 6 — Testing
16. Typecheck + lint pada semua file baru.
17. Manual smoke test 3 skenario: personal scope (rep), team scope (manager), brick scope (drilldown).
18. Test edge case: date range menghasilkan empty state, field null render dengan benar, 403 saat drilldown ke luar scope (simulasi via mock).
19. Test query key isolation: pindah filter/drilldown tidak menampilkan data cache dari scope sebelumnya (bisa dites manual dengan React Query devtools).

---

