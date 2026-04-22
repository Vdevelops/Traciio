"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

interface ManualPolygonInputProps {
  onPolygonChange: (coordinates: number[][][]) => void;
  currentPolygon?: number[][][];
}

export function ManualPolygonInput({ onPolygonChange, currentPolygon }: ManualPolygonInputProps) {
  const [open, setOpen] = useState(false);
  const [coordinates, setCoordinates] = useState(() => {
    if (currentPolygon && currentPolygon[0]) {
      return currentPolygon[0].map(coord => `${coord[0]}, ${coord[1]}`).join('\n');
    }
    return "";
  });

  const handleSubmit = () => {
    try {
      const lines = coordinates.trim().split('\n');
      const coords = lines.map(line => {
        const [lng, lat] = line.split(',').map(s => parseFloat(s.trim()));
        if (isNaN(lng) || isNaN(lat)) {
          throw new Error(`Invalid coordinate: ${line}`);
        }
        return [lng, lat];
      });

      if (coords.length < 3) {
        throw new Error("At least 3 points are required");
      }

      // Ensure polygon is closed
      const firstCoord = coords[0];
      const lastCoord = coords[coords.length - 1];
      if (firstCoord[0] !== lastCoord[0] || firstCoord[1] !== lastCoord[1]) {
        coords.push(firstCoord);
      }

      onPolygonChange([coords]);
      setOpen(false);
    } catch (error) {
      alert(`Error parsing coordinates: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          Manual Input
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Manual Polygon Input</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <div>
            <Label htmlFor="coordinates">Coordinates (longitude, latitude)</Label>
            <Textarea
              id="coordinates"
              placeholder="106.8, -6.2&#10;106.9, -6.2&#10;106.9, -6.3&#10;106.8, -6.3"
              value={coordinates}
              onChange={(e) => setCoordinates(e.target.value)}
              className="h-32 font-mono text-sm"
            />
            <p className="text-xs text-muted-foreground mt-1">
              Enter coordinates one per line in format: longitude, latitude
            </p>
          </div>
          <div className="flex gap-2 justify-end">
            <Button variant="outline" onClick={() => setOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleSubmit}>
              Apply
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}