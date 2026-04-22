"use client";

import { Activity, Clock, User, Building2, Contact, FileText } from "lucide-react";
import { DataTable, type Column } from "@/components/ui/data-table";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import type { Activity as ActivityType } from "../types/activity";
import { useTranslations } from "next-intl";
import { renderIcon } from "../lib/icon-utils";

interface ActivityTimelineProps {
  readonly activities: ActivityType[];
  readonly isLoading: boolean;
  readonly accountId?: string;
}

export function ActivityTimeline({ activities, isLoading, accountId }: ActivityTimelineProps) {
  const t = useTranslations("visitReportActivityTimeline");

  const columns: Column<ActivityType>[] = [
    {
      id: "type",
      header: t("table.type"),
      accessor: (row) => {
        // Use ActivityType if available, otherwise fallback to type string
        const activityType = row.activity_type;
        const iconName = activityType?.icon;
        const badgeColor = activityType?.badge_color as "default" | "secondary" | "destructive" | "outline" | undefined;
        const typeName = activityType?.name ?? row.type;
        
        return (
          <div className="flex items-center gap-2">
            <div className="text-muted-foreground">
              {iconName ? renderIcon(iconName, "h-4 w-4") : <Activity className="h-4 w-4" />}
            </div>
            <Badge variant={badgeColor ?? "outline"} className="text-xs">
              {typeName}
            </Badge>
          </div>
        );
      },
      className: "w-[120px]",
    },
    {
      id: "description",
      header: t("table.description"),
      accessor: (row) => {
        // Format metadata nicely instead of raw JSON
        const formatMetadata = (metadata: Record<string, unknown> | undefined): string => {
          if (!metadata || Object.keys(metadata).length === 0) return "";
          
          const parts: string[] = [];
          if (metadata.status && typeof metadata.status === "string") {
            parts.push(`${t("meta.statusLabel")}: ${metadata.status}`);
          }
          if (metadata.action && typeof metadata.action === "string") {
            parts.push(`${t("meta.actionLabel")}: ${metadata.action}`);
          }
          if (metadata.visit_date && typeof metadata.visit_date === "string") {
            parts.push(`${t("meta.visitDateLabel")}: ${metadata.visit_date}`);
          }
          if (metadata.duration && typeof metadata.duration === "string") {
            parts.push(`${t("meta.durationLabel")}: ${metadata.duration}`);
          }
          if (metadata.outcome && typeof metadata.outcome === "string") {
            parts.push(`${t("meta.outcomeLabel")}: ${metadata.outcome}`);
          }
          if (metadata.subject && typeof metadata.subject === "string") {
            parts.push(`${t("meta.subjectLabel")}: ${metadata.subject}`);
          }
          if (metadata.priority && typeof metadata.priority === "string") {
            parts.push(`${t("meta.priorityLabel")}: ${metadata.priority}`);
          }
          if (metadata.value) {
            const valueStr = typeof metadata.value === "number" 
              ? metadata.value.toLocaleString("id-ID")
              : String(metadata.value);
            parts.push(`${t("meta.valueLabel")}: ${valueStr}`);
          }
          
          return parts.length > 0 ? parts.join(" • ") : "";
        };

        const metadataText = formatMetadata(row.metadata as Record<string, unknown> | undefined);
        
        return (
          <div className="flex-1">
            <p className="text-sm font-medium line-clamp-2">{row.description}</p>
            {metadataText && (
              <p className="text-xs text-muted-foreground mt-1 line-clamp-1">
                {metadataText}
              </p>
            )}
          </div>
        );
      },
      className: "min-w-[300px]",
    },
    {
      id: "account",
      header: t("table.account"),
      accessor: (row) => (
        row.account ? (
          <div className="flex items-center gap-2">
            <Building2 className="h-3 w-3 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              {typeof row.account === "object" && "name" in row.account
                ? row.account.name
                : "N/A"}
            </span>
          </div>
        ) : (
          <span className="text-sm text-muted-foreground">-</span>
        )
      ),
      className: "w-[150px]",
    },
    {
      id: "contact",
      header: t("table.contact"),
      accessor: (row) => (
        row.contact ? (
          <div className="flex items-center gap-2">
            <Contact className="h-3 w-3 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              {typeof row.contact === "object" && "name" in row.contact
                ? row.contact.name
                : "N/A"}
            </span>
          </div>
        ) : (
          <span className="text-sm text-muted-foreground">-</span>
        )
      ),
      className: "w-[150px]",
    },
    {
      id: "deal",
      header: t("table.deal"),
      accessor: (row) => (
        row.deal ? (
          <div className="flex items-center gap-2">
            <FileText className="h-3 w-3 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              {typeof row.deal === "object" && "title" in row.deal
                ? row.deal.title
                : "N/A"}
            </span>
          </div>
        ) : (
          <span className="text-sm text-muted-foreground">-</span>
        )
      ),
      className: "w-[150px]",
    },
    {
      id: "user",
      header: t("table.user"),
      accessor: (row) => (
        row.user ? (
          <div className="flex items-center gap-2">
            <User className="h-3 w-3 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              {typeof row.user === "object" && "name" in row.user
                ? row.user.name
                : "N/A"}
            </span>
          </div>
        ) : (
          <span className="text-sm text-muted-foreground">-</span>
        )
      ),
      className: "w-[150px]",
    },
    {
      id: "timestamp",
      header: t("table.timestamp"),
      accessor: (row) => {
        const date = new Date(row.timestamp);
        const formatted = date.toLocaleString("id-ID", {
          year: "numeric",
          month: "short",
          day: "numeric",
          hour: "2-digit",
          minute: "2-digit",
        });
        return (
          <div className="flex items-center gap-2">
            <Clock className="h-3 w-3 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">{formatted}</span>
          </div>
        );
      },
      className: "w-[180px]",
    },
  ];

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 5 }, (_, i) => (
          <Skeleton key={`skeleton-${i}`} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  if (activities.length === 0) {
    return (
      <div className="text-center text-muted-foreground py-8">
        {t("empty")}
      </div>
    );
  }

  return (
    <DataTable
      columns={columns}
      data={activities}
      isLoading={isLoading}
      emptyMessage={t("empty")}
      itemName="activity"
    />
  );
}

