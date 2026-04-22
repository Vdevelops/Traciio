# Mobile Route Optimization - Feature Documentation

## Overview

Mobile Route Optimization adalah fitur untuk sales rep untuk membuat dan mengelola rute yang dioptimalkan untuk kunjungan ke multiple accounts/contacts. Fitur ini menggunakan algoritma 2-Opt untuk menghasilkan rute terpendek dan paling efisien, dengan dukungan time windows dan priority-based scheduling.

**Note:** Istilah "Route Optimization" tetap menggunakan bahasa Inggris di semua bahasa (tidak diterjemahkan) untuk konsistensi.

## Fitur Utama

### 1. **Create Optimized Route**

Sales rep dapat membuat route yang dioptimalkan dengan memilih start location dan multiple waypoints.

**Form Fields:**
- **Route Name** (Optional): Nama route untuk identifikasi
- **Start Location** (Required): Lokasi awal (dapat menggunakan current GPS location atau manual input)
- **Waypoints** (Required, minimal 1): Daftar tujuan yang akan dikunjungi

**Waypoint Selection:**
- **From Accounts**: Pilih dari daftar accounts dengan:
  - Search functionality untuk mencari accounts
  - GPS location indicator jika account memiliki visit report dengan GPS coordinates
  - Geocoding otomatis untuk accounts tanpa GPS location (menggunakan address)
  - Priority: Visit report location > Geocoded address
- **From Visit Reports**: Pilih dari daftar visit reports yang sudah approved dengan:
  - Filter: Hanya visit reports dengan status `approved` dan memiliki GPS location (check-in atau check-out)
  - Search functionality untuk mencari visit reports
  - GPS coordinates langsung dari visit report (lebih akurat)

**API Endpoint:**
- `POST /api/v1/mobile/route-optimization/optimize`

**Features:**
- **GPS Location**: Menggunakan `LocationAccuracy.bestForNavigation` untuk akurasi maksimal
- **Reverse Geocoding**: Menggunakan Nominatim (OpenStreetMap) untuk mendapatkan address dari GPS coordinates
- **Geocoding**: Menggunakan backend API dengan fallback ke Nominatim untuk mendapatkan coordinates dari address
- **Request Debouncing**: 500ms debounce untuk mencegah duplicate requests
- **Local Caching**: Cache route results untuk instant response (<100ms) untuk same waypoints
- **Retry Logic**: Automatic retry dengan exponential backoff untuk network errors

**Validations:**
- Start location wajib di-set (dapat menggunakan current GPS location)
- Minimal 1 waypoint wajib ditambahkan
- Route name opsional

**Permissions:**
- Memerlukan permission `route-optimization.create`

### 2. **View Route Details**

Sales rep dapat melihat detail lengkap route yang sudah dioptimalkan, termasuk:
- Route name (atau "Unnamed Route" jika tidak ada)
- Route statistics (distance, duration, number of stops)
- Interactive map dengan:
  - Route polyline (jika tersedia)
  - Route steps (jika tersedia)
  - Waypoint markers (start location dan destinations)
  - Optimized order visualization
- Waypoints list dengan:
  - Order number (start location = "S", destinations = 1, 2, 3, ...)
  - Account name (jika terkait dengan account)
  - Address atau GPS coordinates
  - Account badge (jika terkait dengan account)
- Route steps (jika tersedia):
  - Step-by-step navigation instructions
  - Distance dan duration per step

**API Endpoint:**
- `GET /api/v1/mobile/route-optimization/route/:id`

**Features:**
- Offline caching untuk akses cepat
- Interactive map dengan `flutter_map` dan `latlong2`
- Polyline decoding untuk menampilkan route di map
- Route steps visualization
- Waypoint markers dengan color coding (start: green, destinations: primary color)

**Permissions:**
- Memerlukan permission `route-optimization.view`
- Hanya dapat melihat routes milik user yang login

### 3. **List Routes**

Sales rep dapat melihat daftar routes yang sudah dibuat dengan:
- Pagination (page-based untuk infinite scroll)
- Pull-to-refresh untuk update data

**API Endpoint:**
- `GET /api/v1/mobile/route-optimization/my-routes?page=1&per_page=20`

**Features:**
- Infinite scroll support dengan page-based pagination
- Offline caching untuk akses cepat (TTL: 5 minutes)
- Pull-to-refresh untuk update data
- Optimized list rendering dengan `addRepaintBoundaries` dan `cacheExtent`

**UI Features:**
- **Route Card**: Menampilkan route name, creation date, distance, duration, dan number of stops
- **Delete Button**: Delete route button (jika memiliki permission)
- **Empty State**: Message dengan create button jika tidak ada routes

**Performance Optimizations:**
- List caching dengan TTL 5 minutes
- Lazy loading dengan infinite scroll
- Optimized widget rendering
- Stable keys untuk list items

**Permissions:**
- Memerlukan permission `route-optimization.view`
- Hanya menampilkan routes milik user yang login

### 4. **Delete Route**

Sales rep dapat menghapus route yang sudah dibuat.

**API Endpoint:**
- `DELETE /api/v1/mobile/route-optimization/route/:id`

**Restrictions:**
- Hanya dapat dihapus oleh owner (user yang membuat route)
- Konfirmasi dialog sebelum delete

**Security:**
- Ownership validation di backend
- Backend memvalidasi bahwa route belongs to logged-in user

**Permissions:**
- Memerlukan permission `route-optimization.delete`

### 5. **Route Optimization Algorithm**

Backend menggunakan algoritma 2-Opt untuk optimasi route, yang menghasilkan:
- **15-25% shorter routes** dibandingkan dengan urutan asli
- **Distance Matrix Caching**: 80% reduction OSRM calls
- **Parallel OSRM Requests**: 5-10x faster distance matrix calculation
- **Time Windows Support**: Smart scheduling dengan priority-based ordering

**Time Windows Features:**
- **Earliest Arrival**: Waktu tercepat customer tersedia
- **Latest Arrival**: Waktu terakhir customer tersedia
- **Service Duration**: Durasi layanan per waypoint (dalam menit)
- **Priority**: Priority level (1 = highest, 5 = lowest, default: 3)

**Note:** Time windows features sudah didukung di model, tetapi UI untuk input time windows belum diimplementasikan di mobile app.

## Alur Sales (Sales Flow)

### 1. **Create Route Flow**

1. Buka aplikasi mobile
2. Navigate ke "Route" menu (atau "Route Optimization" di navbar)
3. Tap "+" button (FloatingActionButton) jika memiliki permission
4. Set start location:
   - Tap "Use Current Location" button untuk menggunakan GPS location saat ini
   - Sistem meminta permission GPS (jika belum granted)
   - Sistem mendapatkan lokasi GPS dengan akurasi tinggi (`bestForNavigation`)
   - Sistem melakukan reverse geocoding untuk mendapatkan address
   - Start location ditampilkan dengan address dan coordinates
5. Add waypoints:
   - Tap "Add" button untuk membuka waypoint selector dialog
   - Pilih tab "Accounts" atau "Visit Reports"
   - **From Accounts Tab:**
     - Search accounts (jika perlu)
     - Pilih accounts yang ingin dikunjungi
     - Accounts dengan GPS location dari visit report ditandai dengan badge "GPS"
     - Accounts tanpa GPS location akan di-geocode menggunakan address
   - **From Visit Reports Tab:**
     - Hanya menampilkan visit reports dengan status `approved` dan memiliki GPS location
     - Search visit reports (jika perlu)
     - Pilih visit reports yang ingin dikunjungi
   - Tap "Add X waypoints" button
   - Jika ada accounts yang perlu geocoding, loading dialog ditampilkan
   - Waypoints ditambahkan ke form
6. (Optional) Set route name
7. Tap "Optimize" button
8. Sistem mengoptimalkan route:
   - Request di-debounce (500ms) untuk mencegah duplicate calls
   - Check cache terlebih dahulu (jika same waypoints)
   - Jika tidak ada di cache, hit API untuk optimasi
   - Route di-optimize menggunakan 2-Opt algorithm
   - Result di-cache untuk akses cepat di masa depan
9. Route berhasil dibuat dan ditambahkan ke list
10. Navigate kembali ke route list

**Status:** Route created dan saved

### 2. **View Route Details Flow**

1. Dari route list, tap route card
2. Route detail screen ditampilkan dengan:
   - Interactive map dengan route visualization
   - Route statistics (distance, duration, stops)
   - Waypoints list dengan optimized order
   - Route steps (jika tersedia)
3. User dapat:
   - Melihat route di map
   - Melihat waypoints dalam urutan yang dioptimalkan
   - Melihat step-by-step navigation instructions
   - Delete route (jika memiliki permission)

### 3. **Delete Route Flow**

1. Dari route list atau route detail, tap delete button
2. Konfirmasi dialog ditampilkan
3. Tap "Delete" untuk konfirmasi
4. Route berhasil dihapus dan kembali ke route list

### 4. **Waypoint Selection Flow**

#### 4.1. **From Accounts Tab**
1. Buka waypoint selector dialog
2. Pilih tab "Accounts"
3. Search accounts (jika perlu) - search dengan debounce
4. Pilih accounts yang ingin dikunjungi (checkbox)
5. Accounts dengan GPS location dari visit report ditandai dengan badge "GPS" (hijau)
6. Accounts tanpa GPS location akan di-geocode menggunakan address saat confirm
7. Tap "Add X waypoints" button
8. Jika ada accounts yang perlu geocoding:
   - Loading dialog ditampilkan
   - Geocoding dilakukan secara parallel (batch size: 3)
   - Accounts dengan visit report location ditambahkan langsung (tidak perlu geocoding)
   - Accounts tanpa visit report location di-geocode menggunakan backend API atau Nominatim
   - Error messages ditampilkan untuk accounts yang gagal di-geocode
9. Waypoints ditambahkan ke form

#### 4.2. **From Visit Reports Tab**
1. Buka waypoint selector dialog
2. Pilih tab "Visit Reports"
3. Hanya visit reports dengan status `approved` dan memiliki GPS location ditampilkan
4. Search visit reports (jika perlu) - search dengan debounce
5. Pilih visit reports yang ingin dikunjungi (checkbox)
6. Tap "Add X waypoints" button
7. Waypoints ditambahkan langsung (tidak perlu geocoding karena sudah ada GPS coordinates)

## Validations

### Form Validations

#### Create Route
- **Start Location**: Wajib di-set (dapat menggunakan current GPS location)
- **Waypoints**: Wajib minimal 1 waypoint
- **Route Name**: Opsional

#### Waypoint Selection
- **Accounts**: 
  - Account harus memiliki address atau city/province, atau memiliki visit report dengan GPS location
  - Accounts tanpa address dan tanpa visit report location tidak dapat dipilih
- **Visit Reports**:
  - Hanya visit reports dengan status `approved` dan memiliki GPS location (check-in atau check-out) yang dapat dipilih

### Business Rule Validations

1. **Permission-based Actions:**
   - CREATE: Memerlukan permission `route-optimization.create`
   - VIEW: Memerlukan permission `route-optimization.view`
   - DELETE: Memerlukan permission `route-optimization.delete`

2. **Ownership:**
   - Semua operasi (view, delete) hanya dapat dilakukan oleh owner
   - Backend memvalidasi ownership untuk semua operasi
   - API filter berdasarkan `user_id` dari JWT token

3. **Waypoint Requirements:**
   - Waypoint harus memiliki valid coordinates (lat, lng)
   - Waypoint dapat terkait dengan account atau visit report
   - Geocoding harus berhasil untuk accounts tanpa GPS location

4. **Route Optimization:**
   - Minimal 1 waypoint diperlukan untuk optimasi
   - Start location wajib di-set
   - Route optimization memerlukan internet connection

## Security

### Authentication & Authorization

1. **JWT Token:**
   - Semua API calls menggunakan JWT token untuk authentication
   - Token di-refresh otomatis jika expired
   - Token disimpan securely di device

2. **Ownership Validation:**
   - Backend memvalidasi ownership untuk semua operasi
   - User hanya dapat mengakses routes miliknya sendiri
   - API filter berdasarkan `user_id` dari JWT token

3. **Permission-based Access:**
   - CREATE permission diperlukan untuk membuat route
   - VIEW permission diperlukan untuk melihat routes
   - DELETE permission diperlukan untuk menghapus route
   - UI elements (buttons, FAB) hanya ditampilkan jika user memiliki permission yang sesuai

### Data Security

1. **GPS Data:**
   - GPS coordinates dikirim ke server untuk optimasi
   - GPS accuracy check (recommended < 20m untuk navigation)
   - GPS permission handling dengan clear error messages

2. **Geocoding Security:**
   - Geocoding menggunakan backend API dengan fallback ke Nominatim
   - Rate limiting handling untuk Nominatim (429 errors)
   - Retry logic untuk temporary failures

3. **Input Validation:**
   - Semua input divalidasi di client dan server
   - SQL injection prevention (parameterized queries)
   - XSS prevention (output escaping)

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

2. **Sensitive Data:**
   - GPS coordinates tidak disimpan di cache (hanya route results)
   - Route results di-cache untuk akses cepat

## Performance Optimizations

### Route Optimization Performance

1. **Local Caching:**
   - Route results di-cache dengan TTL 1 hour
   - Cache key berdasarkan waypoints (hash of waypoint coordinates)
   - Instant response (<100ms) untuk cached routes
   - Background refresh untuk update cache

2. **Request Debouncing:**
   - 500ms debounce untuk rapid requests
   - Prevents duplicate API calls
   - Shares pending request result

3. **Retry Logic:**
   - Automatic retry untuk network errors
   - Exponential backoff (1s, 2s, ...)
   - Max 2 retries
   - No retry untuk client errors (4xx)

### List Performance

1. **Caching:**
   - List data di-cache dengan TTL 5 minutes
   - Cache key berdasarkan page
   - Cache di-invalidate saat create/delete

2. **Lazy Loading:**
   - Infinite scroll dengan page-based pagination
   - Load more saat scroll mencapai 80% dari bottom
   - Loading indicator saat load more

3. **Widget Optimization:**
   - `addRepaintBoundaries: true` untuk mengurangi repaints
   - `cacheExtent: 500` untuk cache lebih banyak items
   - Stable keys (`ValueKey`) untuk list items
   - `addAutomaticKeepAlives: false` untuk mengurangi memory usage

### Detail Screen Performance

1. **Caching:**
   - Detail data di-cache untuk akses cepat
   - Background refresh untuk update cache

2. **Map Performance:**
   - Lazy loading untuk map tiles
   - Polyline decoding di-background
   - Optimized marker rendering
   - Map bounds calculation dengan padding

### Geocoding Performance

1. **Parallel Processing:**
   - Geocoding dilakukan secara parallel (batch size: 3)
   - Priority: Visit report location > Geocoded address
   - Error handling untuk failed geocoding

2. **Caching:**
   - Geocoded addresses dapat di-cache (jika diimplementasikan)
   - Visit report locations tidak perlu geocoding

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
   - Map markers menggunakan theme-aware colors

### Language Support

1. **Supported Languages:**
   - English (en)
   - Indonesian (id)

2. **Localization:**
   - Semua text menggunakan `AppLocalizations` (l10n)
   - Date/time formatting menggunakan locale-aware `DateFormat`
   - Number formatting menggunakan locale-aware formatting

3. **Special Cases:**
   - "Route Optimization" tetap menggunakan bahasa Inggris di semua bahasa (tidak diterjemahkan)
   - Technical terms tetap menggunakan bahasa Inggris untuk konsistensi

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
   - Error message jika tidak ada cached data

2. **Error Messages:**
   - User-friendly error messages
   - Retry mechanism untuk failed requests
   - Error logging untuk debugging

### GPS Errors

1. **Permission Errors:**
   - Clear messages untuk permission denials
   - Instructions untuk grant permissions
   - Deep links ke settings (jika supported)

2. **Location Service Errors:**
   - Clear messages jika location services disabled
   - Instructions untuk enable location services

3. **GPS Accuracy:**
   - Warning jika GPS accuracy > 20m
   - Attempt to get better GPS fix untuk real devices
   - Debug logging untuk GPS information

### Geocoding Errors

1. **Geocoding Failures:**
   - Clear error messages untuk failed geocoding
   - Suggestions untuk fix errors (verify address, use visit reports)
   - Error messages untuk accounts tanpa address

2. **Rate Limiting:**
   - Handling untuk Nominatim rate limiting (429 errors)
   - Retry dengan exponential backoff
   - Fallback strategies

### Validation Errors

1. **Form Validation:**
   - Real-time validation dengan error messages
   - Field-level error messages
   - Summary error messages

2. **Business Rule Errors:**
   - Clear error messages untuk business rule violations
   - Suggestions untuk fix errors

## Best Practices

1. **Use GPS Location**: Gunakan current GPS location untuk start location yang akurat
2. **Use Visit Reports**: Prioritaskan visit reports dengan GPS location untuk waypoints yang lebih akurat
3. **Verify Addresses**: Pastikan accounts memiliki address yang valid untuk geocoding
4. **Check Cache**: Gunakan cached routes untuk instant response
5. **Use Search**: Gunakan search untuk memudahkan pencarian accounts/visit reports
6. **Check Permissions**: Pastikan user memiliki permission sebelum melakukan action
7. **Check Theme**: Pastikan theme sesuai dengan preferensi (light/dark)
8. **Check Language**: Pastikan language sesuai dengan preferensi (English/Indonesian)
9. **Wait for GPS Fix**: Tunggu beberapa detik untuk GPS mendapatkan fix yang lebih akurat
10. **Use Offline Mode**: Gunakan cached routes saat offline untuk akses cepat

## Future Enhancements

1. **Time Windows UI**: Implementasi UI untuk input time windows (earliest arrival, latest arrival, service duration, priority)
2. **Route Sharing**: Share route dengan sales rep lain
3. **Route Templates**: Save dan reuse route templates
4. **Route History**: History untuk routes yang sudah digunakan
5. **Navigation Integration**: Integrasi dengan navigation apps (Google Maps, Waze)
6. **Real-time Tracking**: Real-time tracking untuk route execution
7. **Route Analytics**: Analytics untuk route performance
8. **Bulk Operations**: Bulk create/delete untuk routes
9. **Route Export**: Export routes ke PDF atau Excel
10. **Offline Route Creation**: Create routes saat offline dengan sync saat online
11. **Route Optimization Settings**: Custom settings untuk optimization algorithm
12. **Multi-vehicle Support**: Support untuk multiple vehicles dengan capacity constraints
13. **Traffic Integration**: Real-time traffic data untuk route optimization
14. **Weather Integration**: Weather data untuk route planning
15. **Route Comparison**: Compare multiple route options
