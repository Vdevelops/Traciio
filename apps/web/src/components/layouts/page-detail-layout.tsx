"use client";

import React from "react";
import { Button } from "@/components/ui/button";
import { ChevronLeft } from "lucide-react";
import { useRouter } from "@/i18n/routing";

export interface PageDetailLayoutProps {
  title: React.ReactNode;
  subtitle?: React.ReactNode;
  actions?: React.ReactNode;
  children: React.ReactNode;
  onBack?: () => void;
  backHref?: string;
  className?: string;
}

export function PageDetailLayout({
  title,
  subtitle,
  actions,
  children,
  onBack,
  backHref,
  className = "",
}: PageDetailLayoutProps) {
  const router = useRouter();

  const handleBack = () => {
    if (onBack) {
      onBack();
    } else if (backHref) {
      router.push(backHref);
    } else {
      router.back();
    }
  };

  return (
    <div className={`space-y-6 ${className}`}>
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div className="flex items-start gap-4">
          <Button
            variant="outline"
            size="icon"
            onClick={handleBack}
            className="shrink-0 mt-1 sm:mt-0"
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-2xl sm:text-3xl font-medium tracking-tight">
              {title}
            </h1>
            {subtitle && (
              <div className="text-muted-foreground mt-1 text-sm sm:text-base">
                {subtitle}
              </div>
            )}
          </div>
        </div>
        {actions && (
          <div className="flex items-center gap-2 shrink-0">{actions}</div>
        )}
      </div>
      <div className="w-full">
        {children}
      </div>
    </div>
  );
}
