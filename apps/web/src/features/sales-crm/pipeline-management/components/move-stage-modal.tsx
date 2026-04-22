"use client";

import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { AlertCircle } from "lucide-react";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { moveStageSchema, type MoveStageFormData } from "../schemas/deal.schema";
import { useMoveStage } from "../hooks/useDeals";
import type { PipelineStage } from "../types";

interface MoveStageModalProps {
  readonly dealId: string;
  readonly currentStageId: string;
  readonly availableStages: readonly PipelineStage[];
  readonly isOpen: boolean;
  readonly onClose: () => void;
  readonly onSuccess?: () => void;
}

export function MoveStageModal({
  dealId,
  currentStageId,
  availableStages,
  isOpen,
  onClose,
  onSuccess,
}: MoveStageModalProps) {
  const moveStageMutation = useMoveStage();

  const {
    control,
    handleSubmit,
    reset,
    watch,
    formState: { errors },
  } = useForm<MoveStageFormData>({
    resolver: zodResolver(moveStageSchema),
    defaultValues: {
      to_stage_id: "",
      reason: "",
      notes: "",
    },
  });

  const selectedStageId = watch("to_stage_id");

  // Filter out current stage and get only next valid stages
  const nextStages = availableStages.filter(
    (stage) => stage.id !== currentStageId
  );

  const selectedStage = availableStages.find((s) => s.id === selectedStageId);

  const onSubmit = async (data: MoveStageFormData) => {
    try {
      await moveStageMutation.mutateAsync({
        deal_id: dealId,
        stage_id: data.to_stage_id,
      });

      toast.success("Deal stage updated", {
        description: `Moved to ${selectedStage?.name ?? "new stage"}`,
      });

      reset();
      onClose();
      onSuccess?.();
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Failed to move deal stage";
      toast.error("Stage movement failed", {
        description: errorMessage,
      });
    }
  };

  const handleClose = () => {
    if (!moveStageMutation.isPending) {
      reset();
      onClose();
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Move Deal to Different Stage</DialogTitle>
          <DialogDescription>
            Select the next stage for this deal and provide a reason for the
            movement.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {nextStages.length === 0 ? (
            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                No available stages to move to. This deal might already be in
                the final stage.
              </AlertDescription>
            </Alert>
          ) : (
            <>
              <div className="space-y-2">
                <Label htmlFor="to_stage_id">Target Stage *</Label>
                <Controller
                  name="to_stage_id"
                  control={control}
                  render={({ field }) => (
                    <Select
                      onValueChange={field.onChange}
                      value={field.value}
                    >
                      <SelectTrigger id="to_stage_id">
                        <SelectValue placeholder="Select target stage" />
                      </SelectTrigger>
                      <SelectContent>
                        {nextStages.map((stage) => (
                          <SelectItem
                            key={stage.id}
                            value={stage.id}
                            className="cursor-pointer"
                          >
                            <div className="flex items-center space-x-2">
                              <span>{stage.name}</span>
                              {stage.probability !== undefined &&
                                stage.probability !== null && (
                                  <span className="text-xs text-muted-foreground">
                                    ({stage.probability}% probability)
                                  </span>
                                )}
                            </div>
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                />
                {errors.to_stage_id?.message && (
                  <p className="text-sm text-destructive">
                    {errors.to_stage_id.message}
                  </p>
                )}
              </div>

              {selectedStage?.requirements && (
                <Alert>
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    <div className="font-medium mb-1">
                      Stage Requirements:
                    </div>
                    <div className="text-sm">
                      {selectedStage.requirements}
                    </div>
                  </AlertDescription>
                </Alert>
              )}

              <div className="space-y-2">
                <Label htmlFor="reason">Reason for Movement *</Label>
                <Controller
                  name="reason"
                  control={control}
                  render={({ field }) => (
                    <Textarea
                      {...field}
                      id="reason"
                      placeholder="Explain why this deal is moving to the next stage..."
                      className="min-h-[80px] resize-none"
                    />
                  )}
                />
                {errors.reason?.message && (
                  <p className="text-sm text-destructive">
                    {errors.reason.message}
                  </p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="notes">Additional Notes</Label>
                <Controller
                  name="notes"
                  control={control}
                  render={({ field }) => (
                    <Textarea
                      {...field}
                      value={field.value ?? ""}
                      id="notes"
                      placeholder="Any additional notes about this stage transition..."
                      className="min-h-[80px] resize-none"
                    />
                  )}
                />
                {errors.notes?.message && (
                  <p className="text-sm text-destructive">
                    {errors.notes.message}
                  </p>
                )}
              </div>
            </>
          )}

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={moveStageMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={
                moveStageMutation.isPending || nextStages.length === 0
              }
              className="cursor-pointer"
            >
              {moveStageMutation.isPending ? "Moving..." : "Move Stage"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
