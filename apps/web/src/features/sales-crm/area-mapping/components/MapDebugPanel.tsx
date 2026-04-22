"use client";

import { useEffect, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface MapDebugPanelProps {
  mapRef: React.RefObject<L.Map | null>;
  drawnItemsRef: React.RefObject<L.FeatureGroup | null>;
  polygon?: number[][][];
}

export function MapDebugPanel({ mapRef, drawnItemsRef, polygon }: MapDebugPanelProps) {
  const [debugInfo, setDebugInfo] = useState({
    mapReady: false,
    featureGroupReady: false,
    layerCount: 0,
    polygonPoints: 0,
    leafletDrawAvailable: false,
  });

  const updateDebugInfo = () => {
    setDebugInfo({
      mapReady: !!mapRef.current,
      featureGroupReady: !!drawnItemsRef.current,
      layerCount: drawnItemsRef.current?.getLayers().length || 0,
      polygonPoints: polygon?.[0]?.length || 0,
      leafletDrawAvailable: typeof window !== "undefined" && !!(window as any).L?.Control?.Draw,
    });
  };

  useEffect(() => {
    updateDebugInfo();
    const interval = setInterval(updateDebugInfo, 1000);
    return () => clearInterval(interval);
  }, [mapRef, drawnItemsRef, polygon]);

  const forceVisibility = () => {
    if (mapRef.current && drawnItemsRef.current) {
      const layers = drawnItemsRef.current.getLayers();
      
      layers.forEach((layer, index) => {
        try {
          // Add layer directly to map if not already there
          if (!mapRef.current?.hasLayer(layer)) {
            mapRef.current?.addLayer(layer);
          }
          
          // Bring to front
          if ('bringToFront' in layer && typeof layer.bringToFront === 'function') {
            layer.bringToFront();
          }
          
          // Apply styling
          if ('setStyle' in layer && typeof layer.setStyle === 'function') {
            (layer as any).setStyle({
              color: '#3B82F6',
              fillColor: '#3B82F6',
              fillOpacity: 0.3,
              weight: 3,
              opacity: 0.9,
            });
          }
        } catch (err) {
        }
      });
      
      // Force map redraw
      mapRef.current.invalidateSize(false);
    }
  };

  const forceRedraw = () => {
    if (mapRef.current) {
      mapRef.current.invalidateSize(false);
    }
  };

  const showLayers = () => {
    if (drawnItemsRef.current && mapRef.current) {
      const layers = drawnItemsRef.current.getLayers();
      layers.forEach((layer, index) => {
        if ('getBounds' in layer && typeof layer.getBounds === 'function') {
          try {
          } catch (e) {
          }
        }
      });
      
      // Also check what's directly on the map
      const mapLayers = (mapRef.current as any)._layers || {};
      const mapLayerCount = Object.keys(mapLayers).length;
    }
  };

  if (process.env.NODE_ENV !== "development") {
    return null;
  }

  return (
    <Card className="mt-2">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm">Map Debug Panel</CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        <div className="grid grid-cols-2 gap-2 text-xs">
          <div className="flex items-center justify-between">
            <span>Map Ready:</span>
            <Badge variant={debugInfo.mapReady ? "default" : "destructive"}>
              {debugInfo.mapReady ? "✅" : "❌"}
            </Badge>
          </div>
          <div className="flex items-center justify-between">
            <span>FeatureGroup:</span>
            <Badge variant={debugInfo.featureGroupReady ? "default" : "destructive"}>
              {debugInfo.featureGroupReady ? "✅" : "❌"}
            </Badge>
          </div>
          <div className="flex items-center justify-between">
            <span>Layers:</span>
            <Badge variant="outline">{debugInfo.layerCount}</Badge>
          </div>
          <div className="flex items-center justify-between">
            <span>Polygon Points:</span>
            <Badge variant="outline">{debugInfo.polygonPoints}</Badge>
          </div>
          <div className="flex items-center justify-between col-span-2">
            <span>Leaflet Draw:</span>
            <Badge variant={debugInfo.leafletDrawAvailable ? "default" : "destructive"}>
              {debugInfo.leafletDrawAvailable ? "✅" : "❌"}
            </Badge>
          </div>
        </div>
        <div className="flex gap-2">
          <Button size="sm" variant="outline" onClick={forceRedraw}>
            Force Redraw
          </Button>
          <Button size="sm" variant="outline" onClick={forceVisibility}>
            Force Visible
          </Button>
          <Button size="sm" variant="outline" onClick={showLayers}>
            Log Layers
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}