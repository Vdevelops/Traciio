"use client";

import dynamic from "next/dynamic";
import { useState, useCallback, useEffect, useMemo, useRef } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { Layers, Grid3X3 } from "lucide-react";
import { useTheme } from "next-themes";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import "leaflet/dist/leaflet.css";

// Fix for default marker icons in Next.js
import L from "leaflet";
if (typeof window !== "undefined") {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  delete (L.Icon.Default.prototype as any)._getIconUrl;
  L.Icon.Default.mergeOptions({
    iconRetinaUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png",
    iconUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png",
    shadowUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png",
  });
}

import type { FeatureCollection } from "geojson";
import type { Brick } from "../types";
import { useBrickPerformanceList } from "../hooks/useBrickAnalytics";
import { useBrickPeriodQueryParams } from "../hooks/useBrickPeriodQueryParams";
import { useTranslations } from "next-intl";
import type { LatLngExpression } from "leaflet";
import { SmartTileLayer, TILE_SOURCES, LIGHT_FALLBACK_CHAIN, DARK_FALLBACK_CHAIN, type TileSource } from "@/components/ui/smart-tile-layer";
import type { BrickPerformanceMetrics } from "../types/analytics";
import { useAccountsForMap } from "@/features/sales-crm/account-management/hooks/useAccounts";
import type { Account } from "@/features/sales-crm/account-management/types";
import { formatCurrency } from "@/lib/utils";

// Dynamic imports for SSR compatibility
const MapContainer = dynamic(
  () => import("react-leaflet").then((mod) => mod.MapContainer),
  { ssr: false }
);

const GeoJSON = dynamic(
  () => import("react-leaflet").then((mod) => mod.GeoJSON),
  { ssr: false }
);

const MarkerCluster = dynamic(
  () => Promise.all([
    import("react-leaflet"),
    import("leaflet"),
    import("leaflet.markercluster"),
    import("react"),
  ]).then(([reactLeaflet, leaflet, _, react]) => {
    const { useMap } = reactLeaflet;
    const L = leaflet.default;

    return function MarkerClusterInner({
      accounts,
      selectedBrickId,
      labels,
    }: {
      accounts: Account[];
      selectedBrickId?: string | null;
      labels: {
        noCategory: string;
        city: string;
        coordinates: string;
      };
    }) {
      const map = useMap();
      const clusterRef = react.useRef<any>(null);

      const createMarkerIcon = react.useCallback((account: Account) => {
        const isSelectedBrick = selectedBrickId != null && account.brick_id === selectedBrickId;
        const size = isSelectedBrick ? 36 : 28;
        const border = isSelectedBrick ? 4 : 3;
        const color = isSelectedBrick ? "#F59E0B" : "#16A34A";

        return L.divIcon({
          className: "brick-account-marker",
          html: `
            <div style="
              background:${color};
              color:white;
              width:${size}px;
              height:${size}px;
              border-radius:9999px;
              display:flex;
              align-items:center;
              justify-content:center;
              border:${border}px solid white;
              box-shadow:0 6px 16px rgba(0,0,0,0.22);
            ">
              <svg xmlns="http://www.w3.org/2000/svg" width="${isSelectedBrick ? 16 : 13}" height="${isSelectedBrick ? 16 : 13}" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 22s8-4 8-11a8 8 0 1 0-16 0c0 7 8 11 8 11z"></path>
                <circle cx="12" cy="11" r="2.5"></circle>
              </svg>
            </div>
          `,
          iconSize: [size, size],
          iconAnchor: [size / 2, size / 2],
          popupAnchor: [0, -size / 2],
        });
      }, [selectedBrickId]);

      react.useEffect(() => {
        if (!map) return;

        if (clusterRef.current) {
          map.removeLayer(clusterRef.current);
          clusterRef.current = null;
        }

        const markerClusterGroup = (L as any).markerClusterGroup;
        if (!markerClusterGroup) return;

        const cluster = markerClusterGroup({
          maxClusterRadius: 48,
          spiderfyOnMaxZoom: true,
          showCoverageOnHover: false,
          zoomToBoundsOnClick: true,
          disableClusteringAtZoom: 14,
          chunkedLoading: true,
          iconCreateFunction: (clusterObj: { getChildCount: () => number }) => {
            const count = clusterObj.getChildCount();
            const size = count >= 100 ? 52 : count >= 10 ? 44 : 36;
            const fontSize = count >= 100 ? 16 : count >= 10 ? 14 : 12;

            return L.divIcon({
              className: "brick-account-cluster",
              html: `<div style="
                background:linear-gradient(135deg, #16A34A 0%, #15803D 100%);
                color:white;
                width:${size}px;
                height:${size}px;
                border-radius:9999px;
                display:flex;
                align-items:center;
                justify-content:center;
                font-weight:700;
                font-size:${fontSize}px;
                border:3px solid white;
                box-shadow:0 4px 12px rgba(22,163,74,0.35);
              ">${count}</div>`,
              iconSize: [size, size],
              iconAnchor: [size / 2, size / 2],
            });
          },
        });

        accounts.forEach((account) => {
          if (account.latitude == null || account.longitude == null) return;

          const marker = L.marker([account.latitude, account.longitude], {
            icon: createMarkerIcon(account),
          });

          const categoryName = account.category?.name || labels.noCategory;
          const cityText = account.city?.trim() || "-";
          marker.bindPopup(
            `<div class="brick-account-popup">
              <div class="brick-account-popup__title">${account.name}</div>
              <div class="brick-account-popup__meta">${categoryName}</div>
              <div class="brick-account-popup__row"><span>${labels.city}</span><strong>${cityText}</strong></div>
              <div class="brick-account-popup__row"><span>${labels.coordinates}</span><strong>${account.latitude.toFixed(6)}, ${account.longitude.toFixed(6)}</strong></div>
            </div>`,
            { maxWidth: 320, closeButton: true, autoPan: true }
          );

          cluster.addLayer(marker);
        });

        map.addLayer(cluster);
        clusterRef.current = cluster;

        return () => {
          if (clusterRef.current) {
            map.removeLayer(clusterRef.current);
            clusterRef.current = null;
          }
        };
      }, [accounts, createMarkerIcon, labels, map]);

      return null;
    };
  }),
  { ssr: false }
);

// Map styles using shared TILE_SOURCES for consistency
type MapStyle = "light" | "dark" | "satellite" | "streets";

const mapStyles: Record<MapStyle, { name: string; source: TileSource }> = {
  light: {
    name: "Light",
    source: TILE_SOURCES.cartoLight,
  },
  dark: {
    name: "Dark",
    source: TILE_SOURCES.cartoDark,
  },
  satellite: {
    name: "Satellite",
    source: TILE_SOURCES.esriSatellite,
  },
  streets: {
    name: "Streets",
    source: TILE_SOURCES.openStreetMap,
  },
};

interface BrickMapFullProps {
  bricks: Brick[];
  onBrickClick?: (brick: Brick) => void;
  onCreateBrickFromMap?: (regency: string, province: string) => void;
  searchQuery?: string;
  selectedBrickId?: string | null;
  isLoading?: boolean;
  colorBy?: "revenue" | "achievement";
}

// Component to control map view
const MapFocus = dynamic(
  () => Promise.all([
    import("react-leaflet"),
    import("react"),
  ]).then(([reactLeaflet, react]) => {
    const MapFocusInner = ({ center, zoom, shouldAnimate }: { center: LatLngExpression; zoom: number; shouldAnimate: boolean }) => {
      const map = reactLeaflet.useMap();
      const prevCenterRef = react.useRef<LatLngExpression | null>(null);
      const prevZoomRef = react.useRef<number | null>(null);

      react.useEffect(() => {
        if (!map) return;

        const centerStr = Array.isArray(center) ? `${center[0]},${center[1]}` : JSON.stringify(center);
        const prevCenterStr = prevCenterRef.current
          ? (Array.isArray(prevCenterRef.current) ? `${prevCenterRef.current[0]},${prevCenterRef.current[1]}` : JSON.stringify(prevCenterRef.current))
          : null;
        const centerChanged = prevCenterStr !== centerStr;
        const zoomChanged = prevZoomRef.current !== zoom;

        if (!centerChanged && !zoomChanged) return;

        if (shouldAnimate) {
          const timeoutId = setTimeout(() => {
            if (typeof (map as any).flyTo === "function") {
              (map as any).flyTo(center, zoom, { duration: 1.8, easeLinearity: 0.15 });
            } else {
              map.setView(center, zoom, { animate: true, duration: 1.8 });
            }
            prevCenterRef.current = center;
            prevZoomRef.current = zoom;
          }, 50);
          return () => clearTimeout(timeoutId);
        }

        map.setView(center, zoom, { animate: false });
        prevCenterRef.current = center;
        prevZoomRef.current = zoom;
      }, [map, center, zoom, shouldAnimate]);

      return null;
    };
    return { default: MapFocusInner };
  }),
  { ssr: false }
);

/**
 * Compute a fill color for a brick region using heatmap-style intensity scaling.
 * Uses the project's primary color (orange, hue ~36) scaled from light to dark.
 * Higher metric values produce deeper/darker primary tones.
 */
function getHeatmapColor(value: number, min: number, max: number): string {
  if (max <= min) return "hsl(36, 70%, 90%)";

  // Normalize to 0..1
  const ratio = Math.min(Math.max((value - min) / (max - min), 0), 1);

  // HSL interpolation: light primary (h=36 s=70% l=90%) -> dark primary (h=36 s=100% l=35%)
  const lightness = 90 - ratio * 55; // 90 -> 35
  const saturation = 70 + ratio * 30; // 70 -> 100
  return `hsl(36, ${saturation}%, ${lightness}%)`;
}

/** Returns a color based on achievement percentage thresholds */
function getAchievementColor(achievement: number): string {
  if (achievement >= 80) return "#16A34A";
  if (achievement >= 50) return "#2563EB";
  return "#DC2626";
}

function BrickMapFullComponent({
  bricks,
  onBrickClick,
  onCreateBrickFromMap,
  searchQuery,
  selectedBrickId,
  isLoading,
  colorBy = "revenue",
}: Readonly<BrickMapFullProps>) {
  const [geoJsonData, setGeoJsonData] = useState<FeatureCollection | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedBrick, setSelectedBrick] = useState<Brick | null>(null);
  const [mapCenter, setMapCenter] = useState<LatLngExpression>([-2.5, 118.0]);
  const [mapZoom, setMapZoom] = useState(5);
  const [shouldAnimate, setShouldAnimate] = useState(false);
  // Counter incremented every time the GeoJSON layer mounts (or remounts after key change)
  const [geoJsonVersion, setGeoJsonVersion] = useState(0);
  const [mapStyle, setMapStyle] = useState<MapStyle>("light");
  const brickToFeatureMapRef = useRef<Map<string, any>>(new Map());
  // Stores actual Leaflet layer objects so we can apply highlight animation per brick
  const layerMapRef = useRef<Map<string, any>>(new Map());
  const t = useTranslations("brickManagement.list");
  const { theme, resolvedTheme } = useTheme();
  const { mode, periodStart, periodEnd } = useBrickPeriodQueryParams();
  const { data: accountsMapData } = useAccountsForMap({ status: "active" });

  // Key that forces GeoJSON layer to remount (re-run onEachFeature) whenever the bricks list changes.
  // This ensures newly created or deleted bricks are immediately reflected on the map.
  const geoJsonKey = useMemo(
    () => bricks.map((b) => b.id).join(",") || "empty",
    [bricks]
  );

  // Clear stale layer references whenever the brick list changes so the highlight
  // effect always binds to the freshly-mounted Leaflet layers.
  useEffect(() => {
    layerMapRef.current.clear();
    brickToFeatureMapRef.current.clear();
  }, [geoJsonKey]);

  // Period query string used for navigation links
  const _periodQueryString = useMemo(() => {
    const params = new URLSearchParams();
    params.set("period_mode", mode);
    params.set("period_start", periodStart);
    params.set("period_end", periodEnd);
    return params.toString();
  }, [mode, periodStart, periodEnd]);

  // Auto-select map style based on theme
  useEffect(() => {
    const currentTheme = resolvedTheme || theme;
    setMapStyle(currentTheme === "dark" ? "dark" : "light");
  }, [theme, resolvedTheme]);

  // Batch-fetch performance data for ALL bricks in a single request.
  // This is the ONLY API call for the map. All hover/click interactions use this cache.
  const brickIdsForPerformance = useMemo(() => bricks.map((b) => b.id), [bricks]);
  const requestPeriodStart = periodStart || undefined;
  const requestPeriodEnd = periodEnd || undefined;
  const { data: performanceListData } = useBrickPerformanceList(
    brickIdsForPerformance,
    requestPeriodStart,
    requestPeriodEnd,
    {
      staleTime: 5 * 60 * 1000,
      refetchOnWindowFocus: false,
    }
  );

  // Build performance lookup map and compute min/max for heatmap scaling
  const { performanceByBrickId, metricMin, metricMax } = useMemo(() => {
    const map = new Map<string, BrickPerformanceMetrics>();
    let min = Infinity;
    let max = -Infinity;

    for (const metrics of performanceListData?.data ?? []) {
      if (!metrics?.brick_id) continue;
      map.set(metrics.brick_id, metrics);
      const val = colorBy === "revenue" ? (metrics.total_revenue ?? 0) : (metrics.achievement_percentage ?? 0);
      if (val < min) min = val;
      if (val > max) max = val;
    }

    if (!Number.isFinite(min)) min = 0;
    if (!Number.isFinite(max)) max = 0;

    return { performanceByBrickId: map, metricMin: min, metricMax: max };
  }, [performanceListData, colorBy]);

  // Filter bricks based on search query
  const filteredBricks = useMemo(() => {
    if (!searchQuery) return bricks;
    const query = searchQuery.toLowerCase().trim();
    return bricks.filter(
      (brick) =>
        brick.name.toLowerCase().includes(query) ||
        brick.code.toLowerCase().includes(query) ||
        brick.province.toLowerCase().includes(query) ||
        brick.regency.toLowerCase().includes(query)
    );
  }, [bricks, searchQuery]);

  const accountsWithCoordinates = useMemo(() => {
    const brickIds = new Set(bricks.map((brick) => brick.id));
    return (accountsMapData?.data ?? []).filter((account) =>
      account.latitude != null &&
      account.longitude != null &&
      account.brick_id != null &&
      brickIds.has(account.brick_id)
    );
  }, [accountsMapData, bricks]);

  // Create map of existing bricks for quick lookup
  const bricksMap = useMemo(() => {
    const map = new Map<string, Brick>();
    bricks.forEach((brick) => {
      const regency = brick.regency.toLowerCase().trim();
      const province = brick.province.toLowerCase().trim();
      const key = `${regency}_${province}`;
      map.set(key, brick);

      const regencyWithoutPrefix = regency
        .replace(/^(kota|kabupaten)\s+/i, "")
        .replace(/\s+(kota|kabupaten)$/i, "")
        .trim();
      if (regencyWithoutPrefix !== regency) {
        map.set(`${regencyWithoutPrefix}_${province}`, brick);
      }
      if (!regency.includes("kota") && !regency.includes("kabupaten")) {
        map.set(`kota ${regency}_${province}`, brick);
        map.set(`kabupaten ${regency}_${province}`, brick);
      }
    });
    return map;
  }, [bricks]);

  // Auto-focus to selected brick
  useEffect(() => {
    if (!geoJsonVersion || !geoJsonData) return;

    setShouldAnimate(false);

    const findFeatureAndZoom = (brick: Brick, retryCount = 0) => {
      const feature = brickToFeatureMapRef.current.get(brick.id);

      if (feature?.geometry) {
        const center = getCenter(feature.geometry);
        setSelectedBrick(brick);
        requestAnimationFrame(() => {
          setMapCenter(center);
          setMapZoom(10);
          setTimeout(() => setShouldAnimate(true), 50);
        });
      } else if (retryCount < 10) {
        setTimeout(() => findFeatureAndZoom(brick, retryCount + 1), 100 * (retryCount + 1));
      }
    };

    if (selectedBrickId && filteredBricks.length > 0) {
      const brick = filteredBricks.find((b) => b.id === selectedBrickId);
      if (brick) findFeatureAndZoom(brick);
    } else if (searchQuery && filteredBricks.length === 1) {
      findFeatureAndZoom(filteredBricks[0]);
    } else if (!searchQuery && !selectedBrickId) {
      setSelectedBrick(null);
      setMapCenter([-2.5, 118.0]);
      setMapZoom(5);
      setShouldAnimate(false);
    }
  }, [searchQuery, filteredBricks, selectedBrickId, geoJsonVersion, geoJsonData]);

  // Load GeoJSON data
  useEffect(() => {
    const loadGeoJSON = async () => {
      try {
        const sources = [
          "/geojson/indonesia-provinces-simple.geojson",
          "/geojson/indonesia-regencies.geojson",
          "https://raw.githubusercontent.com/superpikar/indonesia-geojson/master/indonesia-province-city-regency.geojson",
        ];

        for (const url of sources) {
          try {
            const controller = new AbortController();
            const timeoutId = setTimeout(() => controller.abort(), 10000);
            const response = await fetch(url, { signal: controller.signal });
            clearTimeout(timeoutId);

            if (!response.ok) {
              if (response.status === 404) continue;
              throw new Error(`HTTP ${response.status}`);
            }

            const data = await response.json();
            if (data.type === "FeatureCollection" && Array.isArray(data.features)) {
              setGeoJsonData(data);
              setGeoJsonVersion(0); // reset; will be incremented when layer fires 'add'
              setError(null);
              setLoading(false);
              return;
            }
          } catch {
            continue;
          }
        }
        throw new Error("No GeoJSON source available");
      } catch {
        setError("Map data is not available");
        setLoading(false);
      }
    };
    loadGeoJSON();
  }, []);

  const extractRegencyAndProvince = useCallback((properties: any) => {
    if (properties?.WADMPR) {
      let regencyName = "";
      if (properties.WADMKK) {
        regencyName = String(properties.WADMKK).trim();
      } else if (properties.WADMKD) {
        const wadmkd = String(properties.WADMKD).trim();
        if (wadmkd.includes("Kota") || wadmkd.includes("Kabupaten")) regencyName = wadmkd;
      } else if (properties.NAMOBJ) {
        regencyName = String(properties.NAMOBJ).trim();
      }
      const province = String(properties.WADMPR || "").trim();
      return { regencyName, province };
    }
    const regencyName = String(properties?.name || properties?.regency || properties?.NAMOBJ || "").trim();
    const province = String(properties?.province || properties?.WADMPR || "").trim();
    return { regencyName, province };
  }, []);

  const getCenter = (geometry: any): [number, number] => {
    let coords = null;
    if (geometry.type === "Polygon") {
      coords = geometry.coordinates?.[0];
    } else if (geometry.type === "MultiPolygon") {
      coords = geometry.coordinates?.[0]?.[0];
    }

    if (!coords) return [-2.5, 118.0];

    let latSum = 0;
    let lngSum = 0;
    for (const coord of coords) {
      lngSum += coord[0];
      latSum += coord[1];
    }
    return [latSum / coords.length, lngSum / coords.length];
  };

  const hasBrick = useCallback(
    (regencyName: string, province: string): Brick | null => {
      if (!regencyName || !province) return null;

      const normalizedRegency = regencyName.toLowerCase().trim();
      const normalizedProvince = province.toLowerCase().trim();

      let key = `${normalizedRegency}_${normalizedProvince}`;
      let brick = bricksMap.get(key);
      if (brick) return brick;

      const regencyWithoutPrefix = normalizedRegency
        .replace(/^(kota|kabupaten)\s+/i, "")
        .replace(/\s+(kota|kabupaten)$/i, "")
        .trim();
      if (regencyWithoutPrefix !== normalizedRegency) {
        brick = bricksMap.get(`${regencyWithoutPrefix}_${normalizedProvince}`);
        if (brick) return brick;
      }

      if (!normalizedRegency.includes("kota") && !normalizedRegency.includes("kabupaten")) {
        brick = bricksMap.get(`kota ${normalizedRegency}_${normalizedProvince}`);
        if (brick) return brick;
        brick = bricksMap.get(`kabupaten ${normalizedRegency}_${normalizedProvince}`);
        if (brick) return brick;
      }

      return null;
    },
    [bricksMap]
  );

  // GeoJSON style with heatmap color intensity
  const style = useCallback(
    (feature: any) => {
      const { regencyName, province } = extractRegencyAndProvince(feature.properties);
      const brick = hasBrick(regencyName, province);

      if (brick) {
        const isSelected = selectedBrick?.id === brick.id;
        const metrics = performanceByBrickId.get(brick.id);

        if (brick.status !== "active") {
          return {
            fillColor: "#9CA3AF",
            fillOpacity: isSelected ? 0.65 : 0.35,
            color: "#FFFFFF",
            weight: isSelected ? 2 : 1.2,
            opacity: 1,
          };
        }

        // Heatmap: color intensity based on revenue or achievement
        let metricValue = 0;
        if (metrics) {
          metricValue = colorBy === "revenue" ? (metrics.total_revenue ?? 0) : (metrics.achievement_percentage ?? 0);
        }
        const heatColor = getHeatmapColor(metricValue, metricMin, metricMax);

        return {
          fillColor: isSelected ? "#F59E0B" : heatColor,
          fillOpacity: isSelected ? 0.7 : 0.6,
          color: "#FFFFFF",
          weight: isSelected ? 2.5 : 1.2,
          opacity: 1,
        };
      }

      // Unassigned area - subtle styling
      return {
        fillColor: "#F3F4F6",
        fillOpacity: 0.15,
        color: "#D1D5DB",
        weight: 0.8,
        opacity: 0.4,
      };
    },
    [hasBrick, extractRegencyAndProvince, selectedBrick, performanceByBrickId, metricMin, metricMax, colorBy]
  );

  // Apply pulsing highlight to the currently selected brick layer element
  useEffect(() => {
    // Remove highlight from all layers first
    layerMapRef.current.forEach((layer) => {
      const el: SVGElement | undefined = layer.getElement?.();
      if (el) {
        el.classList.remove("brick-selected-pulse");
      }
    });

    if (selectedBrick) {
      const layer = layerMapRef.current.get(selectedBrick.id);
      const el: SVGElement | undefined = layer?.getElement?.();
      if (el) {
        el.classList.add("brick-selected-pulse");
      }
    }
  }, [selectedBrick, geoJsonVersion]);

  const onEachFeature = useCallback(
    (feature: any, layer: any) => {
      const { regencyName, province } = extractRegencyAndProvince(feature.properties);
      const brick = hasBrick(regencyName, province);

      if (brick) {
        brickToFeatureMapRef.current.set(brick.id, feature);
        // Register the Leaflet layer so we can apply CSS animations to it
        layerMapRef.current.set(brick.id, layer);
      }

      layer.on({
        mouseover: () => {
          if (brick) {
            // Pure client-side hover highlight (no API call)
            layer.setStyle({
              fillOpacity: 0.85,
              weight: 2.5,
            });
          } else {
            // Unassigned area hover: show "+" cursor
            layer.setStyle({
              fillColor: "#DBEAFE",
              fillOpacity: 0.35,
              weight: 1.5,
            });
            if (layer.getElement?.()) {
              layer.getElement().style.cursor = "crosshair";
            }
          }
        },
        mouseout: () => {
          layer.setStyle(style(feature));
          if (!brick && layer.getElement?.()) {
            layer.getElement().style.cursor = "";
          }
        },
        click: () => {
          if (brick) {
            setSelectedBrick(brick);
            if (onBrickClick) onBrickClick(brick);
          } else if (onCreateBrickFromMap && regencyName && province) {
            // Unassigned area click -> trigger create flow
            onCreateBrickFromMap(regencyName, province);
          }
        },
      });

      // Tooltip for bricks: pure client-side, uses pre-loaded performance data
      if (brick) {
        const metrics = performanceByBrickId.get(brick.id);
        const revenue = metrics?.total_revenue ?? 0;
        const achievement = metrics?.achievement_percentage ?? 0;
        const revenueGrowth = metrics?.revenue_growth_percentage ?? 0;

        layer.bindTooltip(
          `<div style="text-align:center;min-width:160px;padding:4px;">
            <strong style="font-size:12px;">${brick.name}</strong>
            <div style="font-size:10px;color:#6B7280;margin:2px 0 6px;">${brick.code}</div>
            <div style="display:flex;gap:12px;justify-content:center;">
              <div>
                <div style="font-size:10px;color:#6B7280;">Revenue</div>
                <div style="font-size:12px;font-weight:600;color:${revenueGrowth >= 0 ? '#16A34A' : '#DC2626'};">
                  ${formatCurrency(revenue)}
                </div>
              </div>
              <div>
                <div style="font-size:10px;color:#6B7280;">Achievement</div>
                <div style="font-size:12px;font-weight:600;color:${getAchievementColor(achievement)};">
                  ${achievement.toFixed(1)}%
                </div>
              </div>
            </div>
          </div>`,
          {
            sticky: true,
            direction: "top",
            offset: [0, -10],
            className: "brick-tooltip-clean",
          }
        );
      } else if (regencyName) {
        // Unassigned area tooltip
        layer.bindTooltip(
          `<div style="text-align:center;padding:4px;">
            <strong style="font-size:11px;">${regencyName}</strong>
            <div style="font-size:10px;color:#6B7280;">${province}</div>
            <div style="margin-top:4px;font-size:10px;color:#2563EB;display:flex;align-items:center;gap:4px;justify-content:center;">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
              Click to create brick
            </div>
          </div>`,
          {
            sticky: true,
            direction: "top",
            offset: [0, -10],
            className: "brick-tooltip-clean",
          }
        );
      }
    },
    [hasBrick, style, onBrickClick, onCreateBrickFromMap, extractRegencyAndProvince, performanceByBrickId, formatCurrency]
  );

  const currentMapStyle = mapStyles[mapStyle];

  if (loading) {
    return (
      <div className="h-full w-full flex items-center justify-center bg-muted/30">
        <div className="flex flex-col items-center gap-4">
          <div className="relative w-16 h-16">
            <div className="absolute inset-0 border-4 border-primary/20 rounded-full" />
            <div className="absolute inset-0 border-4 border-primary border-t-transparent rounded-full animate-spin" />
            <Grid3X3 className="absolute inset-0 m-auto w-6 h-6 text-primary/60" />
          </div>
          <p className="text-sm text-muted-foreground font-medium">Loading map...</p>
        </div>
      </div>
    );
  }

  if (error || !geoJsonData) {
    return (
      <div className="h-full w-full flex items-center justify-center bg-muted/30">
        <div className="text-center space-y-4">
          <Grid3X3 className="h-12 w-12 text-muted-foreground mx-auto" />
          <p className="text-sm font-medium text-muted-foreground">{error || "Failed to load map"}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full w-full relative">
      {/* Map Style Selector */}
      <div className="absolute top-4 right-4 z-30">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="secondary"
              size="icon"
              className="h-10 w-10 rounded-lg shadow-lg bg-card/95 backdrop-blur-sm border-0 hover:bg-card cursor-pointer"
            >
              <Layers className="h-5 w-5" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="min-w-[140px]">
            {(Object.keys(mapStyles) as MapStyle[]).map((ms) => (
              <DropdownMenuItem
                key={ms}
                onClick={() => setMapStyle(ms)}
                className={`cursor-pointer ${mapStyle === ms ? "bg-primary/10 text-primary font-medium" : ""}`}
              >
                {mapStyles[ms].name}
              </DropdownMenuItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {/* Heatmap Legend */}
      <div className="absolute bottom-6 right-6 z-30">
        <div className="bg-card/95 backdrop-blur-md rounded-xl shadow-lg border border-border/30 p-4">
          <h4 className="text-xs font-medium text-foreground mb-3">
            {colorBy === "revenue" ? "Revenue" : "Achievement"}
          </h4>
          <div className="flex items-center gap-1 mb-1">
            {[0, 0.2, 0.4, 0.6, 0.8, 1].map((ratio) => (
              <div
                key={`heatmap-legend-${ratio}`}
                className="w-5 h-3 rounded-sm"
                style={{ backgroundColor: getHeatmapColor(ratio, 0, 1) }}
              />
            ))}
          </div>
          <div className="flex justify-between">
            <span className="text-[10px] text-muted-foreground">Low</span>
            <span className="text-[10px] text-muted-foreground">High</span>
          </div>
          <div className="mt-3 space-y-1.5">
            <div className="flex items-center gap-2">
              <div className="w-4 h-4 rounded bg-gray-400" />
              <span className="text-[10px] text-muted-foreground">{t("inactive")}</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-4 h-4 rounded bg-amber-500" />
              <span className="text-[10px] text-muted-foreground">{t("selected")}</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-4 h-4 rounded-full bg-green-600" />
              <span className="text-[10px] text-muted-foreground">
                {t("accountPins", { count: accountsWithCoordinates.length })}
              </span>
            </div>
          </div>
        </div>
      </div>

      <MapContainer
        center={mapCenter}
        zoom={mapZoom}
        style={{ height: "100%", width: "100%", zIndex: 0 }}
        scrollWheelZoom={true}
        zoomControl={false}
        preferCanvas={true}
        wheelPxPerZoomLevel={120}
        zoomSnap={0.5}
        zoomDelta={0.5}
      >
        <SmartTileLayer
          key={`smart-tile-${mapStyle}`}
          source={currentMapStyle.source}
          fallbackSources={mapStyle === "dark" ? DARK_FALLBACK_CHAIN : LIGHT_FALLBACK_CHAIN}
          maxRetries={2}
          retryDelay={200}
          priorityMode="viewport"
        />
        <MapFocus
          key={`focus-${shouldAnimate ? "animate" : "static"}-${Array.isArray(mapCenter) ? `${mapCenter[0]}-${mapCenter[1]}` : "default"}-${mapZoom}`}
          center={mapCenter}
          zoom={mapZoom}
          shouldAnimate={shouldAnimate}
        />
        {geoJsonData && (
          <GeoJSON
            key={`geojson-${geoJsonKey}`}
            data={geoJsonData}
            style={style}
            onEachFeature={onEachFeature}
            eventHandlers={{
              // Increment version on every (re)mount so dependent effects re-run
              add: () => setGeoJsonVersion((v) => v + 1),
            }}
          />
        )}
        {accountsWithCoordinates.length > 0 && (
          <MarkerCluster
            accounts={accountsWithCoordinates}
            selectedBrickId={selectedBrick?.id ?? selectedBrickId}
            labels={{
              noCategory: t("noCategory"),
              city: t("city"),
              coordinates: t("coordinates"),
            }}
          />
        )}
      </MapContainer>

      {/* Tooltip and selected-brick highlight styles */}
      <style jsx global>{`
        .brick-tooltip-clean {
          background: var(--card) !important;
          border: 1px solid var(--border) !important;
          border-radius: 8px !important;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12) !important;
          padding: 6px 8px !important;
          color: var(--card-foreground) !important;
          font-size: 12px !important;
        }
        .brick-tooltip-clean::before {
          border-top-color: var(--card) !important;
        }
        @keyframes brick-pulse-stroke {
          0%, 100% { stroke-opacity: 1; stroke-width: 3px; }
          50%       { stroke-opacity: 0.4; stroke-width: 5px; }
        }
        .brick-selected-pulse {
          animation: brick-pulse-stroke 1.6s ease-in-out infinite !important;
          stroke: hsl(36, 100%, 50%) !important;
          stroke-width: 3px !important;
          filter: drop-shadow(0 0 6px hsla(36, 100%, 50%, 0.7)) !important;
        }
        .brick-account-popup {
          min-width: 220px;
          display: flex;
          flex-direction: column;
          gap: 6px;
          padding: 2px 0;
        }
        .brick-account-popup__title {
          font-size: 13px;
          font-weight: 600;
          color: var(--foreground);
        }
        .brick-account-popup__meta {
          font-size: 11px;
          color: var(--muted-foreground);
        }
        .brick-account-popup__row {
          display: flex;
          flex-direction: column;
          gap: 2px;
          font-size: 11px;
        }
        .brick-account-popup__row span {
          color: var(--muted-foreground);
        }
        .brick-account-popup__row strong {
          color: var(--foreground);
          font-weight: 600;
        }
      `}</style>
    </div>
  );
}

export const BrickMapFull = dynamic(
  () => Promise.resolve({ default: BrickMapFullComponent }),
  {
    ssr: false,
    loading: () => (
      <div className="h-full w-full flex items-center justify-center bg-muted/30">
        <Skeleton className="h-full w-full" />
      </div>
    ),
  }
);
