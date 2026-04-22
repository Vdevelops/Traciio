import { apiClient } from "@/lib/api-client";

export interface GoogleCalendarAuthURLResponse {
  success: boolean;
  data: {
    auth_url: string;
    state: string;
  };
  error?: string;
}

export interface GoogleCalendarCallbackResponse {
  success: boolean;
  data: {
    message: string;
  };
  error?: string;
}

export interface GoogleCalendarDisconnectResponse {
  success: boolean;
  data: {
    message: string;
  };
  error?: string;
}

export interface GoogleCalendarStatusResponse {
  success: boolean;
  data: {
    connected: boolean;
  };
  error?: string;
}

export const googleCalendarService = {
  /**
   * Get Google Calendar connection status
   */
  getStatus: async (): Promise<GoogleCalendarStatusResponse> => {
    try {
      const response = await apiClient.get<GoogleCalendarStatusResponse>(
        "/google-calendar/status"
      );
      return response.data;
    } catch (error: any) {
      throw error;
    }
  },
  /**
   * Get Google Calendar OAuth2 authorization URL
   */
  getAuthURL: async (): Promise<GoogleCalendarAuthURLResponse> => {
    try {
      const response = await apiClient.get<GoogleCalendarAuthURLResponse>(
        "/google-calendar/auth-url"
      );
      return response.data;
    } catch (error: any) {
      throw error;
    }
  },

  /**
   * Handle Google Calendar OAuth2 callback
   * Note: This is typically handled by redirect, but can be called manually if needed
   */
  handleCallback: async (code: string, state: string): Promise<GoogleCalendarCallbackResponse> => {
    const response = await apiClient.get<GoogleCalendarCallbackResponse>(
      `/google-calendar/callback?code=${code}&state=${state}`
    );
    return response.data;
  },

  /**
   * Disconnect Google Calendar integration
   */
  disconnect: async (): Promise<GoogleCalendarDisconnectResponse> => {
    const response = await apiClient.delete<GoogleCalendarDisconnectResponse>(
      "/google-calendar/disconnect"
    );
    return response.data;
  },
};

