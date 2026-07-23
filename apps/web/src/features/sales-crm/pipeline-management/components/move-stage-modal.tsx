"use client";

import { useEffect } from "react";
import { useForm, Controller, useFieldArray, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import { AlertCircle, Trash2 } from "lucide-react";
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
import { Label } from "@/components/ui/label";
import { NumberInput } from "@/components/ui/number-input";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { moveStageSchema, type MoveStageFormData } from "../schemas/deal.schema";
import { useMoveStage } from "../hooks/useDeals";
import { useProducts } from "@/features/sales-crm/product-management/hooks/useProducts";
import type { PipelineStage } from "../types";

interface MoveStageModalProps {
  readonly dealId: string;
  readonly currentStageId: string;
  readonly availableStages: readonly PipelineStage[];
  readonly isOpen: boolean;
  readonly initialStageId?: string;
  readonly onClose: () => void;
  readonly onSuccess?: () => void;
}

export function MoveStageModal({
  dealId,
  currentStageId,
  availableStages,
  isOpen,
  initialStageId,
  onClose,
  onSuccess,
}: MoveStageModalProps) {
  const t = useTranslations("pipelineManagement.statusReason");
  const moveStageMutation = useMoveStage();
  const { data: productsData, refetch: refetchProducts } = useProducts({ per_page: 100, status: "active" });
  const products = productsData?.data ?? [];

  const {
    register,
    control,
    handleSubmit,
    reset,
    setValue,
    formState: { errors },
  } = useForm<MoveStageFormData>({
    resolver: zodResolver(moveStageSchema),
    defaultValues: {
      to_stage_id: "",
      reason: "",
      notes: "",
      product_items: [],
    },
  });

  const selectedStageId = useWatch({ control, name: "to_stage_id" });
  const reasonValue = useWatch({ control, name: "reason" });
  const productItems = useWatch({ control, name: "product_items" });
  const productItemsFieldArray = useFieldArray({
    control,
    name: "product_items",
  });

  // Filter out current stage and get only next valid stages
  const nextStages = availableStages.filter(
    (stage) => stage.id !== currentStageId
  );

  const selectedStage = availableStages.find((s) => s.id === selectedStageId);
  const requiresReason = Boolean(selectedStage?.is_won || selectedStage?.is_lost);
  const requiresSoldProducts = Boolean(selectedStage?.is_won);

  useEffect(() => {
    if (isOpen) {
      void refetchProducts();
    }
    if (isOpen && initialStageId) {
      setValue("to_stage_id", initialStageId, { shouldValidate: true });
    }
  }, [initialStageId, isOpen, refetchProducts, setValue]);

  useEffect(() => {
    const currentProductItems = productItems ?? [];
    if (!requiresSoldProducts && currentProductItems.length > 0) {
      setValue("product_items", [], { shouldValidate: false });
    }
    if (requiresSoldProducts && currentProductItems.length === 0) {
      productItemsFieldArray.append({
        product_id: "",
        quantity: 1,
        unit_price: 0,
        discount_amount: 0,
        notes: "",
      });
    }
  }, [productItems, productItemsFieldArray, requiresSoldProducts, setValue]);

  const onSubmit = async (data: MoveStageFormData) => {
    if (requiresReason && !data.reason?.trim()) {
      return;
    }
    const normalizedItems = (data.product_items ?? [])
      .filter((item) => item.product_id && item.quantity > 0)
      .map((item) => ({
        ...item,
        unit_price: Math.round((item.unit_price ?? 0) * 100),
        discount_amount: Math.round((item.discount_amount ?? 0) * 100),
      }));
    if (requiresSoldProducts && normalizedItems.length === 0) {
      toast.error("Produk deal wajib diisi sebelum Closed Won");
      return;
    }

    try {
      await moveStageMutation.mutateAsync({
        deal_id: dealId,
        stage_id: data.to_stage_id,
        reason: data.reason?.trim(),
        product_items: requiresSoldProducts ? normalizedItems : [],
      });

      toast.success(t("successTitle"), {
        description: `Moved to ${selectedStage?.name ?? "new stage"}`,
      });

      reset();
      onClose();
      onSuccess?.();
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to move deal stage";
      toast.error(t("errorTitle"), {
        description: errorMessage,
      });
    }
  };

  const handleClose = () => {
    if (!moveStageMutation.isPending) {
      reset();
      onClose();
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[720px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Move Deal to Different Stage</DialogTitle>
          <DialogDescription>
            Select the next stage for this deal and provide a reason for the
            movement.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {nextStages.length === 0 ? (
            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                No available stages to move to. This deal might already be in
                the final stage.
              </AlertDescription>
            </Alert>
          ) : (
            <>
              <div className="space-y-2">
                <Label htmlFor="to_stage_id">Target Stage *</Label>
                <Controller
                  name="to_stage_id"
                  control={control}
                  render={({ field }) => (
                    <Select
                      onValueChange={field.onChange}
                      value={field.value}
                    >
                      <SelectTrigger id="to_stage_id">
                        <SelectValue placeholder="Select target stage" />
                      </SelectTrigger>
                      <SelectContent>
                        {nextStages.map((stage) => (
                          <SelectItem
                            key={stage.id}
                            value={stage.id}
                            className="cursor-pointer"
                          >
                            <div className="flex items-center space-x-2">
                              <span>{stage.name}</span>
                              {stage.probability !== undefined &&
                                stage.probability !== null && (
                                  <span className="text-xs text-muted-foreground">
                                    ({stage.probability}% probability)
                                  </span>
                                )}
                            </div>
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                />
                {errors.to_stage_id?.message && (
                  <p className="text-sm text-destructive">
                    {errors.to_stage_id.message}
                  </p>
                )}
              </div>

              {selectedStage?.requirements && (
                <Alert>
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    <div className="font-medium mb-1">
                      Stage Requirements:
                    </div>
                    <div className="text-sm">
                      {selectedStage.requirements}
                    </div>
                  </AlertDescription>
                </Alert>
              )}

              {requiresSoldProducts && (
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <Label>Produk Deal</Label>
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

                  {productItemsFieldArray.fields.length === 0 ? (
                    <p className="text-sm text-muted-foreground">
                      Tambahkan produk yang benar-benar terjual untuk menutup deal.
                    </p>
                  ) : (
                    <div className="space-y-3">
                      {productItemsFieldArray.fields.map((field, index) => (
                        <div key={field.id} className="rounded-md border p-3">
                          <div className="grid grid-cols-1 gap-3 md:grid-cols-12">
                            <div className="md:col-span-5">
                              <Label>Produk</Label>
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
                            </div>
                            <div className="md:col-span-2">
                              <Label>Qty</Label>
                              <Input
                                type="number"
                                min={1}
                                {...register(`product_items.${index}.quantity`, { valueAsNumber: true })}
                              />
                            </div>
                            <div className="md:col-span-2">
                              <Label>Harga</Label>
                              <Controller
                                control={control}
                                name={`product_items.${index}.unit_price`}
                                render={({ field }) => (
                                  <NumberInput
                                    value={field.value ?? 0}
                                    onChange={(value) => field.onChange(value ?? 0)}
                                    min={0}
                                  />
                                )}
                              />
                            </div>
                            <div className="md:col-span-2">
                              <Label>Diskon</Label>
                              <Controller
                                control={control}
                                name={`product_items.${index}.discount_amount`}
                                render={({ field }) => (
                                  <NumberInput
                                    value={field.value ?? 0}
                                    onChange={(value) => field.onChange(value ?? 0)}
                                    min={0}
                                  />
                                )}
                              />
                            </div>
                            <div className="flex items-end justify-end md:col-span-1">
                              <Button
                                type="button"
                                variant="outline"
                                size="icon"
                                onClick={() => productItemsFieldArray.remove(index)}
                              >
                                <Trash2 className="h-4 w-4" />
                              </Button>
                            </div>
                            <div className="md:col-span-12">
                              <Label>Catatan</Label>
                              <Input
                                {...register(`product_items.${index}.notes`)}
                                placeholder="Catatan produk"
                              />
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}

                  {errors.product_items?.message && (
                    <p className="text-sm text-destructive">{errors.product_items.message}</p>
                  )}
                </div>
              )}

              <div className="space-y-2">
                <Label htmlFor="reason">{t("label")}</Label>
                <Controller
                  name="reason"
                  control={control}
                  render={({ field }) => (
                    <Textarea
                      {...field}
                      id="reason"
                      placeholder={t("placeholder")}
                      className="min-h-[80px] resize-none"
                    />
                  )}
                />
                {errors.reason?.message && (
                  <p className="text-sm text-destructive">
                    {errors.reason.message}
                  </p>
                )}
                {requiresReason && !reasonValue?.trim() && (
                  <p className="text-sm text-destructive">{t("required")}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="notes">Additional Notes</Label>
                <Controller
                  name="notes"
                  control={control}
                  render={({ field }) => (
                    <Textarea
                      {...field}
                      value={field.value ?? ""}
                      id="notes"
                      placeholder="Any additional notes about this stage transition..."
                      className="min-h-[80px] resize-none"
                    />
                  )}
                />
                {errors.notes?.message && (
                  <p className="text-sm text-destructive">
                    {errors.notes.message}
                  </p>
                )}
              </div>
            </>
          )}

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={moveStageMutation.isPending}
            >
              {t("cancel")}
            </Button>
            <Button
              type="submit"
              disabled={
                moveStageMutation.isPending || nextStages.length === 0
              }
              className="cursor-pointer"
            >
              {moveStageMutation.isPending ? "Moving..." : t("confirm")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
