# Sprint Planning - Master
## CRM Healthcare/Pharmaceutical Platform - Sales CRM

**Versi**: 2.0  
**Status**: Active  
**Last Updated**: 2025-01-15  
**Product Type**: Sales CRM untuk Perusahaan Farmasi

---

## 📋 Daftar Isi

1. [Overview](#overview)
2. [Team Structure](#team-structure)
3. [Sprint Overview](#sprint-overview)
4. [Developer Sprint Plans](#developer-sprint-plans)
5. [Coordination & Dependencies](#coordination--dependencies)
6. [Timeline](#timeline)
7. [Sprint Checklist](#sprint-checklist)
8. [Dependencies](#dependencies)
9. [📊 Project Diagrams](#-project-diagrams)

> **📊 Visualisasi Project**: Untuk memahami scope, fitur, dan user flow secara visual, lihat [**PROJECT_DIAGRAMS.md**](./PROJECT_DIAGRAMS.md) yang berisi:
> - Project Scope Overview
> - Feature Modules Diagram
> - User Flow Diagrams (Sales Rep, Supervisor, Admin)
> - Input/Output Diagrams per Feature
> - System Architecture Diagram

---

## Overview

Dokumen ini adalah **master planning** untuk development Sales CRM dengan **3 developers**:
- **Developer 1**: Web Developer (Full-stack, fokus Web)
- **Developer 2**: Backend Developer (Fokus BE, sedikit FE)
- **Developer 3**: Mobile Developer (Flutter)

Setiap developer memiliki sprint planning terpisah yang detail.

### Product Scope

**Sales CRM untuk Perusahaan Farmasi** dengan fitur:
- Account & Contact Management
- Visit Report & Activity Tracking
- Sales Pipeline Management
- Task & Reminder
- Product Management
- Dashboard & Reports
- Mobile App (Flutter)

### Technology Stack

- **Backend**: Go (Gin framework)
- **Web Frontend**: Next.js 16 (App Router, Server Components)
- **Mobile App**: Flutter
- **Database**: PostgreSQL
- **State Management (Web)**: Zustand
- **Styling (Web)**: Tailwind CSS v4
- **UI Components (Web)**: shadcn/ui v4

---

## Team Structure

### Developer 1: Fullstack Developer
- **Focus**: Fullstack development (Go Backend + Next.js Frontend)
- **Responsibilities**:
  - Perbaiki sprint 0-2 yang sudah ada
  - Develop modul-modul yang ditugaskan secara fullstack
  - Backend APIs + Frontend untuk modul yang ditugaskan
- **Modul yang ditugaskan**:
  - Account & Contact Management (Fullstack)
  - Visit Report & Activity Tracking (Fullstack)
  - Dashboard & Reports (Fullstack)
  - Settings (Fullstack)
- **Sprint Planning**: [`SPRINT_PLANNING_DEV1.md`](./SPRINT_PLANNING_DEV1.md)

### Developer 2: Fullstack Developer
- **Focus**: Fullstack development (Go Backend + Next.js Frontend)
- **Responsibilities**:
  - Develop modul-modul yang ditugaskan secara fullstack
  - Backend APIs + Frontend untuk modul yang ditugaskan
  - Database design dan migration
  - Update Postman collection
- **Modul yang ditugaskan**:
  - Sales Pipeline Management (Fullstack)
  - Task & Reminder Management (Fullstack)
  - Product Management (Fullstack)
- **Sprint Planning**: [`SPRINT_PLANNING_DEV2.md`](./SPRINT_PLANNING_DEV2.md)

### Developer 3: Mobile Developer
- **Focus**: Mobile app (Flutter)
- **Responsibilities**:
  - Develop Flutter mobile app
  - Integrate dengan backend APIs
  - Mobile-specific features (GPS, camera, push notifications)
- **Sprint Planning**: [`SPRINT_PLANNING_DEV3.md`](./SPRINT_PLANNING_DEV3.md)

---

## Sprint Overview

### Master Timeline (100 Hari / ~14 Minggu)

| Week | Developer 1 (Fullstack) | Developer 2 (Fullstack) | Developer 3 (Mobile) |
|------|-------------------------|------------------------|---------------------|
| 1-2 | Foundation Review, User Management Review, Master Data Cleanup | Foundation Review, User Management Review | Flutter Setup |
| 3-4 | Account & Contact (Fullstack) | - | Account & Contact Mobile |
| 5-6 | Visit Report (Fullstack) | - | Visit Report Mobile |
| 7-8 | - | Sales Pipeline (Fullstack) | Task & Reminder Mobile |
| 9 | - | Task & Reminder (Fullstack) | Dashboard Mobile |
| 10 | - | Product Management (Fullstack) | Mobile Polish |
| 11-12 | Dashboard & Reports (Fullstack) | - | Mobile Integration |
| 13 | Settings (Fullstack) | - | Final Testing |
| 14 | Integration & Testing | Integration & Testing | Final Testing |

---

## Developer Sprint Plans

### Developer 1: Fullstack Developer
📄 **Detail Sprint Planning**: [`SPRINT_PLANNING_DEV1.md`](./SPRINT_PLANNING_DEV1.md)

**Sprint Summary**:
- Sprint 0: Foundation Review (3-4 days)
- Sprint 1: User Management Review (3-4 days)
- Sprint 2: Master Data Cleanup (1-2 days)
- Sprint 3: Account & Contact (Fullstack) (6-7 days)
- Sprint 4: Visit Report (Fullstack) (7-8 days)
- Sprint 5: Dashboard & Reports (Fullstack) (6-7 days)
- Sprint 6: Settings (Fullstack) (3-4 days)
- Sprint 7: Integration & Testing (3-4 days)

**Total**: 33-40 days (4.7-5.7 weeks)

**Modul yang dikerjakan**:
- ✅ Account & Contact Management (Backend + Frontend)
- ✅ Visit Report & Activity Tracking (Backend + Frontend)
- ✅ Dashboard & Reports (Backend + Frontend)
- ✅ Settings (Backend + Frontend)

### Developer 2: Fullstack Developer
📄 **Detail Sprint Planning**: [`SPRINT_PLANNING_DEV2.md`](./SPRINT_PLANNING_DEV2.md)

**Sprint Summary**:
- Sprint 0: Foundation Review (2-3 days)
- Sprint 1: User Management Review (2-3 days)
- Sprint 2: Sales Pipeline (Fullstack) (6-7 days)
- Sprint 3: Task & Reminder (Fullstack) (4-5 days)
- Sprint 4: Product Management (Fullstack) (3-4 days)
- Sprint 5: Integration & Testing (3-4 days)

**Total**: 20-26 days (2.9-3.7 weeks)

**Modul yang dikerjakan**:
- ✅ Sales Pipeline Management (Backend + Frontend)
- ✅ Task & Reminder Management (Backend + Frontend)
- ✅ Product Management (Backend + Frontend)

### Developer 3: Mobile Developer
📄 **Detail Sprint Planning**: [`SPRINT_PLANNING_DEV3.md`](./SPRINT_PLANNING_DEV3.md)

**Sprint Summary**:
- Sprint 0: Flutter Setup (4-5 days) ✅ **Completed**
- Sprint 1: Account & Contact Mobile (5-6 days) ✅ **Completed**
- Sprint 2: Visit Report Mobile (7-8 days) ✅ **Completed**
- Sprint 3: Task & Reminder Mobile (5-6 days) ✅ **Completed** (menggunakan API versi Web, push notification menunggu backend)
- Sprint 4: Dashboard Mobile (3-4 days) ⚠️ **UI Navigation selesai, menunggu APIs**
- Sprint 5: Polish & Optimization (5-6 days)
- Sprint 6: Integration & Testing (4-5 days)

**Total**: 33-42 days (4.5-6 weeks)

**Catatan API Integration**:
- ⚠️ **Sprint 1**: Menggunakan endpoint API versi **Web** karena endpoint mobile belum tersedia.
- ⚠️ **Sprint 2**: Menggunakan endpoint API versi **Web** karena endpoint mobile belum tersedia.
- ⚠️ **Sprint 3**: Menggunakan endpoint API versi **Web** karena endpoint mobile belum tersedia. Push notification menunggu backend service.
- Response parsing dibuat fleksibel untuk kompatibilitas dengan format API yang berbeda.

---

## Coordination & Dependencies

### Parallel Development Strategy

1. **Developer 1 & Developer 2**: Bekerja paralel, tidak ada dependencies
   - Setiap developer mengerjakan modul secara fullstack (backend + frontend)
   - Tidak perlu menunggu satu sama lain
   - Integration dilakukan di akhir (Week 14)

2. **Integration Points** (Week 14):
   - Dashboard & Reports (Dev1) perlu data dari Pipeline (Dev2)
   - Visit Report (Dev1) bisa link ke Task (Dev2)
   - Deal (Dev2) perlu link ke Account (Dev1)
   - Deal (Dev2) link ke Product (Dev2)

3. **Developer 1 ↔ Developer 3**: Coordinate untuk UI/UX consistency
4. **Developer 2 ↔ Developer 3**: Coordinate untuk API consistency

### Coordination Points

1. **Week 1**: Kickoff meeting - align scope dan timeline
2. **Week 3**: API design review - Developer 2 present API design ke Developer 1 & 3
3. **Week 5**: Mid-sprint review - check progress dan blockers
4. **Week 8**: Integration planning - plan integration antara web, mobile, dan backend
5. **Week 11**: Pre-demo review - prepare untuk demo
6. **Week 14**: Final review - final testing dan delivery

### API Design Coordination

- **Developer 1 & Developer 2** design APIs untuk modul masing-masing
- Semua APIs harus mengikuti API response standards
- Coordinate integration points di Week 3 dan Week 11
- Postman collection harus terpisah untuk Web dan Mobile

---

## Timeline

### Overall Timeline (100 Hari)

**Week 1-2**: Foundation & Setup
- Developer 1: Review dan perbaiki sprint 0-2
- Developer 2: Review foundation dan API standards
- Developer 3: Flutter project setup

**Week 3-4**: Account & Contact
- Developer 1: Account & Contact (Fullstack - Backend + Frontend)
- Developer 2: Foundation Review, User Management Review
- Developer 3: Account & Contact Mobile ✅ **Completed** (menggunakan API versi Web)

**Week 5-6**: Visit Report
- Developer 1: Visit Report (Fullstack - Backend + Frontend)
- Developer 2: Foundation Review, User Management Review
- Developer 3: Visit Report Mobile

**Week 7-8**: Pipeline
- Developer 1: (Buffer/Polish)
- Developer 2: Sales Pipeline (Fullstack - Backend + Frontend)
- Developer 3: Task & Reminder Mobile ✅ **Completed**

**Week 9-10**: Task & Product
- Developer 1: (Buffer/Polish)
- Developer 2: Task & Reminder (Fullstack), Product Management (Fullstack)
- Developer 3: Dashboard Mobile, Mobile Polish

**Week 11-12**: Dashboard & Reports
- Developer 1: Dashboard & Reports (Fullstack - Backend + Frontend)
- Developer 2: (Buffer/Polish)
- Developer 3: Mobile Integration

**Week 13**: Settings
- Developer 1: Settings (Fullstack - Backend + Frontend)
- Developer 2: (Buffer/Polish)
- Developer 3: Final Testing

**Week 14**: Integration & Delivery
- All Developers: Integration, Testing, Final Delivery

---

## Sprint Details

> **Note**: Detail sprint untuk setiap developer ada di file terpisah:
> - [`SPRINT_PLANNING_DEV1.md`](./SPRINT_PLANNING_DEV1.md) - Web Developer (10 sprints)
> - [`SPRINT_PLANNING_DEV2.md`](./SPRINT_PLANNING_DEV2.md) - Backend Developer (9 sprints)
> - [`SPRINT_PLANNING_DEV3.md`](./SPRINT_PLANNING_DEV3.md) - Mobile Developer (6 sprints)

### Completed Sprints (Legacy)

**Sprint 0**: Project Setup & Foundation ✅ (COMPLETED)  
**Sprint 1**: User Management Module ✅ (COMPLETED)  
**Sprint 2**: Master Data - Diagnosis & Procedures ✅ (COMPLETED - ARCHIVED)

> **Note**: Sprint 2 (Diagnosis & Procedures) di-archive karena tidak relevan untuk Sales CRM.  
> Dapat digunakan sebagai optional module di future jika diperlukan.

---

## Sprint Checklist

### Before Starting Sprint
- [ ] Review sprint goal dan tasks
- [ ] Check dependencies dari previous sprints
- [ ] Coordinate dengan developers lain
- [ ] Setup development environment
- [ ] Create feature branch

### During Sprint
- [ ] Follow coding standards
- [ ] Write tests untuk new features
- [ ] Update documentation
- [ ] Commit frequently dengan clear messages
- [ ] Coordinate dengan developers lain jika ada blockers

### After Sprint
- [ ] Test all acceptance criteria
- [ ] Code review
- [ ] Update sprint status
- [ ] Deploy to staging (jika applicable)
- [ ] Demo ke stakeholders

---

## Dependencies

### Sprint Dependency Graph (Sales CRM)

```
Sprint 0 (Foundation) ✅
    ↓
Sprint 1 (User Management) ✅
    ↓
Sprint 2 (Master Data Cleanup) ✅
    ↓
┌─────────────────────────────────────────┐
│  PARALLEL DEVELOPMENT (No Dependencies) │
└─────────────────────────────────────────┘
    ↓                    ↓
Dev1:                 Dev2:
- Account & Contact   - Sales Pipeline
- Visit Report        - Task & Reminder
- Dashboard & Reports - Product Management
- Settings
    ↓                    ↓
┌─────────────────────────────────────────┐
│  Integration & Testing (Week 14)        │
└─────────────────────────────────────────┘
```

---

## Estimated Timeline

- **Total Duration**: 100 hari (~14 minggu)
- **Team Size**: 3 developers
  - Developer 1 (Fullstack): 33-40 days (4.7-5.7 weeks)
  - Developer 2 (Fullstack): 20-26 days (2.9-3.7 weeks)
  - Developer 3 (Mobile): 33-42 days (4.5-6 weeks)

---

## Notes

1. **Flexibility**: Sprint dapat di-adjust sesuai kebutuhan
2. **Parallel Work**: Beberapa sprints dapat dikerjakan parallel jika tidak ada dependency
3. **Testing**: Setiap sprint harus include testing
4. **Documentation**: Update documentation setelah setiap sprint
5. **Code Review**: Lakukan code review sebelum merge
6. **Coordination**: Coordinate dengan developers lain untuk dependencies

---

**Dokumen ini adalah master planning. Untuk detail sprint, lihat:**
- [`SPRINT_PLANNING_DEV1.md`](./SPRINT_PLANNING_DEV1.md) - Web Developer
- [`SPRINT_PLANNING_DEV2.md`](./SPRINT_PLANNING_DEV2.md) - Backend Developer
- [`SPRINT_PLANNING_DEV3.md`](./SPRINT_PLANNING_DEV3.md) - Mobile Developer

---

**Dokumen ini akan diupdate sesuai dengan progress development.**
