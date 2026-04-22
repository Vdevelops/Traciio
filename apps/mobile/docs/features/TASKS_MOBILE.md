# Mobile Tasks - Feature Documentation

## Overview

Mobile Tasks adalah fitur untuk sales rep untuk melihat, mengelola, dan melacak tugas-tugas mereka. Fitur ini dirancang khusus untuk mobile dengan fokus pada kemudahan akses dan manajemen tugas yang efisien.

**Note:** Istilah "Tasks" tetap menggunakan bahasa Inggris di semua bahasa (tidak diterjemahkan) untuk konsistensi.

## Fitur Utama

### 1. **List Tasks**

Sales rep dapat melihat daftar tugas mereka dengan:
- Pagination (page-based untuk infinite scroll)
- Filter by status (all, pending, in_progress, completed, cancelled)
- Filter by priority (all, low, medium, high, urgent)
- Filter by type (all, general, call, email, meeting, follow_up)
- Filter by due date (from dan to)
- Search functionality

**API Endpoint:**
- `GET /api/v1/mobile/tasks/my-tasks?page=1&per_page=20&status=...&priority=...&type=...&due_date_from=...&due_date_to=...&search=...`

**Features:**
- Infinite scroll support dengan page-based pagination
- Offline caching untuk akses cepat (60 detik TTL)
- Pull-to-refresh untuk update data
- Optimized list rendering dengan `addRepaintBoundaries` dan `cacheExtent`
- Search dengan debounce (500ms)
- Real-time search results dengan clear search button

**UI Features:**
- **Task Card**: Menampilkan title sebagai judul utama, description (jika ada), priority badge, type badge, status badge, due date, account, dan contact
- **Priority Badge**: Badge dengan warna berbeda untuk setiap priority (urgent: red, high: orange, medium: blue, low: gray)
- **Type Badge**: Badge dengan warna dan icon berbeda untuk setiap type (call: blue, email: purple, meeting: teal, follow_up: orange, general: primary)
- **Status Badge**: Badge dengan warna berbeda untuk setiap status (completed: green, in_progress: blue, pending: orange, cancelled: gray)
- **Search Indicator**: Menampilkan query search yang aktif dan tombol clear search
- **Filter Indicator**: Red dot pada filter icon jika ada filter aktif

**Performance Optimizations:**
- List caching dengan TTL 60 detik
- Lazy loading dengan infinite scroll
- Optimized widget rendering
- Stable keys untuk list items

**Security:**
- Hanya menampilkan tasks milik user yang login
- Filter di backend berdasarkan `assigned_to` dari JWT token

### 2. **View Task Details**

Sales rep dapat melihat detail lengkap task, termasuk:
- Task information (title, description, status, priority, type, due date)
- Related information (account, contact, deal jika ada)
- Reminders (jika ada)
- Action buttons berdasarkan status (Mark In Progress, Complete Task, Add Reminder)

**API Endpoint:**
- `GET /api/v1/tasks/:id`

**Features:**
- **Always Fresh Data**: Setiap kali user membuka detail task, data selalu di-fetch dari API (tidak menggunakan cache saat online)
- **Cache Fallback**: Cache hanya digunakan saat offline atau jika API gagal
- **Status Badge**: Badge dengan warna berbeda untuk setiap status
- **Type Badge**: Badge dengan warna dan icon berbeda untuk setiap type
- **Action Buttons**: Tombol action hanya muncul sesuai status:
  - "Mark In Progress" hanya untuk status `pending`
  - "Complete Task" hanya untuk status yang bukan `completed`
  - "Add Reminder" selalu tersedia

**Security:**
- Ownership validation (hanya owner yang dapat melihat)
- Data selalu fresh dari API untuk memastikan informasi terbaru

### 3. **Mark Task In Progress**

Sales rep dapat menandai task dengan status `pending` menjadi `in_progress`.

**API Endpoint:**
- `POST /api/v1/tasks/:id/mark-in-progress`

**Restrictions:**
- Hanya dapat dilakukan jika status = `pending`
- Hanya dapat dilakukan oleh owner
- Konfirmasi dialog sebelum mark in progress

**Validations:**
- Status validation di backend
- Ownership validation di backend

**Security:**
- Ownership validation di backend
- Status validation di backend

### 4. **Complete Task**

Sales rep dapat menandai task sebagai completed.

**API Endpoint:**
- `POST /api/v1/tasks/:id/complete`

**Restrictions:**
- Hanya dapat dilakukan jika status bukan `completed` atau `cancelled`
- Hanya dapat dilakukan oleh user dengan permission `can_complete_task`
- Konfirmasi dialog sebelum complete

**Validations:**
- Permission validation di backend
- Status validation di backend

**Security:**
- Permission-based access
- Ownership validation di backend

### 5. **Add Reminder**

Sales rep dapat menambahkan reminder untuk task.

**API Endpoint:**
- `POST /api/v1/tasks/:id/reminders`

**Features:**
- Date and time picker untuk reminder
- Optional message untuk reminder
- Reminder ditampilkan di task detail

**Validations:**
- Reminder date wajib dipilih
- Reminder message opsional

**Security:**
- Ownership validation di backend

### 6. **Search Tasks**

Sales rep dapat mencari tasks dengan keyword.

**Features:**
- Search modal dengan overlay design
- Debounce search (500ms) untuk mengurangi API calls
- Real-time search results
- Clear search button dengan indicator query yang aktif
- Empty state yang berbeda untuk search vs no search

**UI Features:**
- Search modal muncul dari atas dengan slide animation
- Search query ditampilkan di atas list dengan tombol clear
- Empty state menampilkan "No results found for {query}" untuk search
- Empty state menampilkan "No tasks found" untuk no search

**Performance:**
- Debounce 500ms untuk mengurangi API calls
- Search results di-cache dengan TTL 60 detik

### 7. **Filter Tasks**

Sales rep dapat memfilter tasks berdasarkan:
- **Status**: All, Pending, In Progress, Completed, Cancelled
- **Priority**: All, Low, Medium, High, Urgent
- **Type**: All, General, Call, Email, Meeting, Follow Up
- **Due Date**: From dan To date range

**Features:**
- Filter sheet dengan bottom sheet design
- Multiple filters dapat aktif bersamaan
- Clear filters button untuk reset semua filter
- Filter indicator (red dot) pada filter icon jika ada filter aktif
- Filter state dipertahankan saat refresh

**UI Features:**
- Filter menu di AppBar dengan popup menu
- Filter sheet untuk setiap jenis filter
- Check mark untuk filter yang aktif
- Clear filters option di popup menu jika ada filter aktif

**Performance:**
- Filter state di-cache
- Filter di-apply saat load tasks

## Alur Sales (Sales Flow)

### 1. **View Tasks**

Sales rep melihat daftar tasks mereka:
1. Buka aplikasi mobile
2. Navigate ke "Tasks" menu di navbar
3. List tasks ditampilkan dengan pagination
4. User dapat scroll untuk load more tasks (infinite scroll)
5. User dapat pull-to-refresh untuk update data

**Status:** Semua status tasks ditampilkan

**UI Flow:**
- List menggunakan TaskCard untuk setiap task
- Badge status, priority, dan type ditampilkan dengan warna berbeda
- Due date ditampilkan dengan format yang user-friendly (Due today, Due tomorrow, Overdue by X days, Due MMM dd, yyyy)
- Account dan contact ditampilkan jika ada

### 2. **Search Tasks**

Sales rep mencari tasks dengan keyword:
1. Tap search icon di AppBar
2. Search modal muncul dari atas
3. User mengetik keyword (debounce 500ms)
4. Search results ditampilkan secara real-time
5. Query search ditampilkan di atas list dengan tombol clear
6. User dapat tap "Clear Search" untuk menghapus search

**UI Flow:**
- Search modal dengan overlay design
- Search query ditampilkan di atas list
- Empty state berbeda untuk search vs no search

### 3. **Filter Tasks**

Sales rep memfilter tasks:
1. Tap filter icon di AppBar
2. Popup menu muncul dengan opsi filter
3. User memilih jenis filter (Status, Priority, Type, Due Date)
4. Filter sheet muncul dari bawah
5. User memilih filter value
6. Filter di-apply dan list di-refresh
7. Filter indicator (red dot) muncul pada filter icon

**UI Flow:**
- Filter menu dengan popup menu design
- Filter sheet untuk setiap jenis filter
- Multiple filters dapat aktif bersamaan
- Clear filters option tersedia

### 4. **View Task Details**

Sales rep melihat detail task:
1. Tap task card di list
2. Task detail screen dibuka
3. Data selalu di-fetch fresh dari API (tidak menggunakan cache saat online)
4. Detail task ditampilkan dengan informasi lengkap
5. Action buttons ditampilkan berdasarkan status

**Status:** Semua status tasks dapat dilihat detailnya

**UI Flow:**
- Task detail dengan sections (Task Information, Related Information, Reminders)
- Badge status dan type dengan warna berbeda
- Action buttons berdasarkan status dan permission

### 5. **Mark Task In Progress**

Sales rep menandai task sebagai in progress:
1. Buka task detail dengan status `pending`
2. Tap "Mark In Progress" button
3. Konfirmasi dialog muncul
4. Tap "Mark In Progress" di dialog
5. Task status berubah menjadi `in_progress`
6. Task detail di-refresh dengan data terbaru
7. List tasks di-refresh untuk update status

**Status:** `pending` → `in_progress`

**Validations:**
- Status harus `pending`
- Ownership validation

### 6. **Complete Task**

Sales rep menandai task sebagai completed:
1. Buka task detail dengan status bukan `completed` atau `cancelled`
2. Tap "Complete Task" button (jika memiliki permission)
3. Konfirmasi dialog muncul
4. Tap "Complete Task" di dialog
5. Task status berubah menjadi `completed`
6. Task detail di-refresh dengan data terbaru
7. List tasks di-refresh untuk update status

**Status:** `pending` atau `in_progress` → `completed`

**Validations:**
- Status bukan `completed` atau `cancelled`
- Permission `can_complete_task` required
- Ownership validation

### 7. **Add Reminder**

Sales rep menambahkan reminder untuk task:
1. Buka task detail
2. Tap "Add Reminder" button
3. Dialog muncul dengan date/time picker dan message field
4. User memilih reminder date dan time
5. (Optional) User memasukkan reminder message
6. Tap "Save"
7. Reminder dibuat dan ditampilkan di task detail

**Features:**
- Date and time picker untuk reminder
- Optional message field
- Reminder ditampilkan di task detail dengan status (sent/not sent)

## Status Flow Diagram

```
[pending]
  ↓ (mark in progress)
[in_progress]
  ↓ (complete)
[completed]

[pending] atau [in_progress]
  ↓ (cancel - jika diimplementasikan)
[cancelled]
```

## Validations

### Form Validations

#### Search
- Search query dapat berupa string kosong (untuk clear search)
- Search query di-debounce 500ms untuk mengurangi API calls

#### Filter
- Status filter: null (all), 'pending', 'in_progress', 'completed', 'cancelled'
- Priority filter: null (all), 'low', 'medium', 'high', 'urgent'
- Type filter: null (all), 'general', 'call', 'email', 'meeting', 'follow_up'
- Due date filter: DateTime? (from dan to, opsional)

### Business Rule Validations

1. **Status-based Actions:**
   - `pending`: Dapat mark in progress
   - `in_progress`: Dapat complete task
   - `completed`: Tidak dapat diubah (hanya view)
   - `cancelled`: Tidak dapat diubah (hanya view)

2. **Ownership:**
   - Semua operasi (view, mark in progress, complete, add reminder) hanya dapat dilakukan oleh owner
   - Backend memvalidasi ownership untuk semua operasi
   - Mobile API hanya menampilkan tasks milik user yang login

3. **Permission-based Actions:**
   - Complete task memerlukan permission `can_complete_task`
   - Permission di-validate di backend

## Security

### Authentication & Authorization

1. **JWT Token:**
   - Semua API calls menggunakan JWT token untuk authentication
   - Token di-refresh otomatis jika expired
   - Token disimpan securely di device

2. **Ownership Validation:**
   - Backend memvalidasi ownership untuk semua operasi
   - User hanya dapat mengakses tasks miliknya sendiri
   - Mobile API filter berdasarkan `assigned_to` dari JWT token

3. **Permission-based Access:**
   - Complete task memerlukan permission `can_complete_task`
   - Permission di-validate di backend

### Data Security

1. **Input Validation:**
   - Semua input divalidasi di client dan server
   - SQL injection prevention (parameterized queries)
   - XSS prevention (output escaping)

2. **Cache Security:**
   - Cached data di-encrypt di device
   - Cache TTL untuk mencegah data stale
   - Cache di-clear saat logout
   - Task detail selalu fetch fresh dari API saat online (tidak menggunakan cache)

### Network Security

1. **HTTPS:**
   - Semua API calls menggunakan HTTPS
   - Certificate pinning (jika diimplementasikan)

2. **Request Validation:**
   - Rate limiting di backend
   - Request size limits
   - Content-Type validation

### Offline Security

1. **Cached Data:**
   - Cached data di-encrypt di device
   - Cache TTL untuk mencegah data stale
   - Cache di-clear saat logout
   - Task detail tidak menggunakan cache saat online (hanya saat offline)

2. **Sensitive Data:**
   - Task data tidak mengandung informasi sensitif
   - Cache hanya digunakan untuk offline fallback

## Performance Optimizations

### List Performance

1. **Caching:**
   - List data di-cache dengan TTL 60 detik
   - Cache key berdasarkan page, search query, dan filters
   - Cache di-invalidate saat refresh

2. **Lazy Loading:**
   - Infinite scroll dengan page-based pagination
   - Load more saat scroll mencapai 80% dari bottom
   - Loading indicator saat load more

3. **Widget Optimization:**
   - `addRepaintBoundaries: true` untuk mengurangi repaints
   - `cacheExtent: 500` untuk cache lebih banyak items
   - Stable keys (`ValueKey`) untuk list items
   - Optimized TaskCard rendering

4. **Search Optimization:**
   - Debounce 500ms untuk mengurangi API calls
   - Search results di-cache dengan TTL 60 detik
   - Clear search untuk reset state

### Detail Screen Performance

1. **Always Fresh Data:**
   - Task detail selalu fetch fresh dari API saat online
   - Cache hanya digunakan sebagai fallback saat offline atau jika API gagal
   - Data di-invalidate saat screen dibuka untuk memastikan data terbaru

2. **Lazy Loading:**
   - Related information di-load secara lazy
   - Reminders di-load secara lazy

### Filter Performance

1. **Filter State:**
   - Filter state di-cache
   - Filter di-apply saat load tasks
   - Multiple filters dapat aktif bersamaan

## Theme & Language Support

### Theme Support

1. **Light/Dark Theme:**
   - Full support untuk light dan dark theme
   - Semua components menggunakan `Theme.of(context)` dan `colorScheme`
   - Automatic theme switching berdasarkan system preference atau user selection

2. **Theme-aware Components:**
   - Cards menggunakan `colorScheme.surface` dan `colorScheme.surfaceContainerHighest`
   - Borders dan shadows disesuaikan untuk dark theme
   - Icons dan text colors menggunakan `colorScheme.onSurface`
   - Badge colors menggunakan warna yang konsisten (tidak bergantung pada theme)

3. **Badge Colors:**
   - **Status Badge:**
     - Completed: Green
     - In Progress: Blue
     - Pending: Orange
     - Cancelled: Gray
   - **Priority Badge:**
     - Urgent: Red
     - High: Orange
     - Medium: Blue
     - Low: Gray
   - **Type Badge:**
     - Call: Blue
     - Email: Purple
     - Meeting: Teal
     - Follow Up: Orange
     - General: Primary color

### Language Support

1. **Supported Languages:**
   - English (en)
   - Indonesian (id)

2. **Localization:**
   - Semua text menggunakan `AppLocalizations` (l10n)
   - Date/time formatting menggunakan locale-aware `DateFormat`
   - Number formatting menggunakan locale-aware formatting
   - Task types dan priorities menggunakan localization keys

3. **Special Cases:**
   - "Tasks" tetap menggunakan bahasa Inggris di semua bahasa (tidak diterjemahkan)
   - Technical terms tetap menggunakan bahasa Inggris untuk konsistensi
   - Task types (General, Call, Email, Meeting, Follow Up) menggunakan localization
   - Task priorities (Low, Medium, High, Urgent) menggunakan localization

4. **Language Switching:**
   - Language dapat diubah di Profile screen
   - Perubahan langsung diterapkan tanpa restart app
   - Locale di-save di local storage

## Error Handling

### Network Errors

1. **Offline Detection:**
   - Connectivity service untuk detect online/offline status
   - Offline indicator di UI
   - Cached data ditampilkan saat offline
   - Task detail menggunakan cache sebagai fallback saat offline

2. **Error Messages:**
   - User-friendly error messages
   - Retry mechanism untuk failed requests
   - Error logging untuk debugging

### Validation Errors

1. **Form Validation:**
   - Real-time validation dengan error messages
   - Field-level error messages
   - Summary error messages

2. **Business Rule Errors:**
   - Clear error messages untuk business rule violations
   - Suggestions untuk fix errors

### API Errors

1. **Error Handling:**
   - Error messages dari API ditampilkan dengan jelas
   - Retry mechanism untuk failed requests
   - Fallback ke cache jika API gagal

## Best Practices

1. **Always Check Status**: Pastikan status task sesuai sebelum melakukan action
2. **Use Filters**: Gunakan filter untuk memudahkan pencarian tasks
3. **Use Search**: Gunakan search untuk mencari tasks dengan keyword
4. **Check Permissions**: Pastikan memiliki permission sebelum melakukan action
5. **Refresh Data**: Pull-to-refresh untuk mendapatkan data terbaru
6. **Check Theme**: Pastikan theme sesuai dengan preferensi (light/dark)
7. **Check Language**: Pastikan language sesuai dengan preferensi (English/Indonesian)

## Future Enhancements

1. **Create/Update/Delete Tasks**: Fitur untuk create, update, dan delete tasks (saat ini tidak tersedia untuk sales users)
2. **Task Templates**: Template untuk tasks yang sering dibuat
3. **Task Recurrence**: Fitur untuk membuat recurring tasks
4. **Task Dependencies**: Fitur untuk membuat task dependencies
5. **Task Comments**: Fitur untuk menambahkan comments pada tasks
6. **Task Attachments**: Fitur untuk menambahkan attachments pada tasks
7. **Task Notifications**: Push notifications untuk task reminders dan updates
8. **Task Analytics**: Dashboard analytics untuk task performance
9. **Task Export**: Export tasks ke PDF atau Excel
10. **Task Sharing**: Fitur untuk sharing tasks dengan team members
