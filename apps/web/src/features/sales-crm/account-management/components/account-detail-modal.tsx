"use client";

import { Edit, Trash2, Mail, Building2, MapPin, Phone, Calendar, Plus, UserCircle, Pencil } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CustomerPurchaseHistory } from "./CustomerPurchaseHistory";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { useAccount, useDeleteAccount, useUpdateAccount } from "../hooks/useAccounts";
import { useContacts, useCreateContact, useUpdateContact, useDeleteContact } from "../hooks/useContacts";
import { toast } from "sonner";
import { useState } from "react";
import type { CreateDealFormData } from "@/features/sales-crm/pipeline-management/schemas/deal.schema";
import type { CreateContactFormData, UpdateContactFormData } from "../schemas/contact.schema";
import { AccountForm } from "./account-form";
import { ContactForm } from "./contact-form";
import { DealForm } from "@/features/sales-crm/pipeline-management/components/deal-form";
import { useCreateDeal } from "@/features/sales-crm/pipeline-management/hooks/useDeals";
import { useTranslations } from "next-intl";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { toBadgeVariant } from "@/lib/badge-variant";
import { useIsMobile } from "@/hooks/use-mobile";
import { cn, formatPhoneNumberToWA, formatEmailToMailto } from "@/lib/utils";
import type { Contact, Account } from "../types";

interface ContactsCardProps {
  readonly contacts: Contact[];
  readonly isLoading: boolean;
  readonly onCreateContact?: () => void;
  readonly onEditContact?: (contact: Contact) => void;
  readonly onDeleteContact?: (contact: Contact) => void;
}

interface ContactsCardContentProps {
  readonly contacts: Contact[];
  readonly onCreateContact?: () => void;
  readonly onEditContact?: (contact: Contact) => void;
  readonly onDeleteContact?: (contact: Contact) => void;
}

function ContactsCardContent({ contacts, onCreateContact, onEditContact, onDeleteContact }: ContactsCardContentProps) {
  if (contacts.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-8 text-muted-foreground">
        <UserCircle className="h-10 w-10 mb-2 opacity-30" />
        <p className="text-sm">No contacts found for this account</p>
        {onCreateContact && (
          <Button variant="outline" size="sm" className="mt-4" onClick={onCreateContact}>
            <Plus className="h-4 w-4 mr-1.5" />
            Add First Contact
          </Button>
        )}
      </div>
    );
  }
  return (
    <div className="space-y-2">
      {contacts.map((contact) => (
        <div
          key={contact.id}
          className="flex items-start gap-3 p-3 rounded-lg border border-border/50 hover:bg-muted/30 transition-colors"
        >
          <div className="h-9 w-9 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
            <UserCircle className="h-5 w-5 text-primary" />
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <span className="text-sm font-medium truncate">{contact.name}</span>
              {contact.role && (
                <Badge
                  variant={toBadgeVariant(contact.role.badge_color, "secondary")}
                  className="text-[10px] h-4 px-1.5 font-normal shrink-0"
                >
                  {contact.role.name}
                </Badge>
              )}
            </div>
            {contact.position && (
              <p className="text-xs text-muted-foreground mt-0.5">{contact.position}</p>
            )}
            <div className="flex items-center gap-3 mt-1 flex-wrap">
              {contact.phone && (
                <a
                  href={formatPhoneNumberToWA(contact.phone)}
                  target="_blank"
                  rel="noreferrer"
                  className="flex items-center gap-1 text-xs text-primary hover:underline"
                >
                  <Phone className="h-3 w-3" />
                  {contact.phone}
                </a>
              )}
              {contact.email && (
                <a
                  href={formatEmailToMailto(contact.email)}
                  className="flex items-center gap-1 text-xs text-primary hover:underline"
                >
                  <Mail className="h-3 w-3" />
                  {contact.email}
                </a>
              )}
            </div>
          </div>
          {(onEditContact || onDeleteContact) && (
            <div className="flex items-center gap-1 shrink-0 ml-auto">
              {onEditContact && (
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7 text-muted-foreground hover:text-foreground"
                  onClick={() => onEditContact(contact)}
                  title="Edit contact"
                >
                  <Pencil className="h-3.5 w-3.5" />
                </Button>
              )}
              {onDeleteContact && (
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7 text-muted-foreground hover:text-destructive"
                  onClick={() => onDeleteContact(contact)}
                  title="Delete contact"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </Button>
              )}
            </div>
          )}
        </div>
      ))}
    </div>
  );
}

function ContactsCard({ contacts, isLoading, onCreateContact, onEditContact, onDeleteContact }: ContactsCardProps) {
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <UserCircle className="h-5 w-5 text-muted-foreground" />
              Contacts
              {!isLoading && (
                <Badge variant="secondary" className="text-xs font-normal">
                  {contacts.length}
                </Badge>
              )}
            </CardTitle>
            <CardDescription className="mt-1">People associated with this account</CardDescription>
          </div>
          {onCreateContact && (
            <Button variant="outline" size="sm" onClick={onCreateContact} className="shrink-0">
              <Plus className="h-4 w-4 mr-1.5" />
              Add Contact
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="space-y-3">
            {Array.from({ length: 3 }, (_, i) => (
              <div key={i} className="flex items-center gap-3 p-3 rounded-lg border">
                <Skeleton className="h-9 w-9 rounded-full shrink-0" />
                <div className="flex-1 space-y-1.5">
                  <Skeleton className="h-4 w-32" />
                  <Skeleton className="h-3 w-24" />
                </div>
              </div>
            ))}
          </div>
        ) : (
          <ContactsCardContent
            contacts={contacts}
            onCreateContact={onCreateContact}
            onEditContact={onEditContact}
            onDeleteContact={onDeleteContact}
          />
        )}
      </CardContent>
    </Card>
  );
}

interface AccountInfoGridProps {
  readonly account: Account;
}

function AccountInfoGrid({ account }: AccountInfoGridProps) {
  const t = useTranslations("accountManagement.accountDetail");
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div className="space-y-2">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <Building2 className="h-4 w-4" />
          <span>{t("infoCard.name")}</span>
        </div>
        <div className="text-base font-medium">{account.name || t("infoCard.notAvailable")}</div>
      </div>

      <div className="space-y-2">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <span>{t("infoCard.category")}</span>
        </div>
        <div>
          {account.category ? (
            <Badge variant={toBadgeVariant(account.category?.badge_color, "secondary")} className="font-normal">
              {account.category?.name || "-"}
            </Badge>
          ) : (
            <span className="text-muted-foreground">{t("infoCard.notAvailable")}</span>
          )}
        </div>
      </div>

      {account.address && (
        <div className="space-y-2 md:col-span-2">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <MapPin className="h-4 w-4" />
            <span>{t("infoCard.address")}</span>
          </div>
          <div className="text-base font-medium">{account.address}</div>
        </div>
      )}

      {account.city && (
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <MapPin className="h-4 w-4" />
            <span>{t("infoCard.city")}</span>
          </div>
          <div className="text-base font-medium">{account.city}</div>
        </div>
      )}

      {account.province && (
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <MapPin className="h-4 w-4" />
            <span>{t("infoCard.province")}</span>
          </div>
          <div className="text-base font-medium">{account.province}</div>
        </div>
      )}

      {account.phone && (
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Phone className="h-4 w-4" />
            <span>{t("infoCard.phone")}</span>
          </div>
          <a
            href={formatPhoneNumberToWA(account.phone)}
            target="_blank"
            rel="noreferrer"
            className="text-sm font-medium text-primary hover:underline cursor-pointer"
          >
            {account.phone}
          </a>
        </div>
      )}

      {account.email && (
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Mail className="h-4 w-4" />
            <span>{t("infoCard.email")}</span>
          </div>
          <a
            href={formatEmailToMailto(account.email)}
            className="text-sm font-medium text-primary hover:underline cursor-pointer"
          >
            {account.email}
          </a>
        </div>
      )}

      {typeof account.latitude === "number" && (
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <MapPin className="h-4 w-4" />
            <span>{t("infoCard.latitude")}</span>
          </div>
          <div className="text-base font-medium">{account.latitude.toFixed(6)}</div>
        </div>
      )}

      {typeof account.longitude === "number" && (
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <MapPin className="h-4 w-4" />
            <span>{t("infoCard.longitude")}</span>
          </div>
          <div className="text-base font-medium">{account.longitude.toFixed(6)}</div>
        </div>
      )}

      <div className="space-y-2">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <span>{t("infoCard.status")}</span>
        </div>
        <div>
          <Badge variant={account.status === "active" ? "active" : "inactive"}>
            {account.status}
          </Badge>
        </div>
      </div>

      <div className="space-y-2">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <Calendar className="h-4 w-4" />
          <span>{t("infoCard.createdAt")}</span>
        </div>
        <div className="text-base font-medium">
          {account.created_at
            ? new Date(account.created_at).toLocaleDateString()
            : t("infoCard.notAvailable")}
        </div>
      </div>

      <div className="space-y-2">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <Calendar className="h-4 w-4" />
          <span>{t("infoCard.updatedAt")}</span>
        </div>
        <div className="text-base font-medium">
          {account.updated_at
            ? new Date(account.updated_at).toLocaleDateString()
            : t("infoCard.notAvailable")}
        </div>
      </div>
    </div>
  );
}

interface AccountDetailModalProps {
  readonly accountId: string | null;
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onAccountUpdated?: () => void;
}

export function AccountDetailModal({ accountId, open, onOpenChange, onAccountUpdated }: AccountDetailModalProps) {
  const { data, isLoading, error } = useAccount(accountId || "");
  const deleteAccount = useDeleteAccount();
  const updateAccount = useUpdateAccount();
  const createContact = useCreateContact();
  const updateContact = useUpdateContact();
  const deleteContact = useDeleteContact();
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isCreateOpportunityDialogOpen, setIsCreateOpportunityDialogOpen] = useState(false);
  const [isCreateContactDialogOpen, setIsCreateContactDialogOpen] = useState(false);
  const [editingContact, setEditingContact] = useState<Contact | null>(null);
  const [deletingContact, setDeletingContact] = useState<Contact | null>(null);
  const t = useTranslations("accountManagement.accountDetail");
  const hasCreateOpportunityPermission = useHasPermission("pipeline.opportunity-create");
  const hasCreateContactPermission = useHasPermission("accounts.edit");
  const hasEditContactPermission = useHasPermission("accounts.edit");
  const hasDeleteContactPermission = useHasPermission("accounts.delete");
  const createDeal = useCreateDeal();
  const isMobile = useIsMobile();

  const account = data?.data;

  const { data: contactsData, isLoading: isContactsLoading } = useContacts(
    { account_id: accountId || "", per_page: 50 },
    { enabled: !!accountId && open }
  );
  const contacts = contactsData?.data || [];

  const handleDeleteConfirm = async () => {
    if (!account || !accountId) return;
    try {
      await deleteAccount.mutateAsync(accountId);
      toast.success(t("toastDeleted"));
      onOpenChange(false);
      onAccountUpdated?.();
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  const handleCreateOpportunitySubmit = async (data: CreateDealFormData) => {
    try {
      await createDeal.mutateAsync(data);
      toast.success(t("toast.opportunityCreated") || "Opportunity created successfully");
      setIsCreateOpportunityDialogOpen(false);
      onAccountUpdated?.();
    } catch {
      // Error handled by interceptor
    }
  };

  const handleCreateContactSubmit = async (data: CreateContactFormData | UpdateContactFormData) => {
    try {
      await createContact.mutateAsync(data as CreateContactFormData);
      toast.success("Contact created successfully");
      setIsCreateContactDialogOpen(false);
      onAccountUpdated?.();
    } catch {
      // Error handled by interceptor
    }
  };

  const handleEditContactSubmit = async (data: CreateContactFormData | UpdateContactFormData) => {
    if (!editingContact) return;
    try {
      await updateContact.mutateAsync({ id: editingContact.id, data: data as UpdateContactFormData });
      toast.success("Contact updated successfully");
      setEditingContact(null);
      onAccountUpdated?.();
    } catch {
      // Error handled by interceptor
    }
  };

  const handleDeleteContactConfirm = async () => {
    if (!deletingContact) return;
    try {
      await deleteContact.mutateAsync(deletingContact.id);
      toast.success("Contact deleted successfully");
      setDeletingContact(null);
      onAccountUpdated?.();
    } catch {
      // Error handled by interceptor
    }
  };


  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className={cn("sm:max-w-[700px] max-h-[90vh] overflow-y-auto", isMobile && "mx-2 max-w-[calc(100vw-1rem)]")}>
          <DialogHeader className="flex flex-row items-center justify-between space-y-0 pr-8">
            <DialogTitle>{t("title")}</DialogTitle>
            {account && (
              <div className="flex items-center gap-1 shrink-0">
                {hasCreateOpportunityPermission && (
                  <Button
                    variant="outline"
                    size="icon"
                    className="h-8 w-8"
                    onClick={() => setIsCreateOpportunityDialogOpen(true)}
                    title={t("header.createOpportunity") || "Create Opportunity"}
                  >
                    <Plus className="h-4 w-4" />
                  </Button>
                )}
                <Button
                  variant="outline"
                  size="icon"
                  className="h-8 w-8"
                  onClick={() => setIsEditDialogOpen(true)}
                  title="Edit"
                >
                  <Edit className="h-4 w-4" />
                </Button>
                <Button
                  variant="destructive"
                  size="icon"
                  className="h-8 w-8"
                  onClick={() => setIsDeleteDialogOpen(true)}
                  disabled={deleteAccount.isPending}
                  title="Delete"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            )}
          </DialogHeader>

          {isLoading && (
            <div className="space-y-6">
              <div className="flex items-center gap-4">
                <Skeleton className="h-16 w-16 rounded-lg" />
                <div className="space-y-2">
                  <Skeleton className="h-6 w-48" />
                  <Skeleton className="h-4 w-32" />
                </div>
              </div>
              <Card>
                <CardHeader>
                  <Skeleton className="h-6 w-32" />
                  <Skeleton className="h-4 w-64 mt-2" />
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
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
              {t("loadError")}
            </div>
          )}

          {!isLoading && !error && account && (
            <div className="space-y-6">
              {/* Account Header */}
              <div className={cn("pb-4 border-b", isMobile ? "space-y-2" : "flex items-center gap-4")}>
                <div className={cn("flex items-center gap-4", isMobile && "flex-col items-start")}>
                  <div className={cn("rounded-lg bg-primary/10 flex items-center justify-center shrink-0", isMobile ? "h-12 w-12" : "h-16 w-16")}>
                    <Building2 className={cn("text-primary", isMobile ? "h-6 w-6" : "h-8 w-8")} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h2 className={cn("font-medium tracking-tight wrap-break-word", isMobile ? "text-xl" : "text-2xl")}>{account.name || t("infoCard.notAvailable")}</h2>
                    <div className="flex flex-wrap items-center gap-2 mt-1">
                      {account.category ? (
                        <Badge variant={toBadgeVariant(account.category?.badge_color, "secondary")} className="font-normal text-xs">
                          {account.category?.name || "-"}
                        </Badge>
                      ) : (
                        <Badge variant="outline" className="text-xs">-</Badge>
                      )}
                      <Badge variant={account.status === "active" ? "active" : "inactive"} className="text-xs">
                        {account.status || "-"}
                      </Badge>
                    </div>
                  </div>
                </div>
              </div>

              <Tabs defaultValue="overview" className="w-full">
                <TabsList className="grid w-full grid-cols-3 mb-4">
                  <TabsTrigger value="overview">Overview</TabsTrigger>
                  <TabsTrigger value="contacts">Contacts</TabsTrigger>
                  <TabsTrigger value="purchases">Purchase History</TabsTrigger>
                </TabsList>

                <TabsContent value="overview" className="space-y-6 mt-0">
                  {/* Account Info Card */}
                  <Card>
                    <CardHeader>
                      <CardTitle>{t("infoCard.title")}</CardTitle>
                      <CardDescription>{t("infoCard.description")}</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-6">
                      <AccountInfoGrid account={account} />
                    </CardContent>
                  </Card>
                </TabsContent>

                <TabsContent value="contacts" className="space-y-6 mt-0">
                  {/* Contacts Card */}
                  <ContactsCard
                    contacts={contacts}
                    isLoading={isContactsLoading}
                    onCreateContact={hasCreateContactPermission ? () => setIsCreateContactDialogOpen(true) : undefined}
                    onEditContact={hasEditContactPermission ? (contact) => setEditingContact(contact) : undefined}
                    onDeleteContact={hasDeleteContactPermission ? (contact) => setDeletingContact(contact) : undefined}
                  />
                </TabsContent>

                <TabsContent value="purchases" className="space-y-6 mt-0">
                  <CustomerPurchaseHistory accountId={accountId!} />
                </TabsContent>
              </Tabs>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      {isEditDialogOpen && account && (
        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogContent className={cn("sm:max-w-[600px] max-h-[90vh] overflow-y-auto", isMobile && "mx-2 max-w-[calc(100vw-1rem)]")}>
            <DialogHeader>
              <DialogTitle>{t("editDialogTitle")}</DialogTitle>
            </DialogHeader>
            <AccountForm
              account={account}
              onSubmit={async (formData) => {
                try {
                  await updateAccount.mutateAsync({ id: accountId!, data: formData });
                  setIsEditDialogOpen(false);
                  toast.success(t("toastUpdated"));
                  onAccountUpdated?.();
                } catch {
                  // Error already handled in api-client interceptor
                }
              }}
              onCancel={() => setIsEditDialogOpen(false)}
              isLoading={updateAccount.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Create Opportunity Dialog */}
      {account && (
        <Dialog open={isCreateOpportunityDialogOpen} onOpenChange={setIsCreateOpportunityDialogOpen}>
          <DialogContent className={cn("sm:max-w-[600px] max-h-[90vh] overflow-y-auto", isMobile && "mx-2 max-w-[calc(100vw-1rem)]")}>
            <DialogHeader>
              <DialogTitle>{t("header.createOpportunity") || "Create Opportunity for Account"}</DialogTitle>
            </DialogHeader>
            <DealForm
              initialAccountId={account.id}
              onSubmit={handleCreateOpportunitySubmit}
              onCancel={() => setIsCreateOpportunityDialogOpen(false)}
              isLoading={createDeal.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Create Contact Dialog */}
      {account && (
        <Dialog open={isCreateContactDialogOpen} onOpenChange={setIsCreateContactDialogOpen}>
          <DialogContent className={cn("sm:max-w-[600px] max-h-[90vh] overflow-y-auto", isMobile && "mx-2 max-w-[calc(100vw-1rem)]")}>
            <DialogHeader>
              <DialogTitle>Add Contact to Account</DialogTitle>
            </DialogHeader>
            <ContactForm
              defaultAccountId={account.id}
              onSubmit={handleCreateContactSubmit}
              onCancel={() => setIsCreateContactDialogOpen(false)}
              isLoading={createContact.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Delete Account Dialog */}
      <DeleteDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        onConfirm={handleDeleteConfirm}
        title={t("deleteDialogTitle")}
        description={
          account
            ? t("deleteDialogDescriptionWithName", { name: account.name })
            : t("deleteDialogDescription")
        }
        itemName={t("deleteDialogItemName")}
        isLoading={deleteAccount.isPending}
      />

      {/* Edit Contact Dialog */}
      {editingContact && (
        <Dialog open={!!editingContact} onOpenChange={(open) => { if (!open) setEditingContact(null); }}>
          <DialogContent className={cn("sm:max-w-[600px] max-h-[90vh] overflow-y-auto", isMobile && "mx-2 max-w-[calc(100vw-1rem)]")}
          >
            <DialogHeader>
              <DialogTitle>Edit Contact</DialogTitle>
            </DialogHeader>
            <ContactForm
              contact={editingContact}
              onSubmit={handleEditContactSubmit}
              onCancel={() => setEditingContact(null)}
              isLoading={updateContact.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Delete Contact Dialog */}
      <DeleteDialog
        open={!!deletingContact}
        onOpenChange={(open) => { if (!open) setDeletingContact(null); }}
        onConfirm={handleDeleteContactConfirm}
        title="Delete Contact"
        description={
          deletingContact
            ? `Are you sure you want to delete "${deletingContact.name}"? This action cannot be undone.`
            : "Are you sure you want to delete this contact?"
        }
        itemName="contact"
        isLoading={deleteContact.isPending}
      />
    </>
  );
}

