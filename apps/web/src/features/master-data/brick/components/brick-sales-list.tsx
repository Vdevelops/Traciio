"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Users, BarChart3, Eye } from "lucide-react";
import { useBrickSales } from "../hooks/useBricks";
import { formatEmailToMailto } from "@/lib/utils";
import Link from "next/link";

interface BrickSalesListProps {
  brickId: string;
}

type BrickSalesUser = {
  id: string;
  name: string;
  email: string;
};

export function BrickSalesList({ brickId }: BrickSalesListProps) {
  const t = useTranslations("brickSales");
  const { data, isLoading } = useBrickSales(brickId);

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Users className="h-5 w-5" />
            {t("title")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  const sales: BrickSalesUser[] = data?.data ?? [];

  if (sales.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Users className="h-5 w-5" />
            {t("title")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground text-center py-8">
            {t("noSales")}
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Users className="h-5 w-5" />
          {t("title")} ({sales.length})
        </CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t("name")}</TableHead>
              <TableHead>{t("email")}</TableHead>
              <TableHead>{t("role")}</TableHead>
              <TableHead>{t("status")}</TableHead>
              <TableHead className="text-right">{t("actions")}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sales.map((sale) => (
              <TableRow key={sale.id}>
                <TableCell className="font-medium">{sale.name}</TableCell>
                <TableCell>
                  <a href={formatEmailToMailto(sale.email)} className="text-muted-foreground hover:text-primary hover:underline cursor-pointer min-w-0">
                    {sale.email}
                  </a>
                </TableCell>
                <TableCell>
                  <Badge variant="outline">-</Badge>
                </TableCell>
                <TableCell>
                  <Badge variant="secondary">-</Badge>
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-2">
                    <Link href={`/sales-overview/sales-rep/${sale.id}`}>
                      <Button variant="ghost" size="icon" title={t("viewPerformance")}>
                        <BarChart3 className="h-4 w-4" />
                      </Button>
                    </Link>
                    <Link href={`/sales-overview/sales-rep/${sale.id}`}>
                      <Button variant="ghost" size="icon" title={t("viewDetails")}>
                        <Eye className="h-4 w-4" />
                      </Button>
                    </Link>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}

