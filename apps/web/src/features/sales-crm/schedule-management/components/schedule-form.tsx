"use client";

import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  createScheduleSchema,
  updateScheduleSchema,
  type CreateScheduleFormData,
  type UpdateScheduleFormData,
} from "../schemas/schedule.schema";
import type { Schedule } from "../types";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { DateTimePicker } from "@/components/ui/date-time-picker";
import { useSchedule } from "../hooks/useSchedules";
import { useGoogleCalendarStatus, useConnectGoogleCalendar } from "@/features/profile/hooks/useGoogleCalendar";
import { useTranslations } from "next-intl";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Link2 } from "lucide-react";

interface ScheduleFormProps {
  readonly schedule?: Schedule;
  readonly scheduleId?: string;
  readonly defaultScheduledAt?: string;
  readonly onSubmit: (data: CreateScheduleFormData | UpdateScheduleFormData) => Promise<void>;
  readonly onCancel: () => void;
  readonly isLoading?: boolean;
}

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
    combined.setHours(0, 0, 0, 0);
    return combined.toISOString();
  }
  const [hours, minutes] = time.split(":").map(Number);
  const combined = new Date(date);
  combined.setHours(hours, minutes, 0, 0);
  return combined.toISOString();
}

export function ScheduleForm({
  schedule,
  scheduleId, 
  defaultScheduledAt,
  onSubmit,
  onCancel,
  isLoading 
}: ScheduleFormProps) {
  const { data: scheduleData } = useSchedule(scheduleId ?? "");
  const scheduleToUse = schedule ?? scheduleData?.data;
  const isEdit = !!scheduleToUse;
  const t = useTranslations("scheduleManagement.form");

  const { data: googleCalendarStatus } = useGoogleCalendarStatus();
  const isGoogleCalendarConnected = googleCalendarStatus?.data?.connected ?? false;
  
  const connectGoogleCalendar = useConnectGoogleCalendar();

  const schema = isEdit ? updateScheduleSchema : createScheduleSchema;
  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
    setValue,
  } = useForm<CreateScheduleFormData | UpdateScheduleFormData>({
    resolver: zodResolver(schema),
    defaultValues: isEdit && scheduleToUse
      ? {
          title: scheduleToUse.title,
          description: scheduleToUse.description ?? "",
          scheduled_at: scheduleToUse.scheduled_at || "",
          scheduled_date: scheduleToUse.scheduled_at ? extractDateFromISO(scheduleToUse.scheduled_at) : null,
          scheduled_time: scheduleToUse.scheduled_at ? extractTimeFromISO(scheduleToUse.scheduled_at) : null,
          reminder_minutes_before: scheduleToUse.reminder_minutes_before ?? undefined,
          sync_to_google_calendar: scheduleToUse.google_calendar_sync_status === "synced" || isGoogleCalendarConnected,
        }
      : {
          title: "",
          description: "",
          scheduled_at: defaultScheduledAt || "",
          scheduled_date: defaultScheduledAt ? extractDateFromISO(defaultScheduledAt) : null,
          scheduled_time: defaultScheduledAt ? extractTimeFromISO(defaultScheduledAt) : null,
          reminder_minutes_before: undefined,
          sync_to_google_calendar: isGoogleCalendarConnected,
        },
  });

  useEffect(() => {
    if (!isEdit && isGoogleCalendarConnected) {
      setValue("sync_to_google_calendar", true);
    }
  }, [isEdit, isGoogleCalendarConnected, setValue]);

  const onFormSubmit = async (data: CreateScheduleFormData | UpdateScheduleFormData) => {
    const submitData: Record<string, unknown> = { ...data };

    const scheduledDate = (data as { scheduled_date?: Date | null; scheduled_time?: string | null }).scheduled_date;
    const scheduledTime = (data as { scheduled_date?: Date | null; scheduled_time?: string | null }).scheduled_time;
    
    if (scheduledDate) {
      const combined = combineDateAndTime(scheduledDate, scheduledTime || null);
      if (combined) {
        submitData.scheduled_at = combined;
      }
    }

    delete submitData.scheduled_date;
    delete submitData.scheduled_time;

    await onSubmit(submitData as CreateScheduleFormData | UpdateScheduleFormData);
  };

  return (
    <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-4">
      <Field>
        <FieldLabel>{t("title")} *</FieldLabel>
        <Input {...register("title")} placeholder={t("titlePlaceholder")} />
        {errors.title && <FieldError>{errors.title.message}</FieldError>}
      </Field>

      <Field>
        <FieldLabel>{t("description")}</FieldLabel>
        <Textarea
          {...register("description")}
          placeholder={t("descriptionPlaceholder")}
          rows={3}
        />
        {errors.description && <FieldError>{errors.description.message}</FieldError>}
      </Field>

      <Field orientation="vertical">
        <FieldLabel>{t("scheduledAt")} *</FieldLabel>
        <DateTimePicker
          date={watch("scheduled_date") as Date | null}
          time={watch("scheduled_time") as string | null}
          onDateChange={(date, time) => {
            setValue("scheduled_date", date, { shouldValidate: true });
            setValue("scheduled_time", time, { shouldValidate: true });
          }}
          disabled={isLoading}
        />
        {(errors.scheduled_date || errors.scheduled_time) && (
          <FieldError>
            {errors.scheduled_date?.message || errors.scheduled_time?.message}
          </FieldError>
        )}
      </Field>

      <Field>
        <FieldLabel>{t("reminderMinutesBefore")}</FieldLabel>
        <Input
          type="number"
          {...register("reminder_minutes_before", { valueAsNumber: true })}
          placeholder={t("reminderMinutesBeforePlaceholder")}
          min={0}
          max={10080}
        />
        <p className="text-xs text-muted-foreground mt-1">
          {t("reminderMinutesBeforeHelp")}
        </p>
        {errors.reminder_minutes_before && (
          <FieldError>{errors.reminder_minutes_before.message}</FieldError>
        )}
      </Field>

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
        {isGoogleCalendarConnected ? (
          <p className="text-xs text-muted-foreground mt-1">
            {t("syncToGoogleCalendarHelp")}
          </p>
        ) : (
          <div className="mt-2 flex flex-col gap-2">
            <p className="text-xs text-muted-foreground">
              {t("syncToGoogleCalendarNotConnected")}
            </p>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => connectGoogleCalendar.mutate()}
              disabled={connectGoogleCalendar.isPending}
              className="w-fit cursor-pointer"
            >
              <Link2 className="mr-2 h-4 w-4" />
              {connectGoogleCalendar.isPending ? t("connecting") : t("connectGoogleCalendar")}
            </Button>
          </div>
        )}
      </Field>

      <div className="flex justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel} className="cursor-pointer">
          {t("cancel")}
        </Button>
        <Button type="submit" disabled={isLoading} className="cursor-pointer">
          {isLoading ? t("saving") : isEdit ? t("update") : t("create")}
        </Button>
      </div>
    </form>
  );
}

