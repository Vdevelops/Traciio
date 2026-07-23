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

const WALL_CLOCK_PATTERN = /^(\d{4})-(\d{2})-(\d{2})(?:[T ](\d{2}):(\d{2})(?::\d{2}(?:\.\d+)?)?)?/;

export function parseWallClockDateTime(value: string | undefined | null): Date | null {
  if (!value) return null;

  const match = value.match(WALL_CLOCK_PATTERN);
  if (match) {
    const [, year, month, day, hour = "00", minute = "00"] = match;
    const date = new Date(Number(year), Number(month) - 1, Number(day), Number(hour), Number(minute));
    if (!Number.isNaN(date.getTime())) {
      return date;
    }
  }

  const fallback = new Date(value);
  return Number.isNaN(fallback.getTime()) ? null : fallback;
}

export function formatWallClockDateTime(value: string | undefined | null, locale = "id-ID"): string {
  const date = parseWallClockDateTime(value);
  if (!date) return "-";

  return date.toLocaleString(locale, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function formatDateTimeWithLocalOffset(date: Date): string {
  const offsetMinutes = -date.getTimezoneOffset();
  const sign = offsetMinutes >= 0 ? "+" : "-";
  const absOffset = Math.abs(offsetMinutes);
  const offsetHours = String(Math.floor(absOffset / 60)).padStart(2, "0");
  const offsetMins = String(absOffset % 60).padStart(2, "0");
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  const seconds = String(date.getSeconds()).padStart(2, "0");

  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}${sign}${offsetHours}:${offsetMins}`;
}

export function parseDateTimeInputToLocalOffset(value: string): string {
  const [datePart, timePart] = value.split("T");
  if (!datePart || !timePart) {
    return formatDateTimeWithLocalOffset(new Date());
  }

  const [year, month, day] = datePart.split("-").map(Number);
  const [hours, minutes] = timePart.split(":").map(Number);
  const date = new Date(year, month - 1, day, hours, minutes, 0, 0);

  if (Number.isNaN(date.getTime())) {
    return formatDateTimeWithLocalOffset(new Date());
  }

  return formatDateTimeWithLocalOffset(date);
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
