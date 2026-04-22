"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useDebounce } from "@/hooks/use-debounce";
import {
  useContacts,
  useDeleteContact,
  useContact,
  useCreateContact,
  useUpdateContact,
} from "./useContacts";
import type { CreateContactFormData, UpdateContactFormData } from "../schemas/contact.schema";

export function useContactList() {
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(20);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);

  useEffect(() => {
    setPage(1);
  }, [debouncedSearch]);
  const [accountId, setAccountId] = useState<string>("");
  const [role, setRole] = useState<string>("");
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingContact, setEditingContact] = useState<string | null>(null);
  const [deletingContactId, setDeletingContactId] = useState<string | null>(null);

  const { data, isLoading } = useContacts({
    page,
    per_page: perPage,
    search: debouncedSearch,
    account_id: accountId,
    role_id: role,
  });
  const { data: editingContactData } = useContact(editingContact || "");
  const deleteContact = useDeleteContact();
  const createContact = useCreateContact();
  const updateContact = useUpdateContact();

  const contacts = data?.data || [];
  const pagination = data?.meta?.pagination;

  const handleCreate = async (formData: CreateContactFormData) => {
    try {
      await createContact.mutateAsync(formData);
      setIsCreateDialogOpen(false);
      toast.success("Contact created successfully");
    } catch (error) {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: UpdateContactFormData) => {
    if (editingContact) {
      try {
        await updateContact.mutateAsync({ id: editingContact, data: formData });
        setEditingContact(null);
        toast.success("Contact updated successfully");
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteContact.mutateAsync(id);
      toast.success("Contact deleted successfully");
    } catch (error) {
      // Error already handled in api-client interceptor
    }
  };

  const handlePerPageChange = (newPerPage: number) => {
    setPerPage(newPerPage);
    setPage(1);
  };

  return {
    // State
    page,
    setPage,
    perPage,
    setPerPage: handlePerPageChange,
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
    deletingContactId,
    setDeletingContactId,
    // Data
    contacts,
    pagination,
    editingContactData,
    isLoading,
    // Actions
    handleCreate,
    handleUpdate,
    handleDelete,
    // Mutations
    deleteContact,
    createContact,
    updateContact,
  };
}

