"use client";

import { useState } from "react";
import { PackagePlus, X } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Field, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { useTranslations } from "next-intl";
import { useProducts } from "../../product-management/hooks/useProducts";

export interface ProductInterestItem {
  product_name: string;
  product_id?: string;
  category_id?: string;
  category_name?: string;
  interest_level: number;
  quantity: number;
  price: number;
}

interface ProductInterestEditorProps {
  readonly value: ProductInterestItem[];
  readonly onChange: (items: ProductInterestItem[]) => void;
  readonly showCommercialFields?: boolean;
  readonly className?: string;
}

export function ProductInterestEditor({
  value,
  onChange,
  showCommercialFields = true,
  className,
}: ProductInterestEditorProps) {
  const t = useTranslations("productInterestEditor");
  const { data: productsData } = useProducts({ per_page: 100, status: "active" });
  const products = productsData?.data ?? [];
  const [currentProduct, setCurrentProduct] = useState<ProductInterestItem>({
    product_name: "",
    interest_level: 3,
    quantity: 1,
    price: 0,
  });

  const addProductInterest = () => {
    if (!currentProduct.product_name.trim()) {
      toast.error(t("validation.productRequired"));
      return;
    }

    if (showCommercialFields && currentProduct.quantity <= 0) {
      toast.error(t("validation.quantityRequired"));
      return;
    }

    onChange([...value, { ...currentProduct }]);
    setCurrentProduct({
      product_name: "",
      interest_level: 3,
      quantity: 1,
      price: 0,
    });
    toast.success(t("toast.added"));
  };

  const removeProductInterest = (index: number) => {
    onChange(value.filter((_, i) => i !== index));
  };

  const handleProductSelect = (productId: string) => {
    const product = products.find((item) => item.id === productId);
    if (!product) return;

    setCurrentProduct((prev) => ({
      ...prev,
      product_id: product.id,
      product_name: product.name,
      price: product.price ? product.price / 100 : prev.price,
    }));
  };

  return (
    <div className={className ?? "crm-stack space-y-3 border-t border-border/70 pt-4"}>
      <div className="flex items-center justify-between">
        <FieldLabel className="mb-0 text-sm font-semibold tracking-tight">{t("title")}</FieldLabel>
        {value.length > 0 && <Badge variant="secondary" className="rounded-full px-3">{t("addedCount", { count: value.length })}</Badge>}
      </div>

      <Card className="crm-panel border-border/70 shadow-none">
        <CardContent className="pt-4 space-y-4">
          <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
            <Field orientation="vertical">
              <FieldLabel className="text-xs">{t("productLabel")} *</FieldLabel>
              {products.length > 0 ? (
                <Select value={currentProduct.product_id ?? ""} onValueChange={handleProductSelect}>
                  <SelectTrigger>
                    <SelectValue placeholder={t("productPlaceholder")} />
                  </SelectTrigger>
                  <SelectContent>
                    {products.map((product) => (
                      <SelectItem key={product.id} value={product.id}>
                        {product.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              ) : (
                <Input
                  placeholder={t("customProductPlaceholder")}
                  value={currentProduct.product_name}
                  onChange={(event) =>
                    setCurrentProduct({
                      ...currentProduct,
                      product_id: undefined,
                      product_name: event.target.value,
                    })
                  }
                />
              )}
            </Field>

            <Field orientation="vertical">
              <FieldLabel className="text-xs">{t("interestLabel")} *</FieldLabel>
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
                  <SelectValue placeholder={t("interestPlaceholder")} />
                </SelectTrigger>
                <SelectContent>
                  {[1, 2, 3, 4, 5].map((level) => (
                    <SelectItem key={level} value={String(level)}>
                      {t("interestValue", { level })}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </Field>
          </div>

          {showCommercialFields && (
            <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
              <Field orientation="vertical">
                <FieldLabel className="text-xs">{t("quantityLabel")} *</FieldLabel>
                <Input
                  type="number"
                  min="1"
                  placeholder={t("quantityPlaceholder")}
                  value={currentProduct.quantity}
                  onChange={(event) =>
                    setCurrentProduct({
                      ...currentProduct,
                      quantity: parseInt(event.target.value, 10) || 1,
                    })
                  }
                />
              </Field>

              <Field orientation="vertical">
                <FieldLabel className="text-xs">{t("priceLabel")}</FieldLabel>
                <Input
                  type="number"
                  min="0"
                  placeholder={t("pricePlaceholder")}
                  value={currentProduct.price}
                  onChange={(event) =>
                    setCurrentProduct({
                      ...currentProduct,
                      price: parseInt(event.target.value, 10) || 0,
                    })
                  }
                />
              </Field>
            </div>
          )}

          <Button type="button" onClick={addProductInterest} size="sm" variant="outline" className="w-full rounded-xl border-dashed">
            <PackagePlus className="h-4 w-4" />
            {t("addButton")}
          </Button>
        </CardContent>
      </Card>

      {value.length > 0 && (
        <div className="space-y-2">
          <p className="text-sm font-medium">{t("addedTitle")}</p>
          {value.map((product, index) => (
            <div
              key={`${product.product_id ?? product.product_name}-${index}`}
              className="crm-list-card flex items-center justify-between rounded-2xl border border-border/70 bg-card/80 p-3"
            >
              <div className="flex-1">
                <p className="font-medium text-sm">{product.product_name}</p>
                <p className="text-xs text-muted-foreground">
                  {showCommercialFields
                    ? t("summary", {
                        interest: product.interest_level,
                        quantity: product.quantity,
                        price: product.price?.toLocaleString("id-ID") ?? "0",
                      })
                    : t("summaryCompact", { interest: product.interest_level })}
                </p>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                onClick={() => removeProductInterest(index)}
                className="h-8 w-8"
                title={t("removeButton")}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
