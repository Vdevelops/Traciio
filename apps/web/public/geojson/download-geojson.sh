#!/bin/bash
# Bash Script to Download Accurate Indonesia GeoJSON
# This script downloads GeoJSON with real geographic coordinates (not simple boxes)

echo "Downloading Indonesia GeoJSON with accurate geographic coordinates..."
echo ""

GEOJSON_URL="https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-city-regency.geojson"
OUTPUT_FILE="indonesia-regencies.geojson"

echo "Downloading from: $GEOJSON_URL"
echo "Saving to: $OUTPUT_FILE"
echo ""

if curl -L "$GEOJSON_URL" -o "$OUTPUT_FILE"; then
    FILE_SIZE=$(du -h "$OUTPUT_FILE" | cut -f1)
    echo "✓ Download successful!"
    echo "  File size: $FILE_SIZE"
    echo ""
    echo "The map will now display accurate geographic shapes instead of simple boxes."
    echo "Refresh your browser to see the updated map."
else
    echo "✗ Download failed!"
    echo ""
    echo "Alternative: Download manually from:"
    echo "  https://github.com/superpikar/indonesia-geojson"
    echo "  Save as: indonesia-regencies.geojson"
fi

