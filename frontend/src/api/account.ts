import apiClient from './axios';

export interface LinkedAccount {
  id: string;
  name: string;
  email: string;
  created_at: string;
}

export interface SwitchAccountRequest {
  target_user_id: string;
}

export interface SwitchAccountResponse {
  token: string;
}

export interface LinkAccountRequest {
  email: string;
  password: string;
}

// Получение списка связанных аккаунтов
export const getLinkedAccounts = async (): Promise<LinkedAccount[]> => {
  const response = await apiClient.get('/api/account/links');
  return response.data;
};

// Переключение на другой аккаунт
export const switchAccount = async (data: SwitchAccountRequest): Promise<SwitchAccountResponse> => {
  const response = await apiClient.post('/api/account/switch', data);
  return response.data;
};

// Привязка нового аккаунта
export const linkAccount = async (data: LinkAccountRequest): Promise<void> => {
  await apiClient.post('/api/account/link', data);
}; 