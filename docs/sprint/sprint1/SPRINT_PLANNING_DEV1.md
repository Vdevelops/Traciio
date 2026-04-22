# Sprint Planning - Developer 1 (Fullstack Developer)

## CRM Healthcare/Pharmaceutical Platform - Sales CRM

**Developer**: Fullstack Developer (Go Backend + Next.js Frontend)  
**Role**: Develop modul-modul Sales CRM secara fullstack (backend + frontend)  
**Versi**: 2.0  
**Status**: Active  
**Last Updated**: 2025-01-15

> **📊 Visualisasi Project**: Lihat [**PROJECT_DIAGRAMS.md**](../PROJECT_DIAGRAMS.md) untuk memahami scope, fitur, dan user flow secara visual.

---

## 📋 Overview

Developer 1 bertanggung jawab untuk:

- **Fullstack Development**: Develop modul-modul yang ditugaskan secara lengkap (backend API + frontend)
- **Backend**: Go (Gin) APIs untuk modul yang ditugaskan
- **Frontend**: Next.js 16 frontend untuk modul yang ditugaskan
- **Database**: Design dan implement database schema untuk modul yang ditugaskan
- **Postman Collection**: Update Postman collection untuk modul yang ditugaskan

**Modul yang ditugaskan ke Dev1**:

1. ✅ Account & Contact Management (Fullstack)
2. ✅ Visit Report & Activity Tracking (Fullstack)
3. ✅ Dashboard & Reports (Fullstack)
4. ✅ Settings (Fullstack)

**Parallel Development Strategy**:

- ✅ **TIDAK bergantung ke Dev2** - bisa dikerjakan paralel
- ✅ Setiap modul dikerjakan fullstack sampai selesai
- ✅ **Hackathon mode** - tidak ada unit test
- ✅ Manual testing saja

---

## 🎯 Sprint Details

### Sprint 0: Foundation Review & Improvement (Week 1)

**Goal**: Review dan perbaiki foundation yang sudah ada

**Tasks**:

- [x] Review authentication flow (login, token refresh)
- [x] Fix bugs di sprint 0 jika ada
- [x] Improve error handling di frontend
- [x] Optimize API client interceptors
- [x] Improve loading states dan skeletons
- [x] Review dan fix auth guard component

**Acceptance Criteria**:

- ✅ Login flow bekerja dengan baik
- ✅ Token refresh otomatis bekerja
- ✅ Error handling konsisten
- ✅ Loading states smooth

**Estimated Time**: 3-4 days

---

### Sprint 1: User Management Review & Improvement (Week 1-2)

**Goal**: Review dan perbaiki user management module

**Tasks**:

- [x] Review user list page dan fix bugs
- [x] Improve user form validation
- [x] Optimize user list table (pagination, search, filter)
- [x] Fix sidebar permission integration
- [x] Improve user form UI/UX
- [x] Add user detail page jika belum ada

**Acceptance Criteria**:

- ✅ User CRUD bekerja dengan baik
- ✅ Search dan filter bekerja optimal
- ✅ Form validation comprehensive
- ✅ UI/UX modern dan intuitive

**Estimated Time**: 3-4 days

---

### Sprint 2: Master Data Cleanup (Week 2)

**Goal**: Archive atau adapt master data yang tidak relevan untuk Sales CRM

**Tasks**:

- [x] Review diagnosis & procedures module
- [x] Archive atau mark sebagai optional module
- [x] Update sidebar untuk remove/hide diagnosis & procedures menu
- [x] Cleanup unused components jika ada
- [x] Update navigation structure

**Acceptance Criteria**:

- ✅ Diagnosis & Procedures tidak muncul di menu utama (atau di archive)
- ✅ Navigation clean dan hanya menampilkan Sales CRM modules
- ✅ Tidak ada broken links

**Estimated Time**: 1-2 days

---

### Sprint 3: Account & Contact Management (Fullstack) (Week 3-4)

**Goal**: Implement Account & Contact Management secara fullstack (backend + frontend)

**Backend Tasks**:

- [x] Create account model dan migration
- [x] Create contact model dan migration
- [x] Create account repository interface dan implementation
- [x] Create contact repository interface dan implementation
- [x] Create account service
- [x] Create contact service
- [x] Implement account list API (`GET /api/v1/accounts`)
- [x] Implement account detail API (`GET /api/v1/accounts/:id`)
- [x] Implement create account API (`POST /api/v1/accounts`)
- [x] Implement update account API (`PUT /api/v1/accounts/:id`)
- [x] Implement delete account API (`DELETE /api/v1/accounts/:id`)
- [x] Implement contact list API (`GET /api/v1/contacts`)
- [x] Implement contact detail API (`GET /api/v1/contacts/:id`)
- [x] Implement create contact API (`POST /api/v1/contacts`)
- [x] Implement update contact API (`PUT /api/v1/contacts/:id`)
- [x] Implement delete contact API (`DELETE /api/v1/contacts/:id`)
- [x] Implement account-contact relationship APIs (via foreign key and preload)
- [x] Add account search API (integrated in list API with search param)
- [x] Add contact search API (integrated in list API with search param)
- [x] Add pagination support
- [x] Add validation
- [x] Add account categories seeder (via menu/permission seeder integration)

**Frontend Tasks**:

- [x] Create account types (`types/account.d.ts`)
- [x] Create contact types (`types/contact.d.ts`)
- [x] Create account service (`accountService`)
- [x] Create contact service (`contactService`)
- [x] Create account list page (`/accounts`)
- [x] Create account form component (`AccountForm`)
- [x] Create account detail page (view functionality in list, detail page can be added later)
- [x] Create contact list page (`/contacts`)
- [x] Create contact form component (`ContactForm`)
- [x] Create contact detail page (view functionality in list, detail page can be added later)
- [x] Add account search and filter
- [x] Add contact search and filter
- [x] Create account selector component (integrated in ContactForm)
- [x] Create contact selector component (can be added when needed for other forms)
- [x] Add account-contact relationship UI (contact form shows account selector, contact list shows account name)

**Postman Collection**:

- [x] Add account APIs ke Postman collection (Web section)
- [x] Add contact APIs ke Postman collection (Web section)

**Acceptance Criteria**:

- ✅ Account CRUD APIs bekerja dengan baik
- ✅ Contact CRUD APIs bekerja dengan baik
- ✅ Account-contact relationship bekerja
- ✅ Frontend terintegrasi dengan backend APIs
- ✅ Search dan filter bekerja optimal
- ✅ Form validation comprehensive
- ✅ UI/UX modern dan intuitive
- ✅ Postman collection updated

**Testing** (Manual testing):

- Test account CRUD (backend + frontend)
- Test contact CRUD (backend + frontend)
- Test account-contact relationship
- Test search and filter

**Estimated Time**: 6-7 days

---

### Sprint 4: Visit Report & Activity Tracking (Fullstack) (Week 5-6)

**Goal**: Implement Visit Report & Activity Tracking secara fullstack (backend + frontend)

**Backend Tasks**:

- [x] Create visit report model dan migration
- [x] Create activity model dan migration
- [x] Create visit report repository interface dan implementation
- [x] Create activity repository interface dan implementation
- [x] Create visit report service
- [x] Create activity service
- [x] Implement visit report list API (`GET /api/v1/visit-reports`)
- [x] Implement visit report detail API (`GET /api/v1/visit-reports/:id`)
- [x] Implement create visit report API (`POST /api/v1/visit-reports`)
- [x] Implement update visit report API (`PUT /api/v1/visit-reports/:id`)
- [x] Implement delete visit report API (`DELETE /api/v1/visit-reports/:id`)
- [x] Implement check-in API (`POST /api/v1/visit-reports/:id/check-in`)
- [x] Implement check-out API (`POST /api/v1/visit-reports/:id/check-out`)
- [x] Implement approve visit report API (`POST /api/v1/visit-reports/:id/approve`)
- [x] Implement reject visit report API (`POST /api/v1/visit-reports/:id/reject`)
- [x] Implement photo upload API (`POST /api/v1/visit-reports/:id/photos`)
- [x] Implement activity list API (`GET /api/v1/activities`)
- [x] Implement activity timeline API (`GET /api/v1/accounts/:id/activities`)
- [x] Add GPS location tracking
- [x] Add pagination support
- [x] Add validation
- [x] Add file storage setup
- [x] Add visit report and activity seeders

**Frontend Tasks**:

- [x] Create visit report types (`types/visit-report.d.ts`)
- [x] Create activity types (`types/activity.d.ts`)
- [x] Create visit report service (`visitReportService`)
- [x] Create activity service (`activityService`)
- [x] Create visit report list page (`/visit-reports`)
- [x] Create visit report form component (`VisitReportForm`)
- [x] Create visit report detail page (`/visit-reports/[id]`)
- [x] Create activity timeline component (`ActivityTimeline`)
- [x] Create photo upload component
- [x] Create visit report status badge component
- [x] Add visit report search and filter
- [x] Create supervisor review UI (approve/reject)
- [x] Add activity timeline di account detail page

**Postman Collection**:

- [x] Add visit report APIs ke Postman collection (Web section)
- [x] Add visit report APIs ke Postman collection (Mobile section)
- [x] Add activity APIs ke Postman collection (Web section)

**Acceptance Criteria**:

- ✅ Visit report CRUD APIs bekerja dengan baik
- ✅ Check-in/out APIs bekerja dengan GPS
- ✅ Photo upload bekerja
- ✅ Activity tracking bekerja
- ✅ Supervisor approve/reject bekerja
- ✅ Frontend terintegrasi dengan backend APIs
- ✅ Search dan filter bekerja optimal
- ✅ Postman collection updated (Web + Mobile)

**Testing** (Manual testing):

- Test visit report CRUD (backend + frontend)
- Test check-in/out dengan GPS
- Test photo upload
- Test activity tracking
- Test supervisor workflow

**Estimated Time**: 7-8 days

---

### Sprint 5: Dashboard & Reports (Fullstack) (Week 11-12)

**Goal**: Implement Dashboard & Reports secara fullstack (backend + frontend)

**Backend Tasks**:

- [x] Create dashboard service
- [x] Create report service
- [x] Implement dashboard overview API (`GET /api/v1/dashboard/overview`)
- [x] Implement visit statistics API (`GET /api/v1/dashboard/visits`)
- [x] Implement pipeline summary API (`GET /api/v1/dashboard/pipeline`)
- [x] Implement top accounts API (`GET /api/v1/dashboard/top-accounts`)
- [x] Implement top sales rep API (`GET /api/v1/dashboard/top-sales-rep`)
- [x] Implement recent activities API (`GET /api/v1/dashboard/recent-activities`)
- [x] Implement visit report API (`GET /api/v1/reports/visit-reports`)
- [x] Implement sales pipeline report API (`GET /api/v1/reports/pipeline`)
- [x] Implement sales performance report API (`GET /api/v1/reports/sales-performance`)
- [x] Implement account activity report API (`GET /api/v1/reports/account-activity`)
- [x] Add date range filtering
- [x] Add export functionality (PDF/Excel) - basic

**Frontend Tasks**:

- [x] Create dashboard service (`dashboardService`)
- [x] Create report service (`reportService`)
- [x] Create dashboard types (`types/dashboard.d.ts`)
- [x] Create main dashboard page (`/dashboard`)
- [x] Create dashboard overview component
- [x] Create visit statistics component
- [x] Create pipeline summary component
- [x] Create top accounts component
- [x] Create top sales rep component
- [x] Create recent activities component
- [x] Create charts (using recharts atau similar)
- [x] Create reports list page (`/reports`)
- [x] Create report generator component
- [x] Create report viewer component
- [x] Add date range picker
- [x] Add export functionality (PDF/Excel)

**Postman Collection**:

- [x] Add dashboard APIs ke Postman collection (Web section)
- [x] Add report APIs ke Postman collection (Web section)

**Menu & Permissions**:

- [x] Add Dashboard menu to menu seeder
- [x] Add Reports menu to menu seeder
- [x] Add Dashboard permissions to permission seeder
- [x] Add Reports permissions to permission seeder

**Acceptance Criteria**:

- ✅ Dashboard APIs bekerja dengan baik
- ✅ Report APIs bekerja dengan baik
- ✅ Frontend terintegrasi dengan backend APIs
- ✅ Dashboard menampilkan key metrics
- ✅ Charts dan graphs ditampilkan dengan benar
- ✅ Date range filtering bekerja
- ✅ Reports dapat di-generate dan di-export
- ✅ Postman collection updated

**Testing** (Manual testing):

- Test dashboard data loading (backend + frontend)
- Test date range filtering
- Test chart rendering
- Test report generation
- Test export functionality

**Estimated Time**: 6-7 days

---

### Sprint 6: Settings (Fullstack) (Week 13)

**Goal**: Implement Settings secara fullstack (backend + frontend)

**Backend Tasks**:

- [x] Create settings model dan migration
- [x] Create settings repository interface dan implementation
- [x] Create settings service
- [x] Implement get settings API (`GET /api/v1/settings`)
- [x] Implement update settings API (`PUT /api/v1/settings`)
- [x] Implement general settings API
- [x] Implement notification settings API
- [x] Implement pipeline settings API
- [x] Add validation

**Frontend Tasks**:

- [x] Create settings service (`settingsService`)
- [x] Create settings types (`types/settings.d.ts`)
- [x] Create settings dashboard (`/settings`)
- [x] Create general settings page (`/settings/general`)
- [x] Create notification settings page (`/settings/notifications`)
- [x] Create pipeline settings page (`/settings/pipeline`)
- [x] Improve UI/UX consistency
- [x] Add loading states di semua pages
- [x] Add error boundaries
- [x] Optimize bundle size
- [x] Cross-browser testing

**Postman Collection**:

- [x] Add settings APIs ke Postman collection (Web section)

**Menu & Permissions**:

- [x] Add Settings menu to menu seeder
- [x] Add Settings permissions to permission seeder

**Seeder**:

- [x] Add settings seeder data with related IDs

**Acceptance Criteria**:

- ✅ Settings APIs bekerja dengan baik
- ✅ Settings tersimpan dengan benar
- ✅ Frontend terintegrasi dengan backend APIs
- ✅ Admin dapat manage settings
- ✅ UI/UX konsisten di semua pages
- ✅ Loading states smooth
- ✅ Error handling comprehensive
- ✅ Performance optimal
- ✅ Postman collection updated

**Testing** (Manual testing):

- Test settings update (backend + frontend)
- Test settings persistence
- Test UI/UX consistency
- Test performance

**Estimated Time**: 3-4 days

---

### Sprint 7: Integration & Final Testing (Week 14)

**Goal**: Integration dengan modul Dev2 dan final testing

**Tasks**:

- [ ] Coordinate dengan Developer 2 untuk integration
- [ ] Test integration antara modul Dev1 dan Dev2
- [ ] Fix integration issues
- [ ] End-to-end testing
- [ ] Performance testing
- [ ] Security testing
- [ ] Final bug fixes
- [ ] Documentation update

**Acceptance Criteria**:

- ✅ Semua modules terintegrasi dengan baik
- ✅ Tidak ada critical bugs
- ✅ Performance acceptable
- ✅ Security audit passed

**Testing**:

- End-to-end testing
- Performance testing
- Security testing

**Estimated Time**: 3-4 days

### 🔍 Analisis Sprint 7 – Integration & Final Testing (Dev1 & Dev2)

**Gambaran integrasi saat ini**

- **Front‑office terpadu**: Sidebar, auth guard, dan permission di Dev1 mengontrol akses ke modul Dev1 (Accounts, Visit Reports, Dashboard, Settings) sekaligus modul Dev2 (Pipeline, Tasks, Products) sehingga user memiliki satu pintu masuk ke seluruh Sales CRM.
- **Integrasi data lintas modul**: Dashboard & Reports Dev1 mengonsumsi data dari modul Pipeline (Dev2); Account & Contact Dev1 menjadi _source of truth_ untuk Deals/Tasks Dev2; Visit Report & Activity Dev1 memanfaatkan Task/Activity Dev2 untuk timeline kunjungan; Settings Dev1 mengatur konfigurasi pipeline yang digunakan Dev2.

**Kekuatan integrasi**

- **Model data konsisten**: ID dan relasi (`account_id`, `contact_id`, `deal_id`, `product_id`) dirancang agar semua laporan dan dashboard dapat menggabungkan data dari kedua tim tanpa _join_ yang aneh.
- **Standarisasi API**: Dev1 dan Dev2 mengikuti `api-response-standards.md` dan `api-error-codes.md`, sehingga frontend dan Postman dapat memakai pola yang sama untuk semua endpoint.
- **Navigasi & izin selaras**: Struktur menu & permission yang sama mengurangi risiko akses tidak sinkron antara modul Dev1 dan Dev2.

**Kekurangan & risiko integrasi**

- **Belum ada skenario E2E yang terdokumentasi dengan baik**: Walaupun rencana Sprint 7 menyebut _end‑to‑end testing_, alur pengguna lintas modul (misalnya _lead → account → deal + product → visit report → task follow‑up → reports_) belum tertulis sebagai skenario baku yang bisa diulang setiap rilis.
- **Tidak ada automated integration/regression test lintas tim**: Integrasi saat ini mengandalkan manual testing dan koordinasi lisan; perubahan di satu sisi (misal skema API pipeline/report) bisa merusak modul lain tanpa deteksi awal.
- **Inkonstistensi UX kecil**: Beberapa halaman Dev1 dan Dev2 memiliki pola filter, pagination, atau tampilan error yang sedikit berbeda sehingga pengalaman user kurang seragam saat melintasi modul.\
- **Kurangnya monitoring lintas layanan**: Belum ada health‑check/monitor yang secara eksplisit memantau endpoint Dev2 yang dipakai Dashboard/Reports, sehingga analisis insiden integrasi masih membutuhkan _manual tracing_.

**Rekomendasi tindak lanjut lintas tim**

- Mendefinisikan 3–5 **user‑journey E2E resmi** dan memasukkannya ke dokumen Sprint 7 sebagai acuan regression test bersama Dev1 & Dev2.
- Menambahkan **API contract test sederhana** (misal melalui pipeline CI atau koleksi Postman otomatis) untuk endpoint yang menjadi _integration points_ utama (dashboard overview, pipeline summary, visit report timeline, dsb.).
- Menyelaraskan guideline **UI/UX lintas modul** (pola filter, loading, error) agar perpindahan antar modul terasa mulus bagi user.
- Menambah **monitoring ringan** (grafik error/latency) untuk endpoint integrasi utama sehingga _health_ integrasi dapat dipantau tanpa harus melakukan inspeksi manual setiap kali ada isu.

---

## 📊 Sprint Summary

| Sprint   | Goal                            | Duration | Status       |
| -------- | ------------------------------- | -------- | ------------ |
| Sprint 0 | Foundation Review               | 3-4 days | ✅ Completed |
| Sprint 1 | User Management Review          | 3-4 days | ✅ Completed |
| Sprint 2 | Master Data Cleanup             | 1-2 days | ✅ Completed |
| Sprint 3 | Account & Contact (Fullstack)   | 6-7 days | ✅ Completed |
| Sprint 4 | Visit Report (Fullstack)        | 7-8 days | ✅ Completed |
| Sprint 5 | Dashboard & Reports (Fullstack) | 6-7 days | ✅ Completed |
| Sprint 6 | Settings (Fullstack)            | 3-4 days | ✅ Completed |
| Sprint 7 | Integration & Testing           | 3-4 days | ⏳ Pending   |

**Total Estimated Time**: 33-40 days (4.7-5.7 weeks)

---

## 🔗 Coordination dengan Dev2

### Modul yang dikerjakan Dev2 (untuk referensi):

- Sales Pipeline (Fullstack)
- Task & Reminder (Fullstack)
- Product Management (Fullstack)

### Integration Points:

- Dashboard & Reports perlu data dari Pipeline (Dev2)
- Visit Report bisa link ke Task (Dev2)
- Deal bisa link ke Product (Dev2)

### Coordination:

- [ ] Week 3: Coordinate API contract untuk integration points
- [ ] Week 7: Mid-sprint review - check integration points
- [ ] Week 11: Pre-integration review
- [ ] Week 14: Final integration testing

---

## 📝 Notes

1. **Fullstack Development**: Setiap modul dikerjakan fullstack sampai selesai
2. **No Dependencies**: Tidak bergantung ke Dev2, bisa dikerjakan paralel
3. **Hackathon Mode**: Tidak ada unit test, manual testing saja
4. **Code Review**: Lakukan code review sebelum merge
5. **Documentation**: Update documentation setelah setiap sprint
6. **Postman Collection**: Update Postman collection untuk setiap modul

---

**Dokumen ini akan diupdate sesuai dengan progress development.**
