/**
 * Type declarations for @mapbox/polyline module
 * This allows importing @mapbox/polyline without TypeScript errors
 */
declare module "@mapbox/polyline" {
  interface Polyline {
    /**
     * Decode a polyline string into an array of coordinates
     * @param encoded - Encoded polyline string
     * @param precision - Precision level (default: 5)
     * @returns Array of [latitude, longitude] coordinate pairs
     */
    decode(encoded: string, precision?: number): [number, number][];

    /**
     * Encode an array of coordinates into a polyline string
     * @param coordinates - Array of [latitude, longitude] coordinate pairs
     * @param precision - Precision level (default: 5)
     * @returns Encoded polyline string
     */
    encode(coordinates: [number, number][], precision?: number): string;
  }

  const polyline: Polyline;
  export default polyline;
}

