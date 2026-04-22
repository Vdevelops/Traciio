"use client";

import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { 
  Briefcase, 
  FolderOpen, 
  CheckCircle2, 
  XCircle, 
  TrendingUp,
  DollarSign 
} from "lucide-react";
import { usePipelineSummary } from "../hooks/usePipelines";
import { formatCurrency } from "@/lib/utils";
import { useTranslations } from "next-intl";

export function PipelineSummary() {
  const { data, isLoading } = usePipelineSummary();
  const t = useTranslations("pipelineManagement.pipelineSummary");

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div>
          <div className="h-8 bg-muted animate-pulse rounded w-64 mb-2" />
          <div className="h-4 bg-muted animate-pulse rounded w-96" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {[...Array(4)].map((_, i) => (
            <Card key={i} className="p-6">
              <div className="space-y-4">
                <div className="h-10 w-10 bg-muted animate-pulse rounded-lg" />
                <div className="space-y-2">
                  <div className="h-4 bg-muted animate-pulse rounded w-24" />
                  <div className="h-8 bg-muted animate-pulse rounded w-16" />
                  <div className="h-4 bg-muted animate-pulse rounded w-32" />
                </div>
              </div>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="space-y-6">
        <div>
          <h3 className="text-2xl font-medium tracking-tight">{t("title")}</h3>
          <p className="text-sm text-muted-foreground mt-1">
            {t("description")}
          </p>
        </div>
        <Card className="p-6 border-border">
          <div className="text-center py-12">
            <p className="text-muted-foreground">{t("noData")}</p>
          </div>
        </Card>
      </div>
    );
  }

  const summary = data;

  const summaryCards = [
    {
      label: t("totalDeals"),
      value: summary.total_deals,
      amount: summary.total_value_formatted || formatCurrency(summary.total_value),
      icon: Briefcase,
      color: "text-primary",
      bgColor: "bg-primary/10",
    },
    {
      label: t("openDeals"),
      value: summary.open_deals,
      amount: summary.open_value_formatted || formatCurrency(summary.open_value),
      icon: FolderOpen,
      color: "text-accent-foreground",
      bgColor: "bg-accent/10",
    },
    {
      label: t("wonDeals"),
      value: summary.won_deals,
      amount: summary.won_value_formatted || formatCurrency(summary.won_value),
      icon: CheckCircle2,
      color: "text-chart-1",
      bgColor: "bg-chart-1/10",
    },
    {
      label: t("lostDeals"),
      value: summary.lost_deals,
      amount: summary.lost_value_formatted || formatCurrency(summary.lost_value),
      icon: XCircle,
      color: "text-destructive",
      bgColor: "bg-destructive/10",
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-2xl font-medium tracking-tight">{t("title")}</h3>
        <p className="text-sm text-muted-foreground mt-1">
          {t("description")}
        </p>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {summaryCards.map((card) => {
          const Icon = card.icon;
          return (
            <Card key={card.label} className="p-6 hover:shadow-md transition-shadow border-border">
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className={`p-2 rounded-lg ${card.bgColor}`}>
                    <Icon className={`h-5 w-5 ${card.color}`} />
                  </div>
                </div>
                <div className="space-y-1">
                  <p className="text-sm font-medium text-muted-foreground">{card.label}</p>
                  <p className={`text-3xl font-medium ${card.color}`}>{card.value}</p>
                  <p className="text-sm font-medium text-foreground mt-2">
                    {card.amount}
                  </p>
                </div>
              </div>
            </Card>
          );
        })}
      </div>

      {summary.by_stage && summary.by_stage.length > 0 && (
        <Card className="p-6 border-border">
          <div className="flex items-center gap-2 mb-6">
            <TrendingUp className="h-5 w-5 text-muted-foreground" />
            <h4 className="text-lg font-medium">{t("performanceByStageTitle")}</h4>
          </div>
          <div className="space-y-3">
            {summary.by_stage.map((stage) => (
              <div 
                key={stage.stage_id} 
                className="flex flex-col sm:flex-row items-start sm:items-center justify-between p-4 rounded-lg bg-muted/50 hover:bg-muted/70 transition-colors gap-3"
              >
                <div className="min-w-0 flex-1">
                  <p className="font-medium text-foreground mb-1">{stage.stage_name}</p>
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary" className="text-xs">
                      {stage.deal_count}{" "}
                      {stage.deal_count === 1 ? t("stageDealsSingular") : t("stageDealsPlural")}
                    </Badge>
                  </div>
                </div>
                <div className="text-left sm:text-right shrink-0">
                  <p className="text-lg font-medium text-foreground">
                    {stage.total_value_formatted || formatCurrency(stage.total_value)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  );
}

