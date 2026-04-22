"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslations } from "next-intl";
import { createCategorySchema, updateCategorySchema, type CreateCategoryFormData, type UpdateCategoryFormData } from "../schemas/category.schema";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import type { ProductCategory } from "../types";

interface CategoryFormProps {
  readonly category?: ProductCategory;
  readonly onSubmit: (data: CreateCategoryFormData | UpdateCategoryFormData) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

export function CategoryForm({ category, onSubmit, onCancel, isLoading }: CategoryFormProps) {
  const isEdit = !!category;
  const t = useTranslations("productManagement.categoryForm");

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreateCategoryFormData | UpdateCategoryFormData>({
    resolver: zodResolver(isEdit ? updateCategorySchema : createCategorySchema),
    defaultValues: category
      ? {
          name: category.name,
          slug: category.slug || "",
          description: category.description || "",
          status: category.status,
        }
      : {
          status: "active",
        },
  });

  const selectedStatus = watch("status");

  const handleFormSubmit = async (data: CreateCategoryFormData | UpdateCategoryFormData) => {
    await onSubmit(data);
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <Field orientation="vertical">
        <FieldLabel>{t("nameLabel")}</FieldLabel>
        <Input
          {...register("name")}
          placeholder={t("namePlaceholder")}
        />
        {errors.name && <FieldError>{errors.name.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("slugLabel")}</FieldLabel>
        <Input
          {...register("slug")}
          placeholder={t("slugPlaceholder")}
        />
        <p className="text-xs text-muted-foreground mt-1">
          {t("slugHint")}
        </p>
        {errors.slug && <FieldError>{errors.slug.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("descriptionLabel")}</FieldLabel>
        <Textarea
          {...register("description")}
          placeholder={t("descriptionPlaceholder")}
          rows={3}
        />
        {errors.description && <FieldError>{errors.description.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <div className="flex items-center space-x-3 rounded-md border p-4 bg-muted/20">
          <Checkbox
            id="status"
            checked={selectedStatus === "active"}
            onCheckedChange={(checked) => setValue("status", checked ? "active" : "inactive", { shouldValidate: true })}
            disabled={isLoading}
          />
          <div className="space-y-1 leading-none">
            <Label htmlFor="status" className="cursor-pointer font-medium">
              {t("statusActive")}
            </Label>
            <p className="text-xs text-muted-foreground">
              {selectedStatus === "active" 
                ? t("statusActiveDescription") 
                : t("statusInactiveDescription")}
            </p>
          </div>
        </div>
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

