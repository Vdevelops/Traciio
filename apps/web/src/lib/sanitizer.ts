/**
 * Sanitizer utility for input validation and XSS prevention
 * This module provides comprehensive sanitization functions to protect against XSS, SQL injection, and other security threats
 */

/**
 * Sanitizes HTML content by removing potentially dangerous elements
 * @param input - The HTML string to sanitize
 * @returns Sanitized HTML string
 */
export function sanitizeHTML(input: string): string {
  if (!input) return "";

  // Create a temporary div to parse HTML
  const temp = document.createElement("div");
  temp.textContent = input;
  
  return temp.innerHTML;
}

/**
 * Sanitizes user input by HTML-escaping it
 * @param input - The input string to sanitize
 * @returns HTML-escaped string
 */
export function sanitizeInput(input: string): string {
  if (!input) return "";

  const map: Record<string, string> = {
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': "&quot;",
    "'": "&#x27;",
    "/": "&#x2F;",
  };

  return input.replace(/[&<>"'/]/g, (char) => map[char]);
}

/**
 * Validates that input doesn't contain script injection attempts
 * @param input - The input string to validate
 * @returns true if input is safe, false if suspicious
 */
export function validateNoScriptInjection(input: string): boolean {
  if (!input) return true;

  const lowered = input.toLowerCase();

  // Check for script tags
  if (lowered.includes("<script")) return false;

  // Check for javascript: protocol
  if (lowered.includes("javascript:")) return false;

  // Check for event handlers
  if (/<[^>]*\son\w+\s*=/i.test(input)) return false;

  // Check for dangerous data URIs
  if (lowered.includes("data:text/html")) return false;

  // Check for eval and similar dangerous functions
  if (/eval\s*\(|Function\s*\(/i.test(input)) return false;

  return true;
}

/**
 * Sanitizes filename by removing dangerous characters
 * @param filename - The filename to sanitize
 * @returns Sanitized filename
 */
export function sanitizeFilename(filename: string): string {
  if (!filename) return "";

  // Replace path separators and dangerous characters
  let sanitized = filename
    .replace(/[/\\]/g, "_")
    .replace(/\x00/g, "")
    .replace(/[<>:"|?*]/g, "_");

  // Remove control characters
  sanitized = sanitized.replace(/[\x00-\x1F\x7F]/g, "");

  return sanitized;
}

/**
 * Sanitizes URL to prevent dangerous protocols
 * @param url - The URL to sanitize
 * @returns Sanitized URL or empty string if dangerous
 */
export function sanitizeURL(url: string): string {
  if (!url) return "";

  const trimmed = url.trim();
  const lowered = trimmed.toLowerCase();

  // Block dangerous protocols
  if (
    lowered.startsWith("javascript:") ||
    lowered.startsWith("data:text/html") ||
    lowered.startsWith("vbscript:") ||
    lowered.startsWith("file:")
  ) {
    return "";
  }

  // Only allow http, https, mailto, tel
  if (
    !lowered.startsWith("http://") &&
    !lowered.startsWith("https://") &&
    !lowered.startsWith("mailto:") &&
    !lowered.startsWith("tel:") &&
    !lowered.startsWith("/") &&
    !lowered.startsWith("#")
  ) {
    return "";
  }

  return trimmed;
}

/**
 * Truncates string to maximum length (prevents buffer overflow)
 * @param input - The input string
 * @param maxLength - Maximum allowed length
 * @returns Truncated string
 */
export function truncateString(input: string, maxLength: number): string {
  if (!input || input.length <= maxLength) return input;
  return input.substring(0, maxLength);
}

/**
 * Validates email format
 * @param email - The email address to validate
 * @returns true if valid email format
 */
export function validateEmail(email: string): boolean {
  if (!email) return false;

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

/**
 * Sanitizes email address
 * @param email - The email address to sanitize
 * @returns Sanitized email
 */
export function sanitizeEmail(email: string): string {
  if (!email) return "";

  // Trim and lowercase
  let sanitized = email.trim().toLowerCase();

  // Remove any HTML tags
  sanitized = sanitizeInput(sanitized);

  return sanitized;
}

/**
 * Detects potential SQL injection patterns
 * @param input - The input string to check
 * @returns true if suspicious pattern detected
 */
export function detectSQLInjection(input: string): boolean {
  if (!input) return false;

  const sqlPatterns = [
    /(\bunion\b.*\bselect\b)/i,
    /(\bselect\b.*\bfrom\b)/i,
    /(\binsert\b.*\binto\b)/i,
    /(\bupdate\b.*\bset\b)/i,
    /(\bdelete\b.*\bfrom\b)/i,
    /(\bdrop\b.*\btable\b)/i,
    /(--|;|\/\*|\*\/)/,
    /(\bexec\b|\bexecute\b)/i,
    /'\s*(or|and)\s*'/i,
  ];

  return sqlPatterns.some((pattern) => pattern.test(input));
}

/**
 * Removes sensitive data from objects for logging
 * @param obj - The object to sanitize
 * @returns Sanitized object
 */
export function sanitizeForLog<T extends Record<string, unknown>>(obj: T): T {
  const sensitiveKeys = [
    "password",
    "token",
    "secret",
    "apiKey",
    "api_key",
    "authorization",
    "refreshToken",
    "refresh_token",
    "accessToken",
    "access_token",
    "csrfToken",
    "csrf_token",
  ];

  const sanitized = { ...obj };

  for (const key in sanitized) {
    if (sensitiveKeys.some((sk) => key.toLowerCase().includes(sk.toLowerCase()))) {
      sanitized[key] = "***REDACTED***" as T[Extract<keyof T, string>];
    }
  }

  return sanitized;
}

/**
 * Validates that string doesn't exceed maximum length
 * @param input - The input string
 * @param maxLength - Maximum allowed length
 * @returns true if valid length
 */
export function validateMaxLength(input: string, maxLength: number): boolean {
  return !input || input.length <= maxLength;
}

/**
 * Sanitizes phone number
 * @param phone - The phone number to sanitize
 * @returns Sanitized phone number
 */
export function sanitizePhone(phone: string): string {
  if (!phone) return "";

  // Remove all non-numeric characters except +
  return phone.replace(/[^\d+]/g, "");
}

/**
 * Validates Indonesian phone number format
 * @param phone - The phone number to validate
 * @returns true if valid Indonesian phone format
 */
export function validateIndonesianPhone(phone: string): boolean {
  if (!phone) return false;

  // Indonesian phone: starts with +62, 62, or 0, followed by 8-12 digits
  const phoneRegex = /^(\+62|62|0)[0-9]{8,12}$/;
  return phoneRegex.test(phone.replace(/[\s-]/g, ""));
}

/**
 * Sanitizes object by removing undefined and null values
 * @param obj - The object to sanitize
 * @returns Sanitized object
 */
export function sanitizeObject<T extends Record<string, unknown>>(obj: T): Partial<T> {
  const sanitized: Partial<T> = {};

  for (const key in obj) {
    if (obj[key] !== undefined && obj[key] !== null) {
      sanitized[key] = obj[key];
    }
  }

  return sanitized;
}
