
import { Product } from "@/types";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { ShoppingCart, Plus, Minus, Trash, Lock } from "lucide-react";
import { useState } from "react";
import { useCart } from "@/contexts/CartContext";
import { useAuth } from "@/contexts/AuthContext";
import { toast } from "@/components/ui/sonner";

interface ProductCardProps {
  product: Product;
  inCart?: boolean;
}

export default function ProductCard({ product, inCart = false }: ProductCardProps) {
  const { addItem, updateQuantity, removeItem, items } = useCart();
  const { isAuthenticated } = useAuth();
  const cartItem = items.find(item => item.product.id === product.id);
  
  const step = product.isWeighed ? 1 : 10;
  const minQuantity = step;
  
  const [quantity, setQuantity] = useState(cartItem?.quantity || minQuantity);

  const handleAddToCart = () => {
    if (!isAuthenticated) {
      toast.error("Войдите в систему, чтобы добавлять товары в корзину");
      return;
    }
    addItem(product, quantity);
  };

  const handleIncreaseQuantity = () => {
    if (!isAuthenticated) {
      toast.error("Войдите в систему, чтобы изменять количество товаров");
      return;
    }

    if (cartItem) {
      updateQuantity(product.id, cartItem.quantity + step);
    } else {
      setQuantity(q => q + step);
    }
  };

  const handleDecreaseQuantity = () => {
    if (!isAuthenticated) {
      toast.error("Войдите в систему, чтобы изменять количество товаров");
      return;
    }

    if (cartItem) {
      const newQuantity = cartItem.quantity - step;
      if (newQuantity <= 0) {
        removeItem(product.id);
        setQuantity(minQuantity);
      } else {
        updateQuantity(product.id, newQuantity);
      }
    } else {
      setQuantity(q => Math.max(q - step, minQuantity));
    }
  };

  const handleRemoveFromCart = () => {
    if (!isAuthenticated) {
      toast.error("Войдите в систему, чтобы удалять товары из корзины");
      return;
    }

    removeItem(product.id);

    setQuantity(minQuantity);
  };

  const formatQuantity = (qty: number) => {
    if (product.isWeighed) {
      return `${(qty / 10).toFixed(1)} кг`;
    }
    return `${(qty / 10)} шт`;
  };

  const formatPrice = () => {
    if (product.minPrice !== undefined && product.maxPrice !== undefined && product.minPrice !== product.maxPrice) {
      const suffix = product.isWeighed ? " ₽/кг" : " ₽";
      return `${(product.minPrice).toFixed(2)} - ${(product.maxPrice).toFixed(2)}${suffix}`;
    }
    const suffix = product.isWeighed ? " ₽/кг" : " ₽";
    return `${(product.price).toFixed(2)}${suffix}`;
  };

  return (
    <Card className="overflow-hidden">
      <div className="relative">
        <img
          src={product.image}
          alt={product.name}
          className="h-40 w-full object-cover"
        />
        <Badge className="absolute top-2 right-2 bg-primary text-xs">
          {formatPrice()}
        </Badge>
      </div>
      <CardContent className="p-4">
        <div className="flex justify-between items-start mb-2">
          <h3 className="font-semibold text-lg">{product.name}</h3>
          <Badge variant="outline" className="text-xs">
            {product.category}
          </Badge>
        </div>
        <div className="space-y-1 text-sm">
          <p className="text-muted-foreground">
            Бренд: {product.brand}
          </p>
          {product.isWeighed && (
            <p className="text-muted-foreground text-xs">
              Весовой товар
            </p>
          )}
        </div>
      </CardContent>
      <CardFooter className="p-4 pt-0">
        {!isAuthenticated ? (
          <div className="w-full flex justify-center">
            <Button variant="outline" disabled className="flex items-center gap-2">
              <Lock className="h-4 w-4" />
              Войдите, чтобы добавить в корзину
            </Button>
          </div>
        ) : inCart ? (
          <div className="w-full flex justify-center items-center">
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="icon"
                className="h-8 w-8"
                onClick={handleDecreaseQuantity}
              >
                <Minus className="h-4 w-4" />
              </Button>
              <span className="min-w-[60px] text-center">
                {formatQuantity(cartItem?.quantity || 0)}
              </span>
              <Button
                variant="outline"
                size="icon"
                className="h-8 w-8"
                onClick={handleIncreaseQuantity}
              >
                <Plus className="h-4 w-4" />
              </Button>
            </div>
          </div>
        ) : (
          <div className="w-full">
            {cartItem ? (
              <div className="flex justify-between items-center gap-2">
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="icon"
                    className="h-8 w-8"
                    onClick={handleDecreaseQuantity}
                  >
                    <Minus className="h-4 w-4" />
                  </Button>
                  <span className="min-w-[60px] text-center">
                    {formatQuantity(cartItem.quantity)}
                  </span>
                  <Button
                    variant="outline"
                    size="icon"
                    className="h-8 w-8"
                    onClick={handleIncreaseQuantity}
                  >
                    <Plus className="h-4 w-4" />
                  </Button>
                </div>
                <Button
                  variant="destructive"
                  size="icon"
                  className="h-8 w-8"
                  onClick={handleRemoveFromCart}
                >
                  <Trash className="h-4 w-4" />
                </Button>
              </div>
            ) : (
              <div className="flex justify-between items-center">
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="icon"
                    className="h-8 w-8"
                    onClick={() => setQuantity(q => Math.max(q - step, minQuantity))}
                  >
                    <Minus className="h-4 w-4" />
                  </Button>
                  <span className="min-w-[60px] text-center">
                    {formatQuantity(quantity)}
                  </span>
                  <Button
                    variant="outline"
                    size="icon"
                    className="h-8 w-8"
                    onClick={() => setQuantity(q => q + step)}
                  >
                    <Plus className="h-4 w-4" />
                  </Button>
                </div>
                <Button
                  onClick={handleAddToCart}
                  className="flex items-center gap-2"
                >
                  <ShoppingCart className="h-4 w-4" />
                  В корзину
                </Button>
              </div>
            )}
          </div>
        )}
      </CardFooter>
    </Card>
  );
}
