"use client";

import { useCallback } from "react";
import { useRouter } from "@/i18n/routing";
import { useAuthStore } from "../stores/useAuthStore";
import { fullAuthCleanup } from "../utils/clear-auth-cookies";

export function useLogout() {
  const router = useRouter();
  const { logout } = useAuthStore();

  const handleLogout = useCallback(async () => {
    // Call the store's logout method which resets all state including isSessionVerified
    logout();
    
    // Perform full cleanup (localStorage, cookies, and API call)
    await fullAuthCleanup();
    
    // Redirect to login
    router.push("/login");
  }, [router, logout]);

  return handleLogout;
}



