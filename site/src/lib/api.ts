
import { ApiCategory, ApiProduct, ProductsRequest } from '@/types/api';

const API_BASE_URL = 'http://localhost:8000';

export async function fetchCategories(): Promise<ApiCategory[]> {
  const response = await fetch(`${API_BASE_URL}/categories`);
  if (!response.ok) {
    throw new Error('Failed to fetch categories');
  }
  return response.json();
}

export async function fetchProductsByCategory(request: ProductsRequest): Promise<ApiProduct[]> {
  const response = await fetch(`${API_BASE_URL}/products-by-category`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });
  if (!response.ok) {
    throw new Error('Failed to fetch products');
  }
  return response.json();
}
