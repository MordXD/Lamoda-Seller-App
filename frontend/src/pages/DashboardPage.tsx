import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import SalesChart from '../components/SalesChart';
import TabBar from '../components/TabBar';
import OrderItem from '../components/OrderItem';
import { useDashboard } from '../hooks/useDashboard';
import { useOrders } from '../hooks/useOrders';
import type { AnalyticsFilters } from '../types/dashboard';

export default function DashboardPage() {
  const navigate = useNavigate();
  const [selectedDay, setSelectedDay] = useState(4); // Сегодня - четверг (04)
  const [chartMode, setChartMode] = useState<'orders' | 'purchases'>('orders');
  const [filters, setFilters] = useState<AnalyticsFilters>({ period: 'today' });
  const [showDatePicker, setShowDatePicker] = useState(false);

  const { dashboardData, isLoading, error, refetch } = useDashboard(filters);
  
  // Получаем последние заказы для отображения
  const { orders: recentOrders, isLoading: ordersLoading } = useOrders({ 
    limit: 5, 
    sort_by: 'date', 
    sort_order: 'desc' 
  });

  const weekDays = [
    { short: 'пн', full: 'понедельник', date: '02', fullDate: '2024-12-02' },
    { short: 'вт', full: 'вторник', date: '03', fullDate: '2024-12-03' },
    { short: 'ср', full: 'среда', date: '04', fullDate: '2024-12-04' },
    { short: 'чт', full: 'четверг', date: '05', fullDate: '2024-12-05' },
    { short: 'пт', full: 'пятница', date: '06', fullDate: '2024-12-06' },
    { short: 'сб', full: 'суббота', date: '07', fullDate: '2024-12-07' },
    { short: 'вс', full: 'воскресенье', date: '08', fullDate: '2024-12-08' }
  ];

  const handleDaySelect = useCallback((dayIndex: number) => {
    setSelectedDay(dayIndex);
    const selectedDate = weekDays[dayIndex];
    const newFilters: AnalyticsFilters = {
      period: 'today',
      date_from: selectedDate.fullDate,
      date_to: selectedDate.fullDate
    };
    setFilters(newFilters);
    refetch(newFilters);
  }, [refetch, weekDays]);

  const handlePreviousWeek = useCallback(() => {
    // Логика для перехода к предыдущей неделе
    const newFilters: AnalyticsFilters = {
      period: 'week',
      date_from: '2024-11-25',
      date_to: '2024-12-01'
    };
    setFilters(newFilters);
    refetch(newFilters);
  }, [refetch]);

  const handleNextWeek = useCallback(() => {
    // Логика для перехода к следующей неделе
    const newFilters: AnalyticsFilters = {
      period: 'week',
      date_from: '2024-12-09',
      date_to: '2024-12-15'
    };
    setFilters(newFilters);
    refetch(newFilters);
  }, [refetch]);

  const handlePeriodSelect = useCallback((period: string) => {
    let newFilters: AnalyticsFilters;
    
    switch (period) {
      case 'today':
        newFilters = { period: 'today' };
        break;
      case 'yesterday':
        newFilters = { period: 'yesterday' };
        break;
      case 'week':
        newFilters = { period: 'week' };
        break;
      case 'month':
        newFilters = { period: 'month' };
        break;
      default:
        newFilters = { period: 'today' };
    }
    
    setFilters(newFilters);
    refetch(newFilters);
    setShowDatePicker(false);
  }, [refetch]);

  const handleRefresh = useCallback(() => {
    refetch(filters);
  }, [refetch, filters]);

  const handleRevenueClick = useCallback(() => {
    navigate('/analytics?tab=revenue');
  }, [navigate]);

  const handleOrdersClick = useCallback(() => {
    navigate('/orders');
  }, [navigate]);

  const handleOrderClick = useCallback((orderId: string) => {
    navigate(`/orders/${orderId}`);
  }, [navigate]);

  const handleChartDateFilter = useCallback(() => {
    setShowDatePicker(!showDatePicker);
  }, [showDatePicker]);

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* Header with Calendar */}
      <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
        {/* Navigation arrows and Period selector */}
        <div className="flex items-center justify-center mb-4">
          <button 
            onClick={handlePreviousWeek}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          
          <div className="relative mx-8">
            <button 
              onClick={() => setShowDatePicker(!showDatePicker)}
              className="px-6 py-2 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
            >
              <span className="text-lg font-medium text-gray-900">
                {filters.period === 'today' ? 'Сегодня' : 
                 filters.period === 'yesterday' ? 'Вчера' :
                 filters.period === 'week' ? 'Неделя' :
                 filters.period === 'month' ? 'Месяц' : 'Сегодня'}
              </span>
              <svg className="inline w-4 h-4 ml-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
            </button>
            
            {/* Dropdown menu */}
            {showDatePicker && (
              <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-lg shadow-lg z-20">
                <button 
                  onClick={() => handlePeriodSelect('today')}
                  className="w-full px-4 py-2 text-left hover:bg-gray-50 first:rounded-t-lg"
                >
                  Сегодня
                </button>
                <button 
                  onClick={() => handlePeriodSelect('yesterday')}
                  className="w-full px-4 py-2 text-left hover:bg-gray-50"
                >
                  Вчера
                </button>
                <button 
                  onClick={() => handlePeriodSelect('week')}
                  className="w-full px-4 py-2 text-left hover:bg-gray-50"
                >
                  Неделя
                </button>
                <button 
                  onClick={() => handlePeriodSelect('month')}
                  className="w-full px-4 py-2 text-left hover:bg-gray-50 last:rounded-b-lg"
                >
                  Месяц
                </button>
              </div>
            )}
          </div>
          
          <button 
            onClick={handleNextWeek}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>

        {/* Week Calendar */}
        <div className="flex justify-center space-x-4">
          {weekDays.map((day, index) => (
            <button
              key={index}
              onClick={() => handleDaySelect(index)}
              className={`flex flex-col items-center p-2 rounded-lg transition-colors ${
                selectedDay === index 
                  ? 'bg-black text-white' 
                  : 'text-gray-600 hover:bg-gray-100'
              }`}
            >
              <span className="text-xs uppercase mb-1">{day.short}</span>
              <span className="text-lg font-semibold">{day.date}</span>
            </button>
          ))}
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

        {/* Main KPIs */}
        <div className="space-y-4">
          {/* Revenue KPI - Clickable */}
          <button 
            onClick={handleRevenueClick}
            className="w-full bg-white rounded-lg p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow text-left"
          >
            {isLoading ? (
              <div className="animate-pulse">
                <div className="h-8 bg-gray-200 rounded w-1/2 mb-2"></div>
                <div className="h-4 bg-gray-200 rounded w-1/3"></div>
              </div>
            ) : (
              <>
                <div className="flex items-baseline space-x-2 mb-2">
                  <span className="text-3xl font-bold text-gray-900">
                    {dashboardData?.revenue?.current?.toLocaleString('ru-RU') || 0}
                  </span>
                  <span className="text-lg text-gray-600">₽</span>
                  <div className="flex items-center space-x-1 ml-4">
                    <svg className={`w-4 h-4 ${(dashboardData?.revenue?.change_percent || 0) >= 0 ? 'text-green-500' : 'text-red-500'}`} fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d={`${(dashboardData?.revenue?.change_percent || 0) >= 0 ? 'M5.293 7.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L6.707 7.707a1 1 0 01-1.414 0z' : 'M14.707 12.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 14.586V3a1 1 0 112 0v11.586l2.293-2.293a1 1 0 011.414 0z'}`} clipRule="evenodd" />
                    </svg>
                    <span className={`text-sm font-medium ${(dashboardData?.revenue?.change_percent || 0) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                      {(dashboardData?.revenue?.change_percent || 0) >= 0 ? '+' : ''}{Math.round(dashboardData?.revenue?.change_percent || 0)}%
                    </span>
                  </div>
                </div>
                <div className="flex items-center justify-between">
                  <p className="text-sm text-gray-500">Заказали на сумму</p>
                  <svg className="w-5 h-5 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
                  </svg>
                </div>
              </>
            )}
          </button>

          {/* Orders KPI - Clickable */}
          <button 
            onClick={handleOrdersClick}
            className="w-full bg-white rounded-lg p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow text-left"
          >
            {isLoading ? (
              <div className="animate-pulse">
                <div className="h-8 bg-gray-200 rounded w-1/2 mb-2"></div>
                <div className="h-4 bg-gray-200 rounded w-1/3"></div>
              </div>
            ) : (
              <>
                <div className="flex items-baseline space-x-2 mb-2">
                  <span className="text-3xl font-bold text-gray-900">
                    {dashboardData?.orders?.current || 0}
                  </span>
                  <span className="text-lg text-gray-600">шт.</span>
                  <div className="flex items-center space-x-1 ml-4">
                    <svg className={`w-4 h-4 ${(dashboardData?.orders?.change_absolute || 0) >= 0 ? 'text-green-500' : 'text-red-500'}`} fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d={`${(dashboardData?.orders?.change_absolute || 0) >= 0 ? 'M5.293 7.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L6.707 7.707a1 1 0 01-1.414 0z' : 'M14.707 12.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 14.586V3a1 1 0 112 0v11.586l2.293-2.293a1 1 0 011.414 0z'}`} clipRule="evenodd" />
                    </svg>
                    <span className={`text-sm font-medium ${(dashboardData?.orders?.change_absolute || 0) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                      {(dashboardData?.orders?.change_absolute || 0) >= 0 ? '+' : ''}{Math.round(dashboardData?.orders?.change_absolute || 0)}шт.
                    </span>
                  </div>
                </div>
                <div className="flex items-center justify-between">
                  <p className="text-sm text-gray-500">Кол-во товаров</p>
                  <svg className="w-5 h-5 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
                  </svg>
                </div>
              </>
            )}
          </button>
        </div>

        {/* Chart Section */}
        <div className="bg-white rounded-lg p-6 shadow-sm border border-gray-100">
          {/* Chart Controls */}
          <div className="flex items-center justify-between mb-6">
            <div className="flex bg-black rounded-lg p-1">
              <button
                onClick={() => setChartMode('orders')}
                className={`px-4 py-2 text-sm font-medium rounded-md transition-colors ${
                  chartMode === 'orders'
                    ? 'bg-white text-black'
                    : 'text-white hover:text-gray-300'
                }`}
              >
                Заказали / Выкупили
              </button>
            </div>
            
            <div className="flex items-center space-x-4">
              <button 
                onClick={handleChartDateFilter}
                className="px-4 py-2 bg-gray-100 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-200 transition-colors"
              >
                Дата
                <svg className="inline w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
              </button>
            </div>
          </div>

          {/* Chart */}
          <div className="h-64 relative">
            <div className="absolute left-0 top-0 bottom-0 flex flex-col justify-between text-xs text-gray-500">
              <span>200 тыс.</span>
              <span>150 тыс.</span>
              <span>100 тыс.</span>
              <span>50 тыс.</span>
              <span>0 тыс.</span>
            </div>
            
            <div className="ml-12 h-full">
              <SalesChart
                data={dashboardData?.hourly_sales || []}
                isLoading={isLoading}
              />
            </div>
          </div>

          {/* Chart Legend */}
          <div className="flex items-center justify-center space-x-6 mt-4">
            <div className="flex items-center space-x-2">
              <div className="w-3 h-3 bg-black rounded-full"></div>
              <span className="text-sm text-gray-600">Заказали</span>
            </div>
            <div className="flex items-center space-x-2">
              <div className="w-3 h-3 bg-gray-400 rounded-full"></div>
              <span className="text-sm text-gray-600">Выкупили</span>
            </div>
          </div>
        </div>

        {/* Recent Orders */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-100">
          <div className="p-4 border-b border-gray-100 flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">Последние заказы</h2>
            <button 
              onClick={() => navigate('/orders')}
              className="text-sm text-gray-600 hover:text-gray-900 transition-colors"
            >
              Все заказы →
            </button>
          </div>
          <div className="divide-y divide-gray-100">
            {ordersLoading ? (
              <div className="p-4">
                <div className="animate-pulse space-y-3">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className="flex items-center space-x-4">
                      <div className="w-12 h-12 bg-gray-200 rounded-lg"></div>
                      <div className="flex-1 space-y-2">
                        <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                        <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ) : recentOrders.length > 0 ? (
              <div className="p-4 space-y-3">
                {recentOrders.slice(0, 5).map((order) => (
                  <OrderItem
                    key={order.id}
                    order={{
                      id: order.order_number,
                      date: order.date,
                      amount: order.totals.total,
                      status: order.status === 'new' ? 'Новый' :
                              order.status === 'confirmed' ? 'Подтвержден' :
                              order.status === 'in_transit' ? 'В пути' :
                              order.status === 'delivered' ? 'Доставлен' :
                              order.status === 'returned' ? 'Возвращен' :
                              order.status === 'cancelled' ? 'Отменен' : order.status,
                      firstProductImage: order.items[0]?.image
                    }}
                    onClick={() => handleOrderClick(order.id)}
                  />
                ))}
              </div>
            ) : (
              <div className="p-8 text-center text-gray-500">
                <svg className="w-12 h-12 text-gray-300 mx-auto mb-4" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M6 2a2 2 0 00-2 2v12a2 2 0 002 2h8a2 2 0 002-2V4a2 2 0 00-2-2H6zm1 2a1 1 0 000 2h6a1 1 0 100-2H7zm6 7a1 1 0 011 1v3a1 1 0 11-2 0v-3a1 1 0 011-1zm-3 3a1 1 0 100 2h.01a1 1 0 100-2H10zm-4 1a1 1 0 011-1h.01a1 1 0 110 2H7a1 1 0 01-1-1zm1-4a1 1 0 100 2h.01a1 1 0 100-2H7zm2 0a1 1 0 100 2h.01a1 1 0 100-2H9zm2 0a1 1 0 100 2h.01a1 1 0 100-2H11z" clipRule="evenodd" />
                </svg>
                <p className="text-sm">Заказов пока нет</p>
              </div>
            )}
          </div>
        </div>
      </div>

      <TabBar />
    </div>
  );
} 