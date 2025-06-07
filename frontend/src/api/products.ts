import apiClient from './axios';
import type { 
  Product, 
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
  
  files.forEach((file, index) => {
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

// Удаление товара (если есть такой эндпоинт)
export const deleteProduct = async (productId: string): Promise<{ message: string }> => {
  const response = await apiClient.delete(`/api/products/${productId}`);
  return response.data;
}; 