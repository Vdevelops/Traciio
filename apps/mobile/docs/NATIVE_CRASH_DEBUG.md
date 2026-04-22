# Debug Native Crash (SIGQUIT / Tombstone) – Halaman Accounts

Crash dengan log:
- `Signal Catcher: reacting to signal 3` (SIGQUIT)
- `Wrote stack traces to tombstoned`
- `Lost connection to device. Exited.`

berarti **crash di kode native (Android)**, bukan di Dart. Perbaikan harus mengikuti langkah di bawah ini.

---

## Prioritas 1: Ambil file Tombstone (paling penting)

Stack trace asli ada di file tombstone. Tanpa ini kita hanya menebak.

### 1. Sambungkan device/emulator dan pastikan app sudah crash sekali

```bash
# Cek device
adb devices

# Lihat daftar tombstone (file terbaru = crash terakhir)
adb shell ls -l /data/tombstones/
```

### 2. Pull tombstone terbaru

**Cara cepat (Windows):** Jalankan script setelah app crash sekali:

```bash
# Dari repo root atau apps/mobile
apps\mobile\scripts\pull_tombstone.bat
```

Script akan coba pull `tombstone_00` ke `apps/mobile/tombstone_00` atau fallback ke `logcat_crash.txt`.

**Manual:** Di banyak device/emulator **tanpa root**, path bisa berbeda atau akses dibatasi. Coba:

```bash
# Coba pull langsung (emulator sering bisa)
adb pull /data/tombstones/tombstone_00 ./tombstone_00

# Jika "Permission denied", coba:
adb shell "run-as com.example.mobile cat /data/tombstones/tombstone_00" > tombstone_00

# Atau lihat logcat saat crash (backup jika tombstone tidak bisa diambil)
adb logcat -d > logcat_crash.txt
```

### 3. Baca isi tombstone

Buka `tombstone_00` (atau nama file terbaru):

- **Cek dulu process-nya**: Cari baris **`>>> com.example.mobile <<<`** (package kamu).  
  Jika yang muncul **bukan** package kamu (mis. `>>> com.google.android.bluetooth <<<`), artinya tombstone itu **bukan crash dari app kamu** — itu crash dari proses sistem/emulator. Abaikan untuk perbaikan app, atau lihat [Tombstone dari proses sistem](#tombstone-dari-proses-sistem-bukan-app).
- Cari **signal** (mis. SIGSEGV, SIGABRT, SIGILL) — ini penyebab crash.
- Lihat **backtrace** / stack trace native — itu lokasi crash (plugin / engine / lib).

Dari sini bisa ditentukan: plugin mana, atau apakah Flutter engine, dll.

#### Tombstone dari proses sistem (bukan app)

Jika di tombstone tertulis proses lain, contoh:

- `>>> com.google.android.bluetooth <<<`  
  **Abort message:** `Hardware Error Event with code 0x42` (HCI layer)

itu **crash di stack Bluetooth emulator**, bukan di kode app. Tidak perlu diperbaiki di kode Flutter.

Yang bisa dilakukan:

- **Emulator**: Nonaktifkan Bluetooth di emulator, atau buat AVD baru tanpa Bluetooth / gunakan image lain.
- **Device fisik**: Pastikan Bluetooth stabil; crash seperti ini sering spesifik emulator.
- Jika log yang kamu lihat hanya "Signal Catcher: reacting to signal 3" dan "Wrote stack traces to tombstoned", itu bisa **dump stack** (SIGQUIT), bukan pasti crash app. Crash sesungguhnya bisa di proses lain (mis. Bluetooth) — cek tombstone yang berisi **`>>> com.example.mobile <<<`** untuk memastikan crash app.

---

## Prioritas 2: Tes di mode Release

Banyak native crash hanya muncul di debug (atau hanya di release).

```bash
cd apps/mobile
flutter run --release
```

Lalu buka halaman Accounts, lakukan aksi yang biasa memicu crash (mis. scroll, ganti tab, tekan back).

- Jika crash **hilang di release** → sering terkait debug/tooling (hot reload, observatory, dll.).
- Jika crash **tetap di release** → lebih mungkin bug plugin/engine/skia; tombstone jadi acuan utama.

---

## Prioritas 3: Catat aksi user saat crash

Catat persis apa yang dilakukan:

- [ ] Baru buka Accounts dari dashboard?
- [ ] Scroll daftar account?
- [ ] Ganti tab Accounts ↔ Contacts?
- [ ] Tekan back (keluar dari Accounts)?
- [ ] Pull-to-refresh?
- [ ] Baru hot reload / hot restart?

Ini membantu korelasikan dengan stack trace di tombstone (mis. “crash saat pop route” → cari frame yang berhubungan dengan teardown).

---

## Prioritas 4: Isolasi plugin (tes satu per satu)

Plugin yang sering terkait native crash di Flutter:

| Plugin (di `pubspec.yaml`) | Dipakai di | Cara tes |
|----------------------------|------------|----------|
| **image_picker** | Visit reports, selfie | Comment di `pubspec.yaml` → `flutter pub get` → tes Accounts |
| **geolocator** | Visit reports, route form | Idem |
| **flutter_map** | Route optimization | Idem |
| **connectivity_plus** | Global (ConnectivityService) | Idem |
| **hive** / **hive_flutter** | Storage | Idem |

Langkah:

1. Backup `pubspec.yaml`.
2. Comment satu dependency (dan yang tergantung padanya jika ada).
3. `flutter pub get`, lalu `flutter run`.
4. Buka Accounts dan lakukan aksi yang biasa memicu crash.
5. Jika crash **hilang** → plugin yang di-comment kemungkinan terlibat. Cek changelog / issue plugin, upgrade, atau cari alternatif.
6. Restore `pubspec.yaml`, lanjut ke plugin berikutnya.

Contoh comment:

```yaml
dependencies:
  flutter:
    sdk: flutter
  # image_picker: ^1.1.2   # <-- comment untuk tes
```

---

## Prioritas 5: Logging di MainActivity (opsional)

Sudah ditambahkan log di `MainActivity.kt`. Setelah crash:

```bash
adb logcat -s FlutterCrash:* flutter:*
```

Lalu jalankan app, buka Accounts, dan lakukan aksi sampai crash. Cek log:

- Jika `FlutterCrash: Engine configured` muncul tapi lalu crash → crash terjadi **setelah** engine siap (mis. saat navigasi/teardown).
- Jika log tidak muncul sama sekali → crash sangat awal (sebelum `configureFlutterEngine`).

---

## Ringkasan langkah perbaikan

1. **Ambil tombstone** → baca signal + backtrace → tentukan lokasi crash (plugin / engine).
2. **Tes `flutter run --release`** → apakah crash hanya di debug atau juga di release.
3. **Catat aksi user** saat crash → korelasi dengan stack trace.
4. **Isolasi plugin** (comment satu per satu) → cari plugin yang memicu crash.
5. **Gunakan log MainActivity** untuk memastikan kapan crash relatif terhadap inisialisasi engine.

Setelah dapat tombstone (atau minimal logcat lengkap saat crash), kita bisa memastikan penyebab dan perbaikan yang tepat (update plugin, workaround, atau laporkan ke Flutter/plugin).
