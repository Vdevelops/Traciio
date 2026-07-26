"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
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
import { createActivitySchema, type CreateActivityFormData } from "../schemas/activity.schema";
import { useCreateActivity, useUpdateActivity } from "../hooks/useVisitReports";
import { useActivityTypes } from "../hooks/useActivityTypes";
import { toast } from "sonner";
import { useEffect, useMemo, useState } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { useTranslations } from "next-intl";
import {
  ProductInterestEditor,
  type ProductInterestItem,
} from "./product-interest-editor";
import type { Activity } from "../types/activity";
import {
  formatDateTimeWithLocalOffset,
  parseDateTimeInputToLocalOffset,
  parseWallClockDateTime,
} from "@/lib/utils";

function getCurrentTimestampIso(): string {
  return formatDateTimeWithLocalOffset(new Date());
}

function formatIsoForDateTimeInput(value?: string): string {
  const date = value ? parseWallClockDateTime(value) : new Date();

  if (!date || Number.isNaN(date.getTime())) {
    return "";
  }

  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");

  return `${year}-${month}-${day}T${hours}:${minutes}`;
}

interface CreateActivityDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly accountId?: string;
  readonly contactId?: string;
  readonly dealId?: string;
  readonly leadId?: string;
  readonly onSuccess?: () => void;
  readonly showProductInterests?: boolean;
  readonly activity?: Activity | null;
}

export function CreateActivityDialog({
  open,
  onOpenChange,
  accountId,
  contactId,
  dealId,
  leadId,
  onSuccess,
  showProductInterests = true,
  activity,
}: CreateActivityDialogProps) {
  const t = useTranslations("createActivityDialog");
  const isEdit = !!activity;
  const createActivity = useCreateActivity();
  const updateActivity = useUpdateActivity();
  const { data: activityTypesData, isLoading: isLoadingTypes } = useActivityTypes({ status: "active" });
  const [productInterests, setProductInterests] = useState<ProductInterestItem[]>([]);

  const activityTypes = useMemo(() => {
    return activityTypesData?.data ?? [];
  }, [activityTypesData]);

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    reset,
    formState: { errors },
  } = useForm<CreateActivityFormData>({
    resolver: zodResolver(createActivitySchema),
    defaultValues: {
      activity_type_id: "",
      account_id: accountId,
      contact_id: contactId,
      deal_id: dealId,
      lead_id: leadId,
      description: "",
      timestamp: activity?.timestamp ?? getCurrentTimestampIso(),
    },
  });

  // Reset form when dialog opens
  useEffect(() => {
    if (open) {
      const metadata = activity?.metadata;
      reset({
        activity_type_id: activity?.activity_type_id ?? "",
        account_id: activity?.account_id ?? accountId,
        contact_id: activity?.contact_id ?? contactId,
        deal_id: activity?.deal_id ?? dealId,
        lead_id: activity?.lead_id ?? leadId,
        description: activity?.description ?? "",
        timestamp: activity?.timestamp ?? getCurrentTimestampIso(),
      });
      setProductInterests(
        Array.isArray(metadata?.product_interests)
          ? (metadata.product_interests as ProductInterestItem[])
          : [],
      );
    }
  }, [open, accountId, contactId, dealId, leadId, reset, activity]);

  // Set default activity type when types are loaded
  useEffect(() => {
    const currentTypeId = watch("activity_type_id");
    if (activityTypes.length > 0 && !currentTypeId) {
      setValue("activity_type_id", activityTypes[0].id);
    }
  }, [activityTypes, watch, setValue]);

  const onSubmit = async (data: CreateActivityFormData) => {
    try {
      const finalAccountId = accountId || data.account_id;
      const finalContactId = contactId || data.contact_id;

      // Prepare request payload
      const payload: {
        activity_type_id: string;
        account_id?: string;
        contact_id?: string;
        deal_id?: string;
        lead_id?: string;
        description: string;
        timestamp: string;
        metadata?: Record<string, unknown>;
      } = {
        activity_type_id: data.activity_type_id,
        description: data.description,
        timestamp: data.timestamp,
        metadata: productInterests.length > 0 ? { product_interests: productInterests } : {},
      };

      // Include account_id if available. For lead-stage activity, lead_id alone is valid.
      if (finalAccountId) {
        payload.account_id = finalAccountId;
      }

      // Include contact_id if available
      if (finalContactId) {
        payload.contact_id = finalContactId;
      }

      // Include deal_id if available
      const finalDealId = dealId || data.deal_id;
      if (finalDealId) {
        payload.deal_id = finalDealId;
      }

      // Include lead_id if available
      const finalLeadId = leadId || data.lead_id;
      if (finalLeadId) {
        payload.lead_id = finalLeadId;
      }

      if (isEdit && activity) {
        await updateActivity.mutateAsync({ id: activity.id, data: payload });
        toast.success(t("toast.updated"));
      } else {
        await createActivity.mutateAsync(payload);
        toast.success(t("toast.created"));
      }
      const defaultTypeId = activityTypes.length > 0 ? activityTypes[0]?.id : "";
      reset({
        activity_type_id: defaultTypeId,
        account_id: accountId,
        contact_id: contactId,
        deal_id: dealId,
        lead_id: leadId,
        description: "",
        timestamp: getCurrentTimestampIso(),
      });
      setProductInterests([]);
      onOpenChange(false);
      onSuccess?.();
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl border-border/70 bg-card/95">
        <DialogHeader>
          <DialogTitle>{isEdit ? t("editTitle") : t("title")}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Field orientation="vertical">
            <FieldLabel>{t("activityTypeLabel")} *</FieldLabel>
            {isLoadingTypes ? (
              <Skeleton className="h-10 w-full" />
            ) : (
              <Select
                value={watch("activity_type_id") ?? ""}
                onValueChange={(value) => setValue("activity_type_id", value, { shouldValidate: true })}
              >
                <SelectTrigger>
                  <SelectValue placeholder={t("activityTypePlaceholder")} />
                </SelectTrigger>
                <SelectContent>
                  {activityTypes.length === 0 ? (
                    <div className="p-2 text-sm text-muted-foreground text-center">
                      No activity types available
                    </div>
                  ) : (
                    activityTypes.map((type) => (
                      <SelectItem key={type.id} value={type.id}>
                        {type.name}
                      </SelectItem>
                    ))
                  )}
                </SelectContent>
              </Select>
            )}
            {errors.activity_type_id && <FieldError>{errors.activity_type_id.message}</FieldError>}
          </Field>

          <Field orientation="vertical">
            <FieldLabel>{t("descriptionLabel")}</FieldLabel>
            <Textarea
              {...register("description")}
              placeholder={t("descriptionPlaceholder")}
              rows={4}
            />
            {errors.description && <FieldError>{errors.description.message}</FieldError>}
          </Field>

          <Field orientation="vertical">
            <FieldLabel>{t("timestampLabel")} *</FieldLabel>
            <Input
              type="datetime-local"
              value={
                watch("timestamp")
                  ? formatIsoForDateTimeInput(watch("timestamp"))
                  : formatIsoForDateTimeInput()
              }
              onChange={(e) => {
                const value = e.target.value;
                if (value) {
                  setValue("timestamp", parseDateTimeInputToLocalOffset(value), { shouldValidate: true });
                }
              }}
            />
            {errors.timestamp && <FieldError>{errors.timestamp.message}</FieldError>}
          </Field>

          {showProductInterests && (
            <ProductInterestEditor
              value={productInterests}
              onChange={setProductInterests}
              className="crm-stack rounded-2xl border border-border/70 bg-muted/20 p-4"
            />
          )}

          <div className="flex justify-end gap-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                reset();
                setProductInterests([]);
                onOpenChange(false);
              }}
              disabled={createActivity.isPending || updateActivity.isPending}
            >
              {t("cancel")}
            </Button>
            <Button type="submit" disabled={createActivity.isPending || updateActivity.isPending}>
              {isEdit
                ? (updateActivity.isPending ? t("updating") : t("update"))
                : (createActivity.isPending ? t("creating") : t("create"))}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
