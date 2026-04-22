# Sprint Planning - Developer 3 (Mobile Developer)
## CRM Healthcare/Pharmaceutical Platform - Sales CRM

**Developer**: Mobile Developer (Flutter)  
**Role**: Develop mobile app (Flutter) untuk sales rep  
**Versi**: 1.0  
**Status**: Active  
**Last Updated**: 2025-01-15

**⚠️ Catatan API Integration**:
- Mobile app saat ini menggunakan endpoint API versi **Web** karena endpoint khusus untuk mobile belum tersedia.
- Response parsing dibuat fleksibel untuk menangani berbagai format response API.
- Endpoint mobile-specific akan digunakan ketika sudah tersedia dari Developer 2.

> **📊 Visualisasi Project**: Lihat [**PROJECT_DIAGRAMS.md**](../PROJECT_DIAGRAMS.md) untuk memahami scope, fitur, dan user flow secara visual.

---

## 📋 Overview

Developer 3 bertanggung jawab untuk:
- **Mobile App Development**: Develop Flutter mobile app untuk sales rep
- **API Integration**: Integrate dengan backend APIs yang dibuat Developer 2
- **Mobile-Specific Features**: GPS tracking, photo upload, push notifications
- **UI/UX**: Modern dan intuitive mobile UI/UX

---

## 🎯 Sprint Details

### Sprint 0: Flutter Project Setup (Week 1-2)

**Goal**: Setup Flutter project dan foundation

**Tasks**:
- [x] Setup Flutter project structure
- [x] Setup state management (Provider / Riverpod)
- [x] Setup HTTP client (Dio)
- [x] Setup local storage (Hive / SharedPreferences)
- [x] Setup dependency injection
- [x] Create API client dengan interceptors
- [x] Setup authentication flow
- [x] Create login screen
- [x] Create auth guard / route protection
- [x] Setup error handling
- [x] Setup loading states
- [x] Setup theme (light/dark)
- [x] Setup navigation

**Acceptance Criteria**:
- ✅ Flutter project structure ready
- ✅ Authentication flow bekerja
- ✅ API client setup dengan interceptors
- ✅ Error handling comprehensive
- ✅ Navigation structure ready

**Estimated Time**: 4-5 days

---

### Sprint 1: Account & Contact Mobile (Week 3-4)

**Goal**: Implement Account & Contact di mobile app

**Backend Dependencies** (Coordinate dengan Developer 2):
- Account APIs ready
- Contact APIs ready

**Note**: 
- ⚠️ **API yang digunakan**: Menggunakan endpoint API versi **Web** (`/api/v1/accounts`, `/api/v1/contacts`) karena endpoint khusus untuk mobile belum tersedia.
- Response format API di-handle secara fleksibel untuk menangani format array langsung atau object dengan `items` & `pagination`.

**Mobile Tasks**:
- [x] Create account service (`account_service.dart`)
- [x] Create contact service (`contact_service.dart`)
- [x] Create account models (`models/account.dart`)
- [x] Create contact models (`models/contact.dart`)
- [x] Create account list screen (`screens/accounts/account_list.dart`)
- [x] Create account detail screen (`screens/accounts/account_detail.dart`)
- [x] Create contact list screen (`screens/contacts/contact_list.dart`)
- [x] Create contact detail screen (`screens/contacts/contact_detail.dart`)
- [x] Add account search functionality
- [x] Add contact search functionality
- [x] Add pull-to-refresh
- [x] Add pagination
- [x] Create account card widget
- [x] Create contact card widget

**Acceptance Criteria**:
- ✅ Sales rep dapat melihat list accounts
- ✅ Sales rep dapat melihat account detail
- ✅ Sales rep dapat melihat list contacts
- ✅ Sales rep dapat melihat contact detail
- ✅ Search dan pagination bekerja
- ✅ UI/UX modern dan intuitive

**Testing**:
- Test account list dan detail
- Test contact list dan detail
- Test search functionality
- Test pagination

**Estimated Time**: 5-6 days

---

### Sprint 2: Visit Report Mobile (Week 5-6)

**Goal**: Implement Visit Report di mobile app dengan GPS dan photo

**Backend Dependencies** (Coordinate dengan Developer 2):
- Visit report APIs ready
- Check-in/out APIs ready
- Photo upload API ready

**Note**: 
- ⚠️ **API yang akan digunakan**: Akan menggunakan endpoint API versi **Web** (`/api/v1/visit-reports`) jika endpoint mobile belum tersedia.
- Response parsing akan mengikuti pattern yang sama dengan Sprint 1 (fleksibel untuk berbagai format).

**Mobile Tasks**:
- [x] Create visit report service (`visit_report_service.dart`)
- [x] Create visit report models (`models/visit_report.dart`)
- [x] Create visit report list screen (`screens/visit_reports/visit_report_list.dart`)
- [x] Create visit report detail screen (`screens/visit_reports/visit_report_detail.dart`)
- [x] Create visit report form screen (`screens/visit_reports/visit_report_form.dart`)
- [x] Integrate GPS location (geolocator package)
- [x] Implement check-in dengan GPS
- [x] Implement check-out dengan GPS
- [x] Integrate camera untuk photo upload (image_picker package)
- [x] Implement photo upload functionality
- [x] Add visit report status badge
- [x] Add pull-to-refresh
- [x] Add pagination
- [x] Create visit report card widget

**Acceptance Criteria**:
- ✅ Sales rep dapat melihat list visit reports
- ✅ Sales rep dapat create visit report
- ✅ Check-in/out dengan GPS bekerja
- ✅ Photo upload bekerja
- ✅ UI/UX modern dan intuitive

**Testing**:
- Test visit report creation
- Test check-in/out dengan GPS
- Test photo upload
- Test visit report list

**Estimated Time**: 7-8 days

---

### Sprint 3: Task & Reminder Mobile (Week 7-8)

**Goal**: Implement Task & Reminder di mobile app dengan push notifications

**Backend Dependencies** (Coordinate dengan Developer 2):
- Task APIs ready
- Reminder APIs ready
- Push notification service ready

**Note**: 
- ⚠️ **API yang akan digunakan**: Akan menggunakan endpoint API versi **Web** jika endpoint mobile belum tersedia.
- Push notification service harus ready untuk implementasi reminder notifications.

**Mobile Tasks**:
- [x] Create task repository (`task_repository.dart`)
- [x] Create task models (`models/task.dart`)
- [x] Create task list screen (`screens/tasks/task_list.dart`)
- [x] Create task detail screen (`screens/tasks/task_detail.dart`)
- [x] Create task form screen (`screens/tasks/task_form.dart`)
- [ ] Integrate push notifications (firebase_messaging package) - ⏳ **Menunggu Backend Push Notification Service**
- [x] Implement task reminder management (create/delete reminders)
- [x] Add task filter (status, priority)
- [x] Add pull-to-refresh
- [x] Add pagination
- [x] Create task card widget
- [x] Add task completion functionality

**Acceptance Criteria**:
- [x] Sales rep dapat melihat list tasks
- [x] Sales rep dapat create task
- [x] Sales rep dapat complete task
- [ ] Push notifications untuk task reminder bekerja - ⏳ **Menunggu Backend Push Notification Service**
- [x] UI/UX modern dan intuitive

**Testing**:
- Test task list dan detail
- Test task creation
- Test task completion
- Test push notifications

**Estimated Time**: 5-6 days

---

### Sprint 4: Dashboard Mobile (Week 9)

**Goal**: Implement Dashboard di mobile app

**Backend Dependencies** (Coordinate dengan Developer 2):
- Dashboard APIs ready

**Note**: 
- ⚠️ **Status Saat Ini**: Dashboard screen sudah dibuat dengan UI sederhana untuk navigasi ke semua fitur (Accounts, Contacts, Visit Reports, Tasks).
- UI dashboard lengkap dengan metrics & statistics akan diimplementasikan setelah Dashboard APIs ready.

**Note**: 
- ⚠️ **API yang digunakan**: Menggunakan endpoint API versi **Web** (`/api/v1/dashboard/overview`, `/api/v1/dashboard/recent-activities`) karena endpoint khusus untuk mobile belum tersedia.
- Response parsing dibuat fleksibel untuk menangani format response API.

**Mobile Tasks**:
- [x] Create dashboard screen (`screens/dashboard/dashboard.dart`) - ✅ **Selesai dengan data integration**
- [x] Create dashboard repository (`dashboard_repository.dart`) - ✅ **Menggunakan API versi Web**
- [x] Create dashboard models (`models/dashboard.dart`) - ✅ **Selesai**
- [x] Create stat cards widget - ✅ **Selesai**
- [x] Create pipeline summary widget - ✅ **Selesai**
- [x] Create recent activities widget - ✅ **Selesai**
- [x] Add pull-to-refresh - ✅ **Selesai**
- [x] Add period selector (today/week/month/year) - ✅ **Selesai**

**Acceptance Criteria**:
- ✅ Sales rep dapat melihat dashboard
- ✅ Key metrics ditampilkan (visits, accounts, activities, revenue)
- ✅ Pipeline summary ditampilkan
- ✅ Recent activities ditampilkan
- ✅ Period selector (today/week/month/year) bekerja
- ✅ Pull-to-refresh bekerja
- ✅ UI/UX modern dan intuitive

**Testing**:
- Test dashboard data loading
- Test dashboard widgets

**Estimated Time**: 3-4 days

---

### Sprint 5: Mobile App Polish & Optimization (Week 10-11)

**Goal**: Polish mobile app dan optimize performance

**Tasks**:
- [x] Improve UI/UX consistency
- [x] Add loading states di semua screens
- [x] Add error handling di semua screens
- [x] Optimize app performance
- [x] Add offline support (basic read-only) - ✅ **Phase 1 Completed**
  - [x] Setup Hive storage infrastructure
  - [x] Implement connectivity service
  - [x] Create offline storage helper
  - [x] Add offline indicator UI
  - [ ] Update repositories for offline-first (pending - see docs/mobile/OFFLINE_SUPPORT_IMPLEMENTATION.md)
- [x] Add app icon dan splash screen
- [x] Add app versioning
- [ ] Test di multiple devices (Android & iOS)
- [x] Fix bugs
- [x] Performance optimization

**Acceptance Criteria**:
- ✅ UI/UX konsisten
- ✅ Performance optimal
- ✅ Tidak ada critical bugs
- ✅ App ready untuk testing

**Testing**:
- Test di multiple devices
- Performance testing
- Bug testing

**Estimated Time**: 5-6 days

---

### Sprint 6: Integration & Final Testing (Week 12-13)

**Goal**: Integration dengan backend dan final testing

**Tasks**:
- [ ] Coordinate dengan Developer 2 untuk API integration
- [ ] Integration testing
- [ ] Fix integration issues
- [ ] End-to-end testing
- [ ] Test semua features
- [ ] Final bug fixes
- [ ] Prepare APK untuk demo
- [ ] Prepare iOS build (jika diperlukan)

**Acceptance Criteria**:
- ✅ Semua features terintegrasi dengan baik
- ✅ Tidak ada critical bugs
- ✅ APK ready untuk demo
- ✅ App ready untuk production

**Testing**:
- Integration testing
- End-to-end testing
- User acceptance testing

**Estimated Time**: 4-5 days

---

## 📊 Sprint Summary

| Sprint | Goal | Duration | Status |
|--------|------|----------|--------|
| Sprint 0 | Flutter Setup | 4-5 days | ✅ Completed |
| Sprint 1 | Account & Contact | 5-6 days | ✅ Completed |
| Sprint 2 | Visit Report | 7-8 days | ✅ Completed |
| Sprint 3 | Task & Reminder | 5-6 days | ✅ Completed |
| Sprint 4 | Dashboard | 3-4 days | ✅ Completed |
| Sprint 5 | Polish & Optimization | 5-6 days | ⏳ Pending |
| Sprint 6 | Integration & Testing | 4-5 days | ⏳ Pending |

**Total Estimated Time**: 33-42 days (4.5-6 weeks)

---

## 🔗 Dependencies

### Dependencies dari Developer 2 (Backend)
- Account & Contact APIs (Sprint 1) ✅ **Menggunakan API versi Web**
- Visit Report APIs (Sprint 2) ✅ **Menggunakan API versi Web**
- Task APIs (Sprint 3) ✅ **Menggunakan API versi Web**
- Dashboard APIs (Sprint 4) ✅ **Menggunakan API versi Web**
- Push notification service (Sprint 3) ⏳ **Menunggu Backend Push Notification Service**

### ⚠️ Catatan Penting: API Integration
- **Sprint 1 (Account & Contact)**: Menggunakan endpoint API versi **Web** (`/api/v1/accounts`, `/api/v1/contacts`) karena endpoint khusus untuk mobile belum tersedia.
- **Sprint 2 (Visit Report)**: Menggunakan endpoint API versi **Web** (`/api/v1/visit-reports`) karena endpoint khusus untuk mobile belum tersedia.
- **Sprint 3 (Task & Reminder)**: Menggunakan endpoint API versi **Web** (`/api/v1/tasks`) karena endpoint khusus untuk mobile belum tersedia.
- **Sprint 4 (Dashboard)**: Menggunakan endpoint API versi **Web** (`/api/v1/dashboard/overview`, `/api/v1/dashboard/recent-activities`) karena endpoint khusus untuk mobile belum tersedia.
- Response parsing dibuat fleksibel untuk menangani format array langsung atau object dengan `items` & `pagination`.
- Endpoint mobile-specific akan digunakan ketika sudah tersedia dari Developer 2.

### Dependencies dari Developer 1 (Web)
- Coordinate untuk UI/UX consistency
- Coordinate untuk feature requirements

---

## 📱 Mobile App Features

### Core Features (MVP)
1. **Authentication**
   - Login
   - Logout
   - Token refresh

2. **Account & Contact**
   - List accounts
   - Account detail
   - List contacts
   - Contact detail
   - Search

3. **Visit Report**
   - List visit reports
   - Create visit report
   - Check-in/out dengan GPS
   - Photo upload
   - Visit report detail

4. **Task & Reminder**
   - List tasks
   - Create task
   - Complete task
   - Push notifications

5. **Dashboard**
   - Basic dashboard
   - Visit statistics
   - Recent activities

### Optional Features (Future)
- Offline support
- Advanced dashboard
- Sales pipeline view
- Advanced search

---

## 📝 Notes

1. **API Integration**: Coordinate dengan Developer 2 untuk API design
2. **UI/UX**: Coordinate dengan Developer 1 untuk consistency
3. **Testing**: Test di multiple devices (Android & iOS)
4. **Performance**: Optimize untuk mobile performance
5. **Push Notifications**: Setup Firebase Cloud Messaging

---

**Dokumen ini akan diupdate sesuai dengan progress development.**

