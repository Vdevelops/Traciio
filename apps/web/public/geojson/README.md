# GeoJSON Data for Indonesia Map

Place your Indonesia GeoJSON files here.

## ⚠️ Important

File `indonesia-provinces-simple.geojson` yang ada saat ini **hanya untuk testing** dengan koordinat kotak sederhana. Untuk peta yang akurat dengan bentuk geografis yang benar, Anda perlu download file GeoJSON lengkap.

**Lihat file `DOWNLOAD_GEOJSON.md` untuk instruksi lengkap download GeoJSON yang akurat.**

## Recommended Files

1. **indonesia-regencies.geojson** - GeoJSON file containing all regencies/cities in Indonesia dengan koordinat geografis yang akurat
   - Should be a FeatureCollection with properties: `name` (regency name) and `province` (province name)
   - File size: Usually 5-20MB for complete data
   - Coordinate system: WGS84 (EPSG:4326)

2. **indonesia-provinces-simple.geojson** - File sederhana untuk testing (sudah ada)
   - Hanya berisi 4 provinsi dengan koordinat kotak sederhana
   - **TIDAK AKURAT** - hanya untuk testing

## Quick Download

Untuk download file lengkap dengan koordinat akurat, jalankan:

**PowerShell:**
```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-city-regency.geojson" -OutFile "indonesia-regencies.geojson"
```

**Bash/Linux/Mac:**
```bash
curl -L "https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-city-regency.geojson" -o "indonesia-regencies.geojson"
```

## Sources

You can obtain accurate GeoJSON data from:
- https://github.com/superpikar/indonesia-geojson (Recommended)
- https://github.com/ans-4175/peta-indonesia-geojson
- Or create your own GeoJSON file

## File Structure

The GeoJSON should follow this structure:

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
        "coordinates": [...]
      }
    }
  ]
}
```

## Notes

- File size: Keep the file as small as possible for better performance
- If the file is too large (>5MB), consider using a simplified version or province-level data instead
- The map component will automatically load from `/geojson/indonesia-regencies.geojson`

