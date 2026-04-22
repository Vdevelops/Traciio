"use client";

import React, { useState } from "react";
import { MoreHorizontal, Pencil, Trash2, Star, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { DataTable, Column } from "@/components/ui/data-table";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { Input } from "@/components/ui/input";
import { StatusSwitch } from "@/components/ui/status-switch";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { LeadStatusForm } from "./lead-status-form";
import {
  useLeadStatuses,
  useCreateLeadStatus,
  useUpdateLeadStatus,
  useDeleteLeadStatus,
  useSetDefaultLeadStatus,
} from "../hooks/useLeadStatuses";
import type { LeadStatus, CreateLeadStatusRequest, UpdateLeadStatusRequest } from "../types/lead-status";

export function LeadStatusList(): React.JSX.Element {
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const [isActiveFilter, setIsActiveFilter] = useState<boolean | undefined>(undefined);
  const [sortBy, setSortBy] = useState("order");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("asc");

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedStatus, setSelectedStatus] = useState<LeadStatus | null>(null);

  const { data, isLoading } = useLeadStatuses({
    page,
    per_page: perPage,
    search,
    is_active: isActiveFilter,
    sort_by: sortBy,
    sort_order: sortOrder,
  });

  const createMutation = useCreateLeadStatus();
  const updateMutation = useUpdateLeadStatus();
  const deleteMutation = useDeleteLeadStatus();
  const setDefaultMutation = useSetDefaultLeadStatus();

  const getScoreBadgeColor = (score: number) => {
    if (score >= 60) return "default";
    if (score >= 30) return "secondary";
    return "destructive";
  };

  const columns: Column<LeadStatus>[] = [
    {
      id: "name",
      header: "Name",
      accessor: (row) => (
        <div className="flex items-center gap-2">
          <div
            className="w-3 h-3 rounded-full"
            style={{ backgroundColor: row.color }}
          />
          <span className="font-medium">{row.name}</span>
          {row.is_default && (
            <Star className="w-4 h-4 fill-yellow-500 text-yellow-500" />
          )}
        </div>
      ),
    },
    {
      id: "code",
      header: "Code",
      accessor: (row) => <span className="font-mono text-sm">{row.code}</span>,
    },
    {
      id: "score",
      header: "Score",
      accessor: (row) => (
        <Badge variant={getScoreBadgeColor(row.score)}>
          {row.score}%
        </Badge>
      ),
    },
    {
      id: "order",
      header: "Order",
      accessor: (row) => <span>{row.order}</span>,
    },
    {
      id: "is_active",
      header: "Status",
      accessor: (row) => (
        <StatusSwitch
          checked={row.is_active}
          onCheckedChange={(checked) => {
            updateMutation.mutate({
              id: row.id,
              data: { is_active: checked },
            });
          }}
          disabled={row.is_converted}
        />
      ),
    },
    {
      id: "actions",
      header: "Actions",
      accessor: (row) => (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0 cursor-pointer">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>
            <DropdownMenuItem
              onClick={() => {
                setSelectedStatus(row);
                setIsEditDialogOpen(true);
              }}
              className="cursor-pointer"
            >
              <Pencil className="mr-2 h-4 w-4" />
              Edit
            </DropdownMenuItem>
            {!row.is_default && (
              <DropdownMenuItem
                onClick={() => {
                  setDefaultMutation.mutate(row.id);
                }}
                className="cursor-pointer"
              >
                <Star className="mr-2 h-4 w-4" />
                Set as Default
              </DropdownMenuItem>
            )}
            {(!row.lead_count || row.lead_count === 0) && !row.is_default && !row.is_converted && (
              <>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  onClick={() => {
                    setSelectedStatus(row);
                    setIsDeleteDialogOpen(true);
                  }}
                  className="text-red-600 cursor-pointer"
                >
                  <Trash2 className="mr-2 h-4 w-4" />
                  Delete
                </DropdownMenuItem>
              </>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      ),
    },
  ];

  const handleCreate = async (data: CreateLeadStatusRequest | UpdateLeadStatusRequest) => {
    await createMutation.mutateAsync(data as CreateLeadStatusRequest);
    setIsCreateDialogOpen(false);
  };

  const handleUpdate = async (data: CreateLeadStatusRequest | UpdateLeadStatusRequest) => {
    if (selectedStatus) {
      await updateMutation.mutateAsync({
        id: selectedStatus.id,
        data: data as UpdateLeadStatusRequest,
      });
      setIsEditDialogOpen(false);
      setSelectedStatus(null);
    }
  };

  const handleDelete = async () => {
    if (selectedStatus) {
      await deleteMutation.mutateAsync(selectedStatus.id);
      setIsDeleteDialogOpen(false);
      setSelectedStatus(null);
    }
  };

  // Calculate pagination metadata
  const pagination = data && data.meta
    ? {
        page: data.meta.pagination?.page || 1,
        per_page: data.meta.pagination?.per_page || 10,
        total: data.meta.pagination?.total || 0,
        total_pages: data.meta.pagination?.total_pages || 0,
        has_next: (data.meta.pagination?.page || 0) < (data.meta.pagination?.total_pages || 0),
        has_prev: (data.meta.pagination?.page || 0) > 1,
      }
    : undefined;

  return (
    <div className="space-y-4">
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
        <h2 className="text-xl sm:text-2xl font-medium">Lead Statuses</h2>
        <Button onClick={() => setIsCreateDialogOpen(true)} className="w-full sm:w-auto cursor-pointer">
          <Plus className="mr-2 h-4 w-4" />
          Create Status
        </Button>
      </div>

      <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
        <Input
          placeholder="Search statuses..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="flex-1 sm:max-w-sm"
        />
        <Select
          value={isActiveFilter === undefined ? "all" : isActiveFilter ? "active" : "inactive"}
          onValueChange={(value) => {
            if (value === "all") setIsActiveFilter(undefined);
            else setIsActiveFilter(value === "active");
          }}
        >
          <SelectTrigger className="w-full sm:w-[180px]">
            <SelectValue placeholder="Filter by status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Statuses</SelectItem>
            <SelectItem value="active">Active Only</SelectItem>
            <SelectItem value="inactive">Inactive Only</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <DataTable
        columns={columns}
        data={data?.data || []}
        isLoading={isLoading}
        pagination={pagination}
        onPageChange={setPage}
        onPerPageChange={setPerPage}
        itemName="status"
      />

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto mx-2 sm:mx-auto">
          <DialogHeader>
            <DialogTitle>Create Lead Status</DialogTitle>
          </DialogHeader>
          <LeadStatusForm
            onSubmit={handleCreate}
            isSubmitting={createMutation.isPending}
            onCancel={() => setIsCreateDialogOpen(false)}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto mx-2 sm:mx-auto">
          <DialogHeader>
            <DialogTitle>Edit Lead Status</DialogTitle>
          </DialogHeader>
          {selectedStatus && (
            <LeadStatusForm
              initialData={selectedStatus}
              onSubmit={handleUpdate}
              isSubmitting={updateMutation.isPending}
              onCancel={() => {
                setIsEditDialogOpen(false);
                setSelectedStatus(null);
              }}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Delete Dialog */}
      <DeleteDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        onConfirm={handleDelete}
        title="Delete Lead Status"
        description={
          selectedStatus
            ? `Are you sure you want to delete "${selectedStatus.name}"? This action cannot be undone.`
            : ""
        }
        isLoading={deleteMutation.isPending}
      />
    </div>
  );
}
