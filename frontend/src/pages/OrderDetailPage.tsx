import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useOrder } from '../hooks/useOrders';
import type { UpdateOrderStatusRequest } from '../types/order';

export default function OrderDetailPage() {
  const { orderId } = useParams<{ orderId: string }>();
  const navigate = useNavigate();
  const { order, isLoading, error, updateStatus } = useOrder(orderId);
  const [isUpdatingStatus, setIsUpdatingStatus] = useState(false);
  const [showStatusModal, setShowStatusModal] = useState(false);
  const [newStatus, setNewStatus] = useState<UpdateOrderStatusRequest['status']>('new');
  const [statusComment, setStatusComment] = useState('');

  const statusLabels = {
    new: 'Новый',
    confirmed: 'Подтвержден',
    in_transit: 'В пути',
    delivered: 'Доставлен',
    returned: 'Возвращен',
    cancelled: 'Отменен'
  };

  const statusColors = {
    new: 'bg-blue-100 text-blue-800',
    confirmed: 'bg-green-100 text-green-800',
    in_transit: 'bg-yellow-100 text-yellow-800',
    delivered: 'bg-green-100 text-green-800',
    returned: 'bg-red-100 text-red-800',
    cancelled: 'bg-gray-100 text-gray-800'
  };

  const handleStatusUpdate = async () => {
    if (!order) return;

    setIsUpdatingStatus(true);
    try {
      await updateStatus({
        status: newStatus,
        comment: statusComment || undefined,
      });
      setShowStatusModal(false);
      setStatusComment('');
    } catch (error) {
      console.error('Ошибка обновления статуса:', error);
      alert('Не удалось обновить статус заказа');
    } finally {
      setIsUpdatingStatus(false);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 pb-20">
        <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
          <div className="flex items-center">
            <button onClick={() => navigate(-1)} className="mr-4">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <h1 className="text-xl font-bold text-gray-900">Загрузка...</h1>
          </div>
        </header>
        
        <div className="px-4 py-6">
          <div className="animate-pulse space-y-4">
            <div className="h-32 bg-gray-200 rounded-lg"></div>
            <div className="h-24 bg-gray-200 rounded-lg"></div>
            <div className="h-40 bg-gray-200 rounded-lg"></div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !order) {
    return (
      <div className="min-h-screen bg-gray-50 pb-20">
        <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
          <div className="flex items-center">
            <button onClick={() => navigate(-1)} className="mr-4">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <h1 className="text-xl font-bold text-gray-900">Ошибка</h1>
          </div>
        </header>
        
        <div className="px-4 py-6">
          <div className="text-center">
            <p className="text-red-600 mb-4">{error || 'Заказ не найден'}</p>
            <button
              onClick={() => navigate(-1)}
              className="bg-black text-white px-4 py-2 rounded-lg"
            >
              Назад
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <button onClick={() => navigate(-1)} className="mr-4">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <div>
              <h1 className="text-xl font-bold text-gray-900">Заказ {order.order_number}</h1>
              <p className="text-sm text-gray-500">
                {new Date(order.date).toLocaleDateString('ru-RU')}
              </p>
            </div>
          </div>
          
          <button
            onClick={() => setShowStatusModal(true)}
            className="bg-black text-white px-4 py-2 rounded-lg text-sm"
          >
            Изменить статус
          </button>
        </div>
      </header>

      <div className="px-4 py-6 space-y-6">
        {/* Status and Summary */}
        <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
          <div className="flex items-center justify-between mb-4">
            <span className={`px-3 py-1 rounded-full text-sm font-medium ${statusColors[order.status]}`}>
              {statusLabels[order.status]}
            </span>
            <span className="text-2xl font-bold text-gray-900">
              {order.totals.total.toLocaleString('ru-RU')} ₽
            </span>
          </div>
          
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-gray-500">Товары</p>
              <p className="font-semibold">{order.totals.subtotal.toLocaleString('ru-RU')} ₽</p>
            </div>
            <div>
              <p className="text-gray-500">Доставка</p>
              <p className="font-semibold">{order.totals.delivery.toLocaleString('ru-RU')} ₽</p>
            </div>
            <div>
              <p className="text-gray-500">Скидка</p>
              <p className="font-semibold text-red-600">-{order.totals.discount.toLocaleString('ru-RU')} ₽</p>
            </div>
            <div>
              <p className="text-gray-500">К оплате</p>
              <p className="font-semibold">{order.payment.seller_amount.toLocaleString('ru-RU')} ₽</p>
            </div>
          </div>
        </div>

        {/* Customer Info */}
        <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold text-gray-900 mb-3">Покупатель</h3>
          <div className="space-y-2">
            <p><span className="text-gray-500">Имя:</span> {order.customer.name}</p>
            <p><span className="text-gray-500">Email:</span> {order.customer.email}</p>
            <p><span className="text-gray-500">Телефон:</span> {order.customer.phone}</p>
            {order.customer.is_regular && (
              <span className="inline-block bg-green-100 text-green-800 text-xs px-2 py-1 rounded-full">
                Постоянный клиент ({order.customer.orders_count} заказов)
              </span>
            )}
          </div>
        </div>

        {/* Delivery Info */}
        <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold text-gray-900 mb-3">Доставка</h3>
          <div className="space-y-2">
            <p><span className="text-gray-500">Тип:</span> {order.delivery.type}</p>
            <p><span className="text-gray-500">Адрес:</span> {order.delivery.address.city}, {order.delivery.address.street}, {order.delivery.address.house}</p>
            <p><span className="text-gray-500">Ожидаемая дата:</span> {new Date(order.delivery.estimated_date).toLocaleDateString('ru-RU')}</p>
            {order.delivery.tracking_number && (
              <p><span className="text-gray-500">Трек-номер:</span> {order.delivery.tracking_number}</p>
            )}
          </div>
        </div>

        {/* Order Items */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-100">
          <div className="p-4 border-b border-gray-100">
            <h3 className="text-lg font-semibold text-gray-900">Товары ({order.items.length})</h3>
          </div>
          <div className="divide-y divide-gray-100">
            {order.items.map((item) => (
              <div key={item.id} className="p-4">
                <div className="flex items-center space-x-4">
                  <div className="w-16 h-16 bg-gray-200 rounded-lg overflow-hidden flex-shrink-0">
                    {item.image ? (
                      <img src={item.image} alt={item.name} className="w-full h-full object-cover" />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center text-gray-400">
                        <svg className="w-8 h-8" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clipRule="evenodd" />
                        </svg>
                      </div>
                    )}
                  </div>
                  
                  <div className="flex-1 min-w-0">
                    <h4 className="text-sm font-medium text-gray-900 truncate">{item.name}</h4>
                    <p className="text-xs text-gray-500">{item.brand} • {item.sku}</p>
                    <p className="text-xs text-gray-500">Размер: {item.size} • Цвет: {item.color}</p>
                    <p className="text-sm text-gray-900 mt-1">
                      {item.quantity} × {item.price.toLocaleString('ru-RU')} ₽ = {item.total.toLocaleString('ru-RU')} ₽
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Status History */}
        {order.status_history && order.status_history.length > 0 && (
          <div className="bg-white rounded-lg shadow-sm border border-gray-100">
            <div className="p-4 border-b border-gray-100">
              <h3 className="text-lg font-semibold text-gray-900">История статусов</h3>
            </div>
            <div className="divide-y divide-gray-100">
              {order.status_history.map((history, index) => (
                <div key={index} className="p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium text-gray-900">{statusLabels[history.status as keyof typeof statusLabels] || history.status}</p>
                      {history.comment && (
                        <p className="text-xs text-gray-500 mt-1">{history.comment}</p>
                      )}
                    </div>
                    <p className="text-xs text-gray-500">
                      {new Date(history.date).toLocaleString('ru-RU')}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Status Update Modal */}
      {showStatusModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Изменить статус заказа</h3>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Новый статус
                </label>
                <select
                  value={newStatus}
                  onChange={(e) => setNewStatus(e.target.value as UpdateOrderStatusRequest['status'])}
                  className="w-full border border-gray-300 rounded-md px-3 py-2"
                >
                  <option value="new">Новый</option>
                  <option value="confirmed">Подтвержден</option>
                  <option value="in_transit">В пути</option>
                  <option value="delivered">Доставлен</option>
                  <option value="returned">Возвращен</option>
                  <option value="cancelled">Отменен</option>
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Комментарий (опционально)
                </label>
                <textarea
                  value={statusComment}
                  onChange={(e) => setStatusComment(e.target.value)}
                  className="w-full border border-gray-300 rounded-md px-3 py-2"
                  rows={3}
                  placeholder="Добавьте комментарий к изменению статуса..."
                />
              </div>
            </div>
            
            <div className="flex space-x-3 mt-6">
              <button
                onClick={() => setShowStatusModal(false)}
                className="flex-1 bg-gray-200 text-gray-800 py-2 px-4 rounded-lg"
                disabled={isUpdatingStatus}
              >
                Отмена
              </button>
              <button
                onClick={handleStatusUpdate}
                className="flex-1 bg-black text-white py-2 px-4 rounded-lg disabled:opacity-50"
                disabled={isUpdatingStatus}
              >
                {isUpdatingStatus ? 'Обновление...' : 'Обновить'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
} 