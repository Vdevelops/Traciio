"use client";

import * as React from "react";
import { X } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { cn } from "@/lib/utils";
import { Button } from "./button";

interface DrawerProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly children: React.ReactNode;
  readonly title?: string;
  readonly description?: string;
  readonly side?: "left" | "right" | "top" | "bottom";
  readonly className?: string;
  readonly defaultWidth?: number;
  readonly minWidth?: number;
  readonly maxWidth?: number;
  readonly resizable?: boolean;
}

export function Drawer({
  open,
  onOpenChange,
  children,
  title,
  description,
  side = "right",
  className,
  defaultWidth = 672, // max-w-2xl = 672px
  minWidth = 320,
  maxWidth = 1200,
  resizable = true,
}: DrawerProps) {
  const [width, setWidth] = React.useState(defaultWidth);
  const [isResizing, setIsResizing] = React.useState(false);
  const [isHoveringResizeArea, setIsHoveringResizeArea] = React.useState(false);
  const drawerRef = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    if (open) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }
    return () => {
      document.body.style.overflow = "";
    };
  }, [open]);

  React.useEffect(() => {
    if (open) {
      setWidth(defaultWidth);
    }
  }, [open, defaultWidth]);

  const handleMouseDown = React.useCallback(
    (e: React.MouseEvent) => {
      if (!resizable || (side !== "left" && side !== "right")) return;
      e.preventDefault();
      setIsResizing(true);
    },
    [resizable, side]
  );

  React.useEffect(() => {
    if (!isResizing) return;

    const handleMouseMove = (e: MouseEvent) => {
      if (!drawerRef.current) return;

      let newWidth: number;
      if (side === "right") {
        newWidth = window.innerWidth - e.clientX;
      } else {
        // side === "left"
        newWidth = e.clientX;
      }

      // Clamp width between min and max
      newWidth = Math.max(minWidth, Math.min(maxWidth, newWidth));
      setWidth(newWidth);
    };

    const handleMouseUp = () => {
      setIsResizing(false);
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);

    return () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
    };
  }, [isResizing, side, minWidth, maxWidth]);

  const canResize = resizable && (side === "left" || side === "right");

  const getDefaultWidthClass = () => {
    if (canResize) return undefined;
    if (side === "left" || side === "right") return "w-full max-w-2xl";
    return "h-full max-h-[80vh]";
  };

  const variants = {
    left: {
      initial: { x: "-100%" },
      animate: { x: 0 },
      exit: { x: "-100%" },
    },
    right: {
      initial: { x: "100%" },
      animate: { x: 0 },
      exit: { x: "100%" },
    },
    top: {
      initial: { y: "-100%" },
      animate: { y: 0 },
      exit: { y: "-100%" },
    },
    bottom: {
      initial: { y: "100%" },
      animate: { y: 0 },
      exit: { y: "100%" },
    },
  };

  const sideClasses = {
    left: "left-0 top-0 h-full",
    right: "right-0 top-0 h-full",
    top: "top-0 left-0 w-full",
    bottom: "bottom-0 left-0 w-full",
  };

  return (
    <AnimatePresence>
      {open && (
        <>
          {/* Overlay */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="fixed inset-0 z-[1030] bg-black/50 h-full cursor-pointer"
            onClick={() => onOpenChange(false)}
          />

          {/* Drawer */}
          <motion.div
            ref={drawerRef}
            variants={variants[side]}
            initial="initial"
            animate="animate"
            exit="exit"
            transition={{ type: "spring", damping: 25, stiffness: 200 }}
            style={
              canResize
                ? {
                    width: `${width}px`,
                    maxWidth: "none",
                  }
                : undefined
            }
            className={cn(
              "fixed z-[1031] bg-background border shadow-lg",
              getDefaultWidthClass(),
              sideClasses[side],
              className
            )}
          >
            {/* Resize Handle */}
            {canResize && (
              <button
                type="button"
                aria-label="Resize drawer"
                onMouseDown={handleMouseDown}
                onMouseEnter={() => setIsHoveringResizeArea(true)}
                onMouseLeave={() => setIsHoveringResizeArea(false)}
                className={cn(
                  "absolute top-0 bottom-0 w-1 transition-colors cursor-col-resize border-0 bg-transparent p-0 outline-none focus:outline-none",
                  side === "right" ? "left-0" : "right-0",
                  isHoveringResizeArea || isResizing
                    ? "bg-primary/50"
                    : "bg-transparent hover:bg-primary/30"
                )}
              />
            )}

            <div className="flex flex-col h-full">
              {/* Header */}
              {(title || description) && (
                <div className="flex items-center justify-between border-b px-6 py-4 shrink-0">
                  <div className="flex-1 min-w-0">
                    {title && (
                      <h2 className="text-lg font-semibold truncate">{title}</h2>
                    )}
                    {description && (
                      <p className="text-sm text-muted-foreground mt-1">
                        {description}
                      </p>
                    )}
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => onOpenChange(false)}
                    className="h-8 w-8 shrink-0 ml-4"
                  >
                    <X className="h-4 w-4" />
                    <span className="sr-only">Close</span>
                  </Button>
                </div>
              )}

              {/* Content */}
              <div className="flex-1 overflow-y-auto px-6 py-4">
                {children}
              </div>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}

