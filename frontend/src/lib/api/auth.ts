import { apiClient } from './client';

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface User {
  id: number;
  username: string;
  email: string;
  status: string;
}

export const authApi = {
  login: (credentials: LoginCredentials) => 
    apiClient.post<User>('/auth/login', credentials),

  register: (data: { username: string; email: string; password: string }) =>
    apiClient.post<User>('/auth/register', data),

  logout: () => 
    apiClient.post('/auth/logout', {}),

  getProfile: () => 
    apiClient.get<User>('/profile'),
};