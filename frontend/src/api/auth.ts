import apiClient from './axios';

// Интерфейсы для регистрации - упрощенная версия
export interface RegisterData {
  name: string;
  email: string;
}

export interface RegisterResponse {
  success: boolean;
  message: string;
  temporary_password: string;
  user_id: string;
  email: string;
}

// Интерфейсы для авторизации - поддержка username или email
export interface LoginData {
  username_or_email: string; // Может быть как username, так и email
  password: string;
}

export interface LoginResponse {
  token: string;
  expires_in: number;
  user: {
    id: string;
    email: string;
    name: string;
    shop_name: string;
    is_verified: boolean;
    registration_date: string;
  };
}

export interface RefreshResponse {
  token: string;
  expires_in: number;
}

// Регистрация селлера - только имя и email
export const register = async (data: RegisterData): Promise<RegisterResponse> => {
  const response = await apiClient.post('/api/auth/register', data);
  return response.data;
};

// Авторизация - username или email + пароль
export const login = async (data: LoginData): Promise<LoginResponse> => {
  const response = await apiClient.post('/api/auth/login', data);
  return response.data;
};

// Обновление токена
export const refreshToken = async (): Promise<RefreshResponse> => {
  const response = await apiClient.post('/api/auth/refresh');
  return response.data;
}; 