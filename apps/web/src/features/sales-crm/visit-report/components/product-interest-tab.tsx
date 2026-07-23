"use client";

import { Skeleton } from "@/components/ui/skeleton";
import { DataTable, type Column } from "@/components/ui/data-table";
import { Badge } from "@/components/ui/badge";
import { useTranslations } from "next-intl";
import type { Activity as ActivityType } from "../types/activity";
import { formatWallClockDateTime } from "@/lib/utils";

interface ProductInterest {
  id: string;
  activity_id: string;
  product_name: string;
  product_id?: string;
  interest_level: number; // 1-5 stars
  quantity: number;
  price: number;
  created_at: string;
  timestamp: string;
  activity_type?: string;
  metadata?: Record<string, unknown>;
}

interface ProductInterestTabProps {
  readonly activities: ActivityType[];
  readonly isLoading: boolean;
}

export function ProductInterestTab({ activities, isLoading }: ProductInterestTabProps) {
  const t = useTranslations("visitReportProductInterest");

  // Extract product interests from all activities
  const productInterests: ProductInterest[] = [];

  activities.forEach((activity) => {
    const interests = Array.isArray(activity.metadata?.product_interests)
      ? activity.metadata.product_interests
      : [];

    interests.forEach(
      (
        item: {
          product_name?: string;
          product_id?: string;
          interest_level?: number;
          quantity?: number;
          price?: number;
        },
        index: number
      ) => {
        productInterests.push({
          id: `${activity.id}-${index}`,
          activity_id: activity.id,
          product_name: item.product_name || "Unknown Product",
          product_id: item.product_id,
          interest_level: item.interest_level || 0,
          quantity: item.quantity || 0,
          price: item.price || 0,
          created_at: activity.created_at,
          timestamp: activity.timestamp,
          activity_type: activity.type,
          metadata: activity.metadata,
        });
      }
    );
  });

  const columns: Column<ProductInterest>[] = [
    {
      id: "product_name",
      header: t("table.product") || "Product",
      accessor: (row) => (
        <div className="flex items-center gap-2">
          <span className="font-medium">{row.product_name}</span>
        </div>
      ),
      className: "min-w-[200px]",
    },
    {
      id: "interest_level",
      header: t("table.interest") || "Interest",
      accessor: (row) => (
        <div className="flex items-center gap-2">
          {row.interest_level > 0 ? (
            <div className="flex gap-0.5">
              {Array.from({ length: 5 }).map((_, i) => (
                <span key={i} className={i < row.interest_level ? "text-yellow-400" : "text-gray-300"}>
                  ⭐
                </span>
              ))}
            </div>
          ) : (
            <Badge variant="outline" className="text-xs">No rating</Badge>
          )}
        </div>
      ),
      className: "w-[150px]",
    },
    {
      id: "quantity",
      header: t("table.quantity") || "Qty",
      accessor: (row) => (
        <div className="text-center">
          <span className="font-medium">{row.quantity}</span>
        </div>
      ),
      className: "w-[80px]",
    },
    {
      id: "price",
      header: t("table.price") || "Price",
      accessor: (row) => (
        <div className="text-right">
          <span className="font-medium">
            Rp {row.price?.toLocaleString("id-ID") || "0"}
          </span>
        </div>
      ),
      className: "w-[150px]",
    },
    {
      id: "timestamp",
      header: t("table.date") || "Date & Time",
      accessor: (row) => {
        let dateStr = row.timestamp;
        // For VISIT type activities, use visit_date from metadata if available
        if (row.activity_type === "visit" && row.metadata && typeof row.metadata === "object") {
          const meta = row.metadata as Record<string, unknown>;
          if (typeof meta.visit_date === "string") {
            dateStr = meta.visit_date;
          }
        }
        const formatted = formatWallClockDateTime(dateStr);
        return <span className="text-sm text-muted-foreground">{formatted}</span>;
      },
      className: "w-[180px]",
    },
  ];

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 3 }, (_, i) => (
          <Skeleton key={`skeleton-${i}`} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  if (productInterests.length === 0) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        <p className="text-sm">{t("empty") || "No product interests recorded yet"}</p>
      </div>
    );
  }

  return (
    <DataTable
      columns={columns}
      data={productInterests}
      isLoading={isLoading}
      emptyMessage={t("empty") || "No product interests"}
      itemName="product-interest"
    />
  );
}
