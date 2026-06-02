"use client";

import { Card } from "@/components/ui/card";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { Calendar, Building2, User, Circle } from "lucide-react";
import type { Task } from "../types";
import { useTranslations } from "next-intl";

interface TaskKanbanCardProps {
  readonly task: Task;
  readonly onClick?: () => void;
}

const statusColorMap: Record<Task["status"], string> = {
  pending: "#f59e0b",
  completed: "#10b981",
};



export function TaskKanbanCard({ task, onClick }: TaskKanbanCardProps) {
  const t = useTranslations("taskManagement.kanbanCard");

  const dueLabel = task.due_date
    ? new Date(task.due_date).toLocaleDateString("id-ID", {
        day: "numeric",
        month: "short",
        year: "numeric",
      })
    : null;

  const accountName = task.account?.name;
  const contactName = task.contact?.name;
  const statusColor = statusColorMap[task.status];

  return (
    <Card
      className="p-4 cursor-pointer hover:shadow-lg hover:border-primary/50 transition-all duration-200 bg-card border border-border"
      onClick={onClick}
    >
      <div className="space-y-3">
        <div className="flex items-start justify-between gap-2">
          <h4 className="font-medium text-base leading-tight line-clamp-2 flex-1">{task.title}</h4>
        </div>

        {/* Description */}
        {task.description && (
          <p className="text-sm text-muted-foreground line-clamp-2">{task.description}</p>
        )}

        {/* Main Info */}
        <div className="space-y-2.5">
          {accountName && (
            <div className="flex items-center gap-2 text-sm">
              <Building2 className="h-4 w-4 text-muted-foreground shrink-0" />
              <span className="text-muted-foreground shrink-0">{t("accountLabel")}</span>
              <span className="font-medium text-foreground truncate">{accountName}</span>
            </div>
          )}

          {contactName && (
            <div className="flex items-center gap-2 text-sm">
              <User className="h-4 w-4 text-muted-foreground shrink-0" />
              <span className="text-muted-foreground shrink-0">{t("contactLabel")}</span>
              <span className="font-medium text-foreground truncate">{contactName}</span>
            </div>
          )}

          {dueLabel && (
            <div className="flex items-center gap-2 text-sm">
              <Calendar className="h-4 w-4 text-muted-foreground shrink-0" />
              <span className="text-muted-foreground shrink-0">{t("dueLabel")}</span>
              <span className="font-medium text-foreground truncate">{dueLabel}</span>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between pt-3 border-t border-border gap-2">
          <div className="flex items-center gap-2 min-w-0">
            <Circle
              className="h-3 w-3 shrink-0"
              style={{ color: statusColor, fill: statusColor }}
            />
            <span className="text-xs text-muted-foreground truncate font-medium capitalize">
              {task.status.replace("_", " ")}
            </span>
          </div>

          {task.assigned_user && (
            <div className="flex items-center gap-2 shrink-0">
              {task.assigned_user.avatar_url && (
                <Avatar className="h-6 w-6 border border-border">
                  <AvatarImage src={task.assigned_user.avatar_url} alt={task.assigned_user.name} />
                </Avatar>
              )}
              <span className="text-xs text-muted-foreground font-medium">
                {task.assigned_user.name}
              </span>
            </div>
          )}
        </div>
      </div>
    </Card>
  );
}
