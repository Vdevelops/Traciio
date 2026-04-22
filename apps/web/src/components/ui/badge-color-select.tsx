"use client";

import * as React from "react";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

type BadgeColor = "default" | "secondary" | "outline" | "success" | "warning" | "active" | "destructive";

interface BadgeColorOption {
  value: BadgeColor;
  label: string;
}

interface BadgeColorSelectProps {
  value: BadgeColor;
  onValueChange: (value: BadgeColor) => void;
  disabled?: boolean;
  options: BadgeColorOption[];
  placeholder?: string;
}

export function BadgeColorSelect({
  value,
  onValueChange,
  disabled = false,
  options,
  placeholder,
}: BadgeColorSelectProps) {
  return (
    <Select
      value={value}
      onValueChange={(val) => onValueChange(val as BadgeColor)}
      disabled={disabled}
    >
      <SelectTrigger className="w-full">
        <SelectValue placeholder={placeholder}>
          <div className="flex items-center gap-2">
            <Badge variant={value} className="shrink-0">
              Sample
            </Badge>
            <span>
              {options.find((opt) => opt.value === value)?.label || value}
            </span>
          </div>
        </SelectValue>
      </SelectTrigger>
      <SelectContent>
        {options.map((option) => (
          <SelectItem key={option.value} value={option.value} className="cursor-pointer">
            <div className="flex items-center gap-2">
              <Badge variant={option.value} className="shrink-0">
                Sample
              </Badge>
              <span>{option.label}</span>
            </div>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
