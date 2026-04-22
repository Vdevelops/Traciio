"use client";

import { Switch } from "@/components/ui/switch";
import { cn } from "@/lib/utils";

interface StatusSwitchProps {
  checked: boolean;
  onCheckedChange: (checked: boolean) => void;
  disabled?: boolean;
  className?: string;
}

export function StatusSwitch({
  checked,
  onCheckedChange,
  disabled,
  className,
}: StatusSwitchProps) {
  return (
    <div className={cn("flex items-center justify-center w-full h-full", className)}>
      <Switch
        checked={checked}
        onCheckedChange={onCheckedChange}
        disabled={disabled}
        className="cursor-pointer data-disabled:cursor-not-allowed"
      />
    </div>
  );
}
