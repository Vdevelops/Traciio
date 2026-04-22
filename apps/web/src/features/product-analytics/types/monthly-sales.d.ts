// Monthly Sales types

export interface MonthlySalesData {
  month: number;
  month_name: string;
  year: number;
  total_sold: number;
  total_revenue: number;
  total_profit: number;
  sales_count: number;
}

export interface MonthlySalesResponse {
  year: number;
  monthly_sales: MonthlySalesData[];
  total_sold: number;
  total_revenue: number;
  total_profit: number;
  total_sales: number;
}

export interface GetMonthlySalesResponse {
  success: boolean;
  data: MonthlySalesResponse;
  meta?: {
    filters?: {
      year?: number;
    };
  };
  timestamp?: string;
  request_id?: string;
}
