"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslations } from "next-intl";
import { activityTypeSchema, type ActivityTypeFormData } from "../schemas/activity-type.schema";
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
import type { ActivityType } from "../types/activity-type";

interface ActivityTypeFormProps {
  readonly activityType?: ActivityType;
  readonly onSubmit: (data: ActivityTypeFormData) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

export function ActivityTypeForm({
  activityType,
  onSubmit,
  onCancel,
  isLoading,
}: ActivityTypeFormProps) {
  const isEdit = !!activityType;
  const t = useTranslations("visitReportActivityType.form");

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<ActivityTypeFormData>({
    resolver: zodResolver(activityTypeSchema),
    defaultValues: activityType
      ? {
          name: activityType.name,
          code: activityType.code,
          description: activityType.description || "",
          icon: activityType.icon || "",
          badge_color: activityType.badge_color,
          status: activityType.status,
          order: activityType.order,
        }
      : {
          name: "",
          code: "",
          description: "",
          icon: "activity",
          badge_color: "default",
          status: "active",
          order: 0,
        },
  });

  const selectedStatus = watch("status");
  const selectedBadgeColor = watch("badge_color");

  const handleFormSubmit = async (data: ActivityTypeFormData) => {
    await onSubmit(data);
  };

  // Popular icons list
  const popularIcons = [
    { value: "activity", label: "Activity" },
    { value: "phone", label: "Phone" },
    { value: "mail", label: "Mail" },
    { value: "handshake", label: "Handshake" },
    { value: "message-square", label: "Message" },
    { value: "users", label: "Users" },
    { value: "calendar", label: "Calendar" },
    { value: "file-text", label: "File" },
    { value: "clock", label: "Clock" },
    { value: "check-circle", label: "Check Circle" },
    { value: "heart", label: "Heart" },
  ];

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <div className="grid grid-cols-2 gap-4">
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
            disabled={isLoading || isEdit} // Code cannot be edited once created
          />
          {!isEdit && (
            <p className="text-[10px] text-muted-foreground mt-0.5">
              {t("codeHint")}
            </p>
          )}
          {errors.code && <FieldError>{errors.code.message}</FieldError>}
        </Field>
      </div>

      <Field orientation="vertical">
        <FieldLabel>{t("descriptionLabel")}</FieldLabel>
        <Textarea
          {...register("description")}
          placeholder={t("descriptionPlaceholder")}
          rows={2}
          disabled={isLoading}
        />
        {errors.description && <FieldError>{errors.description.message}</FieldError>}
      </Field>

      <div className="grid grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("iconLabel")}</FieldLabel>
          <Select
            disabled={isLoading}
            value={watch("icon") || "activity"}
            onValueChange={(val) => setValue("icon", val, { shouldValidate: true })}
          >
            <SelectTrigger>
              <SelectValue placeholder={t("iconPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              {popularIcons.map((ico) => (
                <SelectItem key={ico.value} value={ico.value}>
                  {ico.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <p className="text-[10px] text-muted-foreground mt-0.5">
            {t("iconHint")}
          </p>
          {errors.icon && <FieldError>{errors.icon.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("badgeColorLabel")}</FieldLabel>
          <Select
            disabled={isLoading}
            value={selectedBadgeColor}
            onValueChange={(val) =>
              setValue("badge_color", val as any, { shouldValidate: true })
            }
          >
            <SelectTrigger>
              <SelectValue placeholder={t("badgeColorPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="default">{t("badgeColorOptions.default")}</SelectItem>
              <SelectItem value="secondary">{t("badgeColorOptions.secondary")}</SelectItem>
              <SelectItem value="destructive">{t("badgeColorOptions.destructive")}</SelectItem>
              <SelectItem value="outline">{t("badgeColorOptions.outline")}</SelectItem>
              <SelectItem value="success">{t("badgeColorOptions.success")}</SelectItem>
              <SelectItem value="warning">{t("badgeColorOptions.warning")}</SelectItem>
              <SelectItem value="active">{t("badgeColorOptions.active")}</SelectItem>
            </SelectContent>
          </Select>
          {errors.badge_color && <FieldError>{errors.badge_color.message}</FieldError>}
        </Field>
      </div>

      <div className="grid grid-cols-2 gap-4 items-center">
        <Field orientation="vertical">
          <FieldLabel>{t("orderLabel")}</FieldLabel>
          <Input
            type="number"
            placeholder={t("orderPlaceholder")}
            disabled={isLoading}
            onChange={(e) => {
              const val = e.target.value === "" ? 0 : parseInt(e.target.value, 10);
              setValue("order", val, { shouldValidate: true });
            }}
            value={watch("order")}
          />
          <p className="text-[10px] text-muted-foreground mt-0.5">
            {t("orderHint")}
          </p>
          {errors.order && <FieldError>{errors.order.message}</FieldError>}
        </Field>

        <div className="pt-5">
          <div className="flex items-center space-x-3 rounded-md border p-3 bg-muted/20">
            <Checkbox
              id="status"
              checked={selectedStatus === "active"}
              onCheckedChange={(checked) =>
                setValue("status", checked ? "active" : "inactive", {
                  shouldValidate: true,
                })
              }
              disabled={isLoading}
            />
            <div className="space-y-0.5 leading-none">
              <Label htmlFor="status" className="cursor-pointer font-medium text-sm">
                {t("statusActive")}
              </Label>
              <p className="text-[10px] text-muted-foreground">
                {t("statusHint")}
              </p>
            </div>
          </div>
          {errors.status && <FieldError>{errors.status.message}</FieldError>}
        </div>
      </div>

      <div className="flex justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading}>
          {t("cancel")}
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading
            ? t("submitting")
            : isEdit
            ? t("update")
            : t("create")}
        </Button>
      </div>
    </form>
  );
}
