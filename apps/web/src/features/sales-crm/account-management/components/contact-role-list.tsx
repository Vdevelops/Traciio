"use client";

import { Edit, Trash2, Plus } from "lucide-react";
import { StatusSwitch } from "@/components/ui/status-switch";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { useContactRoleList } from "../hooks/useContactRoleList";
import { ContactRoleForm } from "./contact-role-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { useTranslations } from "next-intl";
import { toBadgeVariant } from "@/lib/badge-variant";
import type {
  CreateContactRoleFormData,
  UpdateContactRoleFormData,
} from "../schemas/contact-role.schema";
import { useIsMobile } from "@/hooks/use-mobile";
import { cn } from "@/lib/utils";

export function ContactRoleList() {
  const {
    editingContactRole,
    setEditingContactRole,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    deletingContactRoleId,
    setDeletingContactRoleId,
    contactRoles,
    contactRoleForEdit,
    isLoading,
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    deleteContactRole,
    createContactRole,
    updateContactRole,
  } = useContactRoleList();

  const t = useTranslations("accountManagement.contactRoleList");
  const isMobile = useIsMobile();

  return (
    <div className="space-y-4">
      {/* Header with Actions */}
      <div className="flex items-center justify-end">
        <Button
          onClick={() => setIsCreateDialogOpen(true)}
          size="sm"
          className="w-full sm:w-auto"
        >
          <Plus className="h-4 w-4 mr-2" />
          {t("addContactRole")}
        </Button>
      </div>

      {/* Table */}
      <div className="border rounded-lg">
        {isLoading ? (
          <div className="p-4 space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={`skeleton-${i}`} className="h-10 w-full" />
            ))}
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[200px]">{t("table.name")}</TableHead>
                <TableHead>{t("table.code")}</TableHead>
                <TableHead>{t("table.description")}</TableHead>
                <TableHead className="w-[120px]">
                  {t("table.badgeColor")}
                </TableHead>
                <TableHead className="w-[100px]">{t("table.status")}</TableHead>
                <TableHead className="w-[120px] text-right">
                  {t("table.actions")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {contactRoles.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={6}
                    className="text-center text-muted-foreground py-8"
                  >
                    {t("empty")}
                  </TableCell>
                </TableRow>
              ) : (
                contactRoles.map((contactRole) => (
                  <TableRow key={contactRole.id} className="hover:bg-muted/50">
                    <TableCell className="font-medium">
                      {contactRole.name}
                    </TableCell>
                    <TableCell>
                      <code className="text-xs bg-muted px-1.5 py-0.5 rounded">
                        {contactRole.code}
                      </code>
                    </TableCell>
                    <TableCell className="text-muted-foreground max-w-xs truncate">
                      {contactRole.description || "-"}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={toBadgeVariant(
                          contactRole.badge_color,
                          "secondary",
                        )}
                        className="font-normal"
                      >
                        {contactRole.badge_color}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <StatusSwitch
                        checked={contactRole.status === "active"}
                        onCheckedChange={(checked) => {
                          return updateContactRole.mutateAsync({
                            id: contactRole.id,
                            data: { status: checked ? "active" : "inactive" },
                          });
                        }}
                      />
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          onClick={() => setEditingContactRole(contactRole.id)}
                          className="h-8 w-8"
                        >
                          <Edit className="h-3.5 w-3.5" />
                        </Button>
                        {(!contactRole.contact_count ||
                          contactRole.contact_count === 0) && (
                          <Button
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => handleDeleteClick(contactRole.id)}
                            className="h-8 w-8 text-destructive hover:text-destructive"
                          >
                            <Trash2 className="h-3.5 w-3.5" />
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        )}
      </div>

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent
          className={cn(
            "sm:max-w-[500px] max-h-[90vh] overflow-y-auto",
            isMobile && "mx-2 max-w-[calc(100vw-1rem)]",
          )}
        >
          <DialogHeader>
            <DialogTitle>{t("createTitle")}</DialogTitle>
          </DialogHeader>
          <ContactRoleForm
            onSubmit={async (data) => {
              await handleCreate(data as CreateContactRoleFormData);
            }}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createContactRole.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      {editingContactRole && contactRoleForEdit && (
        <Dialog
          open={!!editingContactRole}
          onOpenChange={(open) => !open && setEditingContactRole(null)}
        >
          <DialogContent
            className={cn(
              "sm:max-w-[500px] max-h-[90vh] overflow-y-auto",
              isMobile && "mx-2 max-w-[calc(100vw-1rem)]",
            )}
          >
            <DialogHeader>
              <DialogTitle>{t("editTitle")}</DialogTitle>
            </DialogHeader>
            <ContactRoleForm
              contactRole={contactRoleForEdit}
              onSubmit={(data) =>
                handleUpdate(data as UpdateContactRoleFormData)
              }
              onCancel={() => setEditingContactRole(null)}
              isLoading={updateContactRole.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingContactRoleId}
        onOpenChange={(open) => {
          if (!open) {
            setDeletingContactRoleId(null);
          }
        }}
        onConfirm={handleDeleteConfirm}
        title={t("deleteTitle")}
        description={
          deletingContactRoleId
            ? t("deleteDescriptionWithName", {
                name:
                  contactRoles.find((r) => r.id === deletingContactRoleId)
                    ?.name || "this contact role",
              })
            : t("deleteDescription")
        }
        itemName={t("deleteItemName")}
        isLoading={deleteContactRole.isPending}
      />
    </div>
  );
}
