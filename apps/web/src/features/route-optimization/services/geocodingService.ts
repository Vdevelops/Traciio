import apiClient from "@/lib/api-client";

export interface GeocodeRequest {
  address: string;
}

export interface GeocodeResponse {
  latitude: number;
  longitude: number;
  address: string;
}

export interface GeocodeApiResponse {
  success: boolean;
  data: GeocodeResponse;
}

export const geocodingService = {
  async geocode(address: string): Promise<GeocodeResponse> {
    try {
      const response = await apiClient.post<GeocodeApiResponse>("/geocoding/geocode", {
        address,
      });
      
      // Handle axios response structure: response.data contains the API response
      const apiResponse = response.data;
      
      if (!apiResponse.success || !apiResponse.data) {
        throw new Error("Failed to geocode address: Invalid response format");
      }
      
      return apiResponse.data;
    } catch (error: any) {
      // Handle axios errors
      if (error.response) {
        // Server responded with error status
        const errorData = error.response.data;
        if (errorData?.error?.message) {
          throw new Error(errorData.error.message);
        }
        throw new Error(`Geocoding failed: ${error.response.status} ${error.response.statusText}`);
      } else if (error.request) {
        // Request was made but no response received
        throw new Error("No response from geocoding service. Please check your connection.");
      } else {
        // Something else happened
        throw new Error(`Geocoding error: ${error.message || "Unknown error"}`);
      }
    }
  },
};

