/**
 * Format currency from smallest unit (sen) to formatted string
 */
export function formatCurrency(amount: number | undefined | null): string {
  // Handle undefined/null/NaN
  if (amount === undefined || amount === null || isNaN(amount)) {
    return "Rp 0";
  }
  
  // Convert from sen to rupiah
  const rupiah = amount / 100;
  
  // Format with thousand separator
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(rupiah);
}

/**
 * Format number with thousand separator
 */
export function formatNumber(num: number): string {
  return new Intl.NumberFormat("id-ID").format(num);
}

