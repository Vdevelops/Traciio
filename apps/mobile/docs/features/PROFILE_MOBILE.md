# Mobile Profile - Feature Documentation

## Overview

Mobile Profile adalah fitur untuk sales rep untuk melihat dan mengelola informasi profil mereka, termasuk pengaturan aplikasi seperti theme dan bahasa. Fitur ini dirancang khusus untuk mobile dengan UI yang modern dan intuitif.

**Note:** Istilah "Profile" tetap menggunakan bahasa Inggris di semua bahasa (tidak diterjemahkan) untuk konsistensi.

## Fitur Utama

### 1. **View Profile**

Sales rep dapat melihat informasi profil mereka yang lengkap.

**Informasi yang Ditampilkan:**
- **Avatar**: Foto profil user (dengan fallback ke initial letter jika tidak ada avatar)
- **Name**: Nama lengkap user
- **Email**: Email user
- **Role**: Role user dengan badge styling
- **Status**: Status user (active/inactive)

**UI Features:**
- **Avatar Support**: 
  - Support untuk SVG images (dicebear avatars)
  - Support untuk regular images (JPG, PNG)
  - Fallback ke gradient background dengan initial letter jika avatar tidak tersedia
- **Modern Card Design**: Profile header dengan card design yang modern
- **Theme Support**: Full support untuk light/dark theme
- **Language Support**: Full support untuk bahasa Indonesia dan Inggris

**API Endpoint:**
- `GET /api/v1/auth/mobile/profile`

**Data Source:**
- Data diambil dari API dengan fallback ke auth state jika API gagal
- Avatar URL diambil dari profile response atau auth state

### 2. **Edit Profile**

Sales rep dapat mengupdate nama mereka melalui dialog form.

**Form Fields:**
- **Name** (Required): Nama lengkap user (minimal 1 karakter)

**UI Features:**
- **Dialog Form**: Modern dialog dengan rounded corners
- **Real-time Validation**: Validasi saat submit
- **Error Handling**: Error messages yang jelas
- **Success Feedback**: Snackbar notification setelah berhasil update
- **Auto Refresh**: Profile dan auth state otomatis di-refresh setelah update

**API Endpoint:**
- `PUT /api/v1/auth/mobile/profile`

**Validations:**
- Name wajib diisi (tidak boleh kosong)

**Permissions:**
- Tidak ada permission khusus yang diperlukan (semua user dapat edit profile mereka sendiri)

### 3. **Change Password**

Sales rep dapat mengubah password mereka melalui dialog form (hanya jika memiliki permission).

**Form Fields:**
- **Current Password** (Required): Password saat ini
- **New Password** (Required): Password baru
- **Confirm Password** (Required): Konfirmasi password baru

**UI Features:**
- **Password Visibility Toggle**: Icon untuk show/hide password di setiap field
- **Dialog Form**: Modern dialog dengan rounded corners
- **Real-time Validation**: Validasi saat submit
- **Error Handling**: Error messages yang jelas
- **Success Feedback**: Snackbar notification setelah berhasil change password

**API Endpoint:**
- `PUT /api/v1/auth/mobile/password`

**Validations:**
- Semua field wajib diisi
- New password dan confirm password harus sama
- Password validation sesuai dengan backend requirements

**Permissions:**
- `profile.change-password` - Required untuk mengakses fitur change password
- Jika user tidak memiliki permission, menu "Change Password" tidak akan ditampilkan

### 4. **Theme Settings**

Sales rep dapat mengubah theme aplikasi (Light, Dark, atau System).

**Theme Options:**
- **Light Theme**: Theme terang
- **Dark Theme**: Theme gelap
- **System Theme**: Mengikuti pengaturan sistem device

**UI Features:**
- **Bottom Sheet Modal**: Modern bottom sheet untuk memilih theme
- **Visual Feedback**: Icon dan label yang jelas untuk setiap theme option
- **Selected Indicator**: Check icon untuk theme yang sedang aktif
- **Instant Apply**: Theme langsung diterapkan setelah dipilih
- **Persistent Storage**: Theme preference disimpan dan di-load saat app restart

**Implementation:**
- Menggunakan `ThemeModeProvider` untuk state management
- Theme preference disimpan di local storage
- Theme diterapkan secara global di seluruh aplikasi

### 5. **Language Settings**

Sales rep dapat mengubah bahasa aplikasi (English atau Indonesian).

**Language Options:**
- **English (🇬🇧)**: Bahasa Inggris
- **Indonesian (🇮🇩)**: Bahasa Indonesia

**UI Features:**
- **Bottom Sheet Modal**: Modern bottom sheet untuk memilih bahasa
- **Flag Emoji**: Flag emoji untuk visual identification
- **Selected Indicator**: Check icon untuk bahasa yang sedang aktif
- **Instant Apply**: Bahasa langsung diterapkan setelah dipilih
- **Persistent Storage**: Language preference disimpan dan di-load saat app restart

**Implementation:**
- Menggunakan `LocaleProvider` untuk state management
- Language preference disimpan di local storage
- Language diterapkan secara global di seluruh aplikasi

### 6. **Notifications Settings**

Sales rep dapat mengakses pengaturan notifikasi (placeholder untuk fitur future).

**UI Features:**
- **Settings Tile**: Menu item untuk notifications settings
- **Placeholder**: Saat ini hanya placeholder, akan diimplementasikan di future

### 7. **About Section**

Sales rep dapat melihat informasi tentang aplikasi.

**Informasi yang Ditampilkan:**
- **App Version**: Versi aplikasi saat ini
- **Privacy Policy**: Link ke privacy policy (placeholder)
- **Terms of Service**: Link ke terms of service (placeholder)

**UI Features:**
- **App Version**: Menampilkan versi dari `AppInfo.versionString`
- **Info Icons**: Icon yang sesuai untuk setiap informasi
- **Non-interactive**: App version tidak dapat di-click
- **Interactive**: Privacy Policy dan Terms of Service dapat di-click (placeholder)

### 8. **Logout**

Sales rep dapat logout dari aplikasi.

**UI Features:**
- **Confirmation Dialog**: Dialog konfirmasi sebelum logout
- **Red Styling**: Button dengan warna merah untuk menunjukkan action yang destructive
- **Clear Auth State**: Auth state di-clear setelah logout
- **Navigation**: User diarahkan ke login screen setelah logout

**Implementation:**
- Menggunakan `authProvider.notifier.logout()` untuk logout
- Auth state dan storage di-clear
- User di-redirect ke login screen

## Alur Sales (Sales Flow)

### 1. **View Profile Flow**

```
1. User membuka Profile screen dari bottom navigation
2. System menampilkan loading indicator
3. System fetch profile data dari API
4. System menampilkan profile information:
   - Avatar (dengan fallback jika tidak ada)
   - Name, Email, Role
5. System menampilkan settings options
6. System menampilkan about section
7. System menampilkan logout button
```

### 2. **Edit Profile Flow**

```
1. User tap "Edit Profile" menu
2. System menampilkan dialog form dengan current name
3. User mengubah name
4. User tap "Save"
5. System validasi name (tidak boleh kosong)
6. System mengirim request ke API
7. System update auth state dengan data baru
8. System refresh profile data
9. System menampilkan success message
10. System close dialog
```

### 3. **Change Password Flow**

```
1. User tap "Change Password" menu (jika memiliki permission)
2. System menampilkan dialog form dengan 3 password fields
3. User mengisi current password, new password, dan confirm password
4. User dapat toggle visibility untuk setiap password field
5. User tap "Save"
6. System validasi:
   - Semua field wajib diisi
   - New password dan confirm password harus sama
7. System mengirim request ke API
8. System menampilkan success message
9. System close dialog
```

### 4. **Theme Change Flow**

```
1. User tap "Theme" menu
2. System menampilkan bottom sheet dengan theme options
3. User memilih theme (Light, Dark, atau System)
4. System apply theme secara instant
5. System save theme preference
6. System close bottom sheet
7. Theme tetap tersimpan setelah app restart
```

### 5. **Language Change Flow**

```
1. User tap "Language" menu
2. System menampilkan bottom sheet dengan language options
3. User memilih language (English atau Indonesian)
4. System apply language secara instant
5. System save language preference
6. System close bottom sheet
7. Language tetap tersimpan setelah app restart
```

### 6. **Logout Flow**

```
1. User tap "Logout" button
2. System menampilkan confirmation dialog
3. User confirm logout
4. System clear auth state
5. System clear local storage
6. System redirect ke login screen
```

## Validations

### 1. **Edit Profile Validations**

- **Name**: 
  - Required (tidak boleh kosong)
  - Minimal 1 karakter
  - Error message: "Name is required"

### 2. **Change Password Validations**

- **Current Password**: 
  - Required (tidak boleh kosong)
  - Error message: "All fields are required"
- **New Password**: 
  - Required (tidak boleh kosong)
  - Error message: "All fields are required"
- **Confirm Password**: 
  - Required (tidak boleh kosong)
  - Error message: "All fields are required"
  - Harus sama dengan new password
  - Error message: "Passwords do not match"

### 3. **Backend Validations**

Backend akan melakukan validasi tambahan:
- Password strength requirements
- Current password verification
- Email format (jika email dapat diubah di future)

## Security

### 1. **Authentication**

- **Required**: User harus authenticated untuk mengakses profile screen
- **Auth State Check**: System check auth state sebelum fetch profile
- **Token Validation**: API requests menggunakan JWT token dari auth state

### 2. **Permissions**

- **Edit Profile**: Tidak ada permission khusus (semua user dapat edit profile mereka sendiri)
- **Change Password**: 
  - Permission: `profile.change-password`
  - Menu "Change Password" hanya ditampilkan jika user memiliki permission
  - API endpoint juga melakukan permission check di backend

### 3. **Data Privacy**

- **Profile Data**: Hanya menampilkan data user yang sedang login
- **Password**: Password tidak pernah ditampilkan atau disimpan di client
- **Auth State**: Auth state di-clear saat logout

### 4. **API Security**

- **HTTPS**: Semua API requests menggunakan HTTPS
- **JWT Token**: Semua API requests menggunakan JWT token untuk authentication
- **Error Handling**: Error messages tidak mengekspos informasi sensitif

## Performance Optimizations

### 1. **Data Fetching**

- **Auto-dispose Provider**: `profileProvider` menggunakan `FutureProvider.autoDispose` untuk auto cleanup
- **Caching**: Profile data di-cache di auth state sebagai fallback
- **Error Recovery**: System fallback ke auth state jika API gagal

### 2. **UI Performance**

- **Lazy Loading**: Profile screen menggunakan lazy loading untuk data
- **Optimized Images**: Avatar images menggunakan optimized loading dengan placeholder
- **SVG Support**: Support untuk SVG images dengan proper rendering

### 3. **State Management**

- **Riverpod**: Menggunakan Riverpod untuk efficient state management
- **Provider Reuse**: Reuse providers untuk avoid duplicate API calls
- **Auto Refresh**: Profile di-refresh setelah update untuk ensure data consistency

## Theme & Language Support

### 1. **Theme Support**

- **Light Mode**: Full support untuk light theme
- **Dark Mode**: Full support untuk dark theme
- **System Theme**: Support untuk mengikuti system theme
- **Theme Persistence**: Theme preference disimpan dan di-load saat app restart
- **Instant Apply**: Theme langsung diterapkan tanpa perlu restart app

### 2. **Language Support**

- **English**: Full support untuk bahasa Inggris
- **Indonesian**: Full support untuk bahasa Indonesia
- **Language Persistence**: Language preference disimpan dan di-load saat app restart
- **Instant Apply**: Language langsung diterapkan tanpa perlu restart app
- **Localized Strings**: Semua strings menggunakan `AppLocalizations`

### 3. **Localized Content**

**English:**
- Settings, Edit Profile, Change Password, Notifications, Language, Theme, About, App Version, Privacy Policy, Terms of Service, Logout

**Indonesian:**
- Pengaturan, Edit Profil, Ubah Kata Sandi, Notifikasi, Bahasa, Tema, Tentang, Versi Aplikasi, Kebijakan Privasi, Syarat Layanan, Keluar

## Error Handling

### 1. **API Errors**

- **Network Errors**: Menampilkan error message yang user-friendly
- **404 Errors**: Menampilkan error message khusus untuk endpoint not found
- **Validation Errors**: Menampilkan field-specific error messages
- **Server Errors**: Menampilkan generic error message untuk server errors

### 2. **Error Display**

- **Error Widget**: Menggunakan `ErrorStateWidget` untuk display error dengan retry button
- **Snackbar**: Menggunakan Snackbar untuk error messages di forms
- **Error Extraction**: System extract error message dari API response dengan proper handling

### 3. **Error Recovery**

- **Retry Mechanism**: User dapat retry fetch profile jika terjadi error
- **Fallback Data**: System fallback ke auth state jika API gagal
- **Graceful Degradation**: System tetap menampilkan UI meskipun beberapa data tidak tersedia

## Best Practices

### 1. **Code Organization**

- **Separation of Concerns**: 
  - Repository layer untuk API calls
  - Provider layer untuk state management
  - Presentation layer untuk UI
- **Reusable Components**: 
  - `_SettingsCard`, `_SettingsTile`, `_SectionTitle` untuk reusable UI components
  - `_ThemeSelectorModal`, `_LanguageSelectorModal` untuk reusable modals

### 2. **State Management**

- **Riverpod Providers**: Menggunakan Riverpod untuk efficient state management
- **Auto-dispose**: Menggunakan `autoDispose` untuk auto cleanup
- **Provider Families**: Menggunakan `FutureProvider.family` untuk parameterized providers

### 3. **UI/UX**

- **Loading States**: Menampilkan loading indicator saat fetch data
- **Error States**: Menampilkan error widget dengan retry option
- **Success Feedback**: Menampilkan success message setelah action berhasil
- **Confirmation Dialogs**: Menggunakan confirmation dialog untuk destructive actions

### 4. **Security**

- **Permission Checks**: Selalu check permission sebelum menampilkan menu atau melakukan action
- **Input Validation**: Validasi input di client dan server
- **Error Messages**: Tidak mengekspos informasi sensitif di error messages

## Future Enhancements

### 1. **Profile Features**

- **Avatar Upload**: Allow user untuk upload avatar mereka sendiri
- **Email Update**: Allow user untuk mengubah email (dengan verification)
- **Phone Number**: Add phone number field ke profile
- **Bio/Description**: Add bio/description field ke profile

### 2. **Settings Features**

- **Notifications Settings**: Implement full notifications settings page
- **Privacy Settings**: Add privacy settings (data sharing, etc.)
- **Account Settings**: Add account settings (delete account, etc.)

### 3. **About Features**

- **Privacy Policy**: Implement actual privacy policy page
- **Terms of Service**: Implement actual terms of service page
- **Help & Support**: Add help & support section
- **Feedback**: Add feedback form untuk user feedback

### 4. **Performance**

- **Image Caching**: Implement proper image caching untuk avatars
- **Offline Support**: Add offline support untuk profile data
- **Background Sync**: Sync profile data di background

### 5. **Security**

- **Two-Factor Authentication**: Add 2FA support
- **Session Management**: Add session management (active sessions, logout from all devices)
- **Security Log**: Add security log untuk track security events

## API Endpoints

### 1. **Get Profile**

```
GET /api/v1/auth/mobile/profile

Response:
{
  "success": true,
  "data": {
    "user": {
      "id": "string",
      "email": "string",
      "name": "string",
      "avatar_url": "string | null",
      "role_id": "string | null",
      "role": {
        "id": "string",
        "name": "string",
        "code": "string"
      },
      "status": "string",
      "created_at": "datetime",
      "updated_at": "datetime"
    },
    "stats": {
      "visits": "int | null",
      "deals": "int | null",
      "tasks": "int | null"
    },
    "activities": [...],
    "transactions": [...]
  }
}
```

### 2. **Update Profile**

```
PUT /api/v1/auth/mobile/profile

Request:
{
  "name": "string"
}

Response:
{
  "success": true,
  "data": {
    "id": "string",
    "email": "string",
    "name": "string",
    "avatar_url": "string | null",
    "role_id": "string | null",
    "role": {...},
    "status": "string",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
}
```

### 3. **Change Password**

```
PUT /api/v1/auth/mobile/password

Request:
{
  "current_password": "string",
  "password": "string",
  "confirm_password": "string"
}

Response:
{
  "success": true,
  "data": null
}
```

## Summary

Mobile Profile adalah fitur penting untuk user experience, memberikan akses ke:
- **Profile Information**: View dan edit profile information
- **Security**: Change password dengan permission-based access
- **Settings**: Theme dan language settings dengan persistent storage
- **About**: App version dan legal information
- **Logout**: Secure logout dengan confirmation

Fitur ini dirancang dengan fokus pada:
- **User Experience**: Modern UI dengan smooth interactions
- **Security**: Permission-based access dan secure API calls
- **Performance**: Efficient data fetching dan state management
- **Accessibility**: Theme dan language support untuk semua users
