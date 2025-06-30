
import { useState, useEffect, useRef } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useNavigate } from "react-router-dom";
import { updateMapInfo, updateExchange, updateLoyaltyCards, getAvailableShops } from "@/lib/userApi";
import { useQuery } from "@tanstack/react-query";
import { toast } from "@/components/ui/sonner";

declare global {
  interface Window { DG: any; }
}

export default function Profile() {
  const { currentUser, logout, updateUserData } = useAuth();
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState("profile");

  const [lat, setLat] = useState(currentUser?.map_info.point.lat || 0);
  const [lon, setLon] = useState(currentUser?.map_info.point.lon || 0);
  const [radius, setRadius] = useState(currentUser?.map_info.radius || 0);

  const [exchange, setExchange] = useState(currentUser?.exchange || 0);

  const [selectedCards, setSelectedCards] = useState<string[]>(currentUser?.cards || []);

  useEffect(() => {
    if (!currentUser) return;
  
    setLat(currentUser.map_info.point.lat);
    setLon(currentUser.map_info.point.lon);
    setRadius(currentUser.map_info.radius);
    setExchange(currentUser.exchange);
    setSelectedCards(currentUser.cards);
  }, [currentUser]);

  const { data: availableShops = [] } = useQuery({
    queryKey: ['availableShops'],
    queryFn: getAvailableShops,
    enabled: activeTab === "cards"
  });

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  const handleUpdateMapInfo = async () => {
    try {
      const updatedUser = await updateMapInfo({
        point: { lat, lon },
        radius
      });
      updateUserData(updatedUser);
      toast.success("Geo data updated successfully");
    } catch (error) {
      toast.error("Failed to update geo data");
    }
  };

  const handleUpdateExchange = async () => {
    try {
      const updatedUser = await updateExchange({ exchange });
      updateUserData(updatedUser);
      toast.success("Exchange updated successfully");
    } catch (error) {
      toast.error("Failed to update exchange");
    }
  };

  const handleUpdateCards = async () => {
    try {
      const updatedUser = await updateLoyaltyCards(selectedCards);
      updateUserData(updatedUser);
      toast.success("Loyalty cards updated successfully");
    } catch (error) {
      toast.error("Failed to update loyalty cards");
    }
  };

  const handleCardToggle = (shop: string, checked: boolean) => {
    if (checked) {
      setSelectedCards(prev => [...prev, shop]);
    } else {
      setSelectedCards(prev => prev.filter(card => card !== shop));
    }
  };

  const mapRef = useRef<any>(null);
  const markerRef = useRef<any>(null);
  const circleRef = useRef<any>(null);

  useEffect(() => {
    if (activeTab === "geo" && window.DG) {
      window.DG.then(() => {
        mapRef.current = window.DG.map("map", { center: [lat, lon], zoom: 12 });
        markerRef.current = window.DG.marker([lat, lon]).addTo(mapRef.current);
        circleRef.current = window.DG.circle([lat, lon], radius, {
          color: '#136AEC',
          fillColor: '#136AEC',
          fillOpacity: 0.2,
        }).addTo(mapRef.current);

        mapRef.current.on('click', e => {
          const { lat: clickedLat, lng: clickedLng } = e.latlng;
          setLat(clickedLat);
          setLon(clickedLng);

          markerRef.current.setLatLng(e.latlng);
          circleRef.current.setLatLng(e.latlng);
        });
      });
    }
  }, [activeTab]);
  
  // Обновление позиции и радиуса при изменении состояний
  useEffect(() => {
    if (mapRef.current) {
      // Перемещаем маркер
      markerRef.current.setLatLng([lat, lon]);
      // Центрируем карту
      mapRef.current.panTo([lat, lon]);
      // Обновляем окружность
      if (circleRef.current) {
        circleRef.current.setLatLng([lat, lon]);
        circleRef.current.setRadius(radius);
      }
    }
  }, [lat, lon, radius]);

  if (!currentUser) {
    navigate("/login");
    return null;
  }

  return (
    <div className="container py-8">
      <h1 className="text-3xl font-bold mb-6">Profile</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="md:col-span-1">
          <Card>
            <CardHeader className="flex flex-col items-center">
              <CardTitle className="text-xl">{currentUser.name}</CardTitle>
              <CardDescription>{currentUser.login}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex flex-col space-y-4">
                <Button onClick={handleLogout} variant="outline">
                  Log Out
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="md:col-span-2">
          <Tabs
            defaultValue="profile"
            value={activeTab}
            onValueChange={setActiveTab}
            className="w-full"
          >
            <TabsList className="mb-6">
              <TabsTrigger value="profile">Profile Information</TabsTrigger>
              <TabsTrigger value="geo">Геоданные</TabsTrigger>
              <TabsTrigger value="cards">Карты лояльности</TabsTrigger>
              <TabsTrigger value="exchange">Размен</TabsTrigger>
            </TabsList>
            
            <TabsContent value="profile">
              <Card>
                <CardHeader>
                  <CardTitle>Personal Information</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="text-sm font-medium text-muted-foreground">
                        Name
                      </label>
                      <p className="text-lg">{currentUser.name}</p>
                    </div>
                    <div>
                      <label className="text-sm font-medium text-muted-foreground">
                        Login
                      </label>
                      <p className="text-lg">{currentUser.login}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="geo">
              <Card>
                <CardHeader>
                  <CardTitle>Геоданные</CardTitle>
                  <CardDescription>Настройте ваше местоположение и радиус поиска</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  {/* Контейнер для карты */}
                  <div id="map" style={{ width: "100%", height: "400px" }} />
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="text-sm font-medium">Latitude</label>
                      <Input
                        type="number"
                        step="any"
                        value={lat}
                        onChange={e => setLat(parseFloat(e.target.value) || 0)}
                        placeholder="59.927168"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Longitude</label>
                      <Input
                        type="number"
                        step="any"
                        value={lon}
                        onChange={e => setLon(parseFloat(e.target.value) || 0)}
                        placeholder="30.317502"
                      />
                    </div>
                  </div>
                  <div>
                    <label className="text-sm font-medium">Radius (метры)</label>
                    <Input
                      type="number"
                      value={radius}
                      onChange={e => setRadius(parseInt(e.target.value) || 0)}
                      placeholder="3000"
                    />
                  </div>
                  <Button onClick={handleUpdateMapInfo}>Сохранить</Button>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="cards">
              <Card>
                <CardHeader>
                  <CardTitle>Карты лояльности</CardTitle>
                  <CardDescription>
                    Выберите магазины, в которых у вас есть карты лояльности
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  {availableShops.map((shop) => (
                    <div key={shop} className="flex items-center space-x-2">
                      <Checkbox
                        id={shop}
                        checked={selectedCards.includes(shop)}
                        onCheckedChange={(checked) => handleCardToggle(shop, checked as boolean)}
                      />
                      <label htmlFor={shop} className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                        {shop}
                      </label>
                    </div>
                  ))}
                  <Button onClick={handleUpdateCards}>
                    Сохранить
                  </Button>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="exchange">
              <Card>
                <CardHeader>
                  <CardTitle>Размен</CardTitle>
                  <CardDescription>
                    Настройте размен времени на валюту
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div>
                    <label className="text-sm font-medium">Exchange Rate</label>
                    <Input
                      type="number"
                      value={exchange}
                      onChange={(e) => setExchange(parseInt(e.target.value) || 0)}
                      placeholder="600"
                    />
                  </div>
                  <Button onClick={handleUpdateExchange}>
                    Сохранить
                  </Button>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
