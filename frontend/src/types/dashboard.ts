export interface Metric {
  current: number;
  previous: number;
  change_percent: number;
  change_absolute: number;
  trend: string; // up, down, stable
}

export interface TopCategory {
  category: string;
  name: string;
  revenue: number;
  orders: number;
  items: number;
}

export interface HourlySale {
  hour: number;
  revenue: number;
  orders: number;
}

export interface PeriodInfo {
  type: string;
  date_from: string;
  date_to: string;
  previous_period?: {
    date_from: string;
    date_to: string;
  };
}

export interface DashboardResponse {
  period: PeriodInfo;
  revenue: Metric;
  orders: Metric;
  items_sold: Metric;
  avg_order_value: Metric;
  conversion_rate: Metric;
  return_rate: Metric;
  top_categories: TopCategory[];
  hourly_sales: HourlySale[];
}

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