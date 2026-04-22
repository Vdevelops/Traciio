"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslations } from "next-intl";
import {
  createGroupTargetWithUserAssignmentSchema,
  type CreateGroupTargetWithUserAssignmentFormData,
} from "@/features/master-data/monthly-target/schemas/monthly-target.schema";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { formatCurrency } from "@/lib/utils";
import type { Group } from "../types";
import { Info } from "lucide-react";

interface GroupTargetFormProps {
  readonly group: Group;
  readonly onSubmit: (
    data: CreateGroupTargetWithUserAssignmentFormData
  ) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

export function GroupTargetForm({
  group,
  onSubmit,
  onCancel,
  isLoading,
}: GroupTargetFormProps) {
  const t = useTranslations("groupManagement.targetForm");

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreateGroupTargetWithUserAssignmentFormData>({
    resolver: zodResolver(createGroupTargetWithUserAssignmentSchema),
    defaultValues: {
      group_id: group.id,
      year: new Date().getFullYear(),
      month: new Date().getMonth() + 1,
      target_amount: 0,
    },
  });

  const targetAmount = watch("target_amount");
  const year = watch("year");
  const month = watch("month");

  const currentYear = new Date().getFullYear();
  const years = Array.from({ length: 11 }, (_, i) => currentYear - 5 + i);
  const months = [
    { value: 1, label: "Januari" },
    { value: 2, label: "Februari" },
    { value: 3, label: "Maret" },
    { value: 4, label: "April" },
    { value: 5, label: "Mei" },
    { value: 6, label: "Juni" },
    { value: 7, label: "Juli" },
    { value: 8, label: "Agustus" },
    { value: 9, label: "September" },
    { value: 10, label: "Oktober" },
    { value: 11, label: "November" },
    { value: 12, label: "Desember" },
  ];

  const handleFormSubmit = async (
    data: CreateGroupTargetWithUserAssignmentFormData
  ) => {
    // Ensure target_amount is a number and not NaN
    const submitData = {
      ...data,
      target_amount: Number(data.target_amount) || 0,
    };
    await onSubmit(submitData);
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-6">
      {/* Group Info Card */}
      <Card className="border-border/60 bg-card shadow-sm">
        <CardHeader className="pb-3">
          <CardTitle className="text-lg">{group.name}</CardTitle>
          <CardDescription className="text-sm">{group.code}</CardDescription>
        </CardHeader>
      </Card>

      {/* Period Selection */}
      <div className="grid grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("yearLabel")}</FieldLabel>
          <Select
            value={year?.toString() ?? ""}
            onValueChange={(value) =>
              setValue("year", Number.parseInt(value, 10), {
                shouldValidate: true,
              })
            }
            disabled={isLoading}
          >
            <SelectTrigger className="h-11">
              <SelectValue placeholder={t("yearPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              {years.map((y) => (
                <SelectItem key={y} value={y.toString()}>
                  {y}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.year && <FieldError>{errors.year.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("monthLabel")}</FieldLabel>
          <Select
            value={month?.toString() ?? ""}
            onValueChange={(value) =>
              setValue("month", Number.parseInt(value, 10), {
                shouldValidate: true,
              })
            }
            disabled={isLoading}
          >
            <SelectTrigger className="h-11">
              <SelectValue placeholder={t("monthPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              {months.map((m) => (
                <SelectItem key={m.value} value={m.value.toString()}>
                  {m.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.month && <FieldError>{errors.month.message}</FieldError>}
        </Field>
      </div>

      {/* Target Amount */}
      <Field orientation="vertical">
        <FieldLabel>{t("targetAmountLabel")}</FieldLabel>
        <Input
          type="number"
          {...register("target_amount", { 
            valueAsNumber: true,
            required: "Target amount is required",
            min: {
              value: 0,
              message: "Target amount must be non-negative",
            },
          })}
          placeholder={t("targetAmountPlaceholder")}
          disabled={isLoading}
          min={0}
          step={1000000}
          className="h-11"
        />
        {targetAmount > 0 && (
          <p className="text-sm text-muted-foreground mt-1.5">
            {formatCurrency(targetAmount)}
          </p>
        )}
        {errors.target_amount && (
          <FieldError>{errors.target_amount.message}</FieldError>
        )}
      </Field>

      {/* Info Message Card */}
      <Card className="border-accent/50 bg-accent/5">
        <CardContent>
          <div className="flex gap-3">
            <Info className="h-5 w-5 text-accent-foreground mt-0.5 shrink-0" />
            <p className="text-sm text-muted-foreground leading-relaxed">
              {t("infoMessage")}
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Actions */}
      <div className="flex justify-end gap-3 pt-2">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isLoading}
          className="cursor-pointer"
        >
          {t("cancel")}
        </Button>
        <Button 
          type="submit" 
          disabled={isLoading || !targetAmount || targetAmount <= 0} 
          className="cursor-pointer"
        >
          {isLoading ? t("savingButton") || "Saving..." : t("submit")}
        </Button>
      </div>
    </form>
  );
}
