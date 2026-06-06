"use client";

import { useState, forwardRef, useImperativeHandle, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { 
  Plus, 
  Route as RouteIcon, 
  History, 
  Navigation, 
  List,
  ChevronRight,
  Loader2,
  Sparkles,
  MapPin
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { RouteOptimizationForm } from "./route-optimization-form";
import { RouteMap } from "./route-map";
import { RouteHistoryDialog } from "./route-history-dialog";
import { RouteDetailsPanel } from "./route-details-panel";
import { useOptimizeRoute, useRouteList } from "../hooks/useRouteOptimization";
import { toast } from "sonner";
import type { OptimizeRouteFormData } from "../schemas/route-optimization.schema";
import type { OptimizedRoute, RouteStep, Waypoint } from "../types";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import type { MapFocus } from "@/components/ui/map";

export interface RouteOptimizationManagementRef {
  openForm: () => void;
}

export const RouteOptimizationManagement = forwardRef<RouteOptimizationManagementRef>((props, ref) => {
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isHistoryOpen, setIsHistoryOpen] = useState(false);
  const [isDetailsPanelOpen, setIsDetailsPanelOpen] = useState(false);
  const [selectedRoute, setSelectedRoute] = useState<OptimizedRoute | null>(null);
  const [mapFocus, setMapFocus] = useState<MapFocus | null>(null);
  const user = useAuthStore((state) => state.user);

  const optimizeRoute = useOptimizeRoute();
  const { data: routesData, isLoading: routesLoading } = useRouteList({ user_id: user?.id });

  useImperativeHandle(ref, () => ({
    openForm: () => setIsFormOpen(true),
  }));

  const handleOptimize = async (data: OptimizeRouteFormData) => {
    try {
      const result = await optimizeRoute.mutateAsync(data);
      setSelectedRoute(result.data);
      setIsFormOpen(false);
      setIsDetailsPanelOpen(true);
      toast.success("Route optimized successfully!");
    } catch (error) {
      // Error handled by interceptor
    }
  };

  const handleSelectRoute = (route: OptimizedRoute) => {
    setSelectedRoute(route);
    setIsDetailsPanelOpen(true);
  };

  const handleWaypointClick = useCallback((waypoint: Waypoint, _index: number) => {
    setMapFocus({ lat: waypoint.lat, lng: waypoint.lng, zoom: 17 });
  }, []);

  const handleStepClick = useCallback((step: RouteStep, _index: number) => {
    setMapFocus({
      lat: step.end_location.lat,
      lng: step.end_location.lng,
      zoom: 17,
    });
  }, []);

  // Get current route to display (selected or latest)
  const currentRoute = selectedRoute || routesData?.data?.[0] || null;
  const hasRoutes = (routesData?.data?.length ?? 0) > 0;

  return (
    <div className="relative w-full h-full">
      {/* Full-screen Map */}
      <div className="absolute inset-0">
        {currentRoute ? (
          <RouteMap
            waypoints={currentRoute.waypoints}
            optimizedRoute={currentRoute}
            onMarkerClick={handleWaypointClick}
            focus={mapFocus}
            height="100%"
            showControls
            showRouteInfo
          />
        ) : (
          <RouteMap
            waypoints={[]}
            height="100%"
            showControls={false}
            showRouteInfo={false}
          />
        )}
      </div>

      {/* Floating Controls - Top Left */}
      <div className="absolute top-6 left-6 z-30 flex flex-col gap-3">
        {/* Logo/Title */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="bg-card/95 backdrop-blur-md rounded-2xl shadow-2xl border border-border/30 p-4 flex items-center gap-3"
        >
          <div className="w-10 h-10 rounded-xl bg-linear-to-br from-primary to-primary/80 flex items-center justify-center shadow-lg">
            <Navigation className="w-5 h-5 text-white" />
          </div>
          <div>
            <h1 className="font-medium text-foreground">Route Optimization</h1>
            <p className="text-xs text-muted-foreground">
              {currentRoute 
                ? `${currentRoute.waypoints.length - 1} destinations` 
                : "Plan your optimal route"}
            </p>
          </div>
        </motion.div>

        {/* Action Buttons */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="flex flex-col gap-2"
        >
          {/* Create New Route Button */}
          <Button
            onClick={() => setIsFormOpen(true)}
            className="bg-linear-to-r from-primary to-primary/90 hover:from-primary/90 hover:to-primary !text-white shadow-lg rounded-xl h-12 px-5 font-medium cursor-pointer"
          >
            {optimizeRoute.isPending ? (
              <>
                <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                Optimizing...
              </>
            ) : (
              <>
                <Sparkles className="w-5 h-5 mr-2" />
                New Route
              </>
            )}
          </Button>

          {/* History Button */}
          <Button
            variant="secondary"
            onClick={() => setIsHistoryOpen(true)}
            className="bg-card/95 backdrop-blur-md hover:bg-card shadow-lg rounded-xl h-12 px-5 font-medium border-0 cursor-pointer"
          >
            <History className="w-5 h-5 mr-2" />
            Route History
            {hasRoutes && (
              <span className="ml-2 px-2 py-0.5 bg-primary/10 text-primary rounded-full text-xs font-medium">
                {routesData?.data?.length}
              </span>
            )}
          </Button>

          {/* View Route Details Button */}
          {currentRoute && (
            <Button
              variant="secondary"
              onClick={() => setIsDetailsPanelOpen(true)}
              className="bg-card/95 backdrop-blur-md hover:bg-card shadow-lg rounded-xl h-12 px-5 font-medium border-0 cursor-pointer"
            >
              <List className="w-5 h-5 mr-2" />
              View Details
              <ChevronRight className="w-4 h-4 ml-auto" />
            </Button>
          )}
        </motion.div>
      </div>

      {/* Empty State Overlay - Hide when form is open */}
      <AnimatePresence>
        {!currentRoute && !routesLoading && !isFormOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="absolute inset-0 flex items-center justify-center z-20 pointer-events-none"
          >
            <div className="bg-card/90 backdrop-blur-xl rounded-3xl shadow-2xl border border-border/30 p-8 max-w-sm text-center pointer-events-auto">
              <div className="w-20 h-20 mx-auto mb-6 rounded-2xl bg-linear-to-br from-primary/20 to-primary/10 flex items-center justify-center">
                <RouteIcon className="w-10 h-10 text-primary" />
              </div>
              <h2 className="text-xl font-medium text-foreground mb-2">
                Plan Your Route
              </h2>
              <p className="text-sm text-muted-foreground mb-6">
                Create an optimized route to visit multiple destinations in the most efficient order
              </p>
              <Button
                onClick={() => setIsFormOpen(true)}
                className="w-full h-12 rounded-xl bg-linear-to-r from-primary to-primary/90 hover:from-primary/90 hover:to-primary !text-white font-medium shadow-lg cursor-pointer"
              >
                <Plus className="w-5 h-5 mr-2" />
                Create First Route
              </Button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Loading State */}
      <AnimatePresence>
        {routesLoading && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="absolute inset-0 flex items-center justify-center z-20 bg-background/50 backdrop-blur-sm"
          >
            <div className="flex flex-col items-center gap-4">
              <div className="relative w-16 h-16">
                <div className="absolute inset-0 border-4 border-primary/20 rounded-full" />
                <div className="absolute inset-0 border-4 border-primary border-t-transparent rounded-full animate-spin" />
                <Navigation className="absolute inset-0 m-auto w-6 h-6 text-primary/60" />
              </div>
              <p className="text-sm font-medium text-muted-foreground">Loading routes...</p>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Route Optimization Form Dialog */}
      <RouteOptimizationForm
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
        onSubmit={handleOptimize}
        isLoading={optimizeRoute.isPending}
      />

      {/* Route History Dialog */}
      <RouteHistoryDialog
        isOpen={isHistoryOpen}
        onClose={() => setIsHistoryOpen(false)}
        userId={user?.id}
        onSelectRoute={handleSelectRoute}
      />

      {/* Route Details Panel */}
      <RouteDetailsPanel
        route={currentRoute}
        isOpen={isDetailsPanelOpen}
        onClose={() => setIsDetailsPanelOpen(false)}
        onWaypointClick={handleWaypointClick}
        onStepClick={handleStepClick}
      />
    </div>
  );
});

RouteOptimizationManagement.displayName = "RouteOptimizationManagement";
