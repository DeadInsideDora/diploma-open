
export interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
}

export interface ProductCharacteristic {
  name: string;
  value: string;
}

export interface Product {
  id: string;
  name: string;
  price: number;
  minPrice?: number;
  maxPrice?: number;
  isWeighed?: boolean;
  image: string;
  store: Store;
  category: string;
  brand: string;
  characteristics: ProductCharacteristic[];
}

export interface ApiCartItem {
  info: {
    name: string;
    type: string;
  };
  amount: number;
}

export interface Store {
  id: string;
  name: string;
  address: string;
  location: {
    lat: number;
    lng: number;
  };
}

export interface CartItem {
  product: Product;
  quantity: number;
}

export interface Route {
  stores: Store[];
  distance: number;
  duration: number;
}
