"use client";

import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { FileText, Download, FileSpreadsheet, FileText as FileTextIcon, ChevronDown } from "lucide-react";
import { usePipelineReport } from "../hooks/useReports";
import { SalesFunnelViewer } from "./sales-funnel-viewer";
import { reportService } from "../services/reportService";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useTranslations } from "next-intl";
import { useAccounts } from "@/features/sales-crm/account-management/hooks/useAccounts";
import { useUsers } from "@/features/master-data/user-management/hooks/useUsers";
import { useAllLeadStatuses } from "@/features/sales-crm/lead-management/hooks/useLeadStatuses";
import { usePipelines } from "@/features/sales-crm/pipeline-management/hooks/usePipelines";
import type { PipelineStage } from "@/features/sales-crm/pipeline-management/types";

export function ReportGenerator() {
  const t = useTranslations("reportsFeature.generator");
  const [startDate, setStartDate] = useState<string>(
    new Date(new Date().setMonth(new Date().getMonth() - 1)).toISOString().split("T")[0]
  );
  const [endDate, setEndDate] = useState<string>(new Date().toISOString().split("T")[0]);
  const [accountId, setAccountId] = useState<string>("");
  const [salesRepId, setSalesRepId] = useState<string>("");
  const [entityType, setEntityType] = useState<"lead" | "deal">("deal");
  const [status, setStatus] = useState<string>("all");
  const [exportPopoverOpen, setExportPopoverOpen] = useState(false);

  const pipelineParams = {
    start_date: startDate,
    end_date: endDate,
    entity_type: entityType,
    account_id: accountId || undefined,
    sales_rep_id: salesRepId || undefined,
    status: status === "all" ? undefined : status,
  };

  const { data: pipelineReport, isLoading: pipelineLoading } = usePipelineReport(pipelineParams);
  const { data: accountsData } = useAccounts({ status: "active", per_page: 100 });
  const { data: usersData } = useUsers({ status: "active", per_page: 100 });
  const { data: leadStatusesData } = useAllLeadStatuses();
  const { data: pipelineStagesData } = usePipelines({ is_active: true });

  const accountOptions = accountsData?.data ?? [];
  const salesOptions = usersData?.data ?? [];
  const leadStatusOptions = leadStatusesData?.data ?? [];
  const pipelineStageOptions = pipelineStagesData?.data ?? [];

  const statusOptions = entityType === "lead"
    ? leadStatusOptions.map((statusItem) => ({
        value: statusItem.code,
        label: statusItem.name,
      }))
    : pipelineStageOptions.map((stage: PipelineStage) => ({
        value: stage.id,
        label: stage.name,
      }));

  const exportMutation = useMutation({
    mutationFn: async (format: "csv" | "excel") => {
      const blob = await reportService.exportPipelineReport(pipelineParams, format);
      const filename = `sales-funnel-${entityType}-${startDate}-${endDate}.${format === "csv" ? "csv" : "xlsx"}`;

      // Create download link
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    },
    onSuccess: () => {
      toast.success(t("exportSuccess"));
    },
    onError: (error) => {
      toast.error(t("exportError"));
    },
  });

  const handleExport = (format: "csv" | "excel") => {
    exportMutation.mutate(format);
  };

  const handleEntityTypeChange = (value: "lead" | "deal") => {
    setEntityType(value);
    setStatus("all");
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <FileText className="h-5 w-5" />
            {t("title")}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="start-date">{t("startDateLabel")}</Label>
              <Input
                id="start-date"
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="end-date">{t("endDateLabel")}</Label>
              <Input
                id="end-date"
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
              />
            </div>
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="account-id">{t("accountLabel")}</Label>
              <Select value={accountId || "all"} onValueChange={(value) => setAccountId(value === "all" ? "" : value)}>
                <SelectTrigger id="account-id">
                  <SelectValue placeholder={t("accountPlaceholder")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t("accountAll")}</SelectItem>
                  {accountOptions.map((account) => (
                    <SelectItem key={account.id} value={account.id}>
                      {account.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="sales-rep-id">{t("salesLabel")}</Label>
              <Select value={salesRepId || "all"} onValueChange={(value) => setSalesRepId(value === "all" ? "" : value)}>
                <SelectTrigger id="sales-rep-id">
                  <SelectValue placeholder={t("salesPlaceholder")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t("salesAll")}</SelectItem>
                  {salesOptions.map((sales) => (
                    <SelectItem key={sales.id} value={sales.id}>
                      {sales.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="entity-type">{t("entityTypeLabel")}</Label>
              <Select value={entityType} onValueChange={(value) => handleEntityTypeChange(value as "lead" | "deal")}>
                <SelectTrigger id="entity-type">
                  <SelectValue placeholder={t("entityTypePlaceholder")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="lead">{t("entityTypeLead")}</SelectItem>
                  <SelectItem value="deal">{t("entityTypeDeal")}</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="status">{t("statusLabel")}</Label>
              <Select value={status} onValueChange={setStatus}>
                <SelectTrigger id="status">
                  <SelectValue placeholder={t("statusPlaceholder")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t("statusAll")}</SelectItem>
                  {statusOptions.map((option: { value: string; label: string }) => (
                    <SelectItem key={option.value} value={option.value}>
                      {option.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Popover open={exportPopoverOpen} onOpenChange={setExportPopoverOpen}>
              <PopoverTrigger asChild>
                <Button
                  variant="default"
                  disabled={exportMutation.isPending}
                  className="gap-2"
                >
                  <Download className="h-4 w-4" />
                  {exportMutation.isPending ? t("exportingButton") : t("exportButton")}
                  <ChevronDown className="h-4 w-4" />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-48 p-2" align="start">
                <div className="space-y-1">
                  <button
                    onClick={() => {
                      handleExport("csv");
                      setExportPopoverOpen(false);
                    }}
                    disabled={exportMutation.isPending}
                    className="w-full flex items-center gap-2 px-3 py-2 text-sm rounded-md hover:bg-accent transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <FileTextIcon className="h-4 w-4" />
                    {t("exportCsv")}
                  </button>
                  <button
                    onClick={() => {
                      handleExport("excel");
                      setExportPopoverOpen(false);
                    }}
                    disabled={exportMutation.isPending}
                    className="w-full flex items-center gap-2 px-3 py-2 text-sm rounded-md hover:bg-accent transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <FileSpreadsheet className="h-4 w-4" />
                    {t("exportExcel")}
                  </button>
                </div>
              </PopoverContent>
            </Popover>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{t("viewerTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          <SalesFunnelViewer data={pipelineReport?.data} isLoading={pipelineLoading} />
        </CardContent>
      </Card>
    </div>
  );
}
