import { apiClient } from './client';

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterCredentials {
  username: string;
  email: string;
  password: string;
}

export interface User {
  id: number;
  username: string;
  email: string;
  status: string;
}

export interface AuthResponse {
  status: string;
  message?: string;
}

export const authApi = {
  login: (credentials: LoginCredentials): Promise<AuthResponse> => 
    apiClient.post<AuthResponse>('/auth/login', credentials),

  register: (data: RegisterCredentials): Promise<AuthResponse> =>
    apiClient.post<AuthResponse>('/auth/register', data),

  logout: (): Promise<AuthResponse> => 
    apiClient.post<AuthResponse>('/auth/logout', {}),

  getProfile: (): Promise<User> => 
    apiClient.get<User>('/auth/profile'),
};