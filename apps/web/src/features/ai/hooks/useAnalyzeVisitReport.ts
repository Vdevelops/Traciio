import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { aiService } from "../services/aiService";
import type { AnalyzeVisitReportRequest } from "../types";

export function useAnalyzeVisitReport() {
  return useMutation({
    mutationFn: (data: AnalyzeVisitReportRequest) =>
      aiService.analyzeVisitReport(data),
    onSuccess: () => {
      toast.success("Visit report analyzed successfully");
    },
    onError: (error: unknown) => {
      const message =
        error instanceof Error ? error.message : "Failed to analyze visit report";
      toast.error(message);
    },
  });
}

