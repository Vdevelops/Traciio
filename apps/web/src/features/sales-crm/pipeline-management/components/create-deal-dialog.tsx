"use client";

import { useState } from "react";
import { Building2, UserCheck } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/components/ui/tabs";
import { CreateDealFromAccount } from "./create-deal-from-account";
import { CreateDealFromLead } from "./create-deal-from-lead";
import { useTranslations } from "next-intl";

interface CreateDealDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onSuccess?: () => void;
}

export function CreateDealDialog({ open, onOpenChange, onSuccess }: CreateDealDialogProps) {
  const t = useTranslations("pipelineManagement.createDeal");
  const [activeTab, setActiveTab] = useState<"account" | "lead">("account");

  const handleSuccess = () => {
    onOpenChange(false);
    setActiveTab("account");
    onSuccess?.();
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>

        <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as "account" | "lead")} className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="account" className="gap-2">
              <Building2 className="h-4 w-4" />
              {t("fromAccount")}
            </TabsTrigger>
            <TabsTrigger value="lead" className="gap-2">
              <UserCheck className="h-4 w-4" />
              {t("fromLead")}
            </TabsTrigger>
          </TabsList>

          <TabsContent value="account" className="mt-4">
            <CreateDealFromAccount onSuccess={handleSuccess} onCancel={() => onOpenChange(false)} />
          </TabsContent>

          <TabsContent value="lead" className="mt-4">
            <CreateDealFromLead onSuccess={handleSuccess} onCancel={() => onOpenChange(false)} />
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
}

