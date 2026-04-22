"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { MapPin, TrendingUp, Calendar } from "lucide-react";
import { useHeatmapData, useTerritories } from "../hooks/useAreaMapping";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";

export function HeatmapView() {
  const t = useTranslations("areaMapping.heatmap");
  const [territoryId, setTerritoryId] = useState<string>("");
  const [captureType, setCaptureType] = useState<string>("");
  const [startDate, setStartDate] = useState<string>("");
  const [endDate, setEndDate] = useState<string>("");

  // Get territories for dropdown
  const { data: territoriesData } = useTerritories({ page: 1, page_size: 100 });

  // Get heatmap data
  const { data, isLoading, error } = useHeatmapData({
    territory_id: territoryId || undefined,
    capture_type: captureType || undefined,
    start_date: startDate || undefined,
    end_date: endDate || undefined,
  });

  const heatmapPoints = data?.data ?? [];
  const territories = territoriesData?.territories ?? [];

  // Calculate intensity ranges for legend
  const intensities = heatmapPoints.map((p) => p.intensity);
  const maxIntensity = intensities.length > 0 ? Math.max(...intensities) : 0;
  const minIntensity = intensities.length > 0 ? Math.min(...intensities) : 0;
  const avgIntensity =
    intensities.length > 0
      ? intensities.reduce((a, b) => a + b, 0) / intensities.length
      : 0;

  const getIntensityColor = (intensity: number) => {
    if (maxIntensity === 0) return "bg-gray-200";
    const ratio = intensity / maxIntensity;
    if (ratio < 0.33) return "bg-blue-200";
    if (ratio < 0.66) return "bg-yellow-200";
    return "bg-red-500";
  };

  const getIntensityLabel = (intensity: number) => {
    if (maxIntensity === 0) return t("legend.low");
    const ratio = intensity / maxIntensity;
    if (ratio < 0.33) return t("legend.low");
    if (ratio < 0.66) return t("legend.medium");
    return t("legend.high");
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>{t("title")}</CardTitle>
          <CardDescription>{t("description")}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
            <div className="space-y-2">
              <Label htmlFor="territory">{t("filters.territory")}</Label>
              <select
                id="territory"
                value={territoryId}
                onChange={(e) => setTerritoryId(e.target.value)}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                <option value="">All Territories</option>
                {territories.map((territory) => (
                  <option key={territory.id} value={territory.id}>
                    {territory.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="capture_type">Capture Type</Label>
              <select
                id="capture_type"
                value={captureType}
                onChange={(e) => setCaptureType(e.target.value)}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                <option value="">All Types</option>
                <option value="check_in">Check In</option>
                <option value="check_out">Check Out</option>
                <option value="area">Area</option>
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
        </CardContent>
      </Card>

      {error && (
        <Card>
          <CardContent className="pt-6">
            <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4 text-center">
              <p className="text-sm text-destructive">Failed to load heatmap data</p>
            </div>
          </CardContent>
        </Card>
      )}

      {isLoading && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {Array.from({ length: 3 }).map((_, i) => (
            <Card key={i}>
              <CardContent className="pt-6">
                <Skeleton className="h-20 w-full" />
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {!isLoading && !error && (
        <>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-muted-foreground">
                      Total Points
                    </p>
                    <p className="text-2xl font-medium">{heatmapPoints.length}</p>
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
                      Max Intensity
                    </p>
                    <p className="text-2xl font-medium">{maxIntensity}</p>
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
                      Avg Intensity
                    </p>
                    <p className="text-2xl font-medium">
                      {avgIntensity.toFixed(1)}
                    </p>
                  </div>
                  <Calendar className="h-8 w-8 text-muted-foreground" />
                </div>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Heatmap Data Points</CardTitle>
              <CardDescription>
                Location points with activity intensity. Higher intensity indicates more visits.
              </CardDescription>
            </CardHeader>
            <CardContent>
              {heatmapPoints.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <p>No heatmap data available for the selected filters</p>
                </div>
              ) : (
                <div className="space-y-4">
                  <div className="flex items-center gap-4 mb-4">
                    <span className="text-sm font-medium">Legend:</span>
                    <Badge variant="outline" className="bg-blue-200">
                      {t("legend.low")}
                    </Badge>
                    <Badge variant="outline" className="bg-yellow-200">
                      {t("legend.medium")}
                    </Badge>
                    <Badge variant="outline" className="bg-red-500 text-white">
                      {t("legend.high")}
                    </Badge>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 max-h-[600px] overflow-y-auto">
                    {heatmapPoints.slice(0, 100).map((point, index) => (
                      <div
                        key={index}
                        className="rounded-lg border p-4 space-y-2 cursor-pointer hover:bg-accent"
                      >
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-2">
                            <div
                              className={`h-4 w-4 rounded-full ${getIntensityColor(
                                point.intensity
                              )}`}
                            />
                            <span className="text-sm font-medium">
                              Point {index + 1}
                            </span>
                          </div>
                          <Badge variant="outline">
                            {getIntensityLabel(point.intensity)}
                          </Badge>
                        </div>
                        <div className="text-xs text-muted-foreground space-y-1">
                          <p>
                            <span className="font-medium">Location:</span>{" "}
                            {point.lat.toFixed(6)}, {point.lng.toFixed(6)}
                          </p>
                          <p>
                            <span className="font-medium">Intensity:</span>{" "}
                            {point.intensity} visits
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>

                  {heatmapPoints.length > 100 && (
                    <p className="text-sm text-muted-foreground text-center">
                      Showing first 100 of {heatmapPoints.length} points
                    </p>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
}






