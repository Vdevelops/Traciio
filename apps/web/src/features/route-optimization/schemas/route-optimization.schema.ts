import { z } from "zod";

export const locationSchema = z.object({
  lat: z.number().min(-90).max(90),
  lng: z.number().min(-180).max(180),
  address: z.string().optional(),
});

export const waypointSchema = z.object({
  lat: z.number().min(-90).max(90),
  lng: z.number().min(-180).max(180),
  address: z.string().optional(),
  account_id: z.string().uuid().optional(),
  account_name: z.string().optional(),
  contact_id: z.string().uuid().optional(),
  contact_name: z.string().optional(),
  visit_report_id: z.string().uuid().optional(),
});

export const optimizeRouteSchema = z.object({
  route_name: z.string().max(255).optional(),
  start_location: locationSchema,
  waypoints: z.array(waypointSchema).min(1, "Minimum 1 destination required").max(25, "Maximum 25 destinations allowed"),
});

export type OptimizeRouteFormData = z.infer<typeof optimizeRouteSchema>;
export type WaypointFormData = z.infer<typeof waypointSchema>;
export type LocationFormData = z.infer<typeof locationSchema>;
