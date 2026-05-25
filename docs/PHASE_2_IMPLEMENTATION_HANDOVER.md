# Phase 2 Implementation Handover
## CRM Healthcare / Pharmaceutical Platform

Date: 2026-04-24
Owner: Dev3 (Backend + Mobile integration track)
Status: In Progress

---

## 1. Purpose

Dokumen ini dibuat sebagai panduan fase berikutnya agar tim dapat:
- memahami apa yang sudah diimplementasikan pada Phase 2,
- memahami bagaimana alur sistem berjalan pada area kritikal,
- menjalankan verifikasi teknis secara konsisten,
- melanjutkan implementasi tanpa kehilangan konteks.

Dokumen baseline audit Phase 1 tersedia di [docs/PHASE_1_AUDIT_BASELINE.md](docs/PHASE_1_AUDIT_BASELINE.md).

---

## 2. Scope Phase 2

Target Phase 2 (10 core modules):
1. Activities
2. Leads
3. Visit Reports
4. Pipeline and Deals
5. Contacts
6. Tasks
7. Dashboard
8. Reports
9. Master Data
10. Notifications

Catatan:
- Delivery mode disepakati gabungan backend + web.
- Implementasi saat ini masih berfokus pada backend hardening untuk business rules dan security consistency.

---

## 3. Current Implementation Status

### Completed in this iteration

1. Activity scope enforcement
- Request DTO sudah mendukung scoped user IDs.
- Query list dan timeline activity sudah mengikuti scope.
- Route activity sudah menggunakan scope middleware.

Main files:
- [apps/api/internal/domain/activity/entity.go](apps/api/internal/domain/activity/entity.go)
- [apps/api/internal/repository/postgres/activity/repository.go](apps/api/internal/repository/postgres/activity/repository.go)
- [apps/api/internal/api/handlers/activity_handler.go](apps/api/internal/api/handlers/activity_handler.go)
- [apps/api/internal/api/routes/activity_routes.go](apps/api/internal/api/routes/activity_routes.go)
- [apps/api/cmd/server/main.go](apps/api/cmd/server/main.go)

2. Notification ownership enforcement
- Mark as read dan delete notification sekarang wajib pemilik notification.
- Handler sudah map forbidden dan not-found secara eksplisit.

Main files:
- [apps/api/internal/service/notification/service.go](apps/api/internal/service/notification/service.go)
- [apps/api/internal/api/handlers/notification_handler.go](apps/api/internal/api/handlers/notification_handler.go)
- [apps/api/internal/service/notification/notification_service_test.go](apps/api/internal/service/notification/notification_service_test.go)
- [apps/api/internal/service/notification/integration_test.go](apps/api/internal/service/notification/integration_test.go)

3. Pipeline and Deals move consistency
- Endpoint PATCH move dipusatkan ke jalur validasi stage yang sama dengan move-stage.
- Error mapping deal/stage not found dinormalisasi.
- Inconsistent debug print di deals list dihapus.

Main files:
- [apps/api/internal/api/handlers/deal_handler.go](apps/api/internal/api/handlers/deal_handler.go)
- [apps/api/internal/service/pipeline/service.go](apps/api/internal/service/pipeline/service.go)

4. Visit report submit error normalization
- Submit menggunakan sentinel error untuk not owner dan prerequisite (check-in and check-out).
- Handler web/mobile submit tidak lagi bergantung string matching.

Main files:
- [apps/api/internal/service/visit_report/service.go](apps/api/internal/service/visit_report/service.go)
- [apps/api/internal/api/handlers/visit_report_handler.go](apps/api/internal/api/handlers/visit_report_handler.go)

5. Task auto-reminder D-1 and day-H
- Auto reminder in-app dibuat otomatis saat create dan update task jika due date tersedia.
- Mekanisme dibuat idempotent untuk mencegah reminder auto dobel.

Main file:
- [apps/api/internal/service/task/service.go](apps/api/internal/service/task/service.go)

---

## 4. How the System Works (Critical Flows)

### 4.1 Deal stage transition flow

Entry points:
- PATCH /api/v1/deals/:id/move
- POST /api/v1/deals/:id/move-stage

Runtime flow:
1. Handler melakukan auth dan permission gate pipeline.update_stage.
2. Service menjalankan ValidateStageRequirements.
3. Service memindahkan stage, update probability dan status won/lost/open.
4. Service membuat deal history.
5. Service membuat task transisi stage bila relevan.
6. Jika won, service memicu konversi purchase history.

Source:
- [apps/api/internal/api/handlers/deal_handler.go](apps/api/internal/api/handlers/deal_handler.go)
- [apps/api/internal/api/handlers/pipeline_handler.go](apps/api/internal/api/handlers/pipeline_handler.go)
- [apps/api/internal/service/pipeline/service.go](apps/api/internal/service/pipeline/service.go)

### 4.2 Visit report submit flow

Entry points:
- PATCH /api/v1/visit-reports/:id/submit
- PATCH /api/v1/mobile/visit-reports/:id/submit

Runtime flow:
1. Handler mengambil user_id dari auth context.
2. Service validasi ownership (sales rep owner).
3. Service validasi state transition hanya draft to submitted.
4. Service validasi prerequisite check-in dan check-out.
5. Service menjalankan auto trigger: lead update, task auto-create, manager notification, activity log.

Source:
- [apps/api/internal/api/handlers/visit_report_handler.go](apps/api/internal/api/handlers/visit_report_handler.go)
- [apps/api/internal/service/visit_report/service.go](apps/api/internal/service/visit_report/service.go)

### 4.3 Task auto reminder flow

Entry points:
- Task create
- Task update

Runtime flow:
1. Saat task punya due date, service sinkronkan reminder otomatis.
2. Reminder auto yang lama (prefix sistem) dihapus.
3. Service membuat dua reminder baru:
- D-1 dari due date
- Day-H pada due date
4. Reminder yang timestamp-nya sudah lewat tidak dibuat.

Source:
- [apps/api/internal/service/task/service.go](apps/api/internal/service/task/service.go)
- [apps/api/internal/worker/reminder_worker.go](apps/api/internal/worker/reminder_worker.go)

### 4.4 Notification ownership flow

Entry points:
- PUT /api/v1/notifications/:id/read
- DELETE /api/v1/notifications/:id

Runtime flow:
1. Handler ambil user_id terautentikasi.
2. Service load notification by ID.
3. Service bandingkan notification.user_id dengan requester user_id.
4. Jika mismatch, return forbidden.
5. Jika cocok, lanjut mark as read atau delete.

Source:
- [apps/api/internal/api/handlers/notification_handler.go](apps/api/internal/api/handlers/notification_handler.go)
- [apps/api/internal/service/notification/service.go](apps/api/internal/service/notification/service.go)

### 4.5 Activity scoped access

Entry points:
- GET /api/v1/activities
- GET /api/v1/activities/timeline

Runtime flow:
1. Scope middleware membentuk user context dan scoped user IDs.
2. Handler inject scoped user IDs ke request.
3. Repository apply filtering user_id IN scoped user IDs jika filter user explicit tidak diberikan.

Source:
- [apps/api/internal/api/middleware/scope.go](apps/api/internal/api/middleware/scope.go)
- [apps/api/internal/api/handlers/activity_handler.go](apps/api/internal/api/handlers/activity_handler.go)
- [apps/api/internal/repository/postgres/activity/repository.go](apps/api/internal/repository/postgres/activity/repository.go)

---

## 5. Verification Guide (Phase 2)

### 5.1 Backend checks

From repository root:

- Push-Location apps/api; go test ./...; Pop-Location

Focused checks used in this iteration:
- Push-Location apps/api; go test ./internal/service/notification ./internal/api/handlers ./internal/api/routes ./cmd/server; Pop-Location
- Push-Location apps/api; go test ./internal/service/pipeline ./internal/api/handlers ./internal/api/routes ./cmd/server; Pop-Location
- Push-Location apps/api; go test ./internal/service/visit_report ./internal/api/handlers ./cmd/server; Pop-Location
- Push-Location apps/api; go test ./internal/service/task ./internal/api/handlers ./cmd/server; Pop-Location

### 5.2 Web checks

- Push-Location apps/web; npx pnpm lint; Pop-Location
- Push-Location apps/web; npx pnpm check-types; Pop-Location
- Push-Location apps/web; npx pnpm build; Pop-Location

### 5.3 API contract checks

Run Postman collection:
- [docs/postman/CRM-Healthcare-API.postman_collection.json](docs/postman/CRM-Healthcare-API.postman_collection.json)

Validate:
- success envelope consistency,
- error envelope consistency,
- pagination meta consistency,
- proper 403/404/422 style mapping for business constraints.

---

## 6. Remaining Work for Next Phase

### Priority 1
1. Master Data backend route group masih placeholder dan perlu final route contract yang eksplisit.
2. Reports contract hardening (error and filter consistency), terutama untuk web integration yang stabil.
3. Dashboard role-aware data assertions (security and scope correctness in aggregate endpoints).

### Priority 2
1. Tambah unit test untuk task auto-reminder logic (service-level behavior tests).
2. Tambah regression test untuk dual move endpoints agar behavior tetap konsisten.
3. Tambah E2E flow assertions untuk visit report submit and approval chain.

### Priority 3
1. Web integration sweep terhadap endpoint yang baru dinormalisasi.
2. Update sprint doc dan postman examples jika ada delta contract.

---

## 7. Operational Notes

1. Database schema di project ini hybrid:
- AutoMigrate membuat base table,
- SQL migration digunakan untuk perubahan schema lanjutan.

Source:
- [apps/api/internal/database/database.go](apps/api/internal/database/database.go)

2. Jika melakukan validasi terminal dari root monorepo, selalu pindah ke apps/api untuk command Go module.

3. Jangan asumsikan endpoint stage move lama dan baru bisa divergen; keduanya harus dipertahankan konsisten sampai ada keputusan deprecate.

---

## 8. Definition of Done for Phase 2 Completion

Phase 2 dinyatakan complete jika:
1. 10 modul target mencapai behavior coverage minimal 80 persen sesuai contract internal.
2. Semua flow kritikal pada Section 4 lulus verifikasi.
3. Tidak ada fallback 500 untuk rule validation yang seharusnya 4xx.
4. Security checks untuk ownership, scope, dan permission lulus.
5. Web lint, type-check, build, dan smoke flow utama lulus.

---

End of document.
