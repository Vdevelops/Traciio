"use client";

import { DealForm } from "./deal-form";
import { useCreateDeal } from "../hooks/useDeals";
import { useTranslations } from "next-intl";
import { toast } from "sonner";
import type { CreateDealFormData } from "../schemas/deal.schema";

interface CreateDealFromLeadProps {
  readonly onSuccess: () => void;
  readonly onCancel: () => void;
}

export function CreateDealFromLead({ onSuccess, onCancel }: CreateDealFromLeadProps) {
  const t = useTranslations("pipelineManagement.createDeal");
  const createDeal = useCreateDeal();

  const handleSubmit = async (data: CreateDealFormData) => {
    try {
      await createDeal.mutateAsync({
        ...data,
      });
      toast.success(t("toast.createdFromLead") || "Opportunity created and lead converted successfully");
      onSuccess();
    } catch {
      // Error handled by interceptor
    }
  };

  return (
    <DealForm
      // Create mode (no `deal` prop)
      showQualifiedLeadDropdown={true}
      onSubmit={handleSubmit}
      onCancel={onCancel}
      isLoading={createDeal.isPending}
    />
  );
}

