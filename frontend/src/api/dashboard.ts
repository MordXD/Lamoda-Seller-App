import apiClient from './axios';
import type { 
  DashboardResponse, 
  AnalyticsFilters,
  SizeChartResponse 
} from '../types/dashboard';

// Получение статистики дашборда
export const getDashboardStats = async (filters?: AnalyticsFilters): Promise<DashboardResponse> => {
  const params = new URLSearchParams();
  
  if (filters) {
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, value.toString());
      }
    });
  }
  
  const response = await apiClient.get(`/api/dashboard/stats?${params.toString()}`);
  return response.data;
};

// Получение размерной сетки
export const getSizeChart = async (category: string): Promise<SizeChartResponse> => {
  const response = await apiClient.get(`/api/products/sizes?category=${category}`);
  return response.data;
}; 