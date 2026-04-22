"use client";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { AlertTriangle, X } from "lucide-react";
import { useTranslations } from "next-intl";

interface FakeGPSWarningModalProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly reason?: string;
}

export function FakeGPSWarningModal({
  open,
  onOpenChange,
  reason,
}: FakeGPSWarningModalProps) {
  const t = useTranslations("fakeGPSWarning");

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-destructive/10">
              <AlertTriangle className="h-5 w-5 text-destructive" />
            </div>
            <div className="flex-1">
              <DialogTitle>{t("title")}</DialogTitle>
              <DialogDescription className="mt-1">
                {t("description")}
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <div className="space-y-4">
          {reason && (
            <div className="rounded-md bg-muted p-3">
              <p className="text-sm text-muted-foreground">
                <strong>{t("detectedReason")}:</strong> {reason}
              </p>
            </div>
          )}

          <div className="space-y-3">
            <div>
              <h4 className="text-sm font-medium mb-2">
                {t("instructions.title")}
              </h4>
              <ol className="list-decimal list-inside space-y-2 text-sm text-muted-foreground">
                <li>{t("instructions.step1")}</li>
                <li>{t("instructions.step2")}</li>
                <li>{t("instructions.step3")}</li>
                <li>{t("instructions.step4")}</li>
              </ol>
            </div>

            <div className="rounded-md bg-yellow-500/10 border border-yellow-500/20 p-3">
              <p className="text-sm text-yellow-700 dark:text-yellow-400">
                <strong>{t("important")}:</strong> {t("importantNote")}
              </p>
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button
              variant="outline"
              onClick={() => onOpenChange(false)}
              className="cursor-pointer"
            >
              {t("buttons.close")}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}

