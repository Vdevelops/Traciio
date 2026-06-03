"use client";

import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { createActivitySchema, type CreateActivityFormData } from "../schemas/activity.schema";
import { useCreateActivity } from "../hooks/useVisitReports";
import { useActivityTypes } from "../hooks/useActivityTypes";
import { toast } from "sonner";
import { useEffect, useMemo, useState } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { useTranslations } from "next-intl";
import { Plus, X } from "lucide-react";
import { Badge } from "@/components/ui/badge";

interface ProductInterest {
  product_name: string;
  product_id?: string;
  interest_level: number; // 1-5
  quantity: number;
  price: number;
}

interface CreateActivityWithProductsDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly accountId?: string;
  readonly contactId?: string;
  readonly dealId?: string;
  readonly leadId?: string;
  readonly onSuccess?: () => void;
  readonly showProductInterests?: boolean;
}

export function CreateActivityWithProductsDialog({
  open,
  onOpenChange,
  accountId,
  contactId,
  dealId,
  leadId,
  onSuccess,
  showProductInterests = false,
}: CreateActivityWithProductsDialogProps) {
  const t = useTranslations("createActivityDialog");
  const createActivity = useCreateActivity();
  const { data: activityTypesData, isLoading: isLoadingTypes } = useActivityTypes({
    status: "active",
  });

  const activityTypes = useMemo(() => {
    return activityTypesData?.data ?? [];
  }, [activityTypesData]);

  const [productInterests, setProductInterests] = useState<ProductInterest[]>([]);
  const [currentProduct, setCurrentProduct] = useState<ProductInterest>({
    product_name: "",
    interest_level: 3,
    quantity: 1,
    price: 0,
  });

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
      timestamp: new Date().toISOString(),
    },
  });

  // Reset form when dialog opens
  useEffect(() => {
    if (open) {
      reset({
        activity_type_id: "",
        account_id: accountId,
        contact_id: contactId,
        deal_id: dealId,
        lead_id: leadId,
        description: "",
        timestamp: new Date().toISOString(),
      });
      setProductInterests([]);
      setCurrentProduct({
        product_name: "",
        interest_level: 3,
        quantity: 1,
        price: 0,
      });
    }
  }, [open, accountId, contactId, dealId, leadId, reset]);

  // Set default activity type when types are loaded
  useEffect(() => {
    const currentTypeId = watch("activity_type_id");
    if (activityTypes.length > 0 && !currentTypeId) {
      setValue("activity_type_id", activityTypes[0].id);
    }
  }, [activityTypes, watch, setValue]);

  const addProductInterest = () => {
    if (!currentProduct.product_name.trim()) {
      toast.error("Please enter product name");
      return;
    }
    if (currentProduct.quantity <= 0) {
      toast.error("Quantity must be greater than 0");
      return;
    }

    setProductInterests([...productInterests, { ...currentProduct }]);
    setCurrentProduct({
      product_name: "",
      interest_level: 3,
      quantity: 1,
      price: 0,
    });
    toast.success("Product added");
  };

  const removeProductInterest = (index: number) => {
    setProductInterests(productInterests.filter((_, i) => i !== index));
  };

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
        metadata: {},
      };

      // Add product interests to metadata if available
      if (productInterests.length > 0 && showProductInterests) {
        payload.metadata = {
          product_interests: productInterests,
        };
      }

      // Include account_id if available
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

      await createActivity.mutateAsync(payload);
      toast.success("Activity created successfully");
      const defaultTypeId = activityTypes.length > 0 ? activityTypes[0]?.id : "";
      reset({
        activity_type_id: defaultTypeId,
        account_id: accountId,
        contact_id: contactId,
        deal_id: dealId,
        lead_id: leadId,
        description: "",
        timestamp: new Date().toISOString(),
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
      <DialogContent className={`${showProductInterests ? "max-w-3xl" : "max-w-2xl"} max-h-[90vh] overflow-y-auto`}>
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Field orientation="vertical">
            <FieldLabel>{t("activityTypeLabel")} *</FieldLabel>
            {isLoadingTypes ? (
              <Skeleton className="h-10 w-full" />
            ) : (
              <Select
                value={watch("activity_type_id") ?? ""}
                onValueChange={(value) =>
                  setValue("activity_type_id", value, { shouldValidate: true })
                }
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
            {errors.activity_type_id && (
              <FieldError>{errors.activity_type_id.message}</FieldError>
            )}
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

          <Field orientation="vertical">
            <FieldLabel>Date & Time *</FieldLabel>
            <Input
              type="datetime-local"
              value={
                watch("timestamp")
                  ? new Date(watch("timestamp")).toISOString().slice(0, 16)
                  : new Date().toISOString().slice(0, 16)
              }
              onChange={(e) => {
                const value = e.target.value;
                if (value) {
                  const date = new Date(value);
                  setValue("timestamp", date.toISOString(), { shouldValidate: true });
                }
              }}
            />
            {errors.timestamp && <FieldError>{errors.timestamp.message}</FieldError>}
          </Field>

          {/* Product Interests Section */}
          {showProductInterests && (
            <div className="space-y-3 border-t pt-4">
              <div className="flex items-center justify-between">
                <FieldLabel className="mb-0">{t("productInterest.sectionTitle")}</FieldLabel>
                {productInterests.length > 0 && (
                  <Badge variant="secondary">{productInterests.length} added</Badge>
                )}
              </div>

              {/* Add Product Form */}
              <Card className="bg-muted/30">
                <CardContent className="pt-4 space-y-3">
                  <div className="grid grid-cols-2 gap-3">
                    <Field orientation="vertical">
                      <FieldLabel className="text-xs">{t("productInterest.productNameLabel")} *</FieldLabel>
                      <Input
                        placeholder={t("productInterest.productNamePlaceholder")}
                        value={currentProduct.product_name}
                        onChange={(e) =>
                          setCurrentProduct({
                            ...currentProduct,
                            product_name: e.target.value,
                          })
                        }
                      />
                    </Field>

                    <Field orientation="vertical">
                      <FieldLabel className="text-xs">{t("productInterest.interestLevelLabel")} *</FieldLabel>
                      <Select
                        value={String(currentProduct.interest_level)}
                        onValueChange={(value) =>
                          setCurrentProduct({
                            ...currentProduct,
                            interest_level: parseInt(value),
                          })
                        }
                      >
                        <SelectTrigger className="text-xs">
                          <SelectValue placeholder={t("productInterest.selectInterestLevel")} />
                        </SelectTrigger>
                        <SelectContent>
                          {[1, 2, 3, 4, 5].map((level) => (
                            <SelectItem key={level} value={String(level)}>
                              {"⭐".repeat(level)}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </Field>
                  </div>

                  <div className="grid grid-cols-2 gap-3">
                    <Field orientation="vertical">
                      <FieldLabel className="text-xs">{t("productInterest.quantityLabel")} *</FieldLabel>
                      <Input
                        type="number"
                        min="1"
                        placeholder={t("productInterest.quantityPlaceholder")}
                        value={currentProduct.quantity}
                        onChange={(e) =>
                          setCurrentProduct({
                            ...currentProduct,
                            quantity: parseInt(e.target.value) || 1,
                          })
                        }
                      />
                    </Field>

                    <Field orientation="vertical">
                      <FieldLabel className="text-xs">{t("productInterest.priceLabel")}</FieldLabel>
                      <Input
                        type="number"
                        min="0"
                        placeholder={t("productInterest.pricePlaceholder")}
                        value={currentProduct.price}
                        onChange={(e) =>
                          setCurrentProduct({
                            ...currentProduct,
                            price: parseInt(e.target.value) || 0,
                          })
                        }
                      />
                    </Field>
                  </div>

                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={addProductInterest}
                    className="w-full"
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    {t("productInterest.addButton")}
                  </Button>
                </CardContent>
              </Card>

              {/* Product List */}
              {productInterests.length > 0 && (
                <div className="space-y-2">
                  {productInterests.map((product, index) => (
                    <Card key={index} className="bg-card border-border">
                      <CardContent className="p-3 flex items-center justify-between">
                        <div className="flex-1 space-y-1">
                          <p className="font-medium text-sm">{product.product_name}</p>
                          <div className="flex gap-4 text-xs text-muted-foreground">
                            <span>{"⭐".repeat(product.interest_level)}</span>
                            <span>Qty: {product.quantity}</span>
                            <span>Rp {product.price?.toLocaleString("id-ID") || "0"}</span>
                          </div>
                        </div>
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => removeProductInterest(index)}
                          className="h-8 w-8 p-0 text-destructive hover:text-destructive hover:bg-destructive/10"
                          title={t("productInterest.removeButton")}
                        >
                          <X className="h-4 w-4" />
                        </Button>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </div>
          )}

          <DialogFooter className="pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                reset();
                setProductInterests([]);
                onOpenChange(false);
              }}
              disabled={createActivity.isPending}
            >
              {t("cancel")}
            </Button>
            <Button type="submit" disabled={createActivity.isPending}>
              {createActivity.isPending ? t("creating") : t("create")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
