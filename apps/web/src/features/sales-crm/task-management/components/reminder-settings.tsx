"use client";

import React, { useState } from "react";
import { Clock, Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { ReminderDateTimePicker } from "@/components/ui/reminder-date-time-picker";
import type { Reminder } from "../types";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  createReminderSchema,
  updateReminderSchema,
  type CreateReminderFormData,
  type UpdateReminderFormData,
} from "../schemas/reminder.schema";
import { useTaskReminders } from "../hooks/useTaskList";
import { useTranslations } from "next-intl";

interface ReminderSettingsProps {
  readonly taskId: string;
}

// Separate component for header button to be used in CardTitle
export function ReminderSettingsHeader({ onCreateClick }: { onCreateClick?: () => void }) {
  const t = useTranslations("taskManagement.reminders");
  
  return (
    <Button
      type="button"
      size="sm"
      variant="outline"
      className="h-8 gap-2"
      onClick={onCreateClick}
    >
      <Plus className="h-3.5 w-3.5" />
      <span className="text-xs">{t("addButton")}</span>
    </Button>
  );
}

export const ReminderSettings = React.forwardRef<
  { openDialog: () => void },
  ReminderSettingsProps
>(({ taskId }, ref) => {
  const {
    reminders,
    isLoading,
    handleCreateReminder,
    handleUpdateReminder,
    handleDeleteReminder,
  } = useTaskReminders(taskId);

  const t = useTranslations("taskManagement.reminders");
  const tDialog = useTranslations("taskManagement.reminderDialog");

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingReminder, setEditingReminder] = useState<Reminder | null>(null);
  const [deletingReminderId, setDeletingReminderId] = useState<string | null>(null);
  
  // Expose open dialog function to parent via ref
  React.useImperativeHandle(ref, () => ({
    openDialog: () => setIsCreateDialogOpen(true),
  }), []);

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <div className="space-y-4">

      {isLoading ? (
        <p className="text-sm text-muted-foreground">{t("loading")}</p>
      ) : reminders.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-8 text-center">
          <div className="rounded-full bg-muted p-3 mb-3">
            <Clock className="h-5 w-5 text-muted-foreground" />
          </div>
          <p className="text-sm text-muted-foreground">{t("empty")}</p>
        </div>
      ) : (
        <div className="space-y-3">
          {reminders.map((reminder) => (
            <div
              key={reminder.id}
              className="group flex items-start justify-between rounded-lg border bg-card p-4 hover:border-primary/20 hover:shadow-sm transition-all"
            >
              <div className="space-y-1">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">
                    {formatDateTime(reminder.remind_at)}
                  </span>
                  <span className="text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded">
                    {reminder.reminder_type}
                  </span>
                </div>
                {reminder.message && (
                  <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                    {reminder.message}
                  </p>
                )}
              </div>
              <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                <Button
                  type="button"
                  size="icon-sm"
                  variant="ghost"
                  className="h-8 w-8"
                  onClick={() => setEditingReminder(reminder)}
                  title={tDialog("editTitle")}
                >
                  <Clock className="h-3.5 w-3.5" />
                </Button>
                <Button
                  type="button"
                  size="icon-sm"
                  variant="ghost"
                  className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
                  onClick={() => setDeletingReminderId(reminder.id)}
                  title={t("deleteTitle")}
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create Reminder Dialog */}
      <ReminderDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        taskId={taskId}
        mode="create"
        onCreate={handleCreateReminder}
      />

      {/* Edit Reminder Dialog */}
      {editingReminder && (
        <ReminderDialog
          open={!!editingReminder}
          onOpenChange={(open) => {
            if (!open) {
              setEditingReminder(null);
            }
          }}
          taskId={taskId}
          mode="edit"
          reminder={editingReminder}
          onUpdate={async (data) => {
            await handleUpdateReminder(editingReminder.id, data);
            setEditingReminder(null);
          }}
        />
      )}

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingReminderId}
        onOpenChange={(open) => {
          if (!open) {
            setDeletingReminderId(null);
          }
        }}
        onConfirm={async () => {
          if (!deletingReminderId) return;
          await handleDeleteReminder(deletingReminderId);
          setDeletingReminderId(null);
        }}
        title={t("deleteTitle")}
        description={t("deleteDescription")}
        itemName={t("deleteItemName")}
        isLoading={false}
      />
    </div>
  );
});

ReminderSettings.displayName = "ReminderSettings";

interface ReminderDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly taskId: string;
  readonly mode: "create" | "edit";
  readonly reminder?: Reminder;
  readonly onCreate?: (data: CreateReminderFormData) => Promise<void>;
  readonly onUpdate?: (data: UpdateReminderFormData) => Promise<void>;
}

function ReminderDialog({
  open,
  onOpenChange,
  taskId,
  mode,
  reminder,
  onCreate,
  onUpdate,
}: ReminderDialogProps) {
  const isEdit = mode === "edit" && !!reminder;
  const tDialog = useTranslations("taskManagement.reminderDialog");

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<CreateReminderFormData | UpdateReminderFormData>({
    resolver: zodResolver(isEdit ? updateReminderSchema : createReminderSchema),
    defaultValues: isEdit
      ? {
          task_id: reminder?.task_id ?? taskId,
          reminder_type: reminder?.reminder_type,
          message: reminder?.message,
          remind_at: reminder?.remind_at
            ? new Date(reminder.remind_at).toISOString()
            : "",
        }
      : {
          task_id: taskId,
          reminder_type: "in_app",
        },
  });

  // Watch remind_at value for the picker
  const remindAtValue = watch("remind_at");

  const onSubmit = async (data: CreateReminderFormData | UpdateReminderFormData) => {
    if (isEdit) {
      if (!onUpdate) return;
      const submitData: UpdateReminderFormData = {};

      if ("remind_at" in data && data.remind_at) {
        const date = new Date(data.remind_at);
        if (!Number.isNaN(date.getTime())) {
          submitData.remind_at = date.toISOString();
        }
      }

      if ("reminder_type" in data && data.reminder_type) {
        submitData.reminder_type = data.reminder_type;
      }

      if ("message" in data && data.message !== undefined) {
        submitData.message = data.message;
      }

      await onUpdate(submitData);
    } else {
      if (!onCreate) return;
      const submitData: CreateReminderFormData = {
        task_id: taskId,
        reminder_type: "in_app",
        remind_at: "",
      };

      if ("remind_at" in data && data.remind_at) {
        const date = new Date(data.remind_at);
        if (!Number.isNaN(date.getTime())) {
          submitData.remind_at = date.toISOString();
        }
      }

      if ("reminder_type" in data && data.reminder_type) {
        submitData.reminder_type = data.reminder_type;
      }

      if ("message" in data && data.message) {
        submitData.message = data.message;
      }

      await onCreate(submitData);
    }

    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent overlayClassName="z-[1040]" className="sm:max-w-md z-[1041]">
        <DialogHeader>
          <DialogTitle>
            {isEdit ? tDialog("editTitle") : tDialog("createTitle")}
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {!isEdit && (
            <input type="hidden" {...register("task_id")} value={taskId} />
          )}

          <Field orientation="vertical">
            <FieldLabel>{tDialog("remindAtLabel")} *</FieldLabel>
            <ReminderDateTimePicker
              value={remindAtValue ?? ""}
              onChange={(value) => {
                setValue("remind_at", value);
              }}
              error={errors.remind_at?.message}
            />
          </Field>

          <Field orientation="vertical">
            <FieldLabel>{tDialog("channelLabel")}</FieldLabel>
            <Select
              value={watch("reminder_type") ?? "in_app"}
              onValueChange={(value) => {
                const fieldName = "reminder_type" as keyof (CreateReminderFormData | UpdateReminderFormData);
                setValue(fieldName, value as "in_app" | "email" | "sms");
              }}
            >
              <SelectTrigger className="h-9">
                <SelectValue placeholder={tDialog("channelPlaceholder")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="in_app">{tDialog("channelInApp")}</SelectItem>
                <SelectItem value="email">{tDialog("channelEmail")}</SelectItem>
                <SelectItem value="sms">{tDialog("channelSms")}</SelectItem>
              </SelectContent>
            </Select>
            {errors.reminder_type && <FieldError>{errors.reminder_type.message}</FieldError>}
          </Field>

          <Field orientation="vertical">
            <FieldLabel>{tDialog("messageLabel")}</FieldLabel>
            <Textarea
              {...register("message")}
              placeholder={tDialog("messagePlaceholder")}
              rows={3}
            />
            {errors.message && <FieldError>{errors.message.message}</FieldError>}
          </Field>

          <div className="flex justify-end gap-3 pt-4 border-t">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              {tDialog("cancel")}
            </Button>
            <Button type="submit">
              {(() => {
                if (isSubmitting) return tDialog("submitting");
                return isEdit ? tDialog("submitUpdate") : tDialog("submitCreate");
              })()}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}


