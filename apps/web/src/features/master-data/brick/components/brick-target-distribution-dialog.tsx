"use client";

import { useMemo, useEffect } from "react";
import { useForm, useFieldArray, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslations } from "next-intl";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Badge } from "@/components/ui/badge";
import { Plus, Trash2, AlertCircle } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import {
  createBrickTargetDistributionSchema,
  type CreateBrickTargetDistributionFormData,
} from "../schemas/brick.schema";
import { useBrickSales } from "../hooks/useBricks";
import { useDistributeBrickTarget } from "../hooks/useBricks";
import { toast } from "sonner";
import type { BrickTargetWithDistributions } from "../types";

interface BrickTargetDistributionDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly brickId: string;
  readonly targetData?: BrickTargetWithDistributions;
  readonly onSuccess?: () => void;
}

export function BrickTargetDistributionDialog({
  open,
  onOpenChange,
  brickId,
  targetData,
  onSuccess,
}: BrickTargetDistributionDialogProps) {
  const t = useTranslations("brickTargetDistribution");
  const { data: salesData, isLoading: isLoadingSales } = useBrickSales(brickId);
  const distributeTarget = useDistributeBrickTarget();

  const sales = salesData?.data ?? [];
  const monthlyTarget = targetData?.target?.target_amount ?? 0;
  const existingDistributions = targetData?.distributions ?? [];
  const totalDistributed = targetData?.total_distributed ?? 0;
  const remainingAmount = targetData?.remaining_amount ?? 0;

  // Get sales that don't have distribution yet
  const salesWithoutDistribution = useMemo(() => {
    const distributedSalesIds = new Set(
      existingDistributions.map((d) => d.sales_user_id)
    );
    return sales.filter((s) => !distributedSalesIds.has(s.id));
  }, [sales, existingDistributions]);

  const {
    control,
    handleSubmit,
    watch,
    reset,
    setValue,
    formState: { errors },
  } = useForm<CreateBrickTargetDistributionFormData>({
    resolver: zodResolver(createBrickTargetDistributionSchema),
    defaultValues: {
      distributions: [],
    },
  });

  // Reset form when dialog opens/closes
  useEffect(() => {
    if (!open) {
      reset({ distributions: [] });
    }
  }, [open, reset]);

  const { fields, append, remove } = useFieldArray({
    control,
    name: "distributions",
  });

  const watchedDistributions = watch("distributions");
  const totalDistributedInForm = useMemo(() => {
    return watchedDistributions.reduce(
      (sum, dist) => sum + (dist.distributed_amount || 0),
      0
    );
  }, [watchedDistributions]);

  const availableAmount = remainingAmount;

  const formatCurrency = (value: number) => {
    // Value is in sen (smallest currency unit), convert to rupiah
    const rupiah = value / 100;
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(rupiah);
  };

  const handleAddDistribution = () => {
    if (salesWithoutDistribution.length === 0) {
      toast.error(t("noSalesAvailable"));
      return;
    }
    append({
      sales_user_id: salesWithoutDistribution[0].id,
      distributed_amount: 0,
    });
  };

  const handleRemoveDistribution = (index: number) => {
    remove(index);
  };

  const onSubmit = async (data: CreateBrickTargetDistributionFormData) => {
    if (!targetData?.target?.id) {
      toast.error(t("noTargetFound"));
      return;
    }

    // Validate total doesn't exceed remaining amount
    const total = data.distributions.reduce(
      (sum, dist) => sum + dist.distributed_amount,
      0
    );
    if (total > availableAmount) {
      toast.error(t("exceedsRemainingAmount"));
      return;
    }

    try {
      await distributeTarget.mutateAsync({
        brickId,
        targetId: targetData.target.id,
        data,
      });
      toast.success(t("distributionSuccess"));
      reset();
      onOpenChange(false);
      onSuccess?.();
    } catch (error: any) {
      toast.error(
        error?.response?.data?.error?.message || t("distributionError")
      );
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {/* Target Info */}
          <div className="grid gap-4 grid-cols-3 p-4 bg-muted/50 rounded-lg">
            <div>
              <p className="text-sm text-muted-foreground">{t("monthlyTarget")}</p>
              <p className="text-lg font-medium">{formatCurrency(monthlyTarget)}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">{t("totalDistributed")}</p>
              <p className="text-lg font-medium text-blue-600">
                {formatCurrency(totalDistributed)}
              </p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">{t("remainingAmount")}</p>
              <p className="text-lg font-medium text-yellow-600">
                {formatCurrency(remainingAmount)}
              </p>
            </div>
          </div>

          {/* Existing Distributions */}
          {existingDistributions.length > 0 && (
            <div className="space-y-2">
              <h4 className="text-sm font-medium">{t("existingDistributions")}</h4>
              <div className="space-y-2">
                {existingDistributions.map((dist) => (
                  <div
                    key={dist.id}
                    className="flex items-center justify-between p-3 border rounded-lg"
                  >
                    <div className="flex-1">
                      <p className="text-sm font-medium">
                        {dist.sales_user?.name ?? t("unknownSales")}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {formatCurrency(dist.distributed_amount)}
                      </p>
                    </div>
                    <Badge variant="outline">{t("distributed")}</Badge>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* New Distributions */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h4 className="text-sm font-medium">{t("newDistributions")}</h4>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={handleAddDistribution}
                disabled={salesWithoutDistribution.length === 0}
                className="cursor-pointer"
              >
                <Plus className="h-4 w-4 mr-2" />
                {t("addDistribution")}
              </Button>
            </div>

            {fields.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground border rounded-lg">
                <p className="text-sm">{t("noNewDistributions")}</p>
                <p className="text-xs mt-1">{t("clickAddToStart")}</p>
              </div>
            ) : (
              <ScrollArea className="h-64 rounded-md border p-4">
                <div className="space-y-4">
                  {fields.map((field, index) => {
                    const availableSales = salesWithoutDistribution.filter(
                      (s) =>
                        !watchedDistributions
                          .slice(0, index)
                          .some((d) => d.sales_user_id === s.id)
                    );
                    const selectedSalesId = watch(`distributions.${index}.sales_user_id`);
                    const selectedSales = sales.find((s) => s.id === selectedSalesId);

                    return (
                      <div
                        key={field.id}
                        className="flex items-start gap-4 p-4 border rounded-lg"
                      >
                        <div className="flex-1 space-y-4">
                          <Field>
                            <FieldLabel>{t("salesUser")}</FieldLabel>
                            <Controller
                              control={control}
                              name={`distributions.${index}.sales_user_id`}
                              render={({ field }) => (
                                <Select
                                  value={field.value || ""}
                                  onValueChange={field.onChange}
                                >
                                  <SelectTrigger>
                                    <SelectValue placeholder={t("selectSales")} />
                                  </SelectTrigger>
                                  <SelectContent>
                                    {availableSales.map((sale) => (
                                      <SelectItem key={sale.id} value={sale.id}>
                                        {sale.name} ({sale.email})
                                      </SelectItem>
                                    ))}
                                  </SelectContent>
                                </Select>
                              )}
                            />
                            {errors.distributions?.[index]?.sales_user_id && (
                              <FieldError>
                                {errors.distributions[index]?.sales_user_id?.message}
                              </FieldError>
                            )}
                          </Field>

                          <Field>
                            <FieldLabel>
                              {t("amount")} ({t("remaining")}:{" "}
                              {formatCurrency(
                                availableAmount -
                                  totalDistributedInForm +
                                  (watchedDistributions[index]?.distributed_amount || 0)
                              )}
                              )
                            </FieldLabel>
                            <Controller
                              control={control}
                              name={`distributions.${index}.distributed_amount`}
                              render={({ field }) => (
                                <Input
                                  type="number"
                                  min="0"
                                  step="1000000"
                                  placeholder="0"
                                  value={
                                    field.value
                                      ? Math.floor(field.value / 100).toString()
                                      : ""
                                  }
                                  onChange={(e) => {
                                    const value = e.target.value;
                                    // Convert rupiah input to sen (smallest currency unit)
                                    field.onChange(
                                      value ? parseInt(value, 10) * 100 : 0
                                    );
                                  }}
                                />
                              )}
                            />
                            {errors.distributions?.[index]?.distributed_amount && (
                              <FieldError>
                                {
                                  errors.distributions[index]?.distributed_amount
                                    ?.message
                                }
                              </FieldError>
                            )}
                          </Field>
                        </div>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          onClick={() => handleRemoveDistribution(index)}
                          className="cursor-pointer"
                        >
                          <Trash2 className="h-4 w-4 text-destructive" />
                        </Button>
                      </div>
                    );
                  })}
                </div>
              </ScrollArea>
            )}

            {/* Total Validation */}
            {totalDistributedInForm > 0 && (
              <Alert>
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>
                  <div className="flex items-center justify-between">
                    <span>{t("totalInForm")}:</span>
                    <span className="font-medium">
                      {formatCurrency(totalDistributedInForm)}
                    </span>
                  </div>
                  {totalDistributedInForm > availableAmount && (
                    <p className="text-sm text-destructive mt-2">
                      {t("exceedsRemainingWarning")}
                    </p>
                  )}
                </AlertDescription>
              </Alert>
            )}
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                reset();
                onOpenChange(false);
              }}
              className="cursor-pointer"
            >
              {t("cancel")}
            </Button>
            <Button
              type="submit"
              disabled={
                distributeTarget.isPending ||
                fields.length === 0 ||
                totalDistributedInForm > availableAmount ||
                totalDistributedInForm === 0
              }
              className="cursor-pointer"
            >
              {distributeTarget.isPending ? t("distributing") : t("distribute")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

