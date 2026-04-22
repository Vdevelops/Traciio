# Setup GeoJSON untuk Brick Management Map

## Overview

Fitur Brick Management menggunakan peta interaktif untuk memilih regency/kota. Peta ini memerlukan file GeoJSON yang berisi data geografis Indonesia.

## Lokasi File

File GeoJSON harus ditempatkan di:

```
apps/web/public/geojson/indonesia-regencies.geojson
```

## Cara Mendapatkan GeoJSON File

### Option 1: Download dari GitHub Repository

1. Kunjungi repository GeoJSON Indonesia:
   - https://github.com/superpikar/indonesia-geojson
   - https://github.com/ans-4175/peta-indonesia-geojson

2. Download file yang sesuai:
   - Pilih file yang berisi data regency/kota (bukan hanya provinsi)
   - File yang recommended: `indonesia-province-city-regency.geojson` atau sejenisnya

3. Rename file menjadi `indonesia-regencies.geojson`

4. Copy file ke folder `apps/web/public/geojson/`

### Option 2: Download dari CDN (temporary)

Jika file belum tersedia di public folder, sistem akan mencoba load dari CDN sebagai fallback. Namun, ini tidak reliable dan sebaiknya digunakan hanya untuk development.

### Option 3: Buat Sendiri (Advanced)

Jika Anda memiliki data geografis sendiri, buat GeoJSON file dengan struktur:

```json
{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "properties": {
        "name": "Jakarta Pusat",
        "province": "DKI Jakarta"
      },
      "geometry": {
        "type": "Polygon",
        "coordinates": [
          [
            [106.8, -6.2],
            [106.9, -6.2],
            [106.9, -6.3],
            [106.8, -6.3],
            [106.8, -6.2]
          ]
        ]
      }
    }
  ]
}
```

## Struktur GeoJSON yang Diperlukan

File GeoJSON harus memiliki struktur berikut:

- **type**: Harus `"FeatureCollection"`
- **features**: Array of Feature objects
  - Setiap feature harus memiliki:
    - **type**: `"Feature"`
    - **properties**: Object dengan minimal:
      - `name`: Nama regency/kota (required)
      - `province`: Nama provinsi (required)
    - **geometry**: Object dengan type `"Polygon"` atau `"MultiPolygon"` dan coordinates

## Validasi File

Setelah menempatkan file, pastikan:

1. File dapat diakses di browser: `http://localhost:3000/geojson/indonesia-regencies.geojson`
2. File valid JSON (dapat di-parse)
3. File memiliki structure yang benar (FeatureCollection dengan features array)
4. Setiap feature memiliki properties `name` dan `province`

## Performance Considerations

- **File Size**: Usahakan file tidak lebih dari 5MB untuk performa optimal
- **Simplified Geometry**: Jika file terlalu besar, gunakan versi yang sudah di-simplify (reduced coordinates)
- **Province-level**: Untuk development/testing, bisa menggunakan data level provinsi yang lebih kecil

## Troubleshooting

### Map tidak muncul

1. Check apakah file ada di lokasi yang benar: `apps/web/public/geojson/indonesia-regencies.geojson`
2. Check console browser untuk error messages
3. Verify file dapat diakses via URL: `http://localhost:3000/geojson/indonesia-regencies.geojson`
4. Pastikan file valid JSON format
5. Pastikan file memiliki structure FeatureCollection

### Error 404

- File tidak ditemukan di public folder
- Pastikan file sudah di-copy ke folder yang benar
- Restart dev server setelah menambahkan file baru

### Map muncul tapi tidak bisa di-click

- Pastikan GeoJSON features memiliki properties `name` dan `province`
- Check console untuk error messages

## Alternative: Use Manual Input

Jika GeoJSON file tidak tersedia, user tetap dapat membuat brick menggunakan tab "Manual Input" untuk mengisi province dan regency/city secara manual.
