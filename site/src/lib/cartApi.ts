
import { CartApiRequest, CartApiResponse } from '@/types/cart';

const CART_API_URL = 'http://localhost:8001';

export async function fetchCartData(request: CartApiRequest): Promise<CartApiResponse> {
  const response = await fetch(`${CART_API_URL}/products`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch cart data');
  }
  
  return response.json();
}
