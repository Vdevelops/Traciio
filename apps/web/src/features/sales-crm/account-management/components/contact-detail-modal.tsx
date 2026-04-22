"use client";

import { Calendar, Mail, Phone, User, Building2, StickyNote } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useContact } from "../hooks/useContacts";
import type { Contact } from "../types";
import { useTranslations } from "next-intl";
import { useIsMobile } from "@/hooks/use-mobile";
import { cn, formatPhoneNumberToWA, formatEmailToMailto } from "@/lib/utils";

interface ContactDetailModalProps {
  readonly contactId: string | null;
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onContactUpdated?: () => void;
}

export function ContactDetailModal({
  contactId,
  open,
  onOpenChange,
}: ContactDetailModalProps) {
  const { data, isLoading, error } = useContact(contactId || "");
  const t = useTranslations("accountManagement.contactList");
  const isMobile = useIsMobile();

  const contact: Contact | undefined = data?.data;

  const formatDate = (dateString?: string) => {
    if (!dateString) return "-";
    const date = new Date(dateString);
    return date.toLocaleDateString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className={cn("sm:max-w-[600px] max-h-[90vh] overflow-y-auto", isMobile && "mx-2 max-w-[calc(100vw-1rem)]")}>
        <DialogHeader>
          <DialogTitle>{t("detailTitle")}</DialogTitle>
        </DialogHeader>

        {isLoading && (
          <div className="space-y-6">
            <div className="flex items-center gap-4">
              <Skeleton className="h-14 w-14 rounded-full" />
              <div className="space-y-2">
                <Skeleton className="h-5 w-40" />
                <Skeleton className="h-4 w-32" />
              </div>
            </div>
            <Card>
              <CardHeader>
                <Skeleton className="h-5 w-32" />
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <Skeleton className="h-4 w-full" />
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-4 w-1/2" />
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {error && (
          <div className="text-center text-muted-foreground py-8">
            {t("empty")}
          </div>
        )}

        {!isLoading && !error && contact && (
          <div className="space-y-6">
            {/* Header */}
            <div className={cn("pb-4 border-b", isMobile ? "space-y-3" : "flex items-center gap-4")}>
              <div className={cn("rounded-full bg-primary/10 flex items-center justify-center shrink-0", isMobile ? "h-10 w-10" : "h-12 w-12")}>
                <User className={cn("text-primary", isMobile ? "h-5 w-5" : "h-6 w-6")} />
              </div>
              <div className="flex-1 min-w-0">
                <h2 className={cn("font-medium tracking-tight break-words", isMobile ? "text-lg" : "text-xl")}>{contact.name}</h2>
                <div className="flex flex-wrap items-center gap-2 mt-1 text-xs text-muted-foreground">
                  {contact.position && <span className="break-words">{contact.position}</span>}
                  {contact.role && (
                    <Badge variant={contact.role.badge_color} className="font-normal capitalize text-xs">
                      {contact.role.name}
                    </Badge>
                  )}
                </div>
              </div>
            </div>

            {/* Info card */}
            <Card>
              <CardHeader>
                <CardTitle>{t("table.name")}</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4 text-sm">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {contact.phone && (
                    <div className="space-y-1">
                      <div className="flex items-center gap-2 text-muted-foreground">
                        <Phone className="h-4 w-4" />
                        <span>{t("table.phone")}</span>
                      </div>
                      <a 
                        href={formatPhoneNumberToWA(contact.phone)} 
                        target="_blank" 
                        rel="noreferrer" 
                        className="font-medium text-sm text-primary hover:underline cursor-pointer"
                      >
                        {contact.phone}
                      </a>
                    </div>
                  )}

                  {contact.email && (
                    <div className="space-y-1">
                      <div className="flex items-center gap-2 text-muted-foreground">
                        <Mail className="h-4 w-4" />
                        <span>{t("table.email")}</span>
                      </div>
                      <a 
                        href={formatEmailToMailto(contact.email)} 
                        className="font-medium text-sm text-primary hover:underline cursor-pointer"
                      >
                        {contact.email}
                      </a>
                    </div>
                  )}

                  {contact.account_id && (
                    <div className="space-y-1">
                      <div className="flex items-center gap-2 text-muted-foreground">
                        <Building2 className="h-4 w-4" />
                        <span>{t("allAccounts")}</span>
                      </div>
                      <div className="font-medium">{contact.account_id}</div>
                    </div>
                  )}

                  <div className="space-y-1">
                    <div className="flex items-center gap-2 text-muted-foreground">
                      <Calendar className="h-4 w-4" />
                      <span>Created / Updated</span>
                    </div>
                    <div className="flex flex-col text-xs text-muted-foreground">
                      <span>Created: {formatDate(contact.created_at)}</span>
                      <span>Updated: {formatDate(contact.updated_at)}</span>
                    </div>
                  </div>
                </div>

                {contact.notes && (
                  <div className="space-y-1 pt-2 border-t">
                    <div className="flex items-center gap-2 text-muted-foreground">
                      <StickyNote className="h-4 w-4" />
                      <span>Notes</span>
                    </div>
                    <p className="text-sm whitespace-pre-wrap">{contact.notes}</p>
                  </div>
                )}
              </CardContent>
            </Card>

            <div className="flex justify-end">
              <Button type="button" variant="outline" size="sm" onClick={() => onOpenChange(false)}>
                Close
              </Button>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}


