"use client";

import { Activity, Clock, Building2, Contact, FileText, Trash2, Edit2 } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import type { Activity as ActivityType } from "../types/activity";
import { useTranslations } from "next-intl";
import { renderIcon } from "../lib/icon-utils";
import { formatWallClockDateTime } from "@/lib/utils";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

interface ActivityTimelineCardProps {
  readonly activities: ActivityType[];
  readonly isLoading: boolean;
  readonly accountId?: string;
  readonly onEdit?: (activity: ActivityType) => void;
  readonly onDelete?: (activityId: string) => void;
}

export function ActivityTimelineCard({
  activities,
  isLoading,
  accountId,
  onEdit,
  onDelete,
}: ActivityTimelineCardProps) {
  const t = useTranslations("visitReportActivityTimeline");

  const formatDateTime = (dateString: string): string => {
    return formatWallClockDateTime(dateString);
  };

  const getActivityDateTime = (activity: ActivityType): string => {
    // For VISIT type activities, use visit_date from metadata
    if (activity.type === "visit" && activity.metadata && typeof activity.metadata === "object") {
      const meta = activity.metadata as Record<string, unknown>;
      if (typeof meta.visit_date === "string") {
        return meta.visit_date;
      }
    }
    // For all other activity types, use timestamp
    return activity.timestamp;
  };

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 3 }, (_, i) => (
          <Skeleton key={`skeleton-${i}`} className="h-24 w-full rounded-lg" />
        ))}
      </div>
    );
  }

  if (activities.length === 0) {
    return (
      <div className="text-center py-12 text-muted-foreground">
        <Activity className="h-8 w-8 mx-auto mb-2 opacity-50" />
        <p className="text-sm">{t("empty")}</p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {activities.map((activity, index) => {
        const activityType = activity.activity_type;
        const iconName = activityType?.icon;
        const badgeColor = activityType?.badge_color as
          | "default"
          | "secondary"
          | "destructive"
          | "outline"
          | undefined;
        const typeName = activityType?.name ?? activity.type;

        // Get product interests from metadata if exists
        const productInterests = Array.isArray(activity.metadata?.product_interests)
          ? activity.metadata.product_interests
          : [];

        return (
          <div
            key={activity.id}
            className="relative p-4 border rounded-lg hover:shadow-sm transition-shadow bg-card"
          >
            {/* Timeline line */}
            {index < activities.length - 1 && (
              <div className="absolute left-6 top-full h-3 w-0.5 bg-border" />
            )}

            <div className="space-y-3">
              {/* Header: Type and Timestamp */}
              <div className="flex items-start justify-between gap-3">
                <div className="flex items-center gap-3 flex-1">
                  {/* Timeline dot */}
                  <div className="flex flex-col items-center">
                    <div className="w-4 h-4 rounded-full border-2 border-primary bg-background" />
                  </div>

                  {/* Type Badge */}
                  <div className="flex items-center gap-2">
                    {iconName ? (
                      renderIcon(iconName, "h-4 w-4 text-muted-foreground")
                    ) : (
                      <Activity className="h-4 w-4 text-muted-foreground" />
                    )}
                    <Badge variant={badgeColor ?? "outline"} className="text-xs font-medium">
                      {typeName}
                    </Badge>
                  </div>
                </div>

                {/* Timestamp */}
                <div className="flex items-center gap-2 text-xs text-muted-foreground whitespace-nowrap">
                  <Clock className="h-3.5 w-3.5" />
                  <span>{formatDateTime(getActivityDateTime(activity))}</span>
                </div>
              </div>

              {/* Description */}
              <div className="ml-[30px]">
                <p className="text-sm font-medium line-clamp-2">{activity.description}</p>
              </div>

              {/* Metadata - Show relevant info */}
              {activity.metadata && typeof activity.metadata === "object" && (() => {
                const meta = activity.metadata as Record<string, unknown>;
                const hasMetadata = Object.keys(meta).length > 0;
                
                if (!hasMetadata) return null;
                
                return (
                  <div className="ml-[30px] space-y-2">
                    {typeof meta.subject === "string" && (
                      <div className="flex items-center gap-2 text-xs">
                        <span className="text-muted-foreground">Subject:</span>
                        <span className="font-medium">{meta.subject}</span>
                      </div>
                    )}
                    {typeof meta.status === "string" && (
                      <div className="flex items-center gap-2 text-xs">
                        <span className="text-muted-foreground">Status:</span>
                        <Badge variant="secondary" className="text-xs">
                          {meta.status}
                        </Badge>
                      </div>
                    )}
                    {typeof meta.duration === "string" && (
                      <div className="flex items-center gap-2 text-xs">
                        <Clock className="h-3 w-3 text-muted-foreground" />
                        <span className="text-muted-foreground">{meta.duration}</span>
                      </div>
                    )}
                    {typeof meta.priority === "string" && (
                      <div className="flex items-center gap-2 text-xs">
                        <span className="text-muted-foreground">Priority:</span>
                        <Badge
                          variant={
                            meta.priority === "high"
                              ? "destructive"
                              : "secondary"
                          }
                          className="text-xs"
                        >
                          {meta.priority}
                        </Badge>
                      </div>
                    )}
                  </div>
                );
              })()}

              {/* Related Info */}
              <div className="ml-[30px] flex flex-wrap gap-2 text-xs">
                {activity.account && (
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <div className="flex items-center gap-1 px-2 py-1 rounded bg-muted text-muted-foreground hover:bg-muted/80 cursor-help">
                          <Building2 className="h-3 w-3" />
                          <span>
                            {typeof activity.account === "object" && "name" in activity.account
                              ? activity.account.name
                              : "Account"}
                          </span>
                        </div>
                      </TooltipTrigger>
                      <TooltipContent>Account linked to this activity</TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                )}
                {activity.contact && (
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <div className="flex items-center gap-1 px-2 py-1 rounded bg-muted text-muted-foreground hover:bg-muted/80 cursor-help">
                          <Contact className="h-3 w-3" />
                          <span>
                            {typeof activity.contact === "object" && "name" in activity.contact
                              ? activity.contact.name
                              : "Contact"}
                          </span>
                        </div>
                      </TooltipTrigger>
                      <TooltipContent>Contact linked to this activity</TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                )}
                {activity.deal && (
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <div className="flex items-center gap-1 px-2 py-1 rounded bg-muted text-muted-foreground hover:bg-muted/80 cursor-help">
                          <FileText className="h-3 w-3" />
                          <span>
                            {typeof activity.deal === "object" && "title" in activity.deal
                              ? activity.deal.title
                              : "Deal"}
                          </span>
                        </div>
                      </TooltipTrigger>
                      <TooltipContent>Deal linked to this activity</TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                )}
              </div>

              {/* Product Interests - If this activity has product interests */}
              {productInterests.length > 0 && (
                <div className="ml-[30px] mt-3 pt-3 border-t">
                  <p className="text-xs font-semibold mb-2">Product Interests:</p>
                  <div className="grid grid-cols-2 gap-2">
                    {productInterests.map(
                      (
                        item: {
                          product_name?: string;
                          interest_level?: number;
                          quantity?: number;
                          price?: number;
                        },
                        idx: number
                      ) => (
                        <div
                          key={idx}
                          className="p-2 bg-muted rounded text-xs space-y-1"
                        >
                          <div className="font-medium">{item.product_name}</div>
                          <div className="flex justify-between">
                            <span className="text-muted-foreground">
                              Interest: {"⭐".repeat(item.interest_level || 0)}
                            </span>
                          </div>
                          <div className="text-muted-foreground">
                            Qty: {item.quantity} | Price: Rp {item.price?.toLocaleString()}
                          </div>
                        </div>
                      )
                    )}
                  </div>
                </div>
              )}

              {/* Actions */}
              {(onEdit || onDelete) && (
                <div className="ml-[30px] flex gap-2">
                  {onEdit && (
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => onEdit(activity)}
                            className="h-8 w-8 p-0"
                          >
                            <Edit2 className="h-3.5 w-3.5" />
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>Edit activity</TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  )}
                  {onDelete && (
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Button
                            size="sm"
                            variant="ghost"
                            className="h-8 w-8 p-0 text-destructive hover:text-destructive hover:bg-destructive/10"
                            onClick={() => onDelete(activity.id)}
                          >
                            <Trash2 className="h-3.5 w-3.5" />
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>Delete activity</TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  )}
                </div>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
