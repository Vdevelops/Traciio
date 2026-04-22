// Area Mapping Types
export interface Territory {
  id: string;
  name: string;
  description?: string;
  polygon: GeoJSONPolygon;
  assigned_to?: string;
  color: string;
  created_at: string;
  updated_at: string;
}

export interface AreaCapture {
  id: string;
  visit_report_id: string;
  capture_type: "check_in" | "check_out" | "area";
  location: GeoJSONPoint;
  address?: string;
  accuracy?: number;
  captured_at: string;
  created_at: string;
  updated_at: string;
}

export interface CoverageAnalysis {
  id: string;
  territory_id?: string;
  user_id?: string;
  period_start: string;
  period_end: string;
  visit_count: number;
  coverage_percent?: number;
  analyzed_at: string;
  created_at: string;
  updated_at: string;
}

// GeoJSON Types
export interface GeoJSONPoint {
  type: "Point";
  coordinates: [number, number]; // [lng, lat]
}

export interface GeoJSONPolygon {
  type: "Polygon";
  coordinates: [number, number][][]; // Array of rings, each ring is array of [lng, lat]
}

// API Request Types
export interface CreateTerritoryRequest {
  name: string;
  description?: string;
  polygon: GeoJSONPolygon;
  assigned_to?: string;
  color?: string;
}

export interface UpdateTerritoryRequest {
  name?: string;
  description?: string;
  polygon?: GeoJSONPolygon;
  assigned_to?: string;
  color?: string;
}

export interface CreateAreaCaptureRequest {
  visit_report_id: string;
  capture_type: "check_in" | "check_out" | "area";
  latitude: number;
  longitude: number;
  address?: string;
  accuracy?: number;
}

// API Response Types
export interface TerritoriesResponse {
  territories: Territory[];
  total: number;
  page: number;
  page_size: number;
}

export interface AreaCapturesResponse {
  captures: AreaCapture[];
  total: number;
  page: number;
  page_size: number;
}


export interface HeatmapDataPoint {
  lat: number;
  lng: number;
  intensity: number;
}

export interface HeatmapResponse {
  data: HeatmapDataPoint[];
  bounds: {
    north: number;
    south: number;
    east: number;
    west: number;
  };
}

// UI State Types
export interface TerritoryFormData {
  name: string;
  description: string;
  color: string;
  assigned_to?: string;
  polygon?: GeoJSONPolygon;
}

export interface AreaMappingFilters {
  search: string;
  assigned_to?: string;
  date_from?: string;
  date_to?: string;
}
