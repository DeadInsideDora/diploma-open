
import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { User, ShoppingCart, LogOut, Map, Home, Search } from "lucide-react";
import { useAuth } from "@/contexts/AuthContext";
import { useCart } from "@/contexts/CartContext";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { toast } from "@/components/ui/sonner";

export default function Navigation() {
  const location = useLocation();
  const { currentUser, isAuthenticated, logout } = useAuth();
  const { items } = useCart();
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const handleRestrictedClick = (e: React.MouseEvent, path: string) => {
    if (!isAuthenticated && (path === '/cart' || path === '/map')) {
      e.preventDefault();
      toast.error("Войдите в систему для доступа к этой странице");
    }
  };

  const navigationItems = [
    {
      name: "Главная",
      path: "/",
      icon: <Home className="h-5 w-5" />,
      auth: false,
    },
    {
      name: "Каталог",
      path: "/search",
      icon: <Search className="h-5 w-5" />,
      auth: false,
    },
    {
      name: "Корзина",
      path: "/cart",
      icon: <ShoppingCart className="h-5 w-5" />,
      badge: items.length,
      auth: false,
      restricted: true,
    },
    {
      name: "Карта маршрута",
      path: "/map",
      icon: <Map className="h-5 w-5" />,
      auth: false,
      restricted: true,
    },
    {
      name: "Профиль",
      path: "/profile",
      icon: <User className="h-5 w-5" />,
      auth: true,
    },
  ];

  const closeMenu = () => setIsMenuOpen(false);

  return (
    <nav className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-16 items-center justify-between">
        <Link to="/" className="flex items-center gap-2">
          <ShoppingCart className="h-6 w-6 text-primary" />
          <span className="font-bold text-xl hidden md:block">Умный Шоппинг</span>
        </Link>

        {/* Desktop Navigation */}
        <div className="hidden md:flex items-center gap-6">
          {navigationItems
            .filter((item) => !item.auth || isAuthenticated)
            .map((item) => (
              <Link
                key={item.path}
                to={item.path}
                onClick={(e) => handleRestrictedClick(e, item.path)}
                className={`flex items-center gap-1 text-sm font-medium transition-colors hover:text-primary ${
                  location.pathname === item.path
                    ? "text-primary"
                    : "text-muted-foreground"
                } ${item.restricted && !isAuthenticated ? "opacity-50" : ""}`}
              >
                {item.icon}
                {item.name}
                {item.badge ? (
                  <span className="ml-1 rounded-full bg-primary px-2 py-0.5 text-xs text-white">
                    {item.badge}
                  </span>
                ) : null}
              </Link>
            ))}

          {isAuthenticated ? (
            <Button
              variant="ghost"
              size="sm"
              onClick={logout}
              className="gap-1"
            >
              <LogOut className="h-4 w-4" />
              Выйти
            </Button>
          ) : (
            <Link to="/login">
              <Button size="sm" variant="outline">
                Войти
              </Button>
            </Link>
          )}
        </div>

        {/* Mobile Navigation */}
        <div className="md:hidden">
          <Sheet open={isMenuOpen} onOpenChange={setIsMenuOpen}>
            <SheetTrigger asChild>
              <Button variant="outline" size="icon">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="h-5 w-5"
                >
                  <line x1="4" x2="20" y1="12" y2="12" />
                  <line x1="4" x2="20" y1="6" y2="6" />
                  <line x1="4" x2="20" y1="18" y2="18" />
                </svg>
                <span className="sr-only">Меню</span>
              </Button>
            </SheetTrigger>
            <SheetContent side="right">
              <SheetHeader>
                <SheetTitle>Умный Шоппинг</SheetTitle>
                <SheetDescription>
                  Покупай с умом, экономь время
                </SheetDescription>
              </SheetHeader>
              <div className="mt-6 flex flex-col gap-4">
                {navigationItems
                  .filter((item) => !item.auth || isAuthenticated)
                  .map((item) => (
                    <Link
                      key={item.path}
                      to={item.path}
                      onClick={(e) => {
                        handleRestrictedClick(e, item.path);
                        if (!(item.restricted && !isAuthenticated)) {
                          closeMenu();
                        }
                      }}
                      className={`flex items-center gap-3 text-sm font-medium transition-colors hover:text-primary ${
                        location.pathname === item.path
                          ? "text-primary"
                          : "text-muted-foreground"
                      } ${item.restricted && !isAuthenticated ? "opacity-50" : ""}`}
                    >
                      {item.icon}
                      {item.name}
                      {item.badge ? (
                        <span className="ml-auto rounded-full bg-primary px-2 py-0.5 text-xs text-white">
                          {item.badge}
                        </span>
                      ) : null}
                    </Link>
                  ))}

                {isAuthenticated ? (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => {
                      logout();
                      closeMenu();
                    }}
                    className="justify-start gap-3"
                  >
                    <LogOut className="h-4 w-4" />
                    Выйти
                  </Button>
                ) : (
                  <Link to="/login" onClick={closeMenu}>
                    <Button className="w-full">Войти</Button>
                  </Link>
                )}
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </nav>
  );
}
