"use client";

import { useState, useRef, useCallback, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Camera, X, RotateCcw } from "lucide-react";
import { useTranslations } from "next-intl";
import { detectFakeGPSFromPosition } from "../utils/detectFakeGPS";
import { FakeGPSWarningModal } from "./fake-gps-warning-modal";

interface CheckInCameraDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onCapture: (file: File, deviceGPS: {
    latitude: number;
    longitude: number;
    accuracy?: number;
    timestamp: number;
  }) => Promise<void>;
  readonly isLoading?: boolean;
}

export function CheckInCameraDialog({
  open,
  onOpenChange,
  onCapture,
  isLoading,
}: CheckInCameraDialogProps) {
  const [stream, setStream] = useState<MediaStream | null>(null);
  const [capturedImage, setCapturedImage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [gpsError, setGpsError] = useState<string | null>(null);
  const [isLoadingCamera, setIsLoadingCamera] = useState(false);
  const [deviceGPS, setDeviceGPS] = useState<{
    latitude: number;
    longitude: number;
    accuracy?: number;
    timestamp: number;
  } | null>(null);
  const [isFakeGPSModalOpen, setIsFakeGPSModalOpen] = useState(false);
  const [fakeGPSReason, setFakeGPSReason] = useState<string | undefined>();
  const [previousGPSPosition, setPreviousGPSPosition] = useState<GeolocationPosition | undefined>();
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const capturedBlobRef = useRef<Blob | null>(null);
  const gpsWatchRef = useRef<number | null>(null);
  const t = useTranslations("checkInCameraDialog");

  // Callback ref to set stream immediately when video element is mounted
  const setVideoRef = useCallback((node: HTMLVideoElement | null) => {
    videoRef.current = node;
    if (node && streamRef.current) {
      // Check if stream is still active before setting
      const stream = streamRef.current;
      const isStreamActive = stream.getTracks().some(track => track.readyState === 'live');
      
      if (isStreamActive && open) {
        try {
          node.srcObject = stream;
          node.play().catch((err) => {
            // Ignore AbortError - it's usually from cleanup
            if (err.name !== 'AbortError') {
            }
          });
        } catch (err) {
          // Ignore errors if stream is already stopped
          if (err instanceof Error && err.name !== 'AbortError') {
          }
        }
      }
    } else if (node && !streamRef.current) {
      // Clear srcObject if no stream
      node.srcObject = null;
    }
  }, [open]);

  const startCamera = useCallback(async () => {
    try {
      setError(null);
      setIsLoadingCamera(true);

      // Check if getUserMedia is supported
      if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        throw new Error(t("errors.cameraNotSupported"));
      }

      // Try front camera first, fallback to any camera
      let mediaStream: MediaStream | null = null;
      try {
        mediaStream = await navigator.mediaDevices.getUserMedia({
          video: {
            facingMode: "user", // Front camera for selfie
            width: { ideal: 1280 },
            height: { ideal: 720 },
          },
          audio: false,
        });
      } catch (frontCameraError) {
        // Fallback to any available camera
        try {
          mediaStream = await navigator.mediaDevices.getUserMedia({
            video: {
              width: { ideal: 1280 },
              height: { ideal: 720 },
            },
            audio: false,
          });
        } catch (anyCameraError) {
          // If still fails, try with minimal constraints
          mediaStream = await navigator.mediaDevices.getUserMedia({
            video: true,
            audio: false,
          });
        }
      }

      if (!mediaStream) {
        throw new Error(t("errors.cameraFailed"));
      }

      // Store stream in ref for callback ref
      streamRef.current = mediaStream;
      
      // Set stream to state
      setStream(mediaStream);
      
      // Immediately try to set stream to video element if available
      if (videoRef.current) {
        const video = videoRef.current;
        video.srcObject = mediaStream;
        
        // Function to play video
        const playVideo = async () => {
          try {
            if (video.paused) {
              await video.play();
            }
            setIsLoadingCamera(false);
          } catch (playError: unknown) {
            setIsLoadingCamera(false);
          }
        };

        // Try to play immediately if ready
        if (video.readyState >= 2) {
          playVideo();
        } else {
          // Wait for video to be ready
          const onLoadedMetadata = () => {
            video.removeEventListener('loadedmetadata', onLoadedMetadata);
            playVideo();
          };
          
          const onCanPlay = () => {
            video.removeEventListener('canplay', onCanPlay);
            playVideo();
          };

          const onPlay = () => {
            video.removeEventListener('play', onPlay);
            setIsLoadingCamera(false);
          };

          video.addEventListener('loadedmetadata', onLoadedMetadata, { once: true });
          video.addEventListener('canplay', onCanPlay, { once: true });
          video.addEventListener('play', onPlay, { once: true });
          
          // Fallback timeout
          setTimeout(() => {
            if (video.paused && video.readyState >= 2) {
              playVideo();
            } else {
              setIsLoadingCamera(false);
            }
          }, 1000);
        }
      } else {
        setIsLoadingCamera(false);
      }

      // GPS is fetched independently — see startGPS()
    } catch (err) {
      setIsLoadingCamera(false);
      let errorMessage = t("errors.cameraFailed");
      
      if (err instanceof Error) {
        if (err.name === "NotAllowedError" || err.name === "PermissionDeniedError") {
          errorMessage = t("errors.cameraPermissionDenied");
        } else if (err.name === "NotFoundError" || err.name === "DevicesNotFoundError") {
          errorMessage = t("errors.cameraNotFound");
        } else if (err.name === "NotReadableError" || err.name === "TrackStartError") {
          errorMessage = t("errors.cameraInUse");
        } else {
          errorMessage = err.message || errorMessage;
        }
      }
      
      setError(errorMessage);
    }
  }, [t]);

  // Stop GPS watch
  const stopGPS = useCallback(() => {
    if (gpsWatchRef.current !== null) {
      navigator.geolocation.clearWatch(gpsWatchRef.current);
      gpsWatchRef.current = null;
    }
  }, []);

  // Fetch GPS independently from camera — uses watchPosition to progressively
  // refine accuracy until we get < 50m (target < 10m).
  const startGPS = useCallback(() => {
    if (!navigator.geolocation) {
      setGpsError(t("errors.gpsFailed"));
      return;
    }

    // Clear any previous watch
    stopGPS();
    setGpsError(null);

    const TARGET_ACCURACY = 50; // stop watching once accuracy <= 50m
    const MAX_WATCH_TIME = 20000; // give up refining after 20s
    let bestAccuracy = Infinity;

    // Timeout: if after MAX_WATCH_TIME we still don't have target accuracy,
    // accept best reading so far (if any) or show error
    const timeoutId = setTimeout(() => {
      stopGPS();
      // If we already set a GPS position (even if > target), keep it — it's the best we got
    }, MAX_WATCH_TIME);

    const watchId = navigator.geolocation.watchPosition(
      (position) => {
        const accuracy = position.coords.accuracy;

        // Only update if this reading is better than previous
        if (accuracy >= bestAccuracy) return;
        bestAccuracy = accuracy;

        // Detect Fake GPS
        const fakeGPSDetection = detectFakeGPSFromPosition(position, previousGPSPosition);

        if (fakeGPSDetection.isFakeGPS) {
          setFakeGPSReason(fakeGPSDetection.reason);
          setIsFakeGPSModalOpen(true);
          setDeviceGPS(null);
          setGpsError(t("errors.gpsFailed"));
          stopGPS();
          clearTimeout(timeoutId);
          return;
        }

        // Update GPS with better reading
        setDeviceGPS({
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
          accuracy: position.coords.accuracy,
          timestamp: Math.floor(Date.now() / 1000),
        });
        setGpsError(null);
        setPreviousGPSPosition(position);

        // If accuracy is good enough, stop watching to save battery
        if (accuracy <= TARGET_ACCURACY) {
          stopGPS();
          clearTimeout(timeoutId);
        }
      },
      () => {
        clearTimeout(timeoutId);
        // Only show error if we don't already have a reading
        setGpsError((prev) => prev ?? t("errors.gpsFailed"));
      },
      {
        enableHighAccuracy: true,
        timeout: 15000,
        maximumAge: 0,
      }
    );

    gpsWatchRef.current = watchId;
  }, [t, previousGPSPosition, stopGPS]);

  const stopCamera = useCallback(() => {
    // Clear video srcObject first to prevent AbortError
    if (videoRef.current) {
      try {
        videoRef.current.srcObject = null;
      } catch (err) {
        // Ignore errors during cleanup
      }
    }
    
    // Stop all tracks
    if (streamRef.current) {
      streamRef.current.getTracks().forEach((track) => {
        track.stop();
      });
      streamRef.current = null;
    }
    if (stream) {
      stream.getTracks().forEach((track) => {
        track.stop();
      });
      setStream(null);
    }
  }, [stream]);

  const capturePhoto = useCallback(() => {
    if (!videoRef.current || !canvasRef.current) return;

    const video = videoRef.current;
    const canvas = canvasRef.current;
    const context = canvas.getContext("2d");

    if (!context) return;

    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;
    context.drawImage(video, 0, 0);

    canvas.toBlob((blob) => {
      if (blob) {
        capturedBlobRef.current = blob;
        const imageUrl = URL.createObjectURL(blob);
        setCapturedImage(imageUrl);
        stopCamera();
      }
    }, "image/jpeg", 0.9);
  }, [stopCamera]);

  const retakePhoto = useCallback(() => {
    capturedBlobRef.current = null;
    setCapturedImage(null);
    startCamera();
  }, [startCamera]);

  const handleConfirm = useCallback(async () => {
    if (!capturedImage || !deviceGPS) {
      setError(t("errors.missingData"));
      return;
    }

    try {
      // Use the stored blob directly — no fetch needed (avoids blob URL errors)
      const blob = capturedBlobRef.current;
      if (!blob) {
        setError(t("errors.uploadFailed"));
        return;
      }
      const file = new File([blob], `checkin-${Date.now()}.jpg`, { type: "image/jpeg" });

      await onCapture(file, deviceGPS);
      handleClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : t("errors.uploadFailed"));
    }
  }, [capturedImage, deviceGPS, onCapture, t]);

  const handleClose = useCallback(() => {
    stopCamera();
    stopGPS();
    capturedBlobRef.current = null;
    setCapturedImage(null);
    setError(null);
    setGpsError(null);
    setDeviceGPS(null);
    onOpenChange(false);
  }, [stopCamera, stopGPS, onOpenChange]);

  // Ensure video plays when stream is set
  useEffect(() => {
    if (!stream || !videoRef.current || !open) return;

    const video = videoRef.current;
    const currentStream = streamRef.current || stream;
    
    // Check if stream is still active
    const isStreamActive = currentStream.getTracks().some(track => track.readyState === 'live');
    if (!isStreamActive) {
      return;
    }
    
    // Set stream if not already set
    if (video.srcObject !== currentStream) {
      try {
        video.srcObject = currentStream;
      } catch (err) {
        if (err instanceof Error && err.name !== 'AbortError') {
        }
        return;
      }
    }
    
    // Function to ensure video plays
    const ensurePlay = async () => {
      // Check if still mounted and stream is active
      if (!videoRef.current || !open) return;
      
      const isStillActive = currentStream.getTracks().some(track => track.readyState === 'live');
      if (!isStillActive) return;
      
      try {
        if (video.paused) {
          await video.play();
        } else {
        }
      } catch (err) {
        // Ignore AbortError - it's usually from cleanup
        if (err instanceof Error && err.name !== 'AbortError') {
        }
      }
    };

    // Try to play immediately if ready
    if (video.readyState >= 2) {
      ensurePlay();
    } else {
      // Wait for video to be ready
      const onLoadedMetadata = () => {
        if (!videoRef.current || !open) return;
        video.removeEventListener('loadedmetadata', onLoadedMetadata);
        ensurePlay();
      };
      
      const onCanPlay = () => {
        if (!videoRef.current || !open) return;
        video.removeEventListener('canplay', onCanPlay);
        ensurePlay();
      };

      video.addEventListener('loadedmetadata', onLoadedMetadata, { once: true });
      video.addEventListener('canplay', onCanPlay, { once: true });
      
      // Fallback
      const timeoutId = setTimeout(() => {
        if (!videoRef.current || !open) return;
        if (video.readyState >= 2 && video.paused) {
          ensurePlay();
        }
      }, 500);
      
      return () => {
        clearTimeout(timeoutId);
        video.removeEventListener('loadedmetadata', onLoadedMetadata);
        video.removeEventListener('canplay', onCanPlay);
      };
    }
  }, [stream, open]);

  // Start camera and GPS when dialog opens — independently
  useEffect(() => {
    if (!open) {
      stopCamera();
      return;
    }

    // Small delay to ensure dialog is fully rendered
    const timer = setTimeout(() => {
      startCamera();
      startGPS();
    }, 100);

    return () => {
      clearTimeout(timer);
      stopCamera();
      stopGPS();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open]); // Only depend on open to avoid re-initializing camera unnecessarily

  // Listen for camera permission state changes (e.g. user toggles permission in browser settings)
  // and auto-retry when permission becomes 'granted'
  useEffect(() => {
    if (!open) return;

    let permissionStatus: PermissionStatus | null = null;
    let cancelled = false;

    const listenPermission = async () => {
      try {
        permissionStatus = await navigator.permissions.query({ name: 'camera' as PermissionName });
        if (cancelled) return;

        const onChange = () => {
          if (permissionStatus?.state === 'granted' && !streamRef.current) {
            // Permission was just granted and we don't have a stream — auto-retry
            setError(null);
            startCamera();
          }
        };

        permissionStatus.addEventListener('change', onChange);

        // Cleanup
        return () => {
          permissionStatus?.removeEventListener('change', onChange);
        };
      } catch {
        // Permissions API not supported for camera in this browser — ignored
      }
    };

    const cleanupPromise = listenPermission();

    return () => {
      cancelled = true;
      cleanupPromise.then((cleanup) => cleanup?.());
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, startCamera]);

  // Auto-retry camera when user returns to this tab after changing permissions
  useEffect(() => {
    if (!open) return;

    const handleVisibilityChange = () => {
      if (document.visibilityState === 'visible' && !streamRef.current && error) {
        // User came back to tab and we still have a camera error — retry
        setError(null);
        startCamera();
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    return () => document.removeEventListener('visibilitychange', handleVisibilityChange);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, error, startCamera]);

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {error && (
            <div className="bg-destructive/10 text-destructive text-sm p-3 rounded-md flex items-center justify-between">
              <span>{error}</span>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  setError(null);
                  startCamera();
                }}
                className="ml-2 h-auto py-1 shrink-0"
              >
                {t("buttons.retry")}
              </Button>
            </div>
          )}

          {gpsError && !deviceGPS && (
            <div className="bg-destructive/10 text-destructive text-sm p-3 rounded-md flex items-center justify-between">
              <span>{gpsError}</span>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  setGpsError(null);
                  startGPS();
                }}
                className="ml-2 h-auto py-1 shrink-0"
              >
                {t("buttons.retry")}
              </Button>
            </div>
          )}

          {!deviceGPS && !gpsError && (
            <div className="bg-yellow-500/10 text-yellow-600 dark:text-yellow-400 text-sm p-3 rounded-md">
              {t("warnings.gettingLocation")}
            </div>
          )}

          {deviceGPS && (
            <div className={`text-sm p-3 rounded-md ${
              (deviceGPS.accuracy ?? 999) <= 50
                ? "bg-[color:var(--color-success)]/10 text-[color:var(--color-success)] dark:text-[color:var(--color-success-foreground)]"
                : "bg-yellow-500/10 text-yellow-600 dark:text-yellow-400"
            }`}>
              <div className="flex items-center justify-between">
                <span>
                  {t("success.locationCaptured")} (Accuracy: {deviceGPS.accuracy?.toFixed(0) ?? "N/A"}m)
                </span>
                {(deviceGPS.accuracy ?? 999) > 50 && gpsWatchRef.current !== null && (
                  <span className="text-xs opacity-75 flex items-center gap-1">
                    <span className="inline-block h-2 w-2 rounded-full bg-yellow-500 animate-pulse" />
                    Refining...
                  </span>
                )}
              </div>
              {(deviceGPS.accuracy ?? 999) > 100 && (
                <p className="text-xs mt-1 opacity-75">
                  ⚠ Accuracy &gt;100m — check-in may fail. Move near a window or open area for better GPS signal.
                </p>
              )}
            </div>
          )}

          <div className="relative bg-black rounded-lg overflow-hidden aspect-video">
            {capturedImage ? (
              <img
                src={capturedImage}
                alt="Captured selfie"
                className="w-full h-full object-cover"
              />
            ) : isLoadingCamera ? (
              <div className="w-full h-full flex flex-col items-center justify-center text-white">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mb-2"></div>
                <p className="text-sm opacity-75">{t("warnings.loadingCamera")}</p>
              </div>
            ) : stream ? (
              <video
                ref={setVideoRef}
                autoPlay
                playsInline
                muted
                className="w-full h-full object-cover"
                style={{ transform: 'scaleX(-1)', display: 'block' }} // Mirror for selfie, ensure visible
                onLoadedMetadata={(e) => {
                  const video = e.currentTarget;
                  video.play().catch((err) => {
                  });
                }}
                onCanPlay={(e) => {
                  const video = e.currentTarget;
                  if (video.paused) {
                    video.play().catch((err) => {
                    });
                  }
                }}
                onPlay={() => {
                  setIsLoadingCamera(false);
                }}
                onError={(e) => {
                  setError("Failed to load video stream");
                }}
              />
            ) : (
              <div className="w-full h-full flex flex-col items-center justify-center text-white">
                <Camera className="h-16 w-16 opacity-50 mb-2" />
                <p className="text-sm opacity-75">{t("warnings.cameraNotReady")}</p>
              </div>
            )}
          </div>

          <div className="flex gap-2 justify-end">
            {capturedImage ? (
              <>
                <Button
                  variant="outline"
                  onClick={retakePhoto}
                  disabled={isLoading}
                >
                  <RotateCcw className="h-4 w-4 mr-2" />
                  {t("buttons.retake")}
                </Button>
                <Button
                  onClick={handleConfirm}
                  disabled={isLoading || !deviceGPS}
                >
                  {isLoading ? t("buttons.processing") : t("buttons.confirm")}
                </Button>
              </>
            ) : (
              <>
                <Button
                  variant="outline"
                  onClick={handleClose}
                  disabled={isLoading}
                >
                  {t("buttons.cancel")}
                </Button>
                <Button
                  onClick={capturePhoto}
                  disabled={!stream || isLoading}
                >
                  <Camera className="h-4 w-4 mr-2" />
                  {t("buttons.capture")}
                </Button>
              </>
            )}
          </div>
        </div>

        <canvas ref={canvasRef} className="hidden" />
      </DialogContent>
      
      {/* Fake GPS Warning Modal */}
      <FakeGPSWarningModal
        open={isFakeGPSModalOpen}
        onOpenChange={setIsFakeGPSModalOpen}
        reason={fakeGPSReason}
      />
    </Dialog>
  );
}
