"use client";

import React, { useState } from "react";
import { MoreHorizontal, Pencil, Trash2, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { DataTable, Column } from "@/components/ui/data-table";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { Input } from "@/components/ui/input";
import { StatusSwitch } from "@/components/ui/status-switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { LeadSourceForm } from "./lead-source-form";
import {
  useLeadSources,
  useCreateLeadSource,
  useUpdateLeadSource,
  useDeleteLeadSource,
} from "../hooks/useLeadSources";
import type {
  LeadSource,
  CreateLeadSourceRequest,
  UpdateLeadSourceRequest,
} from "../types/lead-source";

export function LeadSourceList(): React.JSX.Element {
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  const [search, setSearch] = useState("");
  const [isActiveFilter, setIsActiveFilter] = useState<boolean | undefined>(
    undefined,
  );
  const [sortBy] = useState("order");
  const [sortOrder] = useState<"asc" | "desc">("asc");

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedLeadSource, setSelectedLeadSource] =
    useState<LeadSource | null>(null);

  const { data, isLoading } = useLeadSources({
    page,
    per_page: perPage,
    search,
    is_active: isActiveFilter,
    sort_by: sortBy,
    sort_order: sortOrder,
  });

  const createMutation = useCreateLeadSource();
  const updateMutation = useUpdateLeadSource();
  const deleteMutation = useDeleteLeadSource();

  const columns: Column<LeadSource>[] = [
    {
      id: "name",
      header: "Name",
      accessor: (row) => <span className="font-medium">{row.name}</span>,
    },
    {
      id: "code",
      header: "Code",
      accessor: (row) => <span className="font-mono text-sm">{row.code}</span>,
    },
    {
      id: "description",
      header: "Description",
      accessor: (row) => (
        <span className="text-sm text-muted-foreground truncate max-w-xs">
          {row.description || "-"}
        </span>
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
            return updateMutation.mutateAsync({
              id: row.id,
              data: { is_active: checked },
            });
          }}
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
                setSelectedLeadSource(row);
                setIsEditDialogOpen(true);
              }}
              className="cursor-pointer"
            >
              <Pencil className="mr-2 h-4 w-4" />
              Edit
            </DropdownMenuItem>
            {(!row.lead_count || row.lead_count === 0) && (
              <>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  onClick={() => {
                    setSelectedLeadSource(row);
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

  const handleCreate = async (
    data: CreateLeadSourceRequest | UpdateLeadSourceRequest,
  ) => {
    await createMutation.mutateAsync(data as CreateLeadSourceRequest);
    setIsCreateDialogOpen(false);
  };

  const handleUpdate = async (
    data: CreateLeadSourceRequest | UpdateLeadSourceRequest,
  ) => {
    if (selectedLeadSource) {
      await updateMutation.mutateAsync({
        id: selectedLeadSource.id,
        data: data as UpdateLeadSourceRequest,
      });
      setIsEditDialogOpen(false);
      setSelectedLeadSource(null);
    }
  };

  const handleDelete = async () => {
    if (selectedLeadSource) {
      await deleteMutation.mutateAsync(selectedLeadSource.id);
      setIsDeleteDialogOpen(false);
      setSelectedLeadSource(null);
    }
  };

  // Calculate pagination metadata
  const pagination =
    data && data.meta
      ? {
          page: data.meta.pagination?.current_page || 1,
          per_page: data.meta.pagination?.per_page || 10,
          total: data.meta.pagination?.total || 0,
          total_pages: data.meta.pagination?.total_pages || 0,
          has_next:
            (data.meta.pagination?.current_page || 0) <
            (data.meta.pagination?.total_pages || 0),
          has_prev: (data.meta.pagination?.current_page || 0) > 1,
        }
      : undefined;

  return (
    <div className="space-y-4">
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
        <h2 className="text-xl sm:text-2xl font-medium">Lead Sources</h2>
        <Button
          onClick={() => setIsCreateDialogOpen(true)}
          className="cursor-pointer w-full sm:w-auto"
        >
          <Plus className="mr-2 h-4 w-4" />
          Create Lead Source
        </Button>
      </div>

      <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
        <Input
          placeholder="Search lead sources..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="flex-1 sm:max-w-sm"
        />
        <Select
          value={
            isActiveFilter === undefined
              ? "all"
              : isActiveFilter
                ? "active"
                : "inactive"
          }
          onValueChange={(value) => {
            if (value === "all") setIsActiveFilter(undefined);
            else setIsActiveFilter(value === "active");
          }}
        >
          <SelectTrigger className="w-full sm:w-[180px]">
            <SelectValue placeholder="Filter by status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Lead Sources</SelectItem>
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
        itemName="lead source"
      />

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto mx-2 sm:mx-auto">
          <DialogHeader>
            <DialogTitle>Create Lead Source</DialogTitle>
          </DialogHeader>
          <LeadSourceForm
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
            <DialogTitle>Edit Lead Source</DialogTitle>
          </DialogHeader>
          {selectedLeadSource && (
            <LeadSourceForm
              initialData={selectedLeadSource}
              onSubmit={handleUpdate}
              isSubmitting={updateMutation.isPending}
              onCancel={() => {
                setIsEditDialogOpen(false);
                setSelectedLeadSource(null);
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
        title="Delete Lead Source"
        description={
          selectedLeadSource
            ? `Are you sure you want to delete "${selectedLeadSource.name}"? This action cannot be undone.`
            : ""
        }
        isLoading={deleteMutation.isPending}
      />
    </div>
  );
}
