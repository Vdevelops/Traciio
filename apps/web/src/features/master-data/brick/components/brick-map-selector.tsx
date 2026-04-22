"use client";

import dynamic from "next/dynamic";
import { useState, useCallback, useEffect, useMemo, memo, useRef } from "react";
import { motion } from "framer-motion";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { MapPin, AlertCircle } from "lucide-react";
import { useTheme } from "next-themes";
import type { FeatureCollection } from "geojson";
import { SmartTileLayer, TILE_SOURCES, LIGHT_FALLBACK_CHAIN, DARK_FALLBACK_CHAIN } from "@/components/ui/smart-tile-layer";

// Cache GeoJSON data globally to avoid reloading on every component mount
let cachedGeoJsonData: FeatureCollection | null = null;
let geoJsonLoadingPromise: Promise<FeatureCollection | null> | null = null;

// Dynamic imports untuk SSR compatibility
const MapContainer = dynamic(
  () => import("react-leaflet").then((mod) => mod.MapContainer),
  { ssr: false }
);

// TileLayer removed - using SmartTileLayer instead

const GeoJSON = dynamic(
  () => import("react-leaflet").then((mod) => mod.GeoJSON),
  { ssr: false }
);

// Leaflet CSS is imported globally in globals.css

interface RegencyInfo {
  name: string;
  province: string;
  district?: string; // Kecamatan/District
  center?: [number, number];
}

interface BrickMapSelectorProps {
  selectedRegency?: RegencyInfo;
  onRegencySelect: (regency: RegencyInfo) => void;
  existingBricks?: Array<{ regency: string; province: string }>;
}

const BrickMapSelectorComponent = memo(function BrickMapSelectorComponent({
  selectedRegency,
  onRegencySelect,
  existingBricks = [],
}: BrickMapSelectorProps) {
  const [geoJsonData, setGeoJsonData] = useState<FeatureCollection | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hoveredRegency, setHoveredRegency] = useState<RegencyInfo | null>(null);
  const { theme } = useTheme();

  // Load GeoJSON data with caching
  useEffect(() => {
    // If already cached, use it immediately
    if (cachedGeoJsonData) {
      setGeoJsonData(cachedGeoJsonData);
      setLoading(false);
      return;
    }

    // If already loading, wait for existing promise
    if (geoJsonLoadingPromise) {
      geoJsonLoadingPromise.then((data) => {
        if (data) {
          setGeoJsonData(data);
        } else {
          setError(
            "Map visualization is not available. You can still create bricks using the 'Manual Input' tab below."
          );
        }
        setLoading(false);
      }).catch(() => {
        setError("Failed to load map data. Please use manual input instead.");
        setLoading(false);
      });
      return;
    }

    // Try to load from CDN or API
    const loadGeoJSON = async (): Promise<FeatureCollection | null> => {
      try {
        // Try to load from local public folder first (most reliable)
        // Then fallback to CDN sources if local file doesn't exist
        const sources = [
          // User's GeoJSON file with accurate coordinates (priority)
          "/geojson/indonesia-provinces-simple.geojson",
          // Alternative local file name
          "/geojson/indonesia-regencies.geojson",
          // Fallback to CDN sources (if local file not available)
          "https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-city-regency.geojson",
          "https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-simple.geojson",
        ];

        let lastError: Error | null = null;
        
        for (const url of sources) {
          try {
            const controller = new AbortController();
            const timeoutId = setTimeout(() => controller.abort(), 10000); // 10 second timeout
            
            const response = await fetch(url, {
              method: "GET",
              headers: {
                "Accept": "application/json",
              },
              signal: controller.signal,
            });
            
            clearTimeout(timeoutId);
            
            if (!response.ok) {
              // Don't throw for 404, just log and continue to next source
              if (response.status === 404) {
                continue;
              }
              throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            
            const data = await response.json();
            
            // Validate GeoJSON structure
            if (data && data.type === "FeatureCollection" && Array.isArray(data.features) && data.features.length > 0) {
              // Cache the data
              cachedGeoJsonData = data;
              return data; // Success, return data
            } else {
              throw new Error("Invalid GeoJSON structure");
            }
          } catch (err) {
            // Only log warning for non-abort errors, don't throw yet - try next source
            if (err instanceof Error && err.name !== "AbortError") {
              if (!lastError) {
                lastError = err;
              }
            } else if (err instanceof Error && err.name === "AbortError") {
              if (!lastError) {
                lastError = new Error("Request timeout");
              }
            }
            // Continue to next source
          }
        }
        
        // If all sources failed, show error but don't crash
        // Don't log as error in console - this is expected if GeoJSON file is not provided
        console.info("GeoJSON map data not available. Using manual input mode.");
        return null;
      } catch (err) {
        return null;
      }
    };

    // Create and cache the promise
    geoJsonLoadingPromise = loadGeoJSON();
    
    geoJsonLoadingPromise.then((data) => {
      if (data) {
        setGeoJsonData(data);
      } else {
        setError(
          "Map visualization is not available. You can still create bricks using the 'Manual Input' tab below."
        );
      }
      setLoading(false);
    }).catch(() => {
      setError("Failed to load map data. Please use manual input instead.");
      setLoading(false);
    });
  }, []);

  // Helper function to extract regency, province, and district from GeoJSON properties
  // Supports multiple GeoJSON formats:
  // 1. Standard format: { name, province, district } or { regency, province, district }
  // 2. Indonesian format: { WADMKD, WADMPR, WADMKK, WADMKC, KDPKAB } 
  //    - WADMKD: Nama Kecamatan/Desa (district/village) - level lebih rendah
  //    - WADMPR: Nama Provinsi (province)
  //    - WADMKK: Nama Kabupaten/Kota (regency/city) - level yang kita butuhkan
  //    - WADMKC: Nama Kecamatan (district/sub-district) - untuk field district
  //    - KDPKAB: Kode Kabupaten (regency code) - bisa digunakan untuk grouping
  const extractRegencyAndProvince = useCallback((properties: any) => {
    // Try Indonesian format first (format dari file 38_Provinsi_Indonesia_Kabupaten.json)
    if (properties?.WADMPR) {
      // File ini berisi data level desa/kecamatan, tapi kita butuh level kabupaten/kota
      // WADMKK adalah nama kabupaten/kota (ini yang kita butuhkan untuk brick)
      // WADMKD adalah nama kecamatan/desa (level lebih rendah)
      // WADMKC adalah nama kecamatan (untuk field district)
      
      let regencyName = "";
      let district = "";
      
      // Priority: WADMKK (Kabupaten/Kota) - ini adalah level yang benar untuk brick
      if (properties.WADMKK) {
        regencyName = properties.WADMKK;
      } 
      // Jika tidak ada WADMKK, mungkin data sudah di level kabupaten
      // Coba gunakan WADMKD jika tidak ada WADMKK
      else if (properties.WADMKD) {
        // Check jika WADMKD mengandung "Kota" atau "Kabupaten" - mungkin ini adalah nama kabupaten
        const wadmkd = String(properties.WADMKD);
        if (wadmkd.includes("Kota") || wadmkd.includes("Kabupaten")) {
          regencyName = wadmkd;
        } else {
          // Jika tidak, mungkin perlu grouping berdasarkan KDPKAB
          // Untuk sekarang, gunakan WADMKD sebagai fallback
          regencyName = properties.WADMKD;
        }
      } 
      // Fallback ke NAMOBJ
      else if (properties.NAMOBJ) {
        regencyName = properties.NAMOBJ;
      }
      
      // Extract district from WADMKC (Kecamatan) or WADMKD (if WADMKK exists)
      if (properties.WADMKC) {
        district = properties.WADMKC;
      } else if (properties.WADMKD && properties.WADMKK) {
        // Jika ada WADMKK, maka WADMKD adalah kecamatan
        district = properties.WADMKD;
      }
      
      const province = properties.WADMPR || "";
      return { regencyName, province, district };
    }
    
    // Fallback to standard format
    const regencyName = properties?.name || properties?.regency || properties?.NAMOBJ || "";
    const province = properties?.province || properties?.WADMPR || "";
    const district = properties?.district || properties?.WADMKC || "";
    return { regencyName, province, district };
  }, []);

  // Helper function to get center of geometry
  const getCenter = (geometry: any): [number, number] => {
    if (geometry.type === "Polygon" && geometry.coordinates?.[0]) {
      const coords = geometry.coordinates[0];
      let latSum = 0;
      let lngSum = 0;
      for (const coord of coords) {
        lngSum += coord[0];
        latSum += coord[1];
      }
      return [latSum / coords.length, lngSum / coords.length];
    }
    if (geometry.type === "MultiPolygon" && geometry.coordinates?.[0]?.[0]) {
      // For MultiPolygon, use the first polygon
      const coords = geometry.coordinates[0][0];
      let latSum = 0;
      let lngSum = 0;
      for (const coord of coords) {
        lngSum += coord[0];
        latSum += coord[1];
      }
      return [latSum / coords.length, lngSum / coords.length];
    }
    return [-2.5, 118.0]; // Center of Indonesia
  };

  // Memoize existing bricks map for faster lookup
  const existingBricksMap = useMemo(() => {
    const map = new Map<string, boolean>();
    existingBricks.forEach((brick) => {
      const key = `${brick.regency.toLowerCase()}_${brick.province.toLowerCase()}`;
      map.set(key, true);
    });
    return map;
  }, [existingBricks]);

  // Check if regency already has a brick (optimized with Map lookup)
  const hasBrick = useCallback(
    (regencyName: string, province: string) => {
      const key = `${regencyName.toLowerCase()}_${province.toLowerCase()}`;
      return existingBricksMap.has(key);
    },
    [existingBricksMap]
  );

  // Memoize selected regency key for faster comparison
  const selectedRegencyKey = useMemo(() => {
    if (!selectedRegency) return null;
    return `${selectedRegency.name.toLowerCase()}_${selectedRegency.province.toLowerCase()}`;
  }, [selectedRegency]);

  // Style function for GeoJSON features (optimized with memoization)
  const style = useCallback(
    (feature: any) => {
      const { regencyName, province } = extractRegencyAndProvince(feature.properties);
      const isBrick = hasBrick(regencyName, province);
      const featureKey = `${regencyName.toLowerCase()}_${province.toLowerCase()}`;
      const isSelected = selectedRegencyKey === featureKey;

      if (isSelected) {
        return {
          fillColor: "#10B981", // Green for selected
          fillOpacity: 0.7,
          color: "#FFFFFF",
          weight: 2,
          opacity: 1,
        };
      }

      if (isBrick) {
        return {
          fillColor: "#3B82F6", // Blue for existing brick
          fillOpacity: 0.5,
          color: "#FFFFFF",
          weight: 1,
          opacity: 0.8,
        };
      }

      return {
        fillColor: "#9CA3AF", // Gray for available
        fillOpacity: 0.4,
        color: "#FFFFFF",
        weight: 1,
        opacity: 0.6,
      };
    },
    [hasBrick, selectedRegencyKey, extractRegencyAndProvince]
  );

  // Debounce hover to prevent flickering (same technique as brick-map-preview)
  const hoverTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const currentHoveredLayerRef = useRef<any>(null); // Track currently hovered layer

  // Handle feature events (optimized with smooth hover like brick-map-preview)
  const onEachFeature = useCallback(
    (feature: any, layer: any) => {
      const { regencyName, province, district } = extractRegencyAndProvince(feature.properties);
      const isBrick = hasBrick(regencyName, province);
      const featureStyle = style(feature);

      // Store feature reference in layer for style reset
      layer.feature = feature;

      // Hover events with smooth handling (same as brick-map-preview)
      layer.on({
        mouseover: (e: any) => {
          // Clear any pending hover timeout
          if (hoverTimeoutRef.current) {
            clearTimeout(hoverTimeoutRef.current);
          }
          
          // Reset previous hovered layer style if exists (prevents flickering)
          if (currentHoveredLayerRef.current && currentHoveredLayerRef.current !== layer) {
            const prevFeature = currentHoveredLayerRef.current.feature || feature;
            currentHoveredLayerRef.current.setStyle(style(prevFeature));
          }
          
          // Set current hovered layer
          currentHoveredLayerRef.current = layer;
          
          // Debounce hover to prevent flickering
          hoverTimeoutRef.current = setTimeout(() => {
            const center = getCenter(feature.geometry);
            setHoveredRegency({ name: regencyName, province, district, center });
            
            // Apply hover style smoothly
            if (isBrick) {
              layer.setStyle({
                fillColor: "#2563EB", // Blue for existing brick
                fillOpacity: 0.85,
                color: "#FFFFFF",
                weight: 2.5,
                opacity: 1,
              });
            } else {
              // Available area style
              layer.setStyle({
                fillColor: "#FBBF24", // Yellow
                fillOpacity: 0.7,
                color: "#FFFFFF",
                weight: 2,
                opacity: 1,
              });
            }
            
            // Open popup with slight delay for smoother UX
            setTimeout(() => {
              if (layer && layer.isPopupOpen && !layer.isPopupOpen()) {
                layer.openPopup();
              }
            }, 100);
          }, 50); // Small delay to prevent rapid state changes
        },
        mouseout: (e: any) => {
          // Clear hover timeout
          if (hoverTimeoutRef.current) {
            clearTimeout(hoverTimeoutRef.current);
            hoverTimeoutRef.current = null;
          }
          
          // Only reset if this is the currently hovered layer
          if (currentHoveredLayerRef.current === layer) {
            // Reset hover state
            setHoveredRegency(null);
            
            // Reset style smoothly
            layer.setStyle(featureStyle);
            
            // Close popup only if not selected
            if (!selectedRegency || 
                selectedRegency.name.toLowerCase() !== regencyName.toLowerCase() ||
                selectedRegency.province.toLowerCase() !== province.toLowerCase()) {
              layer.closePopup();
            }
            
            currentHoveredLayerRef.current = null;
          }
        },
        click: (e: any) => {
          // Clear hover timeout on click
          if (hoverTimeoutRef.current) {
            clearTimeout(hoverTimeoutRef.current);
            hoverTimeoutRef.current = null;
          }
          
          if (!isBrick) {
            onRegencySelect({ name: regencyName, province, district });
            
            // Keep popup open for selected regency
            layer.openPopup();
          }
        },
      });

      // Add popup (memoized content)
      const popupContent = `
        <div style="text-align: center; min-width: 180px; padding: 8px;">
          <strong style="font-size: 13px; display: block; margin-bottom: 4px; font-weight: 600; color: #111827;">${regencyName || "Unknown"}</strong>
          <small style="display: block; color: #6B7280; margin-bottom: 8px; font-size: 11px;">${province || "Unknown Province"}</small>
          ${district ? `<small style="display: block; color: #9CA3AF; font-size: 10px; margin-bottom: 8px;">${district}</small>` : ""}
          ${isBrick ? '<span style="color: #3B82F6; font-size: 11px; font-weight: 500;">Already a Brick</span>' : '<span style="color: #10B981; font-size: 11px; font-weight: 500;">Available</span>'}
        </div>
      `;
      layer.bindPopup(popupContent, {
        closeOnClick: false,
        autoClose: false,
        closeOnEscapeKey: true,
        className: "brick-selector-popup",
        maxWidth: 220,
        offset: [0, -5],
      });
    },
    [hasBrick, style, onRegencySelect, extractRegencyAndProvince, selectedRegency]
  );

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Select Regency/City from Map</CardTitle>
          <CardDescription>
            Click on a region in the map to select it as a brick
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Skeleton className="h-[400px] w-full" />
        </CardContent>
      </Card>
    );
  }

  if (error || !geoJsonData) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Map Selection</CardTitle>
          <CardDescription>
            Map visualization is currently unavailable. You can still create bricks using manual input.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-3">
            <div className="flex items-start gap-3 text-blue-600 dark:text-blue-400 p-4 border border-blue-200 dark:border-blue-800 rounded-md bg-blue-50 dark:bg-blue-950/20">
              <AlertCircle className="h-5 w-5 mt-0.5 flex-shrink-0" />
              <div className="flex-1">
                <p className="font-medium mb-1">Map data not configured</p>
                <p className="text-sm text-blue-700 dark:text-blue-300">
                  {error || "To enable map visualization, place a GeoJSON file at /public/geojson/indonesia-regencies.geojson. For now, please use the 'Manual Input' tab."}
                </p>
              </div>
            </div>
            <div className="text-sm text-muted-foreground p-3 bg-muted rounded-md">
              <p className="font-medium mb-2">How to proceed:</p>
              <ol className="list-decimal list-inside space-y-1.5 ml-2">
                <li>Switch to the &quot;Manual Input&quot; tab above</li>
                <li>Fill in the <strong>Province</strong> and <strong>Regency/City</strong> fields manually</li>
                <li>Complete the rest of the form and submit</li>
              </ol>
              <div className="mt-3 pt-3 border-t">
                <p className="font-medium mb-1 text-xs">Optional: Enable Map Visualization</p>
                <p className="text-xs text-muted-foreground">
                  Download Indonesia GeoJSON from GitHub and place it at: <code className="text-xs bg-background px-1 py-0.5 rounded">apps/web/public/geojson/indonesia-regencies.geojson</code>
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-lg font-medium">Select Regency/City from Map</CardTitle>
            <CardDescription className="text-xs text-muted-foreground mt-1">
              Click on a region in the map to select it as a brick. Hover to see details.
            </CardDescription>
          </div>
          {hoveredRegency && (
            <motion.div
              initial={{ opacity: 0, y: -5 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -5 }}
              transition={{ duration: 0.2 }}
              className="flex items-center gap-2"
            >
              <Badge variant="outline" className="text-xs font-medium">
                {hoveredRegency.name}
              </Badge>
            </motion.div>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <div className="h-[400px] w-full rounded-lg border border-border/50 overflow-hidden relative bg-muted/30">
          <MapContainer
            center={[-2.5, 118.0]} // Center of Indonesia
            zoom={5}
            style={{ height: "100%", width: "100%", zIndex: 0 }}
            scrollWheelZoom={true}
            zoomControl={true}
          >
            {/* SmartTileLayer with auto-retry and fallback */}
            <SmartTileLayer
              key={`smart-tile-${theme === "dark" ? "dark" : "light"}`}
              source={theme === "dark" ? TILE_SOURCES.cartoDark : TILE_SOURCES.cartoLight}
              fallbackSources={theme === "dark" ? DARK_FALLBACK_CHAIN : LIGHT_FALLBACK_CHAIN}
              maxRetries={2}
              retryDelay={200}
              priorityMode="viewport"
            />
            {geoJsonData && (
              <GeoJSON
                data={geoJsonData}
                style={style}
                onEachFeature={onEachFeature}
              />
            )}
          </MapContainer>
          
          {/* Custom styles for smooth popup (same as brick-map-preview) */}
          <style jsx global>{`
            .brick-selector-popup .leaflet-popup-content-wrapper {
              border-radius: 8px;
              box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
              padding: 0;
              border: 1px solid rgba(0, 0, 0, 0.1);
            }
            .brick-selector-popup .leaflet-popup-content {
              margin: 0;
              padding: 0;
            }
            .brick-selector-popup .leaflet-popup-tip {
              background: white;
              border: 1px solid rgba(0, 0, 0, 0.1);
            }
            .brick-selector-popup .leaflet-popup-close-button {
              color: #6B7280;
              font-size: 18px;
              padding: 4px 8px;
            }
            .brick-selector-popup .leaflet-popup-close-button:hover {
              color: #111827;
            }
          `}</style>
        </div>
        <div className="mt-3 pt-3 border-t border-border/50 flex flex-wrap gap-x-6 gap-y-2 text-xs">
          <div className="flex items-center gap-1.5">
            <div className="w-3 h-3 rounded bg-gray-400"></div>
            <span className="text-muted-foreground">Available</span>
          </div>
          <div className="flex items-center gap-1.5">
            <div className="w-3 h-3 rounded bg-blue-500"></div>
            <span className="text-muted-foreground">Existing Brick</span>
          </div>
          <div className="flex items-center gap-1.5">
            <div className="w-3 h-3 rounded bg-green-500"></div>
            <span className="text-muted-foreground">Selected</span>
          </div>
          <div className="flex items-center gap-1.5">
            <div className="w-3 h-3 rounded bg-yellow-400"></div>
            <span className="text-muted-foreground">Hover</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
});

// Export with loading fallback
export const BrickMapSelector = dynamic(
  () => Promise.resolve(BrickMapSelectorComponent),
  {
    ssr: false,
    loading: () => (
      <Card>
        <CardHeader>
          <CardTitle>Select Regency/City from Map</CardTitle>
        </CardHeader>
        <CardContent>
          <Skeleton className="h-[400px] w-full" />
        </CardContent>
      </Card>
    ),
  }
);

