/**
 * CSRF Token Management
 * Provides utilities for CSRF protection in SPA applications
 */

const CSRF_TOKEN_KEY = "csrf_token";
const CSRF_COOKIE_NAME = "csrf_token";

/**
 * Gets CSRF token from cookie
 * @returns CSRF token or null if not found
 */
export function getCSRFToken(): string | null {
  if (typeof document === "undefined") return null;

  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${CSRF_COOKIE_NAME}=`);

  if (parts.length === 2) {
    return parts.pop()?.split(";").shift() || null;
  }

  return null;
}

/**
 * Sets CSRF token in memory (received from server)
 * @param token - The CSRF token
 */
export function setCSRFToken(token: string): void {
  if (typeof window === "undefined") return;
  
  // Store in sessionStorage for retrieval
  sessionStorage.setItem(CSRF_TOKEN_KEY, token);
}

/**
 * Gets CSRF token for request header
 * @returns CSRF token for X-CSRF-Token header
 */
export function getCSRFTokenForHeader(): string | null {
  if (typeof window === "undefined") return null;

  // Try to get from sessionStorage first
  const storedToken = sessionStorage.getItem(CSRF_TOKEN_KEY);
  if (storedToken) return storedToken;

  // Fallback to cookie
  return getCSRFToken();
}

/**
 * Clears CSRF token (on logout)
 */
export function clearCSRFToken(): void {
  if (typeof window === "undefined") return;

  sessionStorage.removeItem(CSRF_TOKEN_KEY);
}

/**
 * Checks if CSRF token exists
 * @returns true if CSRF token exists
 */
export function hasCSRFToken(): boolean {
  return getCSRFTokenForHeader() !== null;
}
