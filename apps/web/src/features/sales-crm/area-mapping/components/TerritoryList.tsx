"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { MapPin, Search } from "lucide-react";
import { useTerritories } from "../hooks/useAreaMapping";
import { DataTable, type Column } from "@/components/ui/data-table";
import { CreateTerritoryDialog } from "./CreateTerritoryDialog";
import { EditTerritoryDialog } from "./EditTerritoryDialog";
import { DeleteTerritoryDialog } from "./DeleteTerritoryDialog";
import type { Territory } from "../types";

export function TerritoryList() {
  const t = useTranslations("areaMapping.territories.list");
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);

  const { data, isLoading, error } = useTerritories({
    page,
    page_size: perPage,
    search,
  });

  const columns: Column<Territory>[] = [
    {
      id: "name",
      header: t("table.name"),
      accessor: (row) => (
        <div className="flex items-center gap-2">
          <MapPin className="h-4 w-4 text-muted-foreground" />
          <span className="font-medium">{row.name}</span>
        </div>
      ),
    },
    {
      id: "description",
      header: t("table.description"),
      accessor: (row) => (
        <span className="text-sm text-muted-foreground">
          {row.description || "-"}
        </span>
      ),
    },
    {
      id: "color",
      header: t("table.color"),
      accessor: (row) => (
        <div className="flex items-center gap-2">
          <div
            className="h-6 w-6 rounded border"
            style={{ backgroundColor: row.color }}
          />
          <span className="text-xs text-muted-foreground">{row.color}</span>
        </div>
      ),
    },
    {
      id: "assigned_to",
      header: t("table.assignedTo"),
      accessor: (row) => (
        <Badge variant="outline">
          {row.assigned_to ? "Assigned" : "Unassigned"}
        </Badge>
      ),
    },
    {
      id: "actions",
      header: t("table.actions"),
      accessor: (row) => (
        <div className="flex items-center gap-1">
          <EditTerritoryDialog territory={row} />
          <DeleteTerritoryDialog territory={row} />
        </div>
      ),
    },
  ];

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>{t("title")}</CardTitle>
          </div>
          <CreateTerritoryDialog />
        </div>
      </CardHeader>
      <CardContent>
        <div className="mb-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              placeholder={t("searchPlaceholder")}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9"
            />
          </div>
        </div>

        {error ? (
          <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4 text-center">
            <p className="text-sm text-destructive">Failed to load territories</p>
          </div>
        ) : (
          <DataTable
            columns={columns}
            data={data?.territories || []}
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
            itemName="territory"
          />
        )}
      </CardContent>
    </Card>
  );
}

