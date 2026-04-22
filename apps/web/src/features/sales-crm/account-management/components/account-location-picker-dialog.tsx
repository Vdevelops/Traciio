"use client";

import dynamic from "next/dynamic";
import { useMemo, useState } from "react";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { SmartTileLayer, TILE_SOURCES, LIGHT_FALLBACK_CHAIN } from "@/components/ui/smart-tile-layer";

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

const MapContainer = dynamic(
  () => import("react-leaflet").then((mod) => mod.MapContainer),
  { ssr: false }
);

const Marker = dynamic(
  () => import("react-leaflet").then((mod) => mod.Marker),
  { ssr: false }
);

// Hook must be imported directly (not dynamic)
import { useMapEvents } from "react-leaflet";

type PickedLocation = { lat: number; lng: number };

function ClickToPick({ onPick }: { onPick: (loc: PickedLocation) => void }) {
  useMapEvents({
    click: (e) => {
      onPick({ lat: e.latlng.lat, lng: e.latlng.lng });
    },
  });

  return null;
}

interface AccountLocationPickerDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly initialLocation?: PickedLocation;
  readonly onConfirm: (location: PickedLocation) => void;
  readonly title: string;
  readonly confirmText: string;
  readonly cancelText: string;
}

export function AccountLocationPickerDialog({
  open,
  onOpenChange,
  initialLocation,
  onConfirm,
  title,
  confirmText,
  cancelText,
}: AccountLocationPickerDialogProps) {
  const defaultCenter = useMemo<PickedLocation>(() => {
    return initialLocation ?? { lat: -6.2088, lng: 106.8456 };
  }, [initialLocation]);

  const [picked, setPicked] = useState<PickedLocation>(defaultCenter);

  return (
    <Dialog
      open={open}
      onOpenChange={(next) => {
        onOpenChange(next);
        if (next) {
          setPicked(defaultCenter);
        }
      }}
    >
      <DialogContent className="w-[95vw] sm:max-w-4xl lg:max-w-5xl">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
        </DialogHeader>

        <div className="h-[60vh] min-h-[480px] w-full overflow-hidden rounded-md border">
          <MapContainer
            center={[picked.lat, picked.lng]}
            zoom={13}
            style={{ height: "100%", width: "100%" }}
          >
            <SmartTileLayer source={TILE_SOURCES.cartoLight} fallbackSources={LIGHT_FALLBACK_CHAIN.slice(1)} />
            <ClickToPick onPick={setPicked} />
            <Marker position={[picked.lat, picked.lng]} />
          </MapContainer>
        </div>

        <div className="text-sm text-muted-foreground">
          {picked.lat.toFixed(6)}, {picked.lng.toFixed(6)}
        </div>

        <DialogFooter>
          <Button variant="outline" type="button" onClick={() => onOpenChange(false)}>
            {cancelText}
          </Button>
          <Button type="button" onClick={() => onConfirm(picked)}>
            {confirmText}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
