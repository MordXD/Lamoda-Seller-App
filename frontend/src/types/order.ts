export interface Customer {
  id: string;
  name: string;
  email: string;
  phone: string;
  is_regular: boolean;
  orders_count: number;
}

export interface Address {
  city: string;
  street: string;
  house: string;
  apartment: string;
  postal_code: string;
  coordinates: {
    lat: number;
    lng: number;
  };
}

export interface DeliveryInfo {
  type: 'courier' | 'pickup' | 'post';
  address: Address;
  estimated_date: string;
  cost: number;
  tracking_number: string;
}

export interface OrderItem {
  id: string;
  product_id: string;
  variant_id: string;
  name: string;
  brand: string;
  sku: string;
  size: string;
  color: string;
  quantity: number;
  price: number;
  cost_price: number;
  discount: number;
  total: number;
  image: string;
}

export interface PaymentInfo {
  method: 'card' | 'cash' | 'online';
  status: 'pending' | 'paid' | 'failed' | 'refunded';
  amount: number;
  commission_lamoda: number;
  seller_amount: number;
  transaction_id: string;
}

export interface TotalsInfo {
  subtotal: number;
  discount: number;
  delivery: number;
  total: number;
}

export interface StatusHistory {
  status: string;
  date: string;
  comment: string;
}

export interface ReturnInfo {
  is_returnable: boolean;
  return_deadline: string;
  return_reasons: string[];
}

export interface Review {
  id: string;
  product_id: string;
  rating: number;
  comment: string;
  date: string;
  is_verified_purchase: boolean;
}

export interface Logistics {
  warehouse: {
    id: string;
    name: string;
    address: string;
  };
  picking_date: string;
  packing_date: string;
  shipping_date: string;
  courier: {
    name: string;
    phone: string;
    company: string;
  };
}

export interface Order {
  id: string;
  order_number: string;
  date: string;
  status: 'new' | 'confirmed' | 'in_transit' | 'delivered' | 'returned' | 'cancelled';
  status_history: StatusHistory[];
  customer: Customer;
  delivery: DeliveryInfo;
  items: OrderItem[];
  payment: PaymentInfo;
  totals: TotalsInfo;
  notes: string;
  created_date: string;
  updated_date: string;
  return_info?: ReturnInfo;
  reviews?: Review[];
  logistics?: Logistics;
}

export interface OrdersResponse {
  orders: Order[];
  summary: {
    total_orders: number;
    total_amount: number;
    avg_order_value: number;
    status_breakdown: {
      [key: string]: {
        count: number;
        amount: number;
      };
    };
  };
  pagination: {
    total: number;
    limit: number;
    offset: number;
    has_next: boolean;
    has_prev: boolean;
  };
}

export interface OrdersFilters {
  status?: 'new' | 'confirmed' | 'in_transit' | 'delivered' | 'returned' | 'cancelled';
  date_from?: string;
  date_to?: string;
  customer_id?: string;
  product_id?: string;
  min_amount?: number;
  max_amount?: number;
  sort_by?: 'date' | 'amount' | 'status';
  sort_order?: 'asc' | 'desc';
  limit?: number;
  offset?: number;
}

export interface UpdateOrderStatusRequest {
  status: 'new' | 'confirmed' | 'in_transit' | 'delivered' | 'returned' | 'cancelled';
  comment?: string;
  estimated_delivery_date?: string;
} 