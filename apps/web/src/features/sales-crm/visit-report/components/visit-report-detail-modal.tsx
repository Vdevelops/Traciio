"use client";

import { Calendar, MapPin, CheckCircle2, XCircle, Clock, User, Building2, FileText, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Drawer } from "@/components/ui/drawer";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  useVisitReport,
  useCheckIn,
  useCheckOut,
  useApproveVisitReport,
  useRejectVisitReport,
  useActivityTimeline,
  useUploadPhoto,
} from "../hooks/useVisitReports";
import { toast } from "sonner";
import { useState } from "react";
import { ActivityTimeline } from "./activity-timeline";
import { CreateActivityDialog } from "./create-activity-dialog";
import { PhotoUploadDialog } from "./photo-upload-dialog";
import { VisitReportInsightsButton } from "@/features/ai/components/visit-report-insights-button";
import { CheckInCameraDialog } from "./check-in-camera-dialog";
import { FakeGPSWarningModal } from "./fake-gps-warning-modal";
import { detectFakeGPSFromPosition } from "../utils/detectFakeGPS";
import { useTranslations } from "next-intl";

// Helper function to convert relative photo URL to absolute URL
const getPhotoUrl = (photoUrl: string): string => {
  // If already absolute URL, return as is
  if (photoUrl.startsWith("http://") || photoUrl.startsWith("https://")) {
    return photoUrl;
  }
  
  // Get API base URL from environment or default
  const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  
  // Ensure photoUrl starts with /
  const cleanUrl = photoUrl.startsWith("/") ? photoUrl : `/${photoUrl}`;
  
  // Return absolute URL (API_BASE_URL already includes protocol and domain)
  return `${API_BASE_URL}${cleanUrl}`;
};

const statusColors: Record<string, "default" | "secondary" | "destructive" | "outline"> = {
  draft: "outline",
  submitted: "secondary",
  approved: "default",
  rejected: "destructive",
};

interface VisitReportDetailModalProps {
  readonly visitReportId: string | null;
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onVisitReportUpdated?: () => void;
}

export function VisitReportDetailModal({
  visitReportId,
  open,
  onOpenChange,
  onVisitReportUpdated,
}: VisitReportDetailModalProps) {
  const { data, isLoading, error, refetch } = useVisitReport(visitReportId || "");
  const checkIn = useCheckIn();
  const checkOut = useCheckOut();
  const approve = useApproveVisitReport();
  const reject = useRejectVisitReport();
  const uploadPhoto = useUploadPhoto();
  const [isRejectDialogOpen, setIsRejectDialogOpen] = useState(false);
  const [isCreateActivityDialogOpen, setIsCreateActivityDialogOpen] = useState(false);
  const [isPhotoUploadDialogOpen, setIsPhotoUploadDialogOpen] = useState(false);
  const [isCheckInCameraDialogOpen, setIsCheckInCameraDialogOpen] = useState(false);
  const [isFakeGPSModalOpen, setIsFakeGPSModalOpen] = useState(false);
  const [fakeGPSReason, setFakeGPSReason] = useState<string | undefined>();
  const [previousGPSPosition, setPreviousGPSPosition] = useState<GeolocationPosition | undefined>();
  const [rejectReason, setRejectReason] = useState("");

  const visitReport = data?.data;
  
  // Debug: Log photos to see if they're being loaded
  if (visitReport?.photos) {
    visitReport.photos.forEach((photo, index) => {
      const photoUrl = getPhotoUrl(photo);
    });
  }

  const { data: timelineData } = useActivityTimeline({
    account_id: visitReport?.account_id,
    limit: 10,
  });
  const activities = timelineData?.data || [];
  const t = useTranslations("visitReportDetail");

  const formatDateTime = (dateString?: string | null) => {
    if (!dateString) return t("sections.notAvailable");
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return t("sections.invalidDate");
    return date.toLocaleString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const formatDate = (dateString?: string | null) => {
    if (!dateString) return t("sections.notAvailable");
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return t("sections.invalidDate");
    return date.toLocaleDateString("id-ID", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const getCurrentLocation = (): Promise<{ latitude: number; longitude: number; address: string }> => {
    return new Promise((resolve, reject) => {
      if (!navigator.geolocation) {
        reject(new Error(t("errors.geolocationUnsupported")));
        return;
      }

      navigator.geolocation.getCurrentPosition(
        async (position) => {
          // Detect Fake GPS
          const fakeGPSDetection = detectFakeGPSFromPosition(position, previousGPSPosition);
          
          if (fakeGPSDetection.isFakeGPS) {
            // Fake GPS detected - block check-out and show warning
            setFakeGPSReason(fakeGPSDetection.reason);
            setIsFakeGPSModalOpen(true);
            reject(new Error("Fake GPS detected"));
            return;
          }
          
          // GPS is valid
          const { latitude, longitude } = position.coords;
          setPreviousGPSPosition(position);
          
          // Try to get address from reverse geocoding (using a free service)
          let address = "";
          try {
            const response = await fetch(
              `https://nominatim.openstreetmap.org/reverse?format=json&lat=${latitude}&lon=${longitude}&zoom=18&addressdetails=1`
            );
            const data = await response.json();
            if (data.display_name) {
              address = data.display_name;
            } else {
              address = `${latitude.toFixed(6)}, ${longitude.toFixed(6)}`;
            }
          } catch (error) {
            // Fallback to coordinates if reverse geocoding fails
            address = `${latitude.toFixed(6)}, ${longitude.toFixed(6)}`;
          }

          resolve({ latitude, longitude, address });
        },
        (error) => {
          reject(error);
        },
        {
          enableHighAccuracy: true,
          timeout: 10000,
          maximumAge: 0,
        }
      );
    });
  };

  const handleCheckIn = async () => {
    if (!visitReportId || !visitReport) return;
    // Open camera dialog for selfie capture
    setIsCheckInCameraDialogOpen(true);
  };

  const handleCheckInWithPhoto = async (
    photo: File,
    deviceGPS: {
      latitude: number;
      longitude: number;
      accuracy?: number;
      timestamp: number;
    }
  ) => {
    if (!visitReportId) return;
    try {
      toast.loading(t("actions.processingCheckIn"), { id: "checkin-processing" });
      
      // Use GPS from camera dialog instead of getting location again
      // This avoids timeout issues and uses the GPS that was captured when photo was taken
      let location = {
        latitude: deviceGPS.latitude,
        longitude: deviceGPS.longitude,
        address: "",
      };
      
      // Try to get address from reverse geocoding (optional, don't block on failure)
      try {
        const response = await fetch(
          `https://nominatim.openstreetmap.org/reverse?format=json&lat=${deviceGPS.latitude}&lon=${deviceGPS.longitude}&zoom=18&addressdetails=1`
        );
        const data = await response.json();
        if (data.display_name) {
          location.address = data.display_name;
        } else {
          location.address = `${deviceGPS.latitude.toFixed(6)}, ${deviceGPS.longitude.toFixed(6)}`;
        }
      } catch (error) {
        // Fallback to coordinates if reverse geocoding fails
        location.address = `${deviceGPS.latitude.toFixed(6)}, ${deviceGPS.longitude.toFixed(6)}`;
      }
      
      // Call checkIn with photo and GPS metadata
      await checkIn.mutateAsync({
        id: visitReportId,
        data: {
          location,
        },
        options: {
          photo,
          deviceGPS: {
            latitude: deviceGPS.latitude,
            longitude: deviceGPS.longitude,
            accuracy: deviceGPS.accuracy,
            timestamp: deviceGPS.timestamp,
          },
        },
      });
      
      toast.dismiss("checkin-processing");
      toast.success(t("actions.checkInSuccess"));
      setIsCheckInCameraDialogOpen(false);
      
      // Refresh visit report data to show the uploaded photo
      // Refetch immediately and also call callback
      if (visitReportId) {
        // Refetch visit report data to get updated photos
        refetch().catch((err) => {
        });
        onVisitReportUpdated?.();
      }
    } catch (error) {
      toast.dismiss("checkin-processing");
      
      // Extract detailed error message
      let errorMessage = t("actions.checkInFailed");
      let errorDescription = "";
      
      // Log full error for debugging (with safe serialization)
      if (error && typeof error === "object") {
        const errorInfo: Record<string, unknown> = {
          type: typeof error,
          constructor: error.constructor?.name,
        };
        
        // Try to extract safe properties
        if ("message" in error) errorInfo.message = error.message;
        if ("code" in error) errorInfo.code = error.code;
        if ("name" in error) errorInfo.name = error.name;
        if ("response" in error) {
          const response = (error as { response?: unknown }).response;
          if (response && typeof response === "object") {
            const responseInfo: Record<string, unknown> = {};
            if ("status" in response) responseInfo.status = response.status;
            if ("statusText" in response) responseInfo.statusText = response.statusText;
            if ("data" in response) {
              const data = response.data;
              if (data && typeof data === "object") {
                const dataInfo: Record<string, unknown> = {};
                if ("error" in data && data.error && typeof data.error === "object") {
                  if ("code" in data.error) dataInfo.errorCode = data.error.code;
                  if ("message" in data.error) dataInfo.errorMessage = data.error.message;
                }
                if ("message" in data) dataInfo.message = data.message;
                responseInfo.data = dataInfo;
              }
            }
            errorInfo.response = responseInfo;
          }
        }
        
      } else {
      }
      
      // Handle GeolocationPositionError specifically (has code and message but not enumerable)
      if (error && typeof error === "object" && error.constructor?.name === "GeolocationPositionError") {
        const geoError = error as GeolocationPositionError;
      }
      
      // Check if error has originalError property (from mutation)
      if (error && typeof error === "object" && "originalError" in error) {
        const originalError = (error as { originalError?: unknown }).originalError;
      }
      
      // Handle GeolocationPositionError specifically (has code and message but not enumerable)
      if (error && typeof error === "object" && error.constructor?.name === "GeolocationPositionError") {
        const geoError = error as GeolocationPositionError;
        const errorMessages: Record<number, string> = {
          1: "GPS permission denied. Please allow location access.",
          2: "GPS position unavailable. Please check your location settings.",
          3: "GPS timeout. Please try again or check your location signal.",
        };
        errorDescription = geoError.message || errorMessages[geoError.code] || `GPS error (code: ${geoError.code})`;
        errorMessage = "GPS Location Error";
      }
      // Try to extract error message from various possible structures
      else if (error && typeof error === "object") {
        // Check if it's an Axios error with response
        if ("response" in error) {
          const axiosError = error as { 
            response?: { 
              data?: { 
                error?: { message?: string; code?: string; details?: Record<string, unknown> };
                message?: string;
              };
              status?: number;
              statusText?: string;
            };
            message?: string;
            code?: string;
          };
          
          // Extract from error.response.data.error.message
          if (axiosError.response?.data?.error?.message) {
            errorDescription = axiosError.response.data.error.message;
            if (axiosError.response.data.error.code === "INVALID_GPS") {
              errorMessage = "GPS validation failed";
            } else if (axiosError.response.data.error.code) {
              errorMessage = `Error: ${axiosError.response.data.error.code}`;
            }
          } 
          // Extract from error.response.data.message
          else if (axiosError.response?.data?.message) {
            errorDescription = axiosError.response.data.message;
          }
          // Extract from error.response.status
          else if (axiosError.response?.status) {
            errorDescription = `HTTP ${axiosError.response.status}: ${axiosError.response.statusText || "Request failed"}`;
          }
          // Fallback: stringify entire response data
          else if (axiosError.response?.data) {
            errorDescription = JSON.stringify(axiosError.response.data);
          }
          // Extract from error.message
          else if (axiosError.message) {
            errorDescription = axiosError.message;
          }
        }
        // Check if it's a plain Error object
        else if (error instanceof Error) {
          errorDescription = error.message || String(error);
        }
        // Handle GeolocationPositionError (might not be instanceof Error)
        else if ("code" in error && "message" in error && typeof error.code === "number") {
          const geoError = error as { code: number; message?: string };
          const errorMessages: Record<number, string> = {
            1: "GPS permission denied. Please allow location access.",
            2: "GPS position unavailable. Please check your location settings.",
            3: "GPS timeout. Please try again or check your location signal.",
          };
          errorDescription = geoError.message || errorMessages[geoError.code] || `GPS error (code: ${geoError.code})`;
          errorMessage = "GPS Location Error";
        }
        // Check if error has message property
        else if ("message" in error && typeof error.message === "string") {
          errorDescription = error.message;
        }
        // Fallback: try to stringify the error
        else {
          try {
            errorDescription = JSON.stringify(error);
          } catch {
            errorDescription = String(error);
          }
        }
      } else if (error) {
        errorDescription = String(error);
      }
      
      // If still no description, try to extract from error object directly
      if (!errorDescription) {
        // Try to get error message from error object properties
        if (error && typeof error === "object") {
          // Try common error message properties
          const possibleMessageKeys = ["message", "error", "err", "detail", "details", "reason"];
          for (const key of possibleMessageKeys) {
            if (key in error && typeof (error as Record<string, unknown>)[key] === "string") {
              errorDescription = String((error as Record<string, unknown>)[key]);
              break;
            }
          }
          
          // If still no message, try to stringify the error
          if (!errorDescription) {
            try {
              const errorStr = JSON.stringify(error, null, 2);
              if (errorStr && errorStr !== "{}") {
                errorDescription = `Error details: ${errorStr}`;
              }
            } catch {
              // Ignore JSON stringify errors
            }
          }
        }
        
        // Final fallback
        if (!errorDescription) {
          errorDescription = "An unknown error occurred. Please check the browser console (F12) for details.";
        }
      }
      
      toast.error(errorMessage, { 
        description: errorDescription,
        duration: 5000,
      });
    }
  };

  const handleCheckOut = async () => {
    if (!visitReportId || !visitReport) return;
    try {
      toast.loading(t("actions.gettingLocation"), { id: "checkout-location" });
      const location = await getCurrentLocation();
      toast.dismiss("checkout-location");
      
      await checkOut.mutateAsync({
        id: visitReportId,
        data: {
          location,
        },
      });
      toast.success(t("actions.checkOutSuccess"));
      onVisitReportUpdated?.();
    } catch (error) {
      toast.dismiss("checkout-location");
      if (error instanceof Error) {
        toast.error(t("actions.checkOutGetLocationFailed"), { description: error.message });
      } else {
        toast.error(t("actions.checkOutFailed"));
      }
    }
  };

  const handleApprove = async () => {
    if (!visitReportId) return;
    try {
      await approve.mutateAsync(visitReportId);
      toast.success(t("actions.approveSuccess"));
      onVisitReportUpdated?.();
    } catch (error) {
      // Error already handled
    }
  };

  const handleReject = async () => {
    if (!visitReportId || !rejectReason.trim()) return;
    try {
      await reject.mutateAsync({
        id: visitReportId,
        data: { reason: rejectReason },
      });
      toast.success(t("actions.rejectSuccess"));
      setIsRejectDialogOpen(false);
      setRejectReason("");
      onVisitReportUpdated?.();
    } catch (error) {
      // Error already handled
    }
  };

  const handleUploadPhoto = async (file: File) => {
    if (!visitReportId) return;
    try {
      await uploadPhoto.mutateAsync({
        id: visitReportId,
        file: file,
      });
      toast.success(t("actions.photoUploadSuccess"));
      refetch();
      onVisitReportUpdated?.();
    } catch (error) {
      toast.error(t("actions.photoUploadFailed"));
    }
  };


  return (
    <>
      <Drawer
        open={open}
        onOpenChange={onOpenChange}
        title={t("drawerTitle")}
        side="right"
        className="max-w-3xl"
      >

          {isLoading && (
            <div className="space-y-6">
              <Skeleton className="h-8 w-48" />
              <Card>
                <CardHeader>
                  <Skeleton className="h-6 w-32" />
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-4 w-3/4" />
                  </div>
                </CardContent>
              </Card>
            </div>
          )}

          {error && (
            <div className="text-center text-muted-foreground py-8">
              {t("loadError")}
            </div>
          )}

          {!isLoading && !error && visitReport && (
            <div className="space-y-6">
              {/* Header */}
              <div className="flex items-center justify-between pb-4 border-b">
                <div className="flex items-center gap-3">
                  <Badge variant={statusColors[visitReport.status] || "outline"}>
                    {visitReport.status}
                  </Badge>
                  <span className="text-sm text-muted-foreground">
                    {formatDate(visitReport.visit_date)}
                  </span>
                </div>
                <div className="flex gap-2">
                  <VisitReportInsightsButton visitReportId={visitReport.id} iconOnly />
                  {visitReport.status === "submitted" && (
                    <>
                      <Button
                        size="icon"
                        variant="outline"
                        onClick={() => setIsRejectDialogOpen(true)}
                        disabled={reject.isPending}
                        title={t("actions.reject")}
                      >
                        <XCircle className="h-4 w-4" />
                      </Button>
                      <Button
                        size="icon"
                        onClick={handleApprove}
                        disabled={approve.isPending}
                        title={t("actions.approve")}
                      >
                        <CheckCircle2 className="h-4 w-4" />
                      </Button>
                    </>
                  )}
                  {visitReport.status === "draft" && !visitReport.check_in_time && (
                    <Button
                      size="icon"
                      onClick={handleCheckIn}
                      disabled={checkIn.isPending}
                      title={t("actions.checkIn")}
                    >
                      <MapPin className="h-4 w-4" />
                    </Button>
                  )}
                  {visitReport.check_in_time && !visitReport.check_out_time && (
                    <Button
                      size="icon"
                      variant="outline"
                      onClick={handleCheckOut}
                      disabled={checkOut.isPending}
                      title={t("actions.checkOut")}
                    >
                      <MapPin className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              </div>

              {/* Basic Information */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <FileText className="h-5 w-5" />
                    {t("sections.visitInformationTitle")}
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <div className="text-sm text-muted-foreground mb-1">
                        {t("sections.accountLabel")}
                      </div>
                      <div className="flex items-center gap-2">
                        <Building2 className="h-4 w-4 text-muted-foreground" />
                        <span className="font-medium">{visitReport.account?.name || "N/A"}</span>
                      </div>
                    </div>
                    {visitReport.contact && (
                      <div>
                        <div className="text-sm text-muted-foreground mb-1">
                          {t("sections.contactLabel")}
                        </div>
                        <div className="flex items-center gap-2">
                          <User className="h-4 w-4 text-muted-foreground" />
                          <span className="font-medium">{visitReport.contact?.name ?? t("sections.notAvailable")}</span>
                        </div>
                      </div>
                    )}
                    {visitReport.deal && (
                      <div>
                        <div className="text-sm text-muted-foreground mb-1">
                          {t("sections.dealLabel")}
                        </div>
                        <div className="flex items-center gap-2">
                          <FileText className="h-4 w-4 text-muted-foreground" />
                          <span className="font-medium">{visitReport.deal?.title ?? t("sections.notAvailable")}</span>
                        </div>
                      </div>
                    )}
                  </div>

                  <div>
                    <div className="text-sm text-muted-foreground mb-1">
                      {t("sections.purposeLabel")}
                    </div>
                    <p className="text-sm">{visitReport.purpose}</p>
                  </div>

                  {visitReport.notes && (
                    <div>
                      <div className="text-sm text-muted-foreground mb-1">
                        {t("sections.notesLabel")}
                      </div>
                      <p className="text-sm whitespace-pre-wrap">{visitReport.notes}</p>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Check In/Out Information */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Clock className="h-5 w-5" />
                    {t("sections.checkInOutTitle")}
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <div className="text-sm text-muted-foreground mb-1">
                        {t("sections.checkInLabel")}
                      </div>
                      <div className="flex items-center gap-2">
                        {visitReport.check_in_time ? (
                          <>
                            <Calendar className="h-4 w-4 text-muted-foreground" />
                            <span className="text-sm">{formatDateTime(visitReport.check_in_time)}</span>
                          </>
                        ) : (
                          <span className="text-sm text-muted-foreground">
                            {t("sections.notCheckedIn")}
                          </span>
                        )}
                      </div>
                      {visitReport.check_in_location && (
                        <div className="text-xs text-muted-foreground mt-1">
                          <MapPin className="h-3 w-3 inline mr-1" />
                          {visitReport.check_in_location?.address || 
                            (visitReport.check_in_location?.latitude && visitReport.check_in_location?.longitude
                              ? `${visitReport.check_in_location.latitude}, ${visitReport.check_in_location.longitude}`
                              : t("sections.locationNotAvailable"))}
                        </div>
                      )}
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground mb-1">
                        {t("sections.checkOutLabel")}
                      </div>
                      <div className="flex items-center gap-2">
                        {visitReport.check_out_time ? (
                          <>
                            <Calendar className="h-4 w-4 text-muted-foreground" />
                            <span className="text-sm">{formatDateTime(visitReport.check_out_time)}</span>
                          </>
                        ) : (
                          <span className="text-sm text-muted-foreground">
                            {t("sections.notCheckedOut")}
                          </span>
                        )}
                      </div>
                      {visitReport.check_out_location && (
                        <div className="text-xs text-muted-foreground mt-1">
                          <MapPin className="h-3 w-3 inline mr-1" />
                          {visitReport.check_out_location?.address || 
                            (visitReport.check_out_location?.latitude && visitReport.check_out_location?.longitude
                              ? `${visitReport.check_out_location.latitude}, ${visitReport.check_out_location.longitude}`
                              : t("sections.locationNotAvailable"))}
                        </div>
                      )}
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Photos */}
              <Card>
                <CardHeader className="flex flex-row items-center justify-between">
                  <CardTitle>{t("sections.photosTitle")}</CardTitle>
                  {visitReport.status !== "approved" && (
                    <Button
                      size="sm"
                      onClick={() => setIsPhotoUploadDialogOpen(true)}
                    >
                      <Plus className="h-4 w-4 mr-2" />
                      {t("actions.addPhoto")}
                    </Button>
                  )}
                </CardHeader>
                <CardContent>
                  {visitReport.photos && Array.isArray(visitReport.photos) && visitReport.photos.length > 0 ? (
                    <div className="grid grid-cols-3 gap-4">
                      {visitReport.photos.map((photo, index) => {
                        const photoUrl = getPhotoUrl(photo);
                        return (
                        <div key={photo || `photo-${index}`} className="relative group">
                          <img
                            src={photoUrl}
                            alt={`Visit documentation ${index + 1}`}
                            className="w-full h-32 object-cover rounded-md border cursor-pointer"
                            onError={(e) => {
                              // Handle image load error
                              const target = e.target as HTMLImageElement;
                              target.style.display = "none";
                              const parent = target.parentElement;
                              if (parent) {
                                parent.innerHTML = `
                                  <div class="w-full h-32 flex items-center justify-center bg-muted rounded-md border border-dashed">
                                    <span class="text-xs text-muted-foreground">Failed to load image</span>
                                  </div>
                                `;
                              }
                            }}
                            onLoad={(e) => {
                              // Ensure image is visible when loaded
                              const target = e.target as HTMLImageElement;
                              target.style.display = "block";
                            }}
                          />
                          <a
                            href={photoUrl}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="absolute inset-0 flex items-center justify-center bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity rounded-md cursor-pointer"
                            onClick={(e) => {
                              // Open image in new tab
                              e.preventDefault();
                              window.open(photoUrl, "_blank");
                            }}
                          >
                            <span className="text-white text-xs">
                              {t("sections.viewPhoto")}
                            </span>
                          </a>
                        </div>
                        );
                      })}
                    </div>
                  ) : (
                    <div className="text-center text-muted-foreground py-8">
                      {t("sections.noPhotos")}
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Approval Information */}
              {visitReport.approved_at && (
                <Card>
                  <CardHeader>
                    <CardTitle>{t("sections.approvalTitle")}</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-sm">
                      <div className="text-muted-foreground mb-1">
                        {t("sections.approvedAtLabel")}
                      </div>
                      <div>{formatDateTime(visitReport.approved_at)}</div>
                    </div>
                  </CardContent>
                </Card>
              )}

              {visitReport.rejection_reason && (
                <Card>
                  <CardHeader>
                    <CardTitle>{t("sections.rejectionTitle")}</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-sm">
                      <div className="text-muted-foreground mb-1">
                        {t("sections.rejectionReasonLabel")}
                      </div>
                      <div>{visitReport.rejection_reason}</div>
                    </div>
                  </CardContent>
                </Card>
              )}

              {/* Activity Timeline */}
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <CardTitle>{t("sections.activityTimelineTitle")}</CardTitle>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setIsCreateActivityDialogOpen(true)}
                    >
                      <Plus className="h-4 w-4 mr-2" />
                      {t("sections.addActivity")}
                    </Button>
                  </div>
                </CardHeader>
                <CardContent>
                  <ActivityTimeline
                    activities={activities}
                    isLoading={!timelineData}
                    accountId={visitReport.account_id}
                  />
                </CardContent>
              </Card>
            </div>
          )}
      </Drawer>

      {/* Reject Dialog */}
      <Dialog open={isRejectDialogOpen} onOpenChange={setIsRejectDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("rejectDialog.title")}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <label className="text-sm font-medium mb-2 block">
                {t("rejectDialog.reasonLabel")} *
              </label>
              <textarea
                value={rejectReason}
                onChange={(e) => setRejectReason(e.target.value)}
                placeholder={t("rejectDialog.reasonPlaceholder")}
                className="w-full min-h-[100px] rounded-md border border-input bg-background px-3 py-2 text-sm"
                rows={4}
              />
            </div>
            <div className="flex justify-end gap-2">
              <Button
                variant="outline"
                onClick={() => {
                  setIsRejectDialogOpen(false);
                  setRejectReason("");
                }}
                disabled={reject.isPending}
              >
                {t("rejectDialog.cancel")}
              </Button>
              <Button
                onClick={handleReject}
                disabled={reject.isPending || !rejectReason.trim()}
                variant="destructive"
              >
                {reject.isPending ? t("rejectDialog.submitting") : t("rejectDialog.submit")}
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>


      {/* Create Activity Dialog */}
      {visitReport && (
        <CreateActivityDialog
          open={isCreateActivityDialogOpen}
          onOpenChange={setIsCreateActivityDialogOpen}
          accountId={visitReport.account_id}
          contactId={visitReport.contact_id || undefined}
          dealId={visitReport.deal_id || undefined}
          leadId={visitReport.lead_id || undefined}
          onSuccess={() => {
            // Refresh timeline - query will auto-refresh due to invalidation in hook
            onVisitReportUpdated?.();
          }}
        />
      )}

      {/* Photo Upload Dialog */}
      <PhotoUploadDialog
        open={isPhotoUploadDialogOpen}
        onOpenChange={setIsPhotoUploadDialogOpen}
        onUpload={handleUploadPhoto}
        isLoading={uploadPhoto.isPending}
      />

      {/* Check-In Camera Dialog */}
      <CheckInCameraDialog
        open={isCheckInCameraDialogOpen}
        onOpenChange={setIsCheckInCameraDialogOpen}
        onCapture={handleCheckInWithPhoto}
        isLoading={checkIn.isPending}
      />
      
      {/* Fake GPS Warning Modal */}
      <FakeGPSWarningModal
        open={isFakeGPSModalOpen}
        onOpenChange={setIsFakeGPSModalOpen}
        reason={fakeGPSReason}
      />
    </>
  );
}

