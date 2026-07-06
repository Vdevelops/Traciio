"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { WaypointList } from "./waypoint-list";
import { WaypointSelectorDialog } from "./waypoint-selector-dialog";
import { optimizeRouteSchema, type OptimizeRouteFormData } from "../schemas/route-optimization.schema";
import type { Waypoint, Location } from "../types";
import { Calendar, MapPin, Navigation, Loader2, AlertCircle } from "lucide-react";
import { useCurrentLocation } from "../hooks/useRouteOptimization";
import { toast } from "sonner";

interface RouteOptimizationFormProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onSubmit: (data: OptimizeRouteFormData) => void;
  readonly isLoading?: boolean;
  readonly initialWaypoints?: Waypoint[];
}

export function RouteOptimizationForm({
  open,
  onOpenChange,
  onSubmit,
  isLoading = false,
  initialWaypoints = [],
}: RouteOptimizationFormProps) {
  const [waypoints, setWaypoints] = useState<Waypoint[]>(initialWaypoints);
  const [startLocation, setStartLocation] = useState<Location | null>(null);
  const [isWaypointSelectorOpen, setIsWaypointSelectorOpen] = useState(false);
  const getCurrentLocation = useCurrentLocation();

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<OptimizeRouteFormData>({
    resolver: zodResolver(optimizeRouteSchema),
    defaultValues: {
      route_name: "",
      start_location: undefined,
      waypoints: initialWaypoints,
    },
  });

  const handleGetCurrentLocation = async () => {
    try {
      const location = await getCurrentLocation.mutateAsync();
      setStartLocation(location);
      setValue("start_location", location);
      toast.success("Current location captured successfully!");
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to get current location");
    }
  };

  const handleFormSubmit = (data: OptimizeRouteFormData) => {
    if (!startLocation) {
      toast.error("Please capture your current location first");
      return;
    }
    onSubmit({
      ...data,
      start_location: startLocation,
      waypoints,
    });
  };

  const handleAddWaypoint = () => {
    // Open waypoint selector dialog
    setIsWaypointSelectorOpen(true);
  };

  const handleSelectWaypoints = (selectedWaypoints: Waypoint[]) => {
    const newWaypoints = [...waypoints, ...selectedWaypoints];
    setWaypoints(newWaypoints);
    setValue("waypoints", newWaypoints);
  };

  const handleRemoveWaypoint = (index: number) => {
    const newWaypoints = waypoints.filter((_, i) => i !== index);
    setWaypoints(newWaypoints);
    setValue("waypoints", newWaypoints);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Optimize Route</DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="route_name">Route Name (Optional)</Label>
            <Input
              id="route_name"
              {...register("route_name")}
              placeholder="e.g., Jakarta Sales Route"
            />
            {errors.route_name && (
              <p className="text-sm text-destructive">
                {errors.route_name.message}
              </p>
            )}
          </div>

          {/* Start Location Section */}
          <div className="space-y-2">
            <Label>Start Location (Your Current Location)</Label>
            <div className="flex gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={handleGetCurrentLocation}
                disabled={getCurrentLocation.isPending}
                className="flex-1"
              >
                {getCurrentLocation.isPending ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Getting Location...
                  </>
                ) : (
                  <>
                    <Navigation className="mr-2 h-4 w-4" />
                    Use Current Location
                  </>
                )}
              </Button>
            </div>

            {startLocation && (
              <Alert className="bg-primary/5 border-primary/20">
                <MapPin className="h-4 w-4 text-primary" />
                <AlertDescription>
                  <div className="space-y-1">
                    <p className="font-medium text-sm">Starting from:</p>
                    <p className="text-xs text-muted-foreground">
                      {startLocation.address}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {startLocation.lat.toFixed(6)}, {startLocation.lng.toFixed(6)}
                    </p>
                    {startLocation.accuracy != null && (
                      <p className="text-xs text-muted-foreground">
                        Accuracy: ±{Math.round(startLocation.accuracy)} m
                      </p>
                    )}
                  </div>
                </AlertDescription>
              </Alert>
            )}

            {!startLocation && (
              <Alert>
                <AlertCircle className="h-4 w-4" />
                <AlertDescription className="text-xs">
                  Click &quot;Use Current Location&quot; to capture your starting point
                </AlertDescription>
              </Alert>
            )}
          </div>

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label>Destinations ({waypoints.length})</Label>
              <div className="flex gap-2">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => setIsWaypointSelectorOpen(true)}
                >
                  <Calendar className="h-4 w-4 mr-2" />
                  From Visit Reports
                </Button>
              </div>
            </div>
            <WaypointList
              waypoints={waypoints}
              onRemove={handleRemoveWaypoint}
              onAdd={handleAddWaypoint}
            />
            {errors.waypoints && (
              <p className="text-sm text-destructive">
                {errors.waypoints.message}
              </p>
            )}
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isLoading}
            >
              Cancel
            </Button>
            <Button 
              type="submit" 
              disabled={isLoading || !startLocation || waypoints.length < 1}
            >
              {isLoading ? "Optimizing..." : "Optimize Route"}
            </Button>
          </DialogFooter>
        </form>

        <WaypointSelectorDialog
          open={isWaypointSelectorOpen}
          onOpenChange={setIsWaypointSelectorOpen}
          onSelect={handleSelectWaypoints}
        />
      </DialogContent>
    </Dialog>
  );
}
