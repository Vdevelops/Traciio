"use client";

import { Card } from "@/components/ui/card";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { Building2, User, DollarSign, TrendingUp, Calendar, Circle } from "lucide-react";
import type { Deal } from "../types";
import { formatCurrency } from "@/lib/utils";
import { useTranslations } from "next-intl";
import { getProbabilityColor } from "../utils/color";

interface DealCardProps {
  readonly deal: Deal;
  readonly onClick?: () => void;
}

export function DealCard({ deal, onClick }: DealCardProps) {
  const t = useTranslations("pipelineManagement.dealCard");

  const valueFormatted = deal.value_formatted || formatCurrency(deal.value ?? 0);
  const accountName = deal.account?.name || t("unknownAccount");
  const contactName = deal.contact?.name;
  const assignedUserName = deal.assigned_user?.name;
  const stageName = deal.stage?.name || t("unknownStage");
  const stageColor = deal.stage?.color || "#3B82F6";



  return (
    <Card
      className="p-4 cursor-pointer hover:shadow-lg hover:border-primary/50 transition-all duration-200 bg-card border border-border"
      onClick={onClick}
    >
      <div className="space-y-3">
        <div className="flex items-start justify-between gap-2">
          <h4 className="font-medium text-base leading-tight line-clamp-2 flex-1">{deal.title}</h4>
        </div>

        {/* Main Info */}
        <div className="space-y-2.5">
          <div className="flex items-center gap-2 text-sm">
            <Building2 className="h-4 w-4 text-muted-foreground shrink-0" />
            <span className="text-muted-foreground shrink-0">{t("accountLabel")}</span>
            <span className="font-medium text-foreground truncate">{accountName}</span>
          </div>

          {contactName && (
            <div className="flex items-center gap-2 text-sm">
              <User className="h-4 w-4 text-muted-foreground shrink-0" />
              <span className="text-muted-foreground shrink-0">{t("contactLabel")}</span>
              <span className="font-medium text-foreground truncate">{contactName}</span>
            </div>
          )}

          <div className="flex items-center gap-2 text-sm">
            <DollarSign className="h-4 w-4 text-muted-foreground shrink-0" />
            <span className="text-muted-foreground shrink-0">{t("valueLabel")}</span>
            <span className="font-medium text-foreground truncate text-base">
              {valueFormatted || formatCurrency(deal.value ?? 0)}
            </span>
          </div>

          {(deal.probability ?? 0) > 0 && (
            <div className="flex items-center gap-2 text-sm">
              <TrendingUp 
                className="h-4 w-4 shrink-0 transition-colors" 
                style={{ color: getProbabilityColor(deal.probability ?? 0) }}
              />
              <span className="text-muted-foreground shrink-0">{t("probabilityLabel")}</span>
              <span 
                className="font-medium transition-colors"
                style={{ color: getProbabilityColor(deal.probability ?? 0) }}
              >
                {deal.probability ?? 0}%
              </span>
            </div>
          )}

          {deal.expected_close_date && (
            <div className="flex items-center gap-2 text-sm">
              <Calendar className="h-4 w-4 text-muted-foreground shrink-0" />
              <span className="text-muted-foreground shrink-0">{t("expectedCloseLabel")}</span>
              <span className="font-medium text-foreground truncate">
                {new Date(deal.expected_close_date).toLocaleDateString("id-ID", {
                  day: "numeric",
                  month: "short",
                  year: "numeric",
                })}
              </span>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between pt-3 border-t border-border gap-2">
          <div className="flex items-center gap-2 min-w-0">
            <Circle
              className="h-3 w-3 shrink-0"
              style={{ color: stageColor, fill: stageColor }}
            />
            <span className="text-xs text-muted-foreground truncate font-medium">{stageName}</span>
          </div>

          {assignedUserName && deal.assigned_user?.avatar_url && (
            <div className="flex items-center gap-2 shrink-0">
              <Avatar className="h-6 w-6 border border-border">
                <AvatarImage src={deal.assigned_user.avatar_url} alt={assignedUserName} />
              </Avatar>
              <span className="text-xs text-muted-foreground hidden sm:inline font-medium">
                {assignedUserName}
              </span>
            </div>
          )}
        </div>
      </div>
    </Card>
  );
}

