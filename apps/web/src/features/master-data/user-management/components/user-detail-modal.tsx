"use client";

import { Edit, Trash2, Mail, User, Shield, Calendar, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { useUser, useDeleteUser, useUpdateUser } from "../hooks/useUsers";
import { toast } from "sonner";
import { useState } from "react";
import { UserForm } from "./user-form";
import type { User as UserType } from "../types";
import { formatEmailToMailto } from "@/lib/utils";
import { useTranslations } from "next-intl";

interface UserDetailModalProps {
  readonly userId: string | null;
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onUserUpdated?: () => void;
}

export function UserDetailModal({ userId, open, onOpenChange, onUserUpdated }: UserDetailModalProps) {
  const { data, isLoading, error } = useUser(userId || "");
  const deleteUser = useDeleteUser();
  const updateUser = useUpdateUser();
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const t = useTranslations("userManagement.detailModal");

  const user = data?.data;

  const handleDeleteConfirm = async () => {
    if (!user || !userId) return;
    try {
      await deleteUser.mutateAsync(userId);
      toast.success(t("toastDeleted"));
      onOpenChange(false);
      onUserUpdated?.();
    } catch (error) {
      // Error already handled in api-client interceptor
    }
  };

  const getAvatarUrl = (user: UserType) => {
    if (user.avatar_url) {
      return user.avatar_url;
    }
    // Always use dicebear with email as seed
    return `https://api.dicebear.com/7.x/lorelei/svg?seed=${encodeURIComponent(user.email)}`;
  };

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-[700px]">
          <DialogHeader>
            <DialogTitle>{t("title")}</DialogTitle>
          </DialogHeader>

          {isLoading && (
            <div className="space-y-6">
              <div className="flex items-center gap-4">
                <Skeleton className="h-20 w-20 rounded-full" />
                <div className="space-y-2">
                  <Skeleton className="h-6 w-48" />
                  <Skeleton className="h-4 w-32" />
                </div>
              </div>
              <Card>
                <CardHeader>
                  <Skeleton className="h-6 w-32" />
                  <Skeleton className="h-4 w-64 mt-2" />
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-4 w-3/4" />
                    <Skeleton className="h-4 w-1/2" />
                  </div>
                </CardContent>
              </Card>
            </div>
          )}

          {error && (
            <div className="text-center text-muted-foreground py-8">
              {t("loadError")}
            </div>
          )}

          {!isLoading && !error && user && (
            <div className="space-y-6">
              {/* User Header */}
              <div className="flex items-center gap-4 pb-4 border-b">
                <Avatar className="h-20 w-20">
                  <AvatarImage src={getAvatarUrl(user)} alt={user.name} />
                </Avatar>
                <div className="flex-1">
                  <h2 className="text-2xl font-medium tracking-tight">{user.name}</h2>
                  <a href={formatEmailToMailto(user.email)} className="text-sm text-muted-foreground mt-1 hover:text-primary hover:underline cursor-pointer">{user.email}</a>
                </div>
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setIsEditDialogOpen(true)}
                    className="shrink-0"
                  >
                    <Edit className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="destructive"
                    size="sm"
                    onClick={() => setIsDeleteDialogOpen(true)}
                    disabled={deleteUser.isPending}
                    className="shrink-0"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              {/* User Info Card */}
              <Card>
                <CardHeader>
                  <CardTitle>{t("userInfo.title")}</CardTitle>
                  <CardDescription>
                    {t("userInfo.description")}
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <User className="h-4 w-4" />
                        <span>{t("userInfo.name")}</span>
                      </div>
                      <div className="text-base font-medium">{user.name}</div>
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Mail className="h-4 w-4" />
                        <span>{t("userInfo.email")}</span>
                      </div>
                      <a href={formatEmailToMailto(user.email)} className="text-base font-medium text-primary hover:underline cursor-pointer">{user.email}</a>
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Shield className="h-4 w-4" />
                        <span>{t("userInfo.role")}</span>
                      </div>
                      <div>
                        <Badge variant="outline" className="font-normal">
                          {user.role?.name || "N/A"}
                        </Badge>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <span>{t("userInfo.group")}</span>
                      </div>
                      <div>
                        <Badge variant="outline" className="font-normal">
                          {user.group?.name || "N/A"}
                        </Badge>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <span>{t("userInfo.brick")}</span>
                      </div>
                      <div>
                        {user.brick ? (
                          <Badge variant="outline" className="font-normal">
                            {user.brick.name || user.brick.code || "N/A"}
                          </Badge>
                        ) : (
                          <span className="text-base font-medium text-muted-foreground">{t("userInfo.notAvailable")}</span>
                        )}
                      </div>
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <span>{t("userInfo.monthlyTarget")}</span>
                      </div>
                      <div className="text-base font-medium">
                        {user.monthly_target
                          ? user.monthly_target.target_amount_formatted ||
                            new Intl.NumberFormat("id-ID", {
                              style: "currency",
                              currency: "IDR",
                              minimumFractionDigits: 0,
                            }).format(user.monthly_target.target_amount)
                          : t("userInfo.notAvailable")}
                      </div>
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <span>{t("userInfo.status")}</span>
                      </div>
                      <div>
                        <Badge variant={user.status === "active" ? "active" : "inactive"}>
                          {user.status}
                        </Badge>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Calendar className="h-4 w-4" />
                        <span>{t("userInfo.createdAt")}</span>
                      </div>
                      <div className="text-base font-medium">
                        {user.created_at
                          ? new Date(user.created_at).toLocaleDateString()
                          : t("userInfo.notAvailable")}
                      </div>
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Calendar className="h-4 w-4" />
                        <span>{t("userInfo.updatedAt")}</span>
                      </div>
                      <div className="text-base font-medium">
                        {user.updated_at
                          ? new Date(user.updated_at).toLocaleDateString()
                          : t("userInfo.notAvailable")}
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Role Information */}
              {user.role && (
                <Card>
                  <CardHeader>
                    <CardTitle>{t("roleInfo.title")}</CardTitle>
                    <CardDescription>
                      {t("roleInfo.description")}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Shield className="h-4 w-4" />
                        <span>{t("roleInfo.name")}</span>
                      </div>
                      <div className="text-base font-medium">{user.role.name}</div>
                    </div>

                    {user.role.description && (
                      <div className="space-y-2">
                        <div className="text-sm text-muted-foreground">
                          {t("roleInfo.descriptionLabel")}
                        </div>
                        <div className="text-base">{user.role.description}</div>
                      </div>
                    )}

                    <div className="space-y-2">
                      <div className="text-sm text-muted-foreground">
                        {t("roleInfo.code")}
                      </div>
                      <Badge variant="outline" className="font-mono">
                        {user.role.code}
                      </Badge>
                    </div>

                    {user.role.permissions && user.role.permissions.length > 0 && (
                      <div className="space-y-3">
                        <div className="text-sm font-medium">
                          {t("roleInfo.permissions", {
                            count: user.role.permissions.length,
                          })}
                        </div>
                        <div className="flex flex-wrap gap-2">
                          {user.role.permissions.map((permission) => (
                            <Badge key={permission.id} variant="outline" className="font-normal">
                              {permission.name}
                            </Badge>
                          ))}
                        </div>
                      </div>
                    )}

                    {(!user.role.permissions || user.role.permissions.length === 0) && (
                      <div className="text-sm text-muted-foreground">
                        {t("roleInfo.noPermissions")}
                      </div>
                    )}
                  </CardContent>
                </Card>
              )}
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      {isEditDialogOpen && user && (
        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>{t("editDialogTitle")}</DialogTitle>
            </DialogHeader>
            <UserForm
              user={user}
              onSubmit={async (formData) => {
                try {
                  await updateUser.mutateAsync({ id: userId!, data: formData });
                  setIsEditDialogOpen(false);
                  toast.success(t("toastUpdated"));
                  onUserUpdated?.();
                } catch (error) {
                  // Error already handled in api-client interceptor
                }
              }}
              onCancel={() => setIsEditDialogOpen(false)}
              isLoading={updateUser.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Delete Dialog */}
      <DeleteDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        onConfirm={handleDeleteConfirm}
        title={t("deleteDialogTitle")}
        description={
          user
            ? t("deleteDialogDescriptionWithName", { name: user.name })
            : t("deleteDialogDescription")
        }
        itemName={t("deleteDialogItemName")}
        isLoading={deleteUser.isPending}
      />
    </>
  );
}

