"use client";

import { format, formatDistanceToNow } from "date-fns";
import { id as idLocale } from "date-fns/locale";
import { ArrowRight, CheckCircle2, XCircle, Clock } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useDealHistory } from "../hooks/useDeals";
import type { DealHistory } from "../types";
import { cn } from "@/lib/utils";

interface DealHistoryTimelineProps {
  readonly dealId: string;
}

export function DealHistoryTimeline({ dealId }: DealHistoryTimelineProps) {
  const { data, isLoading, isError, error } = useDealHistory(dealId);

  if (isLoading) {
    return <DealHistoryTimelineSkeleton />;
  }

  if (isError) {
    const errorMessage =
      error instanceof Error ? error.message : "Failed to load deal history";
    return (
      <Alert variant="destructive">
        <XCircle className="h-4 w-4" />
        <AlertDescription>{errorMessage}</AlertDescription>
      </Alert>
    );
  }

  const history: DealHistory[] = data ?? [];

  const getStageLabel = (
    stageName: string | undefined,
    stageObjectName: string | undefined,
    stageId: string | undefined,
    fallback: string
  ) => {
    const normalizedStageName = stageName?.trim();
    if (normalizedStageName) {
      return normalizedStageName;
    }

    const normalizedObjectName = stageObjectName?.trim();
    if (normalizedObjectName) {
      return normalizedObjectName;
    }

    return stageId ? "Unknown" : fallback;
  };

  if (history.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Deal History</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <Clock className="h-12 w-12 text-muted-foreground mb-4" />
            <p className="text-sm text-muted-foreground">
              No stage transitions yet
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Deal History</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="relative space-y-6">
          {/* Timeline line */}
          <div className="absolute left-4 top-2 bottom-2 w-0.5 bg-border" />

          {history.map((entry: DealHistory, index: number) => {
            const isLast = index === history.length - 1;
            const changedAt = entry.changed_at
              ? new Date(entry.changed_at)
              : null;
            const fromStageLabel = getStageLabel(
              entry.from_stage_name,
              entry.from_stage?.name,
              entry.from_stage_id,
              "Created"
            );
            const toStageLabel = getStageLabel(
              entry.to_stage_name,
              entry.to_stage?.name,
              entry.to_stage_id,
              "Unknown"
            );

            const getStageIcon = () => {
              if (entry.to_stage?.is_won) {
                return <CheckCircle2 className="h-5 w-5 text-green-500" />;
              }
              if (entry.to_stage?.is_lost) {
                return <XCircle className="h-5 w-5 text-red-500" />;
              }
              return <ArrowRight className="h-5 w-5 text-blue-500" />;
            };

            const getStageColor = () => {
              if (entry.to_stage?.is_won) return "bg-green-500";
              if (entry.to_stage?.is_lost) return "bg-red-500";
              return "bg-blue-500";
            };

            return (
              <div key={entry.id} className="relative flex gap-4">
                {/* Timeline dot */}
                <div
                  className={cn(
                    "relative z-10 flex h-8 w-8 items-center justify-center rounded-full border-2 border-background",
                    getStageColor()
                  )}
                >
                  <div className="flex h-6 w-6 items-center justify-center rounded-full bg-background">
                    {getStageIcon()}
                  </div>
                </div>

                {/* Content */}
                <div className={cn("flex-1 pb-6", isLast && "pb-0")}>
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1 space-y-2">
                      <div className="flex items-center gap-2 flex-wrap">
                        <Badge variant="outline" className="text-xs">
                          {fromStageLabel}
                        </Badge>
                        <ArrowRight className="h-4 w-4 text-muted-foreground" />
                        <Badge
                          variant="default"
                          className={cn(
                            "text-xs",
                            entry.to_stage?.is_won && "bg-green-500",
                            entry.to_stage?.is_lost && "bg-red-500"
                          )}
                        >
                          {toStageLabel}
                        </Badge>
                        {entry.to_probability !== undefined &&
                          entry.to_probability !== null &&
                          entry.from_probability !== entry.to_probability && (
                            <span className="text-xs text-muted-foreground">
                              {entry.from_probability ?? 0}% →{" "}
                              {entry.to_probability}%
                            </span>
                          )}
                      </div>

                      {entry.reason && (
                        <p className="text-sm text-foreground">
                          <span className="font-medium">Reason:</span>{" "}
                          {entry.reason}
                        </p>
                      )}

                      {entry.notes && (
                        <p className="text-sm text-muted-foreground">
                          {entry.notes}
                        </p>
                      )}

                      <div className="flex items-center gap-4 text-xs text-muted-foreground">
                        <span>
                          By{" "}
                          {entry.changed_by_user?.name ??
                            (typeof entry.changed_by === "string"
                              ? entry.changed_by
                              : "Unknown User")}
                        </span>
                        {changedAt && (
                          <>
                            <span>•</span>
                            <span
                              title={format(changedAt, "PPpp", {
                                locale: idLocale,
                              })}
                            >
                              {formatDistanceToNow(changedAt, {
                                addSuffix: true,
                                locale: idLocale,
                              })}
                            </span>
                          </>
                        )}
                        {entry.days_in_prev_stage !== undefined &&
                          entry.days_in_prev_stage !== null &&
                          entry.days_in_prev_stage > 0 && (
                            <>
                              <span>•</span>
                              <span>
                                {entry.days_in_prev_stage} day
                                {entry.days_in_prev_stage !== 1 ? "s" : ""} in
                                previous stage
                              </span>
                            </>
                          )}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}

function DealHistoryTimelineSkeleton() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Deal History</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-6">
          {[1, 2, 3].map((i) => (
            <div key={i} className="flex gap-4">
              <Skeleton className="h-8 w-8 rounded-full shrink-0" />
              <div className="flex-1 space-y-2">
                <Skeleton className="h-4 w-48" />
                <Skeleton className="h-3 w-full" />
                <Skeleton className="h-3 w-32" />
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
