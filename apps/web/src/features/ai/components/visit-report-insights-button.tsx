"use client";

import { useState } from "react";
import { Sparkles, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useAnalyzeVisitReport } from "../hooks/useAnalyzeVisitReport";
import type { VisitReportInsight } from "../types";
import { Badge } from "@/components/ui/badge";
import { CheckCircle2, XCircle, AlertCircle } from "lucide-react";
import { useTranslations } from "next-intl";

interface VisitReportInsightsButtonProps {
  visitReportId: string;
  iconOnly?: boolean;
}

export function VisitReportInsightsButton({
  visitReportId,
  iconOnly = false,
}: VisitReportInsightsButtonProps) {
  const t = useTranslations("visitReportInsights");
  const [open, setOpen] = useState(false);
  const [insight, setInsight] = useState<VisitReportInsight | null>(null);
  const { mutate: analyze, isPending } = useAnalyzeVisitReport();

  const handleAnalyze = () => {
    analyze(
      { visit_report_id: visitReportId },
      {
        onSuccess: (response) => {
          setInsight(response.data.data);
          setOpen(true);
        },
      }
    );
  };

  const getSentimentColor = (sentiment: string) => {
    switch (sentiment) {
      case "positive":
        return "bg-green-500/10 text-green-700 dark:text-green-400";
      case "negative":
        return "bg-red-500/10 text-red-700 dark:text-red-400";
      default:
        return "bg-yellow-500/10 text-yellow-700 dark:text-yellow-400";
    }
  };

  const getSentimentIcon = (sentiment: string) => {
    switch (sentiment) {
      case "positive":
        return <CheckCircle2 className="h-4 w-4" />;
      case "negative":
        return <XCircle className="h-4 w-4" />;
      default:
        return <AlertCircle className="h-4 w-4" />;
    }
  };

  const getSentimentLabel = (sentiment: string) => {
    switch (sentiment) {
      case "positive":
        return t("sentiment.positive");
      case "negative":
        return t("sentiment.negative");
      default:
        return t("sentiment.neutral");
    }
  };

  return (
    <>
      <Button
        onClick={handleAnalyze}
        disabled={isPending}
        variant="outline"
        size={iconOnly ? "icon" : "sm"}
        className={iconOnly ? "" : "gap-2"}
        title={iconOnly ? t("button.getInsights") : undefined}
      >
        {isPending ? (
          <Loader2 className="h-4 w-4 animate-spin" />
        ) : (
          <Sparkles className="h-4 w-4" />
        )}
        {!iconOnly && (
          <>
            {isPending ? t("button.analyzing") : t("button.getInsights")}
          </>
        )}
      </Button>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{t("dialog.title")}</DialogTitle>
            <DialogDescription>
              {t("dialog.description")}
            </DialogDescription>
          </DialogHeader>

          {insight && (
            <div className="space-y-6">
              {/* Summary */}
              <div>
                <h3 className="text-sm font-medium mb-2">{t("sections.summary")}</h3>
                <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                  {insight.summary || t("empty.summary")}
                </p>
              </div>

              {/* Sentiment */}
              {insight.sentiment && (
                <div>
                  <h3 className="text-sm font-medium mb-2">{t("sections.sentiment")}</h3>
                  <Badge
                    className={getSentimentColor(insight.sentiment)}
                    variant="outline"
                  >
                    {getSentimentIcon(insight.sentiment)}
                    <span className="ml-1">{getSentimentLabel(insight.sentiment)}</span>
                  </Badge>
                </div>
              )}

              {/* Key Points */}
              {insight.key_points && insight.key_points.length > 0 && (
                <div>
                  <h3 className="text-sm font-medium mb-2">{t("sections.keyPoints")}</h3>
                  <ul className="list-disc list-inside space-y-1 text-sm text-muted-foreground">
                    {insight.key_points.map((point, idx) => (
                      <li key={idx}>{point}</li>
                    ))}
                  </ul>
                </div>
              )}

              {/* Action Items */}
              {insight.action_items && insight.action_items.length > 0 && (
                <div>
                  <h3 className="text-sm font-medium mb-2">{t("sections.actionItems")}</h3>
                  <ul className="list-disc list-inside space-y-1 text-sm text-muted-foreground">
                    {insight.action_items.map((item, idx) => (
                      <li key={idx}>{item}</li>
                    ))}
                  </ul>
                </div>
              )}

              {/* Recommendations */}
              {insight.recommendations && insight.recommendations.length > 0 && (
                <div>
                  <h3 className="text-sm font-medium mb-2">{t("sections.recommendations")}</h3>
                  <ul className="list-disc list-inside space-y-1 text-sm text-muted-foreground">
                    {insight.recommendations.map((rec, idx) => (
                      <li key={idx}>{rec}</li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          )}
        </DialogContent>
      </Dialog>
    </>
  );
}

