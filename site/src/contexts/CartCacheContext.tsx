
import { createContext, useContext, useState, ReactNode, useRef } from 'react';
import { CartApiRequest, CartApiResponse } from '@/types/cart';
import { fetchCartData } from '@/lib/cartApi';

interface CartCacheContextType {
  getCachedCartData: (request: CartApiRequest) => Promise<CartApiResponse>;
  clearCache: () => void;
}

const CartCacheContext = createContext<CartCacheContextType | undefined>(undefined);

export function CartCacheProvider({ children }: { children: ReactNode }) {
  const [cachedResponse, setCachedResponse] = useState<CartApiResponse | null>(null);
  const lastRequestRef = useRef<CartApiRequest | null>(null);

  const requestsEqual = (req1: CartApiRequest, req2: CartApiRequest): boolean => {
    return JSON.stringify(req1) === JSON.stringify(req2);
  };

  const getCachedCartData = async (request: CartApiRequest): Promise<CartApiResponse> => {
    // Проверяем, есть ли кэшированный ответ и не изменился ли запрос
    if (cachedResponse && lastRequestRef.current && requestsEqual(request, lastRequestRef.current)) {
      console.log('Using cached cart data');
      return cachedResponse;
    }

    console.log('Fetching new cart data');
    const response = await fetchCartData(request);
    
    // Сохраняем в кэш
    setCachedResponse(response);
    lastRequestRef.current = request;
    
    return response;
  };

  const clearCache = () => {
    setCachedResponse(null);
    lastRequestRef.current = null;
    console.log('Cart cache cleared');
  };

  return (
    <CartCacheContext.Provider
      value={{
        getCachedCartData,
        clearCache
      }}
    >
      {children}
    </CartCacheContext.Provider>
  );
}

export function useCartCache() {
  const context = useContext(CartCacheContext);
  if (context === undefined) {
    throw new Error('useCartCache must be used within a CartCacheProvider');
  }
  return context;
}
