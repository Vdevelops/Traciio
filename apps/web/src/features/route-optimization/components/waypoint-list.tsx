"use client";

import { GripVertical, Trash2, Plus, MapPin } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import type { Waypoint } from "../types";

interface WaypointListProps {
  readonly waypoints: Waypoint[];
  readonly onRemove?: (index: number) => void;
  readonly onReorder?: (fromIndex: number, toIndex: number) => void;
  readonly onAdd?: () => void;
  readonly readonly?: boolean;
}

export function WaypointList({
  waypoints,
  onRemove,
  onReorder,
  onAdd,
  readonly = false,
}: WaypointListProps) {
  return (
    <div className="space-y-2">
      {waypoints.map((waypoint, index) => (
        <Card key={index} className="p-3">
          <CardContent className="p-0">
            <div className="flex items-start gap-3">
              {!readonly && (
                <div className="flex flex-col items-center gap-1 pt-1">
                  <GripVertical className="h-4 w-4 text-muted-foreground cursor-move" />
                  <Badge variant="outline" className="text-xs">
                    {index + 1}
                  </Badge>
                </div>
              )}
              {readonly && (
                <Badge variant="outline" className="mt-1">
                  {waypoint.order ?? index + 1}
                </Badge>
              )}

              <div className="flex-1 space-y-1">
                <div className="flex items-center gap-2">
                  <MapPin className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium">
                    {waypoint.address || `${waypoint.lat}, ${waypoint.lng}`}
                  </span>
                </div>
                {waypoint.account?.name && (
                  <p className="text-xs text-muted-foreground">
                    Account: {waypoint.account.name}
                  </p>
                )}
                <p className="text-xs text-muted-foreground">
                  Lat: {waypoint.lat.toFixed(6)}, Lng: {waypoint.lng.toFixed(6)}
                </p>
              </div>

              {!readonly && onRemove && (
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => onRemove(index)}
                  className="h-8 w-8 text-destructive hover:text-destructive"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              )}
            </div>
          </CardContent>
        </Card>
      ))}

      {!readonly && onAdd && (
        <Button
          variant="outline"
          className="w-full"
          onClick={onAdd}
        >
          <Plus className="h-4 w-4 mr-2" />
          Add Waypoint
        </Button>
      )}

      {waypoints.length === 0 && (
        <div className="text-center py-8 text-muted-foreground">
          <MapPin className="h-8 w-8 mx-auto mb-2 opacity-50" />
          <p>No waypoints added yet</p>
        </div>
      )}
    </div>
  );
}


