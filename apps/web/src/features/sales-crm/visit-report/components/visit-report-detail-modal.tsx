"use client";

import { Calendar, MapPin, Clock, User, Building2, FileText, SquarePen } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Drawer } from "@/components/ui/drawer";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import {
	useVisitReport,
	useCheckIn,
	useCheckOut,
	useActivityTimeline,
	useUploadPhoto,
  useUpdateVisitReport,
} from "../hooks/useVisitReports";
import { toast } from "sonner";
import { useMemo, useState } from "react";
import { ProductInterestTab } from "./product-interest-tab";
import { PhotoUploadDialog } from "./photo-upload-dialog";
import { VisitReportInsightsButton } from "@/features/ai/components/visit-report-insights-button";
import { CheckInCameraDialog } from "./check-in-camera-dialog";
import { FakeGPSWarningModal } from "./fake-gps-warning-modal";
import { detectFakeGPSFromPosition } from "../utils/detectFakeGPS";
import { getVisitReportPhotoUrl } from "../utils/photo-url";
import { useTranslations } from "next-intl";
import { VisitReportForm } from "./visit-report-form";
import type { Activity } from "../types/activity";
import type { VisitReport } from "../types";
import { SubmitVisitReportModal } from "./submit-visit-report-modal";

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
  const uploadPhoto = useUploadPhoto();
  const updateVisitReport = useUpdateVisitReport();
  const [isEditVisitDialogOpen, setIsEditVisitDialogOpen] = useState(false);
  const [isPhotoUploadDialogOpen, setIsPhotoUploadDialogOpen] = useState(false);
  const [isCheckInCameraDialogOpen, setIsCheckInCameraDialogOpen] = useState(false);
  const [isFakeGPSModalOpen, setIsFakeGPSModalOpen] = useState(false);
  const [isSubmitVisitDialogOpen, setIsSubmitVisitDialogOpen] = useState(false);
  const [fakeGPSReason, setFakeGPSReason] = useState<string | undefined>();
  const [previousGPSPosition, setPreviousGPSPosition] = useState<GeolocationPosition | undefined>();

  const visitReport = data?.data;
  
  // Debug: Log photos to see if they're being loaded
  if (visitReport?.photos) {
    visitReport.photos.forEach((photo, index) => {
      const photoUrl = getVisitReportPhotoUrl(photo);
    });
  }

  const { data: timelineData, refetch: refetchTimeline } = useActivityTimeline({
    account_id: visitReport?.account_id,
    limit: 10,
  });
  const activities = timelineData?.data || [];
  const productInterestActivities = useMemo(
    () => buildProductInterestActivities(visitReport, activities),
    [visitReport, activities],
  );
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

  const canMarkCompleted = Boolean(
    visitReport &&
    visitReport.check_in_time &&
    visitReport.status !== "completed"
  );

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

  const handleVisitUpdate = async (formData: {
    account_id?: string;
    contact_id?: string;
    deal_id?: string;
    lead_id?: string;
    visit_date?: string;
    purpose?: string;
    notes?: string;
    metadata?: Record<string, unknown>;
  }) => {
    if (!visitReportId) return;

    await updateVisitReport.mutateAsync({
      id: visitReportId,
      data: formData,
    });
    toast.success(t("actions.visitUpdateSuccess"));
    setIsEditVisitDialogOpen(false);
    refetch();
    refetchTimeline();
    onVisitReportUpdated?.();
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
                  <span className="text-sm text-muted-foreground">
                    {formatDate(visitReport.visit_date)}
                  </span>
                </div>
                <div className="flex gap-2">
                  {!visitReport.check_in_time && (
                    <Button variant="outline" size="sm" onClick={handleCheckIn}>
                      {t("actions.checkIn")}
                    </Button>
                  )}
                  {visitReport.check_in_time && !visitReport.check_out_time && (
                    <Button variant="outline" size="sm" onClick={handleCheckOut}>
                      {t("actions.checkOut")}
                    </Button>
                  )}
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setIsPhotoUploadDialogOpen(true)}
                  >
                    {t("actions.addPhoto")}
                  </Button>
                  {canMarkCompleted && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setIsSubmitVisitDialogOpen(true)}
                    >
                      {t("actions.markComplete")}
                    </Button>
                  )}
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setIsEditVisitDialogOpen(true)}
                    className="gap-2"
                  >
                    <SquarePen className="h-4 w-4" />
                    {t("actions.editVisitLog")}
                  </Button>
                  <VisitReportInsightsButton visitReportId={visitReport.id} iconOnly />
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
                <CardHeader>
                  <CardTitle>{t("sections.photosTitle")}</CardTitle>
                </CardHeader>
                <CardContent>
                  {visitReport.photos && Array.isArray(visitReport.photos) && visitReport.photos.length > 0 ? (
                    <div className="grid grid-cols-3 gap-4">
                      {visitReport.photos.map((photo, index) => {
                        const photoUrl = getVisitReportPhotoUrl(photo);
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

              <Card>
                <CardHeader>
                  <CardTitle>{t("sections.productInterestTitle")}</CardTitle>
                </CardHeader>
                <CardContent>
                  <ProductInterestTab
                    activities={productInterestActivities}
                    isLoading={!timelineData}
                  />
                </CardContent>
              </Card>
            </div>
          )}
      </Drawer>

      {visitReport && (
        <Dialog open={isEditVisitDialogOpen} onOpenChange={setIsEditVisitDialogOpen}>
          <DialogContent className="max-w-3xl">
            <DialogHeader>
              <DialogTitle>{t("dialogs.editVisitTitle")}</DialogTitle>
            </DialogHeader>
            <VisitReportForm
              visitReport={visitReport}
              onSubmit={handleVisitUpdate}
              onCancel={() => setIsEditVisitDialogOpen(false)}
              isLoading={updateVisitReport.isPending}
              open={isEditVisitDialogOpen}
            />
          </DialogContent>
        </Dialog>
      )}

      <PhotoUploadDialog
        open={isPhotoUploadDialogOpen}
        onOpenChange={setIsPhotoUploadDialogOpen}
        onUpload={handleUploadPhoto}
        isLoading={uploadPhoto.isPending}
      />

      <CheckInCameraDialog
        open={isCheckInCameraDialogOpen}
        onOpenChange={setIsCheckInCameraDialogOpen}
        onCheckIn={handleCheckInWithPhoto}
        isLoading={checkIn.isPending}
      />

      <FakeGPSWarningModal
        open={isFakeGPSModalOpen}
        onOpenChange={setIsFakeGPSModalOpen}
        reason={fakeGPSReason}
      />

      {visitReportId && (
        <SubmitVisitReportModal
          visitReportId={visitReportId}
          isOpen={isSubmitVisitDialogOpen}
          onClose={() => setIsSubmitVisitDialogOpen(false)}
          onSuccess={() => {
            setIsSubmitVisitDialogOpen(false);
            refetch();
            onVisitReportUpdated?.();
          }}
        />
      )}

    </>
  );
}

function buildProductInterestActivities(
  visitReport: VisitReport | undefined,
  activities: Activity[],
): Activity[] {
  if (!visitReport) return activities;

  const visitProductInterests = Array.isArray(visitReport.metadata?.product_interests)
    ? visitReport.metadata.product_interests
    : [];

  const relatedActivities = activities.filter((activity) => {
    if (!activity.metadata || typeof activity.metadata !== "object") return true;
    const metadata = activity.metadata as Record<string, unknown>;
    return metadata.visit_report_id !== visitReport.id;
  });

  if (visitProductInterests.length === 0) {
    return relatedActivities;
  }

  const visitActivity: Activity = {
    id: `visit-report-${visitReport.id}`,
    type: "visit",
    account_id: visitReport.account_id,
    contact_id: visitReport.contact_id,
    deal_id: visitReport.deal_id,
    lead_id: visitReport.lead_id,
    user_id: visitReport.sales_rep_id,
    description: visitReport.purpose,
    timestamp: visitReport.visit_date,
    metadata: {
      ...(visitReport.metadata ?? {}),
      visit_report_id: visitReport.id,
      visit_date: visitReport.visit_date,
    },
    created_at: visitReport.created_at,
    updated_at: visitReport.updated_at,
    account: visitReport.account,
    contact: visitReport.contact,
  };

  return [visitActivity, ...relatedActivities];
}
