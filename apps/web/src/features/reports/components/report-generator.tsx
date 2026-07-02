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
import {
  useAccountActivityReport,
  usePipelineReport,
  useSalesPerformanceReport,
  useVisitReportReport,
} from "../hooks/useReports";
import { SalesFunnelViewer } from "./sales-funnel-viewer";
import { VisitReportViewer } from "./visit-report-viewer";
import { SalesPerformanceReportViewer } from "./sales-performance-report-viewer";
import { AccountActivityReportViewer } from "./account-activity-report-viewer";
import { reportService } from "../services/reportService";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useTranslations } from "next-intl";
import { useAccounts } from "@/features/sales-crm/account-management/hooks/useAccounts";
import { useUsers } from "@/features/master-data/user-management/hooks/useUsers";
import { useAllLeadStatuses } from "@/features/sales-crm/lead-management/hooks/useLeadStatuses";
import { usePipelines } from "@/features/sales-crm/pipeline-management/hooks/usePipelines";
import type { PipelineStage } from "@/features/sales-crm/pipeline-management/types";
import type { ReportRequestParams } from "../types";

type ReportType = "pipeline" | "visit-report" | "sales-performance" | "account-activity";

export function ReportGenerator() {
  const t = useTranslations("reportsFeature.generator");
  const [startDate, setStartDate] = useState<string>(
    new Date(new Date().setMonth(new Date().getMonth() - 1)).toISOString().split("T")[0]
  );
  const [endDate, setEndDate] = useState<string>(new Date().toISOString().split("T")[0]);
  const [accountId, setAccountId] = useState<string>("");
  const [salesRepId, setSalesRepId] = useState<string>("");
  const [reportType, setReportType] = useState<ReportType>("pipeline");
  const [entityType, setEntityType] = useState<"lead" | "deal">("deal");
  const [status, setStatus] = useState<string>("all");
  const [exportPopoverOpen, setExportPopoverOpen] = useState(false);

  const commonParams: ReportRequestParams = {
    start_date: startDate,
    end_date: endDate,
    account_id: accountId || undefined,
    sales_rep_id: salesRepId || undefined,
  };

  const pipelineParams: ReportRequestParams = {
    ...commonParams,
    entity_type: entityType,
    status: status === "all" ? undefined : status,
  };

  const { data: pipelineReport, isLoading: pipelineLoading } = usePipelineReport(
    pipelineParams,
    { enabled: reportType === "pipeline" }
  );
  const { data: visitReport, isLoading: visitLoading } = useVisitReportReport(
    commonParams,
    { enabled: reportType === "visit-report" }
  );
  const { data: salesPerformanceReport, isLoading: salesPerformanceLoading } = useSalesPerformanceReport(
    commonParams,
    { enabled: reportType === "sales-performance" }
  );
  const { data: accountActivityReport, isLoading: accountActivityLoading } = useAccountActivityReport(
    commonParams,
    { enabled: reportType === "account-activity" && !!accountId }
  );
  const { data: accountsData } = useAccounts({ status: "active", per_page: 100 });
  const { data: usersData } = useUsers({ status: "active", per_page: 100 });
  const { data: leadStatusesData } = useAllLeadStatuses();
  const { data: pipelineStagesData } = usePipelines({ is_active: true });

  const accountOptions = accountsData?.data ?? [];
  const salesOptions = (usersData?.data ?? []).filter(
    (user) => user.role?.code === "sales" || user.role?.code === "sales_manager"
  );
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

  const isPipelineReport = reportType === "pipeline";
  const isAccountActivityReport = reportType === "account-activity";
  const requiresAccount = isAccountActivityReport;

  const exportMutation = useMutation({
    mutationFn: async (format: "csv" | "excel") => {
      let blob: Blob;
      let filenamePrefix: string;

      switch (reportType) {
        case "visit-report":
          blob = await reportService.exportVisitReportReport(commonParams, format);
          filenamePrefix = "visit-report";
          break;
        case "sales-performance":
          blob = await reportService.exportSalesPerformanceReport(commonParams, format);
          filenamePrefix = "sales-performance";
          break;
        case "account-activity":
          blob = await reportService.exportAccountActivityReport(commonParams, format);
          filenamePrefix = "account-activity";
          break;
        case "pipeline":
        default:
          blob = await reportService.exportPipelineReport(pipelineParams, format);
          filenamePrefix = `sales-funnel-${entityType}`;
          break;
      }

      const filename = `${filenamePrefix}-${startDate}-${endDate}.${format === "csv" ? "csv" : "xlsx"}`;

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

  const getViewerTitle = () => {
    switch (reportType) {
      case "visit-report":
        return t("viewerTitles.visitReport");
      case "sales-performance":
        return t("viewerTitles.salesPerformance");
      case "account-activity":
        return t("viewerTitles.accountActivity");
      case "pipeline":
      default:
        return t("viewerTitles.pipeline");
    }
  };

  const renderViewer = () => {
    if (requiresAccount && !accountId) {
      return (
        <div className="py-8 text-center text-sm text-muted-foreground">
          {t("accountRequired")}
        </div>
      );
    }

    switch (reportType) {
      case "visit-report":
        return <VisitReportViewer data={visitReport?.data} isLoading={visitLoading} />;
      case "sales-performance":
        return (
          <SalesPerformanceReportViewer
            data={salesPerformanceReport?.data}
            isLoading={salesPerformanceLoading}
          />
        );
      case "account-activity":
        return (
          <AccountActivityReportViewer
            data={accountActivityReport?.data}
            isLoading={accountActivityLoading}
          />
        );
      case "pipeline":
      default:
        return <SalesFunnelViewer data={pipelineReport?.data} isLoading={pipelineLoading} />;
    }
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
              <Label htmlFor="report-type">{t("reportTypeLabel")}</Label>
              <Select value={reportType} onValueChange={(value) => setReportType(value as ReportType)}>
                <SelectTrigger id="report-type">
                  <SelectValue placeholder={t("reportTypePlaceholder")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="pipeline">{t("reportTypes.pipeline")}</SelectItem>
                  <SelectItem value="visit-report">{t("reportTypes.visitReport")}</SelectItem>
                  <SelectItem value="sales-performance">{t("reportTypes.salesPerformance")}</SelectItem>
                  <SelectItem value="account-activity">{t("reportTypes.accountActivity")}</SelectItem>
                </SelectContent>
              </Select>
            </div>
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

          {isPipelineReport ? (
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
          ) : null}

          <div className="flex items-center gap-2">
            <Popover open={exportPopoverOpen} onOpenChange={setExportPopoverOpen}>
              <PopoverTrigger asChild>
                <Button
                  variant="default"
                  disabled={exportMutation.isPending || (requiresAccount && !accountId)}
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
          <CardTitle>{getViewerTitle()}</CardTitle>
        </CardHeader>
        <CardContent>
          {renderViewer()}
        </CardContent>
      </Card>
    </div>
  );
}
