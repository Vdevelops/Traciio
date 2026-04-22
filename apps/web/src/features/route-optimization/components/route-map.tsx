"use client";

import { useMemo, useCallback } from "react";
import { Map as BaseMap, type MapMarker, type MapRoute, type MapFocus } from "@/components/ui/map";
import type { Waypoint, OptimizedRoute } from "../types";

interface RouteMapProps {
  readonly waypoints: Waypoint[];
  readonly optimizedRoute?: OptimizedRoute;
  readonly onMarkerClick?: (waypoint: Waypoint, index: number) => void;
  readonly focus?: MapFocus | null;
  readonly height?: string;
  readonly showControls?: boolean;
  readonly showRouteInfo?: boolean;
}

/**
 * RouteMap - A specialized map component for route optimization
 * 
 * This is a wrapper around the reusable Map component from @/components/ui/map
 * that transforms route optimization data types to the generic Map types.
 */
export function RouteMap({
  waypoints,
  optimizedRoute,
  onMarkerClick,
  focus,
  height = "100%",
  showControls = true,
  showRouteInfo = true,
}: RouteMapProps) {
  // Backend returns waypoints already in optimized order:
  // [start(order=0), optimized_stop_1(order=1), optimized_stop_2(order=2), ...]
  // We use the `order` field for marker numbering — no re-indexing via optimized_order needed.

  // Transform waypoints to MapMarker format
  const markers: MapMarker[] = useMemo(() => {
    return waypoints.map((waypoint, index) => {
      const isStartLocation = waypoint.order === 0 || (waypoint.order == null && index === 0);
      const displayOrder = waypoint.order ?? index;

      return {
        id: `waypoint-${index}-${waypoint.account_id || 'no-account'}`,
        lat: waypoint.lat,
        lng: waypoint.lng,
        label: waypoint.address || `${waypoint.lat.toFixed(6)}, ${waypoint.lng.toFixed(6)}`,
        description: waypoint.account?.name || waypoint.account_name,
        order: isStartLocation ? 0 : displayOrder,
        isStart: isStartLocation,
        metadata: {
          account_id: waypoint.account_id,
          visit_report_id: waypoint.visit_report_id,
          originalIndex: index,
        },
      };
    });
  }, [waypoints]);

  // Transform optimized route to MapRoute format
  const route: MapRoute | undefined = useMemo(() => {
    if (!optimizedRoute) return undefined;
    
    return {
      polyline: optimizedRoute.route_polyline,
      markers,
      totalDistance: optimizedRoute.total_distance,
      totalDistanceFormatted: optimizedRoute.total_distance_formatted,
      totalDuration: optimizedRoute.total_duration,
      totalDurationFormatted: optimizedRoute.total_duration_formatted,
    };
  }, [optimizedRoute, markers]);

  // Handle marker click - transform back to original waypoint
  const handleMarkerClick = useCallback((marker: MapMarker) => {
    if (!onMarkerClick) return;
    
    const originalIndex = (marker.metadata?.originalIndex as number) ?? 0;
    const waypoint = waypoints[originalIndex];
    if (waypoint) {
      onMarkerClick(waypoint, originalIndex);
    }
  }, [onMarkerClick, waypoints]);

  return (
    <BaseMap
      markers={markers}
      route={route}
      focus={focus}
      height={height}
      showControls={showControls}
      showZoomControl={false}
      showRouteInfo={showRouteInfo}
      onMarkerClick={handleMarkerClick}
    />
  );
}
