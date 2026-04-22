"use client";

import dynamic from "next/dynamic";
import { Skeleton } from "@/components/ui/skeleton";
import { useCallback, useRef, useState } from "react";
import type { GeoJSONPolygon } from "../types";
import { SmartTileLayer, TILE_SOURCES, LIGHT_FALLBACK_CHAIN } from "@/components/ui/smart-tile-layer";
import { useTheme } from "next-themes";

// Dynamic imports untuk SSR compatibility
const MapContainer = dynamic(
  () => import("react-leaflet").then((mod) => mod.MapContainer),
  { ssr: false }
);

// TileLayer removed - using SmartTileLayer instead

const FeatureGroup = dynamic(
  () => import("react-leaflet").then((mod) => mod.FeatureGroup),
  { ssr: false }
);

const GeoJSON = dynamic(
  () => import("react-leaflet").then((mod) => mod.GeoJSON),
  { ssr: false }
);

const EditControl = dynamic(
  () => import("./EditControl").then((mod) => ({ default: mod.EditControl })),
  { ssr: false }
);

interface AreaMapProps {
  center?: [number, number];
  zoom?: number;
  color?: string;
  initialPolygon?: GeoJSONPolygon;
  onPolygonChange: (polygon: GeoJSONPolygon | null) => void;
  className?: string;
}

function AreaMapComponent({
  center = [-6.2, 106.816666], // Jakarta default
  zoom = 12,
  color = "#3B82F6",
  initialPolygon,
  onPolygonChange,
  className = "h-[400px] w-full",
}: AreaMapProps) {
  const featureGroupRef = useRef<L.FeatureGroup>(null);
  const [currentPolygon, setCurrentPolygon] = useState<L.Layer | null>(null);
  const { theme } = useTheme();

  const onCreated = useCallback((e: any) => {
    const layer = e.layer;
    
    // Clear existing polygon (only allow one)
    if (currentPolygon && featureGroupRef.current) {
      featureGroupRef.current.removeLayer(currentPolygon);
    }
    
    // Add new polygon to feature group
    if (featureGroupRef.current) {
      featureGroupRef.current.addLayer(layer);
    }
    
    // Convert to GeoJSON
    const geojson = layer.toGeoJSON();
    const polygon: GeoJSONPolygon = {
      type: "Polygon",
      coordinates: geojson.geometry.coordinates,
    };
    
    
    // Update state
    setCurrentPolygon(layer);
    onPolygonChange(polygon);
  }, [currentPolygon, onPolygonChange]);

  const onEdited = useCallback((e: any) => {
    const layers = e.layers;
    
    layers.eachLayer((layer: any) => {
      const geojson = layer.toGeoJSON();
      const polygon: GeoJSONPolygon = {
        type: "Polygon",
        coordinates: geojson.geometry.coordinates,
      };
      
      onPolygonChange(polygon);
    });
  }, [onPolygonChange]);

  const onDeleted = useCallback(() => {
    setCurrentPolygon(null);
    onPolygonChange(null);
  }, [onPolygonChange]);

  // Draw options for EditControl
  const drawOptions = {
    polygon: {
      shapeOptions: {
        color: color,
        fillColor: color,
        fillOpacity: 0.3,
        weight: 3,
        opacity: 0.9,
      },
      allowIntersection: false,
      drawError: {
        color: "#e74c3c",
        message: "<strong>Error:</strong> Shape edges cannot cross!",
      },
    },
    rectangle: {
      shapeOptions: {
        color: color,
        fillColor: color,
        fillOpacity: 0.3,
        weight: 3,
        opacity: 0.9,
      },
    },
    polyline: false,
    circle: false,
    marker: false,
    circlemarker: false,
  };

  const editOptions = featureGroupRef.current ? {
    featureGroup: featureGroupRef.current,
    remove: true,
  } : undefined;

  return (
    <div className={className}>
      <MapContainer
        center={center}
        zoom={zoom}
        style={{ height: "100%", width: "100%" }}
        className="rounded-md border"
      >
        <SmartTileLayer
          source={theme === "dark" ? TILE_SOURCES.cartoDark : TILE_SOURCES.openStreetMap}
          fallbackSources={LIGHT_FALLBACK_CHAIN}
          maxRetries={2}
          retryDelay={200}
          priorityMode="viewport"
        />
        
        <FeatureGroup ref={featureGroupRef}>
          <EditControl
            position="topright"
            onCreated={onCreated}
            onEdited={onEdited}
            onDeleted={onDeleted}
            draw={drawOptions}
            edit={editOptions}
          />
        </FeatureGroup>
        
        {/* Render initial polygon if provided */}
        {initialPolygon && (
          <GeoJSON
            key={JSON.stringify(initialPolygon)} // Force re-render when polygon changes
            data={initialPolygon}
            style={{
              color: color,
              fillColor: color,
              fillOpacity: 0.3,
              weight: 3,
              opacity: 0.9,
            }}
          />
        )}
      </MapContainer>
    </div>
  );
}

// Export with loading fallback
const DynamicAreaMapComponent = dynamic(() => Promise.resolve(AreaMapComponent), {
  ssr: false,
  loading: () => (
    <Skeleton className="h-[400px] w-full flex items-center justify-center">
      <span className="text-muted-foreground">Loading map...</span>
    </Skeleton>
  ),
});

export function AreaMap(props: AreaMapProps) {
  return <DynamicAreaMapComponent {...props} />;
}