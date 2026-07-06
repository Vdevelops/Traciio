export interface Waypoint {
  order?: number;
  lat: number;
  lng: number;
  address?: string;
  account_id?: string;
  account_name?: string;
  contact_id?: string;
  contact_name?: string;
  visit_report_id?: string;
  account?: {
    id: string;
    name: string;
  };
}

export interface Location {
  lat: number;
  lng: number;
  address?: string;
  accuracy?: number;
}

export interface RouteStep {
  step: number;
  distance: number;
  distance_formatted?: string;
  duration: number;
  duration_formatted?: string;
  instruction?: string;
  polyline?: string;
  maneuver?: string;
  start_location: Location;
  end_location: Location;
}

export interface OptimizedRoute {
  id: string;
  route_name?: string;
  user_id: string;
  waypoints: Waypoint[];
  optimized_order: number[];
  total_distance?: number;
  total_distance_formatted?: string;
  total_duration?: number;
  total_duration_formatted?: string;
  route_polyline?: string;
  route_steps?: RouteStep[];
  created_at: string;
  updated_at: string;
  user?: {
    id: string;
    name: string;
  };
}

export interface OptimizeRouteRequest {
  route_name?: string;
  start_location: Location;
  waypoints: Waypoint[];
}

export interface CalculateDistanceRequest {
  origin: Location;
  destination: Location;
}

export interface CalculateDistanceResponseData {
  distance: number;
  distance_formatted: string;
  duration: number;
  duration_formatted: string;
}

export interface ListRoutesResponse {
  success: boolean;
  data: OptimizedRoute[];
  meta: {
    pagination: {
      page: number;
      per_page: number;
      total: number;
      total_pages: number;
      has_next: boolean;
      has_prev: boolean;
      next_page?: number;
      prev_page?: number;
    };
    filters?: Record<string, unknown>;
  };
  timestamp: string;
  request_id: string;
}

export interface OptimizedRouteResponse {
  success: boolean;
  data: OptimizedRoute;
  meta?: {
    optimization_type?: string;
    waypoints_count?: number;
  };
  timestamp: string;
  request_id: string;
}

export interface CalculateDistanceResponseWrapper {
  success: boolean;
  data: CalculateDistanceResponseData;
  meta?: Record<string, unknown>;
  timestamp: string;
  request_id: string;
}
