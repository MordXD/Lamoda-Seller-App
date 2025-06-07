import { useState, useEffect, useCallback, useRef } from 'react';
import { getDashboardStats } from '../api/dashboard';
import type { DashboardResponse, AnalyticsFilters } from '../types/dashboard';

export const useDashboard = (initialFilters?: AnalyticsFilters) => {
  const [dashboardData, setDashboardData] = useState<DashboardResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const filtersRef = useRef(initialFilters);
  filtersRef.current = initialFilters;

  const loadDashboard = useCallback(async (newFilters?: AnalyticsFilters) => {
    try {
      setIsLoading(true);
      setError(null);
      
      const filtersToUse = newFilters || filtersRef.current;
      const response = await getDashboardStats(filtersToUse);
      setDashboardData(response);
    } catch (err) {
      console.error('Ошибка загрузки данных дашборда:', err);
      setError('Не удалось загрузить данные дашборда');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    loadDashboard();
  }, [loadDashboard]);

  const refetch = useCallback((newFilters?: AnalyticsFilters) => {
    return loadDashboard(newFilters);
  }, [loadDashboard]);

  return {
    dashboardData,
    isLoading,
    error,
    refetch,
  };
}; 