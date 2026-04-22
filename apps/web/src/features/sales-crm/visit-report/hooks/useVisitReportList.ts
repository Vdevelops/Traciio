"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useDebounce } from "@/hooks/use-debounce";
import {
  useVisitReports,
  useDeleteVisitReport,
  useVisitReport,
  useCreateVisitReport,
  useUpdateVisitReport,
  useApproveVisitReport,
  useRejectVisitReport,
} from "./useVisitReports";
import type {
  CreateVisitReportFormData,
  UpdateVisitReportFormData,
} from "../schemas/visit-report.schema";
import { useTranslations } from "next-intl";

export function useVisitReportList() {
  const t = useTranslations("visitReportDetail");
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(20);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);

  useEffect(() => {
    setPage(1);
  }, [debouncedSearch]);
  const [status, setStatus] = useState<string>("");
  const [accountId, setAccountId] = useState<string>("");
  const [dealId, setDealId] = useState<string>("");
  const [startDate, setStartDate] = useState<string>("");
  const [endDate, setEndDate] = useState<string>("");
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingVisitReport, setEditingVisitReport] = useState<string | null>(null);
  const [deletingVisitReportId, setDeletingVisitReportId] = useState<string | null>(null);
  const [approvingVisitReportId, setApprovingVisitReportId] = useState<string | null>(null);
  const [rejectingVisitReportId, setRejectingVisitReportId] = useState<string | null>(null);
  const [rejectReason, setRejectReason] = useState("");

  const { data, isLoading } = useVisitReports({
    page,
    per_page: perPage,
    search: debouncedSearch,
    status,
    account_id: accountId,
    deal_id: dealId,
    start_date: startDate,
    end_date: endDate,
  });
  const { data: editingVisitReportData } = useVisitReport(editingVisitReport || "");
  const deleteVisitReport = useDeleteVisitReport();
  const createVisitReport = useCreateVisitReport();
  const updateVisitReport = useUpdateVisitReport();
  const approveVisitReport = useApproveVisitReport();
  const rejectVisitReport = useRejectVisitReport();

  const visitReports = data?.data || [];
  const pagination = data?.meta?.pagination;

  const handleCreate = async (formData: CreateVisitReportFormData) => {
    try {
      await createVisitReport.mutateAsync(formData);
      setIsCreateDialogOpen(false);
      toast.success("Visit report created successfully");
    } catch (error) {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: UpdateVisitReportFormData) => {
    if (editingVisitReport) {
      try {
        await updateVisitReport.mutateAsync({ id: editingVisitReport, data: formData });
        setEditingVisitReport(null);
        toast.success("Visit report updated successfully");
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handleDeleteClick = (id: string) => {
    setDeletingVisitReportId(id);
  };

  const handleDeleteConfirm = async () => {
    if (deletingVisitReportId) {
      try {
        await deleteVisitReport.mutateAsync(deletingVisitReportId);
        toast.success("Visit report deleted successfully");
        setDeletingVisitReportId(null);
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handlePerPageChange = (newPerPage: number) => {
    setPerPage(newPerPage);
    setPage(1); // Reset to first page when changing per page
  };

  const handleApprove = async (id: string) => {
    setApprovingVisitReportId(id);
    try {
      await approveVisitReport.mutateAsync(id);
      toast.success(t("actions.approveSuccess"));
      setApprovingVisitReportId(null);
    } catch (error) {
      // Error already handled in api-client interceptor
      setApprovingVisitReportId(null);
    }
  };

  const handleRejectClick = (id: string) => {
    setRejectingVisitReportId(id);
  };

  const handleRejectConfirm = async () => {
    if (!rejectingVisitReportId || !rejectReason.trim()) return;
    try {
      await rejectVisitReport.mutateAsync({
        id: rejectingVisitReportId,
        data: { reason: rejectReason },
      });
      toast.success(t("actions.rejectSuccess"));
      setRejectingVisitReportId(null);
      setRejectReason("");
    } catch (error) {
      // Error already handled in api-client interceptor
    }
  };

  const handleSubmit = async (id: string) => {
    try {
      await updateVisitReport.mutateAsync({
        id,
        data: { status: "submitted" },
      });
      toast.success("Visit report submitted successfully");
    } catch (error) {
      // Error already handled in api-client interceptor
    }
  };

  return {
    // State
    page,
    setPage,
    perPage,
    setPerPage: handlePerPageChange,
    search,
    setSearch,
    status,
    setStatus,
    accountId,
    setAccountId,
    dealId,
    setDealId,
    startDate,
    setStartDate,
    endDate,
    setEndDate,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    editingVisitReport,
    setEditingVisitReport,
    deletingVisitReportId,
    setDeletingVisitReportId,
    approvingVisitReportId,
    rejectingVisitReportId,
    setRejectingVisitReportId,
    rejectReason,
    setRejectReason,
    // Data
    visitReports,
    pagination,
    editingVisitReportData,
    isLoading,
    // Actions
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    handleApprove,
    handleRejectClick,
    handleRejectConfirm,
    handleSubmit,
    // Mutations
    deleteVisitReport,
    createVisitReport,
    updateVisitReport,
    approveVisitReport,
    rejectVisitReport,
  };
}

