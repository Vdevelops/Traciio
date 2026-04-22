"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useDebounce } from "@/hooks/use-debounce";
import {
  useLeads,
  useDeleteLead,
  useLead,
  useCreateLead,
  useUpdateLead,
} from "./useLeads";
import type {
  CreateLeadFormData,
  UpdateLeadFormData,
} from "../schemas/lead.schema";
import { useTranslations } from "next-intl";

export function useLeadList() {
  const t = useTranslations("leadManagement.detail");
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(20);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);

  useEffect(() => {
    setPage(1);
  }, [debouncedSearch]);
  const [status, setStatus] = useState<string>("");
  const [source, setSource] = useState<string>("");
  const [assignedTo, setAssignedTo] = useState<string>("");
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingLead, setEditingLead] = useState<string | null>(null);
  const [deletingLeadId, setDeletingLeadId] = useState<string | null>(null);

  const { data, isLoading } = useLeads({
    page,
    per_page: perPage,
    search: debouncedSearch,
    status,
    source,
    assigned_to: assignedTo,
  });
  const { data: editingLeadData } = useLead(editingLead || "");
  const deleteLead = useDeleteLead();
  const createLead = useCreateLead();
  const updateLead = useUpdateLead();

  // Keep table order stable even after edits/refetches.
  // Backend ordering isn't guaranteed, so we sort deterministically on the client.
  const leads = (data?.data ?? []).slice().sort((a, b) => {
    const aTime = a.created_at ? Date.parse(a.created_at) : 0;
    const bTime = b.created_at ? Date.parse(b.created_at) : 0;
    if (bTime !== aTime) return bTime - aTime;
    return (a.id || "").localeCompare(b.id || "");
  });
  const pagination = data?.meta?.pagination;

  const handleCreate = async (formData: CreateLeadFormData) => {
    try {
      await createLead.mutateAsync(formData);
      setIsCreateDialogOpen(false);
      toast.success("Lead created successfully");
    } catch (error) {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: UpdateLeadFormData) => {
    if (editingLead) {
      try {
        await updateLead.mutateAsync({ id: editingLead, data: formData });
        setEditingLead(null);
        toast.success("Lead updated successfully");
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handleDeleteClick = (id: string) => {
    setDeletingLeadId(id);
  };

  const handleDeleteConfirm = async () => {
    if (deletingLeadId) {
      try {
        await deleteLead.mutateAsync(deletingLeadId);
        toast.success("Lead deleted successfully");
        setDeletingLeadId(null);
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handlePerPageChange = (newPerPage: number) => {
    setPerPage(newPerPage);
    setPage(1); // Reset to first page when changing per page
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
    source,
    setSource,
    assignedTo,
    setAssignedTo,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    editingLead,
    setEditingLead,
    deletingLeadId,
    setDeletingLeadId,
    // Data
    leads,
    pagination,
    editingLeadData,
    isLoading,
    // Actions
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    // Mutations
    deleteLead,
    createLead,
    updateLead,
  };
}

