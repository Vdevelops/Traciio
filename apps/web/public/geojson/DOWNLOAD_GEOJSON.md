# Download GeoJSON Indonesia yang Akurat

⚠️ **PENTING:** File GeoJSON sederhana (`indonesia-provinces-simple.geojson`) yang ada hanya untuk testing dengan koordinat kotak sederhana.

Untuk peta yang **akurat dengan bentuk geografis yang benar** (bukan kotak-kotak), Anda perlu download GeoJSON lengkap dengan koordinat geografis yang sebenarnya.

## Opsi 1: Download dari GitHub (Recommended)

### Langkah-langkah:

1. **Kunjungi repository GeoJSON Indonesia:**
   - https://github.com/superpikar/indonesia-geojson
   - Atau: https://github.com/ans-4175/peta-indonesia-geojson

2. **Pilih file yang sesuai:**
   - Untuk regency/kota level: `indonesia-province-city-regency.geojson`
   - Atau file yang berisi data kabupaten/kota lengkap

3. **Download file:**
   - Klik kanan pada file → "Save Link As"
   - Atau gunakan direct download link (lihat di bawah)

4. **Rename dan simpan:**
   - Rename file menjadi: `indonesia-regencies.geojson`
   - Simpan di: `apps/web/public/geojson/indonesia-regencies.geojson`

## Opsi 2: Direct Download Links

### Repository: superpikar/indonesia-geojson

**Regency/City Level (Recommended):**

```
https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-city-regency.geojson
```

**Province Level (Lebih kecil, kurang detail):**

```
https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-simple.geojson
```

### Cara Download via Command Line (PowerShell):

```powershell
# Download GeoJSON regency level
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-city-regency.geojson" -OutFile "apps/web/public/geojson/indonesia-regencies.geojson"
```

### Cara Download via Command Line (Bash/Linux/Mac):

```bash
# Download GeoJSON regency level
curl -L "https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-city-regency.geojson" -o "apps/web/public/geojson/indonesia-regencies.geojson"
```

## Opsi 3: Download via Browser

1. Buka link berikut di browser:

   ```
   https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-city-regency.geojson
   ```

2. Browser akan menampilkan JSON content
3. Klik kanan → "Save As" atau tekan `Ctrl+S` (Windows) / `Cmd+S` (Mac)
4. Simpan sebagai: `indonesia-regencies.geojson`
5. Pindahkan file ke: `apps/web/public/geojson/`

## Verifikasi File

Setelah download, pastikan:

1. **File size:** Harus lebih dari 1MB (file lengkap biasanya 2-10MB)
2. **File structure:** Buka file dan pastikan memiliki struktur:
   ```json
   {
     "type": "FeatureCollection",
     "features": [
       {
         "type": "Feature",
         "properties": {
           "name": "Nama Regency",
           "province": "Nama Provinsi"
         },
         "geometry": {
           "type": "Polygon" atau "MultiPolygon",
           "coordinates": [[[lng, lat], [lng, lat], ...]]
         }
       }
     ]
   }
   ```
3. **Test di browser:** Buka `http://localhost:3000/geojson/indonesia-regencies.geojson` untuk memastikan file dapat diakses

## Catatan Penting

- **File size:** File GeoJSON lengkap bisa cukup besar (5-20MB). Ini normal untuk data geografis yang detail.
- **Performance:** File besar mungkin membutuhkan waktu loading lebih lama. Pertimbangkan menggunakan versi simplified jika performa menjadi masalah.
- **Coordinate system:** Pastikan file menggunakan WGS84 (EPSG:4326) dengan format `[longitude, latitude]`.

## Troubleshooting

### File terlalu besar dan lambat?

Gunakan versi simplified atau province-level:

- Province level: Lebih kecil, kurang detail
- Simplified regency: Koordinat sudah di-reduce, lebih cepat

### File tidak muncul di map?

1. Pastikan file ada di lokasi yang benar
2. Restart dev server: `pnpm dev`
3. Clear browser cache
4. Check console untuk error messages

### Map masih menampilkan kotak-kotak?

- Pastikan file `indonesia-regencies.geojson` sudah ada (bukan hanya file simple)
- File simple hanya untuk fallback jika file lengkap tidak tersedia
- Refresh halaman setelah menambahkan file baru
