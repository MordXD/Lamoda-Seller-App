import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import OrderItem from '../components/OrderItem';
import TabBar from '../components/TabBar';
import { useOrders } from '../hooks/useOrders';
import type { OrdersFilters } from '../types/order';

type OrderStatus = 'new' | 'in_transit' | 'archive';

interface SummaryData {
  title: string;
  description: string;
  count: number;
  totalAmount: number;
}

export default function OrdersPage() {
  const navigate = useNavigate();
  const [activeSegment, setActiveSegment] = useState<OrderStatus>('new');
  const [refreshing, setRefreshing] = useState(false);

  // Маппинг сегментов на статусы API
  const getApiFilters = (segment: OrderStatus): OrdersFilters => {
    switch (segment) {
      case 'new':
        return { status: 'new' };
      case 'in_transit':
        return { status: 'in_transit' };
      case 'archive':
        return { status: 'delivered' }; // или можно добавить несколько статусов
      default:
        return {};
    }
  };

  const { orders, isLoading, error, summary, refetch } = useOrders(getApiFilters(activeSegment));

  const segments = [
    { id: 'new' as OrderStatus, label: 'Новые' },
    { id: 'in_transit' as OrderStatus, label: 'В пути' },
    { id: 'archive' as OrderStatus, label: 'Архив' }
  ];

  // Создаем данные для summary на основе API ответа
  const getSummaryData = (segment: OrderStatus): SummaryData => {
    const summaryTitles = {
      new: 'Новые заказы',
      in_transit: 'Заказы в пути',
      archive: 'Архивные заказы'
    };

    const summaryDescriptions = {
      new: 'Требуют обработки и подтверждения',
      in_transit: 'Находятся в процессе доставки',
      archive: 'Завершённые и отменённые заказы'
    };

    return {
      title: summaryTitles[segment],
      description: summaryDescriptions[segment],
      count: summary?.total_orders || 0,
      totalAmount: summary?.total_amount || 0
    };
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await refetch(getApiFilters(activeSegment));
    } finally {
      setRefreshing(false);
    }
  };

  const handleSegmentChange = async (segment: OrderStatus) => {
    setActiveSegment(segment);
    // Загружаем данные для нового сегмента
    await refetch(getApiFilters(segment));
  };

  const handleOrderClick = (orderId: string) => {
    navigate(`/orders/${orderId}`);
  };

  const summaryData = getSummaryData(activeSegment);

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
      {!isLoading && !error && summary && (
        <div className="bg-white mx-4 mt-4 rounded-lg p-4 shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold text-gray-900 mb-1">{summaryData.title}</h3>
          <p className="text-sm text-gray-600 mb-2">{summaryData.description}</p>
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-500">
              {summaryData.count} {summaryData.count === 1 ? 'заказ' : summaryData.count < 5 ? 'заказа' : 'заказов'}
            </span>
            <span className="text-lg font-semibold text-gray-900">
              {summaryData.totalAmount.toLocaleString('ru-RU')} ₽
            </span>
          </div>
        </div>
      )}

      {/* Error State */}
      {error && (
        <div className="mx-4 mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-600 text-sm">{error}</p>
          <button
            onClick={handleRefresh}
            className="mt-2 text-red-700 text-sm font-medium hover:text-red-800"
          >
            Попробовать снова
          </button>
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
                order={{
                  id: order.order_number,
                  date: order.date,
                  amount: order.totals.total,
                  status: order.status,
                  firstProductImage: order.items[0]?.image
                }}
                onClick={() => handleOrderClick(order.id)}
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