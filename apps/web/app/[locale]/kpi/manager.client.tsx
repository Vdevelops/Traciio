"use client";

import React, { useMemo } from "react";
import { format, startOfMonth, endOfMonth } from "date-fns";
import type { DateRange } from "react-day-picker";
import { useSearchParams } from "next/navigation";
import { usePathname, useRouter } from "@/i18n/routing";
import { Users, Layers, ShieldAlert, ArrowLeft, RefreshCw, ChevronRight } from "lucide-react";

import { useSalesManagerKPI } from "@/features/kpi/hooks/useKPIHooks";
import { useBricks, useBrick } from "@/features/master-data/brick/hooks/useBricks";
import TeamSummaryCard from "@/features/kpi/components/TeamSummaryCard";
import DiagnosticList from "@/features/kpi/components/DiagnosticList";
import SalesRepKPIView from "@/features/kpi/components/sales-rep-kpi-view";
import KPIBarChart from "@/features/kpi/components/KPIBarChart";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { formatCurrency, formatPercent, gradeColor } from "@/features/kpi/utils/formatters";

interface Props {
  lockedBrickId?: string; // Optional, for locking view to a specific brick
  initialStartDate?: string;
  initialEndDate?: string;
}

export default function SalesManagerKPIPageClient({
  lockedBrickId,
  initialStartDate,
  initialEndDate,
}: Props) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  // Read params from URL (single source of truth)
  const startDateParam = searchParams.get("startDate");
  const endDateParam = searchParams.get("endDate");
  const brickIdParam = lockedBrickId ?? searchParams.get("brickId") ?? "all";
  const drilldownUserId = searchParams.get("userId");
  const drilldownUserName = searchParams.get("userName");

  // Get localized strings or fallback defaults
  const { startDate, endDate } = useMemo(() => {
    const today = new Date();
    return {
      startDate: startDateParam ?? initialStartDate ?? format(startOfMonth(today), "yyyy-MM-dd"),
      endDate: endDateParam ?? initialEndDate ?? format(endOfMonth(today), "yyyy-MM-dd"),
    };
  }, [startDateParam, endDateParam, initialStartDate, initialEndDate]);

  const dateRange = useMemo<DateRange | undefined>(() => {
    return {
      from: new Date(startDate),
      to: new Date(endDate),
    };
  }, [startDate, endDate]);

  // Update query parameters in the URL
  const updateUrlParams = (updates: {
    startDate?: string;
    endDate?: string;
    brickId?: string | null;
    userId?: string | null;
    userName?: string | null;
  }) => {
    const params = new URLSearchParams(searchParams.toString());
    
    if (updates.startDate) params.set("startDate", updates.startDate);
    if (updates.endDate) params.set("endDate", updates.endDate);
    
    if (updates.brickId !== undefined) {
      if (updates.brickId === "all" || updates.brickId === null) {
        params.delete("brickId");
      } else {
        params.set("brickId", updates.brickId);
      }
    }
    
    if (updates.userId !== undefined) {
      if (updates.userId === null) {
        params.delete("userId");
        params.delete("userName");
      } else {
        params.set("userId", updates.userId);
        if (updates.userName) params.set("userName", updates.userName);
      }
    }
    
    router.push(`${pathname}?${params.toString()}`);
  };

  // Fetch KPI team data
  const kpiParams = useMemo(() => {
    return {
      startDate,
      endDate,
      brickId: brickIdParam !== "all" ? brickIdParam : undefined,
      includeTeamBreakdown: true,
    };
  }, [startDate, endDate, brickIdParam]);

  const { data, isLoading, error, refetch, isFetching } = useSalesManagerKPI(kpiParams);
  
  // Fetch Brick information for selector
  const { data: bricksResponse } = useBricks({ per_page: 100 });
  const allBricks = useMemo(() => bricksResponse?.data ?? [], [bricksResponse?.data]);

  // Fetch current brick details if view is locked
  const { data: brickDetailsResponse } = useBrick(lockedBrickId ?? "");
  const lockedBrickName = brickDetailsResponse?.data?.name;

  const resp = data ?? null;
  const scopedBrickIds = useMemo(() => new Set(resp?.scope?.bricks ?? []), [resp?.scope?.bricks]);
  const bricks = useMemo(() => {
    if (scopedBrickIds.size === 0) return [];
    return allBricks.filter((brick) => scopedBrickIds.has(brick.id));
  }, [allBricks, scopedBrickIds]);

  const teamHealthChartData = useMemo(() => {
    const summary = resp?.teamSummary;
    if (!summary) return [];

    const taskHealth =
      summary.teamOverdueTaskRate === null || summary.teamOverdueTaskRate === undefined
        ? 0
        : Math.max(0, 100 - summary.teamOverdueTaskRate);

    return [
      { label: "Target", value: Math.min(summary.teamTargetAttainment ?? 0, 100) },
      { label: "Konversi", value: summary.teamConversionRate ?? 0 },
      { label: "Visit", value: summary.teamVisitCompliance ?? 0 },
      { label: "Task", value: taskHealth },
    ];
  }, [resp?.teamSummary]);

  const repScoreChartData = useMemo(() => {
    return (resp?.teamBreakdown ?? []).slice(0, 8).map((item) => ({
      label: item.name.split(" ")[0] ?? item.name,
      value: item.compositeScore,
    }));
  }, [resp?.teamBreakdown]);

  const activeBrickSummary = useMemo(() => {
    if (!resp?.brickBreakdown?.length) return null;

    if (lockedBrickId) {
      return (
        resp.brickBreakdown.find((item) => item.brickId === lockedBrickId) ??
        resp.brickBreakdown[0] ??
        null
      );
    }

    if (brickIdParam !== "all") {
      return resp.brickBreakdown.find((item) => item.brickId === brickIdParam) ?? null;
    }

    return null;
  }, [resp, lockedBrickId, brickIdParam]);

  const visibleDiagnostics = useMemo(() => {
    const diagnostics = resp?.diagnostics ?? [];

    if (!lockedBrickId && brickIdParam === "all") {
      return diagnostics;
    }

    const activeBrickId = lockedBrickId ?? brickIdParam;

    return diagnostics.filter((item) => !item.brickId || item.brickId === activeBrickId);
  }, [resp?.diagnostics, lockedBrickId, brickIdParam]);

  // Handle drilldown click (Personal Rep KPI)
  const handleRepDrilldown = (userId: string, userName: string) => {
    const params = new URLSearchParams();

    params.set("startDate", startDate);
    params.set("endDate", endDate);
    params.set("userName", userName);

    router.push(`/kpi/${userId}?${params.toString()}`);
  };

  const handleClearDrilldown = () => {
    updateUrlParams({ userId: null, userName: null });
  };

  const handleDateRangeChange = (range: DateRange | undefined) => {
    if (range?.from && range.to) {
      updateUrlParams({
        startDate: format(range.from, "yyyy-MM-dd"),
        endDate: format(range.to, "yyyy-MM-dd"),
      });
    }
  };

  const handleBrickChange = (value: string) => {
    updateUrlParams({ brickId: value });
  };

  const handleBrickClick = (id: string) => {
    router.push(`/kpi/brick/${id}?startDate=${startDate}&endDate=${endDate}`);
  };

  const getInitials = (name?: string) => {
    if (!name) return "SR";
    return name
      .split(" ")
      .filter(Boolean)
      .slice(0, 2)
      .map((part) => part[0]?.toUpperCase())
      .join("");
  };

  // If viewing personal drilldown details, render the SalesRepKPIView directly
  if (drilldownUserId) {
    return (
      <SalesRepKPIView
        userId={drilldownUserId}
        userName={drilldownUserName ?? undefined}
        initialStartDate={startDate}
        initialEndDate={endDate}
        onBack={handleClearDrilldown}
      />
    );
  }

  return (
    <div className="space-y-8">
      {/* Header and Controls */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          {lockedBrickId && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => router.push(`/kpi?startDate=${startDate}&endDate=${endDate}`)}
              className="mb-2 -ml-2 text-muted-foreground hover:text-foreground"
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Kembali ke Semua Brick
            </Button>
          )}
          <h1 className="flex items-center gap-2 text-3xl font-medium tracking-tight">
            <Users className="h-8 w-8 text-primary" aria-hidden="true" />
            {lockedBrickId
              ? `KPI Brick: ${lockedBrickName ?? "Memuat..."}`
              : "Evaluasi KPI Tim"}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {lockedBrickId
              ? `Analisis kinerja dan pencapaian target perwakilan penjualan di area ${lockedBrickName ?? ""}`
              : "Pantau pencapaian KPI komposit tim, rincian aktivitas perwakilan, serta analisis area."}
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          {/* Brick Selector - Hidden if view is locked to a specific brick */}
          {!lockedBrickId && (
            <Select value={brickIdParam} onValueChange={handleBrickChange}>
              <SelectTrigger className="h-9 w-full border-border/60 bg-background sm:w-[180px]">
                <SelectValue placeholder="Pilih Brick" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Semua Brick</SelectItem>
                {bricks.map((b) => (
                  <SelectItem key={b.id} value={b.id}>
                    {b.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}

          <DateRangePicker dateRange={dateRange} onDateChange={handleDateRangeChange} />

          <Button
            variant="outline"
            size="icon"
            onClick={() => refetch()}
            disabled={isLoading || isFetching}
            className="h-9 w-9 border-border/60"
          >
            <RefreshCw className={`h-4 w-4 text-muted-foreground ${isFetching ? "animate-spin" : ""}`} />
            <span className="sr-only">Refresh</span>
          </Button>
        </div>
      </div>

      {isLoading ? (
        <div className="space-y-6">
          <div className="h-32 rounded-lg bg-muted" />
          <div className="h-64 rounded-lg bg-muted" />
        </div>
      ) : error ? (
        <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-8 text-center">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10">
            <ShieldAlert className="h-6 w-6 text-destructive" />
          </div>
          <h3 className="mt-4 text-base font-bold text-foreground">
            Gagal Memuat KPI Tim
          </h3>
          <p className="mx-auto mt-2 max-w-md text-sm text-muted-foreground">
            Terjadi kesalahan saat memproses data performa tim. Pastikan koneksi server aman.
          </p>
          <Button onClick={() => refetch()} className="mt-6" variant="outline">
            Coba Lagi
          </Button>
        </div>
      ) : (
        <div className="space-y-6">
          {/* Section 1: Team Summary */}
          <TeamSummaryCard summary={resp?.teamSummary ?? null} />

          <div className="grid gap-6 xl:grid-cols-2">
            <KPIBarChart
              title="Kesehatan KPI Tim"
              description="Normalisasi metrik tim pada periode yang dipilih."
              data={teamHealthChartData}
              valueSuffix="%"
              maxValue={100}
            />
            <KPIBarChart
              title="Ranking Composite Score Sales"
              description="Delapan sales teratas pada scope brick aktif."
              data={repScoreChartData}
              maxValue={100}
            />
          </div>

          {/* Section 2: Team Breakdown Table */}
          <Card>
            <CardHeader>
              <CardTitle>Performa Perwakilan Penjualan</CardTitle>
              <CardDescription>Leaderboard sales yang berada pada brick scope manager.</CardDescription>
            </CardHeader>
            <CardContent className="p-0">
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow className="border-b border-border/60 hover:bg-transparent">
                    <TableHead className="w-[80px] px-4 py-3 text-center text-xs font-medium uppercase tracking-wider text-foreground/90">Peringkat</TableHead>
                    <TableHead className="w-[250px] px-4 py-3 text-xs font-medium uppercase tracking-wider text-foreground/90">Nama Rep</TableHead>
                    <TableHead className="px-4 py-3 text-center text-xs font-medium uppercase tracking-wider text-foreground/90">Composite Score</TableHead>
                    <TableHead className="px-4 py-3 text-center text-xs font-medium uppercase tracking-wider text-foreground/90">Grade</TableHead>
                    <TableHead className="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-foreground/90">Conversion Rate</TableHead>
                    <TableHead className="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-foreground/90">Total Revenue</TableHead>
                    <TableHead className="w-[100px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {!resp?.teamBreakdown || resp.teamBreakdown.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} className="h-24 text-center text-muted-foreground">
                        Tidak ada perwakilan penjualan terdaftar di filter ini.
                      </TableCell>
                    </TableRow>
                  ) : (
                    resp.teamBreakdown.map((t) => {
                      const colors = gradeColor(t.grade);
                      return (
                        <TableRow key={t.userId} className="border-b border-border/30">
                          <TableCell className="px-4 py-3 text-center font-medium text-muted-foreground">
                            #{t.rank}
                          </TableCell>
                          <TableCell className="px-4 py-3 font-medium">
                            <div className="flex items-center gap-3">
                              <Avatar className="h-8 w-8">
                                <AvatarFallback className="bg-muted text-xs font-semibold text-muted-foreground">
                                  {getInitials(t.name)}
                                </AvatarFallback>
                              </Avatar>
                              <div className="flex flex-col">
                                <span className="font-semibold text-foreground">
                                  {t.name}
                                </span>
                              </div>
                            </div>
                          </TableCell>
                          <TableCell className="px-4 py-3 text-center font-semibold text-foreground">
                            {t.compositeScore !== null ? t.compositeScore.toFixed(1) : "N/A"}
                          </TableCell>
                          <TableCell className="px-4 py-3 text-center">
                            {t.grade ? (
                              <span className={`inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-semibold ${colors.bg} ${colors.text} ${colors.border}`}>
                                {t.grade}
                              </span>
                            ) : (
                              <span className="text-muted-foreground">N/A</span>
                            )}
                          </TableCell>
                          <TableCell className="px-4 py-3 text-right">
                            {formatPercent(t.conversionRate)}
                          </TableCell>
                          <TableCell className="px-4 py-3 text-right font-semibold text-foreground">
                            {formatCurrency(t.totalRevenue)}
                          </TableCell>
                          <TableCell className="px-4 py-3">
                            <Button
                              variant="ghost"
                              size="sm"
                              className="cursor-pointer text-primary"
                              onClick={() => handleRepDrilldown(t.userId, t.name)}
                            >
                              Detail
                            </Button>
                          </TableCell>
                        </TableRow>
                      );
                    })
                  )}
                </TableBody>
              </Table>
            </div>
            </CardContent>
          </Card>

          {/* Brick Drilldown Summary */}
          {lockedBrickId && activeBrickSummary && (
            <Card>
              <CardHeader>
                <CardTitle>Ringkasan Brick Aktif</CardTitle>
                <CardDescription>Ringkasan area untuk brick yang sedang dibuka.</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
                  <div className="rounded-lg border bg-card p-4">
                    <div className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
                      Brick
                    </div>
                    <div className="mt-1 text-lg font-semibold text-foreground">
                      {activeBrickSummary.name}
                    </div>
                    <div className="mt-1 text-sm text-muted-foreground">
                      {activeBrickSummary.brickId}
                    </div>
                  </div>

                  <div className="rounded-lg border bg-card p-4">
                    <div className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
                      Composite Score
                    </div>
                    <div className="mt-1 text-lg font-semibold text-foreground">
                      {activeBrickSummary.compositeScore !== null
                        ? activeBrickSummary.compositeScore.toFixed(1)
                        : "Belum ada data"}
                    </div>
                  </div>

                  <div className="rounded-lg border bg-card p-4">
                    <div className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
                      Total Revenue
                    </div>
                    <div className="mt-1 text-lg font-semibold text-foreground">
                      {formatCurrency(activeBrickSummary.totalRevenue)}
                    </div>
                  </div>

                  <div className="rounded-lg border bg-card p-4">
                    <div className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
                      Coverage Penetration
                    </div>
                    <div className="mt-1 text-lg font-semibold text-foreground">
                      {formatPercent(activeBrickSummary.coveragePenetration)}
                    </div>
                    <div className="mt-1 text-sm text-muted-foreground">
                      {activeBrickSummary.repsCount} rep tercakup
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Section 3: Brick Breakdown Table (Hidden if locked to a brick) */}
          {!lockedBrickId && (
            <Card>
              <CardHeader>
                <CardTitle>Peringkat Area (Brick)</CardTitle>
                <CardDescription>Performa setiap brick yang tersedia dalam scope manager.</CardDescription>
              </CardHeader>
              <CardContent className="p-0">
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow className="border-b border-border/60 hover:bg-transparent">
                      <TableHead className="px-4 py-3 text-xs font-medium uppercase tracking-wider text-foreground/90">Nama Brick</TableHead>
                      <TableHead className="px-4 py-3 text-center text-xs font-medium uppercase tracking-wider text-foreground/90">Jumlah Rep</TableHead>
                      <TableHead className="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-foreground/90">Total Revenue</TableHead>
                      <TableHead className="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-foreground/90">Coverage Penetration</TableHead>
                      <TableHead className="px-4 py-3 text-center text-xs font-medium uppercase tracking-wider text-foreground/90">Composite Score</TableHead>
                      <TableHead className="w-[150px]"></TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {!resp?.brickBreakdown || resp.brickBreakdown.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={6} className="h-24 text-center text-muted-foreground">
                          Tidak ada data brick terdaftar.
                        </TableCell>
                      </TableRow>
                    ) : (
                      resp.brickBreakdown.map((b) => {
                        return (
                          <TableRow key={b.brickId} className="border-b border-border/30">
                            <TableCell className="flex items-center gap-2 px-4 py-3 font-medium text-foreground">
                              <Layers className="h-4 w-4 text-muted-foreground" />
                              {b.name}
                            </TableCell>
                            <TableCell className="px-4 py-3 text-center text-muted-foreground">
                              {b.repsCount}
                            </TableCell>
                            <TableCell className="px-4 py-3 text-right font-semibold text-foreground">
                              {formatCurrency(b.totalRevenue)}
                            </TableCell>
                            <TableCell className="px-4 py-3 text-right">
                              {formatPercent(b.coveragePenetration)}
                            </TableCell>
                            <TableCell className="px-4 py-3 text-center font-semibold text-foreground">
                              {b.compositeScore !== null ? b.compositeScore.toFixed(1) : "N/A"}
                            </TableCell>
                            <TableCell className="px-4 py-3">
                              <Button
                                variant="ghost"
                                size="sm"
                                className="cursor-pointer flex items-center gap-1 text-muted-foreground hover:text-foreground"
                                onClick={() => handleBrickClick(b.brickId)}
                              >
                                Lihat Brick
                                <ChevronRight className="h-4 w-4" />
                              </Button>
                            </TableCell>
                          </TableRow>
                        );
                      })
                    )}
                  </TableBody>
                </Table>
              </div>
              </CardContent>
            </Card>
          )}

          {/* Section 4: Team-level DiagnosticList */}
          <DiagnosticList items={visibleDiagnostics} onBrickClick={handleBrickClick} />
        </div>
      )}
    </div>
  );
}
