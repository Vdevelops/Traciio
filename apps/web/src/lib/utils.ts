import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Format currency (from sen to Rupiah)
 * @param amount Amount in smallest currency unit (sen)
 * @returns Formatted currency string (e.g., "Rp 1.000.000")
 */
export function formatCurrency(amount: number | undefined | null): string {
  // Handle undefined/null/NaN
  if (amount === undefined || amount === null || isNaN(amount)) {
    return "Rp 0";
  }
  
  // Convert from sen to Rupiah
  const rupiah = amount / 100;
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(rupiah);
}

/**
 * Format date string to readable format
 * @param dateString ISO date string
 * @returns Formatted date string (e.g., "Jan 15, 2024 10:30 AM")
 */
export function formatDate(dateString: string | undefined | null): string {
  if (!dateString) return "-";
  
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return "-";
    
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date);
  } catch {
    return "-";
  }
}

/**
 * Format phone number to clean WhatsApp format (https://wa.me/...)
 * @param phone Phone number string
 * @returns Format link WhatsApp (e.g. "https://wa.me/628123456789")
 */
export function formatPhoneNumberToWA(phone: string | null | undefined): string {
  if (!phone) return '#';
  let cleaned = phone.replace(/[^\d+]/g, '');
  if (cleaned.startsWith('0')) {
    cleaned = '+62' + cleaned.substring(1);
  }
  return `https://wa.me/${cleaned.replace('+', '')}`;
}

/**
 * Format email to mailto format (mailto:...)
 * @param email Email string
 * @returns mailto link
 */
export function formatEmailToMailto(email: string | null | undefined): string {
  if (!email) return '#';
  return `mailto:${email}`;
}
