import apiClient from './axios';

export interface User {
  id: string;
  name: string;
  email: string;
  balance_kopecks: number;
  created_at: string;
}

export interface UpdateProfileRequest {
  name: string;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

export interface BalanceResponse {
  balance_kopecks: number;
}

export interface AddBalanceRequest {
  amount_kopecks: number;
}

export interface WithdrawBalanceRequest {
  amount_kopecks: number;
}

// Получение профиля пользователя
export const getProfile = async (): Promise<User> => {
  const response = await apiClient.get('/api/profile');
  return response.data;
};

// Обновление профиля пользователя
export const updateProfile = async (data: UpdateProfileRequest): Promise<void> => {
  await apiClient.put('/api/profile', data);
};

// Смена пароля
export const changePassword = async (data: ChangePasswordRequest): Promise<void> => {
  await apiClient.post('/api/password/change', data);
};

// Получение баланса
export const getBalance = async (): Promise<BalanceResponse> => {
  const response = await apiClient.get('/api/balance');
  return response.data;
};

// Пополнение баланса
export const addBalance = async (data: AddBalanceRequest): Promise<void> => {
  await apiClient.post('/api/balance/add', data);
};

// Снятие средств
export const withdrawBalance = async (data: WithdrawBalanceRequest): Promise<void> => {
  await apiClient.post('/api/balance/withdraw', data);
}; 