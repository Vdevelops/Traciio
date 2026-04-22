"use client";

import { useMemo, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { MapPin, Route, Clock, Navigation, ChevronDown, ChevronUp, Building2, Play } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Drawer } from "@/components/ui/drawer";
import type { OptimizedRoute, RouteStep, Waypoint } from "../types";

interface RouteDetailsPanelProps {
  readonly route: OptimizedRoute | null;
  readonly isOpen: boolean;
  readonly onClose: () => void;
  readonly onWaypointClick?: (waypoint: Waypoint, index: number) => void;
  readonly onStepClick?: (step: RouteStep, index: number) => void;
}

export function RouteDetailsPanel({ 
  route, 
  isOpen, 
  onClose,
  onWaypointClick,
  onStepClick,
}: RouteDetailsPanelProps) {
  const [isStepsExpanded, setIsStepsExpanded] = useState(false);

  const waypoints = route?.waypoints ?? [];

  // Backend returns waypoints already in optimized order:
  // [start(order=0), optimized_stop_1(order=1), optimized_stop_2(order=2), ...]
  // We display them in their natural array order — no re-indexing needed.
  const orderedWaypoints = useMemo(() => {
    if (waypoints.length === 0) return [];

    const result: Array<{
      waypoint: Waypoint;
      originalIndex: number;
      isStart: boolean;
      displayOrder: number;
    }> = [];

    waypoints.forEach((wp, index) => {
      const isStart = wp.order === 0 || (wp.order == null && index === 0);
      result.push({
        waypoint: wp,
        originalIndex: index,
        isStart,
        displayOrder: wp.order ?? index,
      });
    });

    // Sort by order field to ensure start is first and destinations follow
    result.sort((a, b) => a.displayOrder - b.displayOrder);

    return result;
  }, [waypoints]);

  if (!route) return null;

  return (
    <Drawer
      open={isOpen}
      onOpenChange={(open) => !open && onClose()}
      side="right"
      title={route.route_name || "Route Details"}
      description={`${route.waypoints.length} waypoints`}
      defaultWidth={420}
      minWidth={360}
      maxWidth={600}
    >
      <div className="space-y-6">
        {/* Route Stats - Using semantic colors */}
        <div className="grid grid-cols-2 gap-3">
          <div className="bg-primary/10 rounded-xl p-4 border border-primary/20">
            <div className="flex items-center gap-2 text-primary mb-1">
              <Navigation className="w-4 h-4" />
              <span className="text-xs font-medium uppercase tracking-wide">Distance</span>
            </div>
            <p className="text-2xl font-medium text-foreground">
              {route.total_distance_formatted ?? `${route.total_distance?.toFixed(2) ?? 0} km`}
            </p>
          </div>
          <div className="bg-success/10 rounded-xl p-4 border border-success/20">
            <div className="flex items-center gap-2 text-success mb-1">
              <Clock className="w-4 h-4" />
              <span className="text-xs font-medium uppercase tracking-wide">Duration</span>
            </div>
            <p className="text-2xl font-medium text-foreground">
              {route.total_duration_formatted ?? `${Math.round((route.total_duration ?? 0) / 60)} min`}
            </p>
          </div>
        </div>

        {/* Waypoints List */}
        <div>
          <h3 className="text-sm font-medium text-muted-foreground mb-3 uppercase tracking-wide">
            Waypoints
          </h3>
          <div className="relative">
            {orderedWaypoints.map((item, index) => {
              const waypoint = item.waypoint;
              const isStartLocation = item.isStart;
              const displayOrder = item.displayOrder;
              const isLastItem = index === orderedWaypoints.length - 1;
              
              return (
                <div key={index} className="relative">
                  {/* Connection Line - Centered with badge */}
                  {/* Badge is w-10 (40px), card has p-4 (16px) padding */}
                  {/* Center = 16px padding + 20px (half badge) = 36px */}
                  {!isLastItem && (
                    <div 
                      className="absolute w-0.5 bg-border"
                      style={{
                        left: '36px', // 16px (padding) + 20px (half of 40px badge) = center
                        top: '60px', // After badge height
                        bottom: '-8px',
                      }}
                    />
                  )}
                  
                  <motion.div
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: index * 0.05 }}
                    onClick={() => onWaypointClick?.(waypoint, item.originalIndex)}
                    className={`
                      relative cursor-pointer rounded-xl p-4 mb-2
                      border hover:shadow-md active:scale-[0.99]
                      ${isStartLocation 
                        ? 'bg-success/5 border-success/30 hover:border-success/50' 
                        : 'bg-card hover:bg-muted/50 border-border/50 hover:border-primary/30'
                      }
                    `}
                  >
                    <div className="flex items-start gap-3">
                      {/* Order Badge - Using semantic colors */}
                      <div className={`
                        w-10 h-10 rounded-xl flex items-center justify-center text-sm font-medium shadow-sm shrink-0
                        ${isStartLocation 
                          ? 'bg-success text-success-foreground' 
                          : 'bg-primary text-primary-foreground'
                        }
                      `}>
                        {isStartLocation ? (
                          <Play className="w-5 h-5 fill-current" />
                        ) : (
                          displayOrder
                        )}
                      </div>

                      {/* Content */}
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1">
                          {isStartLocation ? (
                            <span className="font-medium text-success">
                              Start Location
                            </span>
                          ) : (
                            <span className="font-medium text-foreground">
                              {waypoint.account?.name || `Stop ${displayOrder}`}
                            </span>
                          )}
                        </div>

                        <div className="flex items-start gap-1.5 text-muted-foreground">
                          <MapPin className="w-3.5 h-3.5 mt-0.5 shrink-0" />
                          <span className="text-xs leading-relaxed truncate">
                            {waypoint.address || `${waypoint.lat.toFixed(4)}, ${waypoint.lng.toFixed(4)}`}
                          </span>
                        </div>

                        {waypoint.account?.name && !isStartLocation && (
                          <div className="mt-2">
                            <Badge variant="secondary" className="text-xs">
                              <Building2 className="w-3 h-3 mr-1" />
                              {waypoint.account.name}
                            </Badge>
                          </div>
                        )}
                      </div>
                    </div>
                  </motion.div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Route Steps (Collapsible) */}
        {route.route_steps && route.route_steps.length > 0 && (
          <div>
            <button
              onClick={() => setIsStepsExpanded(!isStepsExpanded)}
              className="flex items-center justify-between w-full text-sm font-medium text-muted-foreground mb-3 uppercase tracking-wide hover:text-foreground transition-colors cursor-pointer"
            >
              <span>Turn-by-Turn Directions ({route.route_steps.length})</span>
              {isStepsExpanded ? (
                <ChevronUp className="w-4 h-4" />
              ) : (
                <ChevronDown className="w-4 h-4" />
              )}
            </button>

            <AnimatePresence>
              {isStepsExpanded && (
                <motion.div
                  initial={{ height: 0, opacity: 0 }}
                  animate={{ height: "auto", opacity: 1 }}
                  exit={{ height: 0, opacity: 0 }}
                  transition={{ duration: 0.2 }}
                  className="overflow-hidden"
                >
                  <div className="space-y-2 max-h-60 overflow-y-auto">
                    {route.route_steps.map((step, index) => (
                      <div 
                        key={index} 
                        onClick={() => onStepClick?.(step, index)}
                        className="bg-muted/50 rounded-lg p-3 text-sm border border-border/30 cursor-pointer hover:bg-muted/70 transition-colors"
                      >
                        <p className="font-medium text-foreground">{step.instruction}</p>
                        <div className="flex gap-4 text-xs text-muted-foreground mt-1">
                          <span>{step.distance_formatted}</span>
                          <span>{step.duration_formatted}</span>
                        </div>
                      </div>
                    ))}
                  </div>
                </motion.div>
              )}
            </AnimatePresence>
          </div>
        )}
      </div>
    </Drawer>
  );
}
