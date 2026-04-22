"use client";

import { useEffect, useRef } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import { useRefreshSession } from "@/features/auth/hooks/useRefreshSession";
import type { WebSocketMessage, Notification } from "../types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
const WS_URL = API_BASE_URL.replace(/^http/, "ws");

// Helper function to check if JWT token is expired
function isTokenExpired(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    const exp = payload.exp * 1000; // Convert to milliseconds
    const now = Date.now();
    // Consider token expired if it expires within 1 minute (buffer for clock skew)
    return exp < now + 60000;
  } catch {
    return true; // If we can't parse, consider it expired
  }
}

export function useWebSocket() {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const isConnectingRef = useRef(false);
  const isRefreshingRef = useRef(false);
  const queryClient = useQueryClient();
  const queryClientRef = useRef(queryClient);
  const { token } = useAuthStore();
  const tokenRef = useRef<string | null>(token);
  const { refreshSession } = useRefreshSession();

  // Update refs when values change
  useEffect(() => {
    tokenRef.current = token;
    queryClientRef.current = queryClient;
  }, [token, queryClient]);

  useEffect(() => {
    // Don't connect if no token
    if (!token) {
      // Close connection if token is removed
      if (wsRef.current && wsRef.current.readyState !== WebSocket.CLOSED) {
        try {
          wsRef.current.close(1000, "Token removed");
        } catch {
          // Ignore errors
        }
        wsRef.current = null;
      }
      return;
    }

    // Check if token is expired and try to refresh
    if (isTokenExpired(token)) {
      if (!isRefreshingRef.current) {
        isRefreshingRef.current = true;
        refreshSession()
          .then(() => {
            // Token refreshed, effect will re-run with new token
            isRefreshingRef.current = false;
          })
          .catch((error) => {
            isRefreshingRef.current = false;
            // Close connection if refresh fails
            if (wsRef.current && wsRef.current.readyState !== WebSocket.CLOSED) {
              try {
                wsRef.current.close(1000, "Token refresh failed");
              } catch {
                // Ignore errors
              }
              wsRef.current = null;
            }
          });
      }
      return; // Wait for token refresh
    }

    // Don't connect if already connecting
    if (isConnectingRef.current) {
      return;
    }

    // Don't reconnect if already connected and token is the same
    const currentWs = wsRef.current;
    if (
      currentWs &&
      (currentWs.readyState === WebSocket.OPEN || currentWs.readyState === WebSocket.CONNECTING) &&
      tokenRef.current === token
    ) {
      return;
    }

    // Close existing connection only if token changed
    if (currentWs && tokenRef.current !== token) {
      // Token changed, close old connection
      try {
        currentWs.close(1000, "Token changed");
      } catch {
        // Ignore errors
      }
      wsRef.current = null;
      reconnectAttemptsRef.current = 0;
    } else if (currentWs) {
      // Connection exists but is closed/closing, clear ref
      if (currentWs.readyState === WebSocket.CLOSED || currentWs.readyState === WebSocket.CLOSING) {
        wsRef.current = null;
      } else {
        // Connection exists and is open/connecting with same token, don't create new one
        return;
      }
    }

    isConnectingRef.current = true;
    // Reset reconnect attempts on new connection attempt
    reconnectAttemptsRef.current = 0;

    // Build WebSocket URL
    // Try cookie first (more secure), but fallback to query parameter if needed
    // Note: WebSocket doesn't always send cookies, so we include token in query as fallback
    const wsUrl = tokenRef.current
      ? `${WS_URL}/api/v1/ws/notifications?token=${encodeURIComponent(tokenRef.current)}`
      : `${WS_URL}/api/v1/ws/notifications`;

    // Store message handler for reuse
    const messageHandler = (event: MessageEvent) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);

        switch (message.type) {
          case "notification.created": {
            const notification = message.data as Notification;
            
            // Invalidate queries to refresh data
            queryClientRef.current.invalidateQueries({ queryKey: ["notifications"] });
            queryClientRef.current.invalidateQueries({ queryKey: ["notifications", "unread-count"] });

            // Show toast notification
            toast.info(notification.title, {
              description: notification.message,
              duration: 5000,
            });
            break;
          }

          case "notification.updated":
          case "notification.deleted": {
            // Invalidate queries to refresh data
            queryClientRef.current.invalidateQueries({ queryKey: ["notifications"] });
            queryClientRef.current.invalidateQueries({ queryKey: ["notifications", "unread-count"] });
            break;
          }

          case "permissions.updated": {
            // Emit custom event for PermissionsProvider to handle
            if (typeof window !== "undefined") {
              window.dispatchEvent(new CustomEvent("permissions:updated"));
            }
            
            // Also invalidate roles query for backward compatibility
            queryClientRef.current.invalidateQueries({ queryKey: ["roles"] });
            break;
          }

          default:
        }
      } catch (error) {
      }
    };

    try {
      const ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        isConnectingRef.current = false;
        reconnectAttemptsRef.current = 0;
        
        // Clear any pending reconnect
        if (reconnectTimeoutRef.current) {
          clearTimeout(reconnectTimeoutRef.current);
          reconnectTimeoutRef.current = null;
        }
      };


      ws.onmessage = messageHandler;

      ws.onerror = () => {
        isConnectingRef.current = false;
      };

      ws.onclose = async (event) => {
        isConnectingRef.current = false;
        
        // Clear ref only if this is the current connection
        if (wsRef.current === ws) {
          wsRef.current = null;
        }

        // If connection closed due to authentication error (401), try to refresh token
        // Code 1006 (Abnormal Closure) often indicates server rejected connection
        if ((event.code === 1006 || event.code === 1008) && tokenRef.current && !isRefreshingRef.current) {
          // Check if token might be expired
          if (isTokenExpired(tokenRef.current)) {
            isRefreshingRef.current = true;
            try {
              await refreshSession();
              // Token refreshed, effect will re-run with new token
              isRefreshingRef.current = false;
              return; // Don't attempt reconnect, let effect handle it
            } catch (error) {
              isRefreshingRef.current = false;
            }
          }
        }

        // Only reconnect if not a normal closure (1000) and we have a token
        // Don't reconnect if it was closed intentionally (code 1000)
        if (event.code !== 1000 && tokenRef.current && !isRefreshingRef.current) {
          // Exponential backoff: 1s, 2s, 4s, 8s, max 30s
          const maxAttempts = 5;
          if (reconnectAttemptsRef.current < maxAttempts) {
            reconnectAttemptsRef.current += 1;
            const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current - 1), 30000);
            
            reconnectTimeoutRef.current = setTimeout(() => {
            }, delay);
          } else {
          }
        }
      };

      wsRef.current = ws;
    } catch (error) {
      isConnectingRef.current = false;
    }

    return () => {
      // Cleanup on unmount or token change
      // Only cleanup reconnect timeout - don't close WebSocket here
      // WebSocket will be closed in the next effect run if token changed
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
        reconnectTimeoutRef.current = null;
      }
      // Don't close WebSocket in cleanup - let the effect handle it based on token comparison
      // This prevents closing a valid connection when effect re-runs for other reasons
    };
  }, [token]); // Only depend on token - queryClient is stable and doesn't need to be in deps
}

