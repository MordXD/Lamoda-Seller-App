import { useState, useEffect, useCallback, useRef } from 'react';
import { getOrders, getOrder, updateOrderStatus } from '../api/orders';
import type { Order, OrdersResponse, OrdersFilters, UpdateOrderStatusRequest } from '../types/order';

export const useOrders = (initialFilters?: OrdersFilters) => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [summary, setSummary] = useState<OrdersResponse['summary'] | null>(null);
  const [pagination, setPagination] = useState<OrdersResponse['pagination'] | null>(null);
  
  // Используем ref для хранения текущих фильтров
  const filtersRef = useRef(initialFilters);
  filtersRef.current = initialFilters;

  const loadOrders = useCallback(async (newFilters?: OrdersFilters) => {
    try {
      setIsLoading(true);
      setError(null);
      
      const filtersToUse = newFilters || filtersRef.current;
      const response = await getOrders(filtersToUse);
      
      setOrders(response.orders);
      setSummary(response.summary);
      setPagination(response.pagination);
    } catch (err) {
      console.error('Ошибка загрузки заказов:', err);
      setError('Не удалось загрузить заказы');
    } finally {
      setIsLoading(false);
    }
  }, []); // Убираем зависимость от initialFilters

  useEffect(() => {
    loadOrders();
  }, [loadOrders]);

  const refetch = useCallback((newFilters?: OrdersFilters) => {
    return loadOrders(newFilters);
  }, [loadOrders]);

  const updateStatus = useCallback(async (orderId: string, statusData: UpdateOrderStatusRequest) => {
    try {
      await updateOrderStatus(orderId, statusData);
      // Обновляем локальное состояние
      setOrders(prevOrders => 
        prevOrders.map(order => 
          order.id === orderId 
            ? { ...order, status: statusData.status, updated_date: new Date().toISOString() }
            : order
        )
      );
      return true;
    } catch (err) {
      console.error('Ошибка обновления статуса заказа:', err);
      throw err;
    }
  }, []);

  return {
    orders,
    isLoading,
    error,
    summary,
    pagination,
    refetch,
    updateStatus,
  };
};

export const useOrder = (orderId: string | undefined) => {
  const [order, setOrder] = useState<Order | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadOrder = useCallback(async () => {
    if (!orderId) {
      setError('ID заказа не указан');
      setIsLoading(false);
      return;
    }

    try {
      setIsLoading(true);
      setError(null);
      
      const orderData = await getOrder(orderId);
      setOrder(orderData);
    } catch (err) {
      console.error('Ошибка загрузки заказа:', err);
      setError('Не удалось загрузить информацию о заказе');
    } finally {
      setIsLoading(false);
    }
  }, [orderId]);

  useEffect(() => {
    loadOrder();
  }, [loadOrder]);

  const refetch = useCallback(() => {
    return loadOrder();
  }, [loadOrder]);

  const updateStatus = useCallback(async (statusData: UpdateOrderStatusRequest) => {
    if (!orderId) return false;
    
    try {
      await updateOrderStatus(orderId, statusData);
      // Обновляем локальное состояние
      setOrder(prevOrder => 
        prevOrder 
          ? { ...prevOrder, status: statusData.status, updated_date: new Date().toISOString() }
          : null
      );
      return true;
    } catch (err) {
      console.error('Ошибка обновления статуса заказа:', err);
      throw err;
    }
  }, [orderId]);

  return {
    order,
    isLoading,
    error,
    refetch,
    updateStatus,
  };
}; 