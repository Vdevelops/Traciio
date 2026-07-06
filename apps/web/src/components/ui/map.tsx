"use client";

import { useEffect, useMemo, useState, useCallback, useRef } from "react";
import dynamic from "next/dynamic";
import polyline from "@mapbox/polyline";
import { MapPin, Route as RouteIcon, Clock, Gauge, Navigation, Layers, Play } from "lucide-react";
import { useTheme } from "next-themes";
import { Button } from "./button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "./dropdown-menu";
import { SmartTileLayer, TILE_SOURCES, LIGHT_FALLBACK_CHAIN, DARK_FALLBACK_CHAIN, type TileSource } from "./smart-tile-layer";

// Dynamically import Leaflet components (client-side only)
const MapContainer = dynamic(
  () => import("react-leaflet").then((mod) => mod.MapContainer),
  { ssr: false }
);
// TileLayer removed - using SmartTileLayer instead
const Marker = dynamic(
  () => import("react-leaflet").then((mod) => mod.Marker),
  { ssr: false }
);
const Polyline = dynamic(
  () => import("react-leaflet").then((mod) => mod.Polyline),
  { ssr: false }
);
const Popup = dynamic(
  () => import("react-leaflet").then((mod) => mod.Popup),
  { ssr: false }
);
const ZoomControl = dynamic(
  () => import("react-leaflet").then((mod) => mod.ZoomControl),
  { ssr: false }
);

// ============== TYPES ==============
export interface MapMarker {
  readonly id: string;
  readonly lat: number;
  readonly lng: number;
  readonly label?: string;
  readonly description?: string;
  readonly order?: number;
  readonly isStart?: boolean;
  readonly metadata?: Record<string, unknown>;
}

export interface MapRoute {
  readonly polyline?: string; // Encoded polyline string
  readonly markers: MapMarker[];
  readonly totalDistance?: number;
  readonly totalDistanceFormatted?: string;
  readonly totalDuration?: number;
  readonly totalDurationFormatted?: string;
}

export interface MapFocus {
  readonly lat: number;
  readonly lng: number;
  readonly zoom?: number;
}

export interface MapProps {
  readonly markers?: MapMarker[];
  readonly route?: MapRoute;
  readonly center?: [number, number];
  readonly zoom?: number;
  readonly focus?: MapFocus | null;
  readonly height?: string;
  readonly width?: string;
  readonly showControls?: boolean;
  readonly showZoomControl?: boolean;
  readonly showRouteInfo?: boolean;
  readonly onMarkerClick?: (marker: MapMarker) => void;
  readonly className?: string;
}

// ============== MAP STYLES ==============
type MapStyle = "light" | "dark" | "satellite" | "streets";

// Map styles using shared TILE_SOURCES for consistency
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

// ============== DYNAMIC COMPONENTS ==============
// Dynamic component for map bounds fitter
const MapBoundsFitterDynamic = dynamic(
  () => Promise.all([
    import("react-leaflet"),
    import("leaflet"),
  ]).then(([reactLeaflet, leaflet]) => {
    const { useMap } = reactLeaflet;
    const L = leaflet.default;
    
    return function MapBoundsFitterInner({
      markers,
      routePositions,
    }: {
      markers: MapMarker[];
      routePositions: [number, number][];
    }) {
      const map = useMap();
      
      useEffect(() => {
        const boundsPoints = routePositions.length > 1
          ? routePositions
          : markers.map((m) => [m.lat, m.lng] as [number, number]);

        if (boundsPoints.length > 0 && L && map) {
          const bounds = L.latLngBounds(boundsPoints);
          map.fitBounds(bounds, { padding: [80, 80] });
        }
      }, [markers, routePositions, map]);
      
      return null;
    };
  }),
  { ssr: false }
);

// ============== MAIN COMPONENT ==============
function MapComponent({ 
  markers = [],
  route,
  center,
  zoom = 13,
  focus = null,
  height = "100%",
  width = "100%",
  showControls = true,
  showZoomControl = true,
  showRouteInfo = true,
  onMarkerClick,
  className,
}: MapProps) {
  const isClient = typeof window !== "undefined";
  const [leafletLoaded, setLeafletLoaded] = useState(false);
  const [L, setL] = useState<any>(null);
  const mapRef = useRef<any>(null);
  const { theme, resolvedTheme } = useTheme();
  const [mapStyle, setMapStyle] = useState<MapStyle>("light");

  // Use route markers if route is provided, otherwise use markers prop
  const displayMarkers = route?.markers ?? markers;

  // Auto-select map style based on theme
  useEffect(() => {
    const currentTheme = resolvedTheme || theme;
    if (currentTheme === "dark") {
      setMapStyle("dark");
    } else {
      setMapStyle("light");
    }
  }, [theme, resolvedTheme]);

  // Load Leaflet on client side
  useEffect(() => {
    if (typeof window !== "undefined") {
      // Load CSS first
      const link = document.createElement("link");
      link.rel = "stylesheet";
      link.href = "https://unpkg.com/leaflet@1.9.4/dist/leaflet.css";
      link.integrity = "sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY=";
      link.crossOrigin = "";
      document.head.appendChild(link);

      // Then load Leaflet library
      import("leaflet").then((leaflet) => {
        const Leaflet = leaflet.default;
        // Fix for default marker icons in Next.js
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        delete (Leaflet.Icon.Default.prototype as any)._getIconUrl;
        Leaflet.Icon.Default.mergeOptions({
          iconRetinaUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png",
          iconUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png",
          shadowUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
        });
        setL(Leaflet);
        setLeafletLoaded(true);
      });
    }
  }, []);

  // Calculate center from markers if not provided
  const mapCenter: [number, number] = useMemo(() => {
    if (center) return center;
    if (displayMarkers.length > 0) {
      return [displayMarkers[0].lat, displayMarkers[0].lng];
    }
    return [-6.2088, 106.8456]; // Default to Jakarta
  }, [center, displayMarkers]);

  // Create polyline from route polyline or markers
  const polylinePositions: [number, number][] = useMemo(() => {
    if (!route || displayMarkers.length === 0) return [];
    
    // If encoded polyline is available, decode it
    if (route.polyline) {
      try {
        const decoded = polyline.decode(route.polyline);
        return decoded.map((coord: [number, number]) => [coord[0], coord[1]]);
      } catch (error) {
        // Fallback to straight line
      }
    }
    
    // Fallback: Create straight line from markers
    return displayMarkers.map((m) => [m.lat, m.lng] as [number, number]);
  }, [route, displayMarkers]);

  // Route info for overlay
  const routeInfo = useMemo(() => {
    if (!route) return null;
    const stopsCount = displayMarkers.filter(m => !m.isStart).length;
    return {
      distance: route.totalDistanceFormatted || "N/A",
      duration: route.totalDurationFormatted || "N/A",
      waypointsCount: stopsCount,
    };
  }, [route, displayMarkers]);

  const getMapInstance = useCallback(() => {
    const current = mapRef.current;
    return current?.leafletElement ?? current;
  }, []);

  useEffect(() => {
    const map = getMapInstance();
    if (!map || !focus) return;

    const nextZoom = focus.zoom ?? Math.max(zoom, 16);
    try {
      map.flyTo([focus.lat, focus.lng], nextZoom, { animate: true, duration: 0.6 });
    } catch {
      map.setView([focus.lat, focus.lng], nextZoom);
    }
  }, [focus, zoom, getMapInstance]);

  // Loading state
  if (!isClient || !leafletLoaded || !L) {
    return (
      <div 
        style={{ height, width }} 
        className={`overflow-hidden bg-linear-to-br from-blue-50/50 to-indigo-50/50 dark:from-gray-900/50 dark:to-gray-800/50 ${className ?? ""}`}
      >
        <div className="h-full w-full flex items-center justify-center">
          <div className="text-center space-y-3">
            <div className="relative w-16 h-16 mx-auto">
              <div className="absolute inset-0 border-4 border-primary/20 rounded-full" />
              <div className="absolute inset-0 border-4 border-primary border-t-transparent rounded-full animate-spin" />
              <Navigation className="absolute inset-0 m-auto w-6 h-6 text-primary/60" />
            </div>
            <p className="text-sm text-muted-foreground font-medium">Loading map...</p>
          </div>
        </div>
      </div>
    );
  }

  // Create numbered icon using primary color (orange theme)
  const createNumberedIcon = (number: number) => {
    if (!L) return null;
    return L.divIcon({
      className: 'custom-div-icon',
      html: `
        <div style="
          background: linear-gradient(135deg, oklch(0.73 0.19 55) 0%, oklch(0.68 0.16 55) 100%);
          color: white;
          width: 36px;
          height: 36px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-weight: 700;
          font-size: 14px;
          border: 3px solid white;
          box-shadow: 0 4px 12px rgba(243, 146, 0, 0.4), 0 2px 4px rgba(0,0,0,0.2);
          transition: transform 0.2s ease;
        ">${number}</div>
      `,
      iconSize: [36, 36],
      iconAnchor: [18, 18],
      popupAnchor: [0, -18]
    });
  };

  // Create start icon with pulse animation - using SVG Play icon instead of emoji
  const createStartIcon = () => {
    if (!L) return null;
    // Success color from CSS variables: oklch(0.65 0.15 150) = approximately #22c55e
    return L.divIcon({
      className: 'custom-div-icon',
      html: `
        <div style="position: relative;">
          <div style="
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: 50px;
            height: 50px;
            background: oklch(0.65 0.15 150);
            border-radius: 50%;
            opacity: 0.3;
            animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
          "></div>
          <div style="
            background: oklch(0.65 0.15 150);
            color: white;
            width: 42px;
            height: 42px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            border: 4px solid white;
            box-shadow: 0 4px 16px rgba(34, 197, 94, 0.4), 0 2px 4px rgba(0,0,0,0.2);
            position: relative;
            z-index: 1;
          ">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="white" stroke="none">
              <polygon points="5 3 19 12 5 21 5 3"></polygon>
            </svg>
          </div>
        </div>
        <style>
          @keyframes pulse {
            0%, 100% { transform: translate(-50%, -50%) scale(1); opacity: 0.3; }
            50% { transform: translate(-50%, -50%) scale(1.5); opacity: 0; }
          }
        </style>
      `,
      iconSize: [42, 42],
      iconAnchor: [21, 21],
      popupAnchor: [0, -21]
    });
  };

  const currentStyle = mapStyles[mapStyle];

  return (
    <div style={{ height, width }} className={`overflow-hidden relative ${className ?? ""}`}>
      {/* Map Style Selector - Top Right */}
      {showControls && (
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
              {(Object.keys(mapStyles) as MapStyle[]).map((style) => (
                <DropdownMenuItem 
                  key={style} 
                  onClick={() => setMapStyle(style)}
                  className={`cursor-pointer ${mapStyle === style ? "bg-primary/10 text-primary font-medium" : ""}`}
                >
                  {mapStyles[style].name}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      )}

      {/* Route Info Overlay - Bottom Left */}
      {showRouteInfo && routeInfo && (
        <div className="absolute bottom-6 left-6 z-30">
          <div className="bg-card/95 backdrop-blur-md rounded-2xl shadow-2xl border border-border/30 p-5 min-w-[220px]">
            <div className="flex items-center gap-2 mb-4">
              <div className="w-8 h-8 rounded-lg bg-linear-to-br from-primary to-primary/80 flex items-center justify-center">
                <RouteIcon className="w-4 h-4 text-white" />
              </div>
              <span className="font-bold text-sm text-foreground">Route Summary</span>
            </div>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2 text-muted-foreground">
                  <Gauge className="w-4 h-4" />
                  <span className="text-xs font-medium">Distance</span>
                </div>
                <span className="text-sm font-bold text-foreground">{routeInfo.distance}</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2 text-muted-foreground">
                  <Clock className="w-4 h-4" />
                  <span className="text-xs font-medium">Duration</span>
                </div>
                <span className="text-sm font-bold text-foreground">{routeInfo.duration}</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2 text-muted-foreground">
                  <MapPin className="w-4 h-4" />
                  <span className="text-xs font-medium">Stops</span>
                </div>
                <span className="text-sm font-bold text-foreground">{routeInfo.waypointsCount}</span>
              </div>
            </div>
          </div>
        </div>
      )}

      <MapContainer 
        center={mapCenter} 
        zoom={zoom} 
        style={{ height: "100%", width: "100%" }} 
        scrollWheelZoom={true} 
        className="z-0"
        zoomControl={false}
        ref={mapRef as any}
        // Performance optimizations
        preferCanvas={true}
        wheelPxPerZoomLevel={120}
        zoomSnap={0.5}
        zoomDelta={0.5}
      >
        <MapBoundsFitterDynamic
          markers={displayMarkers}
          routePositions={polylinePositions}
        />
        {/* Zoom Control - Bottom Right */}
        {showZoomControl && <ZoomControl position="bottomright" />}
        {/* SmartTileLayer with auto-retry and fallback */}
        <SmartTileLayer
          key={`smart-tile-${mapStyle}`}
          source={currentStyle.source}
          fallbackSources={mapStyle === "dark" ? DARK_FALLBACK_CHAIN : LIGHT_FALLBACK_CHAIN}
          maxRetries={2}
          retryDelay={200}
          priorityMode="viewport"
        />
        
        {/* Route Polyline */}
        {polylinePositions.length > 1 && (
          <>
            {/* Shadow/glow effect */}
            <Polyline 
              positions={polylinePositions} 
              color="#F39200" 
              weight={8} 
              opacity={0.2}
              dashArray="10, 5"
            />
            {/* Main route line */}
            <Polyline 
              positions={polylinePositions} 
              color="#F39200" 
              weight={5} 
              opacity={0.9}
              lineCap="round"
              lineJoin="round"
            />
          </>
        )}
        
        {/* Markers */}
        {displayMarkers.map((marker, index) => {
          const isStart = marker.isStart ?? (marker.order === 0 || index === 0);
          const displayOrder = marker.order ?? index;

          return (
            <Marker 
              key={marker.id || `marker-${index}`} 
              position={[marker.lat, marker.lng]} 
              icon={isStart ? createStartIcon() : createNumberedIcon(displayOrder)}
              eventHandlers={{ 
                click: () => onMarkerClick?.(marker),
                mouseover: (e: any) => {
                  if (e?.target) {
                    e.target.openPopup();
                  }
                },
              }}
            >
              <Popup 
                className="custom-popup"
                closeButton={true}
                autoPan={true}
                maxWidth={280}
              >
                <div className="text-sm space-y-2 p-2 bg-card">
                  <div className="flex items-center gap-2 font-bold pb-2 border-b-2 border-border">
                    {isStart ? (
                      <>
                        <div className="w-8 h-8 rounded-full bg-success flex items-center justify-center text-success-foreground shadow-md">
                          <Play className="w-4 h-4 fill-current" />
                        </div>
                        <div>
                          <div className="text-success font-bold">Start Location</div>
                          <div className="text-xs text-muted-foreground font-semibold">Your current location</div>
                        </div>
                      </>
                    ) : (
                      <>
                        <div className="w-8 h-8 rounded-full bg-linear-to-br from-primary to-primary/80 flex items-center justify-center text-white font-bold shadow-md">
                          {displayOrder}
                        </div>
                        <div>
                          <div className="text-foreground font-bold">
                            {marker.label || `Destination ${displayOrder}`}
                          </div>
                          {marker.description && (
                            <div className="text-xs text-primary font-semibold">{marker.description}</div>
                          )}
                        </div>
                      </>
                    )}
                  </div>
                  
                  <div className="space-y-1.5">
                    <div className="flex items-start gap-2 text-xs">
                      <MapPin className="w-3.5 h-3.5 text-destructive mt-0.5 shrink-0" />
                      <span className="text-foreground font-semibold leading-relaxed">
                        {marker.label || `${marker.lat.toFixed(6)}, ${marker.lng.toFixed(6)}`}
                      </span>
                    </div>
                  </div>
                </div>
              </Popup>
            </Marker>
          );
        })}
      </MapContainer>

      {/* Custom CSS for popup - Using CSS variables for semantic colors */}
      <style jsx global>{`
        .custom-popup .leaflet-popup-content-wrapper {
          border-radius: 12px;
          box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
          padding: 0;
          background: var(--card) !important;
          border: 2px solid var(--border);
        }
        .custom-popup .leaflet-popup-content {
          margin: 0;
          padding: 0;
          color: var(--card-foreground);
        }
        .custom-popup .leaflet-popup-tip {
          background: var(--card);
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
          border: 1px solid var(--border);
        }
        .custom-popup .leaflet-popup-close-button {
          color: var(--muted-foreground) !important;
          font-size: 20px !important;
          font-weight: bold !important;
        }
        .custom-div-icon {
          background: transparent !important;
          border: none !important;
        }
        .custom-div-icon:hover {
          transform: scale(1.1);
          transition: transform 0.2s ease;
        }
        /* Zoom control positioning */
        .leaflet-control-zoom {
          position: absolute !important;
          bottom: 100px !important;
          right: 16px !important;
          top: auto !important;
          left: auto !important;
        }
        .leaflet-top.leaflet-left {
          top: auto !important;
          bottom: 100px !important;
          left: auto !important;
          right: 16px !important;
        }
      `}</style>
    </div>
  );
}

// ============== EXPORTED COMPONENT ==============
export function Map(props: MapProps) {
  const isClient = typeof window !== "undefined";
  if (!isClient) {
    return (
      <div 
        style={{ height: props.height || "100%", width: props.width || "100%" }} 
        className={`bg-muted animate-pulse ${props.className ?? ""}`} 
      />
    );
  }
  return <MapComponent {...props} />;
}
