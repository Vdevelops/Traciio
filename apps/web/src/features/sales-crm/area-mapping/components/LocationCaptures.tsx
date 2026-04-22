"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { MapPin, Search, Calendar } from "lucide-react";
import { useAreaCaptures } from "../hooks/useAreaMapping";
import { DataTable, type Column } from "@/components/ui/data-table";
import type { AreaCapture } from "../types";
import { formatDate } from "@/lib/utils";

export function LocationCaptures() {
  const t = useTranslations("areaMapping.captures.list");
  const tFilters = useTranslations("areaMapping.captures.filters");
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const [captureType, setCaptureType] = useState<string>("");
  const [dateFrom, setDateFrom] = useState<string>("");
  const [dateTo, setDateTo] = useState<string>("");

  const { data, isLoading, error } = useAreaCaptures({
    page,
    per_page: perPage,
    capture_type: captureType || undefined,
    captured_after: dateFrom || undefined,
    captured_before: dateTo || undefined,
  });

  const columns: Column<AreaCapture>[] = [
    {
      id: "visit_report_id",
      header: t("table.visitReport"),
      accessor: (row) => (
        <div className="flex items-center gap-2">
          <MapPin className="h-4 w-4 text-muted-foreground" />
          <span className="font-mono text-xs">{row.visit_report_id.slice(0, 8)}...</span>
        </div>
      ),
    },
    {
      id: "capture_type",
      header: t("table.type"),
      accessor: (row) => {
        const typeLabels: Record<string, string> = {
          check_in: t("types.check_in"),
          check_out: t("types.check_out"),
          area: t("types.area"),
        };
        return (
          <Badge variant="outline">
            {typeLabels[row.capture_type] || row.capture_type}
          </Badge>
        );
      },
    },
    {
      id: "location",
      header: t("table.location"),
      accessor: (row) => {
        const coords = row.location?.coordinates;

        if (!Array.isArray(coords) || coords.length < 2) {
          return <span className="text-sm text-muted-foreground">-</span>;
        }

        const [lng, lat] = coords;

        if (typeof lng !== "number" || typeof lat !== "number") {
          return <span className="text-sm text-muted-foreground">-</span>;
        }

        return (
          <span className="text-sm font-mono">
            {lat.toFixed(6)}, {lng.toFixed(6)}
          </span>
        );
      },
    },
    {
      id: "address",
      header: t("table.address"),
      accessor: (row) => (
        <span className="text-sm text-muted-foreground">
          {row.address || "-"}
        </span>
      ),
    },
    {
      id: "accuracy",
      header: t("table.accuracy"),
      accessor: (row) => (
        <span className="text-sm">
          {row.accuracy ? `${row.accuracy.toFixed(1)}m` : "-"}
        </span>
      ),
    },
    {
      id: "captured_at",
      header: t("table.capturedAt"),
      accessor: (row) => (
        <div className="flex items-center gap-2 text-sm">
          <Calendar className="h-4 w-4 text-muted-foreground" />
          <span>{formatDate(row.captured_at)}</span>
        </div>
      ),
    },
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("title")}</CardTitle>
        <CardDescription>
          View and manage location captures from visit reports
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="mb-4 space-y-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              placeholder={t("searchPlaceholder")}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9"
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">{tFilters("allTypes")}</label>
              <select
                value={captureType}
                onChange={(e) => setCaptureType(e.target.value)}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
              >
                <option value="">{tFilters("allTypes")}</option>
                <option value="check_in">{t("types.check_in")}</option>
                <option value="check_out">{t("types.check_out")}</option>
                <option value="area">{t("types.area")}</option>
              </select>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">{tFilters("dateFrom")}</label>
              <Input
                type="date"
                value={dateFrom}
                onChange={(e) => setDateFrom(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">{tFilters("dateTo")}</label>
              <Input
                type="date"
                value={dateTo}
                onChange={(e) => setDateTo(e.target.value)}
              />
            </div>
          </div>
        </div>

        {error ? (
          <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4 text-center">
            <p className="text-sm text-destructive">Failed to load captures</p>
          </div>
        ) : (
          <DataTable
            columns={columns}
            data={data?.captures ?? []}
            isLoading={isLoading}
            emptyMessage={t("empty")}
            pagination={
              data
                ? {
                    page: data.page,
                    per_page: data.page_size,
                    total: data.total,
                    total_pages: Math.ceil(data.total / data.page_size),
                    has_next: data.page < Math.ceil(data.total / data.page_size),
                    has_prev: data.page > 1,
                  }
                : undefined
            }
            onPageChange={setPage}
            onPerPageChange={setPerPage}
            itemName="capture"
          />
        )}
      </CardContent>
    </Card>
  );
}

