"use client";

import { useTranslations } from "next-intl";
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

interface DateFilterControlsProps {
  filterMode: "year" | "range";
  onFilterModeChange: (mode: "year" | "range") => void;
  selectedYear: number;
  onYearChange: (year: number) => void;
  dateRange: DateRange | undefined;
  onDateRangeChange: (range: DateRange | undefined) => void;
}

export function DateFilterControls({
  filterMode,
  onFilterModeChange,
  selectedYear,
  onYearChange,
  dateRange,
  onDateRangeChange,
}: Readonly<DateFilterControlsProps>) {
  const t = useTranslations("salesOverview"); // Reusing salesOverview translations for consistency
  const currentYear = new Date().getFullYear();
  const years = Array.from({ length: 5 }, (_, i) => currentYear - i);

  return (
    <div className="flex flex-col sm:flex-row gap-3">
      <div className="flex bg-muted p-1 rounded-lg">
        <Button
          variant={filterMode === "year" ? "default" : "ghost"}
          size="sm"
          onClick={() => onFilterModeChange("year")}
          className="flex-1 cursor-pointer"
        >
          {t("yearMode")}
        </Button>
        <Button
          variant={filterMode === "range" ? "default" : "ghost"}
          size="sm"
          onClick={() => onFilterModeChange("range")}
          className="flex-1 cursor-pointer"
        >
          {t("rangeMode")}
        </Button>
      </div>

      {filterMode === "year" ? (
        <Select
          value={selectedYear.toString()}
          onValueChange={(value) => onYearChange(Number.parseInt(value))}
        >
          <SelectTrigger className="w-[180px] cursor-pointer">
            <SelectValue placeholder={t("selectYear")} />
          </SelectTrigger>
          <SelectContent>
            {years.map((year) => (
              <SelectItem key={year} value={year.toString()} className="cursor-pointer">
                {year}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      ) : (
        <DateRangePicker
          dateRange={dateRange}
          onDateChange={onDateRangeChange}
          placeholder={t("selectRange") || "Pick a date"}
        />
      )}
    </div>
  );
}
