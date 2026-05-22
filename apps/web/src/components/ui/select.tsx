"use client";

import * as React from "react";
import * as SelectPrimitive from "@radix-ui/react-select";
import { motion } from "framer-motion";
import { CheckIcon, ChevronDownIcon, ChevronUpIcon, SearchIcon } from "lucide-react";

import { cn } from "@/lib/utils";
import { Input } from "./input";

// Context for search functionality
const SelectContext = React.createContext<{
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  itemCount: number;
  incrementItemCount: () => void;
  resetItemCount: () => void;
  showSearch: boolean;
} | null>(null);

function Select({
  onOpenChange,
  ...props
}: React.ComponentProps<typeof SelectPrimitive.Root>) {
  const [searchQuery, setSearchQuery] = React.useState("");
  const [itemCount, setItemCount] = React.useState(0);
  const showSearch = itemCount > 5;

  const incrementItemCount = React.useCallback(() => {
    setItemCount((prev) => prev + 1);
  }, []);

  const resetItemCount = React.useCallback(() => {
    setItemCount(0);
    setSearchQuery("");
  }, []);

  // Reset when select closes
  const handleOpenChange = React.useCallback(
    (open: boolean) => {
      if (!open) {
        resetItemCount();
      }
      onOpenChange?.(open);
    },
    [onOpenChange, resetItemCount]
  );

  return (
    <SelectContext.Provider
      value={{
        searchQuery,
        setSearchQuery,
        itemCount,
        incrementItemCount,
        resetItemCount,
        showSearch,
      }}
    >
      <SelectPrimitive.Root
        data-slot="select"
        onOpenChange={handleOpenChange}
        {...props}
      />
    </SelectContext.Provider>
  );
}

function SelectGroup({
  ...props
}: React.ComponentProps<typeof SelectPrimitive.Group>) {
  return <SelectPrimitive.Group data-slot="select-group" {...props} />;
}

function SelectValue({
  ...props
}: React.ComponentProps<typeof SelectPrimitive.Value>) {
  return <SelectPrimitive.Value data-slot="select-value" {...props} />;
}

function SelectTrigger({
  className,
  size = "default",
  children,
  ...props
}: React.ComponentProps<typeof SelectPrimitive.Trigger> & {
  size?: "sm" | "default";
}) {
  return (
    <SelectPrimitive.Trigger
      data-slot="select-trigger"
      data-size={size}
      className={cn(
        "border-input data-[placeholder]:text-muted-foreground [&_svg:not([class*='text-'])]:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:bg-input/30 dark:hover:bg-input/50 flex w-full items-center justify-between gap-2 rounded-md border bg-transparent px-3 py-2 text-sm whitespace-nowrap shadow-xs transition-[color,box-shadow,transform] outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50 data-[size=default]:h-9 data-[size=sm]:h-8 *:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex *:data-[slot=select-value]:items-center *:data-[slot=select-value]:gap-2 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 active:scale-[0.98]",
        className
      )}
      {...props}
    >
      {children}
      <SelectPrimitive.Icon asChild>
        <ChevronDownIcon className="size-4 opacity-50 transition-transform duration-200 data-[state=open]:rotate-180" />
      </SelectPrimitive.Icon>
    </SelectPrimitive.Trigger>
  );
}

function SelectContent({
  className,
  children,
  position = "popper",
  align = "center",
  ...props
}: React.ComponentProps<typeof SelectPrimitive.Content>) {
  const context = React.useContext(SelectContext);
  const searchInputRef = React.useRef<HTMLInputElement>(null);
  const contentRef = React.useRef<HTMLDivElement>(null);

  // Auto focus search input when content opens and search is enabled
  React.useEffect(() => {
    if (!context?.showSearch || !searchInputRef.current) return;

    // Focus input after a short delay to ensure content is rendered
    const timer = setTimeout(() => {
      searchInputRef.current?.focus();
    }, 150);

    return () => clearTimeout(timer);
  }, [context?.showSearch]);

  return (
    <SelectPrimitive.Portal>
      <SelectPrimitive.Content
        ref={contentRef}
        data-slot="select-content"
        className={cn(
          "bg-popover text-popover-foreground data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 relative z-[1250] max-h-[var(--radix-select-content-available-height)] min-w-[8rem] origin-[var(--radix-select-content-transform-origin)] overflow-x-hidden overflow-y-auto rounded-md border shadow-md transition-all duration-150 ease-[cubic-bezier(0.4,0,0.2,1)]",
          position === "popper" &&
            "data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1",
          className
        )}
        position={position}
        align={align}
        {...props}
      >
        <SelectScrollUpButton />
        {context?.showSearch && (
          <div className="sticky top-0 z-10 bg-popover border-b p-2">
            <div className="relative">
              <SearchIcon className="absolute left-2 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
              <Input
                ref={(node) => {
                  searchInputRef.current = node;
                  // Auto focus when input is mounted
                  if (node) {
                    requestAnimationFrame(() => {
                      node.focus();
                    });
                  }
                }}
                type="text"
                placeholder="Search..."
                value={context.searchQuery}
                onChange={(e) => context.setSearchQuery(e.target.value)}
                className="pl-8 h-8 text-sm"
                onKeyDown={(e) => {
                  // Prevent select from closing when typing
                  e.stopPropagation();
                  // Prevent Enter from selecting first item
                  if (e.key === "Enter") {
                    e.preventDefault();
                  }
                }}
              />
            </div>
          </div>
        )}
        <SelectPrimitive.Viewport
          className={cn(
            "p-1",
            position === "popper" &&
              "h-[var(--radix-select-trigger-height)] w-full min-w-[var(--radix-select-trigger-width)] scroll-my-1"
          )}
        >
          {children}
        </SelectPrimitive.Viewport>
        <SelectScrollDownButton />
      </SelectPrimitive.Content>
    </SelectPrimitive.Portal>
  );
}

function SelectLabel({
  className,
  ...props
}: React.ComponentProps<typeof SelectPrimitive.Label>) {
  return (
    <SelectPrimitive.Label
      data-slot="select-label"
      className={cn("text-muted-foreground px-2 py-1.5 text-xs", className)}
      {...props}
    />
  );
}

function SelectItem({
  className,
  children,
  value,
  ...props
}: React.ComponentProps<typeof SelectPrimitive.Item>) {
  const context = React.useContext(SelectContext);
  const hasIncrementedRef = React.useRef(false);
  const lastItemCountRef = React.useRef(0);
  
  const itemText = React.useMemo(() => {
    if (typeof children === "string") return children;
    // Extract text from React children
    const extractText = (node: React.ReactNode): string => {
      if (typeof node === "string" || typeof node === "number") {
        return String(node);
      }
      if (React.isValidElement(node)) {
        const props = node.props as { children?: React.ReactNode };
        if (props.children) {
          return extractText(props.children);
        }
      }
      if (Array.isArray(node)) {
        return node.map(extractText).join("");
      }
      return "";
    };
    return extractText(children);
  }, [children]);

  // Increment count on mount (only once per render cycle)
  // Reset tracking when itemCount resets (select closed and reopened)
  React.useEffect(() => {
    if (context) {
      const currentCount = context.itemCount;
      // If itemCount was reset (went back to 0), reset our tracking
      if (currentCount === 0 && lastItemCountRef.current > 0) {
        hasIncrementedRef.current = false;
      }
      // Only increment if we haven't counted this render cycle
      if (!hasIncrementedRef.current) {
        context.incrementItemCount();
        hasIncrementedRef.current = true;
        lastItemCountRef.current = context.itemCount;
      }
    }
  }, [context]);

  // Filter based on search query
  const shouldShow = React.useMemo(() => {
    if (!context?.showSearch || !context.searchQuery) return true;
    const query = context.searchQuery.toLowerCase();
    return itemText.toLowerCase().includes(query) || 
           (value && String(value).toLowerCase().includes(query));
  }, [context?.showSearch, context?.searchQuery, itemText, value]);

  if (!shouldShow) return null;

  return (
    <SelectPrimitive.Item
      data-slot="select-item"
      className={cn(
        "focus:bg-accent focus:text-accent-foreground [&_svg:not([class*='text-'])]:text-muted-foreground relative flex w-full cursor-default items-center gap-2 rounded-sm py-1.5 pr-8 pl-2 text-sm outline-hidden select-none data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 *:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2 transition-colors",
        className
      )}
      value={value}
      {...props}
    >
      <span className="absolute right-2 flex size-3.5 items-center justify-center">
        <SelectPrimitive.ItemIndicator>
          <motion.div
            initial={{ scale: 0, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{ duration: 0.15, ease: [0.4, 0, 0.2, 1] }}
          >
            <CheckIcon className="size-4" />
          </motion.div>
        </SelectPrimitive.ItemIndicator>
      </span>
      <SelectPrimitive.ItemText>{children}</SelectPrimitive.ItemText>
    </SelectPrimitive.Item>
  );
}

function SelectSeparator({
  className,
  ...props
}: React.ComponentProps<typeof SelectPrimitive.Separator>) {
  return (
    <SelectPrimitive.Separator
      data-slot="select-separator"
      className={cn("bg-border pointer-events-none -mx-1 my-1 h-px", className)}
      {...props}
    />
  );
}

function SelectScrollUpButton({
  className,
  ...props
}: React.ComponentProps<typeof SelectPrimitive.ScrollUpButton>) {
  return (
    <SelectPrimitive.ScrollUpButton
      data-slot="select-scroll-up-button"
      className={cn(
        "flex cursor-default items-center justify-center py-1",
        className
      )}
      {...props}
    >
      <ChevronUpIcon className="size-4" />
    </SelectPrimitive.ScrollUpButton>
  );
}

function SelectScrollDownButton({
  className,
  ...props
}: React.ComponentProps<typeof SelectPrimitive.ScrollDownButton>) {
  return (
    <SelectPrimitive.ScrollDownButton
      data-slot="select-scroll-down-button"
      className={cn(
        "flex cursor-default items-center justify-center py-1",
        className
      )}
      {...props}
    >
      <ChevronDownIcon className="size-4" />
    </SelectPrimitive.ScrollDownButton>
  );
}

export {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectScrollDownButton,
  SelectScrollUpButton,
  SelectSeparator,
  SelectTrigger,
  SelectValue,
};
