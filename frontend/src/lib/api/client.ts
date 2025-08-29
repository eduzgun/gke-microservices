import { PUBLIC_GO_API_BASE } from '$env/static/public';

const API_BASE = PUBLIC_GO_API_BASE || 'http://localhost:8080';

export const apiClient = {
  async get<T>(url: string, { fetch: customFetch }: { fetch?: typeof globalThis.fetch } = {}): Promise<T> {
    const response = await (customFetch || fetch)(`${API_BASE}${url}`, {
      credentials: 'include' as const
    });

    if (!response.ok) {
      const text = await response.text();
      const errorMessage = text || response.statusText;
      if (response.status === 401) {
        throw new Error('You must be logged in to access this resource.');
      }
      throw new Error(errorMessage);
    }

    return response.json();
  },

  async post<T>(url: string, data: any, { fetch: customFetch }: { fetch?: typeof globalThis.fetch } = {}): Promise<T> {
    const response = await (customFetch || fetch)(`${API_BASE}${url}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
      credentials: 'include' as const
    });

    if (!response.ok) {
      const text = await response.text();
      const errorMessage = text || response.statusText;
      if (response.status === 401) {
        throw new Error('You must be logged in to perform this action.');
      }
      throw new Error(errorMessage);
    }

    return response.json();
  }
};