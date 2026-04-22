"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Calendar, Contact, Factory, CheckCircle2, Edit2, Trash2, Play, X } from "lucide-react";
import type { Task } from "../types";
import { useTranslations } from "next-intl";

interface TaskCardProps {
  readonly task: Task;
  readonly onEdit?: () => void;
  readonly onDelete?: () => void;
  readonly onComplete?: () => void;
  readonly onStart?: () => void;
  readonly onCancel?: () => void;
  readonly onClickTitle?: () => void;
  readonly onClick?: () => void;
  readonly onClickContact?: (contactId: string) => void;
}

const statusColorMap: Record<Task["status"], string> = {
  pending: "bg-amber-400",
  in_progress: "bg-sky-400",
  completed: "bg-emerald-400",
  cancelled: "bg-rose-400",
};

export function TaskCard({ task, onEdit, onDelete, onComplete, onStart, onCancel, onClickTitle, onClick, onClickContact }: TaskCardProps) {
  const t = useTranslations("taskManagement.card");

  const dueLabel =
    task.due_date &&
    new Date(task.due_date).toLocaleDateString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });

  const handleCardClick = (event: React.MouseEvent<HTMLDivElement>) => {
    // Don't trigger card click if clicking on interactive elements or during drag
    const target = event.target as HTMLElement;
    if (
      target.closest("button") ||
      target.closest("a") ||
      target.closest("[role='button']") ||
      (event.nativeEvent as MouseEvent & { dataTransfer?: DataTransfer }).dataTransfer
    ) {
      return;
    }
    onClick?.();
  };

  return (
    <div 
      className={`group relative flex flex-col gap-2 rounded-xl border border-border/60 bg-card/80 p-3 hover:border-primary/40 hover:bg-card transition-colors ${onClick ? "cursor-pointer" : ""}`}
      onClick={handleCardClick}
    >
      {/* Status dot + title */}
      <div className="flex items-start gap-2">
        <span
          className={`mt-1 h-2 w-2 shrink-0 rounded-full ${statusColorMap[task.status]}`}
          aria-hidden="true"
        />
        <div className="min-w-0 flex-1 space-y-1">
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              onClickTitle?.();
            }}
            className="text-left text-sm font-medium text-foreground hover:text-primary focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-primary rounded-xs cursor-pointer"
          >
            {task.title}
          </button>
          {task.description && (
            <p className="text-xs text-muted-foreground line-clamp-2">
              {task.description}
            </p>
          )}
        </div>
      </div>

      {/* Meta row */}
      <div className="flex flex-wrap items-center gap-2 pl-4 text-[11px] text-muted-foreground">
        <Badge
          variant="outline"
          className="border-none bg-primary/5 text-primary px-2 py-0 text-[11px] font-medium rounded-full"
        >
          {task.priority}
        </Badge>
        {dueLabel && (
          <span className="inline-flex items-center gap-1">
            <Calendar className="h-3 w-3" />
            <span>{t("duePrefix", { date: dueLabel })}</span>
          </span>
        )}
        {task.account && (
          <span className="inline-flex items-center gap-1 truncate max-w-[160px]">
            <Factory className="h-3 w-3" />
            <span className="truncate">
              {t("accountLabel")} {task.account?.name || t("unknownAccount")}
            </span>
          </span>
        )}
        {task.contact && (
          <button
            type="button"
            onClick={() => {
              if (onClickContact && task.contact?.id) {
                onClickContact(task.contact.id);
              }
            }}
            className="inline-flex items-center gap-1 truncate max-w-[160px] text-left hover:text-primary cursor-pointer"
          >
            <Contact className="h-3 w-3" />
            <span className="truncate">
              {t("contactLabel")} {task.contact?.name || t("unknownContact")}
            </span>
          </button>
        )}
      </div>

      {/* Actions */}
      <div className="flex items-center justify-end gap-1 pt-1 opacity-0 group-hover:opacity-100 transition-opacity">
        {onStart && task.status === "pending" && (
          <Button
            type="button"
            size="icon-sm"
            variant="ghost"
            className="h-7 w-7 text-sky-500 hover:text-sky-600 hover:bg-sky-500/10"
            onClick={(e) => {
              e.stopPropagation();
              onStart();
            }}
            title={t("markInProgressTooltip")}
          >
            <Play className="h-3.5 w-3.5" />
          </Button>
        )}
        {onCancel && task.status !== "completed" && task.status !== "cancelled" && (
          <Button
            type="button"
            size="icon-sm"
            variant="ghost"
            className="h-7 w-7 text-rose-500 hover:text-rose-600 hover:bg-rose-500/10"
            onClick={(e) => {
              e.stopPropagation();
              onCancel();
            }}
            title={t("markCancelledTooltip")}
          >
            <X className="h-3.5 w-3.5" />
          </Button>
        )}
        {onComplete && task.status !== "completed" && task.status !== "cancelled" && (
          <Button
            type="button"
            size="icon-sm"
            variant="ghost"
            className="h-7 w-7 text-primary hover:text-primary hover:bg-primary/10"
            onClick={(e) => {
              e.stopPropagation();
              onComplete();
            }}
            title={t("markCompletedTooltip")}
          >
            <CheckCircle2 className="h-3.5 w-3.5" />
          </Button>
        )}
        {onEdit && (
          <Button
            type="button"
            size="icon-sm"
            variant="ghost"
            className="h-7 w-7"
            onClick={(e) => {
              e.stopPropagation();
              onEdit();
            }}
            title={t("editTooltip")}
          >
            <Edit2 className="h-3.5 w-3.5" />
          </Button>
        )}
        {onDelete && (
          <Button
            type="button"
            size="icon-sm"
            variant="ghost"
            className="h-7 w-7 text-destructive hover:text-destructive hover:bg-destructive/10"
            onClick={(e) => {
              e.stopPropagation();
              onDelete();
            }}
            title={t("deleteTooltip")}
          >
            <Trash2 className="h-3.5 w-3.5" />
          </Button>
        )}
      </div>
    </div>
  );
}


