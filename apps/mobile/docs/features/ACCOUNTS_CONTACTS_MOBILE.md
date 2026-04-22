# Mobile Accounts & Contacts - Feature Documentation

## Overview

Mobile Accounts & Contacts adalah fitur untuk sales rep untuk mengelola accounts (klien/perusahaan) dan contacts (kontak person) dalam sistem CRM. Fitur ini dirancang khusus untuk mobile dengan UI yang intuitif dan performa yang optimal.

**Note:** Istilah "Accounts" dan "Contacts" tetap menggunakan bahasa Inggris di semua bahasa (tidak diterjemahkan) untuk konsistensi.

## Fitur Utama

### 1. **Accounts Management**

#### 1.1. **Create Account**

Sales rep dapat membuat account baru dengan informasi lengkap.

**Form Fields:**
- **Name** (Required): Nama account (minimal 3 karakter)
- **Category** (Required): Kategori account (clinic, pharmacy, hospital) - dropdown dengan data dari API
- **Address** (Required): Alamat lengkap account
- **City** (Required): Kota
- **Province** (Required): Provinsi
- **Phone** (Required): Nomor telepon (10-15 digit)
- **Email** (Required): Email dengan format valid
- **Status** (Required): Status account (active/inactive) - default: active

**API Endpoint:**
- `POST /api/v1/accounts`

**Validations:**
- Name wajib diisi dan minimal 3 karakter
- Category wajib dipilih
- Address, city, province wajib diisi
- Phone wajib diisi dan 10-15 digit
- Email wajib diisi dan format valid
- Status wajib dipilih

**Permissions:**
- Memerlukan permission `accounts.create`

#### 1.2. **Update Account**

Sales rep dapat mengupdate account yang sudah ada.

**API Endpoint:**
- `PUT /api/v1/accounts/:id`

**Restrictions:**
- Hanya dapat diupdate oleh user dengan permission `accounts.edit`
- Semua field dapat diupdate

**Validations:**
- Semua validasi sama dengan create account

**Permissions:**
- Memerlukan permission `accounts.edit`

#### 1.3. **Delete Account**

Sales rep dapat menghapus account.

**API Endpoint:**
- `DELETE /api/v1/accounts/:id`

**Restrictions:**
- Hanya dapat dihapus oleh user dengan permission `accounts.delete`
- Konfirmasi dialog sebelum delete

**Security:**
- Ownership validation di backend (jika diimplementasikan)

**Permissions:**
- Memerlukan permission `accounts.delete`

#### 1.4. **View Account Details**

Sales rep dapat melihat detail lengkap account, termasuk:
- Informasi account (name, category, status)
- Contact information (phone, email)
- Address information (address, city, province)
- Badge status (active/inactive) dengan color coding
- Badge type (clinic/pharmacy/hospital) dengan color coding
- Action buttons (View Contacts, Edit, Delete) berdasarkan permissions

**API Endpoint:**
- `GET /api/v1/accounts/:id`

**Features:**
- Offline caching untuk akses cepat
- Color-coded badges untuk status dan type
- Navigation ke contacts list untuk account tersebut
- Action buttons berdasarkan permissions

**Permissions:**
- Memerlukan permission `accounts.view`

#### 1.5. **List Accounts**

Sales rep dapat melihat daftar accounts dengan:
- Pagination (page-based untuk infinite scroll)
- Search functionality (debounce 500ms)
- Filter by status (active/inactive)
- Filter by category (clinic/pharmacy/hospital)
- Filter by assigned_to

**API Endpoint:**
- `GET /api/v1/accounts?page=1&per_page=20&search=...&status=...&category_id=...&assigned_to=...`

**Features:**
- Infinite scroll support dengan page-based pagination
- Offline caching untuk akses cepat (hanya untuk first page dan no search)
- Pull-to-refresh untuk update data
- Optimized list rendering dengan `addRepaintBoundaries` dan `cacheExtent`
- Search dengan debounce (500ms)

**UI Features:**
- **Account Card**: Menampilkan name, category badge, status badge, phone, email, dan address
- **Status Badge**: Badge dengan color coding (active: green, inactive: gray)
- **Type Badge**: Badge dengan color coding (clinic: blue, pharmacy: orange, hospital: purple)
- **Search Bar**: Search dengan debounce untuk mengurangi API calls

**Performance Optimizations:**
- List caching dengan TTL (hanya untuk first page dan no search)
- Lazy loading dengan infinite scroll
- Optimized widget rendering
- Stable keys untuk list items

**Permissions:**
- Memerlukan permission `accounts.view`

### 2. **Contacts Management**

#### 2.1. **Create Contact**

Sales rep dapat membuat contact baru yang terkait dengan account.

**Form Fields:**
- **Account** (Required): Pilih account yang terkait - dropdown dengan data dari API (disabled saat edit)
- **Name** (Required): Nama contact (minimal 3 karakter)
- **Role** (Required): Role contact (doctor, nurse, pharmacist, manager, staff, dll) - dropdown dengan data dari API
- **Phone** (Required): Nomor telepon (10-15 digit)
- **Email** (Required): Email dengan format valid
- **Position** (Optional): Posisi/jabatan contact
- **Notes** (Optional): Catatan tambahan

**API Endpoint:**
- `POST /api/v1/contacts`

**Validations:**
- Account wajib dipilih
- Name wajib diisi dan minimal 3 karakter
- Role wajib dipilih
- Phone wajib diisi dan 10-15 digit
- Email wajib diisi dan format valid
- Position dan notes opsional

**Permissions:**
- Memerlukan permission `accounts.create` (contacts menggunakan accounts permission)

#### 2.2. **Update Contact**

Sales rep dapat mengupdate contact yang sudah ada.

**API Endpoint:**
- `PUT /api/v1/contacts/:id`

**Restrictions:**
- Hanya dapat diupdate oleh user dengan permission `accounts.edit`
- Account tidak dapat diubah setelah contact dibuat (dropdown disabled)

**Validations:**
- Semua validasi sama dengan create contact
- Account tidak dapat diubah

**Permissions:**
- Memerlukan permission `accounts.edit`

#### 2.3. **Delete Contact**

Sales rep dapat menghapus contact.

**API Endpoint:**
- `DELETE /api/v1/contacts/:id`

**Restrictions:**
- Hanya dapat dihapus oleh user dengan permission `accounts.delete`
- Konfirmasi dialog sebelum delete

**Security:**
- Ownership validation di backend (jika diimplementasikan)

**Permissions:**
- Memerlukan permission `accounts.delete`

#### 2.4. **View Contact Details**

Sales rep dapat melihat detail lengkap contact, termasuk:
- Informasi contact (name, position, role)
- Contact information (phone, email)
- Account information (account name, city)
- Role badge dengan color coding
- Notes (jika ada)
- Action buttons (Edit, Delete) berdasarkan permissions

**API Endpoint:**
- `GET /api/v1/contacts/:id`

**Features:**
- Offline caching untuk akses cepat
- Color-coded role badge (doctor: teal, nurse: pink, pharmacist: orange, manager: purple, staff: blue)
- Account information dengan navigation ke account detail (jika diimplementasikan)
- Action buttons berdasarkan permissions

**Permissions:**
- Memerlukan permission `accounts.view` (contacts menggunakan accounts permission)

#### 2.5. **List Contacts**

Sales rep dapat melihat daftar contacts dengan:
- Pagination (page-based untuk infinite scroll)
- Search functionality (debounce 500ms)
- Filter by account (jika dibuka dari account detail)
- Filter by role

**API Endpoint:**
- `GET /api/v1/contacts?page=1&per_page=20&search=...&account_id=...&role_id=...`

**Features:**
- Infinite scroll support dengan page-based pagination
- Offline caching untuk akses cepat (hanya untuk first page, no search, dan no account filter)
- Pull-to-refresh untuk update data
- Optimized list rendering dengan `addRepaintBoundaries` dan `cacheExtent`
- Search dengan debounce (500ms)

**UI Features:**
- **Contact Card**: Menampilkan name, role badge, phone, email, dan account name
- **Role Badge**: Badge dengan color coding berdasarkan role (doctor: teal, nurse: pink, pharmacist: orange, manager: purple, staff: blue)
- **Search Bar**: Search dengan debounce untuk mengurangi API calls
- **Floating Action Button**: Create contact button (hanya jika dibuka dari account detail dan memiliki permission)

**Performance Optimizations:**
- List caching dengan TTL (hanya untuk first page, no search, dan no account filter)
- Lazy loading dengan infinite scroll
- Optimized widget rendering
- Stable keys untuk list items

**Permissions:**
- Memerlukan permission `accounts.view` (contacts menggunakan accounts permission)

### 3. **Color Coding Badges**

#### 3.1. **Account Status Badge**

**Warna:**
- **Active**: 
  - Light mode: Hijau gelap (`#2E7D32`) dengan background hijau terang (`#4CAF50` dengan opacity 15%)
  - Dark mode: Hijau terang (`#4CAF50`) dengan background hijau dengan opacity 20%
- **Inactive**: 
  - Light mode: Abu-abu gelap (`#616161`) dengan background abu-abu terang (`#9E9E9E` dengan opacity 15%)
  - Dark mode: Abu-abu terang (`#9E9E9E`) dengan background abu-abu dengan opacity 20%

**Lokasi:**
- Account list card
- Account detail screen

#### 3.2. **Account Type Badge**

**Warna:**
- **Clinic**: 
  - Light mode: Biru gelap (`#1976D2`) dengan background biru terang (`#42A5F5` dengan opacity 15%)
  - Dark mode: Biru terang (`#42A5F5`) dengan background biru dengan opacity 20%
- **Pharmacy**: 
  - Light mode: Orange gelap (`#F57C00`) dengan background orange terang (`#FF9800` dengan opacity 15%)
  - Dark mode: Orange terang (`#FF9800`) dengan background orange dengan opacity 20%
- **Hospital**: 
  - Light mode: Ungu gelap (`#5E35B1`) dengan background ungu terang (`#9575CD` dengan opacity 15%)
  - Dark mode: Ungu terang (`#9575CD`) dengan background ungu dengan opacity 20%

**Lokasi:**
- Account list card
- Account detail screen

#### 3.3. **Contact Role Badge**

**Warna:**
- **Doctor/Physician/Dokter**: 
  - Light mode: Teal gelap (`#00897B`) dengan background teal terang (`#26A69A` dengan opacity 15%)
  - Dark mode: Teal terang (`#26A69A`) dengan background teal dengan opacity 20%
- **Nurse/Perawat**: 
  - Light mode: Pink gelap (`#C2185B`) dengan background pink terang (`#EC407A` dengan opacity 15%)
  - Dark mode: Pink terang (`#EC407A`) dengan background pink dengan opacity 20%
- **Pharmacist/Apoteker**: 
  - Light mode: Orange gelap (`#F57C00`) dengan background orange terang (`#FF9800` dengan opacity 15%)
  - Dark mode: Orange terang (`#FF9800`) dengan background orange dengan opacity 20%
- **Manager/Admin/Manajer**: 
  - Light mode: Ungu gelap (`#5E35B1`) dengan background ungu terang (`#9575CD` dengan opacity 15%)
  - Dark mode: Ungu terang (`#9575CD`) dengan background ungu dengan opacity 20%
- **Staff/Employee/Karyawan**: 
  - Light mode: Biru gelap (`#1976D2`) dengan background biru terang (`#42A5F5` dengan opacity 15%)
  - Dark mode: Biru terang (`#42A5F5`) dengan background biru dengan opacity 20%
- **Default**: Primary color dari theme

**Lokasi:**
- Contact list card
- Contact detail screen

**Note:** Color coding menggunakan role name atau code untuk menentukan warna. Jika role tidak cocok dengan kategori di atas, akan menggunakan primary color.

## Alur Sales (Sales Flow)

### 1. **Account Management Flow**

#### 1.1. **Create Account**
1. Buka aplikasi mobile
2. Navigate ke "Accounts" menu (atau "Accounts" tab di Accounts & Contacts screen)
3. Tap "+" button (FloatingActionButton) jika memiliki permission
4. Isi form create account:
   - Name (wajib, minimal 3 karakter)
   - Category (wajib, pilih dari dropdown)
   - Address (wajib)
   - City (wajib)
   - Province (wajib)
   - Phone (wajib, 10-15 digit)
   - Email (wajib, format valid)
   - Status (wajib, default: active)
5. Tap "Save" button
6. Account berhasil dibuat dan ditambahkan ke list

#### 1.2. **View Account Details**
1. Dari account list, tap account card
2. Account detail screen ditampilkan dengan:
   - Account information (name, category, status)
   - Contact information (phone, email)
   - Address information (address, city, province)
   - Action buttons (View Contacts, Edit, Delete) berdasarkan permissions

#### 1.3. **Update Account**
1. Dari account detail screen, tap "Edit" button (jika memiliki permission)
2. Form edit account ditampilkan dengan data yang sudah terisi
3. Update field yang diperlukan
4. Tap "Save" button
5. Account berhasil diupdate dan detail screen di-refresh

#### 1.4. **Delete Account**
1. Dari account detail screen, tap "Delete" button (jika memiliki permission)
2. Konfirmasi dialog ditampilkan
3. Tap "Delete" untuk konfirmasi
4. Account berhasil dihapus dan kembali ke account list

#### 1.5. **View Contacts dari Account**
1. Dari account detail screen, tap "View Contacts" button
2. Contact list screen ditampilkan dengan filter by account
3. List hanya menampilkan contacts yang terkait dengan account tersebut
4. FloatingActionButton untuk create contact tersedia (jika memiliki permission)

### 2. **Contact Management Flow**

#### 2.1. **Create Contact**
1. Buka aplikasi mobile
2. Navigate ke "Contacts" menu (atau "Contacts" tab di Accounts & Contacts screen) atau dari account detail screen
3. Tap "+" button (FloatingActionButton) jika dibuka dari account detail dan memiliki permission
4. Isi form create contact:
   - Account (wajib, pilih dari dropdown, atau sudah terisi jika dibuka dari account detail)
   - Name (wajib, minimal 3 karakter)
   - Role (wajib, pilih dari dropdown)
   - Phone (wajib, 10-15 digit)
   - Email (wajib, format valid)
   - Position (opsional)
   - Notes (opsional)
5. Tap "Save" button
6. Contact berhasil dibuat dan ditambahkan ke list

#### 2.2. **View Contact Details**
1. Dari contact list, tap contact card
2. Contact detail screen ditampilkan dengan:
   - Contact information (name, position, role)
   - Contact details (phone, email)
   - Account information (account name, city)
   - Notes (jika ada)
   - Action buttons (Edit, Delete) berdasarkan permissions

#### 2.3. **Update Contact**
1. Dari contact detail screen, tap "Edit" button (jika memiliki permission)
2. Form edit contact ditampilkan dengan data yang sudah terisi
3. Update field yang diperlukan (account tidak dapat diubah)
4. Tap "Save" button
5. Contact berhasil diupdate dan detail screen di-refresh

#### 2.4. **Delete Contact**
1. Dari contact detail screen, tap "Delete" button (jika memiliki permission)
2. Konfirmasi dialog ditampilkan
3. Tap "Delete" untuk konfirmasi
4. Contact berhasil dihapus dan kembali ke contact list (atau account detail jika dibuka dari account)

### 3. **Search & Filter Flow**

#### 3.1. **Search Accounts**
1. Dari account list screen, ketik di search bar
2. Search dengan debounce 500ms
3. List di-update dengan hasil search
4. Search dilakukan di backend dengan full-text search

#### 3.2. **Search Contacts**
1. Dari contact list screen, ketik di search bar
2. Search dengan debounce 500ms
3. List di-update dengan hasil search
4. Search dilakukan di backend dengan full-text search

#### 3.3. **Filter by Account (Contacts)**
1. Dari account detail screen, tap "View Contacts" button
2. Contact list screen ditampilkan dengan filter by account
3. List hanya menampilkan contacts yang terkait dengan account tersebut

## Validations

### Form Validations

#### Create/Update Account
- **Name**: Wajib diisi, minimal 3 karakter
- **Category**: Wajib dipilih dari dropdown
- **Address**: Wajib diisi
- **City**: Wajib diisi
- **Province**: Wajib diisi
- **Phone**: Wajib diisi, 10-15 digit (hanya angka)
- **Email**: Wajib diisi, format valid (regex: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
- **Status**: Wajib dipilih (active/inactive)

#### Create/Update Contact
- **Account**: Wajib dipilih dari dropdown (tidak dapat diubah saat edit)
- **Name**: Wajib diisi, minimal 3 karakter
- **Role**: Wajib dipilih dari dropdown
- **Phone**: Wajib diisi, 10-15 digit (hanya angka)
- **Email**: Wajib diisi, format valid (regex: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
- **Position**: Opsional
- **Notes**: Opsional

### Business Rule Validations

1. **Permission-based Actions:**
   - CREATE: Memerlukan permission `accounts.create`
   - UPDATE: Memerlukan permission `accounts.edit`
   - DELETE: Memerlukan permission `accounts.delete`
   - VIEW: Memerlukan permission `accounts.view`
   - **Note:** Contacts menggunakan permission `accounts.*` (jika user bisa create accounts, mereka juga bisa create contacts)

2. **Account-Contact Relationship:**
   - Contact harus terkait dengan account
   - Account tidak dapat diubah setelah contact dibuat
   - Jika account dihapus, contacts terkait juga dihapus (jika diimplementasikan di backend)

3. **Data Consistency:**
   - Category harus valid (dari API)
   - Role harus valid (dari API)
   - Account harus valid (dari API)

## Security

### Authentication & Authorization

1. **JWT Token:**
   - Semua API calls menggunakan JWT token untuk authentication
   - Token di-refresh otomatis jika expired
   - Token disimpan securely di device

2. **Permission-based Access:**
   - **Accounts**: Memerlukan permission `accounts.view`, `accounts.create`, `accounts.edit`, `accounts.delete`
   - **Contacts**: Menggunakan permission `accounts.*` (jika user bisa manage accounts, mereka juga bisa manage contacts)
   - UI elements (buttons, FAB) hanya ditampilkan jika user memiliki permission yang sesuai
   - Backend memvalidasi permissions untuk semua operasi

3. **Route Protection:**
   - Routes dilindungi dengan `AuthGate` widget
   - Contacts route menggunakan custom permission `accounts.view` (bukan `contacts.view`)
   - Unauthorized access akan redirect ke login atau show error

### Data Security

1. **Input Validation:**
   - Semua input divalidasi di client dan server
   - SQL injection prevention (parameterized queries)
   - XSS prevention (output escaping)

2. **Data Filtering:**
   - List hanya menampilkan data yang diizinkan untuk user
   - Backend memfilter data berdasarkan permissions dan ownership (jika diimplementasikan)

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
   - Phone numbers dan emails tidak disimpan di cache (hanya di server)
   - Cache hanya untuk list data, tidak untuk detail data sensitif

## Performance Optimizations

### List Performance

1. **Caching:**
   - List data di-cache dengan TTL (hanya untuk first page dan no search/filters)
   - Cache key berdasarkan page, search query, dan filters
   - Cache di-invalidate saat create/update/delete

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
   - Cache di-invalidate saat update/delete
   - Background refresh untuk update cache

2. **Lazy Loading:**
   - Data di-load secara lazy
   - Skeleton screens saat loading

### Form Performance

1. **Dropdown Data:**
   - Categories dan roles di-load sekali dan di-cache
   - Data di-load saat form dibuka
   - Loading indicator saat data di-load

2. **Search Debouncing:**
   - Search dengan debounce 500ms untuk mengurangi API calls
   - Local filtering untuk items yang sudah di-load (jika memungkinkan)

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
   - Badges menggunakan color coding yang kontras untuk light dan dark theme

### Language Support

1. **Supported Languages:**
   - English (en)
   - Indonesian (id)

2. **Localization:**
   - Semua text menggunakan `AppLocalizations` (l10n)
   - Date/time formatting menggunakan locale-aware `DateFormat`
   - Number formatting menggunakan locale-aware formatting

3. **Special Cases:**
   - "Accounts" dan "Contacts" tetap menggunakan bahasa Inggris di semua bahasa (tidak diterjemahkan)
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

### Permission Errors

1. **Permission Denials:**
   - Clear messages untuk permission denials
   - UI elements (buttons, FAB) disembunyikan jika tidak memiliki permission
   - Error messages untuk unauthorized actions

## Best Practices

1. **Always Validate Input**: Pastikan semua input divalidasi sebelum submit
2. **Use Search**: Gunakan search untuk memudahkan pencarian account/contact
3. **Check Permissions**: Pastikan user memiliki permission sebelum melakukan action
4. **Use Color Coding**: Gunakan color coding badges untuk memudahkan identifikasi status dan type
5. **Fill Required Fields**: Pastikan semua required fields diisi dengan benar
6. **Check Theme**: Pastikan theme sesuai dengan preferensi (light/dark)
7. **Check Language**: Pastikan language sesuai dengan preferensi (English/Indonesian)
8. **Use Offline Mode**: Gunakan cached data saat offline untuk akses cepat

## Future Enhancements

1. **Offline Mode**: Full offline support dengan sync saat online
2. **Bulk Operations**: Bulk create/update/delete untuk accounts dan contacts
3. **Advanced Filters**: Filter by multiple criteria (status, category, role, dll)
4. **Export**: Export accounts dan contacts ke PDF atau Excel
5. **Import**: Import accounts dan contacts dari CSV atau Excel
6. **Map Integration**: Map view untuk lokasi accounts
7. **Photo Support**: Photo untuk accounts dan contacts
8. **Notes History**: History notes untuk accounts dan contacts
9. **Activity Log**: Activity log untuk tracking changes
10. **Custom Fields**: Custom fields untuk accounts dan contacts
11. **Tags**: Tags untuk categorizing accounts dan contacts
12. **Relationships**: Relationships antara accounts dan contacts
13. **Analytics**: Analytics untuk accounts dan contacts performance
14. **Notifications**: Push notifications untuk important updates
15. **Searchable Dropdown**: Searchable dropdown untuk account dan role selection (seperti di visit reports)
