"use client";

import * as React from "react";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

type BadgeColor = "default" | "secondary" | "outline" | "success" | "warning" | "active" | "destructive";

interface BadgeColorOption {
  value: BadgeColor;
  label: string;
  description?: string;
}

interface BadgeColorPickerProps {
  value: BadgeColor;
  onValueChange: (value: BadgeColor) => void;
  disabled?: boolean;
  options: BadgeColorOption[];
}

export function BadgeColorPicker({
  value,
  onValueChange,
  disabled = false,
  options,
}: BadgeColorPickerProps) {
  return (
    <div className="grid grid-cols-2 gap-3">
      {options.map((option) => {
        const isSelected = value === option.value;
        
        return (
          <button
            key={option.value}
            type="button"
            onClick={() => onValueChange(option.value)}
            disabled={disabled}
            className={cn(
              "flex items-center gap-3 p-3 rounded-lg border-2 transition-all cursor-pointer",
              "hover:bg-accent/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
              isSelected
                ? "border-primary bg-accent/30"
                : "border-border bg-background",
              disabled && "opacity-50 cursor-not-allowed"
            )}
          >
            <Badge variant={option.value} className="shrink-0">
              Sample
            </Badge>
            <div className="flex flex-col items-start text-left flex-1 min-w-0">
              <span className="text-sm font-medium text-foreground">
                {option.label}
              </span>
              {option.description && (
                <span className="text-xs text-muted-foreground">
                  {option.description}
                </span>
              )}
            </div>
            {isSelected && (
              <div className="shrink-0 w-4 h-4 rounded-full bg-primary flex items-center justify-center">
                <svg
                  className="w-3 h-3 text-primary-foreground"
                  fill="none"
                  strokeWidth="2"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </div>
            )}
          </button>
        );
      })}
    </div>
  );
}
