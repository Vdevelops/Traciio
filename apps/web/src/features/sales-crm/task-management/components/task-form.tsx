"use client";

import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  createTaskSchema,
  updateTaskSchema,
  type CreateTaskFormData,
  type UpdateTaskFormData,
  taskPriorityValues,
  taskStatusValues,
  taskTypeValues,
} from "../schemas/task.schema";
import type { Task } from "../types";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { DateTimePicker } from "@/components/ui/date-time-picker";
import { Checkbox } from "@/components/ui/checkbox";
import { useAccounts } from "@/features/sales-crm/account-management/hooks/useAccounts";
import { useContacts } from "@/features/sales-crm/account-management/hooks/useContacts";
import { useLeads } from "@/features/sales-crm/lead-management/hooks/useLeads";
import { useUsers } from "@/features/master-data/user-management/hooks/useUsers";
import { useGoogleCalendarStatus } from "@/features/profile/hooks/useGoogleCalendar";
import { useTranslations } from "next-intl";

interface TaskFormProps {
  readonly task?: Task;
  readonly onSubmit: (data: CreateTaskFormData | UpdateTaskFormData) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

export function TaskForm({ task, onSubmit, onCancel, isLoading }: TaskFormProps) {
  const isEdit = !!task;
  const t = useTranslations("taskManagement.form");

  const { data: accountsData } = useAccounts({ status: "active", per_page: 100 });
  const accounts = accountsData?.data ?? [];

  const { data: usersData } = useUsers({ status: "active", per_page: 100 });
  const users = usersData?.data ?? [];

  const { data: leadsData } = useLeads({ per_page: 100 });
  const leads = leadsData?.data ?? [];

  const { data: googleCalendarStatus } = useGoogleCalendarStatus();
  const isGoogleCalendarConnected = googleCalendarStatus?.data?.connected ?? false;

  // Extract time from ISO string in HH:mm format
  function extractTimeFromISO(isoString: string): string | null {
    try {
      const date = new Date(isoString);
      if (Number.isNaN(date.getTime())) return null;
      const hours = date.getHours().toString().padStart(2, "0");
      const minutes = date.getMinutes().toString().padStart(2, "0");
      return `${hours}:${minutes}`;
    } catch {
      return null;
    }
  }

  // Extract date from ISO string
  function extractDateFromISO(isoString: string): Date | null {
    try {
      const date = new Date(isoString);
      if (Number.isNaN(date.getTime())) return null;
      return date;
    } catch {
      return null;
    }
  }

  // Combine date and time into ISO string
  function combineDateAndTime(date: Date | null, time: string | null): string | null {
    if (!date) return null;
    if (!time) {
      const combined = new Date(date);
      combined.setHours(23, 59, 0, 0);
      return combined.toISOString();
    }
    const [hours, minutes] = time.split(":").map(Number);
    const combined = new Date(date);
    combined.setHours(hours, minutes, 0, 0);
    return combined.toISOString();
  }

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<CreateTaskFormData | UpdateTaskFormData>({
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    resolver: zodResolver(isEdit ? updateTaskSchema : createTaskSchema) as any,
    defaultValues: task
      ? {
          title: task.title,
          description: task.description,
          type: task.type,
          priority: task.priority,
          status: task.status,
          due_date: task.due_date ? extractDateFromISO(task.due_date) : null,
          due_time: task.due_date ? extractTimeFromISO(task.due_date) : null,
          assigned_to: task.assigned_to || "",
          lead_id: task.lead_id || "",
          account_id: task.account_id || "",
          contact_id: task.contact_id || "",
          deal_id: task.deal_id || "",
          sync_to_google_calendar: task.google_calendar_sync_status === "synced" || isGoogleCalendarConnected,
        }
      : {
          type: "general",
          priority: "medium",
          due_date: null,
          due_time: null,
          lead_id: "",
          sync_to_google_calendar: isGoogleCalendarConnected,
        },
  });

  useEffect(() => {
    if (!isEdit && isGoogleCalendarConnected) {
      setValue("sync_to_google_calendar", true);
    }
  }, [isEdit, isGoogleCalendarConnected, setValue]);

  const accountId = watch("account_id") as string | undefined;

  const { data: contactsData } = useContacts({
    account_id: accountId || task?.account_id,
    per_page: 100,
  });
  const contacts = contactsData?.data ?? [];

  useEffect(() => {
    if (!isEdit && accountId && accountId !== (task as { account_id?: string } | undefined)?.account_id) {
      setValue("contact_id", "");
    }
  }, [accountId, isEdit, task, setValue]);

  const handleFormSubmit = async (data: CreateTaskFormData | UpdateTaskFormData | Record<string, unknown>) => {
    const submitData: Record<string, unknown> = {};

    const isValidValue = (value: unknown): boolean => {
      if (value === undefined || value === null) return false;
      if (typeof value === "string") {
        return value.trim() !== "" && value !== "none";
      }
      return true;
    };

    if (isValidValue(data.title)) {
      submitData.title = data.title;
    }

    if (isValidValue(data.description)) {
      submitData.description = data.description;
    }

    if (isValidValue(data.type)) {
      submitData.type = data.type;
    }

    if (!isEdit && isValidValue(data.priority)) {
      submitData.priority = data.priority;
    } else if (isEdit && "priority" in data && isValidValue(data.priority)) {
      submitData.priority = data.priority;
    }

    if (isEdit && "status" in data && isValidValue(data.status)) {
      submitData.status = data.status;
    }

    const dueDate = (data as { due_date?: Date | null; due_time?: string | null }).due_date;
    const dueTime = (data as { due_date?: Date | null; due_time?: string | null }).due_time;

    if ("due_date" in data) {
      if (dueDate) {
        const combined = combineDateAndTime(dueDate, dueTime || null);
        if (combined) {
          submitData.due_date = combined;
        }
      } else {
        submitData.due_date = null;
      }
    }

    if (isValidValue(data.assigned_to)) {
      submitData.assigned_to = data.assigned_to;
    }

    if (isValidValue(data.lead_id)) {
      submitData.lead_id = data.lead_id;
    }

    if (isValidValue(data.account_id)) {
      submitData.account_id = data.account_id;
    }

    if (isValidValue(data.contact_id)) {
      submitData.contact_id = data.contact_id;
    }

    if (isValidValue(data.deal_id)) {
      submitData.deal_id = data.deal_id;
    }

    if (isEdit && Object.keys(submitData).length === 0) {
      return;
    }

    Object.keys(submitData).forEach((key) => {
      if (submitData[key] === undefined || submitData[key] === null) {
        delete submitData[key];
      }
    });

    await onSubmit(submitData as CreateTaskFormData | UpdateTaskFormData);
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-6">
      {/* Title */}
      <Field orientation="vertical">
        <FieldLabel>{t("titleLabel")} *</FieldLabel>
        <Input {...register("title")} placeholder={t("titlePlaceholder")} />
        {errors.title && <FieldError>{errors.title.message}</FieldError>}
      </Field>

      {/* Description */}
      <Field orientation="vertical">
        <FieldLabel>{t("descriptionLabel")}</FieldLabel>
        <Textarea
          {...register("description")}
          placeholder={t("descriptionPlaceholder")}
          rows={3}
        />
        {errors.description && <FieldError>{errors.description.message}</FieldError>}
      </Field>

      {/* Type, Priority, Status */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("typeLabel")}</FieldLabel>
          <Select
            value={(watch("type") as string | undefined) ?? "general"}
            onValueChange={(value) => setValue("type", value as (typeof taskTypeValues)[number])}
          >
            <SelectTrigger className="h-9">
              <SelectValue placeholder={t("typePlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              {taskTypeValues.map((value) => (
                <SelectItem key={value} value={value}>
                  {value.replace("_", " ")}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.type && <FieldError>{errors.type.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("priorityLabel")}</FieldLabel>
          <Select
            value={(watch("priority") as string | undefined) ?? "medium"}
            onValueChange={(value) =>
              setValue("priority", value as (typeof taskPriorityValues)[number])
            }
          >
            <SelectTrigger className="h-9">
              <SelectValue placeholder={t("priorityPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              {taskPriorityValues.map((value) => (
                <SelectItem key={value} value={value}>
                  {value}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.priority && <FieldError>{errors.priority.message}</FieldError>}
        </Field>

        {isEdit && (
          <Field orientation="vertical">
            <FieldLabel>{t("statusLabel")}</FieldLabel>
            <Select
              value={(watch("status") as string | undefined) ?? task?.status ?? "pending"}
              onValueChange={(value) =>
                setValue("status", value as (typeof taskStatusValues)[number])
              }
            >
              <SelectTrigger className="h-9">
                <SelectValue placeholder={t("statusPlaceholder")} />
              </SelectTrigger>
              <SelectContent>
                {taskStatusValues.map((value) => (
                  <SelectItem key={value} value={value}>
                    {value.replace("_", " ")}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {(
              errors as {
                status?: { message?: string };
              }
            ).status && (
              <FieldError>
                {(
                  errors as {
                    status?: { message?: string };
                  }
                ).status?.message}
              </FieldError>
            )}
          </Field>
        )}
      </div>

      {/* Due Date, Assigned To */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("dueDateLabel")}</FieldLabel>
          <DateTimePicker
            date={watch("due_date") as Date | null}
            time={watch("due_time") as string | null}
            onDateChange={(date, time) => {
              setValue("due_date", date, { shouldValidate: true });
              setValue("due_time", time, { shouldValidate: true });
            }}
            disabled={isLoading}
          />
          {errors.due_date && <FieldError>{errors.due_date.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("assignedToLabel")}</FieldLabel>
          <Select
            value={(watch("assigned_to") as string | undefined) ?? ""}
            onValueChange={(value) => setValue("assigned_to", value === "none" ? "" : value)}
          >
            <SelectTrigger className="h-9">
              <SelectValue placeholder={t("assignedToPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="none">{t("assignedToNone")}</SelectItem>
              {users.map((user) => (
                <SelectItem key={user.id} value={user.id}>
                  {user.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.assigned_to && <FieldError>{errors.assigned_to.message}</FieldError>}
        </Field>
      </div>

      <Field orientation="vertical">
        <FieldLabel>{t("leadLabel") || "Lead"}</FieldLabel>
        <Select
          value={(watch("lead_id") as string | undefined) ?? ""}
          onValueChange={(value) => setValue("lead_id", value === "none" ? "" : value)}
        >
          <SelectTrigger className="h-9">
            <SelectValue placeholder={t("leadPlaceholder") || "Select lead (optional)"} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="none">{t("leadNone") || "No lead"}</SelectItem>
            {leads.map((lead) => (
              <SelectItem key={lead.id} value={lead.id}>
                {lead.first_name} {lead.last_name} {lead.company_name ? `(${lead.company_name})` : ""}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {errors.lead_id && <FieldError>{errors.lead_id.message}</FieldError>}
      </Field>

      {/* Account & Contact */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <Field orientation="vertical">
          <FieldLabel>{t("accountLabel")}</FieldLabel>
          <Select
            value={(watch("account_id") as string | undefined) ?? ""}
            onValueChange={(value) => setValue("account_id", value === "none" ? "" : value)}
          >
            <SelectTrigger className="h-9">
              <SelectValue placeholder={t("accountPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="none">{t("accountNone")}</SelectItem>
              {accounts.map((account) => (
                <SelectItem key={account.id} value={account.id}>
                  {account.name
                    ? `${account.name}${account.category ? ` (${account.category.name})` : ""}`
                    : (account.category?.name ?? account.id)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.account_id && <FieldError>{errors.account_id.message}</FieldError>}
        </Field>

        <Field orientation="vertical">
          <FieldLabel>{t("contactLabel")}</FieldLabel>
          <Select
            value={(watch("contact_id") as string | undefined) ?? ""}
            onValueChange={(value) => setValue("contact_id", value === "none" ? "" : value)}
            disabled={!watch("account_id") && !task?.account_id}
          >
            <SelectTrigger className="h-9">
              <SelectValue placeholder={t("contactPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="none">{t("contactNone")}</SelectItem>
              {contacts.map((contact) => (
                <SelectItem key={contact.id} value={contact.id}>
                  {contact.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {errors.contact_id && <FieldError>{errors.contact_id.message}</FieldError>}
        </Field>
      </div>

      {/* Sync to Google Calendar */}
      <Field>
        <div className="flex items-center space-x-2">
          <Checkbox
            id="sync_to_google_calendar"
            checked={watch("sync_to_google_calendar") ?? false}
            onCheckedChange={(checked) => setValue("sync_to_google_calendar", checked === true)}
            disabled={isLoading || !isGoogleCalendarConnected}
          />
          <label
            htmlFor="sync_to_google_calendar"
            className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
          >
            {t("syncToGoogleCalendar")}
          </label>
        </div>
        <p className="text-xs text-muted-foreground mt-1">
          {isGoogleCalendarConnected
            ? t("syncToGoogleCalendarHelp")
            : t("syncToGoogleCalendarNotConnected")}
        </p>
      </Field>

      {/* Actions */}
      <div className="flex justify-end gap-3 pt-4 border-t">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading}>
          {t("cancel")}
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? t("submitting") : isEdit ? t("submitUpdate") : t("submitCreate")}
        </Button>
      </div>
    </form>
  );
}


