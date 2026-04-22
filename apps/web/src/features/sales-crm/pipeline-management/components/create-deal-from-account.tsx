"use client";

import { DealForm } from "./deal-form";
import { useCreateDeal } from "../hooks/useDeals";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import type { CreateDealFormData } from "../schemas/deal.schema";

interface CreateDealFromAccountProps {
  readonly onSuccess: () => void;
  readonly onCancel: () => void;
}

export function CreateDealFromAccount({ onSuccess, onCancel }: CreateDealFromAccountProps) {
  const t = useTranslations("pipelineManagement.createDeal");
  const createDeal = useCreateDeal();

  const handleSubmit = async (data: CreateDealFormData) => {
    try {
      await createDeal.mutateAsync({
        ...data,
      });
      toast.success(t("toast.created") || "Opportunity created successfully");
      onSuccess();
    } catch {
      // Error handled by interceptor
    }
  };

  return (
    <DealForm
      // Create mode (no `deal` prop)
      onSubmit={handleSubmit}
      onCancel={onCancel}
      isLoading={createDeal.isPending}
    />
  );
}

