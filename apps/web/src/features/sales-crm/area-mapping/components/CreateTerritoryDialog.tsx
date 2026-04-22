"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useCreateTerritory } from "../hooks/useAreaMapping";
import { AreaMap } from "./AreaMap";
import { ManualPolygonInput } from "./ManualPolygonInput";
import type { GeoJSONPolygon } from "../types";

export function CreateTerritoryDialog() {
  const t = useTranslations("areaMapping.territories.create");
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [color, setColor] = useState("#3B82F6");
  const [polygon, setPolygon] = useState<GeoJSONPolygon | null>(null);
  
  const createTerritory = useCreateTerritory();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // Validate polygon
    if (!polygon || !polygon.coordinates || polygon.coordinates.length === 0) {
      alert("Please draw a polygon on the map to define the territory boundary");
      return;
    }

    const firstRing = polygon.coordinates[0];
    if (!firstRing || firstRing.length < 4) {
      alert("Polygon must have at least 4 points (including closing point)");
      return;
    }

    const territoryData = {
      name,
      description,
      color,
      polygon: polygon,
    };
    

    try {
      await createTerritory.mutateAsync(territoryData);
      
      setOpen(false);
      setName("");
      setDescription("");
      setColor("#3B82F6");
      setPolygon(null);
    } catch (error) {
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          {t("button")}
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">{t("form.name.label")}</Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t("form.name.placeholder")}
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">{t("form.description.label")}</Label>
            <Textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder={t("form.description.placeholder")}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="color">{t("form.color.label")}</Label>
            <div className="flex gap-2">
              <Input
                id="color"
                type="color"
                value={color}
                onChange={(e) => setColor(e.target.value)}
                className="h-10 w-20"
              />
              <Input
                value={color}
                onChange={(e) => setColor(e.target.value)}
                placeholder="#3B82F6"
                className="flex-1"
              />
            </div>
            <p className="text-sm text-muted-foreground">
              {t("form.color.description")}
            </p>
          </div>

          <div className="space-y-2">
            <Label>Territory Boundary</Label>
            <div className="h-[400px] border rounded-lg overflow-hidden">
              <AreaMap
                onPolygonChange={setPolygon}
                initialPolygon={polygon || undefined}
              />
            </div>
            <div className="flex items-center justify-between">
              <div>
                {polygon && polygon.coordinates && polygon.coordinates.length > 0 && polygon.coordinates[0] && (
                  <p className="text-xs text-green-600">
                    ✓ Polygon defined with {polygon.coordinates[0].length - 1} points
                  </p>
                )}
              </div>
              <ManualPolygonInput
                onPolygonChange={(coordinates) => {
                  setPolygon({
                    type: "Polygon",
                    coordinates: coordinates as [number, number][][],
                  });
                }}
                currentPolygon={polygon?.coordinates}
              />
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setOpen(false)}
            >
              {t("form.cancel")}
            </Button>
            <Button 
              type="submit" 
              disabled={createTerritory.isPending || !polygon || !polygon.coordinates || polygon.coordinates.length === 0}
            >
              {createTerritory.isPending ? t("form.creating") : t("form.create")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
