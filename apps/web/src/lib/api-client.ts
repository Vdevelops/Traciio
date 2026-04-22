import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from "axios";
import { toast } from "sonner";
import { setSecureCookie } from "./cookie";
import { formatError } from "./i18n/error-messages";
import { useRateLimitStore } from "./stores/useRateLimitStore";
import { getCSRFTokenForHeader } from "./csrf";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const apiClient: AxiosInstance = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 15000, // 15 seconds for enterprise scale
  withCredentials: true, // Enable cookies for CSRF protection
  // Enterprise: Enable response compression
  decompress: true,
});

// Flag to prevent multiple refresh attempts
let isRefreshing = false;
const failedQueue: Array<{
  resolve: (value?: unknown) => void;
  reject: (error?: unknown) => void;
}> = [];

const processQueue = (error: AxiosError | null, token: string | null = null) => {
  const queue = [...failedQueue];
  failedQueue.splice(0, failedQueue.length); // Clear array
  queue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
};

// Request interceptor untuk menambahkan token dan CSRF
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // Add CSRF token for non-GET requests
    if (config.method && config.method.toLowerCase() !== "get") {
      const csrfToken = getCSRFTokenForHeader();
      if (csrfToken && config.headers) {
        config.headers["X-CSRF-Token"] = csrfToken;
      }
    }

    // Mark this request as first request after app load for rate limit validation
    // This helps validate if rate limit state is still valid after API restart
    if (typeof window !== "undefined") {
      const configWithMeta = config as unknown as Record<string, unknown>;
      if (!configWithMeta._rateLimitValidated) {
        configWithMeta._rateLimitValidated = true;
      }
    }

    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

interface ErrorDetails {
  field?: string;
  resource?: string;
  value?: string;
  [key: string]: unknown;
}

interface ApiErrorResponse {
  success: false;
  error: {
    code: string;
    message: string;
    details?: ErrorDetails;
    field_errors?: Array<{ field: string; message: string }>;
  };
  timestamp: string;
  request_id: string;
}

// Response interceptor untuk handle errors
apiClient.interceptors.response.use(
  (response) => {
    // Clear rate limit reset time on successful response
    // This ensures that if API restarts and rate limit state is reset,
    // frontend will also clear the countdown
    // Only keep reset time if we're actually rate limited (429 error)
    const status = response.status;
    if (status !== 429) {
      // Request succeeded - ALWAYS clear any existing rate limit state
      // This handles the case where API restarts and rate limit is reset
      // If backend rate limit was reset, successful request means we're no longer rate limited
      // This is important because backend uses in-memory rate limiting that resets on restart
      const currentResetTime = useRateLimitStore.getState().resetTime;
      if (currentResetTime) {
        // Clear reset time if we have one
        // This validates that backend rate limit state matches frontend state
        useRateLimitStore.getState().clearResetTime();
      }
    }
    return response;
  },
  async (error: AxiosError<ApiErrorResponse>) => {
    // Network error (tidak terhubung ke server) - lebih jelas dan awam
    if (!error.response) {
      if (error.code === "ECONNABORTED") {
        const msg = formatError("network", "timeout");
        toast.error(msg.title, {
          description: msg.description,
        });
      } else if (error.code === "ERR_NETWORK" || error.message === "Network Error") {
        const msg = formatError("network", "connectionFailed");
        toast.error(msg.title, {
          description: msg.description,
        });
      } else {
        const msg = formatError("network", "generic");
        toast.error(msg.title, {
          description: msg.description,
        });
      }
      return Promise.reject(error);
    }

    // HTTP error responses
    const status = error.response.status;
    const errorData = error.response.data;

    if (!errorData || !errorData.error) {
      const msg = formatError("backend", "invalidFormat");
      toast.error(msg.title, {
        description: msg.description,
      });
      return Promise.reject(error);
    }

    const errorCode = errorData.error.code;
    const errorDetails = errorData.error.details;
    const fieldErrors = errorData.error.field_errors;

    // Handle specific error codes
    if (errorCode === "RESOURCE_ALREADY_EXISTS" || errorCode === "CONFLICT") {
      // Handle duplicate email or other resource conflicts - lebih jelas dan awam
      if (errorDetails?.field === "email" && errorDetails?.resource === "user") {
        const msg = formatError("backend", "emailExists", {
          email: String(errorDetails.value || ""),
        });
        toast.error(msg.title, {
          description: msg.description,
        });
      } else if (errorDetails?.field && errorDetails?.resource) {
        const msg = formatError("backend", "resourceExists", {
          field: errorDetails.field,
          value: String(errorDetails.value || ""),
        });
        toast.error(msg.title, {
          description: msg.description,
        });
      } else {
        const msg = formatError("backend", "conflict");
        toast.error(msg.title, {
          description: msg.description,
        });
      }
      return Promise.reject(error);
    }

    // Handle INTERNAL_SERVER_ERROR with details (e.g., duplicate email from database constraint)
    if (errorCode === "INTERNAL_SERVER_ERROR" && errorDetails) {
      if (errorDetails.field === "email" && errorDetails.resource === "user") {
        const msg = formatError("backend", "emailExists", {
          email: String(errorDetails.value || ""),
        });
        toast.error(msg.title, {
          description: msg.description,
        });
        return Promise.reject(error);
      }
      // Other internal errors with details
      if (errorDetails.field && errorDetails.resource) {
        const msg = formatError("backend", "resourceExists", {
          field: errorDetails.field,
          value: String(errorDetails.value || ""),
        });
        toast.error(msg.title, {
          description: msg.description,
        });
        return Promise.reject(error);
      }
    }

    // Handle validation errors
    if (errorCode === "VALIDATION_ERROR" && fieldErrors && fieldErrors.length > 0) {
      const firstError = fieldErrors[0];
      const msg = formatError("backend", "fieldError", {
        field: firstError.field,
        message: firstError.message,
      });
      toast.error(msg.title, {
        description: msg.description,
      });
      return Promise.reject(error);
    }

    // Handle HTTP status codes
    if (status === 401) {
      const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };
      const requestUrl = originalRequest?.url || "";

      // Skip token refresh for authentication endpoints (login, refresh)
      // These endpoints return 401 for invalid credentials, not expired session
      if (requestUrl.includes("/auth/login") || requestUrl.includes("/auth/refresh")) {
        // For login/refresh endpoints, 401 means invalid credentials
        // Don't try to refresh token, don't show toast, just reject the error
        // Error message will be handled by the login form component
        return Promise.reject(error);
      }

      // Skip refresh if this is already a retry
      if (originalRequest?._retry) {
        // Refresh failed or already retried, full cleanup and logout
        isRefreshing = false;
        if (typeof window !== "undefined") {
          // Import and call full auth cleanup
          import("@/features/auth/utils/clear-auth-cookies").then(({ fullAuthCleanup }) => {
            fullAuthCleanup();
          });
          // Call auth store logout to reset state
          import("@/features/auth/stores/useAuthStore").then(({ useAuthStore }) => {
            useAuthStore.getState().logout();
          });
          const msg = formatError("backend", "unauthorized");
          toast.error(msg.title, {
            description: msg.description,
          });
          setTimeout(() => {
            window.location.href = "/login";
          }, 1000);
        }
        return Promise.reject(error);
      }

      // Try to refresh token
      if (!isRefreshing) {
        isRefreshing = true;
        const refreshToken = typeof window !== "undefined" ? localStorage.getItem("refreshToken") : null;

        if (!refreshToken) {
          // No refresh token, full cleanup and logout
          isRefreshing = false;
          if (typeof window !== "undefined") {
            import("@/features/auth/utils/clear-auth-cookies").then(({ fullAuthCleanup }) => {
              fullAuthCleanup();
            });
            import("@/features/auth/stores/useAuthStore").then(({ useAuthStore }) => {
              useAuthStore.getState().logout();
            });
            const msg = formatError("backend", "unauthorized");
            toast.error(msg.title, {
              description: msg.description,
            });
            setTimeout(() => {
              window.location.href = "/login";
            }, 1000);
          }
          processQueue(error, null);
          return Promise.reject(error);
        }

        // Create separate axios instance for refresh to avoid circular dependency
        const refreshClient = axios.create({
          baseURL: `${API_BASE_URL}/api/v1`,
          headers: {
            "Content-Type": "application/json",
          },
          timeout: 10000,
        });

        // Call refresh token endpoint directly
        return refreshClient
          .post<{
            success: boolean;
            data?: {
              user: {
                id: string;
                email: string;
                name: string;
                role: string;
                permissions: string[];
                created_at: string;
                updated_at: string;
              };
              token: string;
              refresh_token: string;
              expires_in: number;
            };
          }>("/auth/refresh", {
            refresh_token: refreshToken,
          })
          .then((refreshResponse) => {
            const response = refreshResponse.data;
            if (response.success && response.data) {
              const { user, token, refresh_token } = response.data;
              if (typeof window !== "undefined") {
                localStorage.setItem("token", token);
                localStorage.setItem("refreshToken", refresh_token);
                // Update secure cookie
                setSecureCookie("token", token);
              }

              // Update auth store with user data and tokens
              import("@/features/auth/stores/useAuthStore").then(({ useAuthStore }) => {
                useAuthStore.getState().setToken(token);
                useAuthStore.getState().setUser(user);
                useAuthStore.setState({
                  refreshToken: refresh_token,
                  isAuthenticated: true,
                });
                // Mark session as verified after successful refresh
                useAuthStore.getState().setSessionVerified(true);
              });

              // Update original request with new token
              if (originalRequest?.headers) {
                originalRequest.headers.Authorization = `Bearer ${token}`;
              }
              originalRequest._retry = true;

              processQueue(null, token);
              isRefreshing = false;

              // Retry original request
              return apiClient(originalRequest);
            } else {
              throw new Error("Refresh token failed");
            }
          })
          .catch((refreshError) => {
            // Refresh failed, full cleanup and logout
            isRefreshing = false;
            if (typeof window !== "undefined") {
              import("@/features/auth/utils/clear-auth-cookies").then(({ fullAuthCleanup }) => {
                fullAuthCleanup();
              });
              import("@/features/auth/stores/useAuthStore").then(({ useAuthStore }) => {
                useAuthStore.getState().logout();
              });
              const msg = formatError("backend", "unauthorized");
              toast.error(msg.title, {
                description: msg.description,
              });
              setTimeout(() => {
                window.location.href = "/login";
              }, 1000);
            }
            processQueue(refreshError as AxiosError, null);
            return Promise.reject(refreshError);
          });
      } else {
        // Already refreshing, queue this request
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            if (originalRequest?.headers && token) {
              originalRequest.headers.Authorization = `Bearer ${token}`;
            }
            originalRequest._retry = true;
            return apiClient(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }
    } else if (status === 403) {
      const msg = formatError("backend", "forbidden");
      toast.error(msg.title, {
        description: msg.description,
      });
    } else if (status === 404) {
      // Only show 404 toast for mutations, not queries (to avoid showing on page refresh)
      const isMutation = error.config?.method && ["post", "put", "patch", "delete"].includes(error.config.method.toLowerCase());
      if (isMutation) {
        const msg = formatError("backend", "notFound");
        toast.error(msg.title, {
          description: msg.description,
        });
      }
    } else if (status === 409) {
      // Conflict - already handled above but keep as fallback
      const msg = formatError("backend", "conflict");
      toast.error(msg.title, {
        description: msg.description,
      });
    } else if (status === 503) {
      const msg = formatError("backend", "serviceUnavailable");
      toast.error(msg.title, {
        description: msg.description,
      });
    } else if (status === 429) {
      // Extract rate limit reset time from response headers
      // Axios normalizes headers to lowercase, but check both cases for safety
      const headers = error.response?.headers || {};
      const resetHeader = headers["x-ratelimit-reset"] ||
        headers["X-RateLimit-Reset"];

      if (resetHeader) {
        const resetTimeValue = typeof resetHeader === "string"
          ? parseInt(resetHeader, 10)
          : (typeof resetHeader === "number" ? resetHeader : null);

        if (resetTimeValue !== null && !isNaN(resetTimeValue) && resetTimeValue > 0) {
          // Store reset time for countdown display
          // The useRateLimitCountdown hook will show toast with countdown
          useRateLimitStore.getState().setResetTime(resetTimeValue);

          // Debug: Log in development
          if (process.env.NODE_ENV === "development") {
            console.log("Rate limit detected - Reset time:", resetTimeValue, "Headers:", Object.keys(headers));
          }
        }
      } else {
        // Debug: Log if header not found
        if (process.env.NODE_ENV === "development") {
          console.warn("Rate limit header not found. Available headers:", Object.keys(headers));
        }
      }

      // Override error message to prevent default axios message
      // The useRateLimitCountdown hook will show toast with countdown
      // Set a custom error message that will be shown in form
      const rateLimitMessage = errorData?.error?.message ||
        "Too many login attempts. Please wait before trying again.";

      // Override the error message to prevent "Request failed with status code 429"
      // Create new error object to ensure message is properly set
      const customError = { ...error } as AxiosError<ApiErrorResponse>;
      customError.message = rateLimitMessage;

      if (customError.response?.data) {
        if (!customError.response.data.error) {
          customError.response.data.error = {
            code: "RATE_LIMIT_EXCEEDED",
            message: rateLimitMessage
          };
        } else {
          customError.response.data.error.message = rateLimitMessage;
        }
      }

      // Don't show generic toast here - useRateLimitCountdown hook will show countdown toast
      // Just reject the error so it can be handled by the component
      return Promise.reject(customError);
    } else if (status >= 500) {
      const msg = formatError("backend", "serverError");
      toast.error(msg.title, {
        description: msg.description,
      });
    } else {
      // Other 4xx errors
      // Skip toast for login and check-in endpoints - error will be shown in form/component
      const requestUrl = error.config?.url || "";
      if (!requestUrl.includes("/auth/login") && !requestUrl.includes("/visit-reports") && !requestUrl.includes("/check-in")) {
        const msg = formatError("backend", "unexpectedError");
        toast.error(msg.title, {
          description: msg.description,
        });
      }
    }

    return Promise.reject(error);
  }
);

export default apiClient;
