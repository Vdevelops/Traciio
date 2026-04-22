# GPS Setup untuk Mobile App

## ⚠️ Penting: GPS di Emulator

**GPS di emulator Android/iOS tidak otomatis real-time** karena **emulator tidak punya GPS fisik**.

Untuk mendapatkan lokasi GPS yang akurat di emulator, kamu **HARUS inject lokasi secara manual**. Di real device, GPS akan bekerja otomatis dengan akurasi tinggi.

## ✅ Best Practice untuk Flutter

1. **Gunakan `LocationAccuracy.bestForNavigation`** (bukan `low` atau `medium`)
2. **Inject lokasi manual di emulator** via Extended Controls
3. **Gunakan GPX file** untuk simulasi pergerakan yang realistis
4. **Test di real device** untuk akurasi GPS yang sebenarnya

## Android Emulator

### ✅ Cara Paling Akurat (Manual Inject Location)

### 1️⃣ Buka Emulator → Extended Controls

1. Jalankan emulator
2. Klik **⋮ (More / Extended Controls)** di toolbar emulator
3. Pilih **Location** tab

Ada 3 mode:

- **Single points** → manual (paling sering dipakai)
- **Route** → simulasi perjalanan
- **GPX/KML** → paling akurat (recommended)

### 🔹 Opsi A: Manual Latitude & Longitude

Masukkan koordinat asli (misalnya dari Google Maps):

```
Latitude  : -7.795580
Longitude : 110.369490
```

Klik **Send**

✔ Emulator langsung pindah ke lokasi tersebut  
✔ Flutter akan menerima update GPS

### 🔹 Opsi B: GPX (PALING AKURAT & REALISTIS)

Kalau kamu butuh simulasi **bergerak / tracking / realtime**:

1. Buat file `route.gpx`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.1" creator="Flutter">
  <trk>
    <trkseg>
      <trkpt lat="-7.795580" lon="110.369490"></trkpt>
      <trkpt lat="-7.796000" lon="110.370000"></trkpt>
      <trkpt lat="-7.797000" lon="110.371000"></trkpt>
    </trkseg>
  </trk>
</gpx>
```

2. Upload di **Extended Controls → Location → Load GPX**
3. Klik **Play**

✔ Cocok untuk:

- Tracking lokasi
- Nearby hospital (PostGIS / ST_DWithin)
- Map movement
- Ride / delivery simulation

### 🔹 Opsi C: Set Lokasi via ADB Command

```bash
# Set lokasi ke Jakarta, Indonesia
adb emu geo fix 106.8456 -6.2088

# Set lokasi ke Bandung, Indonesia
adb emu geo fix 107.6098 -6.9175

# Set lokasi ke Surabaya, Indonesia
adb emu geo fix 112.7521 -7.2575

# Set lokasi ke Yogyakarta, Indonesia
adb emu geo fix 110.369490 -7.795580
```

### ⚠️ Setting Emulator yang WAJIB Dicek

1. Di emulator, buka **Settings**
2. Pilih **Location**
3. Pastikan **Location** toggle **ON**
4. Pilih **Mode** → **High accuracy** (GPS, Wi-Fi, and mobile networks)

## iOS Simulator

### 1️⃣ Menu Bar Simulator

1. Di Simulator, klik menu **Features** → **Location**
2. Pilih:
   - **Custom Location** → Set manual koordinat
   - **GPX file** → Load file GPX untuk simulasi pergerakan
   - Preset locations (Apple, City Bicycle Ride, dll)

### 2️⃣ Set Custom Location

1. **Features** → **Location** → **Custom Location...**
2. Set:
   - **Latitude**: Contoh `-7.795580`
   - **Longitude**: Contoh `110.369490`
3. Klik **OK**

## Testing GPS di Real Device

Untuk testing GPS yang akurat, **disarankan menggunakan real device** karena:

1. **GPS Hardware**: Real device memiliki GPS hardware yang sebenarnya
2. **Akurasi**: GPS di real device lebih akurat (biasanya < 10 meter dengan `bestForNavigation`)
3. **Network Location**: Real device bisa menggunakan network location untuk akurasi lebih baik
4. **Satellite Fix**: Real device bisa mendapatkan fix dari GPS satellites
5. **Real-time Updates**: GPS di real device update secara real-time tanpa perlu inject manual

## Konfigurasi GPS di Code (BEST PRACTICE)

Aplikasi sudah dikonfigurasi dengan **best practice**:

```dart
Position position = await Geolocator.getCurrentPosition(
  desiredAccuracy: LocationAccuracy.bestForNavigation, // ✅ BEST PRACTICE: Most accurate for navigation
  timeLimit: const Duration(seconds: 20), // Allow time for GPS to get accurate fix
  forceAndroidLocationManager: false, // Use GPS provider (not network location)
);
```

### ⚠️ JANGAN Gunakan LocationAccuracy.low

**LocationAccuracy.low** akan menyebabkan:

- Lokasi loncat-loncat
- Akurasi sangat rendah (100+ meter)
- Tidak cocok untuk navigation atau tracking

### Akurasi GPS (Best to Worst)

- **LocationAccuracy.bestForNavigation** ✅ **RECOMMENDED**: Akurasi terbaik untuk navigation (bisa sampai beberapa meter)
- **LocationAccuracy.best**: Akurasi terbaik (bisa sampai beberapa meter)
- **LocationAccuracy.high**: Akurasi tinggi (sekitar 10-50 meter)
- **LocationAccuracy.medium**: Akurasi sedang (sekitar 50-100 meter)
- **LocationAccuracy.low**: ❌ **JANGAN PAKAI**: Akurasi rendah (sekitar 100+ meter, lokasi loncat-loncat)

### Real-Time Location Tracking (Optional)

Untuk real-time tracking (contoh: live tracking, map movement):

```dart
Geolocator.getPositionStream(
  locationSettings: LocationSettings(
    accuracy: LocationAccuracy.bestForNavigation,
    distanceFilter: 5, // Update setiap 5 meter pergerakan
  ),
).listen((Position position) {
  print('Lat: ${position.latitude}, Lng: ${position.longitude}');
  print('Accuracy: ${position.accuracy} meters');
});
```

### Pengecekan Akurasi

Aplikasi akan:

1. Log akurasi GPS yang didapat
2. Jika akurasi > 50 meter, akan mencoba mendapatkan posisi yang lebih baik
3. Menampilkan warning jika akurasi kurang baik

## ❌ Kesalahan Umum (Ini Bikin Tidak Akurat)

| Kesalahan                 | Dampak              | Solusi                             |
| ------------------------- | ------------------- | ---------------------------------- |
| Mengandalkan IP Location  | Lokasi melenceng    | Inject lokasi manual di emulator   |
| Emulator tanpa inject GPS | Lokasi selalu 0.0   | Set lokasi via Extended Controls   |
| Accuracy = low            | Titik loncat-loncat | Gunakan `bestForNavigation`        |
| Tidak grant permission    | Lokasi null         | Request permission dengan benar    |
| Timeout terlalu pendek    | GPS tidak dapat fix | Tingkatkan `timeLimit` ke 20 detik |

## Troubleshooting

### GPS tidak akurat di emulator

1. **Pastikan Location Services enabled** di emulator settings
2. **Set lokasi manual** menggunakan Extended Controls (Single points atau GPX)
3. **Gunakan ADB command** untuk set lokasi: `adb emu geo fix <lon> <lat>`
4. **Test di real device** untuk akurasi yang lebih baik

### GPS timeout

1. **Tingkatkan timeLimit** (default 20 detik sudah optimal)
2. **Pastikan Location Services enabled** di emulator/device settings
3. **Check permission** sudah granted
4. **Di emulator**: Inject lokasi manual untuk menghindari timeout

### GPS tidak update / lokasi tidak berubah

1. **Restart emulator** dan inject lokasi lagi
2. **Clear app data** dan restart app
3. **Check permission** di app settings
4. **Di emulator**: Pastikan sudah inject lokasi via Extended Controls

## 🎯 Ringkasan Cepat

✔ Emulator **tidak bisa GPS asli** - Harus inject lokasi manual  
✔ Lokasi **HARUS di-inject** via Extended Controls atau ADB  
✔ **GPX > Manual Input > Preset** untuk simulasi pergerakan  
✔ Flutter pakai **`bestForNavigation`** (bukan `low`)  
✔ Cocok untuk testing hospital nearby, maps, tracking  
✔ **Real device** lebih akurat untuk testing GPS sebenarnya

## Catatan Penting

- **Emulator GPS tidak akurat**: Emulator menggunakan lokasi mock, bukan GPS hardware sebenarnya
- **Harus inject lokasi**: Di emulator, lokasi HARUS di-inject manual via Extended Controls atau ADB
- **Real device lebih akurat**: Untuk testing GPS yang akurat, gunakan real device
- **Network Location**: Di real device, GPS bisa menggunakan network location untuk akurasi lebih baik
- **Satellite Fix**: Real device perlu beberapa detik untuk mendapatkan GPS satellite fix
- **Best Practice**: Selalu gunakan `LocationAccuracy.bestForNavigation` untuk akurasi terbaik

## Referensi

- [Geolocator Package](https://pub.dev/packages/geolocator)
- [Android Emulator Location](https://developer.android.com/studio/run/emulator#extended)
- [iOS Simulator Location](https://developer.apple.com/documentation/corelocation/cllocationmanager)
