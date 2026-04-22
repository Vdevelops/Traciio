"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { BadgeColorSelect } from "@/components/ui/badge-color-select";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { createCategorySchema, updateCategorySchema, type CreateCategoryFormData, type UpdateCategoryFormData } from "../schemas/category.schema";
import type { Category } from "../types";
import { useTranslations } from "next-intl";

interface CategoryFormProps {
  readonly category?: Category;
  readonly onSubmit: (data: CreateCategoryFormData | UpdateCategoryFormData) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

export function CategoryForm({ category, onSubmit, onCancel, isLoading = false }: CategoryFormProps) {
  const isEdit = !!category;
  const schema = isEdit ? updateCategorySchema : createCategorySchema;
  const t = useTranslations("accountManagement.categoryForm");
  
  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<CreateCategoryFormData | UpdateCategoryFormData>({
    resolver: zodResolver(schema),
    defaultValues: category
      ? {
          name: category.name,
          code: category.code,
          description: category.description || "",
          badge_color: category.badge_color,
          status: category.status,
        }
      : {
          badge_color: "outline",
          status: "active",
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
        <FieldLabel htmlFor="badge_color">{t("badgeColorLabel")}</FieldLabel>
        <BadgeColorSelect
          value={(watch("badge_color") as any) || "outline"}
          onValueChange={(value) => setValue("badge_color", value)}
          disabled={isLoading}
          placeholder={t("badgeColorPlaceholder")}
          options={[
            { value: "default", label: t("badgeColorOptions.default") },
            { value: "secondary", label: t("badgeColorOptions.secondary") },
            { value: "destructive", label: t("badgeColorOptions.destructive") },
            { value: "outline", label: t("badgeColorOptions.outline") },
            { value: "success", label: t("badgeColorOptions.success") },
            { value: "warning", label: t("badgeColorOptions.warning") },
            { value: "active", label: t("badgeColorOptions.active") },
          ]}
        />
        {errors.badge_color && <FieldError>{errors.badge_color.message}</FieldError>}
      </Field>

      <Field>
        <div className="flex items-center space-x-3 rounded-md border p-4 bg-muted/20">
          <Checkbox
            id="status"
            checked={watch("status") === "active"}
            onCheckedChange={(checked) => setValue("status", checked ? "active" : "inactive")}
            disabled={isLoading}
          />
          <div className="space-y-1 leading-none">
            <Label htmlFor="status" className="cursor-pointer font-medium">
              {t("statusActive")}
            </Label>
            <p className="text-xs text-muted-foreground">
              {watch("status") === "active" 
                ? "This category is active and visible in account forms" 
                : "This category is inactive and will be hidden from account forms"}
            </p>
          </div>
        </div>
        {errors.status && <FieldError>{errors.status.message}</FieldError>}
      </Field>

      <div className="flex flex-col-reverse sm:flex-row justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading} className="w-full sm:w-auto">
          {t("cancel")}
        </Button>
        <Button type="submit" disabled={isLoading} className="w-full sm:w-auto">
          {isLoading ? t("submitting") : isEdit ? t("submitUpdate") : t("submitCreate")}
        </Button>
      </div>
    </form>
  );
}

