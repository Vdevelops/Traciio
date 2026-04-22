import apiClient from "@/lib/api-client";
import type {
  VisitReport,
  ListVisitReportsResponse,
  VisitReportResponse,
  CreateVisitReportFormData,
  UpdateVisitReportFormData,
  CheckInFormData,
  CheckOutFormData,
  RejectFormData,
  UploadPhotoFormData,
  SubmitVisitReportFormData,
} from "../types";

export const visitReportService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    search?: string;
    status?: string;
    account_id?: string;
    deal_id?: string;
    lead_id?: string;
    sales_rep_id?: string;
    start_date?: string;
    end_date?: string;
  }): Promise<ListVisitReportsResponse> {
    const response = await apiClient.get<ListVisitReportsResponse>("/visit-reports", { params });
    return response.data;
  },

  async getById(id: string): Promise<VisitReportResponse> {
    const response = await apiClient.get<VisitReportResponse>(`/visit-reports/${id}`);
    return response.data;
  },

  async create(data: CreateVisitReportFormData): Promise<VisitReportResponse> {
    const response = await apiClient.post<VisitReportResponse>("/visit-reports", data);
    return response.data;
  },

  async update(id: string, data: UpdateVisitReportFormData): Promise<VisitReportResponse> {
    const response = await apiClient.put<VisitReportResponse>(`/visit-reports/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/visit-reports/${id}`);
  },

  async checkIn(
    id: string,
    data: CheckInFormData,
    options?: {
      photo?: File;
      deviceGPS?: {
        latitude: number;
        longitude: number;
        accuracy?: number;
        timestamp?: number;
      };
      photoGPS?: {
        latitude: number;
        longitude: number;
        timestamp?: number;
      };
    }
  ): Promise<VisitReportResponse> {
    // If photo is provided, use multipart/form-data
    if (options?.photo) {
      const formData = new FormData();
      formData.append("location[latitude]", data.location.latitude.toString());
      formData.append("location[longitude]", data.location.longitude.toString());
      if (data.location.address) {
        formData.append("location[address]", data.location.address);
      }
      formData.append("photo", options.photo);

      // Add device GPS metadata if provided
      if (options.deviceGPS) {
        formData.append("device_gps[latitude]", options.deviceGPS.latitude.toString());
        formData.append("device_gps[longitude]", options.deviceGPS.longitude.toString());
        if (options.deviceGPS.accuracy !== undefined) {
          formData.append("device_gps[accuracy]", options.deviceGPS.accuracy.toString());
        }
        if (options.deviceGPS.timestamp !== undefined) {
          formData.append("device_gps[timestamp]", options.deviceGPS.timestamp.toString());
        }
      }

      // Add photo GPS metadata if provided
      if (options.photoGPS) {
        formData.append("photo_gps[latitude]", options.photoGPS.latitude.toString());
        formData.append("photo_gps[longitude]", options.photoGPS.longitude.toString());
        if (options.photoGPS.timestamp !== undefined) {
          formData.append("photo_gps[timestamp]", options.photoGPS.timestamp.toString());
        }
      }

      try {
        const response = await apiClient.post<VisitReportResponse>(
          `/visit-reports/${id}/check-in`,
          formData,
          {
            headers: {
              "Content-Type": "multipart/form-data",
            },
            timeout: 30000, // 30 seconds timeout for file upload
          }
        );
        
        // Backend returns { success: true, data: VisitReport, ... }
        // axios wraps it in response.data, so response.data = { success: true, data: VisitReport, ... }
        // We need to return the VisitReportResponse wrapper
        if (response.data && typeof response.data === 'object' && 'success' in response.data) {
          if (response.data.success && response.data.data) {
            return response.data;
          } else {
            throw new Error("Invalid response: success is false or data is missing");
          }
        }
        
        // Fallback: if response.data is already VisitReport (shouldn't happen but handle it)
        return {
          success: true,
          data: response.data as unknown as VisitReport,
          timestamp: new Date().toISOString(),
          request_id: '',
        };
      } catch (error) {
        // Better error logging - handle circular references and complex objects
        if (error && typeof error === "object") {
          // Try to extract useful information from axios error
          const axiosError = error as {
            message?: string;
            code?: string;
            response?: {
              status?: number;
              statusText?: string;
              data?: {
                error?: {
                  code?: string;
                  message?: string;
                };
                message?: string;
              };
            };
            request?: unknown;
          };

          const errorInfo: Record<string, unknown> = {
            message: axiosError.message || "Unknown error",
            code: axiosError.code,
          };

          if (axiosError.response) {
            errorInfo.status = axiosError.response.status;
            errorInfo.statusText = axiosError.response.statusText;
            if (axiosError.response.data) {
              errorInfo.errorCode = axiosError.response.data.error?.code;
              errorInfo.errorMessage = axiosError.response.data.error?.message || axiosError.response.data.message;
            }
          }

        } else {
        }
        throw error;
      }
    }

    // Otherwise use JSON
    const response = await apiClient.post<VisitReportResponse>(`/visit-reports/${id}/check-in`, data);
    return response.data;
  },

  async checkOut(id: string, data: CheckOutFormData): Promise<VisitReportResponse> {
    const response = await apiClient.post<VisitReportResponse>(`/visit-reports/${id}/check-out`, data);
    return response.data;
  },

  async submit(id: string, data: SubmitVisitReportFormData): Promise<VisitReportResponse> {
    const response = await apiClient.patch<VisitReportResponse>(`/visit-reports/${id}/submit`, data);
    return response.data;
  },

  async approve(id: string): Promise<VisitReportResponse> {
    const response = await apiClient.post<VisitReportResponse>(`/visit-reports/${id}/approve`, {});
    return response.data;
  },

  async reject(id: string, data: RejectFormData): Promise<VisitReportResponse> {
    const response = await apiClient.post<VisitReportResponse>(`/visit-reports/${id}/reject`, data);
    return response.data;
  },

  async uploadPhoto(id: string, file: File): Promise<VisitReportResponse> {
    const formData = new FormData();
    formData.append("photo", file);

    const response = await apiClient.post<VisitReportResponse>(
      `/visit-reports/${id}/photos`,
      formData,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      }
    );
    return response.data;
  },
};

