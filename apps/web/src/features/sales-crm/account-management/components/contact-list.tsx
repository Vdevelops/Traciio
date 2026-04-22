"use client";

import { useState } from "react";
import { Edit, Trash2, Plus, Search, Eye } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { DataTable, type Column } from "@/components/ui/data-table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useContactList } from "../hooks/useContactList";
import { ContactForm } from "./contact-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useAccounts } from "../hooks/useAccounts";
import type { Contact, ContactRole } from "../types";
import { ContactDetailModal } from "./contact-detail-modal";
import type { CreateContactFormData, UpdateContactFormData } from "../schemas/contact.schema";
import { formatPhoneNumberToWA, formatEmailToMailto } from "@/lib/utils";

export function ContactList() {
  const {
    page,
    setPage,
    setPerPage,
    search,
    setSearch,
    accountId,
    setAccountId,
    role,
    setRole,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    editingContact,
    setEditingContact,
    contacts,
    pagination,
    editingContactData,
    isLoading,
    handleCreate,
    handleUpdate,
    handleDelete,
    createContact,
    updateContact,
  } = useContactList();

  const { data: accountsData } = useAccounts({ per_page: 100 });
  const accounts = accountsData?.data || [];

  const [viewingContactId, setViewingContactId] = useState<string | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);

  const handleViewContact = (contactId: string) => {
    setViewingContactId(contactId);
    setIsDetailModalOpen(true);
  };

  const getRoleBadgeVariant = (role: ContactRole | undefined) => {
    if (!role) {
      return "outline";
    }

    switch (role.code) {
      case "doctor":
        return "default";
      case "pic":
        return "secondary";
      case "manager":
        return "outline";
      default:
        return "outline";
    }
  };

  const columns: Column<Contact>[] = [
    {
      id: "name",
      header: "Name",
      accessor: (row) => (
        <button
          onClick={() => handleViewContact(row.id)}
          className="font-medium text-primary hover:underline cursor-pointer"
        >
          {row.name}
        </button>
      ),
      className: "w-[200px]",
    },
    {
      id: "role",
      header: "Role",
      accessor: (row) => (
        <Badge variant={getRoleBadgeVariant(row.role)} className="font-normal capitalize">
          {row.role?.name ?? "-"}
        </Badge>
      ),
    },
    {
      id: "position",
      header: "Position",
      accessor: (row) => <span className="text-muted-foreground">{row.position || "-"}</span>,
    },
    {
      id: "phone",
      header: "Phone",
      accessor: (row) => row.phone ? (
        <a 
          href={formatPhoneNumberToWA(row.phone)} 
          target="_blank" 
          rel="noreferrer" 
          className="text-sm text-primary hover:underline cursor-pointer"
        >
          {row.phone}
        </a>
      ) : <span className="text-muted-foreground">-</span>,
    },
    {
      id: "email",
      header: "Email",
      accessor: (row) => row.email ? (
        <a 
          href={formatEmailToMailto(row.email)} 
          title={row.email}
          className="text-sm text-primary hover:underline cursor-pointer truncate block max-w-[200px]"
        >
          {row.email}
        </a>
      ) : (
        <span className="text-muted-foreground">-</span>
      ),
    },
    {
      id: "actions",
      header: "Actions",
      accessor: (row) => (
        <div className="flex items-center justify-end gap-1">
          <Button
            variant="ghost"
            size="icon-sm"
            className="h-8 w-8"
            title="View Details"
            onClick={() => handleViewContact(row.id)}
          >
            <Eye className="h-3.5 w-3.5" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={() => setEditingContact(row.id)}
            className="h-8 w-8"
            title="Edit"
          >
            <Edit className="h-3.5 w-3.5" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={() => handleDelete(row.id)}
            className="h-8 w-8 text-destructive hover:text-destructive"
            title="Delete"
          >
            <Trash2 className="h-3.5 w-3.5" />
          </Button>
        </div>
      ),
      className: "w-[140px] text-right",
    },
  ];

  return (
    <div className="space-y-4">
      {/* Header with Actions */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3 flex-1">
          <div className="relative flex-1 max-w-sm">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search contacts..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-10 h-9"
            />
          </div>
          <Select
            value={accountId || "all"}
            onValueChange={(value) => setAccountId(value === "all" ? "" : value)}
          >
            <SelectTrigger className="h-9 w-[160px]">
              <SelectValue placeholder="All Accounts" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Accounts</SelectItem>
              {Array.isArray(accounts) && accounts.length > 0
                ? accounts.map((account) => (
                    <SelectItem key={account?.id} value={account?.id ?? ""}>
                      {account?.name ?? "-"}
                    </SelectItem>
                  ))
                : null}
            </SelectContent>
          </Select>
          <Select
            value={role || "all"}
            onValueChange={(value) => setRole(value === "all" ? "" : value)}
          >
            <SelectTrigger className="h-9 w-[140px]">
              <SelectValue placeholder="All Roles" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Roles</SelectItem>
              <SelectItem value="doctor">Doctor</SelectItem>
              <SelectItem value="pic">PIC</SelectItem>
              <SelectItem value="manager">Manager</SelectItem>
              <SelectItem value="other">Other</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <Button onClick={() => setIsCreateDialogOpen(true)} size="sm">
          <Plus className="h-4 w-4 mr-2" />
          Add Contact
        </Button>
      </div>

      {/* Table */}
      <DataTable
        columns={columns}
        data={contacts}
        isLoading={isLoading}
        emptyMessage="No contacts found"
        pagination={
          pagination
            ? {
                page: pagination.page,
                per_page: pagination.per_page,
                total: pagination.total,
                total_pages: pagination.total_pages,
                has_next: pagination.has_next,
                has_prev: pagination.has_prev,
              }
            : undefined
        }
        onPageChange={setPage}
        onPerPageChange={setPerPage}
        itemName="contact"
        perPageOptions={[10, 20, 50, 100]}
        onResetFilters={() => {
          setSearch("");
          setAccountId("");
          setRole("");
          setPage(1);
        }}
      />

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Create Contact</DialogTitle>
          </DialogHeader>
          <ContactForm
            onSubmit={async (data) => {
              await handleCreate(data as CreateContactFormData);
            }}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createContact.isPending}
            defaultAccountId={accountId || undefined}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      {editingContact && editingContactData?.data && (
        <Dialog open={!!editingContact} onOpenChange={(open) => !open && setEditingContact(null)}>
          <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>Edit Contact</DialogTitle>
            </DialogHeader>
            <ContactForm
              contact={editingContactData.data}
              onSubmit={(data) => handleUpdate(data as UpdateContactFormData)}
              onCancel={() => setEditingContact(null)}
              isLoading={updateContact.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      <ContactDetailModal
        contactId={viewingContactId}
        open={isDetailModalOpen}
        onOpenChange={(open) => {
          setIsDetailModalOpen(open);
          if (!open) {
            setViewingContactId(null);
          }
        }}
      />
    </div>
  );
}

