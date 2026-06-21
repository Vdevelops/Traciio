"use client";

import { useForm } from "react-hook-form";
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
import { useLeadQualification } from "../hooks/useLeadQualification";
import type { ConvertLeadResponse, Lead } from "../types";
import { toast } from "sonner";
import { useTranslations } from "next-intl";
import { useEffect, useMemo } from "react";
import { Info } from "lucide-react";
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
  const activeStages = useMemo(() => {
    return (stages ?? []).filter((stage) => stage.is_active);
  }, [stages]);
  const convertibleStages = useMemo(() => {
    const wonStages = activeStages.filter((stage) => stage.is_won);
    return wonStages.length > 0 ? wonStages : activeStages;
  }, [activeStages]);
  const { qualification } = useLeadQualification(lead.id);

  const {
    register,
    handleSubmit,
    setValue,
    watch,
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
    },
  });

  useEffect(() => {
    if (open && convertibleStages.length > 0) {
      const sortedStages = [...convertibleStages].sort((a, b) => a.order - b.order);
      const defaultStage = sortedStages[0];
      const initialStageId = defaultStage?.id || "";
      const initialValue = lead.estimated_value || qualification?.budget_target_amount || undefined;

      reset({
        opportunity_title: `${lead.company_name || lead.first_name} - ${lead.industry || "Opportunity"}`,
        opportunity_description: lead.notes || "",
        stage_id: initialStageId,
        value: initialValue !== undefined ? initialValue / 100 : undefined,
        status_reason: "",
      });
    }
  }, [open, lead, convertibleStages, qualification, reset]);

  const onSubmit = async (data: ConvertLeadFormData) => {
    try {
      const payload = { ...data };

      if (payload.value !== undefined) {
        payload.value = payload.value * 100;
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
              value={watch("stage_id")}
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
              value={watch("value") || 0}
              onChange={(value) => setValue("value", value)}
              placeholder="0"
              min={0}
            />
            {errors.value && (
              <FieldError>{errors.value.message}</FieldError>
            )}
          </Field>

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
