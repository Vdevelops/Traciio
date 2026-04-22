"use client";

import { useRouter } from "next/navigation";
import { ArrowLeft, Edit, Trash2, Mail, User, Shield, Calendar } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useUser, useDeleteUser, useUserPermissions, useUpdateUser } from "../hooks/useUsers";
import { useHasPermission } from "../hooks/useHasPermission";
import { toast } from "sonner";
import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { UserForm } from "./user-form";
import { formatEmailToMailto } from "@/lib/utils";
import { useTranslations } from "next-intl";

interface UserDetailProps {
  readonly userId: string;
}

export function UserDetail({ userId }: UserDetailProps) {
  const router = useRouter();
  const { data, isLoading, error } = useUser(userId);
  const deleteUser = useDeleteUser();
  const updateUser = useUpdateUser();
  const { data: permissionsData } = useUserPermissions(userId);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isPermissionsDialogOpen, setIsPermissionsDialogOpen] = useState(false);
  const t = useTranslations("userManagement.detailPage");

  const user = data?.data;
  const userPermissions = permissionsData?.data || [];

  // Permission checks
  const hasEditPermission = useHasPermission("users.edit");
  const hasPermissionsPermission = useHasPermission("users.permissions");
  const hasDeletePermission = useHasPermission("users.delete");

  const handleDelete = async () => {
    if (!user) return;
    if (confirm(t("list.deleteDescriptionWithName", { name: user.name }))) {
      try {
        await deleteUser.mutateAsync(userId);
        toast.success(t("toastDeleted"));
        router.push("/master-data/users");
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Skeleton className="h-10 w-10" />
          <Skeleton className="h-8 w-48" />
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
    );
  }

  if (error || !user) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => router.back()}>
          <ArrowLeft className="h-4 w-4 mr-2" />
          {t("back")}
        </Button>
        <Card>
          <CardContent className="pt-6">
            <div className="text-center text-muted-foreground">
              {t("notFound")}
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={() => router.back()}>
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-2xl font-medium tracking-tight">{user?.name || "Unknown User"}</h1>
            <p className="text-sm text-muted-foreground mt-1">
              {t("headerSubtitle")}
            </p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          {hasEditPermission && (
            <Button
              variant="outline"
              onClick={() => setIsEditDialogOpen(true)}
            >
              <Edit className="h-4 w-4 mr-2" />
              {t("actions.edit")}
            </Button>
          )}
          {hasPermissionsPermission && (
            <Button
              variant="outline"
              onClick={() => setIsPermissionsDialogOpen(true)}
            >
              <Shield className="h-4 w-4 mr-2" />
              {t("actions.permissions")}
            </Button>
          )}
          {hasDeletePermission && (
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteUser.isPending}
            >
              <Trash2 className="h-4 w-4 mr-2" />
              {t("actions.delete")}
            </Button>
          )}
        </div>
      </div>

      {/* User Info Card */}
      <Card>
        <CardHeader>
          <CardTitle>{t("userInfo.title")}</CardTitle>
          <CardDescription>{t("userInfo.description")}</CardDescription>
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
                  <Badge variant="outline" className="font-normal cursor-pointer" onClick={() => router.push(`/master-data/bricks/${user.brick?.id ?? ""}`)}>
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

      {/* Permissions Card */}
      {userPermissions.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>{t("permissionsCard.title")}</CardTitle>
            <CardDescription>
              {t("permissionsCard.description")}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-2">
              {userPermissions.map((permission) => (
                <Badge key={permission} variant="outline" className="font-normal">
                  {permission}
                </Badge>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Edit Dialog */}
      {isEditDialogOpen && (
        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>{t("actions.edit")}</DialogTitle>
            </DialogHeader>
            <UserForm
              user={user}
              onSubmit={async (formData) => {
                try {
                  await updateUser.mutateAsync({ id: userId, data: formData });
                  setIsEditDialogOpen(false);
                  toast.success(t("toastUpdated"));
                  // Refresh the page to show updated data
                  window.location.reload();
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

      {/* Permissions Dialog - TODO: Implement user permissions assignment dialog */}
      {isPermissionsDialogOpen && (
        <Dialog open={isPermissionsDialogOpen} onOpenChange={setIsPermissionsDialogOpen}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle>{t("permissionsDialog.title")}</DialogTitle>
            </DialogHeader>
            <div className="text-sm text-muted-foreground">
              {t("permissionsDialog.description")}
            </div>
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}

