"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslations } from "next-intl";
import { createUserSchema, updateUserSchema, type CreateUserFormData, type UpdateUserFormData } from "../schemas/user.schema";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useRoles, useGroups } from "../hooks/useUsers";
import { useBricks } from "@/features/master-data/brick/hooks/useBricks";
import type { User } from "../types";

interface UserFormProps {
  readonly user?: User;
  readonly onSubmit: (data: CreateUserFormData | UpdateUserFormData) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

export function UserForm({ user, onSubmit, onCancel, isLoading }: UserFormProps) {
  const isEdit = !!user;
  const { data: rolesData } = useRoles();
  const { data: groupsData } = useGroups();
  const { data: bricksData } = useBricks({ per_page: 100, status: "active" });
  const roles = rolesData?.data || [];
  const groups = groupsData?.data || [];
  const bricks = bricksData?.data || [];
  const t = useTranslations("userManagement.form");

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreateUserFormData | UpdateUserFormData>({
    resolver: zodResolver(isEdit ? updateUserSchema : createUserSchema),
    defaultValues: user
      ? {
          email: user.email,
          name: user.name,
          role_id: user.role_id || undefined,
          group_id: user.group_id || null,
          brick_id: user.brick_id || null,
          status: user.status,
        }
      : {
          status: "active",
        },
  });

  const selectedRoleId = watch("role_id");
  const selectedGroupId = watch("group_id") || "none";
  const selectedBrickId = watch("brick_id") || "none";

  const handleFormSubmit = async (data: CreateUserFormData | UpdateUserFormData) => {
    await onSubmit(data);
  };

  const handleRoleChange = (value: string) => {
    setValue("role_id", value);
  };

  const handleGroupChange = (value: string) => {
    setValue("group_id", value === "none" ? null : value);
  };

  const handleBrickChange = (value: string) => {
    setValue("brick_id", value === "none" ? null : value);
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <Field orientation="vertical">
        <FieldLabel>{t("emailLabel")}</FieldLabel>
        <Input
          type="email"
          {...register("email")}
          disabled={isEdit}
          placeholder={t("emailPlaceholder")}
        />
        {errors.email && <FieldError>{errors.email.message}</FieldError>}
      </Field>

      {!isEdit && (
        <Field orientation="vertical">
          <FieldLabel>{t("passwordLabel")}</FieldLabel>
          <Input
            type="password"
            {...register("password")}
            placeholder={t("passwordPlaceholder")}
          />
          {(
            errors as {
              password?: { message?: string };
            }
          ).password && (
            <FieldError>
              {(
                errors as {
                  password?: { message?: string };
                }
              ).password?.message}
            </FieldError>
          )}
        </Field>
      )}

      <Field orientation="vertical">
        <FieldLabel>{t("nameLabel")}</FieldLabel>
        <Input
          {...register("name")}
          placeholder={t("namePlaceholder")}
        />
        {errors.name && <FieldError>{errors.name.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("roleLabel")}</FieldLabel>
        <Select
          value={selectedRoleId || ""}
          onValueChange={handleRoleChange}
        >
          <SelectTrigger>
            <SelectValue placeholder={t("rolePlaceholder")} />
          </SelectTrigger>
          <SelectContent>
            {roles.map((role) => (
              <SelectItem key={role.id} value={role.id}>
                {role.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {errors.role_id && <FieldError>{errors.role_id.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("groupLabel")}</FieldLabel>
        <Select
          value={selectedGroupId}
          onValueChange={handleGroupChange}
        >
          <SelectTrigger>
            <SelectValue placeholder={t("groupPlaceholder")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="none">{t("noGroup")}</SelectItem>
            {groups.map((group) => (
              <SelectItem key={group.id} value={group.id}>
                {group.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {("group_id" in errors && errors.group_id) && (
          <FieldError>
            {typeof errors.group_id === "object" && "message" in errors.group_id
              ? (errors.group_id.message as string)
              : "Invalid group selection"}
          </FieldError>
        )}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("brickLabel")}</FieldLabel>
        <Select
          value={selectedBrickId}
          onValueChange={handleBrickChange}
        >
          <SelectTrigger>
            <SelectValue placeholder={t("brickPlaceholder")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="none">{t("noBrick")}</SelectItem>
            {bricks.map((brick) => (
              <SelectItem key={brick.id} value={brick.id}>
                {brick.name} ({brick.province} - {brick.regency})
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {("brick_id" in errors && errors.brick_id) && (
          <FieldError>
            {typeof errors.brick_id === "object" && "message" in errors.brick_id
              ? (errors.brick_id.message as string)
              : "Invalid brick selection"}
          </FieldError>
        )}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("statusLabel")}</FieldLabel>
        <Select
          value={watch("status") || "active"}
          onValueChange={(value) => setValue("status", value as "active" | "inactive")}
        >
          <SelectTrigger>
            <SelectValue placeholder={t("statusLabel")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="active">{t("statusActive")}</SelectItem>
            <SelectItem value="inactive">{t("statusInactive")}</SelectItem>
          </SelectContent>
        </Select>
        {errors.status && <FieldError>{errors.status.message}</FieldError>}
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

