# PowerShell Script to Download Accurate Indonesia GeoJSON
# This script downloads GeoJSON with real geographic coordinates (not simple boxes)

Write-Host "Downloading Indonesia GeoJSON with accurate geographic coordinates..." -ForegroundColor Cyan
Write-Host ""

$geojsonUrl = "https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-city-regency.geojson"
$outputFile = "indonesia-regencies.geojson"

try {
    Write-Host "Downloading from: $geojsonUrl" -ForegroundColor Yellow
    Write-Host "Saving to: $outputFile" -ForegroundColor Yellow
    Write-Host ""
    
    Invoke-WebRequest -Uri $geojsonUrl -OutFile $outputFile -UseBasicParsing
    
    $fileSize = (Get-Item $outputFile).Length / 1MB
    Write-Host "✓ Download successful!" -ForegroundColor Green
    Write-Host "  File size: $([math]::Round($fileSize, 2)) MB" -ForegroundColor Green
    Write-Host ""
    Write-Host "The map will now display accurate geographic shapes instead of simple boxes." -ForegroundColor Cyan
    Write-Host "Refresh your browser to see the updated map." -ForegroundColor Cyan
} catch {
    Write-Host "✗ Download failed!" -ForegroundColor Red
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
    Write-Host "Alternative: Download manually from:" -ForegroundColor Yellow
    Write-Host "  https://github.com/superpikar/indonesia-geojson" -ForegroundColor Yellow
    Write-Host "  Save as: indonesia-regencies.geojson" -ForegroundColor Yellow
}

