"use client";

import { useState, useMemo, useCallback, useRef } from "react";
import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Input } from "@/components/ui/input";
import { Check, Search, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { getAvailableIcons, searchIcons, DynamicIcon } from "@/lib/icon-utils";
import { useDebounce } from "@/hooks/use-debounce";

interface IconPickerProps {
  readonly value?: string;
  readonly onValueChange: (value: string) => void;
  readonly placeholder?: string;
  readonly disabled?: boolean;
}

// Constants for performance optimization
const ICONS_PER_PAGE = 60; // Render 60 icons at a time
const INITIAL_ICONS_COUNT = 60; // Show first 60 icons initially
const DEBOUNCE_DELAY = 300; // Debounce delay in ms

export function IconPicker({ value, onValueChange, placeholder = "Select icon", disabled }: IconPickerProps) {
  const [open, setOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [scrollOffset, setScrollOffset] = useState(0); // Track scroll-based increment
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const selectedIcon = value ?? "";

  // Debounce search query to avoid lag
  const debouncedSearchQuery = useDebounce(searchQuery, DEBOUNCE_DELAY);

  // Memoize icon names list (only compute once)
  const iconNames = useMemo(() => getAvailableIcons(), []);

  // Filter icons based on debounced search query
  const filteredIcons = useMemo(() => {
    if (!debouncedSearchQuery.trim()) {
      return iconNames;
    }
    return searchIcons(debouncedSearchQuery);
  }, [iconNames, debouncedSearchQuery]);

  // Reset scroll offset when search changes
  const searchKey = debouncedSearchQuery.trim() || "all";
  
  // Calculate current displayed count: initial + scroll increment
  // Reset scroll offset when search key changes (via key prop on container)
  const currentDisplayedCount = INITIAL_ICONS_COUNT + scrollOffset;

  // Limit displayed icons for performance
  const displayedIcons = useMemo(() => {
    return filteredIcons.slice(0, currentDisplayedCount);
  }, [filteredIcons, currentDisplayedCount]);

  // Load more icons on scroll
  const handleScroll = useCallback((e: React.UIEvent<HTMLDivElement>) => {
    const target = e.currentTarget;
    const scrollBottom = target.scrollHeight - target.scrollTop - target.clientHeight;
    
    // Load more when user scrolls near bottom (within 100px)
    const maxCount = filteredIcons.length;
    if (scrollBottom < 100 && currentDisplayedCount < maxCount) {
      const newOffset = Math.min(scrollOffset + ICONS_PER_PAGE, maxCount - INITIAL_ICONS_COUNT);
      if (newOffset > scrollOffset) {
        setScrollOffset(newOffset);
      }
    }
  }, [currentDisplayedCount, filteredIcons.length, scrollOffset]);

  // Handle popover open change
  const handleOpenChange = useCallback((newOpen: boolean) => {
    setOpen(newOpen);
    if (newOpen) {
      // Reset scroll offset when opening
      setScrollOffset(0);
    } else {
      // Reset state when closing
      setSearchQuery("");
      setScrollOffset(0);
    }
  }, []);

  return (
    <Popover open={open} onOpenChange={handleOpenChange}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          aria-expanded={open}
          className="w-full justify-between"
          disabled={disabled}
        >
          <div className="flex items-center gap-2">
            {selectedIcon ? (
              <>
                <DynamicIcon name={selectedIcon} className="h-4 w-4" />
                <span className="font-normal">{selectedIcon}</span>
              </>
            ) : (
              <span className="text-muted-foreground">{placeholder}</span>
            )}
          </div>
          <Search className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[500px] p-0" align="start">
        <div className="p-2 border-b">
          <div className="relative">
            <Search className="absolute left-2 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search icons..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-8"
            />
          </div>
        </div>
        <div 
          key={searchKey}
          ref={scrollContainerRef}
          className="max-h-[400px] overflow-y-auto p-2"
          onScroll={handleScroll}
        >
          {filteredIcons.length === 0 ? (
            <div className="text-center text-muted-foreground py-8 text-sm">
              {searchQuery.trim() && debouncedSearchQuery !== searchQuery ? (
                <div className="flex items-center justify-center gap-2">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  <span>Searching...</span>
                </div>
              ) : (
                "No icon found."
              )}
            </div>
          ) : (
            <>
              <div className="grid grid-cols-6 gap-2">
                {displayedIcons.map((iconName) => {
                  const isSelected = selectedIcon === iconName;
                  return (
                    <button
                      key={iconName}
                      type="button"
                      onClick={() => {
                        onValueChange(iconName);
                        setOpen(false);
                        setSearchQuery("");
                      }}
                      className={cn(
                        "relative flex flex-col items-center justify-center p-3 cursor-pointer rounded-md border transition-colors",
                        isSelected
                          ? "bg-primary text-primary-foreground border-primary"
                          : "hover:bg-accent hover:text-accent-foreground border-transparent"
                      )}
                    >
                      <div className="mb-1">
                        <DynamicIcon name={iconName} className="h-6 w-6" />
                      </div>
                      <div className="text-[10px] text-center truncate w-full">
                        {iconName}
                      </div>
                      {isSelected && (
                        <Check className="absolute top-1 right-1 h-3 w-3" />
                      )}
                    </button>
                  );
                })}
              </div>
              {currentDisplayedCount < filteredIcons.length && (
                <div className="text-center text-muted-foreground py-4 text-xs">
                  Showing {currentDisplayedCount} of {filteredIcons.length} icons. Scroll for more...
                </div>
              )}
            </>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
}
