import apiClient from './axios';
import type { 
  ProductDetail, 
  ProductsResponse, 
  ProductsFilters, 
  CreateProductData,
  CategoriesResponse 
} from '../types/product';

// Получение списка товаров с фильтрами
export const getProducts = async (filters?: ProductsFilters): Promise<ProductsResponse> => {
  const params = new URLSearchParams();
  
  if (filters) {
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, value.toString());
      }
    });
  }
  
  const response = await apiClient.get(`/api/products?${params.toString()}`);
  return response.data;
};

// Получение детальной информации о товаре
export const getProduct = async (productId: string): Promise<ProductDetail> => {
  const response = await apiClient.get(`/api/products/${productId}`);
  return response.data;
};

// Создание нового товара
export const createProduct = async (productData: CreateProductData): Promise<{ id: string; message: string; product: ProductDetail }> => {
  const response = await apiClient.post('/api/products', productData);
  return response.data;
};

// Обновление товара
export const updateProduct = async (productId: string, productData: Partial<CreateProductData>): Promise<{ message: string; product: ProductDetail }> => {
  const response = await apiClient.put(`/api/products/${productId}`, productData);
  return response.data;
};

// Загрузка изображений товара
export const uploadProductImages = async (
  productId: string, 
  files: File[], 
  altTexts: string[], 
  isMainFlags: boolean[]
): Promise<{ message: string; images: any[] }> => {
  const formData = new FormData();
  
  files.forEach((file) => {
    formData.append('files', file);
  });
  
  formData.append('alt_texts', JSON.stringify(altTexts));
  formData.append('is_main', JSON.stringify(isMainFlags));
  
  const response = await apiClient.post(`/api/products/${productId}/images`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  
  return response.data;
};

// Получение категорий товаров
export const getCategories = async (): Promise<CategoriesResponse> => {
  const response = await apiClient.get('/api/products/categories');
  return response.data;
};

// Получение размерной сетки
export const getSizeChart = async (category: string): Promise<any> => {
  const response = await apiClient.get(`/api/products/sizes?category=${category}`);
  return response.data;
};

// Удаление товара
export const deleteProduct = async (productId: string): Promise<{ message: string }> => {
  const response = await apiClient.delete(`/api/products/${productId}`);
  return response.data;
};

// Обновление статуса товара (активный/неактивный)
export const updateProductStatus = async (productId: string, status: 'active' | 'inactive' | 'draft'): Promise<{ message: string }> => {
  const response = await apiClient.patch(`/api/products/${productId}/status`, { status });
  return response.data;
};

// Массовое обновление товаров
export const bulkUpdateProducts = async (productIds: string[], updates: Partial<CreateProductData>): Promise<{ message: string; updated_count: number }> => {
  const response = await apiClient.patch('/api/products/bulk', { product_ids: productIds, updates });
  return response.data;
};

// Экспорт товаров в CSV/Excel
export const exportProducts = async (filters?: ProductsFilters, format: 'csv' | 'excel' = 'csv'): Promise<Blob> => {
  const params = new URLSearchParams();
  
  if (filters) {
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, value.toString());
      }
    });
  }
  
  params.append('format', format);
  
  const response = await apiClient.get(`/api/products/export?${params.toString()}`, {
    responseType: 'blob',
  });
  
  return response.data;
};

// Импорт товаров из CSV/Excel
export const importProducts = async (file: File): Promise<{ message: string; imported_count: number; errors?: string[] }> => {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await apiClient.post('/api/products/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  
  return response.data;
}; 