// Product Analytics types

export interface TopProduct {
  product_id: string;
  product_name: string;
  product_sku: string;
  quantity_sold: number;
  total_revenue: number;
  avg_price: number;
  growth_rate: number;
  rank: number;
}

export interface PeriodSalesData {
  period: string;
  quantity: number;
  revenue: number;
}

export interface BuyerData {
  user_id: string;
  user_name: string;
  quantity: number;
  revenue: number;
}

export interface ProductPerformance {
  product_id: string;
  product_name: string;
  product_sku: string;
  total_quantity: number;
  total_revenue: number;
  avg_price: number;
  total_sales: number;
  unique_buyers: number;
  growth_rate: number;
  sales_by_period: PeriodSalesData[];
  top_buyers: BuyerData[];
}

export interface ProductTrend {
  product_id: string;
  product_name: string;
  product_sku: string;
  trends: PeriodSalesData[];
}

// Request types
export interface GetTopProductsRequest {
  period?: "day" | "week" | "month" | "year";
  metric?: "quantity" | "revenue" | "growth";
  limit?: number;
}

export interface GetProductPerformanceRequest {
  product_id: string;
  start_date?: string;
  end_date?: string;
}

export interface GetProductComparisonRequest {
  product_ids: string[];
  start_date?: string;
  end_date?: string;
}

export interface GetProductTrendsRequest {
  product_id: string;
  start_date?: string;
  end_date?: string;
  group_by?: "day" | "week" | "month" | "year";
}

// Response types
export interface TopProductsResponse {
  success: boolean;
  data: TopProduct[];
  meta?: {
    filters?: {
      period?: string;
      metric?: string;
      limit?: number;
    };
  };
  timestamp?: string;
  request_id?: string;
}

export interface ProductPerformanceResponse {
  success: boolean;
  data: ProductPerformance;
  meta?: {
    filters?: {
      start_date?: string;
      end_date?: string;
    };
  };
  timestamp?: string;
  request_id?: string;
}

export interface ProductComparisonResponse {
  success: boolean;
  data: {
    products: ProductPerformance[];
  };
  meta?: {
    filters?: {
      product_ids?: string[];
      start_date?: string;
      end_date?: string;
    };
  };
  timestamp?: string;
  request_id?: string;
}

export interface ProductTrendsResponse {
  success: boolean;
  data: ProductTrend;
  meta?: {
    filters?: {
      start_date?: string;
      end_date?: string;
      group_by?: string;
    };
  };
  timestamp?: string;
  request_id?: string;
}

// Product List types
export interface ProductListItem {
  product_id: string;
  product_name: string;
  product_sku: string;
  category_id: string;
  category_name: string;
  unit_price: number;
  total_sold: number;
  total_revenue: number;
  total_profit: number;
  avg_unit_price: number;
  sales_count: number;
  last_sold_at?: string;
  rank?: number;
  image_url?: string;
}

// Filter types
export interface ProductAnalyticsFilters {
  period: "day" | "week" | "month" | "year";
  metric: "quantity" | "revenue" | "growth";
  startDate: Date | null;
  endDate: Date | null;
  limit: number;
  groupBy: "day" | "week" | "month" | "year";
}
