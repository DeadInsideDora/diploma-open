import { useCart } from "@/contexts/CartContext";
import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ShoppingCart, Map, Lock } from "lucide-react";
import { Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { fetchCartData } from "@/lib/cartApi";
import { CartApiRequest } from "@/types/cart";
import { useCartCache } from "@/contexts/CartCacheContext";

export default function Cart() {
  const { items, clearCart, getApiCartItems } = useCart();
  const { currentUser, isAuthenticated } = useAuth();
  const { getCachedCartData } = useCartCache();

  if (!isAuthenticated) {
    return (
      <div className="container py-16 text-center">
        <Lock className="mx-auto h-16 w-16 text-muted-foreground mb-4" />
        <h1 className="text-2xl font-bold mb-2">Доступ ограничен</h1>
        <p className="text-muted-foreground mb-6">
          Войдите в систему, чтобы получить доступ к корзине
        </p>
        <Link to="/login">
          <Button>Войти в систему</Button>
        </Link>
      </div>
    );
  }

  const cartRequest: CartApiRequest = {
    products: getApiCartItems(),
    discount_cards: currentUser?.cards || [],
    point: currentUser?.map_info.point || { lat: 59.927168, lon: 30.317502 },
    radius: currentUser?.map_info.radius || 6000,
    exchange: currentUser?.exchange || 600
  };

  const { data: cartData, isLoading, error } = useQuery({
    queryKey: ['cart', items],
    queryFn: () => getCachedCartData(cartRequest),
    enabled: items.length > 0,
  });

  const formatQuantity = (amount: number, isWeighed: boolean) => {
    const quantity = amount / 10;
    return isWeighed ? `${quantity.toFixed(1)} кг` : `${quantity} шт`;
  };

  if (items.length === 0) {
    return (
      <div className="container py-16 text-center">
        <ShoppingCart className="mx-auto h-16 w-16 text-muted-foreground mb-4" />
        <h1 className="text-2xl font-bold mb-2">Ваша корзина пуста</h1>
        <p className="text-muted-foreground mb-6">
          Добавьте продукты в корзину, чтобы начать покупки
        </p>
        <Link to="/search">
          <Button>Перейти в каталог</Button>
        </Link>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="container py-16 text-center">
        <div className="text-lg">Загрузка корзины...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container py-16 text-center">
        <div className="text-lg text-red-600">Ошибка загрузки корзины</div>
        <p className="text-muted-foreground mt-2">
          Попробуйте обновить страницу
        </p>
      </div>
    );
  }

  return (
    <div className="container py-8">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6">
        <h1 className="text-3xl font-bold">Корзина</h1>
        <div className="flex gap-4 mt-4 md:mt-0">
          <Button variant="outline" onClick={clearCart}>
            Очистить корзину
          </Button>
          <Link to="/map">
            <Button className="flex items-center gap-2">
              <Map className="h-4 w-4" />
              Просмотр маршрута
            </Button>
          </Link>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          {cartData?.stores.map((store, index) => (
            <div key={index} className="space-y-4">
              <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold flex items-center">
                  <span className="mr-2">🏪</span> {store.store}
                </h2>
                <div className="text-lg font-semibold">
                  {(store.total_price / 1000).toFixed(2)} ₽
                </div>
              </div>
              <div className="grid grid-cols-1 gap-4">
                {store.products.map((product, productIndex) => (
                  <Card key={productIndex} className="p-4">
                    <div className="flex gap-4">
                      <img
                        src={product.info.url}
                        alt={product.info.name}
                        className="h-16 w-16 object-cover rounded"
                      />
                      <div className="flex-1 flex justify-between items-center">
                        <div>
                          <h3 className="font-semibold">{product.info.name}</h3>
                          <p className="text-sm text-muted-foreground">{product.info.type}</p>
                          <p className="text-sm text-muted-foreground">
                            {formatQuantity(product.info.amount, product.info.isWeighed)}
                          </p>
                        </div>
                        <div className="text-right">
                          <div className="font-semibold">
                            {(product.price / 1000).toFixed(2)} ₽
                          </div>
                        </div>
                      </div>
                    </div>
                  </Card>
                ))}
              </div>
            </div>
          ))}
        </div>

        <div className="lg:col-span-1">
          <Card className="p-6 sticky top-20">
            <h2 className="text-xl font-semibold mb-4">Итого</h2>
            <div className="space-y-2 mb-4">
              {cartData?.stores.map((store, index) => (
                <div key={index} className="flex justify-between text-sm">
                  <span>{store.store}:</span>
                  <span>{(store.total_price / 1000).toFixed(2)} ₽</span>
                </div>
              ))}
            </div>
            <div className="border-t pt-4 mb-6">
              <div className="flex justify-between font-semibold text-lg">
                <span>Общая стоимость:</span>
                <span>{cartData ? (cartData.price / 1000).toFixed(2) : '0.00'} ₽</span>
              </div>
            </div>
            <div className="mt-4">
              <Link to="/map" className="text-center block w-full text-primary hover:underline">
                Просмотр маршрута
              </Link>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
