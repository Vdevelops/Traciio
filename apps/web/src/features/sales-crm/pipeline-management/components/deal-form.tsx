"use client";

import { useFieldArray, useForm, Controller, type ControllerRenderProps, type Resolver } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  createDealSchema,
  updateDealSchema,
  type CreateDealFormData,
  type UpdateDealFormData,
} from "../schemas/deal.schema";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { NumberInput } from "@/components/ui/number-input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { DatePicker } from "@/components/ui/date-picker";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { usePipelines } from "../hooks/usePipelines";
import { useAccounts } from "@/features/sales-crm/account-management/hooks/useAccounts";
import { useContacts } from "@/features/sales-crm/account-management/hooks/useContacts";
import { useLead, useLeads } from "@/features/sales-crm/lead-management/hooks/useLeads";
import { useLeadSources } from "@/features/sales-crm/lead-management/hooks/useLeadSources";
import { useProducts } from "@/features/sales-crm/product-management/hooks/useProducts";

import { Trash2 } from "lucide-react";
import type { Deal } from "../types";
import { useEffect, useMemo } from "react";
import { useTranslations } from "next-intl";

type DealFormProps =
  | {
      readonly deal?: Deal;
      readonly initialLeadId?: string;
      readonly initialAccountId?: string;
      readonly showQualifiedLeadDropdown?: boolean;
      readonly onSubmit: (data: CreateDealFormData) => Promise<void>;
      readonly onCancel: () => void;
      readonly isLoading?: boolean;
    }
  | {
      readonly deal: Deal;
      readonly initialLeadId?: string;
      readonly initialAccountId?: string;
      readonly showQualifiedLeadDropdown?: boolean;
      readonly onSubmit: (data: UpdateDealFormData) => Promise<void>;
      readonly onCancel: () => void;
      readonly isLoading?: boolean;
    };

export function DealForm({ deal, initialLeadId, initialAccountId, showQualifiedLeadDropdown = false, onSubmit, onCancel, isLoading }: DealFormProps) {
  const t = useTranslations("pipelineManagement.dealForm");

  const isEdit = !!deal;
  const { data: pipelinesData } = usePipelines({ is_active: true });
  const { data: accountsData } = useAccounts({ status: "active", per_page: 100 });
  const { data: leadData } = useLead(initialLeadId || "");
  const { data: qualifiedLeadsData } = useLeads({ status: "qualified", per_page: 100 });
  const { data: leadSourcesData } = useLeadSources({ is_active: true, per_page: 100 });
  // NOTE: product_items depends on products API. If backend is temporarily failing (e.g., migrations not applied),

  // we keep the UI stable by not hammering the endpoint in a tight loop.
  const { data: productsData } = useProducts({ per_page: 100, status: "active" });
  
  const pipelines = useMemo(() => pipelinesData?.data || [], [pipelinesData?.data]);
  const accounts = accountsData?.data || [];
  const lead = leadData?.data;
  const qualifiedLeads = qualifiedLeadsData?.data || [];
  const leadSources = leadSourcesData?.data || [];
  const products = productsData?.data || [];


  const {
    register,
    handleSubmit,
    control,
    setValue,
    watch,
    reset,
    formState: { errors },
  } = useForm<CreateDealFormData | UpdateDealFormData>({
    resolver: zodResolver(isEdit ? updateDealSchema : createDealSchema) as Resolver<CreateDealFormData | UpdateDealFormData, unknown, CreateDealFormData | UpdateDealFormData>,
    defaultValues: deal
      ? {
          title: deal.title,
          description: deal.description || "",
          account_id: deal.account_id,
          contact_id: deal.contact_id || "",
          stage_id: deal.stage_id,
          value:
            ((deal.product_items || []).reduce((sum, item) => sum + (item.subtotal ?? 0), 0) || deal.value || 0) / 100,
          probability: deal.stage?.probability ?? deal.probability ?? 0,
          expected_close_date: deal.expected_close_date ? deal.expected_close_date.split("T")[0] : "",
          lead_id: deal.lead_id || "",
          notes: deal.notes || "",
          product_items: (deal.product_items || []).map((item) => ({
            product_id: item.product_id,
            quantity: item.quantity ?? 1,
            unit_price: (item.unit_price ?? 0) / 100,
            discount_amount: item.discount_amount ? item.discount_amount / 100 : 0,
            notes: item.notes || "",
          })),
        }
      : {
          value: 0,
          probability: 0,
          account_id: initialAccountId || "",
          lead_id: initialLeadId || "",
          expected_close_date: new Date().toISOString().split("T")[0],
          product_items: [],
        },
  });

    const stageId = watch("stage_id");
    const productItems = watch("product_items") || [];
    const valueFallbackRupiah = watch("value") ?? 0;

    const stageProbability = useMemo(() => {
      const stage = pipelines.find((s) => s.id === stageId);
      return stage?.probability ?? 0;
    }, [pipelines, stageId]);

    const computedValueRupiah = useMemo(() => {
    if (productItems.length === 0) {
      return valueFallbackRupiah;
    }
    return productItems.reduce((sum, item) => {
      const unitPrice = Number(item?.unit_price ?? 0);
      const quantity = Number(item?.quantity ?? 0);
      const discount = Number(item?.discount_amount ?? 0);
      const subtotal = Math.max(0, unitPrice * quantity - discount);
      return sum + subtotal;
    }, 0);
    }, [productItems, valueFallbackRupiah]);

    // Keep probability synced with selected stage
    useEffect(() => {
      setValue("probability", stageProbability, { shouldValidate: true });
    }, [setValue, stageProbability]);

  const productItemsFieldArray = useFieldArray({
    control,
    name: "product_items" as never,
  });

  // Pre-fill form when lead or account data is loaded
  useEffect(() => {
    if (!isEdit && (lead || initialAccountId)) {
      const defaultValues: Partial<CreateDealFormData> = {
        value: 0,
        probability: 0,
      };

      // Pre-fill from lead
      if (lead) {
        defaultValues.title = lead.company_name 
          ? `${lead.company_name} - ${lead.industry || "Opportunity"}`
          : `${lead.first_name} ${lead.last_name || ""}`.trim() || "New Opportunity";
        defaultValues.description = lead.notes || "";
        defaultValues.account_id = lead.account_id || initialAccountId || "";
        defaultValues.contact_id = lead.contact_id || "";
        defaultValues.expected_close_date = lead.expected_close_date 
          ? lead.expected_close_date.split("T")[0] 
          : new Date().toISOString().split("T")[0];
        defaultValues.notes = lead.notes || "";
        defaultValues.lead_id = lead.id;
        defaultValues.source = lead.lead_source || "";
      } else if (initialAccountId) {
        // Pre-fill from account
        defaultValues.account_id = initialAccountId;
      }

      // Set default stage (first active stage)
      if (pipelines.length > 0 && !defaultValues.stage_id) {
        const sortedStages = [...pipelines]
          .filter((stage) => stage.is_active)
          .sort((a, b) => a.order - b.order);
        if (sortedStages.length > 0) {
          defaultValues.stage_id = sortedStages[0].id;
        }
      }

      reset(defaultValues);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lead?.id, initialAccountId, pipelines.length, isEdit]);

  // Convert expected_close_date string to Date for DatePicker
  // Create date as local midnight to avoid timezone issues
  const expectedCloseDateValue = watch("expected_close_date");
  const expectedCloseDate = expectedCloseDateValue 
    ? (() => {
        const [year, month, day] = expectedCloseDateValue.split("-").map(Number);
        return new Date(year, month - 1, day);
      })()
    : undefined;

  // Load contacts when account is selected
  const accountId = watch("account_id");
  const { data: contactsData } = useContacts({ 
    account_id: accountId || deal?.account_id, 
    per_page: 100 
  });
  const contacts = contactsData?.data || [];

  // Reset contact when account changes
  const initialDealAccountId = (deal as Deal | undefined)?.account_id;
  useEffect(() => {
    if (!isEdit && accountId && accountId !== initialDealAccountId) {
      setValue("contact_id", "");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [accountId, isEdit, initialDealAccountId]);

  // Handle qualified lead selection from dropdown
  const selectedLeadId = watch("lead_id");
  const selectedQualifiedLead = qualifiedLeads.find(l => l.id === selectedLeadId);
  
  useEffect(() => {
    if (!isEdit && showQualifiedLeadDropdown && selectedQualifiedLead && selectedLeadId && !initialLeadId) {
      // Pre-fill form when qualified lead is selected
      const defaultValues: Partial<CreateDealFormData> = {
        title: selectedQualifiedLead.company_name 
          ? `${selectedQualifiedLead.company_name} - ${selectedQualifiedLead.industry || "Opportunity"}`
          : `${selectedQualifiedLead.first_name} ${selectedQualifiedLead.last_name || ""}`.trim() || "New Opportunity",
        description: selectedQualifiedLead.notes || "",
        account_id: selectedQualifiedLead.account_id || "",
        contact_id: selectedQualifiedLead.contact_id || "",
        expected_close_date: selectedQualifiedLead.expected_close_date 
          ? selectedQualifiedLead.expected_close_date.split("T")[0] 
          : new Date().toISOString().split("T")[0],
        notes: selectedQualifiedLead.notes || "",
        lead_id: selectedQualifiedLead.id,
        source: selectedQualifiedLead.lead_source || "",
      };

      // Set default stage if pipelines available
      if (pipelines.length > 0) {
        const sortedStages = [...pipelines]
          .filter((stage) => stage.is_active)
          .sort((a, b) => a.order - b.order);
        if (sortedStages.length > 0) {
          defaultValues.stage_id = sortedStages[0].id;
        }
      }

      // Update form values
      Object.entries(defaultValues).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== "") {
          setValue(key as keyof CreateDealFormData, value as never, { shouldValidate: false });
        }
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedQualifiedLead?.id, selectedLeadId, showQualifiedLeadDropdown, isEdit, initialLeadId]);

  const handleFormSubmit = async (data: CreateDealFormData | UpdateDealFormData) => {
    const normalizedItems = (data.product_items || [])
      .filter((item) => item.product_id && item.quantity > 0)
      .map((item) => {
        const unitPrice = Math.round((item.unit_price ?? 0) * 100);
        const discountAmount = Math.round((item.discount_amount ?? 0) * 100);
        return {
          ...item,
          unit_price: unitPrice,
          discount_amount: discountAmount,
        };
      });

  const computedValueSen =
    normalizedItems.length > 0
      ? normalizedItems.reduce((sum, item) => {
          const subtotal = Math.max(0, (item.unit_price ?? 0) * (item.quantity ?? 0) - (item.discount_amount ?? 0));
          return sum + subtotal;
        }, 0)
      : Math.round((data.value ?? 0) * 100);

    const submitData: Record<string, unknown> = {
      ...data,
      // Value is always derived from product_items if present, otherwise manual value
      value: computedValueSen,
      // Use manually entered probability if provided, otherwise fall back to stage probability
      probability: data.probability ?? stageProbability,
      product_items: normalizedItems,
    };
    delete submitData.assigned_to;

    // Format expected_close_date to ISO if present
    if (data.expected_close_date) {
      const date = new Date(data.expected_close_date + "T00:00:00");
      if (!Number.isNaN(date.getTime())) {
        submitData.expected_close_date = date.toISOString();
      }
    } else {
      delete submitData.expected_close_date;
    }

    // Clean up empty optional fields for backend
    const optionalFields = ["contact_id", "lead_id", "description", "notes", "source"];
    optionalFields.forEach((field) => {
      if (!submitData[field]) {
        delete submitData[field];
      }
    });

    if (isEdit) {
      await (onSubmit as (data: UpdateDealFormData) => Promise<void>)(submitData as UpdateDealFormData);
      return;
    }

    await (onSubmit as (data: CreateDealFormData) => Promise<void>)(submitData as CreateDealFormData);
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <Field orientation="vertical">
        <FieldLabel>{t("titleRequired")}</FieldLabel>
        <Input {...register("title")} placeholder={t("titlePlaceholder")} />
        {errors.title && <FieldError>{errors.title.message}</FieldError>}
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

      {showQualifiedLeadDropdown && (
        <Field orientation="vertical">
          <FieldLabel>{t("qualifiedLeadLabel") || "Qualified Lead"}</FieldLabel>
          <Controller
            control={control}
            name="lead_id"
            render={({ field }) => (
              <Select
                value={field.value || "none"}
                onValueChange={(value) => field.onChange(value === "none" ? "" : value)}
              >
                <SelectTrigger>
                  <SelectValue placeholder={t("qualifiedLeadPlaceholder") || "Select qualified lead (optional)"} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">{t("none") || "None"}</SelectItem>
                  {qualifiedLeads.map((lead) => (
                    <SelectItem key={lead.id} value={lead.id}>
                      {lead.company_name || `${lead.first_name} ${lead.last_name || ""}`.trim()}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          />
          {errors.lead_id && <FieldError>{errors.lead_id.message}</FieldError>}
        </Field>
      )}

      <Field orientation="vertical">
        <FieldLabel>{t("accountRequired")}</FieldLabel>
        <Controller
          control={control}
          name="account_id"
          render={({ field }) => (
            <Select
              value={field.value || ""}
              onValueChange={field.onChange}
            >
              <SelectTrigger>
                <SelectValue placeholder={t("accountPlaceholder")} />
              </SelectTrigger>
              <SelectContent>
                {accounts.map((account) => (
                  <SelectItem key={account.id} value={account.id}>
                    {account.name} {account.category && `(${account.category.name})`}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        />
        {errors.account_id && <FieldError>{errors.account_id.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("contactLabel")} *</FieldLabel>
        <Controller
          control={control}
          name="contact_id"
          render={({ field }) => (
            <Select
              value={field.value || ""}
              onValueChange={field.onChange}
              disabled={!watch("account_id") && !deal?.account_id}
            >
              <SelectTrigger>
                <SelectValue placeholder={t("contactPlaceholder")} />
              </SelectTrigger>
              <SelectContent>
                {contacts.map((contact) => (
                  <SelectItem key={contact.id} value={contact.id}>
                    {contact.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        />
        {errors.contact_id && <FieldError>{errors.contact_id.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("stageRequired")}</FieldLabel>
        <Controller
          control={control}
          name="stage_id"
          render={({ field }) => (
            <Select
              value={field.value || ""}
              onValueChange={field.onChange}
            >
              <SelectTrigger>
                <SelectValue placeholder={t("stagePlaceholder")} />
              </SelectTrigger>
              <SelectContent>
                {pipelines
                  .filter((stage) => stage.is_active)
                  .sort((a, b) => a.order - b.order)
                  .map((stage) => (
                    <SelectItem key={stage.id} value={stage.id}>
                      {stage.name}
                    </SelectItem>
                  ))}
              </SelectContent>
            </Select>
          )}
        />
        {errors.stage_id && <FieldError>{errors.stage_id.message}</FieldError>}
      </Field>

      <div className="grid grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("valueRequired")}</FieldLabel>
          <NumberInput
            value={computedValueRupiah}
            onChange={(val) => {
              // Allow manual input only when no product items are managing the value
              if (productItems.length === 0) {
                setValue("value", val, { shouldValidate: true });
              }
            }}
            placeholder={t("valuePlaceholder")}
            allowDecimal
            decimalPlaces={2}
            min={0}
            disabled={productItems.length > 0}
          />
          {productItems.length > 0 && (
            <p className="text-xs text-muted-foreground">{t("valueFromProducts") || "Derived from product items"}</p>
          )}
          {errors.value && <FieldError>{errors.value.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("probabilityLabel")}</FieldLabel>
          <Controller
            control={control}
            name="probability"
            render={({ field }) => (
              <Input
                type="number"
                value={field.value ?? stageProbability}
                onChange={(e) => field.onChange(Number(e.target.value))}
                placeholder={t("probabilityPlaceholder")}
                min={0}
                max={100}
              />
            )}
          />
          {errors.probability && <FieldError>{errors.probability.message}</FieldError>}
        </Field>
      </div>

      <Field orientation="vertical">
        <FieldLabel>{t("expectedCloseDateLabel")} *</FieldLabel>
        <DatePicker
          date={expectedCloseDate}
          onDateChange={(date) => {
            if (date) {
              // Convert Date to YYYY-MM-DD format using local time methods
              // Calendar component works with local dates, so we use local methods for consistency
              const year = date.getFullYear();
              const month = String(date.getMonth() + 1).padStart(2, "0");
              const day = String(date.getDate()).padStart(2, "0");
              setValue("expected_close_date", `${year}-${month}-${day}`);
            } else {
              setValue("expected_close_date", "");
            }
          }}
          placeholder={t("expectedCloseDateLabel")}
        />
        {errors.expected_close_date && <FieldError>{errors.expected_close_date.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("sourceLabel")}</FieldLabel>
        <Controller
          control={control}
          name="source"
          render={({ field }) => (
            <Select
              value={field.value || "none"}
              onValueChange={(value) => field.onChange(value === "none" ? "" : value)}
            >
              <SelectTrigger>
                <SelectValue placeholder={t("sourcePlaceholder")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">{t("sourceNone") || "None"}</SelectItem>
                {leadSources.map((source) => (
                  <SelectItem key={source.id} value={source.name}>
                    {source.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        />
        {errors.source && <FieldError>{errors.source.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("notesLabel")}</FieldLabel>
        <Textarea
          {...register("notes")}
          placeholder={t("notesPlaceholder")}
          rows={3}
        />
        {errors.notes && <FieldError>{errors.notes.message}</FieldError>}
      </Field>

      {/* Product Items */}
      <div className="space-y-3">
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium">{t("productItemsTitle")}</p>
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
                } as never)
              }
            >
              {t("addProductItem")}
            </Button>
          </div>

          {productItemsFieldArray.fields.length === 0 ? (
            <p className="text-sm text-muted-foreground">{t("productItemsEmpty")}</p>
          ) : (
            <div className="space-y-3">
              {productItemsFieldArray.fields.map((field, index) => (
                <div key={field.id} className="rounded-md border p-3">
                  <div className="grid grid-cols-1 gap-3 md:grid-cols-12">
                    <div className="md:col-span-5">
                      <Field orientation="vertical">
                        <FieldLabel>{t("productLabel")}</FieldLabel>
                        <Controller
                          control={control}
                          name={`product_items.${index}.product_id` as never}
                          render={({ field }: { field: ControllerRenderProps }) => (
                            <Select value={(field.value as unknown as string) ?? ""} onValueChange={field.onChange}>
                              <SelectTrigger>
                                <SelectValue placeholder={t("productPlaceholder")} />
                              </SelectTrigger>
                              <SelectContent>
                                {products.map((p) => (
                                  <SelectItem key={p.id} value={p.id}>
                                    {p.name} {p.sku ? `(${p.sku})` : ""}
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          )}
                        />
                        {errors.product_items?.[index]?.product_id && (
                          <FieldError>{errors.product_items[index].product_id.message}</FieldError>
                        )}
                      </Field>
                    </div>

                    <div className="md:col-span-2">
                      <Field orientation="vertical">
                        <FieldLabel>{t("quantityLabel")}</FieldLabel>
                        <Input
                          type="number"
                          {...register(`product_items.${index}.quantity` as never, { valueAsNumber: true })}
                          min={1}
                        />
                      </Field>
                    </div>

                    <div className="md:col-span-2">
                      <Field orientation="vertical">
                        <FieldLabel>{t("unitPriceLabel")}</FieldLabel>
                        <Controller
                          control={control}
                          name={`product_items.${index}.unit_price` as never}
                          render={({ field }: { field: ControllerRenderProps }) => (
                            <NumberInput
                              value={(field.value as unknown as number) ?? 0}
                              onChange={(value) => field.onChange(value ?? 0)}
                              onBlur={field.onBlur}
                              allowDecimal
                              decimalPlaces={2}
                              min={0}
                            />
                          )}
                        />
                      </Field>
                    </div>

                    <div className="md:col-span-2">
                      <Field orientation="vertical">
                        <FieldLabel>{t("discountLabel")}</FieldLabel>
                        <Controller
                          control={control}
                          name={`product_items.${index}.discount_amount` as never}
                          render={({ field }: { field: ControllerRenderProps }) => (
                            <NumberInput
                              value={(field.value as unknown as number) ?? 0}
                              onChange={(value) => field.onChange(value ?? 0)}
                              onBlur={field.onBlur}
                              allowDecimal
                              decimalPlaces={2}
                              min={0}
                            />
                          )}
                        />
                      </Field>
                    </div>

                    <div className="md:col-span-1 flex items-end justify-end">
                      <Button
                        type="button"
                        variant="outline"
                        size="icon"
                        className="cursor-pointer"
                        onClick={() => productItemsFieldArray.remove(index)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>

                    <div className="md:col-span-12">
                      <Field orientation="vertical">
                        <FieldLabel>{t("itemNotesLabel")}</FieldLabel>
                        <Input
                          {...register(`product_items.${index}.notes` as never)}
                          placeholder={t("itemNotesPlaceholder")}
                        />
                      </Field>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
  </div>

      <div className="flex justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading}>
          {t("cancel")}
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? t("submitting") : isEdit ? t("submitUpdate") : t("submitCreate")}
        </Button>
      </div>
      

    </form>
  );
}
