"use client";

import type { Lead } from "../types";
import { ConvertLeadDialog } from "./convert-lead-dialog";

interface LeadConvertDialogProps {
  readonly open: boolean;
  readonly onClose: () => void;
  readonly lead: Lead;
  readonly onSuccess?: () => void;
}

export function LeadConvertDialog({ open, onClose, lead, onSuccess }: LeadConvertDialogProps) {
  return <ConvertLeadDialog open={open} onOpenChange={(nextOpen) => !nextOpen && onClose()} lead={lead} onSuccess={onSuccess} />;
}
