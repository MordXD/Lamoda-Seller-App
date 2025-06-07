import { useState, useEffect } from 'react';
import type { ProductsFilters } from '../types/product';

interface ProductFiltersProps {
  filters: ProductsFilters;
  onFiltersChange: (filters: ProductsFilters) => void;
  availableFilters?: {
    categories: Array<{ id: string; name: string; count: number }>;
    brands: Array<{ id: string; name: string; count: number }>;
    price_range: { min: number; max: number };
  };
}

export default function ProductFilters({ filters, onFiltersChange, availableFilters }: ProductFiltersProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [localFilters, setLocalFilters] = useState<ProductsFilters>(filters);

  useEffect(() => {
    setLocalFilters(filters);
  }, [filters]);

  const handleFilterChange = (key: keyof ProductsFilters, value: any) => {
    const newFilters = { ...localFilters, [key]: value };
    setLocalFilters(newFilters);
  };

  const applyFilters = () => {
    onFiltersChange(localFilters);
    setIsOpen(false);
  };

  const clearFilters = () => {
    const clearedFilters: ProductsFilters = {};
    setLocalFilters(clearedFilters);
    onFiltersChange(clearedFilters);
    setIsOpen(false);
  };

  const hasActiveFilters = Object.keys(filters).some(key => 
    filters[key as keyof ProductsFilters] !== undefined && 
    filters[key as keyof ProductsFilters] !== ''
  );

  return (
    <div className="relative">
      {/* Filter Button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className={`flex items-center space-x-2 px-4 py-2 rounded-lg border transition-colors ${
          hasActiveFilters 
            ? 'bg-gray-900 text-white border-gray-900' 
            : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
        }`}
      >
        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
          <path fillRule="evenodd" d="M3 3a1 1 0 011-1h12a1 1 0 011 1v3a1 1 0 01-.293.707L12 11.414V15a1 1 0 01-.293.707l-2 2A1 1 0 018 17v-5.586L3.293 6.707A1 1 0 013 6V3z" clipRule="evenodd" />
        </svg>
        <span>Фильтры</span>
        {hasActiveFilters && (
          <span className="bg-white text-gray-900 text-xs px-1.5 py-0.5 rounded-full">
            {Object.keys(filters).filter(key => filters[key as keyof ProductsFilters]).length}
          </span>
        )}
      </button>

      {/* Filter Panel */}
      {isOpen && (
        <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-lg shadow-lg z-20 p-4 space-y-4">
          {/* Category Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Категория
            </label>
            <select
              value={localFilters.category || ''}
              onChange={(e) => handleFilterChange('category', e.target.value || undefined)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500"
            >
              <option value="">Все категории</option>
              {availableFilters?.categories.map((category) => (
                <option key={category.id} value={category.id}>
                  {category.name} ({category.count})
                </option>
              ))}
            </select>
          </div>

          {/* Brand Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Бренд
            </label>
            <select
              value={localFilters.brand || ''}
              onChange={(e) => handleFilterChange('brand', e.target.value || undefined)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500"
            >
              <option value="">Все бренды</option>
              {availableFilters?.brands.map((brand) => (
                <option key={brand.id} value={brand.id}>
                  {brand.name} ({brand.count})
                </option>
              ))}
            </select>
          </div>

          {/* Price Range */}
          {availableFilters?.price_range && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Цена
              </label>
              <div className="grid grid-cols-2 gap-2">
                <input
                  type="number"
                  placeholder={`От ${availableFilters.price_range.min}`}
                  value={localFilters.min_price || ''}
                  onChange={(e) => handleFilterChange('min_price', e.target.value ? Number(e.target.value) : undefined)}
                  className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500"
                />
                <input
                  type="number"
                  placeholder={`До ${availableFilters.price_range.max}`}
                  value={localFilters.max_price || ''}
                  onChange={(e) => handleFilterChange('max_price', e.target.value ? Number(e.target.value) : undefined)}
                  className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500"
                />
              </div>
            </div>
          )}

          {/* Stock Status */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Наличие
            </label>
            <select
              value={localFilters.stock_status || ''}
              onChange={(e) => handleFilterChange('stock_status', e.target.value || undefined)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500"
            >
              <option value="">Все товары</option>
              <option value="in_stock">В наличии</option>
              <option value="low_stock">Заканчивается</option>
              <option value="out_of_stock">Нет в наличии</option>
            </select>
          </div>

          {/* Sort */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Сортировка
            </label>
            <div className="grid grid-cols-2 gap-2">
              <select
                value={localFilters.sort_by || ''}
                onChange={(e) => handleFilterChange('sort_by', e.target.value || undefined)}
                className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500"
              >
                <option value="">По умолчанию</option>
                <option value="name">По названию</option>
                <option value="price">По цене</option>
                <option value="stock">По остаткам</option>
                <option value="sales">По продажам</option>
                <option value="created_date">По дате</option>
              </select>
              <select
                value={localFilters.sort_order || ''}
                onChange={(e) => handleFilterChange('sort_order', e.target.value || undefined)}
                className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500"
              >
                <option value="asc">По возрастанию</option>
                <option value="desc">По убыванию</option>
              </select>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex space-x-2 pt-2">
            <button
              onClick={applyFilters}
              className="flex-1 px-4 py-2 bg-gray-900 text-white rounded-lg hover:bg-gray-800 transition-colors"
            >
              Применить
            </button>
            <button
              onClick={clearFilters}
              className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Сбросить
            </button>
          </div>
        </div>
      )}

      {/* Overlay */}
      {isOpen && (
        <div 
          className="fixed inset-0 bg-black bg-opacity-25 z-10"
          onClick={() => setIsOpen(false)}
        />
      )}
    </div>
  );
} 