export interface DashboardStats {
  revenue: {
    today: number;
    yesterday: number;
    week: number;
    month: number;
    change_percent: number;
  };
  orders: {
    today: number;
    yesterday: number;
    week: number;
    month: number;
    change_count: number;
  };
  products: {
    total_active: number;
    low_stock: number;
    out_of_stock: number;
    bestsellers: number;
  };
  customers: {
    new_today: number;
    returning: number;
    total_active: number;
  };
}

export interface SalesData {
  date: string;
  orders: number;
  purchases: number;
  revenue: number;
}

export interface TopProduct {
  id: string;
  name: string;
  brand: string;
  image?: string;
  sales_count: number;
  revenue: number;
  margin: number;
  return_rate: number;
}

export interface CategoryPerformance {
  category: string;
  sales_count: number;
  revenue: number;
  avg_price: number;
  growth_percent: number;
}

export interface DashboardResponse {
  stats: DashboardStats;
  sales_chart: SalesData[];
  top_products: TopProduct[];
  category_performance: CategoryPerformance[];
  recent_orders: Array<{
    id: string;
    order_number: string;
    customer_name: string;
    amount: number;
    status: string;
    date: string;
  }>;
}

export interface AnalyticsFilters {
  period?: 'today' | 'yesterday' | 'week' | 'month' | 'quarter' | 'year';
  date_from?: string;
  date_to?: string;
  compare_with_previous?: boolean;
}

export interface SizeChart {
  category: string;
  size_chart: {
    type: 'clothing' | 'shoes' | 'accessories';
    sizes: Array<{
      size: string;
      measurements: {
        [key: string]: string;
      };
      international?: string;
      us?: string;
    }>;
  };
}

export interface SizeChartResponse {
  category: string;
  size_chart: {
    type: 'clothing' | 'shoes' | 'accessories';
    sizes: Array<{
      size: string;
      measurements: {
        [key: string]: string;
      };
      international?: string;
      us?: string;
    }>;
  };
} 