# Mobile Dashboard - Feature Documentation

## Overview

Mobile Dashboard adalah halaman utama untuk sales rep yang menampilkan informasi penting dalam format yang sederhana dan mudah diakses. Dashboard dirancang khusus untuk mobile dengan fokus pada quick access ke fitur utama dan informasi yang paling relevan.

**Note:** Istilah "Visits" tetap menggunakan bahasa Inggris di semua bahasa (tidak diterjemahkan) untuk konsistensi.

## Fitur Utama

### 1. **Target Progress Card**

Menampilkan progress target penjualan user dengan visualisasi yang jelas.

**Components:**
- **Target Amount**: Jumlah target yang harus dicapai (format: Rp)
- **Achieved Amount**: Jumlah yang sudah dicapai (format: Rp)
- **Progress Bar**: Visualisasi progress dengan color-coded indicator
  - Green (>80%): Target tercapai dengan baik
  - Yellow (50-80%): Target sedang dalam progress
  - Red (<50%): Target perlu perhatian lebih
- **Progress Percentage**: Persentase progress dengan 1 desimal
- **Remaining Amount**: Jumlah yang masih harus dicapai (format: Rp dengan abbreviation untuk millions/billions)

**UI Features:**
- **Theme Support**: Full support untuk light/dark theme
- **Language Support**: Full support untuk bahasa Indonesia dan Inggris
- **Currency Formatting**: 
  - Millions: "Rp 10M remaining" (EN) atau "Rp 10Jt tersisa" (ID)
  - Billions: "Rp 1B remaining" (EN) atau "Rp 1M tersisa" (ID)

**API Endpoint:**
- `GET /api/v1/dashboard/mobile/overview?period=today`

**Data Source:**
- Data target diambil dari `target` object dalam response API
- Hanya menampilkan jika `target.targetAmount > 0`

**Validations:**
- Progress percentage dihitung dari `achievedAmount / targetAmount * 100`
- Progress bar value di-clamp antara 0.0 dan 1.0

### 2. **Quick Actions Widget**

Widget untuk quick access ke fitur utama: Create Visit dan Create Task.

**Components:**
- **Create Visit Button**: Navigate ke form create visit report
- **Create Task Button**: Navigate ke task list (untuk sales users yang biasanya tidak bisa create task)

**UI Features:**
- **Permission-based Display**: Create Visit button hanya muncul jika user memiliki CREATE permission untuk visit-reports
- **Theme Support**: Full support untuk light/dark theme
- **Language Support**: Full support untuk bahasa Indonesia dan Inggris
- **Color Contrast**: Create Task button menggunakan warna indigo untuk kontras yang lebih baik

**Permissions:**
- CREATE permission untuk visit-reports (diperlukan untuk Create Visit button)
- Task button selalu muncul (navigate ke task list, bukan create task)

### 3. **Visits Section**

Menampilkan daftar visit reports dengan horizontal scrollable list.

**Components:**
- **Section Header**: 
  - Icon dan title "Visits" (tetap dalam bahasa Inggris)
  - "See All" button untuk navigate ke full visits list
- **Filter Tabs**: 
  - All (Semua)
  - Planned (Terencana)
  - Completed (Selesai)
  - Cancelled (Dibatalkan)
- **Visit Cards**: Horizontal scrollable list dengan maksimal items yang ditampilkan
- **Empty State**: Menampilkan pesan jika tidak ada visits

**UI Features:**
- **Horizontal Scrollable**: Cards di-scroll secara horizontal (left-right)
- **Filter Support**: Filter berdasarkan status (all, planned, completed, cancelled)
- **Theme Support**: Full support untuk light/dark theme
- **Language Support**: Full support untuk bahasa Indonesia dan Inggris (kecuali "Visits" tetap dalam bahasa Inggris)
- **Optimized Rendering**: 
  - `addRepaintBoundaries: true` untuk mengurangi repaints
  - `cacheExtent: 500` untuk cache lebih banyak items
  - Stable keys (`ValueKey`) untuk list items

**Visit Card Components:**
- **Account Name**: Nama account yang dikunjungi
- **Status Badge**: Badge dengan warna sesuai status (Planned, Completed, Cancelled)
- **Time Range**: Waktu check-in dan check-out (jika ada)
- **Location**: Alamat check-in location (jika ada)

**API Endpoint:**
- `GET /api/v1/dashboard/mobile/visits?status=all&limit=5`

**Data Source:**
- Data visits diambil dari API dan difilter berdasarkan status
- Hanya menampilkan visits milik user yang login
- Default limit: 5 items untuk performance

**Validations:**
- Status filter validation (all, planned, completed, cancelled)
- Empty state handling jika tidak ada visits

### 4. **Upcoming Tasks Section**

Menampilkan daftar upcoming tasks yang di-assign ke user.

**Components:**
- **Section Header**: 
  - Icon dan title "Upcoming Tasks" (diterjemahkan)
  - "See All" button untuk navigate ke full tasks list
- **Task Cards**: List vertical dengan maksimal 5 items
- **Empty State**: Menampilkan pesan jika tidak ada upcoming tasks

**UI Features:**
- **Fixed List**: List tidak scrollable (fixed height), maksimal 5 items
- **Task Type Icons**: Icon berbeda berdasarkan task type (general, call, email, meeting, follow up)
- **Priority Badge**: Badge dengan warna sesuai priority (Urgent/High: red, Medium: orange, Low: green)
- **Theme Support**: Full support untuk light/dark theme
- **Language Support**: Full support untuk bahasa Indonesia dan Inggris
- **Locale-aware Date Formatting**: Date formatting menggunakan locale user

**Task Card Components:**
- **Task Type Icon**: Icon circular dengan warna sesuai task type
- **Task Title**: Judul task
- **Task Description**: Deskripsi task (jika ada, max 1 line dengan ellipsis)
- **Due Date & Time**: Tanggal dan waktu due date (locale-aware formatting)
- **Priority Badge**: Badge dengan label priority (Urgent, High, Medium, Low)

**API Endpoint:**
- `GET /api/v1/dashboard/mobile/tasks?limit=5`

**Data Source:**
- Data tasks diambil dari API
- Hanya menampilkan upcoming tasks yang di-assign ke user yang login
- Default limit: 5 items

**Validations:**
- Task type validation (general, call, email, meeting, follow up)
- Priority validation (urgent, high, medium, low)
- Empty state handling jika tidak ada tasks

### 5. **Dashboard Header**

Header dengan navigation dan quick actions.

**Components:**
- **Profile Avatar**: Avatar user dari profile API
- **Search Icon**: Quick access ke task search
- **Notifications**: Notification count badge

**UI Features:**
- **Floating Bubble Header**: Header tidak fixed, scroll dengan konten
- **Theme Support**: Full support untuk light/dark theme
- **Language Support**: Full support untuk bahasa Indonesia dan Inggris

## Alur User (User Flow)

### 1. **Initial Load**

Saat user membuka dashboard:
1. Dashboard header ditampilkan (dengan avatar, search, notifications)
2. Loading state ditampilkan untuk semua sections
3. Data di-fetch secara parallel:
   - Overview (target progress)
   - Visits (dengan filter default: all)
   - Tasks (upcoming tasks)
4. Data ditampilkan setelah semua sections selesai loading
5. Jika ada error pada satu section, section lain tetap ditampilkan (partial data)

**Status:** Loading → Data Loaded / Partial Data / Error

### 2. **View Target Progress**

User dapat melihat progress target:
1. Target card ditampilkan di bagian atas dashboard
2. User dapat melihat:
   - Target amount dan achieved amount
   - Progress bar dengan color-coded indicator
   - Progress percentage
   - Remaining amount dengan format yang user-friendly

**Actions:**
- Tidak ada action yang dapat dilakukan di target card (read-only)

### 3. **Quick Actions**

User dapat melakukan quick actions:
1. **Create Visit**:
   - Tap "Create Visit" button
   - Navigate ke form create visit report
   - Setelah create, kembali ke dashboard dan data di-refresh
2. **Create Task**:
   - Tap "Create Task" button
   - Navigate ke task list screen
   - User dapat melihat dan manage tasks

### 4. **View and Filter Visits**

User dapat melihat dan filter visits:
1. Visits section ditampilkan dengan default filter "All"
2. User dapat:
   - Tap filter tab untuk mengubah filter (All, Planned, Completed, Cancelled)
   - Scroll horizontal untuk melihat semua visit cards
   - Tap visit card untuk navigate ke visit detail
   - Tap "See All" untuk navigate ke full visits list
3. Data di-refresh saat filter berubah

**Filter Behavior:**
- Filter "All": Menampilkan semua visits
- Filter "Planned": Menampilkan visits dengan status planned
- Filter "Completed": Menampilkan visits dengan status completed
- Filter "Cancelled": Menampilkan visits dengan status cancelled

### 5. **View Upcoming Tasks**

User dapat melihat upcoming tasks:
1. Upcoming tasks section ditampilkan dengan maksimal 5 items
2. User dapat:
   - Melihat task details (title, description, due date/time, priority)
   - Tap task card untuk navigate ke task detail
   - Tap "See All" untuk navigate ke full tasks list
3. Tasks ditampilkan dengan:
   - Task type icon (general, call, email, meeting, follow up)
   - Priority badge dengan warna sesuai priority
   - Locale-aware date formatting

### 6. **Refresh Dashboard**

User dapat refresh dashboard:
1. Pull down pada dashboard untuk trigger refresh
2. Semua sections di-refresh secara parallel
3. Data di-update setelah refresh selesai
4. Loading indicator ditampilkan saat refresh

## Validations

### Data Validations

1. **Target Progress:**
   - Target amount harus > 0 untuk menampilkan target card
   - Progress percentage dihitung dari `achievedAmount / targetAmount * 100`
   - Progress bar value di-clamp antara 0.0 dan 1.0

2. **Visits:**
   - Status filter validation (all, planned, completed, cancelled)
   - Empty state handling jika tidak ada visits
   - Visit card validation (account name, status, time, location)

3. **Tasks:**
   - Task type validation (general, call, email, meeting, follow up)
   - Priority validation (urgent, high, medium, low)
   - Empty state handling jika tidak ada tasks
   - Due date validation (harus valid date format)

### Business Rule Validations

1. **Data Ownership:**
   - Semua data (target, visits, tasks) hanya menampilkan data milik user yang login
   - API filter berdasarkan `sales_rep_id` dari JWT token

2. **Permission-based Access:**
   - Create Visit button hanya muncul jika user memiliki CREATE permission untuk visit-reports
   - Create Task button selalu muncul (navigate ke task list, bukan create task)

3. **Period Filter:**
   - Period filter (today, week, month) mempengaruhi data target dan overview
   - Default period: today

## Security

### Authentication & Authorization

1. **JWT Token:**
   - Semua API calls menggunakan JWT token untuk authentication
   - Token di-refresh otomatis jika expired
   - Token disimpan securely di device

2. **Data Ownership:**
   - Backend memvalidasi ownership untuk semua data
   - User hanya dapat melihat data miliknya sendiri
   - API filter berdasarkan `sales_rep_id` dari JWT token

3. **Permission-based Access:**
   - CREATE permission diperlukan untuk Create Visit button
   - Semua data lainnya berdasarkan ownership

### Data Security

1. **Cached Data:**
   - Cached data di-encrypt di device
   - Cache TTL untuk mencegah data stale
   - Cache di-clear saat logout

2. **Sensitive Data:**
   - Target amount dan achieved amount ditampilkan dengan format yang user-friendly
   - Tidak ada data sensitive yang di-expose di dashboard

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

2. **Offline Indicators:**
   - Offline indicator ditampilkan jika tidak ada koneksi
   - Cached data ditampilkan saat offline

## Performance Optimizations

### Data Fetching

1. **Parallel Loading:**
   - Overview, visits, dan tasks di-fetch secara parallel menggunakan `Future.wait`
   - `eagerError: false` untuk allow partial data jika satu section gagal

2. **Caching:**
   - Data di-cache dengan TTL untuk mengurangi API calls
   - Cache di-invalidate saat refresh
   - Background refresh jika online dan data cached tersedia

3. **Offline Support:**
   - Cached data ditampilkan saat offline
   - Background sync saat online kembali

### Widget Optimization

1. **List Performance:**
   - `addRepaintBoundaries: true` untuk mengurangi repaints
   - `cacheExtent: 500` untuk cache lebih banyak items
   - Stable keys (`ValueKey`) untuk list items
   - Horizontal scrollable list dengan optimized rendering

2. **Lazy Loading:**
   - Data di-load saat section pertama kali ditampilkan
   - Lazy loading untuk visit cards dan task cards

3. **Memory Management:**
   - Widget di-dispose dengan benar saat tidak digunakan
   - Cache management untuk mencegah memory leak

### Error Handling

1. **Partial Data:**
   - Jika satu section gagal, section lain tetap ditampilkan
   - Error message ditampilkan untuk section yang gagal
   - Retry mechanism untuk failed requests

2. **Network Errors:**
   - Offline detection dengan connectivity service
   - Cached data ditampilkan saat offline
   - Error messages yang user-friendly

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
   - Status badges menggunakan warna yang sesuai dengan theme

### Language Support

1. **Supported Languages:**
   - English (en)
   - Indonesian (id)

2. **Localization:**
   - Semua text menggunakan `AppLocalizations` (l10n)
   - Date/time formatting menggunakan locale-aware `DateFormat`
   - Number formatting menggunakan locale-aware formatting
   - Currency formatting dengan abbreviation untuk millions/billions

3. **Special Cases:**
   - "Visits" tetap menggunakan bahasa Inggris di semua bahasa (tidak diterjemahkan)
   - Technical terms tetap menggunakan bahasa Inggris untuk konsistensi

4. **Language Switching:**
   - Language dapat diubah di Profile screen
   - Perubahan langsung diterapkan tanpa restart app
   - Locale di-save di local storage

## API Endpoints Summary

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/dashboard/mobile/overview` | Get dashboard overview (target, stats) | Yes |
| GET | `/api/v1/dashboard/mobile/visits` | Get visits list (with filters) | Yes (owner only) |
| GET | `/api/v1/dashboard/mobile/tasks` | Get upcoming tasks | Yes (assigned to user only) |

## Error Handling

### Network Errors

1. **Offline Detection:**
   - Connectivity service untuk detect online/offline status
   - Offline indicator di UI
   - Cached data ditampilkan saat offline

2. **Error Messages:**
   - User-friendly error messages
   - Retry mechanism untuk failed requests
   - Error logging untuk debugging

### Validation Errors

1. **Data Validation:**
   - Empty state handling untuk sections yang tidak memiliki data
   - Error messages untuk invalid data

2. **Business Rule Errors:**
   - Clear error messages untuk business rule violations
   - Suggestions untuk fix errors

## Best Practices

1. **Always Refresh**: Pull down untuk refresh data secara manual
2. **Check Target Progress**: Monitor target progress secara berkala
3. **Use Quick Actions**: Gunakan quick actions untuk akses cepat ke fitur utama
4. **Filter Visits**: Gunakan filter untuk melihat visits berdasarkan status
5. **View All Tasks**: Gunakan "See All" untuk melihat semua tasks jika ada lebih dari 5
6. **Check Theme**: Pastikan theme sesuai dengan preferensi (light/dark)
7. **Check Language**: Pastikan language sesuai dengan preferensi (English/Indonesian)

## Future Enhancements

1. **Period Filter**: Filter data berdasarkan period (today, week, month)
2. **Charts & Analytics**: Visualisasi data dengan charts
3. **Notifications**: Push notifications untuk reminders dan updates
4. **Offline Mode**: Full offline support dengan sync saat online
5. **Customization**: User dapat customize dashboard layout
6. **Widgets**: Additional widgets untuk informasi penting lainnya
