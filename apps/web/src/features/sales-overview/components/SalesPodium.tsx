"use client";

import { useTranslations } from "next-intl";
import { Trophy, Medal, Award } from "lucide-react";
import type { SalesPerformanceListItem } from "../types";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

interface SalesPodiumProps {
  readonly topThree: readonly SalesPerformanceListItem[];
}

const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(value);
};

const formatNumber = (value: number): string => {
  return value.toLocaleString("id-ID");
};

export function SalesPodium({ topThree }: SalesPodiumProps) {
  const t = useTranslations("salesOverview.podium");
  
  // Ensure we have at least 3 items to avoid crashes, though UI might look empty
  const first = topThree[0];
  const second = topThree[1];
  const third = topThree[2];

  if (!first && !second && !third) return null;

  return (
    <div className="relative mb-8">
      {/* Subtle Modern Background Gradient */}
      <div className="absolute inset-0 -mx-6 -my-4 rounded-xl bg-gradient-to-br from-secondary/30 via-background to-accent/20 dark:from-secondary/20 dark:via-background dark:to-accent/10" />

      <div className="relative flex items-end justify-center gap-6 pt-8 pb-4">
        {/* 2nd Place - Using Secondary (Blue) */}
        <div className="flex flex-col items-center z-10">
          <div className="relative mb-3">
             <Avatar className="h-20 w-20 border-2 border-secondary/40 dark:border-secondary/30 shadow-lg ring-1 ring-secondary/20 dark:ring-secondary/10">
                <AvatarImage src={second?.avatar_url} alt={second?.user_name} />
                <AvatarFallback className="bg-gradient-to-br from-secondary/20 to-secondary/40 dark:from-secondary/30 dark:to-secondary/50 text-xl font-bold text-secondary-foreground">
                    {second?.user_name?.substring(0, 2).toUpperCase() ?? "2"}
                </AvatarFallback>
            </Avatar>
            <div className="absolute -top-1 -right-1 bg-secondary rounded-full p-1.5 shadow-md z-20">
              <Medal className="h-3.5 w-3.5 text-secondary-foreground" />
            </div>
          </div>
          <div className="relative bg-gradient-to-b from-secondary/10 to-secondary/30 dark:from-secondary/20 dark:to-secondary/40 rounded-t-xl px-6 py-8 min-w-[140px] text-center border border-secondary/30 dark:border-secondary/40 shadow-lg">
            <div className="absolute inset-0 bg-gradient-to-b from-background/40 to-transparent rounded-t-xl" />
            <div className="relative">
              <div className="text-xs font-medium text-secondary-foreground/80 dark:text-secondary-foreground uppercase mb-1 tracking-wide">
                {t("second")}
              </div>
              <div className="font-medium text-sm mb-1 text-secondary-foreground dark:text-secondary-foreground truncate max-w-[120px]">
                {second?.user_name ?? "-"}
              </div>
              {second && (
                <>
                  <div className="text-xs text-secondary-foreground/70 dark:text-secondary-foreground/80 font-medium">
                    {formatNumber(second.deals_closed)} {t("sold")}
                  </div>
                  <div className="text-xs text-secondary-foreground/80 dark:text-secondary-foreground font-medium mt-1">
                    {formatCurrency(second.total_revenue)}
                  </div>
                </>
              )}
            </div>
          </div>
        </div>

        {/* 1st Place - Using Primary (Orange/Brand) */}
        <div className="flex flex-col items-center z-20">
          <div className="relative mb-3">
            <Avatar className="h-28 w-28 border-2 border-primary/50 dark:border-primary/40 shadow-xl ring-2 ring-primary/20 dark:ring-primary/10">
                <AvatarImage src={first?.avatar_url} alt={first?.user_name} />
                <AvatarFallback className="bg-gradient-to-br from-primary/20 to-primary/40 dark:from-primary/30 dark:to-primary/50 text-3xl font-bold text-primary">
                    {first?.user_name?.substring(0, 2).toUpperCase() ?? "1"}
                </AvatarFallback>
            </Avatar>
            <div className="absolute -top-1 -right-1 bg-primary rounded-full p-2 shadow-lg z-20">
              <Trophy className="h-4 w-4 text-primary-foreground" />
            </div>
          </div>
          <div className="relative bg-gradient-to-b from-primary/10 via-primary/20 to-primary/30 dark:from-primary/20 dark:via-primary/30 dark:to-primary/40 rounded-t-xl px-6 py-12 min-w-[160px] text-center border border-primary/30 dark:border-primary/40 shadow-xl">
            <div className="absolute inset-0 bg-gradient-to-b from-background/50 to-transparent rounded-t-xl" />
            <div className="relative">
              <div className="text-xs font-medium text-primary/90 dark:text-primary uppercase mb-2 tracking-wide">
                {t("first")}
              </div>
              <div className="font-medium text-base mb-2 text-primary dark:text-primary truncate max-w-[140px]">
                {first?.user_name ?? "-"}
              </div>
              {first && (
                <>
                  <div className="text-xs text-primary/80 dark:text-primary/90 font-medium">
                    {formatNumber(first.deals_closed)} {t("sold")}
                  </div>
                  <div className="text-sm text-primary dark:text-primary font-medium mt-1">
                    {formatCurrency(first.total_revenue)}
                  </div>
                </>
              )}
            </div>
          </div>
        </div>

        {/* 3rd Place - Using Muted/Neutral */}
        <div className="flex flex-col items-center z-10">
          <div className="relative mb-3">
             <Avatar className="h-16 w-16 border-2 border-muted-foreground/20 dark:border-muted-foreground/30 shadow-lg ring-1 ring-muted-foreground/10 dark:ring-muted-foreground/10">
                <AvatarImage src={third?.avatar_url} alt={third?.user_name} />
                <AvatarFallback className="bg-gradient-to-br from-muted to-muted-foreground/10 dark:from-muted dark:to-muted-foreground/20 text-lg font-bold text-muted-foreground">
                    {third?.user_name?.substring(0, 2).toUpperCase() ?? "3"}
                </AvatarFallback>
            </Avatar>
            <div className="absolute -top-1 -right-1 bg-muted-foreground/20 dark:bg-muted-foreground/30 rounded-full p-1.5 shadow-md z-20 border border-muted-foreground/30">
              <Award className="h-3.5 w-3.5 text-muted-foreground dark:text-muted-foreground" />
            </div>
          </div>
          <div className="relative bg-gradient-to-b from-muted to-muted-foreground/10 dark:from-muted dark:to-muted-foreground/20 rounded-t-xl px-4 py-6 min-w-[120px] text-center border border-muted-foreground/20 dark:border-muted-foreground/30 shadow-lg">
            <div className="absolute inset-0 bg-gradient-to-b from-background/40 to-transparent rounded-t-xl" />
            <div className="relative">
              <div className="text-xs font-medium text-muted-foreground/80 dark:text-muted-foreground uppercase mb-1 tracking-wide">
                {t("third")}
              </div>
              <div className="font-medium text-xs mb-1 text-muted-foreground dark:text-muted-foreground truncate max-w-[100px]">
                {third?.user_name ?? "-"}
              </div>
              {third && (
                <>
                  <div className="text-xs text-muted-foreground/70 dark:text-muted-foreground/80 font-medium">
                    {formatNumber(third.deals_closed)} {t("sold")}
                  </div>
                  <div className="text-xs text-muted-foreground/80 dark:text-muted-foreground font-medium mt-1">
                    {formatCurrency(third.total_revenue)}
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
