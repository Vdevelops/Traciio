"use client";

import { DynamicIcon, isValidIcon } from "@/lib/icon-utils";
import type { IconProps } from "@/lib/icon-utils";

/**
 * Render icon component from icon name string
 * @param iconName - Name of the icon (e.g., "Activity", "Phone", "Mail")
 * @param className - Optional className for the icon
 * @returns React node with the icon or null if icon not found
 */
export function renderIcon(iconName: string | null | undefined, className?: string): React.ReactNode {
  if (!iconName) {
    return null;
  }

  if (!isValidIcon(iconName)) {
    return null;
  }

  return <DynamicIcon name={iconName} className={className} />;
}

/**
 * Check if an icon name exists in lucide-react
 * @param iconName - Name of the icon to check
 * @returns true if icon exists, false otherwise
 */
export function iconExists(iconName: string | null | undefined): boolean {
  if (!iconName) {
    return false;
  }

  return isValidIcon(iconName);
}


