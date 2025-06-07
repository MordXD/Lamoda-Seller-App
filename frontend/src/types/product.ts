export interface Product {
  id: string;
  name: string;
  brand: string;
  category: string;
  subcategory: string;
  sku: string;
  price: number;
  cost_price: number;
  margin_percent: number;
  currency: string;
  main_image?: string;
  image?: string;
  total_stock: number;
  stock: number;
  available_sizes: string[];
  available_colors: string[];
  sales_count_30d: number;
  revenue_30d: number;
  rating: number;
  reviews_count: number;
  return_rate: number;
  created_date: string;
  updated_date: string;
  status: 'active' | 'inactive' | 'draft';
  seasonal_demand: string;
  is_bestseller: boolean;
  is_new: boolean;
  discount_percent: number;
}

export interface ProductVariant {
  id: string;
  sku: string;
  size: string;
  color: string;
  color_hex: string;
  stock: number;
  reserved: number;
  available: number;
  price: number;
  weight: number;
  dimensions: {
    length: number;
    width: number;
    height: number;
  };
}

export interface ProductImage {
  id: string;
  url: string;
  alt: string;
  is_main: boolean;
  order: number;
  size_bytes?: number;
  dimensions?: {
    width: number;
    height: number;
  };
}

export interface ProductDetail extends Product {
  description: string;
  barcode: string;
  images: ProductImage[];
  variants: ProductVariant[];
  total_reserved: number;
  total_available: number;
  sales_stats: {
    total_sold: number;
    revenue_total: number;
    sales_30d: number;
    revenue_30d: number;
    sales_7d: number;
    revenue_7d: number;
    avg_daily_sales: number;
    peak_sales_month: string;
  };
  return_reasons: Array<{
    reason: string;
    count: number;
    percent: number;
  }>;
  tags: string[];
  material: string;
  care_instructions: string;
  country_origin: string;
  supplier: {
    id: string;
    name: string;
    contact: string;
  };
}

export interface ProductsResponse {
  products: Product[];
  pagination: {
    total: number;
    limit: number;
    offset: number;
    has_next: boolean;
    has_prev: boolean;
  };
  filters: {
    categories: Array<{
      id: string;
      name: string;
      count: number;
    }>;
    brands: Array<{
      id: string;
      name: string;
      count: number;
    }>;
    price_range: {
      min: number;
      max: number;
    };
  };
}

export interface ProductsFilters {
  search?: string;
  category?: string;
  brand?: string;
  min_price?: number;
  max_price?: number;
  stock_status?: 'in_stock' | 'low_stock' | 'out_of_stock';
  sort_by?: 'name' | 'price' | 'stock' | 'sales' | 'created_date';
  sort_order?: 'asc' | 'desc';
  limit?: number;
  offset?: number;
}

export interface CreateProductData {
  name: string;
  description: string;
  brand: string;
  category: string;
  subcategory: string;
  sku: string;
  barcode?: string;
  price: number;
  cost_price: number;
  currency: string;
  material?: string;
  care_instructions?: string;
  country_origin?: string;
  tags?: string[];
  variants: Array<{
    sku: string;
    size: string;
    color: string;
    color_hex: string;
    stock: number;
    weight: number;
    dimensions: {
      length: number;
      width: number;
      height: number;
    };
  }>;
  seasonal_demand?: string;
  supplier_id?: string;
}

export interface Category {
  id: string;
  name: string;
  subcategories?: Category[];
}

export interface CategoriesResponse {
  categories: Category[];
} 