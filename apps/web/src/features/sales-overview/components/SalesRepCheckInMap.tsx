"use client";

import { useTranslations } from "next-intl";
import { useState, useEffect, useMemo } from "react";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { MapPin, ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import dynamic from "next/dynamic";
import type { SalesRepCheckInLocation } from "../types";
import { SmartTileLayer, TILE_SOURCES, LIGHT_FALLBACK_CHAIN } from "@/components/ui/smart-tile-layer";

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
const Popup = dynamic(
  () => import("react-leaflet").then((mod) => mod.Popup),
  { ssr: false }
);
const Polyline = dynamic(
  () => import("react-leaflet").then((mod) => mod.Polyline),
  { ssr: false }
);
// useMap hook must be imported directly (not dynamic)
import { useMap } from "react-leaflet";

// Import Leaflet CSS
import "leaflet/dist/leaflet.css";
import L from "leaflet";

// Fix for default marker icons in Next.js
if (typeof window !== "undefined") {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  delete (L.Icon.Default.prototype as any)._getIconUrl;
  L.Icon.Default.mergeOptions({
    iconRetinaUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png",
    iconUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png",
    shadowUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png",
  });
}

interface SalesRepCheckInMapProps {
  readonly locations: readonly SalesRepCheckInLocation[];
  readonly isLoading?: boolean;
  readonly totalVisits?: number;
  readonly page?: number;
  readonly perPage?: number;
  readonly onPageChange?: (page: number) => void;
  readonly onPerPageChange?: (perPage: number) => void;
}

// Component to handle map focus when location is selected
// Must be inside MapContainer to use useMap hook
function MapFocusHandler({ location }: { location: SalesRepCheckInLocation | null }) {
  const map = useMap();
  
  useEffect(() => {
    if (location?.location?.latitude && location.location.longitude) {
      map.setView(
        [location.location.latitude, location.location.longitude],
        15, // Zoom level when focusing on a location
        { animate: true, duration: 0.5 }
      );
    }
  }, [location, map]);
  
  return null;
}

// Create custom numbered marker icon
function createNumberedMarkerIcon(number: number) {
  return L.divIcon({
    className: "custom-numbered-marker",
    html: `
      <div style="
        background-color: #3b82f6;
        color: white;
        width: 32px;
        height: 32px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: bold;
        font-size: 14px;
        border: 3px solid white;
        box-shadow: 0 2px 4px rgba(0,0,0,0.3);
      ">
        ${number}
      </div>
    `,
    iconSize: [32, 32],
    iconAnchor: [16, 16],
  });
}

export function SalesRepCheckInMap({
  locations,
  isLoading,
  totalVisits = 0,
  page = 1,
  perPage = 50,
  onPageChange,
  onPerPageChange,
}: SalesRepCheckInMapProps) {
  const t = useTranslations("salesOverview");
  const [selectedLocation, setSelectedLocation] = useState<SalesRepCheckInLocation | null>(null);

  // Filter valid locations and sort by visit_number in descending order to show latest first (3-2-1)
  const validLocations = useMemo(() => {
    return locations
      .filter((loc) => loc.location?.latitude && loc.location?.longitude)
      .sort((a, b) => b.visit_number - a.visit_number); // Sort descending to show latest first
  }, [locations]);

  // Calculate center point from all locations
  const centerLat = useMemo(() => {
    if (validLocations.length === 0) return -6.2088; // Default to Jakarta
    return validLocations.reduce((sum, loc) => sum + (loc.location?.latitude ?? 0), 0) / validLocations.length;
  }, [validLocations]);

  const centerLng = useMemo(() => {
    if (validLocations.length === 0) return 106.8456; // Default to Jakarta
    return validLocations.reduce((sum, loc) => sum + (loc.location?.longitude ?? 0), 0) / validLocations.length;
  }, [validLocations]);

  // Calculate bounds for map to fit all markers
  const bounds = useMemo(() => {
    if (validLocations.length === 0) return null;
    const lats = validLocations.map((loc) => loc.location?.latitude ?? 0);
    const lngs = validLocations.map((loc) => loc.location?.longitude ?? 0);
    return L.latLngBounds(
      [Math.min(...lats), Math.min(...lngs)],
      [Math.max(...lats), Math.max(...lngs)]
    );
  }, [validLocations]);

  // Create polyline coordinates (ordered by visit_number)
  const polylineCoordinates = useMemo(() => {
    return validLocations.map((loc) => [
      loc.location?.latitude ?? 0,
      loc.location?.longitude ?? 0,
    ] as [number, number]);
  }, [validLocations]);

  if (isLoading) {
    return <Skeleton className="h-[500px] w-full" />;
  }

  if (locations.length === 0) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        {t("no_check_in_locations")}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Map Container */}
      <div className="relative h-[500px] w-full rounded-lg border overflow-hidden bg-muted">
        {globalThis.window !== undefined && (
          <MapContainer
            center={[centerLat, centerLng]}
            zoom={bounds ? undefined : 12}
            bounds={bounds ?? undefined}
            boundsOptions={{ padding: [50, 50] }}
            className="h-full w-full z-0"
            scrollWheelZoom={true}
          >
            <SmartTileLayer
              source={TILE_SOURCES.openStreetMap}
              fallbackSources={LIGHT_FALLBACK_CHAIN}
              maxRetries={2}
              retryDelay={200}
              priorityMode="viewport"
            />
            
            {/* Polyline connecting all points in order */}
            {polylineCoordinates.length > 1 && (
              <Polyline
                positions={polylineCoordinates}
                pathOptions={{
                  color: "#3b82f6",
                  weight: 3,
                  opacity: 0.6,
                }}
              />
            )}

            {/* Markers with numbers */}
            {validLocations.map((location) => (
              <Marker
                key={location.visit_report_id}
                position={[
                  location.location?.latitude ?? 0,
                  location.location?.longitude ?? 0,
                ]}
                icon={createNumberedMarkerIcon(location.visit_number)}
                eventHandlers={{
                  click: () => {
                    setSelectedLocation(location);
                  },
                }}
              >
                <Popup>
                  <div className="text-sm">
                    <div className="font-medium">Visit #{location.visit_number}</div>
                    {location.location?.address && (
                      <div className="mt-1">{location.location.address}</div>
                    )}
                    {location.account?.name && (
                      <div className="mt-1 text-muted-foreground">{location.account.name}</div>
                    )}
                    {location.purpose && (
                      <div className="mt-1 text-xs text-muted-foreground">{location.purpose}</div>
                    )}
                  </div>
                </Popup>
              </Marker>
            ))}

            {/* Map focus handler for auto-focus on card click */}
            <MapFocusHandler location={selectedLocation} />
          </MapContainer>
        )}
        
        {/* Overlay info */}
        {validLocations.length > 0 && (
          <div className="absolute top-2 left-2 bg-background/90 backdrop-blur-sm px-3 py-2 rounded-md border shadow-sm z-[1000]">
            <p className="text-xs font-medium">
              {validLocations.length} {validLocations.length === 1 ? "location" : "locations"}
            </p>
          </div>
        )}
      </div>

      {/* Locations List - Horizontal Slider */}
      <div className="space-y-4">
        {/* Header with pagination info */}
        {totalVisits > 0 && (
          <div className="flex items-center justify-between">
            <p className="text-sm text-muted-foreground">
              {t("showing_locations", {
                start: (page - 1) * perPage + 1,
                end: Math.min(page * perPage, totalVisits),
                total: totalVisits,
              })}
            </p>
            {onPageChange && (
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => onPageChange(Math.max(1, page - 1))}
                  disabled={page <= 1 || isLoading}
                  className="cursor-pointer"
                >
                  <ChevronLeft className="h-4 w-4" />
                </Button>
                <span className="text-sm text-muted-foreground">
                  {t("page_of", { page, total: Math.ceil(totalVisits / perPage) })}
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => onPageChange(page + 1)}
                  disabled={page * perPage >= totalVisits || isLoading}
                  className="cursor-pointer"
                >
                  <ChevronRight className="h-4 w-4" />
                </Button>
              </div>
            )}
          </div>
        )}

        {/* Horizontal Scrollable Cards */}
        <ScrollArea className="w-full whitespace-nowrap">
          <div className="flex gap-4 pb-4">
            {validLocations.map((location) => (
              <Card
                key={location.visit_report_id}
                className={`min-w-[320px] max-w-[320px] p-4 cursor-pointer hover:bg-accent/50 dark:hover:bg-accent/30 transition-colors flex-shrink-0 ${
                  selectedLocation?.visit_report_id === location.visit_report_id
                    ? "bg-accent border-primary"
                    : ""
                }`}
                onClick={() => setSelectedLocation(location)}
              >
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-2">
                    <Badge variant="default" className="font-medium">
                      #{location.visit_number}
                    </Badge>
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                  </div>
                  <span className="text-xs text-muted-foreground">
                    {new Date(location.check_in_time).toLocaleDateString()}
                  </span>
                </div>
                {location.account?.name && (
                  <div className="mt-2 font-medium truncate">{location.account.name}</div>
                )}
                {location.location?.address && (
                  <div className="mt-1 text-sm text-muted-foreground line-clamp-2">
                    {location.location.address}
                  </div>
                )}
                {location.purpose && (
                  <div className="mt-2 text-xs text-muted-foreground truncate">{location.purpose}</div>
                )}
                {location.location?.latitude !== undefined && location.location?.longitude !== undefined && (
                  <div className="mt-2 text-xs text-muted-foreground font-mono">
                    {location.location.latitude.toFixed(6)}, {location.location.longitude.toFixed(6)}
                  </div>
                )}
              </Card>
            ))}
          </div>
          <ScrollBar orientation="horizontal" />
        </ScrollArea>
      </div>
    </div>
  );
}
