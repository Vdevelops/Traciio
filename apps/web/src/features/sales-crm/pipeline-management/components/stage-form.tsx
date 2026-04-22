"use client";

import { useMemo } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { NumberInput } from "@/components/ui/number-input";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { createStageSchema, updateStageSchema, type CreateStageFormData, type UpdateStageFormData } from "../schemas/pipeline.schema";
import type { PipelineStage } from "../types";
import { useTranslations } from "next-intl";
import { usePipelines } from "../hooks/usePipelines";

interface StageFormProps {
  readonly stage?: PipelineStage;
  readonly onSubmit: (data: CreateStageFormData | UpdateStageFormData) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

export function StageForm({ stage, onSubmit, onCancel, isLoading = false }: StageFormProps) {
  const isEdit = !!stage;
  const schema = isEdit ? updateStageSchema : createStageSchema;
  const t = useTranslations("pipelineManagement.stages");
  const { data: stagesData } = usePipelines();
  const stages = stagesData?.data || [];
  
  // Calculate default order on client side only
  const defaultOrder = useMemo(() => {
    if (stage) {
      return stage.order;
    }
    if (stages.length > 0) {
      return Math.max(...stages.map((s) => s.order)) + 1;
    }
    return 1;
  }, [stage, stages]);
  
  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<CreateStageFormData | UpdateStageFormData>({
    resolver: zodResolver(schema),
    defaultValues: stage
      ? {
          name: stage.name,
          code: stage.code,
          order: stage.order,
          color: stage.color,
          is_active: stage.is_active,
          is_won: stage.is_won,
          is_lost: stage.is_lost,
          probability: stage.probability ?? undefined,
          description: stage.description || "",
        }
      : {
          order: defaultOrder,
          color: "#3B82F6",
          is_active: true,
          is_won: false,
          is_lost: false,
          probability: undefined,
        },
  });

  const isActive = watch("is_active");
  const isWon = watch("is_won");
  const isLost = watch("is_lost");
  const colorValue = watch("color");

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <Field>
        <FieldLabel htmlFor="name">{t("nameRequired")}</FieldLabel>
        <Input
          id="name"
          {...register("name")}
          placeholder={t("namePlaceholder")}
          disabled={isLoading}
        />
        {errors.name && <FieldError>{errors.name.message}</FieldError>}
      </Field>

      <Field>
        <FieldLabel htmlFor="code">{t("codeRequired")}</FieldLabel>
        <Input
          id="code"
          {...register("code")}
          placeholder={t("codePlaceholder")}
          disabled={isLoading || isEdit}
        />
        {isEdit && (
          <p className="text-xs text-muted-foreground mt-1">
            Code cannot be changed after creation
          </p>
        )}
        {errors.code && <FieldError>{errors.code.message}</FieldError>}
      </Field>

      <Field>
        <FieldLabel htmlFor="order">{t("orderRequired")}</FieldLabel>
        <Input
          id="order"
          type="number"
          {...register("order", { valueAsNumber: true })}
          disabled={isLoading}
          min={0}
        />
        {errors.order && <FieldError>{errors.order.message}</FieldError>}
      </Field>

      <Field>
        <FieldLabel htmlFor="color">{t("colorLabel")}</FieldLabel>
        <div className="flex items-center gap-2">
          <Input
            id="color-picker"
            type="color"
            value={colorValue || "#3B82F6"}
            onChange={(e) => setValue("color", e.target.value, { shouldValidate: true })}
            disabled={isLoading}
            className="w-20 h-10 cursor-pointer"
          />
          <Input
            id="color"
            {...register("color")}
            placeholder={t("colorPlaceholder")}
            disabled={isLoading}
            className="flex-1"
          />
        </div>
        {errors.color && <FieldError>{errors.color.message}</FieldError>}
      </Field>

      <div className="flex items-center space-x-3 rounded-md border p-4">
        <Checkbox
          id="is_active"
          checked={isActive}
          onCheckedChange={(checked) => setValue("is_active", checked === true)}
          disabled={isLoading}
        />
        <div className="space-y-1 leading-none">
          <Label htmlFor="is_active" className="cursor-pointer">
            {t("isActiveLabel")}
          </Label>
          <p className="text-xs text-muted-foreground">
            Active stages are visible and available for deals
          </p>
        </div>
      </div>

      <div className="flex items-center space-x-3 rounded-md border p-4">
        <Checkbox
          id="is_won"
          checked={isWon}
          onCheckedChange={(checked) => setValue("is_won", checked === true)}
          disabled={isLoading}
        />
        <div className="space-y-1 leading-none">
          <Label htmlFor="is_won" className="cursor-pointer">
            {t("isWonLabel")}
          </Label>
          <p className="text-xs text-muted-foreground">
            Mark this stage as a winning stage (deals are considered won)
          </p>
        </div>
      </div>

      <div className="flex items-center space-x-3 rounded-md border p-4">
        <Checkbox
          id="is_lost"
          checked={isLost}
          onCheckedChange={(checked) => setValue("is_lost", checked === true)}
          disabled={isLoading}
        />
        <div className="space-y-1 leading-none">
          <Label htmlFor="is_lost" className="cursor-pointer">
            {t("isLostLabel")}
          </Label>
          <p className="text-xs text-muted-foreground">
            Mark this stage as a lost stage (deals are considered lost)
          </p>
        </div>
      </div>

      <Field>
        <FieldLabel htmlFor="probability">{t("probabilityLabel")}</FieldLabel>
        <NumberInput
          id="probability"
          value={watch("probability") ?? undefined}
          onChange={(value) => setValue("probability", value ?? undefined, { shouldValidate: true })}
          placeholder={t("probabilityPlaceholder")}
          min={0}
          max={100}
          disabled={isLoading}
        />
        {errors.probability && <FieldError>{errors.probability.message}</FieldError>}
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

      <div className="flex justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading}>
          {t("cancel")}
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? t("saving") : t("save")}
        </Button>
      </div>
    </form>
  );
}

