import { createContext, useContext, useState, ReactNode, useEffect } from 'react';
import { CartItem, Product, Route, Store } from '@/types';
import { CartApiItem } from '@/types/cart';
import { calculateRoute } from '@/lib/mockData';
import { toast } from "@/components/ui/sonner";
import { useCartCache } from './CartCacheContext';

interface CartContextType {
  items: CartItem[];
  addItem: (product: Product, quantity?: number) => void;
  removeItem: (productId: string) => void;
  updateQuantity: (productId: string, quantity: number) => void;
  clearCart: () => void;
  getTotal: () => number;
  getStores: () => Store[];
  getRoute: () => Route;
  getStoreItems: (storeId: string) => CartItem[];
  getApiCartItems: () => CartApiItem[];
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export function CartProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<CartItem[]>([]);

  // Load cart from local storage
  useEffect(() => {
    const savedCart = localStorage.getItem('cart');
    if (savedCart) {
      try {
        setItems(JSON.parse(savedCart));
      } catch (e) {
        console.error('Failed to parse saved cart', e);
      }
    }
  }, []);

  // Save cart to local storage
  useEffect(() => {
    localStorage.setItem('cart', JSON.stringify(items));
  }, [items]);

  const addItem = (product: Product, quantity = (product.isWeighed ? 1 : 10)) => {
    setItems(currentItems => {
      const existingItem = currentItems.find(item => item.product.id === product.id);
      
      if (existingItem) {
        toast.success(`Updated ${product.name} quantity`);
        return currentItems.map(item => 
          item.product.id === product.id
            ? { ...item, quantity: item.quantity + quantity }
            : item
        );
      }
      
      toast.success(`Added ${product.name} to cart`);
      return [...currentItems, { product, quantity }];
    });
  };

  const removeItem = (productId: string) => {
    setItems(currentItems => {
      const item = currentItems.find(item => item.product.id === productId);
      if (item) {
        toast.success(`Removed ${item.product.name} from cart`);
      }
      return currentItems.filter(item => item.product.id !== productId);
    });
  };

  const updateQuantity = (productId: string, quantity: number) => {
    if (quantity <= 0) {
      removeItem(productId);
      return;
    }
    
    setItems(currentItems => 
      currentItems.map(item => 
        item.product.id === productId
          ? { ...item, quantity }
          : item
      )
    );
  };

  const clearCart = () => {
    setItems([]);
    toast.success('Cart cleared');
  };

  const getTotal = () => {
    return items.reduce((total, item) => {
      const price = item.product.price;
      const actualQuantity = item.quantity / 10;
      return total + (price * actualQuantity);
    }, 0);
  };

  const getStores = (): Store[] => {
    const storeMap = new Map<string, Store>();
    
    items.forEach(item => {
      storeMap.set(item.product.store.id, item.product.store);
    });
    
    return Array.from(storeMap.values());
  };

  const getRoute = (): Route => {
    const stores = getStores();
    const { distance, duration } = calculateRoute(stores);
    
    return {
      stores,
      distance,
      duration
    };
  };

  const getStoreItems = (storeId: string): CartItem[] => {
    return items.filter(item => item.product.store.id === storeId);
  };

  const getApiCartItems = (): CartApiItem[] => {
    return items.map(item => ({
      info: {
        name: item.product.name,
        type: item.product.category,
        isWeighed: item.product.isWeighed || false,
        url: item.product.image
      },
      amount: item.quantity
    }));
  };

  return (
    <CartContext.Provider
      value={{
        items,
        addItem,
        removeItem,
        updateQuantity,
        clearCart,
        getTotal,
        getStores,
        getRoute,
        getStoreItems,
        getApiCartItems
      }}
    >
      {children}
    </CartContext.Provider>
  );
}

export function useCart() {
  const context = useContext(CartContext);
  if (context === undefined) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
}
