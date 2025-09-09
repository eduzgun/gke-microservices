import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { authApi, type User, type AuthResponse } from '$lib/api/auth';

export type AuthUser = User | null;

interface AuthState {
  user: AuthUser;
  loading: boolean;
  initialized: boolean;
}

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>({
    user: null,
    loading: false,
    initialized: false
  });

  return {
    subscribe,
    
    async init() {
      if (!browser) {
        return;
      }
      
      update(state => ({ ...state, loading: true }));
      
      try {
        const user = await authApi.getProfile();
        set({ user, loading: false, initialized: true });
      } catch (error) {
        set({ user: null, loading: false, initialized: true });
      }
    },

    async login(email: string, password: string): Promise<{ success: boolean; error?: string }> {
      if (!browser) {
        return { success: false, error: 'Not in browser' };
      }
      
      update(state => ({ ...state, loading: true }));
      
      try {
        const response: AuthResponse = await authApi.login({ email, password });
        
        
        if (response.status === 'success') {
          const user = await authApi.getProfile();
          update(state => ({ ...state, user, loading: false }));
          return { success: true };
        } else {
          update(state => ({ ...state, loading: false }));
          return { success: false, error: response.message || 'Login failed' };
        }
      } catch (error) {
        console.error('Auth store: Login error:', error);
        update(state => ({ ...state, loading: false }));
        const errorMessage = error instanceof Error ? error.message : 'Login failed';
        return { success: false, error: errorMessage };
      }
    },

    async register(username: string, email: string, password: string): Promise<{ success: boolean; error?: string }> {
      if (!browser) {
        return { success: false, error: 'Not in browser' };
      }
      
      update(state => ({ ...state, loading: true }));
      
      try {
        const response: AuthResponse = await authApi.register({ username, email, password });
        
        
        if (response.status === 'success') {
          // After registration, login automatically or redirect to login
          update(state => ({ ...state, loading: false }));
          return { success: true };
        } else {
          update(state => ({ ...state, loading: false }));
          return { success: false, error: response.message || 'Registration failed' };
        }
      } catch (error) {
        console.error('Auth store: Registration error:', error);
        update(state => ({ ...state, loading: false }));
        const errorMessage = error instanceof Error ? error.message : 'Registration failed';
        return { success: false, error: errorMessage };
      }
    },

    async logout() {
      if (!browser) {
        return;
      }
      
      
      try {
        await authApi.logout();
      } catch (error) {
        console.error('Auth store: Logout API error:', error);
      }
      
      set({ user: null, loading: false, initialized: true });
    }
  };
}

export const authStore = createAuthStore();