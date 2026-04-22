"use client";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";

interface DeleteDialogProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onConfirm: () => void | Promise<void>;
  readonly title?: string;
  readonly description?: string;
  readonly itemName?: string;
  readonly isLoading?: boolean;
}

export function DeleteDialog({
  open,
  onOpenChange,
  onConfirm,
  title,
  description,
  itemName = "item",
  isLoading = false,
}: DeleteDialogProps) {
  const handleConfirm = async () => {
    await onConfirm();
    onOpenChange(false);
  };

  const handleOpenChange = (newOpen: boolean) => {
    // Prevent closing if loading
    if (isLoading && !newOpen) {
      return;
    }
    onOpenChange(newOpen);
  };

  return (
    <AlertDialog open={open} onOpenChange={handleOpenChange}>
      <AlertDialogContent
        onOverlayClick={() => {
          // Close on overlay click if not loading
          if (!isLoading) {
            handleOpenChange(false);
          }
        }}
      >
        <AlertDialogHeader>
          <AlertDialogTitle>
            {title || `Delete ${itemName}?`}
          </AlertDialogTitle>
          <AlertDialogDescription>
            {description ||
              `Are you sure you want to delete this ${itemName}? This action cannot be undone.`}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel asChild>
            <Button variant="outline" disabled={isLoading}>
              Cancel
            </Button>
          </AlertDialogCancel>
          <AlertDialogAction asChild>
            <Button
              variant="destructive"
              onClick={handleConfirm}
              disabled={isLoading}
            >
              {isLoading ? "Deleting..." : "Delete"}
            </Button>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

