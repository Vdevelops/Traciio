export interface Product {
  id: string;
  name: string;
  sku: string;
  barcode?: string;
  price: number;
  price_formatted?: string;
  cost: number;
  category_id: string;
  category?: ProductCategory;
  description?: string;
  status: "active" | "inactive";
  image_url?: string;
  created_at: string;
  updated_at: string;
}

export interface ProductCategory {
  id: string;
  name: string;
  slug: string;
  description?: string;
  status: "active" | "inactive";
  product_count?: number;
  created_at: string;
  updated_at: string;
}

export interface ListProductsResponse {
  success: boolean;
  data: Product[];
  meta: {
    pagination: {
      page: number;
      per_page: number;
      total: number;
      total_pages: number;
      has_next: boolean;
      has_prev: boolean;
    };
    filters?: Record<string, unknown>;
  };
  timestamp: string;
  request_id: string;
}

export interface ProductResponse {
  success: boolean;
  data: Product;
  timestamp: string;
  request_id: string;
}

export interface ListProductCategoriesResponse {
  success: boolean;
  data: ProductCategory[];
  timestamp: string;
  request_id: string;
}

