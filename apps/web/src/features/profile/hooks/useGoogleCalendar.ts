"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { googleCalendarService } from "../services/googleCalendarService";
import { toast } from "sonner";

/**
 * Hook to get Google Calendar connection status
 */
export function useGoogleCalendarStatus() {
  return useQuery({
    queryKey: ["google-calendar", "status"],
    queryFn: () => googleCalendarService.getStatus(),
    retry: false,
  });
}

/**
 * Hook to get Google Calendar OAuth2 authorization URL
 */
export function useGoogleCalendarAuthURL() {
  return useQuery({
    queryKey: ["google-calendar", "auth-url"],
    queryFn: () => googleCalendarService.getAuthURL(),
    enabled: false, // Only fetch when explicitly called
    retry: false,
  });
}

/**
 * Hook to initiate Google Calendar OAuth2 flow
 * Opens the authorization URL in a new window
 */
export function useConnectGoogleCalendar() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      const response = await googleCalendarService.getAuthURL();
      if (!response.success || !response.data?.auth_url) {
        throw new Error(response.error || "Failed to get authorization URL");
      }
      return response.data;
    },
    onSuccess: (data) => {
      const width = 500;
      const height = 600;
      const left = window.screen.width / 2 - width / 2;
      const top = window.screen.height / 2 - height / 2;

      const authWindow = window.open(
        data.auth_url,
        "Google Calendar Authorization",
        `width=${width},height=${height},left=${left},top=${top},resizable=yes,scrollbars=yes`
      );

      if (!authWindow) {
        toast.error("Failed to open authorization window. Please check your popup blocker settings.");
        return;
      }

      let pollCount = 0;
      const maxPolls = 60;

      const messageHandler = (event: MessageEvent) => {
        if (event.origin !== window.location.origin) {
          return;
        }

        if (event.data?.type === "GOOGLE_CALENDAR_CALLBACK_SUCCESS") {
          clearInterval(checkStatus);
          window.removeEventListener("message", messageHandler);
          queryClient.invalidateQueries({ queryKey: ["google-calendar", "status"] });
          queryClient.invalidateQueries({ queryKey: ["profile"] });
          toast.success("Google Calendar connected successfully!");

          setTimeout(() => {
            if (!authWindow.closed) {
              authWindow.close();
            }
          }, 500);
        }
      };

      window.addEventListener("message", messageHandler);

      const checkStatus = setInterval(() => {
        pollCount++;

        if (authWindow.closed) {
          clearInterval(checkStatus);
          window.removeEventListener("message", messageHandler);
          setTimeout(() => {
            queryClient.invalidateQueries({ queryKey: ["google-calendar", "status"] });
            queryClient.invalidateQueries({ queryKey: ["profile"] });
          }, 2000);
          return;
        }

        if (pollCount >= maxPolls) {
          clearInterval(checkStatus);
          window.removeEventListener("message", messageHandler);
          queryClient.invalidateQueries({ queryKey: ["google-calendar", "status"] });
          queryClient.invalidateQueries({ queryKey: ["profile"] });
        }
      }, 2000);

      toast.success("Opening Google Calendar authorization...");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to connect Google Calendar");
    },
  });
}

/**
 * Hook to disconnect Google Calendar
 */
export function useDisconnectGoogleCalendar() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      const response = await googleCalendarService.disconnect();
      if (!response.success) {
        throw new Error(response.error || "Failed to disconnect Google Calendar");
      }
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["profile"] });
      toast.success("Google Calendar disconnected successfully");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to disconnect Google Calendar");
    },
  });
}

