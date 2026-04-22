"use client";

import { useEffect, useState } from "react";
import { useRateLimitStore } from "../stores/useRateLimitStore";

/**
 * Hook to get real-time countdown text for rate limit
 * Updates every second to show current countdown
 * 
 * Use this in components that need to display countdown text
 */
export function useRateLimitCountdownText(): string {
  const resetTime = useRateLimitStore((state) => state.resetTime);
  const getCountdownText = useRateLimitStore((state) => state.getCountdownText);
  const [countdownText, setCountdownText] = useState<string>("a moment");

  useEffect(() => {
    if (!resetTime) {
      setCountdownText("a moment");
      return;
    }

    // Update immediately
    setCountdownText(getCountdownText());

    // Update every second
    const interval = setInterval(() => {
      setCountdownText(getCountdownText());
    }, 1000);

    return () => {
      clearInterval(interval);
    };
  }, [resetTime, getCountdownText]);

  return countdownText;
}


