"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslations } from "next-intl";
import { Smartphone } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { createRoleSchema, updateRoleSchema, type CreateRoleFormData, type UpdateRoleFormData } from "../schemas/role.schema";
import type { Role } from "../types";

interface RoleFormProps {
  readonly role?: Role;
  readonly onSubmit: (data: CreateRoleFormData | UpdateRoleFormData) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

export function RoleForm({ role, onSubmit, onCancel, isLoading = false }: RoleFormProps) {
  const isEdit = !!role;
  const schema = isEdit ? updateRoleSchema : createRoleSchema;
  const t = useTranslations("userManagement.roleForm");
  
  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<CreateRoleFormData | UpdateRoleFormData>({
    resolver: zodResolver(schema),
    defaultValues: role
      ? {
          name: role.name,
          code: role.code,
          description: role.description || "",
          status: role.status,
          mobile_access: role.mobile_access ?? false,
        }
      : {
          status: "active",
          mobile_access: false,
        },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <Field>
        <FieldLabel htmlFor="name">{t("nameLabel")}</FieldLabel>
        <Input
          id="name"
          {...register("name")}
          placeholder={t("namePlaceholder")}
          disabled={isLoading}
        />
        {errors.name && <FieldError>{errors.name.message}</FieldError>}
      </Field>

      <Field>
        <FieldLabel htmlFor="code">{t("codeLabel")}</FieldLabel>
        <Input
          id="code"
          {...register("code")}
          placeholder={t("codePlaceholder")}
          disabled={isLoading || isEdit}
        />
        {isEdit && (
          <p className="text-xs text-muted-foreground mt-1">
            {t("codeHint")}
          </p>
        )}
        {errors.code && <FieldError>{errors.code.message}</FieldError>}
      </Field>

      <Field>
        <FieldLabel htmlFor="description">{t("descriptionLabel")}</FieldLabel>
        <Input
          id="description"
          {...register("description")}
          placeholder={t("descriptionPlaceholder")}
          disabled={isLoading}
        />
        {errors.description && <FieldError>{errors.description.message}</FieldError>}
      </Field>

      <Field>
        <FieldLabel htmlFor="status">{t("statusLabel")}</FieldLabel>
        <Select
          value={watch("status") || "active"}
          onValueChange={(value) => setValue("status", value as "active" | "inactive")}
          disabled={isLoading}
        >
          <SelectTrigger className="h-10">
            <SelectValue placeholder={t("statusLabel")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="active">{t("statusActive")}</SelectItem>
            <SelectItem value="inactive">{t("statusInactive")}</SelectItem>
          </SelectContent>
        </Select>
        {errors.status && <FieldError>{errors.status.message}</FieldError>}
      </Field>

      <Field>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Smartphone className="h-4 w-4 text-muted-foreground" />
            <FieldLabel htmlFor="mobile_access" className="cursor-pointer">
              {t("mobileAccessLabel") || "Mobile Access"}
            </FieldLabel>
          </div>
          <Switch
            id="mobile_access"
            checked={watch("mobile_access") ?? false}
            onCheckedChange={(checked) => setValue("mobile_access", checked)}
            disabled={isLoading}
          />
        </div>
        <p className="text-xs text-muted-foreground mt-1">
          {t("mobileAccessHint") || "Allow users with this role to access mobile app"}
        </p>
        {errors.mobile_access && <FieldError>{errors.mobile_access.message}</FieldError>}
      </Field>

      <div className="flex justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading}>
          {t("cancel")}
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading
            ? t("submitting")
            : isEdit
              ? t("submitUpdate")
              : t("submitCreate")}
        </Button>
      </div>
    </form>
  );
}

