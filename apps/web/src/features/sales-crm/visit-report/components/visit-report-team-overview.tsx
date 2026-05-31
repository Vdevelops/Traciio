/* eslint-disable @next/next/no-img-element */
"use client";

import { useState, useMemo } from "react";
import {
  CheckCircle2,
  FileText,
  XCircle,
  ChevronRight,
  Users,
} from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useVisitReports } from "../hooks/useVisitReports";
import { VisitReportDetailModal } from "./visit-report-detail-modal";
import { UserVisitReportsDialog } from "./visit-report-user-dialog";
import { useUsers } from "@/features/master-data/user-management/hooks/useUsers";
import type { VisitReport } from "../types";
import { useTranslations } from "next-intl";

// ─── Constants ────────────────────────────────────────────────────────────────

const MAIN_SKELETONS = ["m1", "m2", "m3", "m4", "m5"];

// ─── Types ────────────────────────────────────────────────────────────────────

export interface SalesRepGroup {
  id: string;
  name: string;
  email: string;
  avatarUrl?: string;
  total: number;
  approved: number;
  submitted: number;
  draft: number;
  rejected: number;
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

/** Always returns a DiceBear URL — uses user's avatar if available, otherwise generates one by seed. */
function getDicebearSrc(seed: string, avatarUrl?: string): string {
  if (avatarUrl && avatarUrl.trim() !== "") return avatarUrl;
  return `https://api.dicebear.com/7.x/lorelei/svg?seed=${encodeURIComponent(seed)}`;
}

/** Increments a single status counter on an existing group entry. */
function incrementStatus(group: SalesRepGroup, status: string): void {
  if (status === "approved") group.approved += 1;
  else if (status === "submitted") group.submitted += 1;
  else if (status === "draft") group.draft += 1;
  else if (status === "rejected") group.rejected += 1;
}

/** Builds a sorted array of SalesRepGroup from raw visit reports and a user lookup map. */
function buildRepGroups(
  reports: VisitReport[],
  userMap: Map<string, { avatar_url?: string; email?: string }>
): SalesRepGroup[] {
  const map = new Map<string, SalesRepGroup>();

  for (const report of reports) {
    const id = report.sales_rep?.id ?? report.sales_rep_id ?? "unknown";
    const name = report.sales_rep?.name ?? "Unknown";
    const userData = userMap.get(id);
    const existing = map.get(id);

    if (existing) {
      existing.total += 1;
      incrementStatus(existing, report.status);
    } else {
      const newGroup: SalesRepGroup = {
        id,
        name,
        email: userData?.email ?? "",
        avatarUrl: userData?.avatar_url,
        total: 1,
        approved: 0,
        submitted: 0,
        draft: 0,
        rejected: 0,
      };
      incrementStatus(newGroup, report.status);
      map.set(id, newGroup);
    }
  }

  return Array.from(map.values()).sort((a, b) => b.total - a.total);
}

// ─── UserCard ─────────────────────────────────────────────────────────────────

interface UserCardProps {
  readonly rep: SalesRepGroup;
  readonly onClick: (rep: SalesRepGroup) => void;
}

function UserCard({ rep, onClick }: UserCardProps) {
  const t = useTranslations("visitReportTeamOverview");
  const avatarSrc = getDicebearSrc(rep.email || rep.name, rep.avatarUrl);
  const hasDraft = rep.draft > 0;

  return (
    <Card
      className="cursor-pointer transition-all hover:shadow-sm hover:border-primary/30 active:scale-[0.99]"
      onClick={() => onClick(rep)}
    >
      <CardContent className="flex items-center gap-3 px-4 py-3">
        {/* DiceBear avatar with relative wrapper for notification dot */}
        <div className="relative shrink-0">
          <img
            src={avatarSrc}
            alt={rep.name}
            className="h-10 w-10 rounded-full bg-primary/5 object-cover"
          />
        </div>

        {/* Name + stats */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-1.5 min-w-0">
            <p className="text-sm font-semibold truncate">{rep.name}</p>
            {/* Draft alert badge — soft indicator for unsubmitted work */}
            {hasDraft && (
              <span className="inline-flex items-center gap-0.5 rounded-full bg-muted px-1.5 py-0 text-[10px] font-medium text-muted-foreground shrink-0">
                <FileText className="h-2.5 w-2.5" />
                {t("card.draftBadge", { count: rep.draft })}
              </span>
            )}
          </div>
          <div className="flex items-center gap-2 mt-1 flex-wrap">
            <span className="text-xs text-muted-foreground">{t("card.visits", { count: rep.total })}</span>
            {rep.approved > 0 && (
              <span className="inline-flex items-center gap-0.5 text-xs text-green-600 dark:text-green-400">
                <CheckCircle2 className="h-3 w-3" />
                {rep.approved}
              </span>
            )}
            {rep.submitted > 0 && (
              <span className="inline-flex items-center gap-0.5 text-xs text-blue-600 dark:text-blue-400 font-medium">
                <FileText className="h-3 w-3" />
                {rep.submitted} logged
              </span>
            )}
            {rep.rejected > 0 && (
              <span className="inline-flex items-center gap-0.5 text-xs text-destructive">
                <XCircle className="h-3 w-3" />
                {rep.rejected}
              </span>
            )}
          </div>
        </div>

        <ChevronRight className="h-4 w-4 text-muted-foreground shrink-0" />
      </CardContent>
    </Card>
  );
}

// ─── VisitReportTeamOverview ─────────────────────────────────────────────────

/**
 * Orchestrates 3 fully independent UI layers rendered as siblings:
 *   Layer 1 — this grid: list of sales rep cards
 *   Layer 2 — UserVisitReportsDialog: visit report list for a selected rep
 *   Layer 3 — VisitReportDetailModal: detail drawer for a single report
 *
 * viewingReportId lives here (not inside the dialog) so the drawer is always a
 * sibling of the dialog in the DOM — eliminating pointer-events / scroll-lock conflicts.
 */
export function VisitReportTeamOverview() {
  const t = useTranslations("visitReportManagement");

  // Layer 2 state
  const [selectedRep, setSelectedRep] = useState<SalesRepGroup | null>(null);
  // Layer 3 state — owned here so the drawer is never nested inside the dialog
  const [viewingReportId, setViewingReportId] = useState<string | null>(null);

  const { data: reportsData, isLoading } = useVisitReports({ per_page: 100 });
  const { data: usersData } = useUsers({ per_page: 100 });

  const repGroups = useMemo<SalesRepGroup[]>(() => {
    const reports: VisitReport[] = reportsData?.data ?? [];
    const users = usersData?.data ?? [];

    const userMap = new Map<string, { avatar_url?: string; email?: string }>();
    for (const u of users) {
      userMap.set(u.id, { avatar_url: u.avatar_url, email: u.email });
    }

    return buildRepGroups(reports, userMap);
  }, [reportsData?.data, usersData?.data]);

  return (
    <div className="space-y-3">
      {/* Layer 1: overview grid */}
      {isLoading && (
        <div className="space-y-2">
          {MAIN_SKELETONS.map((id) => (
            <Skeleton key={id} className="h-16 rounded-xl" />
          ))}
        </div>
      )}
      {!isLoading && repGroups.length === 0 && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16 text-center">
            <Users className="h-12 w-12 text-muted-foreground/30 mb-4" />
            <p className="text-sm text-muted-foreground">
              {t("empty") ?? "No team visit data available"}
            </p>
          </CardContent>
        </Card>
      )}
      {!isLoading && repGroups.length > 0 && (
        <div className="grid grid-cols-1 gap-2">
          {repGroups.map((rep) => (
            <UserCard key={rep.id} rep={rep} onClick={setSelectedRep} />
          ))}
        </div>
      )}

      {/* Layer 2: visit report list for the selected sales rep */}
      {selectedRep && (
        <UserVisitReportsDialog
          open={!!selectedRep}
          onOpenChange={(open) => {
            if (!open) setSelectedRep(null);
          }}
          rep={selectedRep}
          onViewReport={(id) => setViewingReportId(id)}
          isDrawerOpen={!!viewingReportId}
        />
      )}

      {/* Layer 3: detail drawer — sibling of the dialog, never its child */}
      <VisitReportDetailModal
        visitReportId={viewingReportId}
        open={!!viewingReportId}
        onOpenChange={(open) => {
          if (!open) setViewingReportId(null);
        }}
        onVisitReportUpdated={() => {}}
      />
    </div>
  );
}
