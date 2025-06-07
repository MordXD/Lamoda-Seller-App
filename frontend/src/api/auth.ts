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

// Интерфейсы для авторизации - только email
export interface LoginData {
  email: string; // Бэкенд ожидает именно email
  password: string;
}

export interface LoginResponse {
  token: string;
  expires_in?: number;
  user?: {
    id: string;
    email: string;
    name: string;
    shop_name?: string;
    is_verified?: boolean;
    registration_date?: string;
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

// Авторизация - только email + пароль
export const login = async (data: LoginData): Promise<LoginResponse> => {
  console.log('API: Отправляем запрос логина:', {
    url: '/api/auth/login',
    data: {
      email: data.email,
      password_length: data.password.length
    }
  });
  
  try {
    const response = await apiClient.post('/api/auth/login', data);
    console.log('API: Успешный ответ логина:', response.data);
    return response.data;
  } catch (error) {
    console.error('API: Ошибка логина:', error);
    throw error;
  }
};

// Обновление токена
export const refreshToken = async (): Promise<RefreshResponse> => {
  const response = await apiClient.post('/api/auth/refresh');
  return response.data;
}; 