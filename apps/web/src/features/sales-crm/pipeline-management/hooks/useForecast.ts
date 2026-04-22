import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import * as forecastService from "../services/forecastService";

export const forecastKeys = {
  all: ["pipeline-forecast"] as const,
  forecast: (period: "month" | "quarter" | "year") => 
    [...forecastKeys.all, period] as const,
};

export function useRevenueForecast(period: "month" | "quarter" | "year" = "month") {
  return useQuery({
    queryKey: forecastKeys.forecast(period),
    queryFn: () => forecastService.getRevenueForecast(period),
    select: (response) => response.data,
    staleTime: 60000,
  });
}

// Enhanced hook for Forecast component
export function useForecast() {
  const [period, setPeriod] = useState<"month" | "quarter" | "year">("month");
  const { data: forecast, isLoading } = useRevenueForecast(period);

  const formattedPeriod = (() => {
    const now = new Date();
    switch (period) {
      case "month":
        return now.toLocaleDateString("id-ID", { month: "long", year: "numeric" });
      case "quarter": {
        const quarter = Math.floor(now.getMonth() / 3) + 1;
        return `Q${quarter} ${now.getFullYear()}`;
      }
      case "year":
        return now.getFullYear().toString();
      default:
        return "";
    }
  })();

  return {
    forecast: forecast ?? null,
    isLoading,
    period,
    setPeriod,
    formattedPeriod,
  };
}
