"use client";

import { useEffect, useState } from "react";
import { Switch } from "@/components/ui/switch";
import { cn } from "@/lib/utils";

interface StatusSwitchProps {
  checked: boolean;
  onCheckedChange: (checked: boolean) => void | Promise<unknown>;
  disabled?: boolean;
  className?: string;
}

export function StatusSwitch({
  checked,
  onCheckedChange,
  disabled,
  className,
}: StatusSwitchProps) {
  const [optimisticChecked, setOptimisticChecked] = useState(checked);
  const [isUpdating, setIsUpdating] = useState(false);

  useEffect(() => {
    setOptimisticChecked(checked);
  }, [checked]);

  const handleCheckedChange = async (nextChecked: boolean) => {
    const previousChecked = optimisticChecked;
    setOptimisticChecked(nextChecked);
    setIsUpdating(true);

    try {
      await onCheckedChange(nextChecked);
    } catch {
      setOptimisticChecked(previousChecked);
    } finally {
      setIsUpdating(false);
    }
  };

  return (
    <div
      className={cn(
        "flex items-center justify-center w-full h-full",
        className,
      )}
    >
      <Switch
        checked={optimisticChecked}
        onCheckedChange={handleCheckedChange}
        disabled={disabled || isUpdating}
        className="cursor-pointer data-disabled:cursor-not-allowed"
      />
    </div>
  );
}
