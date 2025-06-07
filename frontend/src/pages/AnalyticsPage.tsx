import { useState } from 'react';
import TabBar from '../components/TabBar';
import { useDashboard } from '../hooks/useDashboard';
import type { AnalyticsFilters } from '../types/dashboard';

export default function AnalyticsPage() {
  const [filters, setFilters] = useState<AnalyticsFilters>({ period: 'week' });
  const { dashboardData, isLoading, error, refetch } = useDashboard(filters);

  const topCategories = dashboardData?.top_categories || [];
  const totalRevenue = topCategories.reduce((sum, category) => sum + category.revenue, 0);
  const totalOrders = topCategories.reduce((sum, category) => sum + category.orders, 0);

  const handleRefresh = () => {
    refetch(filters);
  };

  const handlePeriodChange = (period: AnalyticsFilters['period']) => {
    const newFilters = { ...filters, period };
    setFilters(newFilters);
    refetch(newFilters);
  };

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Аналитика</h1>
            <p className="text-sm text-gray-500 mt-1">
              {filters.period === 'week' ? 'За последние 7 дней' : 
               filters.period === 'month' ? 'За последний месяц' : 
               filters.period === 'today' ? 'За сегодня' : 'За выбранный период'}
            </p>
          </div>
          
          {/* Period Selector */}
          <div className="flex space-x-2">
            <select
              value={filters.period || 'week'}
              onChange={(e) => handlePeriodChange(e.target.value as AnalyticsFilters['period'])}
              className="text-sm border border-gray-300 rounded-md px-3 py-1 bg-white"
            >
              <option value="today">Сегодня</option>
              <option value="week">Неделя</option>
              <option value="month">Месяц</option>
              <option value="quarter">Квартал</option>
            </select>
          </div>
        </div>
      </header>

      {/* Content */}
      <div className="px-4 py-6 space-y-6">
        {/* Error State */}
        {error && (
          <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-600 text-sm">{error}</p>
            <button
              onClick={handleRefresh}
              className="mt-2 text-red-700 text-sm font-medium hover:text-red-800"
            >
              Попробовать снова
            </button>
          </div>
        )}

        {/* Summary Cards */}
        {!isLoading && !error && (
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
              <h3 className="text-sm font-medium text-gray-500 mb-1">Общая выручка</h3>
              <p className="text-2xl font-bold text-gray-900">
                {totalRevenue.toLocaleString('ru-RU')} ₽
              </p>
            </div>
            <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
              <h3 className="text-sm font-medium text-gray-500 mb-1">Количество заказов</h3>
              <p className="text-2xl font-bold text-gray-900">
                {totalOrders} шт
              </p>
            </div>
          </div>
        )}

        {/* Top Categories Widget */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-100">
          <div className="p-4 border-b border-gray-100">
            <h2 className="text-lg font-semibold text-gray-900">Топ-5 категорий</h2>
            <p className="text-sm text-gray-500">Отсортированы по выручке</p>
          </div>

          {isLoading ? (
            // Loading skeleton
            <div className="p-4 space-y-4">
              {[1, 2, 3, 4, 5].map((i) => (
                <div key={i} className="animate-pulse flex items-center space-x-4 py-3">
                  <div className="w-8 h-8 bg-gray-200 rounded-full flex-shrink-0"></div>
                  <div className="w-12 h-12 bg-gray-200 rounded-lg flex-shrink-0"></div>
                  <div className="flex-1 space-y-2">
                    <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                    <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                  </div>
                  <div className="text-right space-y-2">
                    <div className="h-4 bg-gray-200 rounded w-20"></div>
                    <div className="h-3 bg-gray-200 rounded w-16"></div>
                  </div>
                </div>
              ))}
            </div>
          ) : topCategories.length > 0 ? (
            <div className="divide-y divide-gray-100">
              {topCategories.map((category, index) => (
                <div key={category.category} className="p-4 hover:bg-gray-50 transition-colors">
                  <div className="flex items-center space-x-4">
                    {/* Rank */}
                    <div className="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center flex-shrink-0">
                      <span className="text-sm font-semibold text-gray-600">
                        {index + 1}
                      </span>
                    </div>

                    {/* Category Icon */}
                    <div className="w-12 h-12 bg-gray-200 rounded-lg flex-shrink-0 overflow-hidden">
                      <div className="w-full h-full flex items-center justify-center text-gray-400">
                        <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clipRule="evenodd" />
                        </svg>
                      </div>
                    </div>

                    {/* Category Info */}
                    <div className="flex-1 min-w-0">
                      <h3 className="text-sm font-medium text-gray-900 truncate">{category.name}</h3>
                      <p className="text-xs text-gray-500">{category.orders} заказов • {category.items} товаров</p>
                    </div>

                    {/* Revenue */}
                    <div className="text-right flex-shrink-0">
                      <p className="text-sm font-semibold text-gray-900">
                        {category.revenue.toLocaleString('ru-RU')} ₽
                      </p>
                      <p className="text-xs text-gray-500">
                        {totalRevenue > 0 ? Math.round((category.revenue / totalRevenue) * 100) : 0}% от общей
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            // Empty state
            <div className="p-8 text-center">
              <svg className="w-12 h-12 text-gray-300 mx-auto mb-4" fill="currentColor" viewBox="0 0 20 20">
                <path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zM8 7a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zM14 4a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z" />
              </svg>
              <h3 className="text-lg font-medium text-gray-900 mb-2">Нет данных</h3>
              <p className="text-gray-500">За выбранный период нет продаж</p>
            </div>
          )}
        </div>

        {/* Additional Stats */}
        {!isLoading && !error && (
          <div className="bg-white rounded-lg shadow-sm border border-gray-100">
            <div className="p-4 border-b border-gray-100">
              <h2 className="text-lg font-semibold text-gray-900">Дополнительная статистика</h2>
            </div>
            <div className="p-4 grid grid-cols-2 gap-4">
              <div className="text-center">
                <span className="text-sm text-gray-600">Общая выручка</span>
                <span className="text-lg font-semibold text-gray-900 block">
                  {dashboardData?.revenue?.current?.toLocaleString('ru-RU') || 0} ₽
                </span>
              </div>
              <div className="text-center">
                <span className="text-sm text-gray-600">Общее количество заказов</span>
                <span className="text-lg font-semibold text-gray-900 block">
                  {dashboardData?.orders?.current || 0} шт
                </span>
              </div>
              <div className="text-center">
                <span className="text-sm text-gray-600">Средний чек</span>
                <span className="text-lg font-semibold text-gray-900 block">
                  {dashboardData?.avg_order_value?.current?.toLocaleString('ru-RU') || 0} ₽
                </span>
              </div>
              <div className="text-center">
                <span className="text-sm text-gray-600">Товаров продано</span>
                <span className="text-lg font-semibold text-gray-900 block">
                  {dashboardData?.items_sold?.current || 0} шт
                </span>
              </div>
            </div>
          </div>
        )}
      </div>

      <TabBar />
    </div>
  );
} 