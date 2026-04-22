"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { MapPin, TrendingUp, Calendar, Users } from "lucide-react";
import { useCoverageAnalysis, useTerritories } from "../hooks/useAreaMapping";
import { formatDate } from "@/lib/utils";
import { Skeleton } from "@/components/ui/skeleton";

export function CoverageAnalysis() {
  const t = useTranslations("areaMapping.coverage");
  const [territoryId, setTerritoryId] = useState<string>("");
  const [startDate, setStartDate] = useState<string>("");
  const [endDate, setEndDate] = useState<string>("");
  const [isAnalyzing, setIsAnalyzing] = useState(false);

  // Get territories for dropdown
  const { data: territoriesData } = useTerritories({ page: 1, page_size: 100 });

  // Get coverage analysis
  const { data, isLoading, error, refetch } = useCoverageAnalysis(
    territoryId && startDate && endDate
      ? {
          territory_id: territoryId,
          start_date: startDate,
          end_date: endDate,
        }
      : undefined
  );

  const handleAnalyze = async () => {
    if (!territoryId || !startDate || !endDate) {
      return;
    }
    setIsAnalyzing(true);
    try {
      await refetch();
    } finally {
      setIsAnalyzing(false);
    }
  };

  const analysis = data?.analysis;
  const territories = territoriesData?.territories ?? [];
  const selectedTerritory = territories.find((t) => t.id === territoryId);

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>{t("title")}</CardTitle>
          <CardDescription>
            Analyze visit coverage for territories within a date range
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
            <div className="space-y-2">
              <Label htmlFor="territory">Territory</Label>
              <select
                id="territory"
                value={territoryId}
                onChange={(e) => setTerritoryId(e.target.value)}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                <option value="">Select territory...</option>
                {territories.map((territory) => (
                  <option key={territory.id} value={territory.id}>
                    {territory.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="start_date">Start Date</Label>
              <Input
                id="start_date"
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="end_date">End Date</Label>
              <Input
                id="end_date"
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
              />
            </div>
          </div>

          <Button
            onClick={handleAnalyze}
            disabled={!territoryId || !startDate || !endDate || isAnalyzing}
            className="w-full md:w-auto"
          >
            {isAnalyzing ? "Analyzing..." : "Analyze Coverage"}
          </Button>
        </CardContent>
      </Card>

      {error && (
        <Card>
          <CardContent className="pt-6">
            <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4 text-center">
              <p className="text-sm text-destructive">Failed to load coverage data</p>
            </div>
          </CardContent>
        </Card>
      )}

      {isLoading && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <Card key={i}>
              <CardContent className="pt-6">
                <Skeleton className="h-20 w-full" />
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {analysis && !isLoading && (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">
                      {t("stats.totalVisits")}
                    </p>
                    <p className="text-2xl font-medium">{analysis.visit_count}</p>
                  </div>
                  <MapPin className="h-8 w-8 text-muted-foreground" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">
                      {t("stats.avgCoverage")}
                    </p>
                    <p className="text-2xl font-medium">
                      {analysis.coverage_percent
                        ? `${analysis.coverage_percent.toFixed(1)}%`
                        : "N/A"}
                    </p>
                  </div>
                  <TrendingUp className="h-8 w-8 text-muted-foreground" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">
                      {t("stats.topTerritory")}
                    </p>
                    <p className="text-lg font-medium truncate">
                      {selectedTerritory?.name ?? "-"}
                    </p>
                  </div>
                  <Users className="h-8 w-8 text-muted-foreground" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">
                      {t("stats.lastAnalyzed")}
                    </p>
                    <p className="text-sm font-medium">
                      {formatDate(analysis.analyzed_at)}
                    </p>
                  </div>
                  <Calendar className="h-8 w-8 text-muted-foreground" />
                </div>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Analysis Details</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">
                      {t("table.territory")}
                    </p>
                    <p className="text-base font-medium">
                      {selectedTerritory?.name ?? "-"}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">
                      {t("table.period")}
                    </p>
                    <p className="text-base font-medium">
                      {formatDate(analysis.period_start)} - {formatDate(analysis.period_end)}
                    </p>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">
                      {t("table.visitCount")}
                    </p>
                    <p className="text-2xl font-medium">{analysis.visit_count}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">
                      {t("table.coveragePercent")}
                    </p>
                    <p className="text-2xl font-medium">
                      {analysis.coverage_percent
                        ? `${analysis.coverage_percent.toFixed(1)}%`
                        : "N/A"}
                    </p>
                  </div>
                </div>

                <div>
                  <p className="text-sm font-medium text-muted-foreground">
                    {t("table.analyzedAt")}
                  </p>
                  <p className="text-base">{formatDate(analysis.analyzed_at)}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </>
      )}

      {!analysis && !isLoading && !error && territoryId && startDate && endDate && (
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8 text-muted-foreground">
              <p>{t("empty")}</p>
              <p className="text-sm mt-2">Click "Analyze Coverage" to generate analysis</p>
            </div>
          </CardContent>
        </Card>
      )}

      {!territoryId && !startDate && !endDate && (
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8 text-muted-foreground">
              <p>Select territory and date range to analyze coverage</p>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

