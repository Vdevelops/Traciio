"use client";

import * as React from "react";
import { cn } from "@/lib/utils";

interface TooltipContextValue {
  open: boolean;
  setOpen: (open: boolean) => void;
}

const TooltipContext = React.createContext<TooltipContextValue | undefined>(
  undefined
);

function useTooltipContext() {
  const context = React.useContext(TooltipContext);
  if (!context) {
    throw new Error("Tooltip components must be used within TooltipProvider");
  }
  return context;
}

interface TooltipProviderProps {
  children: React.ReactNode;
  delayDuration?: number;
}

function TooltipProvider({
  children,
  delayDuration = 300,
}: TooltipProviderProps) {
  return <>{children}</>;
}

interface TooltipProps {
  children: React.ReactNode;
  delayDuration?: number;
}

function Tooltip({ children, delayDuration = 150 }: TooltipProps) {
  const [open, setOpen] = React.useState(false);
  const timeoutRef = React.useRef<NodeJS.Timeout | null>(null);

  const handleMouseEnter = () => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    timeoutRef.current = setTimeout(() => {
      setOpen(true);
    }, delayDuration);
  };

  const handleMouseLeave = () => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    // Add small delay before closing to allow moving to tooltip
    timeoutRef.current = setTimeout(() => {
      setOpen(false);
    }, 100);
  };

  React.useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <TooltipContext.Provider value={{ open, setOpen }}>
      <div 
        data-tooltip 
        onMouseEnter={handleMouseEnter} 
        onMouseLeave={handleMouseLeave} 
      >
        {children}
      </div>
    </TooltipContext.Provider>
  );
}

interface TooltipTriggerProps {
  children: React.ReactNode;
  asChild?: boolean;
}

const TooltipTrigger = React.forwardRef<
  HTMLDivElement,
  TooltipTriggerProps & React.HTMLAttributes<HTMLDivElement>
>(({ children, asChild, className, ...props }, ref) => {
  const { setOpen } = useTooltipContext();
  const internalRef = React.useRef<HTMLElement | null>(null);

  // Store ref for TooltipContent
  React.useEffect(() => {
    if (internalRef.current) {
      // Access parent Tooltip's triggerRef
      const tooltipParent = internalRef.current.closest('[data-tooltip]');
      if (tooltipParent) {
        (tooltipParent as any).triggerRef = internalRef.current;
      }
    }
  }, []);

  if (asChild && React.isValidElement(children)) {
    return React.cloneElement(children, {
      ref: (node: HTMLElement) => {
        internalRef.current = node;
        if (typeof ref === "function") {
          ref(node as HTMLDivElement);
        } else if (ref) {
          (ref as React.MutableRefObject<HTMLDivElement | null>).current = node as HTMLDivElement;
        }
      },
      onMouseEnter: () => setOpen(true),
      onMouseLeave: () => setOpen(false),
      ...props,
    } as React.HTMLAttributes<HTMLElement>);
  }

  return (
    <div
      ref={(node) => {
        internalRef.current = node;
        if (typeof ref === "function") {
          ref(node);
        } else if (ref) {
          (ref as React.MutableRefObject<HTMLDivElement | null>).current = node;
        }
      }}
      className={className}
      onMouseEnter={() => setOpen(true)}
      onMouseLeave={() => setOpen(false)}
      {...props}
    >
      {children}
    </div>
  );
});
TooltipTrigger.displayName = "TooltipTrigger";

interface TooltipContentProps extends React.HTMLAttributes<HTMLDivElement> {
  side?: "top" | "right" | "bottom" | "left";
  sideOffset?: number;
}

const TooltipContent = React.forwardRef<HTMLDivElement, TooltipContentProps>(
  ({ className, side = "right", sideOffset = 4, children, ...props }, ref) => {
    const { open, setOpen } = useTooltipContext();
    const triggerRef = React.useRef<HTMLElement | null>(null);
    const contentRef = React.useRef<HTMLDivElement>(null);

    // Keep tooltip open when hovering over it
    const handleTooltipMouseEnter = () => {
      setOpen(true);
    };

    const handleTooltipMouseLeave = () => {
      setOpen(false);
    };

    // Use requestAnimationFrame for smoother positioning
    React.useEffect(() => {
      if (!open) return;

      const updatePosition = () => {
        const tooltipParent = contentRef.current?.closest('[data-tooltip]') as HTMLElement;
        const trigger = (tooltipParent as any)?.triggerRef || tooltipParent?.querySelector('[data-tooltip-trigger]') as HTMLElement || tooltipParent;
        if (!trigger) return;

        const rect = trigger.getBoundingClientRect();
        const tooltip = contentRef.current;
        if (!tooltip) return;

        const offset = sideOffset;

        switch (side) {
          case "right":
            tooltip.style.top = `${rect.top + rect.height / 2}px`;
            tooltip.style.left = `${rect.right + offset}px`;
            tooltip.style.transform = "translate(0, -50%)";
            break;
          case "left":
            tooltip.style.top = `${rect.top + rect.height / 2}px`;
            tooltip.style.left = `${rect.left - offset}px`;
            tooltip.style.transform = "translate(-100%, -50%)";
            break;
          case "top":
            tooltip.style.top = `${rect.top - offset}px`;
            tooltip.style.left = `${rect.left + rect.width / 2}px`;
            tooltip.style.transform = "translate(-50%, -100%)";
            break;
          case "bottom":
            tooltip.style.top = `${rect.bottom + offset}px`;
            tooltip.style.left = `${rect.left + rect.width / 2}px`;
            tooltip.style.transform = "translate(-50%, 0)";
            break;
        }
      };

      // Use requestAnimationFrame for initial positioning
      requestAnimationFrame(updatePosition);

      // Throttle scroll/resize events
      let ticking = false;
      const handleUpdate = () => {
        if (!ticking) {
          requestAnimationFrame(() => {
            updatePosition();
            ticking = false;
          });
          ticking = true;
        }
      };

      window.addEventListener("scroll", handleUpdate, true);
      window.addEventListener("resize", handleUpdate);

      return () => {
        window.removeEventListener("scroll", handleUpdate, true);
        window.removeEventListener("resize", handleUpdate);
      };
    }, [open, side, sideOffset]);

    if (!open) return null;

    return (
      <div
        ref={(node) => {
          contentRef.current = node;
          if (typeof ref === "function") {
            ref(node);
          } else if (ref) {
            (ref as React.MutableRefObject<HTMLDivElement | null>).current = node;
          }
        }}
        onMouseEnter={handleTooltipMouseEnter}
        onMouseLeave={handleTooltipMouseLeave}
        className={cn(
          "fixed z-[100] overflow-hidden rounded-md border bg-popover text-popover-foreground shadow-md will-change-[top,left] p-1",
          className
        )}
        {...props}
      >
        {children}
      </div>
    );
  }
);
TooltipContent.displayName = "TooltipContent";

export { Tooltip, TooltipTrigger, TooltipContent, TooltipProvider };

