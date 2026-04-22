"use client";

import { useState } from "react";
import { Package } from "lucide-react";
import { cn } from "@/lib/utils";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

function toAbsoluteImageUrl(src: string): string {
  // If the API returns a relative path like "/uploads/...", it must be loaded
  // from the API host, not from the Next.js origin.
  if (src.startsWith("/uploads/")) {
    return `${API_BASE_URL}${src}`;
  }
  return src;
}

interface ProductImageProps {
  readonly src?: string | null;
  readonly alt: string;
  readonly className?: string;
  readonly fallbackClassName?: string;
  readonly loading?: "lazy" | "eager";
  readonly size?: "sm" | "md" | "lg";
}

export function ProductImage({ 
  src, 
  alt, 
  className,
  fallbackClassName,
  loading = "lazy",
  size = "md"
}: ProductImageProps) {
  const [error, setError] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const resolvedSrc = src ? toAbsoluteImageUrl(src) : src;

  // Size mappings for consistent dimensions.
  const sizeClasses = {
    sm: "h-8 w-8",
    md: "h-10 w-10",
    lg: "h-16 w-16",
  };

  const iconSizes = {
    sm: "h-4 w-4",
    md: "h-5 w-5",
    lg: "h-8 w-8",
  };

  // Show fallback if no src or error occurred.
  const showFallback = !resolvedSrc || error;

  if (showFallback) {
    return (
      <div className={cn("flex items-center justify-center rounded bg-muted", sizeClasses[size], fallbackClassName, className)}>
        <Package className={cn(iconSizes[size], "text-muted-foreground")} />
      </div>
    );
  }

  return (
  <div key={resolvedSrc} className={cn("relative overflow-hidden rounded", sizeClasses[size], className)}>
      {isLoading && (
        <div className="absolute inset-0 flex items-center justify-center bg-muted">
          <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
        </div>
      )}
      <img
        src={resolvedSrc}
        alt={alt}
        loading={loading}
        decoding="async"
        className={cn(
          "h-full w-full object-cover transition-opacity duration-200",
          isLoading ? "opacity-0" : "opacity-100"
        )}
        onLoad={() => setIsLoading(false)}
        onError={() => {
          setError(true);
          setIsLoading(false);
        }}
      />
    </div>
  );
}
