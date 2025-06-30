import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Index from "./pages/Index";
import NotFound from "./pages/NotFound";
import Login from "./pages/Login";
import Search from "./pages/Search";
import Cart from "./pages/Cart";
import Profile from "./pages/Profile";
import Map from "./pages/Map";
import Navigation from "./components/Navigation";
import { AuthProvider } from "./contexts/AuthContext";
import { CartProvider } from "./contexts/CartContext";
import { CartCacheProvider } from "./contexts/CartCacheContext";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <AuthProvider>
      <CartProvider>
        <CartCacheProvider>
          <TooltipProvider>
            <Toaster />
            <Sonner />
            <BrowserRouter>
              <div className="flex flex-col min-h-screen">
                <Navigation />
                <main className="flex-1">
                  <Routes>
                    <Route path="/" element={<Index />} />
                    <Route path="/login" element={<Login />} />
                    <Route path="/search" element={<Search />} />
                    <Route path="/cart" element={<Cart />} />
                    <Route path="/profile" element={<Profile />} />
                    <Route path="/map" element={<Map />} />
                    <Route path="*" element={<NotFound />} />
                  </Routes>
                </main>
              </div>
            </BrowserRouter>
          </TooltipProvider>
        </CartCacheProvider>
      </CartProvider>
    </AuthProvider>
  </QueryClientProvider>
);

export default App;
