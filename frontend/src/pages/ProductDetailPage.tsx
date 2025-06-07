import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useProduct } from '../hooks/useProducts';

export default function ProductDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { product, isLoading, error } = useProduct(id);
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="px-4 py-6">
          <div className="animate-pulse">
            <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
            <div className="h-64 bg-gray-200 rounded-lg mb-6"></div>
            <div className="space-y-4">
              <div className="h-6 bg-gray-200 rounded w-3/4"></div>
              <div className="h-4 bg-gray-200 rounded w-1/2"></div>
              <div className="h-4 bg-gray-200 rounded w-2/3"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-gray-900 mb-2">
            {error || 'Товар не найден'}
          </h2>
          <button
            onClick={() => navigate('/products')}
            className="px-4 py-2 bg-gray-900 text-white rounded-lg hover:bg-gray-800 transition-colors"
          >
            Вернуться к товарам
          </button>
        </div>
      </div>
    );
  }

  const currentImage = product.images?.[selectedImageIndex] || null;
  const isLowStock = product.total_available < 10;
  const isOutOfStock = product.total_available === 0;

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
        <div className="flex items-center">
          <button
            onClick={() => navigate('/products')}
            className="mr-3 p-2 -ml-2 rounded-lg hover:bg-gray-100"
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clipRule="evenodd" />
            </svg>
          </button>
          <h1 className="text-xl font-semibold text-gray-900 truncate">
            {product.name}
          </h1>
        </div>
      </header>

      <div className="px-4 py-6 space-y-6">
        {/* Product Images */}
        <div className="bg-white rounded-lg p-4">
          <div className="aspect-square bg-gray-100 rounded-lg overflow-hidden mb-4">
            {currentImage ? (
              <img
                src={currentImage.url}
                alt={currentImage.alt}
                className="w-full h-full object-cover"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-gray-400">
                <svg className="w-16 h-16" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clipRule="evenodd" />
                </svg>
              </div>
            )}
          </div>
          
          {/* Image thumbnails */}
          {product.images && product.images.length > 1 && (
            <div className="flex space-x-2 overflow-x-auto">
              {product.images.map((image, index) => (
                <button
                  key={image.id}
                  onClick={() => setSelectedImageIndex(index)}
                  className={`flex-shrink-0 w-16 h-16 rounded-lg overflow-hidden border-2 ${
                    index === selectedImageIndex ? 'border-gray-900' : 'border-gray-200'
                  }`}
                >
                  <img
                    src={image.url}
                    alt={image.alt}
                    className="w-full h-full object-cover"
                  />
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Product Info */}
        <div className="bg-white rounded-lg p-4 space-y-4">
          <div>
            <div className="flex items-center gap-2 mb-2">
              <h2 className="text-xl font-semibold text-gray-900">{product.name}</h2>
              {product.is_bestseller && (
                <span className="text-xs bg-yellow-100 text-yellow-800 px-2 py-1 rounded-full">
                  ХИТ
                </span>
              )}
              {product.is_new && (
                <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded-full">
                  NEW
                </span>
              )}
            </div>
            <p className="text-gray-600">{product.brand} • {product.sku}</p>
          </div>

          <div className="flex items-center justify-between">
            <div>
              <p className="text-2xl font-bold text-gray-900">
                {product.price.toLocaleString('ru-RU')} ₽
              </p>
              <p className="text-sm text-gray-500">
                Себестоимость: {product.cost_price.toLocaleString('ru-RU')} ₽
              </p>
              <p className="text-sm text-green-600">
                Маржа: {product.margin_percent.toFixed(1)}%
              </p>
            </div>
            
            <div className="text-right">
              <span className={`inline-block px-3 py-1 rounded-full text-sm font-medium ${
                isOutOfStock 
                  ? 'bg-red-100 text-red-700'
                  : isLowStock 
                    ? 'bg-yellow-100 text-yellow-700'
                    : 'bg-green-100 text-green-700'
              }`}>
                {isOutOfStock ? 'Нет в наличии' : `${product.total_available} шт`}
              </span>
              <p className="text-xs text-gray-500 mt-1">
                Всего: {product.total_stock} шт
              </p>
            </div>
          </div>

          {product.description && (
            <div>
              <h3 className="font-medium text-gray-900 mb-2">Описание</h3>
              <p className="text-gray-600 text-sm">{product.description}</p>
            </div>
          )}
        </div>

        {/* Sales Statistics */}
        <div className="bg-white rounded-lg p-4">
          <h3 className="font-medium text-gray-900 mb-4">Статистика продаж</h3>
          <div className="grid grid-cols-2 gap-4">
            <div className="text-center p-3 bg-gray-50 rounded-lg">
              <p className="text-2xl font-bold text-gray-900">{product.sales_stats.sales_30d}</p>
              <p className="text-sm text-gray-600">Продано за месяц</p>
            </div>
            <div className="text-center p-3 bg-gray-50 rounded-lg">
              <p className="text-2xl font-bold text-gray-900">
                {product.sales_stats.revenue_30d.toLocaleString('ru-RU')} ₽
              </p>
              <p className="text-sm text-gray-600">Выручка за месяц</p>
            </div>
            <div className="text-center p-3 bg-gray-50 rounded-lg">
              <p className="text-2xl font-bold text-gray-900">{product.rating}</p>
              <p className="text-sm text-gray-600">Рейтинг</p>
            </div>
            <div className="text-center p-3 bg-gray-50 rounded-lg">
              <p className="text-2xl font-bold text-gray-900">{product.return_rate}%</p>
              <p className="text-sm text-gray-600">Возвраты</p>
            </div>
          </div>
        </div>

        {/* Variants */}
        {product.variants && product.variants.length > 0 && (
          <div className="bg-white rounded-lg p-4">
            <h3 className="font-medium text-gray-900 mb-4">Варианты товара</h3>
            <div className="space-y-3">
              {product.variants.map((variant) => (
                <div key={variant.id} className="flex items-center justify-between p-3 border border-gray-200 rounded-lg">
                  <div>
                    <p className="font-medium text-gray-900">
                      {variant.size} • {variant.color}
                    </p>
                    <p className="text-sm text-gray-600">{variant.sku}</p>
                  </div>
                  <div className="text-right">
                    <p className="font-medium text-gray-900">{variant.available} шт</p>
                    <p className="text-sm text-gray-500">из {variant.stock}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Product Details */}
        <div className="bg-white rounded-lg p-4">
          <h3 className="font-medium text-gray-900 mb-4">Детали товара</h3>
          <div className="space-y-3 text-sm">
            {product.material && (
              <div className="flex justify-between">
                <span className="text-gray-600">Материал:</span>
                <span className="text-gray-900">{product.material}</span>
              </div>
            )}
            {product.care_instructions && (
              <div className="flex justify-between">
                <span className="text-gray-600">Уход:</span>
                <span className="text-gray-900">{product.care_instructions}</span>
              </div>
            )}
            {product.country_origin && (
              <div className="flex justify-between">
                <span className="text-gray-600">Страна:</span>
                <span className="text-gray-900">{product.country_origin}</span>
              </div>
            )}
            <div className="flex justify-between">
              <span className="text-gray-600">Категория:</span>
              <span className="text-gray-900">{product.category}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">Подкатегория:</span>
              <span className="text-gray-900">{product.subcategory}</span>
            </div>
          </div>
        </div>

        {/* Tags */}
        {product.tags && product.tags.length > 0 && (
          <div className="bg-white rounded-lg p-4">
            <h3 className="font-medium text-gray-900 mb-3">Теги</h3>
            <div className="flex flex-wrap gap-2">
              {product.tags.map((tag, index) => (
                <span
                  key={index}
                  className="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded-full"
                >
                  {tag}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
} 