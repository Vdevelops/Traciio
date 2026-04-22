"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslations } from "next-intl";
import {
  createGroupSchema,
  updateGroupSchema,
  type CreateGroupFormData,
  type UpdateGroupFormData,
} from "../schemas/group.schema";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { Group } from "../types";

interface GroupFormProps {
  readonly group?: Group;
  readonly onSubmit: (
    data: CreateGroupFormData | UpdateGroupFormData
  ) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

type GroupFormData = CreateGroupFormData | UpdateGroupFormData;

export function GroupForm({
  group,
  onSubmit,
  onCancel,
  isLoading,
}: GroupFormProps) {
  const isEdit = !!group;
  const t = useTranslations("groupManagement.form");

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<GroupFormData>({
    resolver: zodResolver(isEdit ? updateGroupSchema : createGroupSchema),
    defaultValues: group
      ? {
          name: group.name,
          code: group.code,
          description: group.description || "",
          status: group.status,
        }
      : {
          status: "active",
        },
  });

  const handleFormSubmit = async (data: GroupFormData) => {
    await onSubmit(data);
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <Field orientation="vertical">
        <FieldLabel>{t("nameLabel")}</FieldLabel>
        <Input
          {...register("name")}
          placeholder={t("namePlaceholder")}
          disabled={isLoading}
        />
        {errors.name && <FieldError>{errors.name.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("codeLabel")}</FieldLabel>
        <Input
          {...register("code")}
          placeholder={t("codePlaceholder")}
          disabled={isLoading || isEdit}
        />
        {errors.code && <FieldError>{errors.code.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("descriptionLabel")}</FieldLabel>
        <Textarea
          {...register("description")}
          placeholder={t("descriptionPlaceholder")}
          disabled={isLoading}
          rows={3}
        />
        {errors.description && (
          <FieldError>{errors.description.message}</FieldError>
        )}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("statusLabel")}</FieldLabel>
        <Select
          value={watch("status") || "active"}
          onValueChange={(value) =>
            setValue("status", value as "active" | "inactive")
          }
          disabled={isLoading}
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
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isLoading}
          className="cursor-pointer"
        >
          {t("cancelButton")}
        </Button>
        <Button type="submit" disabled={isLoading} className="cursor-pointer">
          {isLoading
            ? t("savingButton")
            : isEdit
              ? t("updateButton")
              : t("createButton")}
        </Button>
      </div>
    </form>
  );
}

