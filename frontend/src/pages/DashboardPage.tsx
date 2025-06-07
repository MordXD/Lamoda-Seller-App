import { useState } from 'react';
import SalesChart from '../components/SalesChart';
import TabBar from '../components/TabBar';
import { useDashboard } from '../hooks/useDashboard';
import type { AnalyticsFilters } from '../types/dashboard';

export default function DashboardPage() {
  const [selectedDay, setSelectedDay] = useState(4); // Сегодня - четверг (04)
  const [chartMode, setChartMode] = useState<'orders' | 'purchases'>('orders');
  const [filters] = useState<AnalyticsFilters>({ period: 'today' });

  const { dashboardData, isLoading, error, refetch } = useDashboard(filters);

  const weekDays = [
    { short: 'пн', full: 'понедельник', date: '02' },
    { short: 'вт', full: 'вторник', date: '03' },
    { short: 'ср', full: 'среда', date: '04' },
    { short: 'чт', full: 'четверг', date: '05' },
    { short: 'пт', full: 'пятница', date: '06' },
    { short: 'сб', full: 'суббота', date: '07' },
    { short: 'вс', full: 'воскресенье', date: '08' }
  ];

  const handleDaySelect = (dayIndex: number) => {
    setSelectedDay(dayIndex);
    // Здесь можно добавить логику для изменения фильтров по дате
    // const selectedDate = weekDays[dayIndex];
    // setFilters({ ...filters, date_from: selectedDate.date });
  };

  const handleRefresh = () => {
    refetch(filters);
  };

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* Header with Calendar */}
      <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
        {/* Navigation arrows and "Сегодня" */}
        <div className="flex items-center justify-center mb-4">
          <button className="p-2">
            <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          
          <div className="mx-8 px-6 py-2 bg-gray-100 rounded-lg">
            <span className="text-lg font-medium text-gray-900">Сегодня</span>
            <svg className="inline w-4 h-4 ml-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
          </div>
          
          <button className="p-2">
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
          {/* Amount KPI */}
          <div className="bg-white rounded-lg p-6 shadow-sm border border-gray-100">
            {isLoading ? (
              <div className="animate-pulse">
                <div className="h-8 bg-gray-200 rounded w-1/2 mb-2"></div>
                <div className="h-4 bg-gray-200 rounded w-1/3"></div>
              </div>
            ) : (
              <>
                <div className="flex items-baseline space-x-2 mb-2">
                  <span className="text-3xl font-bold text-gray-900">
                    {dashboardData?.stats.revenue.today.toLocaleString('ru-RU') || 0}
                  </span>
                  <span className="text-lg text-gray-600">₽</span>
                  <div className="flex items-center space-x-1 ml-4">
                    <svg className={`w-4 h-4 ${(dashboardData?.stats.revenue.change_percent || 0) >= 0 ? 'text-green-500' : 'text-red-500'}`} fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d={`${(dashboardData?.stats.revenue.change_percent || 0) >= 0 ? 'M5.293 7.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L6.707 7.707a1 1 0 01-1.414 0z' : 'M14.707 12.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 14.586V3a1 1 0 112 0v11.586l2.293-2.293a1 1 0 011.414 0z'}`} clipRule="evenodd" />
                    </svg>
                    <span className={`text-sm font-medium ${(dashboardData?.stats.revenue.change_percent || 0) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                      {(dashboardData?.stats.revenue.change_percent || 0) >= 0 ? '+' : ''} {dashboardData?.stats.revenue.change_percent || 0}%
                    </span>
                  </div>
                </div>
                <p className="text-sm text-gray-500">Заказали на сумму</p>
              </>
            )}
          </div>

          {/* Orders KPI */}
          <div className="bg-white rounded-lg p-6 shadow-sm border border-gray-100">
            {isLoading ? (
              <div className="animate-pulse">
                <div className="h-8 bg-gray-200 rounded w-1/2 mb-2"></div>
                <div className="h-4 bg-gray-200 rounded w-1/3"></div>
              </div>
            ) : (
              <>
                <div className="flex items-baseline space-x-2 mb-2">
                  <span className="text-3xl font-bold text-gray-900">
                    {dashboardData?.stats.orders.today || 0}
                  </span>
                  <span className="text-lg text-gray-600">шт.</span>
                  <div className="flex items-center space-x-1 ml-4">
                    <svg className={`w-4 h-4 ${(dashboardData?.stats.orders.change_count || 0) >= 0 ? 'text-green-500' : 'text-red-500'}`} fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d={`${(dashboardData?.stats.orders.change_count || 0) >= 0 ? 'M5.293 7.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L6.707 7.707a1 1 0 01-1.414 0z' : 'M14.707 12.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 14.586V3a1 1 0 112 0v11.586l2.293-2.293a1 1 0 011.414 0z'}`} clipRule="evenodd" />
                    </svg>
                    <span className={`text-sm font-medium ${(dashboardData?.stats.orders.change_count || 0) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                      {(dashboardData?.stats.orders.change_count || 0) >= 0 ? '+' : ''} {dashboardData?.stats.orders.change_count || 0}шт.
                    </span>
                  </div>
                </div>
                <p className="text-sm text-gray-500">Кол-во товаров</p>
              </>
            )}
          </div>
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
              <button className="px-4 py-2 bg-gray-100 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-200 transition-colors">
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
                data={dashboardData?.sales_chart || []}
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
        {dashboardData?.recent_orders && dashboardData.recent_orders.length > 0 && (
          <div className="bg-white rounded-lg shadow-sm border border-gray-100">
            <div className="p-4 border-b border-gray-100">
              <h2 className="text-lg font-semibold text-gray-900">Последние заказы</h2>
            </div>
            <div className="divide-y divide-gray-100">
              {dashboardData.recent_orders.slice(0, 5).map((order) => (
                <div key={order.id} className="p-4 hover:bg-gray-50 transition-colors">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium text-gray-900">
                        {order.order_number}
                      </p>
                      <p className="text-xs text-gray-500">
                        {order.customer_name}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-semibold text-gray-900">
                        {order.amount.toLocaleString('ru-RU')} ₽
                      </p>
                      <p className="text-xs text-gray-500">
                        {order.status}
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      <TabBar />
    </div>
  );
} 