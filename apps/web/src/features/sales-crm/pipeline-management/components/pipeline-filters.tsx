"use client";

import { Search, Filter, RotateCcw, Plus } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import type { DateRange } from "react-day-picker";
import { useTranslations } from "next-intl";
import { useStages } from "../hooks/useStages";
import { useAccounts } from "@/features/sales-crm/account-management/hooks/useAccounts";
import { useUsers } from "@/features/master-data/user-management/hooks/useUsers";
import type { DealFilters } from "../types";

interface PipelineFiltersProps {
  readonly filters: DealFilters;
  readonly onFiltersChange: (filters: DealFilters) => void;
  readonly onReset: () => void;
  readonly onAdd?: () => void;
  readonly addLabel?: string;
}

export function PipelineFilters({
  filters,
  onFiltersChange,
  onReset,
  onAdd,
  addLabel,
}: PipelineFiltersProps) {
  const t = useTranslations("pipelineManagement.filters");

  const { data: stagesData } = useStages();
  const { data: accountsData } = useAccounts({ status: "active", per_page: 100 });
  const { data: usersData } = useUsers({ status: "active", per_page: 100 });

  const stages = Array.isArray(stagesData) ? stagesData : [];
  const accounts = accountsData?.data ?? [];
  const users = usersData?.data ?? [];

  const currentRange: DateRange | undefined =
    filters.date_from && filters.date_to
      ? {
          from: new Date(`${filters.date_from}T00:00:00`),
          to: new Date(`${filters.date_to}T00:00:00`),
        }
      : filters.date_from
        ? {
            from: new Date(`${filters.date_from}T00:00:00`),
            to: undefined,
          }
        : undefined;

  const handleDateRangeChange = (range: DateRange | undefined) => {
    onFiltersChange({
      ...filters,
      date_from: range?.from
        ? range.from.toISOString().split("T")[0]
        : undefined,
      date_to: range?.to
        ? range.to.toISOString().split("T")[0]
        : undefined,
    });
  };

  const hasActiveFilters =
    filters.search ||
    filters.stage_id ||
    filters.account_id ||
    filters.assigned_to ||
    filters.date_from ||
    filters.date_to;

  return (
    <div className="flex flex-col xl:flex-row items-stretch xl:items-center justify-between gap-4">
      <div className="flex flex-1 flex-col sm:flex-row items-stretch sm:items-center gap-3 flex-wrap">
        {/* Search */}
        <div className="relative flex-1 min-w-[200px] max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder={t("searchPlaceholder")}
            value={filters.search || ""}
            onChange={(e) =>
              onFiltersChange({ ...filters, search: e.target.value || undefined })
            }
            className="pl-10 h-9"
          />
        </div>

        {/* Stage Filter */}
        <Select
          value={filters.stage_id || "all"}
          onValueChange={(value) =>
            onFiltersChange({
              ...filters,
              stage_id: value === "all" ? undefined : value,
            })
          }
        >
          <SelectTrigger className="w-full sm:w-[130px] h-9 text-xs cursor-pointer">
            <SelectValue placeholder={t("stagePlaceholder")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all" className="cursor-pointer">{t("stageAll")}</SelectItem>
            {stages.map((stage: { id: string; name: string }) => (
              <SelectItem key={stage.id} value={stage.id} className="cursor-pointer">
                {stage.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        {/* Account Filter */}
        <Select
          value={filters.account_id || "all"}
          onValueChange={(value) =>
            onFiltersChange({
              ...filters,
              account_id: value === "all" ? undefined : value,
            })
          }
        >
          <SelectTrigger className="w-full sm:w-[150px] h-9 text-xs cursor-pointer">
            <SelectValue placeholder={t("accountPlaceholder")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all" className="cursor-pointer">{t("accountAll")}</SelectItem>
            {accounts.map((account) => (
              <SelectItem key={account.id} value={account.id} className="cursor-pointer">
                {account.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        {/* Assigned To Filter */}
        <Select
          value={filters.assigned_to || "all"}
          onValueChange={(value) =>
            onFiltersChange({
              ...filters,
              assigned_to: value === "all" ? undefined : value,
            })
          }
        >
          <SelectTrigger className="w-full sm:w-[140px] h-9 text-xs cursor-pointer">
            <SelectValue placeholder={t("assignedToPlaceholder")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all" className="cursor-pointer">{t("assignedToAll")}</SelectItem>
            {users.map((user) => (
              <SelectItem key={user.id} value={user.id} className="cursor-pointer">
                {user.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        {/* Date Range Filter */}
        <DateRangePicker
          dateRange={currentRange}
          onDateChange={handleDateRangeChange}
          placeholder={t("dateRangePlaceholder")}
        />

        {/* Reset Button */}
        {hasActiveFilters && (
          <Button
            variant="ghost"
            size="sm"
            onClick={onReset}
            className="h-9 px-2 text-muted-foreground hover:text-primary cursor-pointer"
            title={t("reset")}
          >
            <RotateCcw className="h-4 w-4" />
          </Button>
        )}
      </div>

      <div className="flex items-center gap-2 shrink-0">
        {onAdd && (
          <Button onClick={onAdd} size="sm" className="h-9 cursor-pointer shadow-sm">
            <Plus className="h-4 w-4 mr-2" />
            {addLabel || "Add Item"}
          </Button>
        )}
      </div>
    </div>
  );
}


