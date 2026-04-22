/**
 * Utility functions for converting between different geometry formats
 */

import type { GeoJSONPolygon } from "../types";

/**
 * Convert GeoJSON Polygon to WKT (Well-Known Text) format
 * WKT format: POLYGON((lng lat, lng lat, lng lat, lng lat))
 */
export function geoJSONToWKT(geoJSON: GeoJSONPolygon): string {
  if (geoJSON.type !== "Polygon") {
    throw new Error("Only Polygon type is supported");
  }

  const coordinates = geoJSON.coordinates[0]; // Get exterior ring
  if (!coordinates || coordinates.length < 4) {
    throw new Error("Polygon must have at least 4 points (including closing point)");
  }

  // Convert [lng, lat] pairs to "lng lat" strings
  const wktCoords = coordinates.map(coord => `${coord[0]} ${coord[1]}`).join(", ");
  
  return `POLYGON((${wktCoords}))`;
}

/**
 * Convert WKT Polygon to GeoJSON format
 * Example WKT: POLYGON((lng lat, lng lat, lng lat, lng lat))
 */
export function wktToGeoJSON(wkt: string): GeoJSONPolygon {
  // Simple regex to extract coordinates from WKT POLYGON
  const match = wkt.match(/POLYGON\s*\(\s*\(\s*([^)]+)\s*\)\s*\)/i);
  if (!match) {
    throw new Error("Invalid WKT Polygon format");
  }

  const coordString = match[1];
  const coordinates = coordString
    .split(",")
    .map(pair => {
      const [lng, lat] = pair.trim().split(/\s+/).map(Number);
      return [lng, lat] as [number, number];
    });

  return {
    type: "Polygon",
    coordinates: [coordinates],
  };
}

/**
 * Ensure polygon is properly closed (first point equals last point)
 */
export function ensurePolygonClosed(coordinates: [number, number][]): [number, number][] {
  if (coordinates.length < 3) {
    throw new Error("Polygon must have at least 3 points");
  }

  const first = coordinates[0];
  const last = coordinates[coordinates.length - 1];

  // Check if already closed
  if (first[0] === last[0] && first[1] === last[1]) {
    return coordinates;
  }

  // Close the polygon by adding first point at the end
  return [...coordinates, first];
}

/**
 * Validate GeoJSON Polygon format
 */
export function validateGeoJSONPolygon(geoJSON: GeoJSONPolygon): boolean {
  if (geoJSON.type !== "Polygon") {
    return false;
  }

  if (!geoJSON.coordinates || geoJSON.coordinates.length === 0) {
    return false;
  }

  const exteriorRing = geoJSON.coordinates[0];
  if (!exteriorRing || exteriorRing.length < 4) {
    return false;
  }

  // Check if polygon is closed
  const first = exteriorRing[0];
  const last = exteriorRing[exteriorRing.length - 1];
  return first[0] === last[0] && first[1] === last[1];
}

/**
 * Calculate polygon area (approximate, for validation)
 */
export function calculatePolygonArea(coordinates: [number, number][]): number {
  if (coordinates.length < 3) return 0;

  let area = 0;
  for (let i = 0; i < coordinates.length - 1; i++) {
    const [x1, y1] = coordinates[i];
    const [x2, y2] = coordinates[i + 1];
    area += (x1 * y2) - (x2 * y1);
  }
  return Math.abs(area) / 2;
}