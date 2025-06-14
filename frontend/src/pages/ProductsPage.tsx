import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import ProductItem from '../components/ProductItem';
import ProductFilters from '../components/ProductFilters';
import TabBar from '../components/TabBar';
import { useProducts } from '../hooks/useProducts';
import type { ProductsFilters } from '../types/product';

export default function ProductsPage() {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [currentFilters, setCurrentFilters] = useState<ProductsFilters>({});
  const { products, isLoading, error, pagination, filters: availableFilters, refetch } = useProducts();
  const [filteredProducts, setFilteredProducts] = useState(products);

  useEffect(() => {
    if (searchQuery.trim()) {
      const filtered = products.filter(product =>
        product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        product.brand.toLowerCase().includes(searchQuery.toLowerCase()) ||
        product.sku.toLowerCase().includes(searchQuery.toLowerCase())
      );
      setFilteredProducts(filtered);
    } else {
      setFilteredProducts(products);
    }
  }, [searchQuery, products]);

  const handleProductClick = (productId: string) => {
    navigate(`/products/${productId}`);
  };

  const handleFiltersChange = async (newFilters: ProductsFilters) => {
    setCurrentFilters(newFilters);
    await refetch(newFilters);
  };

  const handleSearch = async () => {
    if (searchQuery.trim()) {
      const searchFilters = { ...currentFilters, search: searchQuery.trim() };
      await refetch(searchFilters);
    } else {
      await refetch(currentFilters);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">Товары</h1>
        
        {/* Search */}
        <div className="relative mb-4">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg className="h-5 w-5 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clipRule="evenodd" />
            </svg>
          </div>
          <input
            type="text"
            placeholder="Поиск по названию, бренду или артикулу"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
            className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-transparent"
          />
          {searchQuery && (
            <button
              onClick={() => {
                setSearchQuery('');
                refetch(currentFilters);
              }}
              className="absolute inset-y-0 right-0 pr-3 flex items-center"
            >
              <svg className="h-5 w-5 text-gray-400 hover:text-gray-600" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
              </svg>
            </button>
          )}
        </div>

        {/* Filters */}
        <div className="flex items-center justify-between">
          <ProductFilters
            filters={currentFilters}
            onFiltersChange={handleFiltersChange}
            availableFilters={availableFilters || undefined}
          />
          
          {/* Results count */}
          {pagination && (
            <span className="text-sm text-gray-600">
              {pagination.total} {
                pagination.total === 1 ? 'товар' : 
                pagination.total < 5 ? 'товара' : 'товаров'
              }
            </span>
          )}
        </div>
        
        {/* Error message */}
        {error && (
          <div className="mt-3 bg-red-50 border border-red-200 rounded-lg p-3">
            <p className="text-sm text-red-800">{error}</p>
          </div>
        )}
      </header>

      {/* Content */}
      <div className="px-4 py-6">
        {isLoading ? (
          // Loading skeleton
          <div className="space-y-4">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
                <div className="animate-pulse flex space-x-4">
                  <div className="w-16 h-16 bg-gray-200 rounded-lg"></div>
                  <div className="flex-1 space-y-2">
                    <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                    <div className="h-4 bg-gray-200 rounded w-1/2"></div>
                    <div className="h-3 bg-gray-200 rounded w-1/3"></div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : filteredProducts.length > 0 ? (
          <div className="space-y-3">
            {/* Results summary */}
            {searchQuery && (
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 mb-4">
                <p className="text-sm text-blue-800">
                  Найдено {filteredProducts.length} {
                    filteredProducts.length === 1 ? 'товар' : 
                    filteredProducts.length < 5 ? 'товара' : 'товаров'
                  } по запросу "{searchQuery}"
                </p>
              </div>
            )}
            
            {filteredProducts.map((product) => (
              <ProductItem
                key={product.id}
                product={product}
                onClick={() => handleProductClick(product.id)}
              />
            ))}
          </div>
        ) : searchQuery ? (
          // No search results
          <div className="text-center py-12">
            <svg className="w-16 h-16 text-gray-300 mx-auto mb-4" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clipRule="evenodd" />
            </svg>
            <h3 className="text-lg font-medium text-gray-900 mb-2">Ничего не найдено</h3>
            <p className="text-gray-500 mb-4">Попробуйте изменить запрос поиска</p>
            <button
              onClick={() => {
                setSearchQuery('');
                refetch(currentFilters);
              }}
              className="px-4 py-2 bg-gray-900 text-white rounded-lg hover:bg-gray-800 transition-colors"
            >
              Показать все товары
            </button>
          </div>
        ) : (
          // Empty state - no products at all
          <div className="text-center py-12">
            <svg className="w-16 h-16 text-gray-300 mx-auto mb-4" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 2L3 7v11a1 1 0 001 1h12a1 1 0 001-1V7l-7-5zM8 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm4 0a1 1 0 112 0v6a1 1 0 11-2 0V8z" clipRule="evenodd" />
            </svg>
            <h3 className="text-lg font-medium text-gray-900 mb-2">У вас пока нет товаров</h3>
            <p className="text-gray-500 mb-4">Добавьте первый товар для начала работы</p>
            <button className="px-4 py-2 bg-gray-900 text-white rounded-lg hover:bg-gray-800 transition-colors">
              Добавить товар
            </button>
          </div>
        )}
      </div>

      <TabBar />
    </div>
  );
} 