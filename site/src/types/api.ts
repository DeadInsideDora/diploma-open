
export interface ApiFilter {
  name: string;
  values: string[];
}

export interface ApiCategory {
  type: string;
  filters: ApiFilter[];
}

export interface ApiMasterData {
  key: string;
  value: string;
}

export interface ApiPrice {
  price_discount: number;
  price_regular: number;
  shop_name: string;
}

export interface ApiProduct {
  title: string;
  category: string;
  image: string;
  master_data: ApiMasterData[];
  prices: ApiPrice[];
}

export interface ProductsRequest {
  type: string;
  filters: ApiFilter[];
}
