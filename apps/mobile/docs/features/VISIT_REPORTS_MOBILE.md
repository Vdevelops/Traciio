# Mobile Visit Reports - Feature Documentation

## Overview

Mobile Visit Reports adalah fitur untuk sales rep untuk membuat, mengelola, dan melacak kunjungan ke klien. Fitur ini dirancang khusus untuk mobile dengan workflow yang disederhanakan namun tetap komprehensif.

**Note:** Istilah "Visit Reports" tetap menggunakan bahasa Inggris di semua bahasa (tidak diterjemahkan) untuk konsistensi.

## Fitur Utama

### 1. **Create Visit Report**

Sales rep dapat membuat visit report baru dengan memilih salah satu dari tiga tipe:
- **Account & Contact**: Kunjungan ke account dan contact tertentu
- **Deal**: Kunjungan terkait dengan deal tertentu
- **Lead**: Kunjungan untuk follow-up lead

**Form Fields:**
- **Visit Type Selection** (Tab): Account, Deal, atau Lead
- **Account/Deal/Lead Selection** (Required): Pilih salah satu berdasarkan tab aktif dengan **searchable dropdown**
- **Contact Selection** (Optional): Pilih contact jika menggunakan Account atau Deal dengan **searchable dropdown** dan opsi "None"
- **Visit Date** (Required): Tanggal kunjungan (default: hari ini)
- **Visit Time** (Required): Waktu kunjungan (default: waktu saat ini)
- **Purpose** (Required): Tujuan kunjungan (minimal 3 karakter)
- **Notes** (Optional): Catatan tambahan

**UI Features:**
- **Searchable Dropdown**: Semua dropdown selection (Account, Contact, Deal, Lead) memiliki fitur search untuk memudahkan pencarian
- **Tab-based Selection**: Modern tab switch untuk memilih tipe kunjungan
- **Form Validation**: Real-time validation dengan error messages yang jelas
- **Theme Support**: Full support untuk light/dark theme
- **Language Support**: Full support untuk bahasa Indonesia dan Inggris

**API Endpoint:**
- `POST /api/v1/mobile/visit-reports`

**Validations:**
- Account/Deal/Lead harus dipilih sesuai dengan tab aktif
- Purpose wajib diisi dan minimal 3 karakter
- Visit date dan time wajib diisi

### 2. **Update Visit Report**

Sales rep dapat mengupdate visit report yang masih dalam status `draft`. Field yang dapat diupdate:
- Visit Date & Time
- Purpose
- Notes
- Account/Deal/Lead (jika belum check-in)
- Contact (jika belum check-in)

**API Endpoint:**
- `PUT /api/v1/mobile/visit-reports/:id`

**Restrictions:**
- Hanya dapat diupdate jika status = `draft`
- Hanya dapat diupdate oleh owner (sales rep yang membuat)
- Tidak dapat mengubah account/deal/lead setelah check-in

**Validations:**
- Semua validasi sama dengan create visit report
- Ownership validation di backend

### 3. **Delete Visit Report**

Sales rep dapat menghapus visit report yang masih dalam status `draft`.

**API Endpoint:**
- `DELETE /api/v1/mobile/visit-reports/:id`

**Restrictions:**
- Hanya dapat dihapus jika status = `draft`
- Hanya dapat dihapus oleh owner
- Konfirmasi dialog sebelum delete

**Security:**
- Ownership validation di backend
- Soft delete (jika diimplementasikan di backend)

### 4. **Check-In**

Sales rep harus melakukan check-in saat tiba di lokasi kunjungan. Check-in memerlukan:
- **GPS Location**: Lokasi saat ini (wajib)
- **Selfie Photo**: Foto selfie untuk verifikasi (wajib)
- **Fake GPS Detection**: Sistem akan mendeteksi dan memblokir jika menggunakan fake GPS

**API Endpoint:**
- `POST /api/v1/mobile/visit-reports/:id/check-in` (multipart/form-data)

**Features:**
- GPS validation dengan accuracy check
- Fake GPS detection (mendeteksi mock location apps)
- Selfie preview sebelum submit
- Photo metadata (GPS dari foto, timestamp, device GPS)
- Permission handling untuk GPS dan Camera

**Restrictions:**
- Hanya dapat check-in jika status = `draft` dan belum check-in
- Lokasi harus valid (tidak fake GPS)
- Accuracy GPS harus < 100m (recommended)

**Validations:**
- GPS permission harus granted
- Camera permission harus granted
- Location services harus enabled
- Fake GPS detection (block jika terdeteksi)
- Photo harus diambil (tidak boleh dari gallery)

**Security:**
- GPS coordinates validation
- Photo EXIF metadata validation
- Fake GPS detection dengan multiple checks
- Ownership validation

### 5. **Check-Out**

Sales rep dapat melakukan check-out setelah selesai kunjungan. Check-out memerlukan:
- **GPS Location**: Lokasi saat ini (wajib)
- **Selfie Photo**: Foto selfie untuk verifikasi (opsional, seperti versi web)

**API Endpoint:**
- `POST /api/v1/mobile/visit-reports/:id/check-out` (multipart/form-data)

**Features:**
- GPS validation dengan accuracy check
- Fake GPS detection
- Optional photo upload (user dapat memilih dengan atau tanpa foto)
- Selfie preview jika user memilih untuk upload foto

**Restrictions:**
- Hanya dapat check-out jika sudah check-in dan belum check-out
- Lokasi harus valid (tidak fake GPS)

**Validations:**
- GPS permission harus granted
- Location services harus enabled
- Fake GPS detection (block jika terdeteksi)
- Photo opsional (jika dipilih, harus diambil, tidak boleh dari gallery)

**Security:**
- GPS coordinates validation
- Photo EXIF metadata validation (jika ada foto)
- Fake GPS detection
- Ownership validation

### 6. **Submit Visit Report**

Setelah check-in dan check-out, sales rep dapat submit visit report untuk approval. Submit memerlukan:
- **Outcome** (Required): Hasil kunjungan (positive, very_positive, neutral, negative)
- **Next Steps** (Optional): Langkah selanjutnya

**API Endpoint:**
- `POST /api/v1/mobile/visit-reports/:id/submit`

**Restrictions:**
- Hanya dapat submit jika sudah check-in dan check-out
- Status akan berubah menjadi `submitted` dan menunggu approval

**Validations:**
- Outcome wajib dipilih
- Next steps opsional

**Security:**
- Ownership validation
- Status validation (harus sudah check-in dan check-out)

### 7. **View Visit Report Details**

Sales rep dapat melihat detail lengkap visit report, termasuk:
- Informasi account/deal/lead dan contact
- Visit information (date, purpose, notes)
- Check-in/out status dan lokasi
- Photos (termasuk selfie dari check-in/out)
- Outcome dan next steps (jika sudah submit)
- Approval status (submitted, approved, rejected)

**API Endpoint:**
- `GET /api/v1/mobile/visit-reports/:id`

**Features:**
- Offline caching untuk akses cepat
- Photo gallery dengan tap to view
- Map view untuk lokasi check-in/out (jika diimplementasikan)
- Action buttons berdasarkan status (Check In, Check Out, Submit, Update, Delete)

**Security:**
- Ownership validation (hanya owner yang dapat melihat)

### 8. **List Visit Reports**

Sales rep dapat melihat daftar visit reports mereka dengan:
- Pagination (offset-based untuk infinite scroll)
- Filter by status (all, planned, completed, cancelled)
- Filter by account, deal, atau lead
- Search functionality

**API Endpoint:**
- `GET /api/v1/mobile/visit-reports?offset=0&per_page=20&status=...`

**Features:**
- Infinite scroll support dengan offset-based pagination
- Offline caching untuk akses cepat (60 detik TTL)
- Pull-to-refresh untuk update data
- Optimized list rendering dengan `addRepaintBoundaries` dan `cacheExtent`
- Search dengan debounce (500ms)

**UI Features:**
- **Visit Report Card**: Menampilkan purpose sebagai title utama, account/deal/lead info, contact info, dan date/time di bagian bawah
- **Type Badge**: Badge untuk menunjukkan tipe visit (ACCOUNT, DEAL, LEAD)
- **Status Badge**: Badge untuk menunjukkan status (DRAFT, IN_PROGRESS, COMPLETED, APPROVED, REJECTED)
- **Check-in/out Indicators**: Icon untuk menunjukkan status check-in/out
- **Photo Count**: Menampilkan jumlah foto jika ada

**Performance Optimizations:**
- List caching dengan TTL 60 detik
- Lazy loading dengan infinite scroll
- Optimized widget rendering
- Stable keys untuk list items

**Security:**
- Hanya menampilkan visit reports milik user yang login
- Filter di backend berdasarkan `sales_rep_id`

### 9. **Form Data API**

API untuk mendapatkan data selection fields (accounts, contacts, deals, leads) untuk form create/update.

**API Endpoint:**
- `GET /api/v1/mobile/visit-reports/form-data`

**Response:**
- `accounts`: List semua accounts (general, accessible to all users)
- `contacts`: List contacts yang terkait dengan accounts
- `deals`: List semua deals (general, accessible to all users)
- `leads`: List semua leads (general, accessible to all users)

**Note:** Data accounts, deals, dan leads adalah general (tidak difilter by assigned_to) untuk konsistensi dengan versi web.

**Caching:**
- Data form di-cache untuk mengurangi API calls
- Cache di-invalidate saat create/update visit report

## Alur Sales (Sales Flow)

### 1. **Planning Phase (Perencanaan)**

Sales rep merencanakan kunjungan dengan membuat visit report:
1. Buka aplikasi mobile
2. Navigate ke "Visits" menu (atau "Visit Reports" di navbar)
3. Tap "Create Visit" atau "+" button
4. Pilih tipe kunjungan (Account/Deal/Lead) menggunakan tab switch
5. Pilih account/deal/lead yang akan dikunjungi menggunakan searchable dropdown
6. (Optional) Pilih contact jika menggunakan Account/Deal menggunakan searchable dropdown dengan opsi "None"
7. Pilih visit date dan time
8. Isi purpose (wajib, minimal 3 karakter) dan notes (opsional)
9. Tap "Create Visit Report"
10. Visit report dibuat dengan status `draft`

**Status:** `draft`

**UI Flow:**
- Form menggunakan tab-based selection untuk tipe kunjungan
- Semua dropdown menggunakan searchable dropdown dengan fitur search
- Real-time validation dengan error messages
- Theme dan language support penuh

### 2. **Execution Phase (Eksekusi - Check-In)**

Saat tiba di lokasi kunjungan:
1. Buka visit report detail
2. Tap "Check In" button
3. Sistem meminta permission GPS (jika belum granted)
4. Sistem meminta permission camera (jika belum granted)
5. Sistem mendapatkan lokasi GPS saat ini
6. Sistem validasi GPS (deteksi fake GPS)
7. Jika fake GPS terdeteksi, tampilkan warning dan block check-in
8. Ambil foto selfie (wajib)
9. Preview foto dan konfirmasi
10. Submit check-in
11. Status tetap `draft`, tetapi sudah ada check-in time dan location

**Status:** `draft` (dengan check-in data)

**Validations:**
- GPS permission harus granted
- Camera permission harus granted
- Location services harus enabled
- Fake GPS detection
- Photo harus diambil (tidak boleh dari gallery)

### 3. **Detail Filling Phase (Pengisian Detail)**

Setelah check-in, sales rep dapat:
- Update visit report (jika perlu mengubah purpose, notes, atau visit date) - **Note:** Account/Deal/Lead tidak dapat diubah setelah check-in
- Upload additional photos (jika perlu)
- Melakukan aktivitas kunjungan

**Status:** `draft` (dengan check-in data)

**Restrictions:**
- Tidak dapat mengubah account/deal/lead setelah check-in
- Tidak dapat delete setelah check-in

### 4. **Check-Out Phase**

Setelah selesai kunjungan:
1. Buka visit report detail
2. Tap "Check Out" button
3. Sistem meminta permission GPS (jika belum granted)
4. Sistem mendapatkan lokasi GPS saat ini
5. Sistem validasi GPS (deteksi fake GPS)
6. Jika fake GPS terdeteksi, tampilkan warning dan block check-out
7. User memilih apakah ingin upload foto selfie (opsional)
8. Jika memilih upload foto:
   - Sistem meminta permission camera (jika belum granted)
   - Ambil foto selfie
   - Preview foto dan konfirmasi
9. Submit check-out
10. Status tetap `draft`, tetapi sudah ada check-out time dan location

**Status:** `draft` (dengan check-in dan check-out data)

**Validations:**
- GPS permission harus granted
- Location services harus enabled
- Fake GPS detection
- Photo opsional (jika dipilih, harus diambil, tidak boleh dari gallery)

### 5. **Submission Phase (Submit untuk Approval)**

Setelah check-out, sales rep dapat submit visit report:
1. Buka visit report detail
2. Tap "Submit Visit Report" button
3. Isi outcome (wajib): positive, very_positive, neutral, atau negative
4. (Optional) Isi next steps
5. Tap "Submit"
6. Status berubah menjadi `submitted`
7. Visit report menunggu approval dari manager

**Status:** `submitted` (awaiting approval)

**Validations:**
- Outcome wajib dipilih
- Next steps opsional

### 6. **Approval Phase (Manager Review)**

Manager dapat:
- Approve visit report → Status menjadi `approved`
- Reject visit report → Status menjadi `rejected` (dengan rejection reason)

**Status:** `approved` atau `rejected`

**Note:** Approval dilakukan di web version, mobile hanya menampilkan status.

### 7. **Completed Phase**

Setelah approved, visit report selesai dan tidak dapat diubah lagi.

**Status:** `approved` (final)

## Status Flow Diagram

```
[draft] 
  ↓ (check-in)
[draft with check-in]
  ↓ (check-out)
[draft with check-in & check-out]
  ↓ (submit)
[submitted] 
  ↓ (manager approve/reject)
[approved] atau [rejected]
```

## Validations

### Form Validations

#### Create/Update Visit Report
- **Account/Deal/Lead**: Wajib dipilih sesuai dengan tab aktif
- **Contact**: Opsional, dapat memilih "None"
- **Visit Date**: Wajib diisi
- **Visit Time**: Wajib diisi
- **Purpose**: Wajib diisi, minimal 3 karakter
- **Notes**: Opsional

#### Check-In
- **GPS Permission**: Harus granted
- **Camera Permission**: Harus granted
- **Location Services**: Harus enabled
- **GPS Accuracy**: Recommended < 100m
- **Fake GPS Detection**: Block jika terdeteksi
- **Photo**: Wajib diambil (tidak boleh dari gallery)

#### Check-Out
- **GPS Permission**: Harus granted
- **Location Services**: Harus enabled
- **GPS Accuracy**: Recommended < 100m
- **Fake GPS Detection**: Block jika terdeteksi
- **Photo**: Opsional (jika dipilih, harus diambil, tidak boleh dari gallery)

#### Submit
- **Outcome**: Wajib dipilih (positive, very_positive, neutral, negative)
- **Next Steps**: Opsional

### Business Rule Validations

1. **Status-based Actions:**
   - `draft`: Dapat update, delete, check-in
   - `draft (with check-in)`: Dapat update (kecuali account/deal/lead), delete (tidak bisa), check-out, upload photos
   - `draft (with check-in & check-out)`: Dapat update (kecuali account/deal/lead), delete (tidak bisa), submit
   - `submitted`: Tidak dapat diubah (hanya view)
   - `approved`: Tidak dapat diubah (hanya view)
   - `rejected`: Tidak dapat diubah (hanya view)

2. **Ownership:**
   - Semua operasi (view, update, delete, check-in, check-out, submit) hanya dapat dilakukan oleh owner
   - Backend memvalidasi ownership untuk semua operasi

3. **Sequence Validation:**
   - Check-in hanya dapat dilakukan jika status = `draft` dan belum check-in
   - Check-out hanya dapat dilakukan jika sudah check-in dan belum check-out
   - Submit hanya dapat dilakukan jika sudah check-in dan check-out

## Security

### Authentication & Authorization

1. **JWT Token:**
   - Semua API calls menggunakan JWT token untuk authentication
   - Token di-refresh otomatis jika expired
   - Token disimpan securely di device

2. **Ownership Validation:**
   - Backend memvalidasi ownership untuk semua operasi
   - User hanya dapat mengakses visit reports miliknya sendiri
   - API filter berdasarkan `sales_rep_id` dari JWT token

3. **Permission-based Access:**
   - CREATE permission diperlukan untuk membuat visit report
   - Semua operasi lainnya berdasarkan ownership

### Data Security

1. **GPS Data:**
   - GPS coordinates dikirim ke server untuk validasi
   - Fake GPS detection dengan multiple checks:
     - Mock location app detection
     - GPS accuracy validation
     - Photo EXIF metadata validation (jika tersedia)

2. **Photo Security:**
   - Photo di-upload sebagai multipart/form-data
   - Photo metadata (GPS, timestamp) divalidasi
   - Photo tidak dapat diambil dari gallery (harus diambil langsung)

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
   - GPS coordinates tidak disimpan di cache (hanya di server)
   - Photo tidak disimpan di cache (hanya URL)

## Performance Optimizations

### List Performance

1. **Caching:**
   - List data di-cache dengan TTL 60 detik
   - Cache key berdasarkan page, search query, dan filters
   - Cache di-invalidate saat create/update/delete

2. **Lazy Loading:**
   - Infinite scroll dengan offset-based pagination
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
   - Cache di-invalidate saat update/delete

2. **Lazy Loading:**
   - Photos di-load secara lazy
   - Map view di-load secara lazy (jika diimplementasikan)

### Form Performance

1. **Searchable Dropdown:**
   - Search dengan debounce untuk mengurangi API calls
   - Local filtering untuk items yang sudah di-load
   - Dialog-based search untuk better UX

2. **Form Data Caching:**
   - Form data (accounts, contacts, deals, leads) di-cache
   - Cache di-invalidate saat create/update visit report

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

### Language Support

1. **Supported Languages:**
   - English (en)
   - Indonesian (id)

2. **Localization:**
   - Semua text menggunakan `AppLocalizations` (l10n)
   - Date/time formatting menggunakan locale-aware `DateFormat`
   - Number formatting menggunakan locale-aware formatting

3. **Special Cases:**
   - "Visit Reports" tetap menggunakan bahasa Inggris di semua bahasa (tidak diterjemahkan)
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

### GPS & Camera Errors

1. **Permission Errors:**
   - Clear messages untuk permission denials
   - Instructions untuk grant permissions
   - Deep links ke settings (jika supported)

2. **Fake GPS Detection:**
   - Warning dialog dengan clear explanation
   - Instructions untuk disable fake GPS
   - Block check-in/out jika fake GPS terdeteksi

## Best Practices

1. **Always Check GPS**: Pastikan GPS aktif dan akurat sebelum check-in/out
2. **Take Clear Selfies**: Foto selfie harus jelas untuk verifikasi
3. **Fill Purpose Clearly**: Purpose harus jelas dan deskriptif
4. **Add Notes When Needed**: Notes membantu untuk follow-up
5. **Submit Promptly**: Submit visit report segera setelah check-out untuk approval cepat
6. **Use Search**: Gunakan searchable dropdown untuk memudahkan pencarian account/contact/deal/lead
7. **Check Theme**: Pastikan theme sesuai dengan preferensi (light/dark)
8. **Check Language**: Pastikan language sesuai dengan preferensi (English/Indonesian)

## Future Enhancements

1. **Offline Mode**: Full offline support dengan sync saat online
2. **Photo Compression**: Auto-compress photos untuk menghemat bandwidth
3. **Route Integration**: Integrasi dengan route optimization untuk planned visits
4. **Notifications**: Push notifications untuk reminders dan approvals
5. **Analytics**: Dashboard analytics untuk sales rep performance
6. **Map Integration**: Map view untuk lokasi check-in/out
7. **Photo Gallery**: Enhanced photo gallery dengan zoom dan swipe
8. **Export**: Export visit reports ke PDF atau Excel