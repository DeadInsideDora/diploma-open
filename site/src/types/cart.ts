
export interface CartApiRequest {
  products: CartApiItem[];
  discount_cards: string[];
  point: {
    lat: number;
    lon: number;
  };
  radius: number;
  exchange: number;
}

export interface CartApiItem {
  info: {
    type: string;
    name: string;
    isWeighed: boolean;
    url: string;
  };
  amount: number;
}

export interface CartStoreProduct {
  info: {
    type: string;
    name: string;
    isWeighed: boolean;
    url: string;
    amount: number;
  };
  price: number;
}

export interface CartStore {
  products: CartStoreProduct[];
  store: string;
  point: {
    lat: number;
    lon: number;
  };
  total_price: number;
}

export interface CartApiResponse {
  stores: CartStore[];
  price: number;
}
