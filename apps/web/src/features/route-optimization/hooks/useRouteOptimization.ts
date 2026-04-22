"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { routeOptimizationService } from "../services/routeOptimizationService";
import type {
  OptimizeRouteRequest,
  OptimizedRoute,
  CalculateDistanceRequest,
  ListRoutesResponse,
} from "../types";

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
    mutationFn: async (): Promise<{
      lat: number;
      lng: number;
      address: string;
    }> => {
      return new Promise((resolve, reject) => {
        if (!navigator.geolocation) {
          reject(new Error("Geolocation is not supported by your browser"));
          return;
        }

        navigator.geolocation.getCurrentPosition(
          async (position) => {
            const lat = position.coords.latitude;
            const lng = position.coords.longitude;

            // Reverse geocode to get address using Nominatim (OpenStreetMap)
            try {
              const response = await fetch(
                `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}`,
                {
                  headers: {
                    "User-Agent": "CRM-Healthcare-App/1.0",
                  },
                }
              );
              const data = await response.json();

              resolve({
                lat,
                lng,
                address: data.display_name || "Current Location",
              });
            } catch (error) {
              // Fallback without address if geocoding fails
              resolve({
                lat,
                lng,
                address: "Current Location",
              });
            }
          },
          (error) => {
            let errorMessage = "Unable to retrieve your location";
            switch (error.code) {
              case error.PERMISSION_DENIED:
                errorMessage = "Location permission denied. Please enable location access.";
                break;
              case error.POSITION_UNAVAILABLE:
                errorMessage = "Location information unavailable.";
                break;
              case error.TIMEOUT:
                errorMessage = "Location request timed out.";
                break;
            }
            reject(new Error(errorMessage));
          },
          {
            enableHighAccuracy: true,
            timeout: 10000,
            maximumAge: 0,
          }
        );
      });
    },
  });
}

