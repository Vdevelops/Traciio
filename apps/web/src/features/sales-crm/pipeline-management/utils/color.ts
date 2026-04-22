/**
 * Get color based on probability percentage (0-100)
 * Returns a color ranging from red (0%) to green (100%) using HSL
 */
export function getProbabilityColor(probability: number = 0): string {
  // Clamp probability between 0 and 100
  const p = Math.max(0, Math.min(100, probability));
  
  // HSL Hue: 0 is Red, 120 is Green
  // We want 0% -> 0 (Red), 100% -> 120 (Green)
  const hue = (p / 100) * 120;
  
  // Return HSL string with fixed Saturation and Lightness
  // Adjust Saturation and Lightness for better readability on both themes
  return `hsl(${hue}, 80%, 45%)`;
}
