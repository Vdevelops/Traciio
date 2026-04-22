"use client";

import * as React from "react";
import { useEffect, useRef } from "react";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";

export interface CommandPaletteItem {
  readonly name: string;
  readonly href: string;
  readonly icon: React.ReactNode;
  readonly group: string;
}

interface CommandPaletteProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly items: CommandPaletteItem[];
  readonly onSelectItem: (href: string) => void;
}

export function CommandPalette({
  open,
  onOpenChange,
  items,
  onSelectItem,
}: CommandPaletteProps) {
  const inputRef = useRef<HTMLInputElement>(null);

  // Auto-focus search input when opened
  useEffect(() => {
    if (open) {
      // Small delay to ensure the component is fully rendered
      const timer = setTimeout(() => {
        inputRef.current?.focus();
      }, 0);
      return () => clearTimeout(timer);
    }
  }, [open]);

  // Handle escape key
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape" && open) {
        onOpenChange(false);
      }
    };

    if (open) {
      document.addEventListener("keydown", handleKeyDown);
      return () => document.removeEventListener("keydown", handleKeyDown);
    }
  }, [open, onOpenChange]);

  if (!open) return null;

  // Group items by category
  const groupedItems = items.reduce<Record<string, CommandPaletteItem[]>>(
    (acc, item) => {
      const group = item.group || "Main";
      if (!acc[group]) {
        acc[group] = [];
      }
      acc[group].push(item);
      return acc;
    },
    {}
  );

  return (
    <div
      className="fixed inset-0 z-50 flex items-start justify-center pt-[20vh] animate-in fade-in-0"
      onClick={() => onOpenChange(false)}
      role="dialog"
      aria-modal="true"
      aria-label="Command palette"
    >
      {/* Backdrop with blur */}
      <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" />

      {/* Command Palette Container */}
      <div
        className="relative w-full max-w-2xl mx-4 bg-background border border-border rounded-lg shadow-2xl animate-in zoom-in-95 slide-in-from-top-2"
        onClick={(e) => e.stopPropagation()}
      >
        <Command className="bg-transparent border-0 shadow-none">
          <div className="border-b border-border">
            <CommandInput
              ref={inputRef}
              placeholder="Type a command or search..."
              className="h-12 border-0 bg-transparent focus:ring-0 text-foreground placeholder:text-muted-foreground"
            />
          </div>
          <CommandList className="max-h-[60vh] overflow-y-auto scrollbar-thin">
            <CommandEmpty className="py-6 text-center text-sm text-muted-foreground">
              No menu found.
            </CommandEmpty>
            {Object.entries(groupedItems).map(([group, groupItems]) => (
              <CommandGroup key={group} heading={group} className="px-2 py-2">
                {groupItems.map((item) => (
                  <CommandItem
                    key={`${group}-${item.href}`}
                    value={item.name}
                    onSelect={() => onSelectItem(item.href)}
                    className="mx-1 mb-0.5 rounded-md px-3 py-2.5 cursor-pointer aria-selected:bg-accent aria-selected:text-accent-foreground hover:bg-accent/50 transition-colors"
                  >
                    <span className="mr-3 text-muted-foreground">
                      {item.icon}
                    </span>
                    <span className="flex-1 truncate text-sm font-medium text-foreground">
                      {item.name}
                    </span>
                    <span className="text-xs text-muted-foreground/40 ml-2">
                      {item.href}
                    </span>
                  </CommandItem>
                ))}
              </CommandGroup>
            ))}
          </CommandList>
        </Command>
      </div>
    </div>
  );
}
