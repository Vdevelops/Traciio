"use client";

import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { User, AuthState } from "../types";
import { setSecureCookie } from "@/lib/cookie";

interface AuthStore extends AuthState {
  setUser: (user: User | null) => void;
  setToken: (token: string | null) => void;
  setSessionVerified: (verified: boolean) => void;
  clearError: () => void;
  logout: () => void;
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,
      isSessionVerified: false, // Backend verification - NOT persisted
      isLoading: false,
      error: null,

      setUser: (user: User | null) => {
        set({ user, isAuthenticated: !!user });
      },

      setToken: (token: string | null) => {
        set({ token });
        if (typeof window !== "undefined" && token) {
          localStorage.setItem("token", token);
          setSecureCookie("token", token);
        }
      },

      clearError: () => {
        set({ error: null });
      },

      setSessionVerified: (verified: boolean) => {
        set({ isSessionVerified: verified });
      },

      logout: () => {
        set({
          user: null,
          token: null,
          refreshToken: null,
          isAuthenticated: false,
          isSessionVerified: false,
          error: null,
        });
        if (globalThis.window !== undefined) {
          localStorage.removeItem("token");
          localStorage.removeItem("refreshToken");
          document.cookie = "token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
        }
      },
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        // After rehydration, sync state with localStorage
        if (typeof window !== "undefined" && state) {
          // CRITICAL: Always reset isSessionVerified to false on rehydration
          // This forces backend verification before trusting localStorage auth state
          state.isSessionVerified = false;

          const token = localStorage.getItem("token");
          const refreshToken = localStorage.getItem("refreshToken");

          // Priority: Zustand persisted state > localStorage
          // If Zustand has token, use it and sync to localStorage
          if (state.token) {
            // Sync Zustand state to localStorage
            if (!token || token !== state.token) {
              localStorage.setItem("token", state.token);
            }
            if (state.refreshToken && (!refreshToken || refreshToken !== state.refreshToken)) {
              localStorage.setItem("refreshToken", state.refreshToken);
            }
            // Set authenticated if we have token
            if (!state.isAuthenticated) {
              state.isAuthenticated = true;
            }
            // Set cookie if not exists
            if (!document.cookie.includes("token=")) {
              setSecureCookie("token", state.token);
            }
          } else if (token) {
            // No Zustand token but localStorage has it, restore from localStorage
            state.token = token;
            if (refreshToken) {
              state.refreshToken = refreshToken;
            }
            state.isAuthenticated = true;
            // Set cookie if not exists
            if (!document.cookie.includes("token=")) {
              document.cookie = `token=${token}; path=/; max-age=${7 * 24 * 60 * 60}; SameSite=Lax`;
            }
          } else if (!token && state.isAuthenticated) {
            // No token anywhere but store says authenticated, clear it
            state.isAuthenticated = false;
            state.user = null;
            state.token = null;
            state.refreshToken = null;
          }
        }
      },
    }
  )
);

