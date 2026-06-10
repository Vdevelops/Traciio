"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  createLeadSchema,
  updateLeadSchema,
  type CreateLeadFormData,
  type UpdateLeadFormData,
} from "../schemas/lead.schema";
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
import { useLeadFormData } from "../hooks/useLeads";
import { useAllLeadStatuses } from "../hooks/useLeadStatuses";
import { useAllIndustries } from "../hooks/useIndustries";
import { useAllLeadSources } from "../hooks/useLeadSources";
import type { Lead } from "../types";
import { useEffect, useMemo } from "react";
import { useWatch } from "react-hook-form";
import { useTranslations } from "next-intl";
import { Skeleton } from "@/components/ui/skeleton";
import type { LeadStatus } from "../types/lead-status";

interface LeadFormProps {
  readonly lead?: Lead;
  readonly onSubmit: (data: CreateLeadFormData | UpdateLeadFormData) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

export function LeadForm({ lead, onSubmit, onCancel, isLoading }: LeadFormProps) {
  const isEdit = !!lead;
  const { data: formData, isLoading: isLoadingFormData } = useLeadFormData();
  const t = useTranslations("leadManagement.leadForm.fields");
  const tButtons = useTranslations("leadManagement.leadForm.buttons");

  // Fetch industries and lead sources from dedicated APIs (dynamic data)
  const { data: industriesData } = useAllIndustries();
  const { data: leadSourcesData } = useAllLeadSources();
  
  // Lead sources: Use dynamic data from API, fallback to form-data if needed
  const leadSourcesFromAPI = leadSourcesData?.data?.map((ls) => ({
    value: ls.code,
    label: ls.name,
  })) ?? [];
  const leadSources = leadSourcesFromAPI.length > 0 ? leadSourcesFromAPI : (formData?.data?.lead_sources ?? []);
  
  // Industries: Use dynamic data from API, fallback to form-data if needed
  const industriesFromAPI = industriesData?.data?.map((ind) => ind.name) ?? [];
  const industries = industriesFromAPI.length > 0 ? industriesFromAPI : (formData?.data?.industries ?? []);
  const provinces = formData?.data?.provinces ?? [];
  const defaults = formData?.data?.defaults;

  // Fetch full lead statuses (id + code + name)
  const { data: allLeadStatusesData } = useAllLeadStatuses();
  const allLeadStatuses: LeadStatus[] = useMemo(() => allLeadStatusesData?.data ?? [], [allLeadStatusesData]);
  const leadStatusOptions: Array<{ value: string; label: string; code: string }> = allLeadStatuses
    .filter((s: LeadStatus) => s.is_active)
    .sort((a: LeadStatus, b: LeadStatus) => (a.order ?? 0) - (b.order ?? 0))
    .map((s: LeadStatus) => ({ value: s.id, label: s.name, code: s.code }));

  const {
    register,
    handleSubmit,
    setValue,
    control,
    formState: { errors },
  } = useForm<CreateLeadFormData | UpdateLeadFormData>({
    resolver: zodResolver(isEdit ? updateLeadSchema : createLeadSchema),
    defaultValues: lead
      ? {
          first_name: lead.first_name,
          last_name: lead.last_name || "",
          company_name: lead.company_name || "",
          email: lead.email,
          phone: lead.phone || "",
          job_title: lead.job_title || "",
          industry: lead.industry || "",
          lead_source: lead.lead_source,
          // If API returns lead_status_id, keep it; else leave undefined
          lead_status_id: (lead as Lead & { lead_status_id?: string }).lead_status_id || undefined,
          notes: lead.notes || "",
          address: lead.address || "",
          city: lead.city || "",
          province: lead.province || "",
          postal_code: lead.postal_code || "",
          country: lead.country || "",
          website: lead.website || "",
          budget_confirmed: lead.budget_confirmed ?? false,
          budget_amount: lead.budget_amount ?? undefined,
          authority_confirmed: lead.authority_confirmed ?? false,
          authority_person: lead.authority_person || "",
          need_confirmed: lead.need_confirmed ?? false,
          need_description: lead.need_description || "",
          timeline_confirmed: lead.timeline_confirmed ?? false,
          probability: lead.probability ?? undefined,
          estimated_value: lead.estimated_value ?? undefined,
          expected_close_date: lead.expected_close_date || "",
        }
      : {
          // Default status resolved server-side; if we already have statuses list, choose default.
          lead_status_id: allLeadStatuses.find((s) => s.is_default)?.id,
          country: defaults?.country || "Indonesia",
          budget_confirmed: false,
          authority_confirmed: false,
          need_confirmed: false,
          timeline_confirmed: false,
        },
  });

  const industryValue = useWatch({ control, name: "industry" });
  const leadSourceValue = useWatch({ control, name: "lead_source" });
  const leadStatusIdValue = useWatch({ control, name: "lead_status_id" });
  const provinceValue = useWatch({ control, name: "province" });

  useEffect(() => {
    if (!isEdit && defaults) {
      // Keep best-effort default selection
      const defaultStatusId = allLeadStatuses.find((s) => s.is_default)?.id;
      if (defaultStatusId) setValue("lead_status_id", defaultStatusId);
      setValue("country", defaults.country);
    }
  }, [defaults, isEdit, setValue, allLeadStatuses]);

  const handleFormSubmit = async (data: CreateLeadFormData | UpdateLeadFormData) => {
    const submitData = { ...data };
    delete submitData.assigned_to;
    await onSubmit(submitData);
  };

  if (isLoadingFormData) {
    return (
      <div className="space-y-4">
        {[...Array(5)].map((_, i) => (
          <Skeleton key={i} className="h-10 w-full" />
        ))}
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>
            {t("firstNameLabel")} *
          </FieldLabel>
          <Input {...register("first_name")} placeholder={t("firstNamePlaceholder")} />
          {errors.first_name && <FieldError>{errors.first_name.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("lastNameLabel")}</FieldLabel>
          <Input {...register("last_name")} placeholder={t("lastNamePlaceholder")} />
          {errors.last_name && <FieldError>{errors.last_name.message}</FieldError>}
        </Field>
      </div>

      <Field orientation="vertical">
        <FieldLabel>{t("emailLabel")} *</FieldLabel>
        <Input type="email" {...register("email")} placeholder={t("emailPlaceholder")} />
        {errors.email && <FieldError>{errors.email.message}</FieldError>}
      </Field>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("phoneLabel")}</FieldLabel>
          <Input {...register("phone")} placeholder={t("phonePlaceholder")} />
          {errors.phone && <FieldError>{errors.phone.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("companyNameLabel")}</FieldLabel>
          <Input {...register("company_name")} placeholder={t("companyNamePlaceholder")} />
          {errors.company_name && <FieldError>{errors.company_name.message}</FieldError>}
        </Field>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("jobTitleLabel")}</FieldLabel>
          <Input {...register("job_title")} placeholder={t("jobTitlePlaceholder")} />
          {errors.job_title && <FieldError>{errors.job_title.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("industryLabel")}</FieldLabel>
            <Select
            value={industryValue || undefined}
            onValueChange={(value) => setValue("industry", value || undefined)}
          >
            <SelectTrigger>
              <SelectValue placeholder={t("industryPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              {industries.map((industry) => (
                <SelectItem key={industry} value={industry}>
                  {industry}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.industry && <FieldError>{errors.industry.message}</FieldError>}
        </Field>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>
            {t("leadSourceLabel")} *
          </FieldLabel>
          <Select
            value={leadSourceValue || ""}
            onValueChange={(value) => setValue("lead_source", value)}
          >
            <SelectTrigger>
              <SelectValue placeholder={t("leadSourcePlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              {leadSources.map((source) => (
                <SelectItem key={source.value} value={source.value}>
                  {source.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.lead_source && <FieldError>{errors.lead_source.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("leadStatusLabel")}</FieldLabel>
          <Select
            value={leadStatusIdValue || ""}
            onValueChange={(value) => setValue("lead_status_id", value || undefined)}
          >
            <SelectTrigger>
              <SelectValue placeholder={t("leadStatusPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              {leadStatusOptions.map((status) => (
                <SelectItem key={status.value} value={status.value}>
                  {status.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.lead_status_id && <FieldError>{errors.lead_status_id.message as string}</FieldError>}
        </Field>
      </div>

      <Field orientation="vertical">
        <FieldLabel>{t("addressLabel")}</FieldLabel>
        <Textarea {...register("address")} placeholder={t("addressPlaceholder")} rows={3} />
        {errors.address && <FieldError>{errors.address.message}</FieldError>}
      </Field>

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("cityLabel")}</FieldLabel>
          <Input {...register("city")} placeholder={t("cityPlaceholder")} />
          {errors.city && <FieldError>{errors.city.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("provinceLabel")}</FieldLabel>
          <Select
            value={provinceValue || undefined}
            onValueChange={(value) => setValue("province", value || undefined)}
          >
            <SelectTrigger>
              <SelectValue placeholder={t("provincePlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              {provinces.map((province) => (
                <SelectItem key={province} value={province}>
                  {province}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.province && <FieldError>{errors.province.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("postalCodeLabel")}</FieldLabel>
          <Input {...register("postal_code")} placeholder={t("postalCodePlaceholder")} />
          {errors.postal_code && <FieldError>{errors.postal_code.message}</FieldError>}
        </Field>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("countryLabel")}</FieldLabel>
          <Input {...register("country")} placeholder={t("countryPlaceholder")} />
          {errors.country && <FieldError>{errors.country.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("websiteLabel")}</FieldLabel>
          <Input {...register("website")} placeholder={t("websitePlaceholder")} type="url" />
          {errors.website && <FieldError>{errors.website.message}</FieldError>}
        </Field>
      </div>

      <Field orientation="vertical">
        <FieldLabel>{t("notesLabel")}</FieldLabel>
        <Textarea {...register("notes")} placeholder={t("notesPlaceholder")} rows={4} />
        {errors.notes && <FieldError>{errors.notes.message}</FieldError>}
      </Field>

      <div className="flex flex-col-reverse sm:flex-row justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading} className="w-full sm:w-auto">
          {tButtons("cancel")}
        </Button>
        <Button type="submit" disabled={isLoading} className="w-full sm:w-auto">
          {isLoading
            ? tButtons("submitting")
            : isEdit
              ? tButtons("submitUpdate")
              : tButtons("submitCreate")}
        </Button>
      </div>
    </form>
  );
}
