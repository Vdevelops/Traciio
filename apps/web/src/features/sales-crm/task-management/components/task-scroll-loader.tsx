"use client";

import { useEffect, useRef } from "react";

interface TaskScrollLoaderProps {
  readonly onLoadMore: () => void;
  readonly hasMore: boolean;
  readonly isLoading: boolean;
}

export function TaskScrollLoader({ 
  onLoadMore, 
  hasMore, 
  isLoading 
}: TaskScrollLoaderProps) {
  const ref = useRef<HTMLDivElement>(null);
  
  useEffect(() => {
    if (!hasMore || isLoading) {
      return;
    }
    
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting) {
          onLoadMore();
        }
      },
      { 
        threshold: 0.1,
        rootMargin: "100px"
      }
    );
    
    if (ref.current) {
      observer.observe(ref.current);
    }
    
    return () => {
      observer.disconnect();
    };
  }, [hasMore, isLoading, onLoadMore]);
  
  if (!hasMore) return null;
  
  return (
    <div ref={ref} className="py-2 flex items-center justify-center">
      {isLoading && (
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <div className="h-3 w-3 animate-spin rounded-full border-2 border-muted-foreground border-t-transparent" />
          <span>Loading more...</span>
        </div>
      )}
    </div>
  );
}
