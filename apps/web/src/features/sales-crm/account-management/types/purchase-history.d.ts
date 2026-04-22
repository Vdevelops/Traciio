export interface CustomerPurchaseHistory {
  id: string;
  account_id: string;
  deal_id: string;
  purchase_date: string;
  total_amount: number;
  total_amount_formatted: string;
  total_items: number;
  products: PurchaseProductItem[];
  sales_rep_id?: string;
  sales_rep_name?: string;
  source_lead_id?: string;
  source_type: string;
  customer_lifetime_value: number;
  clv_formatted: string;
  purchase_number: number;
  created_at: string;
}

export interface PurchaseProductItem {
  product_id: string;
  product_name: string;
  product_sku?: string;
  product_category_id?: string;
  product_category_name?: string;
  quantity: number;
  unit_price: number;
  discount_amount?: number;
  subtotal: number;
}

export interface CustomerProductAnalytics {
  account_id: string;
  product_id: string;
  product_name: string;
  product_category_id?: string;
  product_category_name?: string;
  total_quantity_purchased: number;
  total_amount_purchased: number;
  total_amount_formatted: string;
  purchase_count: number;
  first_purchase_date: string;
  last_purchase_date: string;
}

export interface CustomerPurchaseSummary {
  account_id: string;
  account_name: string;
  total_purchases: number;
  total_amount: number;
  total_amount_formatted: string;
  average_order_value: number;
  aov_formatted: string;
  last_purchase_date?: string;
  favorite_product_category?: string;
}
