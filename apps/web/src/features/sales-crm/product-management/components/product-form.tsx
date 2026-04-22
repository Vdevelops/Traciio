"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslations } from "next-intl";
import { createProductSchema, updateProductSchema, type CreateProductFormData, type UpdateProductFormData } from "../schemas/product.schema";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { NumberInput } from "@/components/ui/number-input";
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

import { useProductCategories } from "../hooks/useProducts";
import type { Product } from "../types";
import { ImageUploadField } from "./image-upload-field";

interface ProductFormProps {
  readonly product?: Product;
  readonly onSubmit: (data: CreateProductFormData | UpdateProductFormData) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

export function ProductForm({ product, onSubmit, onCancel, isLoading }: ProductFormProps) {
  const isEdit = !!product;
  const { data: categoriesData } = useProductCategories({ status: "active" });
  const categories = categoriesData?.data || [];
  const t = useTranslations("productManagement.form");

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreateProductFormData | UpdateProductFormData>({
    resolver: zodResolver(isEdit ? updateProductSchema : createProductSchema),
    defaultValues: product
      ? {
          name: product.name,
          sku: product.sku,
          barcode: product.barcode || "",
          price: product.price / 100, // Convert from sen to rupiah
          cost: product.cost / 100,
          category_id: product.category_id,
          description: product.description || "",
          status: product.status || "active",
          image_url: product.image_url || "",
        }
      : {
          cost: 0,
          status: "active",
          image_url: "",
        },
  });

  const selectedCategoryId = watch("category_id");
  const selectedStatus = watch("status");
  const imageUrl = watch("image_url");

  const handleFormSubmit = async (data: CreateProductFormData | UpdateProductFormData) => {
    // Convert price and cost from rupiah to sen
    const submitData: CreateProductFormData | UpdateProductFormData = {
      ...data,
      price: Math.round((data.price ?? 0) * 100),
      cost: data.cost ? Math.round(data.cost * 100) : undefined,
    };
    await onSubmit(submitData);
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
        <FieldLabel>{t("skuLabel")}</FieldLabel>
        <Input
          {...register("sku")}
          placeholder={t("skuPlaceholder")}
          disabled={isEdit}
        />
        {errors.sku && <FieldError>{errors.sku.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("barcodeLabel")}</FieldLabel>
        <Input
          {...register("barcode")}
          placeholder={t("barcodePlaceholder")}
        />
        {errors.barcode && <FieldError>{errors.barcode.message}</FieldError>}
      </Field>

      <div className="grid grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("priceLabel")}</FieldLabel>
          <NumberInput
            value={watch("price")}
            onChange={(value) => setValue("price", value ?? 0, { shouldValidate: true })}
            placeholder={t("pricePlaceholder")}
            allowDecimal
            decimalPlaces={2}
          />
          {errors.price && <FieldError>{errors.price.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("costLabel")}</FieldLabel>
          <NumberInput
            value={watch("cost")}
            onChange={(value) => setValue("cost", value ?? 0, { shouldValidate: true })}
            placeholder={t("costPlaceholder")}
            allowDecimal
            decimalPlaces={2}
          />
          {errors.cost && <FieldError>{errors.cost.message}</FieldError>}
        </Field>
      </div>

      <Field orientation="vertical">
        <FieldLabel>{t("categoryLabel")}</FieldLabel>
        <Select
          value={selectedCategoryId || ""}
          onValueChange={(value) => setValue("category_id", value)}
        >
          <SelectTrigger>
            <SelectValue placeholder={t("categoryPlaceholder")} />
          </SelectTrigger>
          <SelectContent>
            {categories.map((category) => (
              <SelectItem key={category.id} value={category.id}>
                {category.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {errors.category_id && <FieldError>{errors.category_id.message}</FieldError>}
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

      <ImageUploadField
        value={imageUrl}
        onChange={(url) => setValue("image_url", url, { shouldValidate: true })}
        disabled={isLoading}
        error={errors.image_url?.message}
      />

      <Field orientation="vertical">
        <div className="flex items-center space-x-3 rounded-md border p-4 bg-muted/20 mt-2">
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
