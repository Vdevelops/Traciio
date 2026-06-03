"use client";

import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useTranslations } from "next-intl";
import { LeadForm } from "./lead-form";
import { LeadQualificationCard } from "./LeadQualificationCard";
import { useCreateLead, useUpdateLead } from "../hooks/useLeads";
import type { Lead } from "../types";
import type { CreateLeadFormData, UpdateLeadFormData } from "../schemas/lead.schema";

interface LeadFormDialogProps {
  readonly open: boolean;
  readonly onClose: () => void;
  readonly lead?: Lead | null;
  readonly onSuccess?: () => void;
}

export function LeadFormDialog({ open, onClose, lead, onSuccess }: LeadFormDialogProps) {
  const createLead = useCreateLead();
  const updateLead = useUpdateLead();
  const isEdit = !!lead;
  const tLeadDetail = useTranslations("leadManagement.leadDetail");
  const tLeadList = useTranslations("leadManagement.leadList");

  const handleSubmit = async (data: CreateLeadFormData | UpdateLeadFormData) => {
    if (lead) {
      await updateLead.mutateAsync({ id: lead.id, data: data as UpdateLeadFormData });
    } else {
      await createLead.mutateAsync(data as CreateLeadFormData);
    }

    onSuccess?.();
    onClose();
  };

  return (
    <Dialog open={open} onOpenChange={(nextOpen) => !nextOpen && onClose()}>
      <DialogContent className="max-h-[90vh] overflow-y-auto border-border/70 bg-card/95 sm:max-w-3xl">
        <DialogHeader>
          <DialogTitle>
            {isEdit ? tLeadDetail("editDialog.title") : tLeadList("dialogs.createTitle")}
          </DialogTitle>
        </DialogHeader>

        {isEdit && lead ? (
          <Tabs defaultValue="profile" className="w-full">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="profile">{tLeadDetail("editDialog.tabs.profile")}</TabsTrigger>
              <TabsTrigger value="bant">{tLeadDetail("editDialog.tabs.bant")}</TabsTrigger>
            </TabsList>
            <TabsContent value="profile" className="mt-4">
              <LeadForm
                lead={lead}
                onSubmit={handleSubmit}
                onCancel={onClose}
                isLoading={createLead.isPending || updateLead.isPending}
              />
            </TabsContent>
            <TabsContent value="bant" className="mt-4">
              <LeadQualificationCard leadId={lead.id} />
            </TabsContent>
          </Tabs>
        ) : (
          <LeadForm
            lead={lead ?? undefined}
            onSubmit={handleSubmit}
            onCancel={onClose}
            isLoading={createLead.isPending || updateLead.isPending}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}
