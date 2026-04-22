"use client";

import { useState } from "react";
import { MapPin, Navigation, Loader2, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

interface LocationShareCardProps {
  readonly onLocationShared: (message: string) => void;
}

type LocationState = "idle" | "loading" | "error";

export function LocationShareCard({ onLocationShared }: LocationShareCardProps) {
  const [state, setState] = useState<LocationState>("idle");
  const [errorMsg, setErrorMsg] = useState<string>("");

  const handleShareLocation = () => {
    if (!navigator.geolocation) {
      setErrorMsg("Browser Anda tidak mendukung layanan lokasi.");
      setState("error");
      return;
    }

    setState("loading");
    setErrorMsg("");

    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords;
        const message = `Saya berada di ${latitude.toFixed(6)}, ${longitude.toFixed(6)}`;
        setState("idle");
        onLocationShared(message);
      },
      (error) => {
        setState("error");
        switch (error.code) {
          case error.PERMISSION_DENIED:
            setErrorMsg("Akses lokasi ditolak. Izinkan akses lokasi di browser Anda lalu coba lagi.");
            break;
          case error.POSITION_UNAVAILABLE:
            setErrorMsg("Informasi lokasi tidak tersedia saat ini.");
            break;
          case error.TIMEOUT:
            setErrorMsg("Permintaan lokasi habis waktu. Silakan coba lagi.");
            break;
          default:
            setErrorMsg("Gagal mendapatkan lokasi. Silakan coba lagi.");
        }
      },
      { enableHighAccuracy: true, timeout: 10000, maximumAge: 60000 }
    );
  };

  return (
    <div className="mt-3 rounded-xl border border-primary/20 bg-primary/5 p-4 flex flex-col gap-3">
      <div className="flex items-start gap-3">
        <div className="shrink-0 w-9 h-9 rounded-lg bg-primary/10 flex items-center justify-center">
          <MapPin className="w-4 h-4 text-primary" />
        </div>
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium text-foreground">Berbagi Lokasi</p>
          <p className="text-xs text-muted-foreground mt-0.5">
            Izinkan akses lokasi untuk merencanakan rute kunjungan yang optimal
          </p>
        </div>
      </div>

      {state === "error" && (
        <div className="flex items-start gap-2 rounded-lg bg-destructive/10 px-3 py-2 text-xs text-destructive">
          <AlertCircle className="w-3.5 h-3.5 mt-0.5 shrink-0" />
          <span>{errorMsg}</span>
        </div>
      )}

      <Button
        size="sm"
        variant="outline"
        className="w-full border-primary/30 text-primary hover:bg-primary/10 hover:text-primary gap-2"
        onClick={handleShareLocation}
        disabled={state === "loading"}
      >
        {state === "loading" ? (
          <>
            <Loader2 className="w-3.5 h-3.5 animate-spin" />
            Mendapatkan Lokasi...
          </>
        ) : (
          <>
            <Navigation className="w-3.5 h-3.5" />
            Aktifkan &amp; Kirim Lokasi Saya
          </>
        )}
      </Button>
    </div>
  );
}

/**
 * Parses the LOCATION_NEEDED marker from an AI message.
 * Returns cleaned message (marker removed) and whether the card should be shown.
 */
export function parseLocationNeeded(message: string): {
  cleanMessage: string;
  needsLocation: boolean;
} {
  const marker = "<!-- LOCATION_NEEDED -->";
  const needsLocation = message.includes(marker);
  return {
    cleanMessage: message.replace(marker, "").trim(),
    needsLocation,
  };
}
