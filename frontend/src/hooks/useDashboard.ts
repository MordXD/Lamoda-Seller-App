import { useState, useEffect, useCallback, useRef } from 'react';
import { getDashboardStats } from '../api/dashboard';
import type { DashboardResponse, AnalyticsFilters } from '../types/dashboard';

// Кэш для хранения данных дашборда
const dashboardCache = new Map<string, { data: DashboardResponse; timestamp: number }>();
const CACHE_DURATION = 5 * 60 * 1000; // 5 минут

// Генерация ключа кэша на основе фильтров
const getCacheKey = (filters?: AnalyticsFilters): string => {
  if (!filters) return 'default';
  return JSON.stringify(filters);
};

// Проверка актуальности кэша
const isCacheValid = (timestamp: number): boolean => {
  return Date.now() - timestamp < CACHE_DURATION;
};

export const useDashboard = (initialFilters?: AnalyticsFilters) => {
  const [dashboardData, setDashboardData] = useState<DashboardResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastFetch, setLastFetch] = useState<number>(0);

  const filtersRef = useRef(initialFilters);
  const abortControllerRef = useRef<AbortController | null>(null);
  
  filtersRef.current = initialFilters;

  const loadDashboard = useCallback(async (newFilters?: AnalyticsFilters, forceRefresh = false) => {
    const filtersToUse = newFilters || filtersRef.current;
    const cacheKey = getCacheKey(filtersToUse);
    
    // Отменяем предыдущий запрос, если он еще выполняется
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    
    // Создаем новый AbortController для текущего запроса
    abortControllerRef.current = new AbortController();
    
    // Проверяем кэш, если не требуется принудительное обновление
    if (!forceRefresh) {
      const cached = dashboardCache.get(cacheKey);
      if (cached && isCacheValid(cached.timestamp)) {
        setDashboardData(cached.data);
        setIsLoading(false);
        setError(null);
        return;
      }
    }

    try {
      setIsLoading(true);
      setError(null);
      
      const response = await getDashboardStats(filtersToUse);
      
      // Проверяем, не был ли запрос отменен
      if (abortControllerRef.current?.signal.aborted) {
        return;
      }
      
      // Сохраняем в кэш
      dashboardCache.set(cacheKey, {
        data: response,
        timestamp: Date.now()
      });
      
      setDashboardData(response);
      setLastFetch(Date.now());
    } catch (err: any) {
      // Игнорируем ошибки отмененных запросов
      if (err.name === 'AbortError') {
        return;
      }
      
      console.error('Ошибка загрузки данных дашборда:', err);
      setError('Не удалось загрузить данные дашборда');
    } finally {
      setIsLoading(false);
      abortControllerRef.current = null;
    }
  }, []);

  // Очистка кэша
  const clearCache = useCallback(() => {
    dashboardCache.clear();
  }, []);

  // Проверка, нужно ли обновить данные
  const shouldRefresh = useCallback(() => {
    return Date.now() - lastFetch > CACHE_DURATION;
  }, [lastFetch]);

  useEffect(() => {
    loadDashboard();
    
    // Очистка при размонтировании компонента
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, [loadDashboard]);

  // Автоматическое обновление данных каждые 5 минут
  useEffect(() => {
    const interval = setInterval(() => {
      if (shouldRefresh()) {
        loadDashboard(filtersRef.current, true);
      }
    }, CACHE_DURATION);

    return () => clearInterval(interval);
  }, [loadDashboard, shouldRefresh]);

  const refetch = useCallback((newFilters?: AnalyticsFilters) => {
    return loadDashboard(newFilters, true);
  }, [loadDashboard]);

  return {
    dashboardData,
    isLoading,
    error,
    refetch,
    clearCache,
    shouldRefresh: shouldRefresh(),
    lastFetch,
  };
}; 