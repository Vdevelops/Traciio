import { useQuery } from "@tanstack/react-query";
import * as summaryService from "../services/summaryService";

export const summaryKeys = {
  all: ["pipeline-summary"] as const,
  summary: (dateFrom?: string, dateTo?: string) => 
    [...summaryKeys.all, { dateFrom, dateTo }] as const,
};

export function usePipelineSummary(dateFrom?: string, dateTo?: string) {
  return useQuery({
    queryKey: summaryKeys.summary(dateFrom, dateTo),
    queryFn: () => summaryService.getPipelineSummary(dateFrom, dateTo),
    select: (response) => response.data,
    staleTime: 30000,
  });
}
