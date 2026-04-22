"use client";

import * as React from "react";
import { cn } from "@/lib/utils";

export interface NumberInputProps
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, "type" | "value" | "onChange" | "onBlur"> {
  value?: number | string;
  onChange?: (value: number | undefined) => void;
  onBlur?: (e: React.FocusEvent<HTMLInputElement>) => void;
  allowDecimal?: boolean;
  decimalPlaces?: number;
}

/**
 * Format number with thousand separators (Indonesian format: 1.000)
 */
const formatNumber = (value: number | string | undefined, allowDecimal = true, decimalPlaces?: number): string => {
  if (value === undefined || value === null || value === "") {
    return "";
  }

  // Convert to string and remove existing formatting
  let numStr = String(value).replaceAll(".", "").replaceAll(",", ".");

  // Handle empty string
  if (numStr === "" || numStr === "-") {
    return numStr;
  }

  // Parse the number
  const num = Number.parseFloat(numStr);
  if (Number.isNaN(num)) {
    return "";
  }

  // Format with decimal if allowed
  let formatted: string;
  if (allowDecimal && numStr.includes(".")) {
    const parts = numStr.split(".");
    const integerPart = parts[0] || "0";
    let decimalPart = parts[1] || "";

    // Limit decimal places if specified
    if (decimalPlaces !== undefined && decimalPart.length > decimalPlaces) {
      decimalPart = decimalPart.slice(0, decimalPlaces);
    }

    // Format integer part with thousand separators
    const formattedInteger = formatIntegerPart(integerPart);
    formatted = decimalPart ? `${formattedInteger},${decimalPart}` : formattedInteger;
  } else {
    // Format integer part with thousand separators
    formatted = formatIntegerPart(numStr.split(".")[0] || "0");
  }

  // Preserve negative sign
  if (numStr.startsWith("-")) {
    formatted = `-${formatted}`;
  }

  return formatted;
};

/**
 * Format integer part with thousand separators (dots for Indonesian format)
 */
const formatIntegerPart = (integerStr: string): string => {
  // Remove negative sign temporarily
  const isNegative = integerStr.startsWith("-");
  const cleanStr = isNegative ? integerStr.slice(1) : integerStr;

  // Add thousand separators using regex (replaceAll doesn't support regex)
  // eslint-disable-next-line sonarjs/prefer-replace-all
  const formatted = cleanStr.replace(/\B(?=(\d{3})+(?!\d))/g, ".");

  return isNegative ? `-${formatted}` : formatted;
};

/**
 * Parse formatted string back to number
 */
const parseFormattedNumber = (value: string): number | undefined => {
  if (!value || value.trim() === "" || value === "-") {
    return undefined;
  }

  // Replace dot (thousand separator) with empty string
  // Replace comma (decimal separator) with dot
  const cleaned = value.replaceAll(".", "").replaceAll(",", ".");

  const parsed = Number.parseFloat(cleaned);
  return Number.isNaN(parsed) ? undefined : parsed;
};

export const NumberInput = React.forwardRef<HTMLInputElement, NumberInputProps>(
  (
    {
      className,
      value,
      onChange,
      onBlur,
      allowDecimal = true,
      decimalPlaces,
      placeholder,
      disabled,
      ...props
    },
    ref
  ) => {
    const [displayValue, setDisplayValue] = React.useState<string>("");
    const inputRef = React.useRef<HTMLInputElement | null>(null);

    // Update display value when prop value changes
    React.useEffect(() => {
      const formatted = formatNumber(value, allowDecimal, decimalPlaces);
      setDisplayValue(formatted);
    }, [value, allowDecimal, decimalPlaces]);

    // Auto-select all text on focus
    const handleFocus = (e: React.FocusEvent<HTMLInputElement>) => {
      e.target.select();
      props.onFocus?.(e);
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const inputValue = e.target.value;

      // Allow empty input
      if (inputValue === "") {
        setDisplayValue("");
        onChange?.(undefined);
        return;
      }

      // Allow negative sign at the start
      if (inputValue === "-") {
        setDisplayValue("-");
        return;
      }

      // Remove all non-numeric characters except decimal separator and negative sign
      // For Indonesian format: allow comma as decimal separator
      // Using replace with regex as replaceAll doesn't support regex patterns
      // eslint-disable-next-line sonarjs/prefer-replace-all
      let cleaned = inputValue.replace(/[^\d,.-]/g, "");

      // Only allow one decimal separator
      const commaCount = (cleaned.match(/,/g) || []).length;
      if (commaCount > 1) {
        // Keep only the first comma
        const firstCommaIndex = cleaned.indexOf(",");
        cleaned = cleaned.slice(0, firstCommaIndex + 1) + cleaned.slice(firstCommaIndex + 1).replaceAll(",", "");
      }

      // Only allow negative sign at the start
      if (cleaned.includes("-") && !cleaned.startsWith("-")) {
        cleaned = cleaned.replaceAll("-", "");
        if (cleaned.length > 0) {
          cleaned = `-${cleaned}`;
        }
      }

      // If decimal not allowed, remove decimal separator
      if (!allowDecimal) {
        cleaned = cleaned.replaceAll(",", "");
      }

      // Limit decimal places if specified
      if (allowDecimal && decimalPlaces !== undefined && cleaned.includes(",")) {
        const parts = cleaned.split(",");
        if (parts[1] && parts[1].length > decimalPlaces) {
          cleaned = `${parts[0]},${parts[1].slice(0, decimalPlaces)}`;
        }
      }

      // Format the value
      const formatted = formatNumber(cleaned, allowDecimal, decimalPlaces);
      setDisplayValue(formatted);

      // Parse and call onChange with numeric value
      const numericValue = parseFormattedNumber(formatted);
      onChange?.(numericValue);
    };

    const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
      // Re-format on blur to ensure consistent formatting
      const numericValue = parseFormattedNumber(displayValue);
      if (numericValue !== undefined) {
        const formatted = formatNumber(numericValue, allowDecimal, decimalPlaces);
        setDisplayValue(formatted);
      } else if (displayValue !== "" && displayValue !== "-") {
        // If there's invalid input, clear it
        setDisplayValue("");
        onChange?.(undefined);
      }
      onBlur?.(e);
    };

    // Merge refs
    React.useImperativeHandle(ref, () => inputRef.current as HTMLInputElement);

    return (
      <input
        ref={inputRef}
        type="text"
        inputMode="decimal"
        className={cn(
          "flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground placeholder:opacity-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
          className
        )}
        value={displayValue}
        onChange={handleChange}
        onFocus={handleFocus}
        onBlur={handleBlur}
        placeholder={placeholder}
        disabled={disabled}
        {...props}
      />
    );
  }
);

NumberInput.displayName = "NumberInput";
