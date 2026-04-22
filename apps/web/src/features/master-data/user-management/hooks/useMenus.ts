"use client";

import { useQuery } from "@tanstack/react-query";
import { menuService } from "../services/userService";

// Helper to check if user has token (client-side only)
function hasToken(): boolean {
  if (typeof window === "undefined") return false;
  return !!localStorage.getItem("token");
}

export function useMenus() {
  return useQuery({
    queryKey: ["menus"],
    queryFn: () => menuService.list(),
    enabled: hasToken(), // Only fetch when authenticated
    retry: (failureCount, error) => {
      if (error && typeof error === "object" && "response" in error) {
        const axiosError = error as { response?: { status?: number } };
        // Don't retry on 401 (unauthorized)
        if (axiosError.response?.status === 401) {
          return false;
        }
      }
      return failureCount < 1;
    },
  });
}

export function useMenu(id: string) {
  return useQuery({
    queryKey: ["menu", id],
    queryFn: () => menuService.getById(id),
    enabled: !!id && hasToken(), // Only fetch when authenticated and id is provided
    retry: (failureCount, error) => {
      if (error && typeof error === "object" && "response" in error) {
        const axiosError = error as { response?: { status?: number } };
        // Don't retry on 401 (unauthorized)
        if (axiosError.response?.status === 401) {
          return false;
        }
      }
      return failureCount < 1;
    },
  });
}
