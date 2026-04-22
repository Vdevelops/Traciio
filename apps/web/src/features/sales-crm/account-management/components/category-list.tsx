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
import { useCategoryList } from "../hooks/useCategoryList";
import { CategoryForm } from "./category-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { useTranslations } from "next-intl";
import { toBadgeVariant } from "@/lib/badge-variant";
import type { CreateCategoryFormData, UpdateCategoryFormData } from "../schemas/category.schema";
import { useIsMobile } from "@/hooks/use-mobile";
import { cn } from "@/lib/utils";

export function CategoryList() {
  const {
    editingCategory,
    setEditingCategory,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    deletingCategoryId,
    setDeletingCategoryId,
    categories,
    categoryForEdit,
    isLoading,
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    deleteCategory,
    createCategory,
    updateCategory,
  } = useCategoryList();

  const t = useTranslations("accountManagement.categoryList");
  const isMobile = useIsMobile();

  return (
    <div className="space-y-4">
      {/* Header with Actions */}
      <div className="flex items-center justify-end">
        <Button onClick={() => setIsCreateDialogOpen(true)} size="sm" className="w-full sm:w-auto">
          <Plus className="h-4 w-4 mr-2" />
          {t("addCategory")}
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
                <TableHead className="w-[120px]">{t("table.badgeColor")}</TableHead>
                <TableHead className="w-[100px]">{t("table.status")}</TableHead>
                <TableHead className="w-[120px] text-right">{t("table.actions")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {categories.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center text-muted-foreground py-8">
                      {t("empty")}
                  </TableCell>
                </TableRow>
              ) : (
                categories.map((category) => (
                  <TableRow key={category.id} className="hover:bg-muted/50">
                    <TableCell className="font-medium">{category.name}</TableCell>
                    <TableCell>
                      <code className="text-xs bg-muted px-1.5 py-0.5 rounded">
                        {category.code}
                      </code>
                    </TableCell>
                    <TableCell className="text-muted-foreground max-w-xs truncate">
                      {category.description || "-"}
                    </TableCell>
                    <TableCell>
                      <Badge variant={toBadgeVariant(category.badge_color, "secondary")} className="font-normal">
                        {category.badge_color}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <StatusSwitch
                        checked={category.status === "active"}
                        onCheckedChange={(checked) => {
                          updateCategory.mutate({
                            id: category.id,
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
                          onClick={() => setEditingCategory(category.id)}
                          className="h-8 w-8"
                        >
                          <Edit className="h-3.5 w-3.5" />
                        </Button>
                        {(!category.account_count || category.account_count === 0) && (
                          <Button
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => handleDeleteClick(category.id)}
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
        <DialogContent className={cn("sm:max-w-[500px] max-h-[90vh] overflow-y-auto", isMobile && "mx-2 max-w-[calc(100vw-1rem)]")}>
          <DialogHeader>
            <DialogTitle>{t("createTitle")}</DialogTitle>
          </DialogHeader>
          <CategoryForm
            onSubmit={async (data) => {
              await handleCreate(data as CreateCategoryFormData);
            }}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createCategory.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      {editingCategory && categoryForEdit && (
        <Dialog open={!!editingCategory} onOpenChange={(open) => !open && setEditingCategory(null)}>
          <DialogContent className={cn("sm:max-w-[500px] max-h-[90vh] overflow-y-auto", isMobile && "mx-2 max-w-[calc(100vw-1rem)]")}>
            <DialogHeader>
              <DialogTitle>{t("editTitle")}</DialogTitle>
            </DialogHeader>
            <CategoryForm
              category={categoryForEdit}
              onSubmit={(data) => handleUpdate(data as UpdateCategoryFormData)}
              onCancel={() => setEditingCategory(null)}
              isLoading={updateCategory.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingCategoryId}
        onOpenChange={(open) => {
          if (!open) {
            setDeletingCategoryId(null);
          }
        }}
        onConfirm={handleDeleteConfirm}
        title={t("deleteTitle")}
        description={
          deletingCategoryId
            ? t("deleteDescriptionWithName", {
                name: categories.find((c) => c.id === deletingCategoryId)?.name || "this category",
              })
            : t("deleteDescription")
        }
        itemName="category"
        isLoading={deleteCategory.isPending}
      />
    </div>
  );
}

