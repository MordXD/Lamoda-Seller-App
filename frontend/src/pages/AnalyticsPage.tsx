import { useState, useEffect } from 'react';
import TabBar from '../components/TabBar';

interface TopProduct {
  id: string;
  name: string;
  salesCount: number;
  revenue: number;
  image?: string;
}

export default function AnalyticsPage() {
  const [topProducts, setTopProducts] = useState<TopProduct[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Имитация загрузки данных
    setTimeout(() => {
      setTopProducts([
        {
          id: '1',
          name: 'Пальто шерстяное классическое',
          salesCount: 24,
          revenue: 600000
        },
        {
          id: '2',
          name: 'Свитер кашемировый оверсайз',
          salesCount: 18,
          revenue: 333000
        },
        {
          id: '3',
          name: 'Джинсы широкие с высокой посадкой',
          salesCount: 35,
          revenue: 311500
        },
        {
          id: '4',
          name: 'Сапоги кожаные на каблуке',
          salesCount: 22,
          revenue: 264000
        },
        {
          id: '5',
          name: 'Куртка пуховая зимняя',
          salesCount: 16,
          revenue: 254400
        }
      ]);
      setIsLoading(false);
    }, 800);
  }, []);

  const totalRevenue = topProducts.reduce((sum, product) => sum + product.revenue, 0);
  const totalSales = topProducts.reduce((sum, product) => sum + product.salesCount, 0);

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
        <h1 className="text-2xl font-bold text-gray-900">Аналитика</h1>
        <p className="text-sm text-gray-500 mt-1">За последние 7 дней</p>
      </header>

      {/* Content */}
      <div className="px-4 py-6 space-y-6">
        {/* Summary Cards */}
        {!isLoading && (
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
              <h3 className="text-sm font-medium text-gray-500 mb-1">Общая выручка</h3>
              <p className="text-2xl font-bold text-gray-900">
                {totalRevenue.toLocaleString('ru-RU')} ₽
              </p>
            </div>
            <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
              <h3 className="text-sm font-medium text-gray-500 mb-1">Продано товаров</h3>
              <p className="text-2xl font-bold text-gray-900">
                {totalSales} шт
              </p>
            </div>
          </div>
        )}

        {/* Top Products Widget */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-100">
          <div className="p-4 border-b border-gray-100">
            <h2 className="text-lg font-semibold text-gray-900">Топ-5 товаров</h2>
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
          ) : topProducts.length > 0 ? (
            <div className="divide-y divide-gray-100">
              {topProducts.map((product, index) => (
                <div key={product.id} className="p-4 hover:bg-gray-50 transition-colors">
                  <div className="flex items-center space-x-4">
                    {/* Rank */}
                    <div className="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center flex-shrink-0">
                      <span className="text-sm font-semibold text-gray-600">
                        {index + 1}
                      </span>
                    </div>

                    {/* Product Image */}
                    <div className="w-12 h-12 bg-gray-200 rounded-lg flex-shrink-0 overflow-hidden">
                      {product.image ? (
                        <img 
                          src={product.image} 
                          alt={product.name} 
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

                    {/* Product Info */}
                    <div className="flex-1 min-w-0">
                      <h3 className="text-sm font-medium text-gray-900 truncate">
                        {product.name}
                      </h3>
                      <p className="text-xs text-gray-500">
                        {product.salesCount} продаж
                      </p>
                    </div>

                    {/* Revenue */}
                    <div className="text-right flex-shrink-0">
                      <p className="text-sm font-semibold text-gray-900">
                        {product.revenue.toLocaleString('ru-RU')} ₽
                      </p>
                      <p className="text-xs text-gray-500">
                        {Math.round((product.revenue / totalRevenue) * 100)}% от общей
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

        {/* Additional Analytics Cards */}
        <div className="grid grid-cols-1 gap-4">
          <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
            <h3 className="text-lg font-semibold text-gray-900 mb-3">Ключевые метрики</h3>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-600">Средний чек</span>
                <span className="text-sm font-semibold text-gray-900">
                  {totalSales > 0 ? Math.round(totalRevenue / totalSales).toLocaleString('ru-RU') : 0} ₽
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-600">Конверсия в покупку</span>
                <span className="text-sm font-semibold text-gray-900">12.5%</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-600">Повторные покупки</span>
                <span className="text-sm font-semibold text-gray-900">28%</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <TabBar />
    </div>
  );
} 