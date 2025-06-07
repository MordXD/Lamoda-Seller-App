import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import KPIWidget from '../components/KPIWidget';
import SalesChart from '../components/SalesChart';
import TabBar from '../components/TabBar';

interface SalesData {
  date: string;
  orders: number;
  purchases: number;
}

export default function DashboardPage() {
  const { activeAccount } = useAuth();
  const [isLoading, setIsLoading] = useState(true);
  const [dashboardData, setDashboardData] = useState<{
    todayAmount: number;
    todayAmountChange: number;
    todayOrders: number;
    todayOrdersChange: number;
    salesData: SalesData[];
  }>({
    todayAmount: 0,
    todayAmountChange: 0,
    todayOrders: 0,
    todayOrdersChange: 0,
    salesData: []
  });

  useEffect(() => {
    // Имитация загрузки данных
    setTimeout(() => {
      setDashboardData({
        todayAmount: 125000,
        todayAmountChange: 12.5,
        todayOrders: 18,
        todayOrdersChange: -5.2,
        salesData: [
          { date: '2024-01-15', orders: 85000, purchases: 78000 },
          { date: '2024-01-16', orders: 92000, purchases: 85000 },
          { date: '2024-01-17', orders: 78000, purchases: 72000 },
          { date: '2024-01-18', orders: 105000, purchases: 95000 },
          { date: '2024-01-19', orders: 120000, purchases: 110000 },
          { date: '2024-01-20', orders: 135000, purchases: 125000 },
          { date: '2024-01-21', orders: 125000, purchases: 115000 }
        ]
      });
      setIsLoading(false);
    }, 1000);
  }, []);

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold text-gray-900">Дашборд</h1>
          <div className="text-right">
            <p className="text-sm text-gray-500">{activeAccount?.shopName}</p>
            <p className="text-xs text-gray-400">
              {new Date().toLocaleDateString('ru-RU', { 
                day: 'numeric', 
                month: 'long',
                weekday: 'long'
              })}
            </p>
          </div>
        </div>
      </header>

      {/* Content */}
      <div className="px-4 py-6 space-y-6">
        {/* KPI Widgets */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <KPIWidget
            title="Заказали на сумму"
            value={`${dashboardData.todayAmount.toLocaleString('ru-RU')} ₽`}
            change={dashboardData.todayAmountChange}
            isLoading={isLoading}
          />
          <KPIWidget
            title="Количество заказов"
            value={dashboardData.todayOrders.toString()}
            change={dashboardData.todayOrdersChange}
            isLoading={isLoading}
          />
        </div>

        {/* Sales Chart */}
        <div>
          <h2 className="text-lg font-semibold text-gray-900 mb-4">График продаж</h2>
          <SalesChart
            data={dashboardData.salesData}
            isLoading={isLoading}
          />
        </div>

        {/* Quick Actions */}
        <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Быстрые действия</h3>
          <div className="grid grid-cols-2 gap-3">
            <button className="flex flex-col items-center p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
              <svg className="w-8 h-8 text-gray-600 mb-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V7z" clipRule="evenodd" />
              </svg>
              <span className="text-sm font-medium text-gray-700">Добавить товар</span>
            </button>
            
            <button className="flex flex-col items-center p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
              <svg className="w-8 h-8 text-gray-600 mb-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M6 2a2 2 0 00-2 2v12a2 2 0 002 2h8a2 2 0 002-2V4a2 2 0 00-2-2H6zm1 2a1 1 0 000 2h6a1 1 0 100-2H7zm6 7a1 1 0 011 1v3a1 1 0 11-2 0v-3a1 1 0 011-1zm-3 3a1 1 0 100 2h.01a1 1 0 100-2H10zm-4 1a1 1 0 011-1h.01a1 1 0 110 2H7a1 1 0 01-1-1zm1-4a1 1 0 100 2h.01a1 1 0 100-2H7zm2 0a1 1 0 100 2h.01a1 1 0 100-2H9zm2 0a1 1 0 100 2h.01a1 1 0 100-2H11z" clipRule="evenodd" />
              </svg>
              <span className="text-sm font-medium text-gray-700">Новый заказ</span>
            </button>
          </div>
        </div>
      </div>

      <TabBar />
    </div>
  );
} 