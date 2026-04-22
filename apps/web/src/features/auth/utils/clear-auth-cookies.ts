/**
 * Auth Cookie Cleanup Utilities
 * Provides comprehensive cleanup for both localStorage and cookies
 * to prevent zombie auth state
 */

import { deleteCookie } from "@/lib/cookie";
import { clearCSRFToken } from "@/lib/csrf";
import { authService } from "../services/authService";

/**
 * Clears all auth-related cookies that can be accessed via JavaScript
 * Note: HttpOnly cookies must be cleared by the server via API call
 */
export function clearAuthCookies(): void {
  if (globalThis.window === undefined) return;

  // Clear token cookie (non-HttpOnly)
  deleteCookie("token");
  
  // Clear CSRF token from sessionStorage and cookie
  clearCSRFToken();
  deleteCookie("csrf_token");
}

/**
 * Clears all auth-related data from localStorage
 */
export function clearAuthLocalStorage(): void {
  if (globalThis.window === undefined) return;

  localStorage.removeItem("token");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("auth-storage"); // Zustand persisted storage
}

/**
 * Full authentication cleanup - clears both localStorage and cookies
 * Also calls the logout API to clear HttpOnly cookies server-side
 * 
 * Use this when:
 * - Session verification fails
 * - 401 error after refresh token attempt fails
 * - Manual logout
 */
export async function fullAuthCleanup(): Promise<void> {
  // Clear localStorage first
  clearAuthLocalStorage();
  
  // Clear client-side cookies
  clearAuthCookies();
  
  // Call server to clear HttpOnly cookies
  try {
    await authService.logout();
  } catch {
    // Ignore logout API errors - we're already cleaning up
    // The server might be unavailable or token already invalid
  }
}

/**
 * Synchronous cleanup for use in contexts where async is not possible
 * Note: This does NOT call the logout API, so HttpOnly cookies remain
 * Use fullAuthCleanup() when possible
 */
export function clearAuthSync(): void {
  clearAuthLocalStorage();
  clearAuthCookies();
}
