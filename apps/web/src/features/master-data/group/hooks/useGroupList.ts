"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useDebounce } from "@/hooks/use-debounce";
import {
  useGroups,
  useDeleteGroup,
  useGroup,
  useCreateGroup,
  useUpdateGroup,
} from "./useGroups";
import { useCreateGroupTargetWithUserAssignment } from "@/features/master-data/group/hooks/useGroupTargets";
import type {
  CreateGroupFormData,
  UpdateGroupFormData,
} from "../schemas/group.schema";
import type { CreateGroupTargetWithUserAssignmentFormData } from "@/features/master-data/group/schemas/group-target.schema";

export function useGroupList() {
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(20);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);

  useEffect(() => {
    setPage(1);
  }, [debouncedSearch]);
  const [status, setStatus] = useState<string>("");
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingGroup, setEditingGroup] = useState<string | null>(null);
  const [deletingGroupId, setDeletingGroupId] = useState<string | null>(
    null
  );
  const [targetGroupId, setTargetGroupId] = useState<string | null>(
    null
  );

  const { data, isLoading } = useGroups({
    page,
    per_page: perPage,
    search: debouncedSearch,
    status,
  });
  const { data: editingGroupData } = useGroup(editingGroup || "");
  const { data: targetGroupData } = useGroup(targetGroupId || "");
  const deleteGroup = useDeleteGroup();
  const createGroup = useCreateGroup();
  const updateGroup = useUpdateGroup();
  const createGroupTarget = useCreateGroupTargetWithUserAssignment();

  const groups = data?.data || [];
  const pagination = data?.meta?.pagination;

  const handleCreate = async (formData: CreateGroupFormData) => {
    try {
      await createGroup.mutateAsync(formData);
      setIsCreateDialogOpen(false);
      toast.success("Group created successfully");
    } catch (error) {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: UpdateGroupFormData) => {
    if (editingGroup) {
      try {
        await updateGroup.mutateAsync({
          id: editingGroup,
          data: formData,
        });
        setEditingGroup(null);
        toast.success("Group updated successfully");
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handleDeleteClick = (id: string) => {
    setDeletingGroupId(id);
  };

  const handleDeleteConfirm = async () => {
    if (deletingGroupId) {
      try {
        await deleteGroup.mutateAsync(deletingGroupId);
        toast.success("Group deleted successfully");
        setDeletingGroupId(null);
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handleSetTarget = (groupId: string) => {
    setTargetGroupId(groupId);
  };

  const handleCreateTarget = async (
    formData: CreateGroupTargetWithUserAssignmentFormData
  ) => {
    try {
      const result = await createGroupTarget.mutateAsync(formData);
      setTargetGroupId(null);
      const totalUsers = result.data?.total_users ?? 0;
      toast.success(
        `Target berhasil dibuat untuk group dan ${totalUsers} user`
      );
    } catch (error) {
      // Error already handled in api-client interceptor
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
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    editingGroup,
    setEditingGroup,
    deletingGroupId,
    setDeletingGroupId,
    targetGroupId,
    setTargetGroupId,
    // Data
    groups,
    pagination,
    editingGroupData: editingGroupData?.data,
    targetGroupData: targetGroupData?.data,
    isLoading,
    // Actions
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    handleSetTarget,
    handleCreateTarget,
    deleteGroup,
    createGroup,
    updateGroup,
    createGroupTarget,
  };
}

