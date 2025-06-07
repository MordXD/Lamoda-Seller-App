import { useState, useEffect, useCallback } from 'react';
import { getProducts, getProduct } from '../api/products';
import type { Product, ProductDetail, ProductsFilters, ProductsResponse } from '../types/product';

export const useProducts = (initialFilters?: ProductsFilters) => {
  const [products, setProducts] = useState<Product[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState<ProductsResponse['pagination'] | null>(null);
  const [filters, setFilters] = useState<ProductsResponse['filters'] | null>(null);

  const loadProducts = useCallback(async (newFilters?: ProductsFilters) => {
    try {
      setIsLoading(true);
      setError(null);
      
      const response = await getProducts(newFilters || initialFilters);
      
      // Преобразуем данные для совместимости
      const transformedProducts = response.products.map(product => ({
        ...product,
        stock: product.total_stock,
        image: product.main_image,
      }));
      
      setProducts(transformedProducts);
      setPagination(response.pagination);
      setFilters(response.filters);
    } catch (err) {
      console.error('Ошибка загрузки товаров:', err);
      setError('Не удалось загрузить товары');
    } finally {
      setIsLoading(false);
    }
  }, [initialFilters]);

  useEffect(() => {
    loadProducts();
  }, [loadProducts]);

  const refetch = useCallback((newFilters?: ProductsFilters) => {
    return loadProducts(newFilters);
  }, [loadProducts]);

  return {
    products,
    isLoading,
    error,
    pagination,
    filters,
    refetch,
  };
};

export const useProduct = (productId: string | undefined) => {
  const [product, setProduct] = useState<ProductDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadProduct = useCallback(async () => {
    if (!productId) {
      setError('ID товара не указан');
      setIsLoading(false);
      return;
    }

    try {
      setIsLoading(true);
      setError(null);
      
      const productData = await getProduct(productId);
      setProduct(productData);
    } catch (err) {
      console.error('Ошибка загрузки товара:', err);
      setError('Не удалось загрузить информацию о товаре');
    } finally {
      setIsLoading(false);
    }
  }, [productId]);

  useEffect(() => {
    loadProduct();
  }, [loadProduct]);

  const refetch = useCallback(() => {
    return loadProduct();
  }, [loadProduct]);

  return {
    product,
    isLoading,
    error,
    refetch,
  };
}; 