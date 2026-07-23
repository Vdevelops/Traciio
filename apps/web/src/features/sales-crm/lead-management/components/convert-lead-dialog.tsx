"use client";

import { Controller, useFieldArray, useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
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
import { NumberInput } from "@/components/ui/number-input";
import { convertLeadSchema, type ConvertLeadFormData } from "../schemas/lead.schema";
import { useConvertLead } from "../hooks/useLeads";
import { useStages } from "../../pipeline-management/hooks/useStages";
import { useProducts } from "@/features/sales-crm/product-management/hooks/useProducts";
import type { ConvertLeadResponse, Lead } from "../types";
import { toast } from "sonner";
import { useTranslations } from "next-intl";
import { useEffect, useMemo } from "react";
import { Info, Trash2 } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";

interface ConvertLeadDialogProps {
  readonly lead: Lead;
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onSuccess?: (response: ConvertLeadResponse) => void;
}

export function ConvertLeadDialog({
  lead,
  open,
  onOpenChange,
  onSuccess,
}: ConvertLeadDialogProps) {
  const t = useTranslations("leadManagement.convertLead");
  const convertLead = useConvertLead();
  const { data: stages } = useStages();
  const { data: productsData, refetch: refetchProducts } = useProducts({ per_page: 100, status: "active" });
  const products = productsData?.data ?? [];
  const activeStages = useMemo(() => {
    return (stages ?? []).filter((stage) => stage.is_active);
  }, [stages]);
  const convertibleStages = useMemo(() => {
    return activeStages.filter((stage) => stage.is_won || stage.code === "closed_won");
  }, [activeStages]);
  const {
    register,
    handleSubmit,
    control,
    setValue,
    reset,
    formState: { errors },
  } = useForm<ConvertLeadFormData>({
    resolver: zodResolver(convertLeadSchema),
    defaultValues: {
      opportunity_title: "",
      opportunity_description: "",
      stage_id: "",
      value: undefined,
      status_reason: "",
      product_items: [],
    },
  });
  const selectedStageId = useWatch({ control, name: "stage_id" });
  const valueAmount = useWatch({ control, name: "value" });
  const productItems = useWatch({ control, name: "product_items" });
  const productItemsFieldArray = useFieldArray({
    control,
    name: "product_items",
  });

  useEffect(() => {
    if (open && convertibleStages.length > 0) {
      void refetchProducts();
      const sortedStages = [...convertibleStages].sort((a, b) => a.order - b.order);
      const defaultStage = sortedStages[0];
      const initialStageId = defaultStage?.id || "";
      const initialValue = lead.estimated_value || undefined;

      reset({
        opportunity_title: `${lead.company_name || lead.first_name} - ${lead.industry || "Opportunity"}`,
        opportunity_description: lead.notes || "",
        stage_id: initialStageId,
        value: initialValue !== undefined ? initialValue / 100 : undefined,
        status_reason: "",
        product_items: [
          {
            product_id: "",
            quantity: 1,
            unit_price: 0,
            discount_amount: 0,
            notes: "",
          },
        ],
      });
    }
  }, [open, lead, convertibleStages, refetchProducts, reset]);

  useEffect(() => {
    const currentItems = productItems ?? [];
    const total = currentItems.reduce((sum, item) => {
      const subtotal = Math.max(0, (item.unit_price ?? 0) * (item.quantity ?? 0) - (item.discount_amount ?? 0));
      return sum + subtotal;
    }, 0);
    if (currentItems.length > 0) {
      setValue("value", Math.round(total), { shouldValidate: true });
    }
  }, [productItems, setValue]);

  const onSubmit = async (data: ConvertLeadFormData) => {
    try {
      const payload = { ...data };

      if (payload.value !== undefined) {
        payload.value = payload.value * 100;
      }
      payload.product_items = (payload.product_items ?? [])
        .filter((item) => item.product_id && item.quantity > 0)
        .map((item) => ({
          ...item,
          unit_price: Math.round((item.unit_price ?? 0) * 100),
          discount_amount: Math.round((item.discount_amount ?? 0) * 100),
        }));
      if (payload.product_items.length === 0) {
        toast.error("Produk deal wajib diisi");
        return;
      }

      const response = await convertLead.mutateAsync({ id: lead.id, data: payload });
      toast.success(t("toast.success"));
      onOpenChange(false);
      onSuccess?.(response);
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto mx-2 sm:mx-auto">
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>

        {/* Info Alert */}
        <Alert>
          <Info className="h-4 w-4" />
          <AlertDescription>
            {t("autoCreateInfo")}
          </AlertDescription>
        </Alert>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Field orientation="vertical">
            <FieldLabel>{t("fields.opportunityTitle")} *</FieldLabel>
            <Input
              {...register("opportunity_title")}
              placeholder={t("fields.opportunityTitlePlaceholder")}
            />
            {errors.opportunity_title && (
              <FieldError>{errors.opportunity_title.message}</FieldError>
            )}
          </Field>

          <Field orientation="vertical">
            <FieldLabel>{t("fields.opportunityDescription")}</FieldLabel>
            <Textarea
              {...register("opportunity_description")}
              placeholder={t("fields.opportunityDescriptionPlaceholder")}
              rows={3}
            />
            {errors.opportunity_description && (
              <FieldError>{errors.opportunity_description.message}</FieldError>
            )}
          </Field>

          <Field orientation="vertical">
            <FieldLabel>{t("fields.stage")} *</FieldLabel>
            <Select
              value={selectedStageId}
              onValueChange={(value) => setValue("stage_id", value)}
            >
              <SelectTrigger>
                <SelectValue placeholder={t("fields.stagePlaceholder")} />
              </SelectTrigger>
              <SelectContent>
                {convertibleStages.map((stage) => (
                  <SelectItem key={stage.id} value={stage.id}>
                    {stage.name} - {stage.code}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {errors.stage_id && (
              <FieldError>{errors.stage_id.message}</FieldError>
            )}
          </Field>

          <Field orientation="vertical">
            <FieldLabel>{t("fields.value")}</FieldLabel>
            <NumberInput
              value={valueAmount || 0}
              onChange={(value) => setValue("value", value)}
              placeholder="0"
              min={0}
            />
            {errors.value && (
              <FieldError>{errors.value.message}</FieldError>
            )}
          </Field>

          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <FieldLabel>Produk Deal</FieldLabel>
              <Button
                type="button"
                variant="outline"
                onClick={() =>
                  productItemsFieldArray.append({
                    product_id: "",
                    quantity: 1,
                    unit_price: 0,
                    discount_amount: 0,
                    notes: "",
                  })
                }
              >
                Tambah Produk
              </Button>
            </div>
            {productItemsFieldArray.fields.map((field, index) => (
              <div key={field.id} className="rounded-md border p-3">
                <div className="grid grid-cols-1 gap-3 sm:grid-cols-12">
                  <div className="sm:col-span-5">
                    <Field orientation="vertical">
                      <FieldLabel>Produk</FieldLabel>
                      <Controller
                        control={control}
                        name={`product_items.${index}.product_id`}
                        render={({ field }) => (
                          <Select
                            value={field.value ?? ""}
                            onValueChange={(value) => {
                              field.onChange(value);
                              const product = products.find((item) => item.id === value);
                              if (product) {
                                setValue(`product_items.${index}.unit_price`, product.price / 100, {
                                  shouldValidate: true,
                                });
                              }
                            }}
                          >
                            <SelectTrigger>
                              <SelectValue placeholder="Pilih produk" />
                            </SelectTrigger>
                            <SelectContent>
                              {products.map((product) => (
                                <SelectItem key={product.id} value={product.id}>
                                  {product.name} {product.sku ? `(${product.sku})` : ""}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        )}
                      />
                    </Field>
                  </div>
                  <div className="sm:col-span-2">
                    <Field orientation="vertical">
                      <FieldLabel>Qty</FieldLabel>
                      <Input type="number" min={1} {...register(`product_items.${index}.quantity`, { valueAsNumber: true })} />
                    </Field>
                  </div>
                  <div className="sm:col-span-2">
                    <Field orientation="vertical">
                      <FieldLabel>Harga</FieldLabel>
                      <Controller
                        control={control}
                        name={`product_items.${index}.unit_price`}
                        render={({ field }) => (
                          <NumberInput value={field.value ?? 0} onChange={(value) => field.onChange(value ?? 0)} min={0} />
                        )}
                      />
                    </Field>
                  </div>
                  <div className="sm:col-span-2">
                    <Field orientation="vertical">
                      <FieldLabel>Diskon</FieldLabel>
                      <Controller
                        control={control}
                        name={`product_items.${index}.discount_amount`}
                        render={({ field }) => (
                          <NumberInput value={field.value ?? 0} onChange={(value) => field.onChange(value ?? 0)} min={0} />
                        )}
                      />
                    </Field>
                  </div>
                  <div className="flex items-end justify-end sm:col-span-1">
                    <Button type="button" variant="outline" size="icon" onClick={() => productItemsFieldArray.remove(index)}>
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </div>
            ))}
            {errors.product_items?.message && <FieldError>{errors.product_items.message}</FieldError>}
          </div>

          <Field orientation="vertical">
            <FieldLabel>{t("fields.statusReason")} *</FieldLabel>
            <Textarea
              {...register("status_reason")}
              placeholder={t("fields.statusReasonPlaceholder")}
              rows={3}
            />
            {errors.status_reason && (
              <FieldError>{errors.status_reason.message}</FieldError>
            )}
          </Field>

          <div className="flex flex-col-reverse sm:flex-row gap-3 justify-end pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={convertLead.isPending}
              className="w-full sm:w-auto"
            >
              {t("buttons.cancel")}
            </Button>
            <Button type="submit" disabled={convertLead.isPending} className="w-full sm:w-auto">
              {convertLead.isPending ? t("buttons.converting") : t("buttons.convert")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
