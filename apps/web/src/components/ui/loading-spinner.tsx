import { cn } from "@/lib/utils";

interface LoadingSpinnerProps {
  readonly size?: "sm" | "md" | "lg";
  readonly className?: string;
}

const sizeClasses = {
  sm: "h-4 w-4 border-2",
  md: "h-8 w-8 border-4",
  lg: "h-12 w-12 border-4",
};

/**
 * Reusable loading spinner component with consistent styling
 */
export function LoadingSpinner({ size = "md", className }: LoadingSpinnerProps) {
  return (
    <output
      className={cn(
        "animate-spin rounded-full border-primary border-t-transparent",
        sizeClasses[size],
        className
      )}
      aria-label="Loading"
    />
  );
}
