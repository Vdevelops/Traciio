"use client";

import * as React from "react";
import {
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronFirstIcon,
  ChevronLastIcon,
  MoreHorizontalIcon,
} from "lucide-react";

import { cn } from "@/lib/utils";

function Pagination({ className, ...props }: React.ComponentProps<"nav">) {
  return (
    <nav
      role="navigation"
      aria-label="pagination"
      data-slot="pagination"
      className={cn("mx-auto flex w-full justify-center", className)}
      {...props}
    />
  );
}

function PaginationContent({
  className,
  ...props
}: React.ComponentProps<"ul">) {
  return (
    <ul
      data-slot="pagination-content"
      className={cn("flex flex-row items-center gap-1.5", className)}
      {...props}
    />
  );
}

function PaginationItem({ ...props }: React.ComponentProps<"li">) {
  return <li data-slot="pagination-item" {...props} />;
}

type PaginationLinkProps = {
  isActive?: boolean;
} & React.ComponentProps<"button">;

function PaginationLink({
  className,
  isActive,
  ...props
}: PaginationLinkProps) {
  return (
    <button
      aria-current={isActive ? "page" : undefined}
      data-slot="pagination-link"
      data-active={isActive}
      className={cn(
        "relative inline-flex items-center justify-center rounded-md text-sm font-medium transition-all duration-200",
        "h-9 min-w-9 px-3",
        "hover:bg-accent hover:text-accent-foreground",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        "disabled:pointer-events-none disabled:opacity-50",
        "cursor-pointer",
        isActive
          ? "bg-primary text-primary-foreground shadow-sm hover:bg-primary/90"
          : "text-muted-foreground",
        className
      )}
      {...props}
    />
  );
}

function PaginationPrevious({
  className,
  ...props
}: React.ComponentProps<typeof PaginationLink>) {
  return (
    <PaginationLink
      aria-label="Go to previous page"
      className={cn(
        "gap-1.5 px-3 sm:px-4 h-9 rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground transition-all duration-200",
        className
      )}
      {...props}
    >
      <ChevronLeftIcon className="h-4 w-4" />
      <span className="hidden sm:inline">Previous</span>
    </PaginationLink>
  );
}

function PaginationNext({
  className,
  ...props
}: React.ComponentProps<typeof PaginationLink>) {
  return (
    <PaginationLink
      aria-label="Go to next page"
      className={cn(
        "gap-1.5 px-3 sm:px-4 h-9 rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground transition-all duration-200",
        className
      )}
      {...props}
    >
      <span className="hidden sm:inline">Next</span>
      <ChevronRightIcon className="h-4 w-4" />
    </PaginationLink>
  );
}

function PaginationFirst({
  className,
  ...props
}: React.ComponentProps<typeof PaginationLink>) {
  return (
    <PaginationLink
      aria-label="Go to first page"
      className={cn(
        "h-9 w-9 p-0 rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground transition-all duration-200",
        className
      )}
      {...props}
    >
      <ChevronFirstIcon className="h-4 w-4" />
    </PaginationLink>
  );
}

function PaginationLast({
  className,
  ...props
}: React.ComponentProps<typeof PaginationLink>) {
  return (
    <PaginationLink
      aria-label="Go to last page"
      className={cn(
        "h-9 w-9 p-0 rounded-md border border-input bg-background hover:bg-accent hover:text-accent-foreground transition-all duration-200",
        className
      )}
      {...props}
    >
      <ChevronLastIcon className="h-4 w-4" />
    </PaginationLink>
  );
}

function PaginationEllipsis({
  className,
  ...props
}: React.ComponentProps<"span">) {
  return (
    <span
      aria-hidden
      data-slot="pagination-ellipsis"
      className={cn(
        "flex h-9 w-9 items-center justify-center text-muted-foreground",
        className
      )}
      {...props}
    >
      <MoreHorizontalIcon className="h-4 w-4" />
      <span className="sr-only">More pages</span>
    </span>
  );
}

export {
  Pagination,
  PaginationContent,
  PaginationLink,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
  PaginationFirst,
  PaginationLast,
  PaginationEllipsis,
};

