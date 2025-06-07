import { useState, useEffect } from 'react';
import OrderItem from '../components/OrderItem';
import TabBar from '../components/TabBar';

type OrderStatus = 'new' | 'in_transit' | 'archive';

interface Order {
  id: string;
  date: string;
  amount: number;
  status: string;
  firstProductImage?: string;
}

interface SummaryData {
  title: string;
  description: string;
  count: number;
  totalAmount: number;
}

export default function OrdersPage() {
  const [activeSegment, setActiveSegment] = useState<OrderStatus>('new');
  const [orders, setOrders] = useState<Order[]>([]);
  const [summary, setSummary] = useState<SummaryData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  const segments = [
    { id: 'new' as OrderStatus, label: 'Новые' },
    { id: 'in_transit' as OrderStatus, label: 'В пути' },
    { id: 'archive' as OrderStatus, label: 'Архив' }
  ];

  const loadOrders = async (segment: OrderStatus, showLoader = true) => {
    if (showLoader) setIsLoading(true);
    
    // Имитация загрузки данных
    await new Promise(resolve => setTimeout(resolve, 800));
    
    const mockOrders: Record<OrderStatus, Order[]> = {
      new: [
        { id: '12345-001', date: '2024-01-21', amount: 15000, status: 'Ожидает сборки' },
        { id: '12346-002', date: '2024-01-21', amount: 8500, status: 'Подтверждён' },
        { id: '12347-003', date: '2024-01-20', amount: 22000, status: 'Оплачен' }
      ],
      in_transit: [
        { id: '12340-001', date: '2024-01-19', amount: 12000, status: 'В доставке' },
        { id: '12341-002', date: '2024-01-18', amount: 9800, status: 'Передан курьеру' }
      ],
      archive: [
        { id: '12330-001', date: '2024-01-15', amount: 18000, status: 'Доставлен' },
        { id: '12331-002', date: '2024-01-14', amount: 25000, status: 'Выкуплен' },
        { id: '12332-003', date: '2024-01-13', amount: 7200, status: 'Возврат' }
      ]
    };

    const mockSummaries: Record<OrderStatus, SummaryData> = {
      new: {
        title: 'Новые заказы',
        description: 'Требуют обработки и подтверждения',
        count: 3,
        totalAmount: 45500
      },
      in_transit: {
        title: 'Заказы в пути',
        description: 'Находятся в процессе доставки',
        count: 2,
        totalAmount: 21800
      },
      archive: {
        title: 'Архивные заказы',
        description: 'Завершённые и отменённые заказы',
        count: 3,
        totalAmount: 50200
      }
    };

    setOrders(mockOrders[segment]);
    setSummary(mockSummaries[segment]);
    setIsLoading(false);
    setRefreshing(false);
  };

  useEffect(() => {
    loadOrders(activeSegment);
  }, [activeSegment]);

  const handleRefresh = () => {
    setRefreshing(true);
    loadOrders(activeSegment, false);
  };

  const handleSegmentChange = (segment: OrderStatus) => {
    setActiveSegment(segment);
  };

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
        <h1 className="text-2xl font-bold text-gray-900">Заказы</h1>
      </header>

      {/* Segmented Control */}
      <div className="bg-white border-b border-gray-200 px-4 py-4">
        <div className="flex space-x-1 bg-gray-100 rounded-lg p-1">
          {segments.map((segment) => (
            <button
              key={segment.id}
              onClick={() => handleSegmentChange(segment.id)}
              className={`flex-1 py-3 px-4 text-sm font-medium rounded-md transition-colors ${
                activeSegment === segment.id
                  ? 'bg-white text-gray-900 shadow-sm'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              {segment.label}
            </button>
          ))}
        </div>
      </div>

      {/* Summary Block */}
      {summary && !isLoading && (
        <div className="bg-white mx-4 mt-4 rounded-lg p-4 shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold text-gray-900 mb-1">{summary.title}</h3>
          <p className="text-sm text-gray-600 mb-2">{summary.description}</p>
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-500">
              {summary.count} {summary.count === 1 ? 'заказ' : summary.count < 5 ? 'заказа' : 'заказов'}
            </span>
            <span className="text-lg font-semibold text-gray-900">
              {summary.totalAmount.toLocaleString('ru-RU')} ₽
            </span>
          </div>
        </div>
      )}

      {/* Orders List */}
      <div className="px-4 py-4">
        {isLoading ? (
          // Loading skeleton
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
                <div className="animate-pulse flex space-x-4">
                  <div className="w-12 h-12 bg-gray-200 rounded-lg"></div>
                  <div className="flex-1 space-y-2">
                    <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                    <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                    <div className="h-3 bg-gray-200 rounded w-1/3"></div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : orders.length > 0 ? (
          <div className="space-y-3">
            {/* Pull to refresh indicator */}
            {refreshing && (
              <div className="flex justify-center py-2">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-500"></div>
              </div>
            )}
            
            {orders.map((order) => (
              <OrderItem
                key={order.id}
                order={order}
                onClick={() => console.log('Order clicked:', order.id)}
              />
            ))}
            
            {/* Manual refresh button */}
            <div className="pt-4">
              <button
                onClick={handleRefresh}
                disabled={refreshing}
                className="w-full py-3 px-4 bg-white border border-gray-200 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors disabled:opacity-50"
              >
                {refreshing ? 'Обновление...' : 'Обновить список'}
              </button>
            </div>
          </div>
        ) : (
          // Empty state
          <div className="text-center py-12">
            <svg className="w-16 h-16 text-gray-300 mx-auto mb-4" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M6 2a2 2 0 00-2 2v12a2 2 0 002 2h8a2 2 0 002-2V4a2 2 0 00-2-2H6zm1 2a1 1 0 000 2h6a1 1 0 100-2H7zm6 7a1 1 0 011 1v3a1 1 0 11-2 0v-3a1 1 0 011-1zm-3 3a1 1 0 100 2h.01a1 1 0 100-2H10zm-4 1a1 1 0 011-1h.01a1 1 0 110 2H7a1 1 0 01-1-1zm1-4a1 1 0 100 2h.01a1 1 0 100-2H7zm2 0a1 1 0 100 2h.01a1 1 0 100-2H9zm2 0a1 1 0 100 2h.01a1 1 0 100-2H11z" clipRule="evenodd" />
            </svg>
            <h3 className="text-lg font-medium text-gray-900 mb-2">Нет заказов</h3>
            <p className="text-gray-500">В этой категории пока нет заказов</p>
          </div>
        )}
      </div>

      <TabBar />
    </div>
  );
} 