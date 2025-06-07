interface Order {
  id: string;
  date: string;
  amount: number;
  status: string;
  firstProductImage?: string;
}

interface OrderItemProps {
  order: Order;
  onClick?: () => void;
}

export default function OrderItem({ order, onClick }: OrderItemProps) {
  return (
    <div 
      className="bg-white rounded-lg p-4 shadow-sm border border-gray-100 cursor-pointer hover:bg-gray-50 transition-colors"
      onClick={onClick}
    >
      <div className="flex items-center space-x-4">
        {/* Product Image */}
        <div className="w-12 h-12 bg-gray-200 rounded-lg flex-shrink-0 overflow-hidden">
          {order.firstProductImage ? (
            <img 
              src={order.firstProductImage} 
              alt="Товар" 
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-gray-400">
              <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clipRule="evenodd" />
              </svg>
            </div>
          )}
        </div>

        {/* Order Info */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between mb-1">
            <p className="text-sm font-medium text-gray-900 truncate">
              №{order.id}
            </p>
            <p className="text-sm font-semibold text-gray-900">
              {order.amount.toLocaleString('ru-RU')} ₽
            </p>
          </div>
          
          <p className="text-xs text-gray-500 mb-1">
            {new Date(order.date).toLocaleDateString('ru-RU', {
              day: '2-digit',
              month: '2-digit',
              year: 'numeric'
            })}
          </p>
          
          <p className="text-xs text-gray-600">
            {order.status}
          </p>
        </div>

        {/* Arrow */}
        <div className="flex-shrink-0">
          <svg className="w-5 h-5 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
          </svg>
        </div>
      </div>
    </div>
  );
} 