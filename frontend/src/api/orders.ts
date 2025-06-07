import apiClient from './axios';
import type { 
  Order, 
  OrdersResponse, 
  OrdersFilters, 
  UpdateOrderStatusRequest 
} from '../types/order';

// Получение списка заказов с фильтрами
export const getOrders = async (filters?: OrdersFilters): Promise<OrdersResponse> => {
  const params = new URLSearchParams();
  
  if (filters) {
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, value.toString());
      }
    });
  }
  
  const response = await apiClient.get(`/api/orders?${params.toString()}`);
  return response.data;
};

// Получение детальной информации о заказе
export const getOrder = async (orderId: string): Promise<Order> => {
  const response = await apiClient.get(`/api/orders/${orderId}`);
  return response.data;
};

// Обновление статуса заказа
export const updateOrderStatus = async (
  orderId: string, 
  statusData: UpdateOrderStatusRequest
): Promise<{ message: string; order: { id: string; status: string; updated_date: string } }> => {
  const response = await apiClient.put(`/api/orders/${orderId}/status`, statusData);
  return response.data;
}; 