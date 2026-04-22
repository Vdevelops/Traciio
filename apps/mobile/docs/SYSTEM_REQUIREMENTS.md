# System Requirements - CRM Healthcare Mobile App

**Version**: 1.0  
**Last Updated**: 2025-01-15  
**App Name**: SalesView - CRM Healthcare Mobile Application

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Minimum System Requirements](#minimum-system-requirements)
3. [Recommended System Requirements](#recommended-system-requirements)
4. [Platform-Specific Requirements](#platform-specific-requirements)
5. [Storage Requirements](#storage-requirements)
6. [Network Requirements](#network-requirements)
7. [Rekomendasi Spesifikasi HP](#rekomendasi-spesifikasi-hp)
8. [Pertimbangan Jangka Panjang](#pertimbangan-jangka-panjang)
9. [Testing & Compatibility](#testing--compatibility)

---

## Overview

Aplikasi **SalesView** adalah aplikasi mobile CRM untuk Sales Representative di industri healthcare/pharmaceutical. Aplikasi ini menggunakan **Flutter** dan memiliki fitur-fitur yang memerlukan resource cukup signifikan:

### Fitur-Fitur Utama yang Memerlukan Resource:

1. **GPS Tracking & Location Services**
   - Real-time location tracking untuk visit reports
   - Route optimization dengan multiple waypoints
   - Maps rendering dengan tile layers

2. **Offline Storage & Caching**
   - Hive database untuk offline data
   - Route optimization cache
   - Image caching
   - Dashboard data cache

3. **Image Processing**
   - Camera capture untuk selfie check-in
   - Image upload dengan multipart form-data
   - Image preview dan display

4. **Maps & Charts**
   - Interactive maps dengan markers dan polylines
   - Charts untuk dashboard statistics
   - Real-time data visualization

5. **Background Operations**
   - Data synchronization
   - Offline queue processing
   - Connectivity monitoring

---

## Minimum System Requirements

### Android

| Komponen       | Spesifikasi Minimal                                  |
| -------------- | ---------------------------------------------------- |
| **OS Version** | Android 5.0 (API Level 21) - Lollipop                |
| **RAM**        | 2 GB                                                 |
| **Storage**    | 500 MB (untuk instalasi) + 200 MB (untuk data cache) |
| **CPU**        | ARMv7 atau x86 (32-bit/64-bit)                       |
| **GPU**        | OpenGL ES 2.0 compatible                             |
| **GPS**        | A-GPS support (untuk location tracking)              |
| **Camera**     | 2 MP (untuk selfie check-in)                         |
| **Network**    | 3G/WiFi (untuk sync data)                            |

**Catatan Penting:**

- ✅ Aplikasi dapat berjalan di Android 5.0+, namun performa optimal di Android 8.0+
- ⚠️ Android 5.0-7.1 (API 21-25) mungkin mengalami performa lebih lambat untuk fitur maps dan charts
- ⚠️ RAM 2 GB adalah minimum mutlak; dengan RAM kurang dari ini, aplikasi mungkin crash saat menggunakan maps atau processing images

### iOS

| Komponen       | Spesifikasi Minimal                                  |
| -------------- | ---------------------------------------------------- |
| **OS Version** | iOS 12.0+                                            |
| **Device**     | iPhone 6s atau lebih baru                            |
| **RAM**        | 2 GB                                                 |
| **Storage**    | 500 MB (untuk instalasi) + 200 MB (untuk data cache) |
| **GPS**        | Built-in GPS (semua iPhone 6s+ memiliki GPS)         |
| **Camera**     | 5 MP (untuk selfie check-in)                         |
| **Network**    | 3G/WiFi (untuk sync data)                            |

**Catatan Penting:**

- ✅ iPhone 6s (2015) adalah device tertua yang didukung
- ⚠️ iPhone 6s mungkin mengalami performa lebih lambat untuk fitur maps

---

## Recommended System Requirements

### Android (Rekomendasi untuk Penggunaan Optimal)

| Komponen       | Spesifikasi Rekomendasi                               |
| -------------- | ----------------------------------------------------- |
| **OS Version** | Android 10.0 (API Level 29) atau lebih baru           |
| **RAM**        | 4 GB atau lebih                                       |
| **Storage**    | 1 GB free space (untuk instalasi + cache + updates)   |
| **CPU**        | Octa-core 1.8 GHz atau lebih cepat                    |
| **GPU**        | Adreno 530 / Mali-G71 / PowerVR atau lebih baru       |
| **GPS**        | A-GPS + GLONASS (untuk akurasi lebih baik)            |
| **Camera**     | 8 MP atau lebih (untuk selfie berkualitas baik)       |
| **Network**    | 4G LTE / WiFi (untuk sync cepat)                      |
| **Screen**     | 5.5" atau lebih besar (untuk UX maps yang lebih baik) |

**Keuntungan dengan Spesifikasi Rekomendasi:**

- ✅ Maps loading lebih cepat
- ✅ Route optimization calculation lebih cepat
- ✅ Image processing lebih smooth
- ✅ Multi-tasking lebih baik (app tetap berjalan di background)
- ✅ Update aplikasi lebih cepat
- ✅ Cache lebih banyak data untuk offline mode

### iOS (Rekomendasi untuk Penggunaan Optimal)

| Komponen       | Spesifikasi Rekomendasi                                        |
| -------------- | -------------------------------------------------------------- |
| **OS Version** | iOS 15.0+                                                      |
| **Device**     | iPhone 8 atau lebih baru (iPhone X, 11, 12, 13, 14, 15 series) |
| **RAM**        | 3 GB atau lebih                                                |
| **Storage**    | 1 GB free space (untuk instalasi + cache + updates)            |
| **CPU**        | A11 Bionic atau lebih baru                                     |
| **GPS**        | Built-in GPS + GLONASS                                         |
| **Camera**     | 7 MP atau lebih (untuk selfie berkualitas baik)                |
| **Network**    | 4G LTE / WiFi (untuk sync cepat)                               |
| **Screen**     | 5.5" atau lebih besar                                          |

**Keuntungan dengan Spesifikasi Rekomendasi:**

- ✅ Performa lebih konsisten
- ✅ Battery life lebih baik
- ✅ Update iOS lebih lama didukung
- ✅ Fitur-fitur baru iOS dapat digunakan

---

## Platform-Specific Requirements

### Android Permissions

Aplikasi memerlukan permission berikut:

```xml
<!-- Location -->
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />

<!-- Camera -->
<uses-permission android:name="android.permission.CAMERA" />

<!-- Storage -->
<uses-permission android:name="android.permission.READ_EXTERNAL_STORAGE" />
<uses-permission android:name="android.permission.WRITE_EXTERNAL_STORAGE" />

<!-- Network -->
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
```

### iOS Permissions

Aplikasi memerlukan permission berikut (didefinisikan di `Info.plist`):

```xml
<!-- Location -->
<key>NSLocationWhenInUseUsageDescription</key>
<string>We need your location to track visit reports and optimize routes</string>
<key>NSLocationAlwaysAndWhenInUseUsageDescription</key>
<string>We need your location to track visit reports and optimize routes</string>

<!-- Camera -->
<key>NSCameraUsageDescription</key>
<string>We need camera access to capture selfie for check-in</string>

<!-- Photo Library -->
<key>NSPhotoLibraryUsageDescription</key>
<string>We need photo library access to save and view visit report photos</string>
```

---

## Storage Requirements

### Instalasi Awal

- **APK Size**: ~50-80 MB (tergantung arsitektur: arm64, arm32, x86)
- **Installed Size**: ~150-200 MB (setelah instalasi dan extract native libraries)

### Data Cache & Offline Storage

| Tipe Data           | Perkiraan Ukuran | Catatan                              |
| ------------------- | ---------------- | ------------------------------------ |
| **Hive Database**   | 50-200 MB        | Tergantung jumlah data yang di-cache |
| **Route Cache**     | 10-50 MB         | Cache untuk route optimization       |
| **Image Cache**     | 20-100 MB        | Cache untuk foto visit reports       |
| **Dashboard Cache** | 5-20 MB          | Cache untuk dashboard data           |
| **Total Cache**     | **85-370 MB**    | Bervariasi berdasarkan usage         |

### Update Requirements

- **Space untuk Update**: Minimal 200 MB free space
- **Incremental Update**: Update biasanya 20-50 MB (hanya perubahan)
- **Full Update**: Update penuh bisa 80-150 MB

### Rekomendasi Total Storage

- **Minimum**: 500 MB free space
- **Recommended**: 1 GB free space (untuk buffer update dan cache growth)

---

## Network Requirements

### Koneksi Minimum

- **3G**: Minimum untuk sync data dasar
- **4G LTE**: Recommended untuk sync cepat dan maps loading
- **WiFi**: Optimal untuk update aplikasi dan sync besar

### Bandwidth Requirements

| Operasi             | Bandwidth Minimum | Bandwidth Recommended |
| ------------------- | ----------------- | --------------------- |
| **Login/Auth**      | 50 Kbps           | 100 Kbps              |
| **Sync Data Dasar** | 100 Kbps          | 500 Kbps              |
| **Load Maps Tiles** | 200 Kbps          | 1 Mbps                |
| **Upload Foto**     | 500 Kbps          | 2 Mbps                |
| **Update Aplikasi** | 1 Mbps            | 5 Mbps                |

### Offline Capability

- ✅ Aplikasi dapat berjalan **100% offline** untuk fitur read-only
- ✅ Data dapat di-sync ketika koneksi tersedia
- ⚠️ Fitur yang memerlukan koneksi: upload foto, sync data, route optimization (jika belum di-cache)

---

## Rekomendasi Spesifikasi HP

### Kategori: Budget (Minimal - Dapat Digunakan)

**Target**: HP dengan harga **Rp 1.5 - 2.5 juta**

#### Android Budget

| Brand       | Model         | RAM    | Storage | OS          | Harga (Approx)  |
| ----------- | ------------- | ------ | ------- | ----------- | --------------- |
| **Xiaomi**  | Redmi 9A / 9C | 2-3 GB | 32 GB   | Android 10+ | Rp 1.5-2 juta   |
| **Samsung** | Galaxy A03    | 3 GB   | 32 GB   | Android 11+ | Rp 1.8-2.2 juta |
| **Realme**  | C25 / C35     | 3-4 GB | 64 GB   | Android 11+ | Rp 2-2.5 juta   |
| **Oppo**    | A16           | 3 GB   | 32 GB   | Android 11+ | Rp 1.8-2.2 juta |

**Catatan untuk Budget:**

- ✅ Aplikasi dapat berjalan, namun performa maps mungkin lebih lambat
- ⚠️ RAM 2 GB adalah minimum mutlak; 3 GB lebih baik
- ⚠️ Storage 32 GB mungkin terbatas jika banyak foto
- ✅ Cocok untuk penggunaan dasar (view data, check-in sederhana)

---

### Kategori: Mid-Range (Rekomendasi - Optimal)

**Target**: HP dengan harga **Rp 3 - 6 juta**

#### Android Mid-Range

| Brand       | Model              | RAM    | Storage   | OS            | Harga (Approx) |
| ----------- | ------------------ | ------ | --------- | ------------- | -------------- |
| **Xiaomi**  | Redmi Note 11 / 12 | 4-6 GB | 64-128 GB | Android 11-13 | Rp 3-4.5 juta  |
| **Samsung** | Galaxy A23 / A33   | 4-6 GB | 64-128 GB | Android 12-13 | Rp 3.5-5 juta  |
| **Realme**  | 9 / 10             | 4-6 GB | 64-128 GB | Android 11-13 | Rp 3-4.5 juta  |
| **Oppo**    | A57 / A77          | 4-6 GB | 64-128 GB | Android 11-13 | Rp 3-4.5 juta  |
| **Vivo**    | Y21 / Y33          | 4 GB   | 64 GB     | Android 11-12 | Rp 3-3.5 juta  |

**Keuntungan Mid-Range:**

- ✅ Performa maps sangat baik
- ✅ Route optimization cepat
- ✅ Multi-tasking lancar
- ✅ Camera quality cukup untuk selfie
- ✅ Battery life lebih baik
- ✅ Update Android lebih lama didukung

#### iOS Mid-Range

| Model                     | RAM    | Storage   | OS        | Harga (Approx)         |
| ------------------------- | ------ | --------- | --------- | ---------------------- |
| **iPhone SE (2020/2022)** | 3 GB   | 64-128 GB | iOS 15-17 | Rp 4-6 juta (second)   |
| **iPhone 8 / 8 Plus**     | 2-3 GB | 64-256 GB | iOS 15-16 | Rp 2.5-4 juta (second) |

**Catatan iOS:**

- ✅ iPhone SE 2022 adalah pilihan terbaik untuk budget iOS
- ⚠️ iPhone 8 sudah tidak didukung update iOS terbaru (tapi masih bisa pakai app)

---

### Kategori: Premium (Sangat Optimal - Future Proof)

**Target**: HP dengan harga **Rp 6 - 15 juta**

#### Android Premium

| Brand       | Model             | RAM    | Storage    | OS            | Harga (Approx) |
| ----------- | ----------------- | ------ | ---------- | ------------- | -------------- |
| **Samsung** | Galaxy A54 / A73  | 6-8 GB | 128-256 GB | Android 13-14 | Rp 5-8 juta    |
| **Xiaomi**  | Redmi Note 12 Pro | 6-8 GB | 128-256 GB | Android 12-13 | Rp 4-6 juta    |
| **OnePlus** | Nord CE / 10T     | 8 GB   | 128-256 GB | Android 12-13 | Rp 5-7 juta    |
| **Google**  | Pixel 6a / 7a     | 6-8 GB | 128 GB     | Android 12-14 | Rp 6-9 juta    |

**Keuntungan Premium:**

- ✅ Performa sangat cepat untuk semua fitur
- ✅ Camera sangat baik untuk selfie
- ✅ Battery life sangat baik
- ✅ Update Android lebih lama (3-4 tahun)
- ✅ Future-proof untuk update aplikasi

#### iOS Premium

| Model              | RAM  | Storage    | OS        | Harga (Approx)         |
| ------------------ | ---- | ---------- | --------- | ---------------------- |
| **iPhone 11**      | 4 GB | 64-256 GB  | iOS 15-17 | Rp 6-8 juta (second)   |
| **iPhone 12**      | 4 GB | 64-256 GB  | iOS 15-17 | Rp 8-10 juta (second)  |
| **iPhone 13**      | 4 GB | 128-256 GB | iOS 15-17 | Rp 10-13 juta (second) |
| **iPhone SE 2022** | 4 GB | 64-256 GB  | iOS 15-17 | Rp 5-7 juta (baru)     |

**Keuntungan iOS Premium:**

- ✅ Performa sangat konsisten
- ✅ Update iOS lebih lama (5-6 tahun)
- ✅ Camera sangat baik
- ✅ Battery life sangat baik
- ✅ Resale value tinggi

---

## Pertimbangan Jangka Panjang

### Update Aplikasi

Aplikasi akan menerima update berkala dengan fitur-fitur baru. Pertimbangan untuk jangka panjang:

#### 1. **OS Update Support**

**Android:**

- ⚠️ **Android 5.0-7.1 (API 21-25)**: Support akan dihentikan dalam 1-2 tahun
- ✅ **Android 8.0-10 (API 26-29)**: Support hingga 2-3 tahun
- ✅ **Android 11+ (API 30+)**: Support jangka panjang (3-5 tahun)

**iOS:**

- ⚠️ **iOS 12-14**: Support akan dihentikan dalam 1-2 tahun
- ✅ **iOS 15+**: Support jangka panjang (3-5 tahun)

**Rekomendasi:**

- Pilih HP dengan OS terbaru saat pembelian (minimal Android 11 / iOS 15)
- Pastikan HP masih menerima update OS dari manufacturer

#### 2. **Storage Growth**

Aplikasi akan terus berkembang dengan fitur-fitur baru:

- **Tahun 1**: ~200 MB (instalasi + cache)
- **Tahun 2**: ~300 MB (dengan fitur baru)
- **Tahun 3**: ~400-500 MB (dengan update besar)

**Rekomendasi:**

- Minimal **64 GB storage** untuk penggunaan jangka panjang
- **128 GB** lebih baik untuk buffer dan foto banyak

#### 3. **RAM Requirements**

Fitur-fitur baru mungkin memerlukan lebih banyak RAM:

- **Saat ini**: 2 GB minimum, 4 GB recommended
- **1-2 tahun**: 3 GB minimum, 6 GB recommended
- **3+ tahun**: 4 GB minimum, 8 GB recommended

**Rekomendasi:**

- Pilih HP dengan **minimal 4 GB RAM** untuk jangka panjang
- **6 GB RAM** lebih baik untuk future-proof

#### 4. **Battery Life**

Aplikasi menggunakan GPS dan maps yang menguras battery:

- **GPS Tracking**: Menggunakan battery signifikan
- **Maps Rendering**: Menggunakan GPU dan battery
- **Background Sync**: Menggunakan battery saat sync

**Rekomendasi:**

- Pilih HP dengan **minimal 4000 mAh battery** untuk penggunaan seharian
- **5000 mAh** lebih baik untuk penggunaan intensif

#### 5. **Camera Quality**

Selfie check-in memerlukan camera yang baik:

- **Saat ini**: 2 MP minimum, 8 MP recommended
- **Masa depan**: Mungkin ada requirement kualitas foto lebih tinggi

**Rekomendasi:**

- Pilih HP dengan **minimal 8 MP front camera**
- **12 MP atau lebih** lebih baik untuk kualitas foto

---

## Testing & Compatibility

### Device Testing Matrix

Aplikasi telah diuji pada device berikut:

#### Android

| Device                   | OS Version | RAM  | Status     | Notes             |
| ------------------------ | ---------- | ---- | ---------- | ----------------- |
| **Samsung Galaxy A23**   | Android 12 | 4 GB | ✅ Optimal | Recommended       |
| **Xiaomi Redmi Note 11** | Android 11 | 4 GB | ✅ Optimal | Recommended       |
| **Samsung Galaxy A03**   | Android 11 | 3 GB | ✅ Good    | Budget option     |
| **Realme C35**           | Android 11 | 4 GB | ✅ Optimal | Budget-Mid option |

#### iOS

| Device             | OS Version | RAM  | Status        | Notes                     |
| ------------------ | ---------- | ---- | ------------- | ------------------------- |
| **iPhone 12**      | iOS 15-17  | 4 GB | ✅ Optimal    | Recommended               |
| **iPhone SE 2022** | iOS 15-17  | 4 GB | ✅ Optimal    | Budget iOS option         |
| **iPhone 8**       | iOS 15-16  | 2 GB | ⚠️ Acceptable | Older device, slower maps |

### Known Issues & Limitations

#### Android 5.0-7.1 (API 21-25)

- ⚠️ Maps loading lebih lambat
- ⚠️ Route optimization calculation lebih lambat
- ⚠️ Battery drain lebih cepat
- ✅ Masih dapat digunakan untuk fitur dasar

#### RAM < 2 GB

- ❌ Aplikasi mungkin crash saat menggunakan maps
- ❌ Multi-tasking sangat terbatas
- ❌ Image processing mungkin gagal

#### Storage < 500 MB Free

- ⚠️ Update aplikasi mungkin gagal
- ⚠️ Cache terbatas, offline mode kurang optimal
- ⚠️ Foto mungkin tidak dapat disimpan

---

## Summary & Quick Reference

### Minimum Requirements (Dapat Digunakan)

- **Android**: Android 5.0+, 2 GB RAM, 500 MB storage
- **iOS**: iOS 12.0+, iPhone 6s+, 2 GB RAM, 500 MB storage
- **Harga HP**: Rp 1.5-2.5 juta

### Recommended Requirements (Optimal)

- **Android**: Android 10.0+, 4 GB RAM, 1 GB storage
- **iOS**: iOS 15.0+, iPhone 8+, 3 GB RAM, 1 GB storage
- **Harga HP**: Rp 3-6 juta

### Premium Requirements (Future-Proof)

- **Android**: Android 11+, 6 GB RAM, 128 GB storage
- **iOS**: iOS 15+, iPhone 11+, 4 GB RAM, 128 GB storage
- **Harga HP**: Rp 6-15 juta

---

## Support & Contact

Jika ada pertanyaan tentang system requirements atau kompatibilitas device, silakan hubungi tim development.

**Last Updated**: 2025-01-15  
**Next Review**: Setiap 6 bulan atau setelah major update aplikasi
