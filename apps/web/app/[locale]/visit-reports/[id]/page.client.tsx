/* eslint-disable @next/next/no-img-element */
"use client";

import type React from "react";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import {
  Building2,
  Calendar,
  Clock,
  FileText,
  MapPin,
  User,
} from "lucide-react";

import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import { PageDetailLayout } from "@/components/layouts/page-detail-layout";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useVisitReport } from "@/features/sales-crm/visit-report/hooks/useVisitReports";
import { getVisitReportPhotoUrl } from "@/features/sales-crm/visit-report/utils/photo-url";

function VisitReportDetailPageContent() {
  const params = useParams();
  const visitReportId = params.id as string;
  const t = useTranslations("visitReportDetail");
  const tStatus = useTranslations("visitReportTeamOverview.status");

  const { data, isLoading, error } = useVisitReport(visitReportId);
  const visitReport = data?.data;

  const formatDate = (dateString?: string | null) => {
    if (!dateString) return t("sections.notAvailable");
    const date = new Date(dateString);
    if (Number.isNaN(date.getTime())) return t("sections.invalidDate");
    return date.toLocaleDateString("id-ID", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const formatDateTime = (dateString?: string | null) => {
    if (!dateString) return t("sections.notAvailable");
    const date = new Date(dateString);
    if (Number.isNaN(date.getTime())) return t("sections.invalidDate");
    return date.toLocaleString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const badgeVariant = (status?: string) => {
    switch (status) {
      case "approved":
      case "completed":
        return "default" as const;
      case "submitted":
        return "secondary" as const;
      case "rejected":
        return "destructive" as const;
      default:
        return "outline" as const;
    }
  };

  if (isLoading) {
    return (
      <PageMotion className="p-2 sm:p-4">
        <div className="space-y-4">
          <Skeleton className="h-10 w-72" />
          <Skeleton className="h-40 w-full" />
          <Skeleton className="h-48 w-full" />
          <Skeleton className="h-48 w-full" />
        </div>
      </PageMotion>
    );
  }

  if (error || !visitReport) {
    return (
      <PageMotion className="p-2 sm:p-4">
        <Card className="border-border/70 bg-card/80">
          <CardContent className="py-10 text-center text-muted-foreground">
            {t("page.loadError")}
          </CardContent>
        </Card>
      </PageMotion>
    );
  }

  return (
    <PageMotion className="p-2 sm:p-4">
      <PageDetailLayout
        title={visitReport.account?.name || visitReport.purpose || t("page.fallbackTitle")}
        subtitle={
          <div className="mt-1 flex flex-wrap items-center gap-2 text-sm">
            <Badge variant={badgeVariant(visitReport.status)}>
              {tStatus(visitReport.status)}
            </Badge>
            <span className="text-muted-foreground">
              {t("header.visitDate")}: {formatDate(visitReport.visit_date)}
            </span>
            {visitReport.sales_rep?.name && (
              <span className="text-muted-foreground">
                {t("header.salesRep")}: {visitReport.sales_rep.name}
              </span>
            )}
          </div>
        }
        backHref={t("page.backHref")}
      >
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <FileText className="h-5 w-5" />
                {t("sections.visitInformationTitle")}
              </CardTitle>
              <CardDescription>{t("page.informationDescription")}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 md:grid-cols-2">
                <InfoRow
                  icon={<Building2 className="h-4 w-4 text-muted-foreground" />}
                  label={t("sections.accountLabel")}
                  value={visitReport.account?.name || t("sections.notAvailable")}
                />
                <InfoRow
                  icon={<User className="h-4 w-4 text-muted-foreground" />}
                  label={t("sections.contactLabel")}
                  value={visitReport.contact?.name || t("sections.notAvailable")}
                />
                <InfoRow
                  icon={<FileText className="h-4 w-4 text-muted-foreground" />}
                  label={t("sections.dealLabel")}
                  value={visitReport.deal?.title || t("sections.notAvailable")}
                />
                <InfoRow
                  icon={<User className="h-4 w-4 text-muted-foreground" />}
                  label={t("header.salesRep")}
                  value={visitReport.sales_rep?.name || t("sections.notAvailable")}
                />
              </div>

              <TextBlock
                label={t("sections.purposeLabel")}
                value={visitReport.purpose || t("sections.notAvailable")}
              />
              <TextBlock
                label={t("sections.notesLabel")}
                value={visitReport.notes || t("sections.notAvailable")}
              />
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Clock className="h-5 w-5" />
                {t("sections.checkInOutTitle")}
              </CardTitle>
              <CardDescription>{t("page.trackingDescription")}</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4 md:grid-cols-2">
              <LocationBlock
                label={t("sections.checkInLabel")}
                timestamp={visitReport.check_in_time}
                location={visitReport.check_in_location?.address ||
                  ((visitReport.check_in_location?.latitude && visitReport.check_in_location?.longitude)
                    ? `${visitReport.check_in_location.latitude}, ${visitReport.check_in_location.longitude}`
                    : undefined)}
                emptyLabel={t("sections.notCheckedIn")}
                formatDateTime={formatDateTime}
                t={t}
              />
              <LocationBlock
                label={t("sections.checkOutLabel")}
                timestamp={visitReport.check_out_time}
                location={visitReport.check_out_location?.address ||
                  ((visitReport.check_out_location?.latitude && visitReport.check_out_location?.longitude)
                    ? `${visitReport.check_out_location.latitude}, ${visitReport.check_out_location.longitude}`
                    : undefined)}
                emptyLabel={t("sections.notCheckedOut")}
                formatDateTime={formatDateTime}
                t={t}
              />
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Calendar className="h-5 w-5" />
                {t("header.statusBadge")}
              </CardTitle>
              <CardDescription>{t("page.systemDescription")}</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4 md:grid-cols-2">
              <InfoRow
                icon={<FileText className="h-4 w-4 text-muted-foreground" />}
                label={t("header.statusBadge")}
                value={tStatus(visitReport.status)}
              />
              <InfoRow
                icon={<Calendar className="h-4 w-4 text-muted-foreground" />}
                label={t("header.visitDate")}
                value={formatDate(visitReport.visit_date)}
              />
              <InfoRow
                icon={<Calendar className="h-4 w-4 text-muted-foreground" />}
                label={t("header.createdAt")}
                value={formatDateTime(visitReport.created_at)}
              />
              <InfoRow
                icon={<Calendar className="h-4 w-4 text-muted-foreground" />}
                label={t("header.updatedAt")}
                value={formatDateTime(visitReport.updated_at)}
              />
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>{t("sections.photosTitle")}</CardTitle>
            </CardHeader>
            <CardContent>
              {Array.isArray(visitReport.photos) && visitReport.photos.length > 0 ? (
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                  {visitReport.photos.map((photo, index) => {
                    const photoUrl = getVisitReportPhotoUrl(photo);

                    return (
                      <a
                        key={photo || `photo-${index}`}
                        href={photoUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="group overflow-hidden rounded-lg border border-border/70"
                      >
                        <img
                          src={photoUrl}
                          alt={`${t("sections.photosTitle")} ${index + 1}`}
                          className="h-44 w-full object-cover transition-transform group-hover:scale-[1.02]"
                        />
                      </a>
                    );
                  })}
                </div>
              ) : (
                <div className="py-8 text-center text-muted-foreground">
                  {t("sections.noPhotos")}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </PageDetailLayout>
    </PageMotion>
  );
}

function InfoRow({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode;
  label: string;
  value: string;
}) {
  return (
    <div className="space-y-1.5">
      <div className="text-sm text-muted-foreground">{label}</div>
      <div className="flex items-center gap-2 text-sm font-medium">
        {icon}
        <span>{value}</span>
      </div>
    </div>
  );
}

function TextBlock({
  label,
  value,
}: {
  label: string;
  value: string;
}) {
  return (
    <div className="space-y-1.5">
      <div className="text-sm text-muted-foreground">{label}</div>
      <p className="text-sm whitespace-pre-wrap">{value}</p>
    </div>
  );
}

function LocationBlock({
  label,
  timestamp,
  location,
  emptyLabel,
  formatDateTime,
  t,
}: {
  label: string;
  timestamp?: string;
  location?: string;
  emptyLabel: string;
  formatDateTime: (value?: string | null) => string;
  t: ReturnType<typeof useTranslations<"visitReportDetail">>;
}) {
  return (
    <div className="space-y-2 rounded-lg border border-border/70 p-4">
      <div className="text-sm text-muted-foreground">{label}</div>
      {timestamp ? (
        <>
          <div className="flex items-center gap-2 text-sm font-medium">
            <Calendar className="h-4 w-4 text-muted-foreground" />
            <span>{formatDateTime(timestamp)}</span>
          </div>
          <div className="flex items-start gap-2 text-sm text-muted-foreground">
            <MapPin className="mt-0.5 h-4 w-4 shrink-0" />
            <span>{location || t("sections.locationNotAvailable")}</span>
          </div>
        </>
      ) : (
        <div className="text-sm text-muted-foreground">{emptyLabel}</div>
      )}
    </div>
  );
}

export default function VisitReportDetailPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="visit-reports.view">
        <VisitReportDetailPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
