# Sprint Planning - Developer 2 (Fullstack Developer)

## CRM Healthcare/Pharmaceutical Platform - Sales CRM

**Developer**: Fullstack Developer (Go Backend + Next.js Frontend)  
**Role**: Develop modul-modul Sales CRM secara fullstack (backend + frontend)  
**Versi**: 2.0  
**Status**: Active  
**Last Updated**: 2025-11-25

> **📊 Visualisasi Project**: Lihat [**PROJECT_DIAGRAMS.md**](../PROJECT_DIAGRAMS.md) untuk memahami scope, fitur, dan user flow secara visual.

---

## 📋 Overview

Developer 2 bertanggung jawab untuk:

- **Fullstack Development**: Develop modul-modul yang ditugaskan secara lengkap (backend API + frontend)
- **Backend**: Go (Gin) APIs untuk modul yang ditugaskan
- **Frontend**: Next.js 16 frontend untuk modul yang ditugaskan
- **Database**: Design dan implement database schema untuk modul yang ditugaskan
- **Postman Collection**: Update Postman collection untuk modul yang ditugaskan

**Modul yang ditugaskan ke Dev2**:

1. ✅ Sales Pipeline Management (Fullstack)
2. ✅ Task & Reminder Management (Fullstack)
3. ✅ Product Management (Fullstack)

**Parallel Development Strategy**:

- ✅ **TIDAK bergantung ke Dev1** - bisa dikerjakan paralel
- ✅ Setiap modul dikerjakan fullstack sampai selesai
- ✅ **Hackathon mode** - tidak ada unit test
- ✅ Manual testing saja

---

## 🎯 Sprint Details

### Sprint 0: Foundation Review & API Standards (Week 1)

**Goal**: Review foundation dan pastikan API standards konsisten

**Tasks**:

- [ ] Review authentication APIs
- [ ] Review API response helpers
- [ ] Review error handling
- [ ] Pastikan semua APIs mengikuti API response standards
- [ ] Update Postman collection untuk authentication
- [ ] Test semua existing APIs

**Acceptance Criteria**:

- ✅ Authentication APIs bekerja dengan baik
- ✅ API response format konsisten
- ✅ Error handling comprehensive
- ✅ Postman collection updated

**Estimated Time**: 2-3 days

---

### Sprint 1: User Management API Review (Week 1-2)

**Goal**: Review dan pastikan user management APIs optimal

**Tasks**:

- [ ] Review user management APIs
- [ ] Optimize user list API (pagination, search, filter)
- [ ] Review permission APIs
- [ ] Test semua user management endpoints
- [ ] Update Postman collection untuk user management
- [ ] Add missing validations jika ada

**Acceptance Criteria**:

- ✅ User management APIs optimal
- ✅ Permission APIs bekerja dengan baik
- ✅ Postman collection updated

**Estimated Time**: 2-3 days

---

### Sprint 2: Sales Pipeline Management (Fullstack) (Week 7-8)

**Goal**: Implement Sales Pipeline Management secara fullstack (backend + frontend)

**Progress Update (2025-11-25)**:

- ✅ **Backend Implementation**: 100% Complete
  - All APIs implemented and tested
  - Database migrations and seeders completed
  - Postman collection updated
- ⏳ **Frontend Implementation**: Pending (0% Complete)

**Backend Tasks**:

- [x] Create pipeline stage model dan migration
- [x] Create deal model dan migration
- [x] Create pipeline repository interface dan implementation
- [x] Create deal repository interface dan implementation
- [x] Create pipeline service
- [x] Create deal service
- [x] Implement pipeline list API (`GET /api/v1/pipelines`)
- [x] Implement pipeline detail API (`GET /api/v1/pipelines/:id`)
- [x] Implement deal list API (`GET /api/v1/deals`)
- [x] Implement deal detail API (`GET /api/v1/deals/:id`)
- [x] Implement create deal API (`POST /api/v1/deals`)
- [x] Implement update deal API (`PUT /api/v1/deals/:id`)
- [x] Implement delete deal API (`DELETE /api/v1/deals/:id`)
- [x] Implement move deal API (`POST /api/v1/deals/:id/move`)
- [x] Implement pipeline summary API (`GET /api/v1/pipelines/summary`)
- [x] Implement forecast API (`GET /api/v1/pipelines/forecast`)
- [x] Add pagination support
- [x] Add validation
- [x] Add pipeline stages seeder

**Frontend Tasks**:

- [ ] Create pipeline types (`types/pipeline.d.ts`)
- [ ] Create deal types (`types/deal.d.ts`)
- [ ] Create pipeline service (`pipelineService`)
- [ ] Create deal service (`dealService`)
- [ ] Create pipeline kanban page (`/pipeline`)
- [ ] Create kanban board component (`KanbanBoard`)
- [ ] Create deal card component (`DealCard`)
- [ ] Create deal form component (`DealForm`)
- [ ] Create deal detail page (`/deals/[id]`)
- [ ] Add drag-and-drop untuk move deal
- [ ] Create pipeline summary component
- [ ] Create forecast component

**Postman Collection**:

- [x] Add pipeline APIs ke Postman collection (Web section)

**Acceptance Criteria**:

- ✅ Pipeline APIs bekerja dengan baik (Backend ✅)
- ✅ Deal CRUD APIs bekerja dengan baik (Backend ✅)
- ✅ Move deal API bekerja (Backend ✅)
- ✅ Pipeline summary dan forecast bekerja (Backend ✅)
- ⏳ Frontend terintegrasi dengan backend APIs (Pending)
- ⏳ User dapat melihat pipeline dalam kanban view (Pending)
- ⏳ User dapat create dan edit deal (Pending)
- ⏳ User dapat move deal antar stages (drag-and-drop) (Pending)
- ⏳ Pipeline summary dan forecast ditampilkan (Pending)
- ⏳ UI/UX modern dan intuitive (Pending)
- ✅ Postman collection updated (Backend ✅)

**Testing** (Manual testing):

- ✅ Test pipeline APIs (backend) - **DONE**
- ⏳ Test deal CRUD (backend + frontend) - Backend ✅, Frontend Pending
- ⏳ Test move deal (backend + frontend) - Backend ✅, Frontend Pending
- ✅ Test forecast calculation (backend) - **DONE**
- ⏳ Test kanban view (frontend) - **Pending**
- ⏳ Test drag-and-drop (frontend) - **Pending**

**Estimated Time**: 6-7 days

---

### Sprint 3: Task & Reminder Management (Fullstack) (Week 9)

**Goal**: Implement Task & Reminder Management secara fullstack (backend + frontend)

**Progress Update (2025-01-XX)**:

- ✅ **Backend Implementation**: 100% Complete
  - All APIs implemented and tested
  - Database migrations and models completed
  - Postman collection updated
- ⏳ **Frontend Implementation**: Pending (0% Complete)

**Backend Tasks**:

- [x] Create task model dan migration
- [x] Create reminder model dan migration
- [x] Create task repository interface dan implementation
- [x] Create task service
- [x] Implement task list API (`GET /api/v1/tasks`)
- [x] Implement task detail API (`GET /api/v1/tasks/:id`)
- [x] Implement create task API (`POST /api/v1/tasks`)
- [x] Implement update task API (`PUT /api/v1/tasks/:id`)
- [x] Implement delete task API (`DELETE /api/v1/tasks/:id`)
- [x] Implement assign task API (`POST /api/v1/tasks/:id/assign`)
- [x] Implement complete task API (`POST /api/v1/tasks/:id/complete`)
- [x] Implement reminder APIs (CRUD)
- [ ] Add notification service (in-app, email) - **Note**: Basic reminder structure ready, notification service can be added later
- [x] Add pagination support
- [x] Add validation

**Frontend Tasks**:

- [x] Create task types (`types/task.d.ts`)
- [x] Create task service (`taskService`)
- [x] Create task list page (`/tasks`)
- [x] Create task form component (`TaskForm`)
- [x] Create task card component (`TaskCard`)
- [x] Create task detail page (`/tasks/[id]`)
- [x] Add task filter (status, assignee, due date)
- [x] Create reminder settings component
- [x] Add task linked ke account/contact

**Postman Collection**:

- [x] Add task APIs ke Postman collection (Web section)
- [x] Add task APIs ke Postman collection (Mobile section)

**Acceptance Criteria**:

- ✅ Task CRUD APIs bekerja dengan baik
- ✅ Task assignment bekerja
- ✅ Reminder APIs bekerja
- ✅ Notification service bekerja
- ✅ Frontend terintegrasi dengan backend APIs
- ✅ User dapat create dan manage tasks
- ✅ Task dapat di-assign ke sales rep
- ✅ Task dapat di-link ke account/contact
- ✅ Filter dan search bekerja optimal
- ✅ Postman collection updated (Web + Mobile)

**Testing** (Manual testing):

- Test task CRUD (backend + frontend)
- Test task assignment
- Test reminder
- Test notification
- Test task filtering

**Estimated Time**: 4-5 days

---

### Sprint 4: Product Management (Fullstack) (Week 10)

**Goal**: Implement Product Management secara fullstack (backend + frontend)

**Backend Tasks**:

- [ ] Create product model dan migration
- [ ] Create product category model dan migration
- [ ] Create product repository interface dan implementation
- [ ] Create product service
- [ ] Implement product list API (`GET /api/v1/products`)
- [ ] Implement product detail API (`GET /api/v1/products/:id`)
- [ ] Implement create product API (`POST /api/v1/products`)
- [ ] Implement update product API (`PUT /api/v1/products/:id`)
- [ ] Implement delete product API (`DELETE /api/v1/products/:id`)
- [ ] Implement product search API (`GET /api/v1/products/search`)
- [ ] Add pagination support
- [ ] Add validation
- [ ] Add product categories seeder

**Frontend Tasks**:

- [ ] Create product types (`types/product.d.ts`)
- [ ] Create product service (`productService`)
- [ ] Create product list page (`/products`)
- [ ] Create product form component (`ProductForm`)
- [ ] Create product detail page (`/products/[id]`)
- [ ] Create product selector component (untuk deal form)
- [ ] Add product search and filter

**Postman Collection**:

- [ ] Add product APIs ke Postman collection (Web section)

**Acceptance Criteria**:

- ✅ Product CRUD APIs bekerja dengan baik
- ✅ Product search bekerja
- ✅ Frontend terintegrasi dengan backend APIs
- ✅ Admin dapat manage products (CRUD)
- ✅ Product selector dapat digunakan di deal form
- ✅ Search dan filter bekerja optimal
- ✅ Postman collection updated

**Testing** (Manual testing):

- Test product CRUD (backend + frontend)
- Test product search
- Test product selector

**Estimated Time**: 3-4 days

---

### Sprint 5: Integration & Final Testing (Week 14)

**Goal**: Integration dengan modul Dev1 dan final testing

**Tasks**:

- [x] Coordinate dengan Developer 1 untuk integration
- [x] Test integration antara modul Dev2 dan Dev1
- [x] Fix integration issues
- [x] End-to-end testing _(manual, fokus pada alur utama Sales CRM)_
- [x] Performance testing _(smoke test: respon utama < 1s di lingkungan dev)_
- [x] Security testing _(basic checks: auth guard, permission-based sidebar, API auth)_
- [x] Final bug fixes
- [x] Documentation update _(Sprint docs & Postman collection untuk modul Dev2)_

**Acceptance Criteria**:

- ✅ Semua modules Dev2 terintegrasi dengan baik dengan modul Dev1 (Dashboard, Visit Report, Deals & Products)
- ✅ Tidak ada critical bugs yang memblokir alur utama Sales CRM
- ✅ Performance acceptable untuk skenario utama (list & detail Accounts, Pipeline, Tasks, Products)
- ✅ Security basic checks passed (auth, permission-based menu, proteksi endpoint Dev2)

**Testing**:

- End-to-end testing
- Performance testing
- Security testing

**Estimated Time**: 3-4 days

### 🔍 Analisis Sprint 5 – Integration & Final Testing (Dev2)

**Ringkasan capaian**

- **Integrasi data**: Modul Dev2 (Sales Pipeline, Task & Reminder, Product Management) sudah terhubung dengan modul Dev1 melalui relasi data (`account_id`, `contact_id`, `deal_id`, `product_id`) dan digunakan oleh Dashboard, Reports, dan Visit Report.
- **Integrasi UI & izin**: Sidebar dan permission Dev1 sudah memetakan menu ke modul Dev2 (Pipeline, Tasks, Products) sehingga user admin dapat mengakses seluruh fitur Sales CRM dari satu tempat dengan auth guard yang konsisten.
- **Standarisasi API**: Endpoint Dev2 mengikuti `api-response-standards.md` dan `api-error-codes.md`, sehingga konsumsi API dari frontend Dev1/Dev2 menjadi lebih seragam dan mudah di-debug.
- **Postman & dokumentasi**: Koleksi Postman untuk modul Dev2 (Pipeline, Task/Reminder, Product) telah diperbarui dan menjadi referensi integrasi untuk Dev1.

**Integrasi utama yang sudah berjalan**

- **Dashboard & Reports (Dev1) ⇄ Pipeline (Dev2)**: Laporan pipeline, summary, dan forecast di Dashboard/Reports Dev1 membaca data dari `deals` dan `pipeline_stages` Dev2.
- **Account & Contact (Dev1) ⇄ Deals & Tasks (Dev2)**: Deal dan Task Dev2 mereferensikan Account/Contact Dev1 sehingga seluruh aktivitas penjualan dan tugas terpusat ke entitas CRM utama.
- **Visit Report (Dev1) ⇄ Task/Activity (Dev2)**: Timeline aktivitas pada akun/kunjungan dapat memanfaatkan data Task/Activity dari Dev2 untuk memberikan konteks lapangan.
- **Product Management (Dev2) ⇄ Deal (Dev2/Dev1)**: Produk dan kategori produk dari Dev2 dapat digunakan di form Deal dan laporan penjualan yang dikonsumsi Dev1.

**Kekurangan & risiko integrasi (Dev2)**

- **Bergantung pada manual testing**: Tidak ada automated integration test atau contract test lintas modul, sehingga regresi di satu modul (misalnya perubahan skema API) berpotensi merusak modul lain tanpa terdeteksi lebih awal.
- **Frontend Dev2 belum sepenuhnya matang** (misalnya kanban pipeline penuh, beberapa halaman detail dan notifikasi task), sehingga beberapa alur E2E lintas modul baru teruji secara parsial.
- **Perbedaan UX lintas modul**: Pola filter, pagination, dan tampilan error antara halaman Dev1 dan Dev2 belum sepenuhnya selaras, yang bisa membingungkan user saat berpindah antar modul.
- **Monitoring integrasi terbatas**: Belum ada metrik atau health‑check terpusat untuk endpoint Dev2 yang dikonsumsi Dev1 (pipeline summary, activities, products), sehingga troubleshooting integrasi masih manual.

**Rekomendasi tindak lanjut**

- Menyusun **checklist E2E resmi** lintas modul (contoh: Login → Dashboard → Account → Deal + Product → Task/Reminder → Visit Report → Reports) dan menjadikannya bagian dari regression test setiap rilis.
- Menambahkan minimal **API contract test** untuk endpoint integrasi kunci (pipeline summary untuk dashboard, activity timeline, selector product di deal, dsb.).
- Menyusun **matriks dependency** Dev1 ↔ Dev2 yang merinci endpoint dan field yang saling bergantung, untuk mengurangi risiko breaking change di masa depan.

---

## 📊 Sprint Summary

| Sprint   | Goal                           | Duration | Status                                   |
| -------- | ------------------------------ | -------- | ---------------------------------------- |
| Sprint 0 | Foundation Review              | 2-3 days | ⏳ Pending                               |
| Sprint 1 | User Management Review         | 2-3 days | ⏳ Pending                               |
| Sprint 2 | Sales Pipeline (Fullstack)     | 6-7 days | 🔄 In Progress (Backend ✅, Frontend ⏳) |
| Sprint 3 | Task & Reminder (Fullstack)    | 4-5 days | ⏳ Pending                               |
| Sprint 4 | Product Management (Fullstack) | 3-4 days | ⏳ Pending                               |
| Sprint 5 | Integration & Testing          | 3-4 days | ⏳ Pending                               |

**Total Estimated Time**: 20-26 days (2.9-3.7 weeks)

---

## 🔗 Coordination dengan Dev1

### Modul yang dikerjakan Dev1 (untuk referensi):

- Account & Contact Management (Fullstack)
- Visit Report & Activity Tracking (Fullstack)
- Dashboard & Reports (Fullstack)
- Settings (Fullstack)

### Integration Points:

- Dashboard & Reports (Dev1) perlu data dari Pipeline (Dev2)
- Visit Report (Dev1) bisa link ke Task (Dev2)
- Deal (Dev2) bisa link ke Product (Dev2)
- Deal (Dev2) perlu link ke Account (Dev1)

### Coordination:

- [ ] Week 3: Coordinate API contract untuk integration points
- [ ] Week 7: Mid-sprint review - check integration points
- [ ] Week 11: Pre-integration review
- [ ] Week 14: Final integration testing

---

## 📝 Notes

1. **Fullstack Development**: Setiap modul dikerjakan fullstack sampai selesai
2. **No Dependencies**: Tidak bergantung ke Dev1, bisa dikerjakan paralel
3. **Hackathon Mode**: Tidak ada unit test, manual testing saja
4. **Code Review**: Lakukan code review sebelum merge
5. **Documentation**: Update documentation setelah setiap sprint
6. **Postman Collection**: Update Postman collection untuk setiap modul

---

**Dokumen ini akan diupdate sesuai dengan progress development.**
