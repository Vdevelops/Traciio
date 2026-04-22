"use client";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { useTranslations } from "next-intl";
import type { PipelineReport } from "../types";
import { formatEmailToMailto } from "@/lib/utils";

interface SalesFunnelTableProps {
  readonly data: PipelineReport;
}

export function SalesFunnelTable({ data }: SalesFunnelTableProps) {
  const t = useTranslations("reportsFeature.salesFunnelTable");
  
  // Use actual data from API with null safety
  const deals = data?.deals ?? [];
  const summary = data?.summary ?? {
    total_deals: 0,
    total_value: 0,
    won_deals: 0,
    lost_deals: 0,
    open_deals: 0,
    expected_revenue: 0,
  };
  const grandTotalValue = summary.total_value ?? 0;
  const grandTotalExpected = summary.expected_revenue ?? 0;

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(value);
  };

  const formatDate = (date?: string | null) => {
    if (!date) return "-";
    try {
      const dateObj = new Date(date);
      if (Number.isNaN(dateObj.getTime())) return "-";
      return dateObj.toLocaleDateString("id-ID", {
        year: "numeric",
        month: "short",
        day: "numeric",
      });
    } catch {
      return "-";
    }
  };

  const getStageColor = (stage: string) => {
    const colors: Record<string, string> = {
      Lead: "bg-gray-100 text-gray-800",
      Contacted: "bg-blue-100 text-blue-800",
      Qualified: "bg-yellow-100 text-yellow-800",
      Proposal: "bg-orange-100 text-orange-800",
      Negotiation: "bg-purple-100 text-purple-800",
      Won: "bg-green-100 text-green-800",
      Lost: "bg-red-100 text-red-800",
    };
    return colors[stage] || "bg-gray-100 text-gray-800";
  };

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>{t("title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t("columns.companyName")}</TableHead>
                  <TableHead>{t("columns.contactName")}</TableHead>
                  <TableHead>{t("columns.contactEmail")}</TableHead>
                  <TableHead>{t("columns.stage")}</TableHead>
                  <TableHead className="text-right">
                    {t("columns.value")}
                  </TableHead>
                  <TableHead className="text-right">
                    {t("columns.probability")}
                  </TableHead>
                  <TableHead className="text-right">
                    {t("columns.expectedRevenue")}
                  </TableHead>
                  <TableHead>{t("columns.creationDate")}</TableHead>
                  <TableHead>{t("columns.expectedCloseDate")}</TableHead>
                  <TableHead>{t("columns.teamMember")}</TableHead>
                  <TableHead>{t("columns.progressToWon")}</TableHead>
                  <TableHead>{t("columns.lastInteractedOn")}</TableHead>
                  <TableHead>{t("columns.nextStep")}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {/* Grand Total Row */}
                <TableRow className="bg-red-50 font-medium">
                  <TableCell colSpan={4}>{t("grandTotal")}</TableCell>
                  <TableCell className="text-right">{formatCurrency(grandTotalValue)}</TableCell>
                  <TableCell className="text-right">-</TableCell>
                  <TableCell className="text-right">{formatCurrency(grandTotalExpected)}</TableCell>
                  <TableCell colSpan={6}>-</TableCell>
                </TableRow>

                {/* Deal Rows */}
                {deals.map((deal) => {
                  const dealValue = deal.value ?? 0;
                  const dealExpectedRevenue = deal.expected_revenue ?? 0;
                  const dealProbability = deal.probability ?? 0;
                  const dealProgress = deal.progress_to_won ?? 0;
                  
                  return (
                    <TableRow key={deal.id}>
                      <TableCell className="font-medium">{deal.company_name || "-"}</TableCell>
                      <TableCell>{deal.contact_name || "-"}</TableCell>
                      <TableCell>
                        {deal.contact_email ? (
                          <a href={formatEmailToMailto(deal.contact_email)} className="hover:text-primary hover:underline cursor-pointer min-w-0">
                            {deal.contact_email}
                          </a>
                        ) : (
                          "-"
                        )}
                      </TableCell>
                      <TableCell>
                        <Badge className={getStageColor(deal.stage || "")}>{deal.stage || "-"}</Badge>
                      </TableCell>
                      <TableCell className="text-right">{formatCurrency(dealValue)}</TableCell>
                      <TableCell className="text-right">{dealProbability}%</TableCell>
                      <TableCell className="text-right">{formatCurrency(dealExpectedRevenue)}</TableCell>
                      <TableCell>{formatDate(deal.creation_date)}</TableCell>
                      <TableCell>{formatDate(deal.expected_close_date)}</TableCell>
                      <TableCell>{deal.team_member || "-"}</TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2 w-24">
                          <Progress value={dealProgress} className="flex-1" />
                          <span className="text-xs text-muted-foreground">{dealProgress}%</span>
                        </div>
                      </TableCell>
                      <TableCell>{formatDate(deal.last_interacted_on)}</TableCell>
                      <TableCell>{deal.next_step || "-"}</TableCell>
                    </TableRow>
                  );
                })}

                {/* Empty state if no deals */}
                {deals.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={13} className="text-center py-8 text-muted-foreground">
                      {t("emptyDeals")}
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

