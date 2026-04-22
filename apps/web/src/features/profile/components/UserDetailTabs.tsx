"use client";

import { Suspense } from "react";
import { useTranslations } from "next-intl";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { MapPin, Package, Building2 } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import dynamic from "next/dynamic";

// Lazy load components with code splitting (reuse from sales-overview)
const SalesRepCheckInMap = dynamic(
  () =>
    import("@/features/sales-overview/components/SalesRepCheckInMap").then((mod) => ({
      default: mod.SalesRepCheckInMap,
    })),
  {
    loading: () => <Skeleton className="h-[500px] w-full" />,
    ssr: false,
  },
);

const SalesRepProductSales = dynamic(
  () =>
    import("@/features/sales-overview/components/SalesRepProductSales").then((mod) => ({
      default: mod.SalesRepProductSales,
    })),
  {
    loading: () => (
      <div className="space-y-4">
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-[400px] w-full" />
      </div>
    ),
    ssr: false,
  },
);

const SalesRepCustomers = dynamic(
  () =>
    import("@/features/sales-overview/components/SalesRepCustomers").then((mod) => ({
      default: mod.SalesRepCustomers,
    })),
  {
    loading: () => (
      <div className="space-y-4">
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-[400px] w-full" />
      </div>
    ),
    ssr: false,
  },
);

interface UserDetailTabsProps {
  readonly userId: string;
  readonly startDate?: string;
  readonly endDate?: string;
  readonly checkInLocationsProps: {
    readonly locations: readonly any[];
    readonly isLoading?: boolean;
    readonly totalVisits?: number;
    readonly page?: number;
    readonly perPage?: number;
    readonly onPageChange?: (page: number) => void;
    readonly onPerPageChange?: (perPage: number) => void;
  };
}

export function UserDetailTabs({
  userId,
  startDate,
  endDate,
  checkInLocationsProps,
}: UserDetailTabsProps) {
  const t = useTranslations("profile");

  return (
    <Tabs defaultValue="locations" className="w-full">
      <TabsList>
        <TabsTrigger value="locations" className="gap-2 cursor-pointer">
          <MapPin className="h-4 w-4" />
          {t("tabs.check_in_locations")}
        </TabsTrigger>
        <TabsTrigger value="products" className="gap-2 cursor-pointer">
          <Package className="h-4 w-4" />
          {t("tabs.products_sold")}
        </TabsTrigger>
        <TabsTrigger value="customers" className="gap-2 cursor-pointer">
          <Building2 className="h-4 w-4" />
          {t("tabs.customers")}
        </TabsTrigger>
      </TabsList>

      <TabsContent value="locations" className="mt-6">
        <Suspense fallback={<Skeleton className="h-[500px] w-full" />}>
          <SalesRepCheckInMap {...checkInLocationsProps} />
        </Suspense>
      </TabsContent>

      <TabsContent value="products" className="mt-6">
        <Suspense
          fallback={
            <div className="space-y-4">
              <Skeleton className="h-10 w-full" />
              <Skeleton className="h-[400px] w-full" />
            </div>
          }
        >
          <SalesRepProductSales userId={userId} startDate={startDate} endDate={endDate} />
        </Suspense>
      </TabsContent>

      <TabsContent value="customers" className="mt-6">
        <Suspense
          fallback={
            <div className="space-y-4">
              <Skeleton className="h-10 w-full" />
              <Skeleton className="h-[400px] w-full" />
            </div>
          }
        >
          <SalesRepCustomers userId={userId} startDate={startDate} endDate={endDate} />
        </Suspense>
      </TabsContent>
    </Tabs>
  );
}
