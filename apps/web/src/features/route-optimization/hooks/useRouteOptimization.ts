"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { routeOptimizationService } from "../services/routeOptimizationService";
import type {
  OptimizeRouteRequest,
  OptimizedRoute,
  CalculateDistanceRequest,
  ListRoutesResponse,
  Location,
} from "../types";

const ACCEPTABLE_LOCATION_ACCURACY_METERS = 50;
const LOCATION_CAPTURE_TIMEOUT_MS = 15000;
const LOCATION_MAXIMUM_AGE_MS = 30000;
const GEOLOCATION_PERMISSION_DENIED_CODE = 1;
const GEOLOCATION_POSITION_UNAVAILABLE_CODE = 2;
const GEOLOCATION_TIMEOUT_CODE = 3;
const REVERSE_GEOCODING_URL =
  process.env.NEXT_PUBLIC_REVERSE_GEOCODING_URL ||
  "https://nominatim.openstreetmap.org/reverse";

export function useOptimizeRoute() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: OptimizeRouteRequest) =>
      routeOptimizationService.optimize(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["route-optimization", "list"] });
    },
  });
}

export function useRouteList(params?: {
  page?: number;
  per_page?: number;
  user_id?: string;
}) {
  return useQuery({
    queryKey: ["route-optimization", "list", params],
    queryFn: () => routeOptimizationService.list(params),
  });
}

export function useRoute(id: string) {
  return useQuery({
    queryKey: ["route-optimization", id],
    queryFn: () => routeOptimizationService.getById(id),
    enabled: !!id,
  });
}

export function useCalculateDistance() {
  return useMutation({
    mutationFn: (data: CalculateDistanceRequest) =>
      routeOptimizationService.calculateDistance(data),
  });
}

export function useDeleteRoute() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => routeOptimizationService.delete(id),
    onMutate: async (id: string) => {
      await queryClient.cancelQueries({ queryKey: ["route-optimization", "list"] });

      const previousLists = queryClient.getQueriesData<ListRoutesResponse>({
        queryKey: ["route-optimization", "list"],
      });

      queryClient.setQueriesData<ListRoutesResponse>(
        { queryKey: ["route-optimization", "list"] },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: (old.data ?? []).filter((route) => route.id !== id),
            meta: {
              ...old.meta,
              pagination: {
                ...old.meta.pagination,
                total: Math.max(0, (old.meta.pagination?.total ?? 0) - 1),
              },
            },
          };
        }
      );

      // Remove any cached single-route entry too
      queryClient.removeQueries({ queryKey: ["route-optimization", id] });

      return { previousLists };
    },
    onError: (_err, _id, context) => {
      // Rollback list caches
      context?.previousLists?.forEach(([key, data]) => {
        queryClient.setQueryData(key, data);
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["route-optimization", "list"], exact: false });
    },
  });
}

// Hook to get user's current location
export function useCurrentLocation() {
  return useMutation({
    mutationFn: async (): Promise<Location> => {
      return new Promise((resolve, reject) => {
        if (!navigator.geolocation) {
          reject(new Error("Geolocation is not supported by your browser"));
          return;
        }

        let bestPosition: GeolocationPosition | null = null;
        let settled = false;
        let watchId: number | null = null;
        let timeoutId: number | null = null;

        const cleanup = () => {
          if (watchId != null) {
            navigator.geolocation.clearWatch(watchId);
          }
          if (timeoutId != null) {
            window.clearTimeout(timeoutId);
          }
        };

        const resolvePosition = async (position: GeolocationPosition) => {
          if (settled) return;
          settled = true;
          cleanup();

          const lat = position.coords.latitude;
          const lng = position.coords.longitude;
          const accuracy = position.coords.accuracy;

          try {
            const reverseUrl = new URL(REVERSE_GEOCODING_URL);
            reverseUrl.searchParams.set("format", "json");
            reverseUrl.searchParams.set("lat", String(lat));
            reverseUrl.searchParams.set("lon", String(lng));
            reverseUrl.searchParams.set("zoom", "18");
            reverseUrl.searchParams.set("addressdetails", "1");

            const response = await fetch(reverseUrl.toString());
            const data = await response.json();

            resolve({
              lat,
              lng,
              accuracy,
              address: data.display_name || "Current Location",
            });
          } catch {
            resolve({
              lat,
              lng,
              accuracy,
              address: "Current Location",
            });
          }
        };

        const rejectWithError = (error: Pick<GeolocationPositionError, "code">) => {
          if (settled) return;
          if (bestPosition) {
            void resolvePosition(bestPosition);
            return;
          }

          settled = true;
          cleanup();

          let errorMessage = "Unable to retrieve your location";
          switch (error.code) {
            case GEOLOCATION_PERMISSION_DENIED_CODE:
              errorMessage = "Location permission denied. Please enable location access.";
              break;
            case GEOLOCATION_POSITION_UNAVAILABLE_CODE:
              errorMessage = "Location information unavailable.";
              break;
            case GEOLOCATION_TIMEOUT_CODE:
              errorMessage = "Location request timed out.";
              break;
          }
          reject(new Error(errorMessage));
        };

        watchId = navigator.geolocation.watchPosition(
          (position) => {
            if (
              !bestPosition ||
              position.coords.accuracy < bestPosition.coords.accuracy
            ) {
              bestPosition = position;
            }

            if (position.coords.accuracy <= ACCEPTABLE_LOCATION_ACCURACY_METERS) {
              void resolvePosition(position);
            }
          },
          rejectWithError,
          {
            enableHighAccuracy: true,
            timeout: LOCATION_CAPTURE_TIMEOUT_MS,
            maximumAge: LOCATION_MAXIMUM_AGE_MS,
          }
        );

        timeoutId = window.setTimeout(() => {
          if (bestPosition) {
            void resolvePosition(bestPosition);
            return;
          }
          rejectWithError({
            code: GEOLOCATION_TIMEOUT_CODE,
          });
        }, LOCATION_CAPTURE_TIMEOUT_MS);
      });
    },
  });
}
