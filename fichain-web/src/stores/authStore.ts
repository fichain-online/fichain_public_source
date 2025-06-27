// src/stores/authStore.ts
'use client'

import { create } from 'zustand'
import { persist } from 'zustand/middleware'

// Define the state's shape
interface AuthState {
  privateKey: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  _hasHydrated: boolean; // <-- ADD THIS STATE
  login: (key: string) => Promise<void>;
  logout: () => void;
  setHasHydrated: (state: boolean) => void; // <-- ADD THIS ACTION
}

// Create the store
export const useAuthStore = create<AuthState>()(
  // Use the 'persist' middleware to save parts of the state to localStorage
  persist(
    (set, get) => ({
      // Initial State
      privateKey: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
      _hasHydrated: false, // <-- INITIAL STATE IS FALSE

      // --- ACTIONS ---
      //
      setHasHydrated: (state) => { // <-- ACTION TO MANUALLY SET HYDRATION
        set({
          _hasHydrated: state,
        });
      },

      // Login Action
      login: async (key: string) => {
        set({ isLoading: true, error: null });
        const normalizedKey = key.startsWith("0x") ? key.slice(2) : key;
        // DEMO VALIDATION: Check if the key is not empty and has a reasonable length
        if (normalizedKey && normalizedKey.length == 64) {
          // On successful "validation", update the state
          set({
            privateKey: normalizedKey, // Store the key for the current session
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
          console.log("Zustand: Login successful, auth state updated.");
        } else {
          // On failure, update the state with an error
          set({
            privateKey: null,
            isAuthenticated: false,
            isLoading: false,
            error: "Khóa Riêng Tư không hợp lệ. Vui lòng thử lại.",
          });
          console.error("Zustand: Login failed, invalid key.");
        }
      },

      // Logout Action
      logout: () => {
        set({
          privateKey: null,
          isAuthenticated: false,
          error: null,
        });
        console.log("Zustand: User logged out.");
      },
    }),
    {
      name: 'auth-storage', // Unique name for localStorage key
      // By default, persist saves everything. We want to be selective.
      // We will only persist the `isAuthenticated` flag to keep the user logged in
      // across page reloads, but we will NEVER persist the private key for security.
      // for demo now: i will store it
      partialize: (state) => ({
        isAuthenticated: state.isAuthenticated,
        privateKey: state.privateKey,
      }),
      onRehydrateStorage: () => (state) => {
        if (state) {
          state.setHasHydrated(true);
        }
      },
    }
  )
);
