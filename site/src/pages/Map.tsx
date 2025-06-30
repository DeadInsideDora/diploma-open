import { useState, useEffect } from "react";
import { useCart } from "@/contexts/CartContext";
import { useAuth } from "@/contexts/AuthContext";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Link } from "react-router-dom";
import { ShoppingCart, Navigation } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { fetchCartData } from "@/lib/cartApi";
import { CartApiRequest } from "@/types/cart";
import { useCartCache } from "@/contexts/CartCacheContext";

export default function Map() {
  const { items, getStores, getRoute, getStoreItems, getApiCartItems } = useCart();
  const { currentUser, isAuthenticated } = useAuth();
  const { getCachedCartData } = useCartCache();
  const stores = getStores();
  const route = getRoute();

  if (!isAuthenticated) {
    return (
      <div className="container py-16 text-center">
        <Lock className="mx-auto h-16 w-16 text-muted-foreground mb-4" />
        <h1 className="text-2xl font-bold mb-2">Доступ ограничен</h1>
        <p className="text-muted-foreground mb-6">
          Войдите в систему, чтобы получить доступ к карте маршрута
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

  if (items.length === 0) {
    return (
      <div className="container py-16 text-center">
        <ShoppingCart className="mx-auto h-16 w-16 text-muted-foreground mb-4" />
        <h1 className="text-2xl font-bold mb-2">Ваша корзина пуста</h1>
        <p className="text-muted-foreground mb-6">
          Добавьте товары в корзину, чтобы сгенерировать маршрут покупок
        </p>
        <Link to="/search">
          <Button>Перейти в каталог</Button>
        </Link>
      </div>
    );
  }

  const generateMapUrl = () => {
    if (!cartData?.stores.length) return "";

    let url = "https://static.maps.2gis.com/1.0?s=800x800";
    
    // Add user location point
    url += `&pt=${cartRequest.point.lat},${cartRequest.point.lon}~k:c~n:1`;
    
    // Add store points
    cartData.stores.forEach((store, index) => {
      url += `&pt=${store.point.lat},${store.point.lon}~k:c~n:${index + 2}`;
    });
    
    return url;
  };

  if (isLoading) {
    return (
      <div className="container py-16 text-center">
        <div className="text-lg">Загрузка маршрута...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container py-16 text-center">
        <div className="text-lg text-red-600">Ошибка загрузки маршрута</div>
        <p className="text-muted-foreground mt-2">
          Попробуйте обновить страницу
        </p>
      </div>
    );
  }

  return (
    <div className="container py-8">
      <h1 className="text-3xl font-bold mb-6">Карта маршрута покупок</h1>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <Card className="h-[600px] relative overflow-hidden">
            {cartData && (
              <img 
                src={generateMapUrl()} 
                alt="Карта маршрута покупок" 
                className="w-full h-full object-cover"
              />
            )}
            
            {/* <div className="absolute bottom-4 left-4 right-4 bg-white p-4 rounded-lg shadow-lg bg-opacity-90">
              <div className="flex items-center justify-between">
                <div>
                  <div className="font-semibold text-sm">Общее расстояние</div>
                  <div className="text-lg font-bold">{route.distance} км</div>
                </div>
                <div>
                  <div className="font-semibold text-sm">Расчетное время</div>
                  <div className="text-lg font-bold">{route.duration} мин</div>
                </div>
              </div>
            </div> */}
          </Card>
        </div>

        <div className="lg:col-span-1">
          <Card className="sticky top-20">
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                Список покупок
                <Badge variant="outline" className="ml-2">
                  {items.length} товаров
                </Badge>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="flex items-center gap-2 mb-4">
                <div className="w-6 h-6 rounded-full bg-primary flex items-center justify-center text-white text-sm font-bold">
                  1
                </div>
                <h3 className="font-semibold">Точка старта</h3>
              </div>
              
              {cartData?.stores.map((store, index) => {
                const storeItems = getStoreItems(store.store);
                return (
                  <div key={store.store}>
                    <div className="flex items-center gap-2 mb-2">
                      <div className="w-6 h-6 rounded-full bg-primary flex items-center justify-center text-white text-sm font-bold">
                        {index + 2}
                      </div>
                      <h3 className="font-semibold">{store.store}</h3>
                    </div>
                    <div className="pl-8 space-y-2">
                      {store.products.map((product, productIndex) => (
                        <div
                          key={productIndex}
                          className="flex justify-between items-center"
                        >
                          <div className="flex items-center gap-2">
                            <span>
                              {product.info.name}
                            </span>
                            <Badge variant="outline" className="text-xs">
                              {product.info.isWeighed 
                                ? `${(product.info.amount / 10).toFixed(1)} кг`
                                : `${(product.info.amount / 10)} шт`
                              }
                            </Badge>
                          </div>
                          <span className="text-sm font-medium">
                            {(product.price / 1000).toFixed(2)} ₽
                          </span>
                        </div>
                      ))}
                    </div>
                    {index < cartData.stores.length - 1 && (
                      <Separator className="my-4" />
                    )}
                  </div>
                );
              })}
              
              <div className="pt-2 flex justify-between items-center font-semibold">
                <span>Итог маршрута</span>
                <div className="text-right">
                  <div>{route.distance} км</div>
                  <div className="text-muted-foreground text-sm">
                    {route.duration} мин
                  </div>
                </div>
              </div>
              
              <div className="pt-4">
                <Link to="/cart">
                  <Button variant="outline" className="w-full">
                    Вернуться в корзину
                  </Button>
                </Link>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
