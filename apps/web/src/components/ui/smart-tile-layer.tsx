"use client";

import { useEffect, useRef, useCallback, useState } from "react";
import dynamic from "next/dynamic";

// Dynamic import for react-leaflet components
const TileLayer = dynamic(
  () => import("react-leaflet").then((mod) => mod.TileLayer),
  { ssr: false }
);

// ============== TYPES ==============
interface TileSource {
  name: string;
  url: string;
  attribution: string;
  subdomains?: string;
  maxZoom?: number;
}

interface SmartTileLayerProps {
  /** Primary tile source to use */
  source: TileSource;
  /** Fallback tile sources if primary fails */
  fallbackSources?: TileSource[];
  /** Number of retry attempts per tile before switching to fallback */
  maxRetries?: number;
  /** Delay between retries in ms */
  retryDelay?: number;
  /** Priority mode for viewport-based loading */
  priorityMode?: "viewport" | "all";
  /** Callback when tile loading starts */
  onLoadingStart?: () => void;
  /** Callback when all tiles are loaded */
  onLoadingComplete?: () => void;
  /** Callback when error occurs */
  onError?: (error: Error, tileUrl: string) => void;
}

// ============== FALLBACK TILE SOURCES ==============
// These are reliable CDN-backed tile servers as fallbacks
export const TILE_SOURCES = {
  // CartoDB CDN - very reliable and fast
  cartoLight: {
    name: "CartoDB Light",
    url: "https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png",
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>',
    subdomains: "abcd",
    maxZoom: 20,
  },
  cartoDark: {
    name: "CartoDB Dark",
    url: "https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png",
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>',
    subdomains: "abcd",
    maxZoom: 20,
  },
  // OpenStreetMap - good but can be rate limited
  openStreetMap: {
    name: "OpenStreetMap",
    url: "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    subdomains: "abc",
    maxZoom: 19,
  },
  // Stamen (now Stadia) - alternative
  stadiaLight: {
    name: "Stadia Light",
    url: "https://tiles.stadiamaps.com/tiles/alidade_smooth/{z}/{x}/{y}{r}.png",
    attribution: '&copy; <a href="https://stadiamaps.com/">Stadia Maps</a>, &copy; <a href="https://openmaptiles.org/">OpenMapTiles</a> &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors',
    maxZoom: 20,
  },
  stadiaDark: {
    name: "Stadia Dark",
    url: "https://tiles.stadiamaps.com/tiles/alidade_smooth_dark/{z}/{x}/{y}{r}.png",
    attribution: '&copy; <a href="https://stadiamaps.com/">Stadia Maps</a>, &copy; <a href="https://openmaptiles.org/">OpenMapTiles</a> &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors',
    maxZoom: 20,
  },
  // Esri Satellite - for satellite imagery
  esriSatellite: {
    name: "Esri Satellite",
    url: "https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}",
    attribution: '&copy; Esri',
    maxZoom: 19,
  },
} as const;

// Default fallback chain for light theme
export const LIGHT_FALLBACK_CHAIN: TileSource[] = [
  TILE_SOURCES.cartoLight,
  TILE_SOURCES.stadiaLight,
  TILE_SOURCES.openStreetMap,
];

// Default fallback chain for dark theme
export const DARK_FALLBACK_CHAIN: TileSource[] = [
  TILE_SOURCES.cartoDark,
  TILE_SOURCES.stadiaDark,
  TILE_SOURCES.openStreetMap,
];

// ============== SMART TILE LAYER COMPONENT ==============
export function SmartTileLayer({
  source,
  fallbackSources = [],
  maxRetries = 2,
  retryDelay = 300,
  priorityMode = "viewport",
  onLoadingStart,
  onLoadingComplete,
  onError,
}: SmartTileLayerProps) {
  const [currentSourceIndex, setCurrentSourceIndex] = useState(0);
  const [tilesLoading, setTilesLoading] = useState(0);
  const retryCountRef = useRef<Map<string, number>>(new Map());
  const failedTilesRef = useRef<Set<string>>(new Set());
  const loadingStartedRef = useRef(false);
  
  // All available sources (primary + fallbacks)
  const allSources = [source, ...fallbackSources];
  const currentSource = allSources[currentSourceIndex] ?? source;
  
  // Reset retry counts when source changes
  useEffect(() => {
    retryCountRef.current.clear();
    failedTilesRef.current.clear();
    loadingStartedRef.current = false;
  }, [currentSourceIndex]);
  
  // Handle loading tracking - debounced to prevent rapid state changes
  useEffect(() => {
    if (tilesLoading > 0 && !loadingStartedRef.current) {
      loadingStartedRef.current = true;
      onLoadingStart?.();
    } else if (tilesLoading === 0 && loadingStartedRef.current) {
      // Debounce completion to avoid flickering
      const timeoutId = setTimeout(() => {
        if (tilesLoading === 0) {
          loadingStartedRef.current = false;
          onLoadingComplete?.();
        }
      }, 100);
      return () => clearTimeout(timeoutId);
    }
  }, [tilesLoading, onLoadingStart, onLoadingComplete]);
  
  // Tile event handlers
  const handleTileLoadStart = useCallback(() => {
    setTilesLoading((prev) => prev + 1);
  }, []);
  
  const handleTileLoad = useCallback(() => {
    setTilesLoading((prev) => Math.max(0, prev - 1));
  }, []);
  
  const handleTileError = useCallback((event: unknown) => {
    setTilesLoading((prev) => Math.max(0, prev - 1));
    
    // Extract tile URL from event
    const errorEvent = event as { tile?: { src?: string }; target?: { _url?: string } };
    const tileUrl = errorEvent?.tile?.src ?? errorEvent?.target?._url ?? "unknown";
    
    // Track retry count for this tile
    const currentRetries = retryCountRef.current.get(tileUrl) ?? 0;
    
    if (currentRetries < maxRetries) {
      // Retry with exponential backoff
      retryCountRef.current.set(tileUrl, currentRetries + 1);
      
      const delay = retryDelay * Math.pow(1.5, currentRetries);
      setTimeout(() => {
        // Try to reload the tile by creating a new image request
        if (typeof window !== "undefined" && tileUrl && tileUrl !== "unknown") {
          const img = new Image();
          img.crossOrigin = "anonymous";
          img.src = tileUrl;
        }
      }, delay);
    } else {
      // Mark tile as failed
      failedTilesRef.current.add(tileUrl);
      
      // If too many tiles failed, switch to fallback source
      const failedCount = failedTilesRef.current.size;
      const shouldSwitchSource = failedCount >= 5 && currentSourceIndex < allSources.length - 1;
      
      if (shouldSwitchSource) {
        // eslint-disable-next-line no-console
        console.warn(
          `[SmartTileLayer] Too many tiles failed (${failedCount}), switching to fallback source`
        );
        setCurrentSourceIndex((prev) => prev + 1);
      }
      
      // Call error callback
      onError?.(new Error(`Failed to load tile after ${maxRetries} retries`), tileUrl);
    }
  }, [maxRetries, retryDelay, currentSourceIndex, allSources.length, onError]);
  
  return (
    <TileLayer
      key={`tile-layer-${currentSourceIndex}-${currentSource.name}`}
      url={currentSource.url}
      attribution={currentSource.attribution}
      subdomains={currentSource.subdomains ?? "abcd"}
      maxZoom={currentSource.maxZoom ?? 19}
      
      // ============== PERFORMANCE OPTIMIZATIONS ==============
      // Cross-origin for better caching
      crossOrigin="anonymous"
      
      // CRITICAL: Only update tiles when user stops panning/zooming
      // This prevents loading tiles that will immediately be out of view
      updateWhenIdle={true}
      
      // Don't update during zoom animation
      updateWhenZooming={false}
      
      // Keep buffer - number of tiles to keep outside visible bounds
      // Lower = faster initial load, Higher = smoother panning
      // 2 is optimal for most use cases (preloads 1 tile ring around viewport)
      keepBuffer={priorityMode === "viewport" ? 2 : 4}
      
      // Native lazy loading for images
      // This uses browser's built-in lazy loading
      className="leaflet-tile-pane"
      
      // Event handlers for smart loading
      eventHandlers={{
        loading: handleTileLoadStart,
        load: handleTileLoad,
        tileloadstart: handleTileLoadStart,
        tileload: handleTileLoad,
        tileerror: handleTileError,
      }}
    />
  );
}

export type { TileSource, SmartTileLayerProps };
