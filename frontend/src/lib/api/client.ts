import { env } from '$env/dynamic/public';

const getApiBase = (): string => {
  // Use explicit env var if provided (from ConfigMap in production)
  if (env.PUBLIC_GO_API_BASE) {
    return env.PUBLIC_GO_API_BASE;
  }
  
  // Development fallback
  if (env.PUBLIC_ENVIRONMENT == "dev") {
    return 'http://localhost:8080';
  }
  
  throw new Error('API_BASE not configured. Check your environment variables.');
};

const API_BASE = getApiBase();
console.log("This is backend URL", API_BASE)

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