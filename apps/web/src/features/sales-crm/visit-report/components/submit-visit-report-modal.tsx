"use client";

import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { submitVisitReportSchema } from "../schemas/visit-report.schema";
import { useSubmitVisitReport } from "../hooks/useVisitReports";
import type { SubmitVisitReportFormData } from "../types";
import { cn } from "@/lib/utils";

interface SubmitVisitReportModalProps {
  readonly visitReportId: string;
  readonly isOpen: boolean;
  readonly onClose: () => void;
  readonly onSuccess?: () => void;
}

const outcomeOptions = [
  {
    value: "very_positive" as const,
    label: "Very Positive",
    description: "Strong interest, ready to proceed",
  },
  {
    value: "positive" as const,
    label: "Positive",
    description: "Good progress, potential opportunity",
  },
  {
    value: "neutral" as const,
    label: "Neutral",
    description: "Need more information or follow-up",
  },
  {
    value: "negative" as const,
    label: "Negative",
    description: "Not interested or not a fit",
  },
];

export function SubmitVisitReportModal({
  visitReportId,
  isOpen,
  onClose,
  onSuccess,
}: SubmitVisitReportModalProps) {
  const submitMutation = useSubmitVisitReport();

  const {
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<SubmitVisitReportFormData>({
    resolver: zodResolver(submitVisitReportSchema),
    defaultValues: {
      outcome: "neutral",
      next_steps: "",
    },
  });

  const onSubmit = async (data: SubmitVisitReportFormData) => {
    try {
      await submitMutation.mutateAsync({
        id: visitReportId,
        data,
      });

      toast.success("Visit report submitted successfully", {
        description: "The report has been sent for manager approval.",
      });

      reset();
      onClose();
      onSuccess?.();
    } catch (error) {
      const errorMessage =
        error instanceof Error
          ? error.message
          : "Failed to submit visit report";
      toast.error("Submission failed", {
        description: errorMessage,
      });
    }
  };

  const handleClose = () => {
    if (!submitMutation.isPending) {
      reset();
      onClose();
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Submit Visit Report</DialogTitle>
          <DialogDescription>
            Please provide the outcome of this visit and outline the next steps
            before submitting for approval.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div className="space-y-3">
            <Label htmlFor="outcome">Visit Outcome *</Label>
            <Controller
              name="outcome"
              control={control}
              render={({ field }) => (
                <div className="flex flex-col space-y-2">
                  {outcomeOptions.map((option) => (
                    <label
                      key={option.value}
                      className={cn(
                        "flex items-start space-x-3 rounded-lg border p-3 cursor-pointer transition-colors",
                        field.value === option.value
                          ? "border-primary bg-primary/5"
                          : "border-border hover:bg-accent"
                      )}
                    >
                      <input
                        type="radio"
                        value={option.value}
                        checked={field.value === option.value}
                        onChange={(e) => field.onChange(e.target.value)}
                        className="mt-1 cursor-pointer"
                      />
                      <div className="flex-1">
                        <div className="font-medium">{option.label}</div>
                        <div className="text-sm text-muted-foreground">
                          {option.description}
                        </div>
                      </div>
                    </label>
                  ))}
                </div>
              )}
            />
            {errors.outcome?.message && (
              <p className="text-sm text-destructive">
                {errors.outcome.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="next_steps">Next Steps *</Label>
            <Controller
              name="next_steps"
              control={control}
              render={({ field }) => (
                <Textarea
                  {...field}
                  id="next_steps"
                  placeholder="Describe the recommended next steps based on this visit..."
                  className="min-h-[120px] resize-none"
                />
              )}
            />
            {errors.next_steps?.message && (
              <p className="text-sm text-destructive">
                {errors.next_steps.message}
              </p>
            )}
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={submitMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={submitMutation.isPending}
              className="cursor-pointer"
            >
              {submitMutation.isPending
                ? "Submitting..."
                : "Submit for Approval"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

