import { z } from "zod";

// GeoJSON Schemas
export const geoJSONPointSchema = z.object({
  type: z.literal("Point"),
  coordinates: z.tuple([z.number(), z.number()]), // [lng, lat]
});

export const geoJSONPolygonSchema = z.object({
  type: z.literal("Polygon"),
  coordinates: z.array(z.array(z.tuple([z.number(), z.number()]))),
});

// Territory Schemas
export const createTerritorySchema = z.object({
  name: z.string().min(1, "Territory name is required").max(255),
  description: z.string().optional(),
  polygon: geoJSONPolygonSchema,
  assigned_to: z.string().uuid().optional(),
  color: z.string().regex(/^#[0-9A-F]{6}$/i, "Invalid color format").optional(),
});

export const updateTerritorySchema = z.object({
  name: z.string().min(1).max(255).optional(),
  description: z.string().optional(),
  polygon: geoJSONPolygonSchema.optional(),
  assigned_to: z.string().uuid().optional().nullable(),
  color: z.string().regex(/^#[0-9A-F]{6}$/i).optional(),
});

// Area Capture Schemas
export const createAreaCaptureSchema = z.object({
  visit_report_id: z.string().uuid("Invalid visit report ID"),
  capture_type: z.enum(["check_in", "check_out", "area"]),
  location: geoJSONPointSchema,
  address: z.string().optional(),
  accuracy: z.number().positive().optional(),
});

// Form Schemas (for UI validation)
export const territoryFormSchema = z.object({
  name: z.string().min(1, "Name is required").max(255),
  description: z.string().max(1000).optional(),
  color: z.string().regex(/^#[0-9A-F]{6}$/i, "Invalid color format"),
  assigned_to: z.string().uuid().optional(),
});

export const areaCaptureFormSchema = z.object({
  visit_report_id: z.string().uuid("Select a valid visit report"),
  capture_type: z.enum(["check_in", "check_out", "area"]),
  address: z.string().optional(),
  accuracy: z.number().positive().optional(),
});

// Query Schemas
export const getNearbyCapturesSchema = z.object({
  lat: z.number().min(-90).max(90),
  lon: z.number().min(-180).max(180),
  radius_meters: z.number().positive().max(50000), // Max 50km
  limit: z.number().positive().max(1000).optional(),
});

export const getCapturesInTerritorySchema = z.object({
  territory_id: z.string().uuid(),
  start_date: z.string().datetime().optional(),
  end_date: z.string().datetime().optional(),
});

// Type exports
export type CreateTerritoryInput = z.infer<typeof createTerritorySchema>;
export type UpdateTerritoryInput = z.infer<typeof updateTerritorySchema>;
export type CreateAreaCaptureInput = z.infer<typeof createAreaCaptureSchema>;
export type TerritoryFormInput = z.infer<typeof territoryFormSchema>;
export type AreaCaptureFormInput = z.infer<typeof areaCaptureFormSchema>;
export type GetNearbyCapturesInput = z.infer<typeof getNearbyCapturesSchema>;
export type GetCapturesInTerritoryInput = z.infer<typeof getCapturesInTerritorySchema>;
