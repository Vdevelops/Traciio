"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import {
  createAccountSchema,
  updateAccountSchema,
  type CreateAccountFormData,
  type UpdateAccountFormData,
} from "../schemas/account.schema";
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
import { useCategories } from "../hooks/useCategories";
import type { Account } from "../types";
import { useTranslations } from "next-intl";
import dynamic from "next/dynamic";

// Dynamic import to avoid Leaflet SSR errors (leaflet accesses `window` at module level)
const AccountLocationPickerDialog = dynamic(
  () => import("./account-location-picker-dialog").then((mod) => mod.AccountLocationPickerDialog),
  { ssr: false }
);

interface AccountFormProps {
  readonly account?: Account;
  readonly onSubmit: (data: CreateAccountFormData | UpdateAccountFormData) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

export function AccountForm({ account, onSubmit, onCancel, isLoading }: AccountFormProps) {
  const isEdit = !!account;
  const { data: categoriesData } = useCategories();
  const categories = categoriesData?.data || [];
  const t = useTranslations("accountManagement.accountForm");
  const [isPickerOpen, setIsPickerOpen] = useState(false);

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreateAccountFormData | UpdateAccountFormData>({
    resolver: zodResolver(isEdit ? updateAccountSchema : createAccountSchema),
    defaultValues: account
      ? {
          name: account.name,
          category_id: account.category_id,
          address: account.address || "",
          city: account.city || "",
          province: account.province || "",
          phone: account.phone || "",
          email: account.email || "",
          latitude: account.latitude,
          longitude: account.longitude,
          status: account.status,
          assigned_to: account.assigned_to || "",
        }
      : {
          status: "active",
          province: "Indonesia",
        },
  });

  const handleFormSubmit = async (data: CreateAccountFormData | UpdateAccountFormData) => {
    await onSubmit(data);
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <Field orientation="vertical">
        <FieldLabel>{t("nameLabel")}{""}*</FieldLabel>
        <Input {...register("name")} placeholder={t("namePlaceholder")} />
        {errors.name && <FieldError>{errors.name.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("categoryLabel")}{""}*</FieldLabel>
        <Select
          value={watch("category_id") || ""}
          onValueChange={(value) => setValue("category_id", value)}
        >
          <SelectTrigger>
            <SelectValue placeholder={t("categoryPlaceholder")} />
          </SelectTrigger>
          <SelectContent>
            {categories
              .filter((cat) => cat.status === "active")
              .map((category) => (
                <SelectItem key={category.id} value={category.id}>
                  {category.name}
                </SelectItem>
              ))}
          </SelectContent>
        </Select>
        {errors.category_id && <FieldError>{errors.category_id.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("addressLabel")}</FieldLabel>
        <Textarea {...register("address")} placeholder={t("addressPlaceholder")} rows={3} />
        {errors.address && <FieldError>{errors.address.message}</FieldError>}
      </Field>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("cityLabel")}</FieldLabel>
          <Input {...register("city")} placeholder={t("cityPlaceholder")} />
          {errors.city && <FieldError>{errors.city.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("provinceLabel")}</FieldLabel>
          <Input {...register("province")} placeholder={t("provincePlaceholder")} />
          {errors.province && <FieldError>{errors.province.message}</FieldError>}
        </Field>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("phoneLabel")}</FieldLabel>
          <Input {...register("phone")} placeholder={t("phonePlaceholder")} />
          {errors.phone && <FieldError>{errors.phone.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("emailLabel")}</FieldLabel>
          <Input type="email" {...register("email")} placeholder={t("emailPlaceholder")} />
          {errors.email && <FieldError>{errors.email.message}</FieldError>}
        </Field>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("latitudeLabel")}</FieldLabel>
          <Input
            type="number"
            step="any"
            inputMode="decimal"
            {...register("latitude", {
              setValueAs: (value) => (value === "" || value === null ? undefined : Number(value)),
            })}
            placeholder="-6.208800"
          />
          {errors.latitude && <FieldError>{errors.latitude.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("longitudeLabel")}</FieldLabel>
          <Input
            type="number"
            step="any"
            inputMode="decimal"
            {...register("longitude", {
              setValueAs: (value) => (value === "" || value === null ? undefined : Number(value)),
            })}
            placeholder="106.845600"
          />
          {errors.longitude && <FieldError>{errors.longitude.message}</FieldError>}
        </Field>
      </div>

      <div>
        <Button
          type="button"
          variant="outline"
          onClick={() => setIsPickerOpen(true)}
          disabled={isLoading}
        >
          {t("pickOnMap")}
        </Button>
      </div>

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

      <div className="flex flex-col-reverse sm:flex-row justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading} className="w-full sm:w-auto">
          {t("cancel")}
        </Button>
        <Button type="submit" disabled={isLoading} className="w-full sm:w-auto">
          {isLoading ? t("submitting") : isEdit ? t("submitUpdate") : t("submitCreate")}
        </Button>
      </div>

      <AccountLocationPickerDialog
        open={isPickerOpen}
        onOpenChange={setIsPickerOpen}
        initialLocation={
          typeof watch("latitude") === "number" && typeof watch("longitude") === "number"
            ? { lat: watch("latitude") as number, lng: watch("longitude") as number }
            : undefined
        }
        onConfirm={(loc) => {
          setValue("latitude", loc.lat, { shouldDirty: true, shouldValidate: true });
          setValue("longitude", loc.lng, { shouldDirty: true, shouldValidate: true });
          setIsPickerOpen(false);
        }}
        title={t("pickOnMapTitle")}
        confirmText={t("pickOnMapConfirm")}
        cancelText={t("cancel")}
      />
    </form>
  );
}

