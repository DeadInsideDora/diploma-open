import { Product, Store, User } from "@/types";

export const mockStores: Store[] = [
  {
    id: "store-1",
    name: "Свежий Маркет",
    address: "ул. Главная, 123",
    location: { lat: 51.505, lng: -0.09 }
  },
  {
    id: "store-2",
    name: "Супер Гастроном",
    address: "пр. Дубовый, 456",
    location: { lat: 51.51, lng: -0.1 }
  },
  {
    id: "store-3",
    name: "Дисконт Фудс",
    address: "ул. Сосновая, 789",
    location: { lat: 51.515, lng: -0.09 }
  }
];

export const mockProducts: Product[] = [
  // Фрукты
  {
    id: "product-1",
    name: "Яблоки",
    price: 89.90,
    image: "/placeholder.svg",
    store: mockStores[0],
    category: "Фрукты",
    brand: "Фруктовый сад",
    characteristics: [
      { name: "Страна", value: "Россия" },
      { name: "Сорт", value: "Голден" },
      { name: "Сладость", value: "Сладкие" }
    ]
  },
  {
    id: "product-6",
    name: "Бананы",
    price: 99.90,
    image: "/placeholder.svg",
    store: mockStores[0],
    category: "Фрукты",
    brand: "Чикита",
    characteristics: [
      { name: "Страна", value: "Эквадор" },
      { name: "Сорт", value: "Кавендиш" },
      { name: "Спелость", value: "Средние" }
    ]
  },
  // Молочные продукты
  {
    id: "product-2",
    name: "Молоко",
    price: 89.90,
    image: "/placeholder.svg",
    store: mockStores[0],
    category: "Молочные продукты",
    brand: "Домик в деревне",
    characteristics: [
      { name: "Жирность", value: "2.5%" },
      { name: "Объем", value: "1 л" },
      { name: "Тип", value: "Пастеризованное" }
    ]
  },
  {
    id: "product-7",
    name: "Творог",
    price: 119.90,
    image: "/placeholder.svg",
    store: mockStores[1],
    category: "Молочные продукты",
    brand: "Простоквашино",
    characteristics: [
      { name: "Жирность", value: "5%" },
      { name: "Вес", value: "200 г" },
      { name: "Тип", value: "Мягкий" }
    ]
  },
  // Хлебобулочные изделия
  {
    id: "product-3",
    name: "Хлеб",
    price: 59.90,
    image: "/placeholder.svg",
    store: mockStores[1],
    category: "Хлебобулочные изделия",
    brand: "Хлебный дом",
    characteristics: [
      { name: "Тип", value: "Пшеничный" },
      { name: "Вес", value: "300 г" },
      { name: "Нарезка", value: "Да" }
    ]
  },
  {
    id: "product-8",
    name: "Багет",
    price: 69.90,
    image: "/placeholder.svg",
    store: mockStores[2],
    category: "Хлебобулочные изделия",
    brand: "Французская пекарня",
    characteristics: [
      { name: "Тип", value: "Классический" },
      { name: "Вес", value: "250 г" },
      { name: "Корочка", value: "Хрустящая" }
    ]
  },
  // Мясо
  {
    id: "product-4",
    name: "Курица",
    price: 249.90,
    image: "/placeholder.svg",
    store: mockStores[1],
    category: "Мясо",
    brand: "Петелинка",
    characteristics: [
      { name: "Тип", value: "Охлажденная" },
      { name: "Вес", value: "1 кг" },
      { name: "Часть", value: "Филе" }
    ]
  },
  {
    id: "product-9",
    name: "Говядина",
    price: 549.90,
    image: "/placeholder.svg",
    store: mockStores[2],
    category: "Мясо",
    brand: "Мираторг",
    characteristics: [
      { name: "Тип", value: "Охлажденная" },
      { name: "Вес", value: "500 г" },
      { name: "Часть", value: "Вырезка" }
    ]
  },
  // Крупы
  {
    id: "krupy-1",
    name: "Гречка Увелка Экстра, 5х80г",
    price: 69.99,
    image: "/placeholder.svg",
    store: {
      id: "store-1",
      name: "Магнит",
      address: "ул. Пушкина, д. 10",
      location: {
        lat: 55.753215,
        lng: 37.622504
      }
    },
    category: "Крупы",
    brand: "Увелка",
    characteristics: [
      { name: "Тип крупы", value: "Гречневая" },
      { name: "Вес", value: "5х80г" },
      { name: "В пакетиках", value: "Да" }
    ]
  },
  {
    id: "krupy-2",
    name: "Крупа гречневая Мистраль, 900г",
    price: 101.99,
    image: "/placeholder.svg",
    store: {
      id: "store-2",
      name: "Пятёрочка",
      address: "ул. Ленина, д. 15",
      location: {
        lat: 55.755814,
        lng: 37.617635
      }
    },
    category: "Крупы",
    brand: "Мистраль",
    characteristics: [
      { name: "Тип крупы", value: "Гречневая" },
      { name: "Вес", value: "900г" },
      { name: "В пакетиках", value: "Нет" }
    ]
  },
  {
    id: "krupy-3",
    name: "Рис Кубанский Экстра Агро-Альянс, 5х80г",
    price: 83.99,
    image: "/placeholder.svg",
    store: {
      id: "store-3",
      name: "Перекрёсток",
      address: "ул. Гагарина, д. 20",
      location: {
        lat: 55.751244,
        lng: 37.618423
      }
    },
    category: "Крупы",
    brand: "Агро-Альянс",
    characteristics: [
      { name: "Тип крупы", value: "Рисовая" },
      { name: "Вес", value: "5х80г" },
      { name: "В пакетиках", value: "Да" }
    ]
  }
];

export const mockUsers: User[] = [
  {
    id: "user-1",
    name: "Иван Иванов",
    email: "ivan@example.com",
    avatar: "https://i.pravatar.cc/150?u=user-1"
  }
];

export function calculateRoute(stores: Store[]): {distance: number; duration: number} {
  // In a real app, this would call a routing API
  // For mock purposes, we'll just calculate some dummy values
  if (stores.length <= 1) {
    return {distance: 0, duration: 0};
  }
  
  let totalDistance = 0;
  
  for (let i = 0; i < stores.length - 1; i++) {
    const store1 = stores[i];
    const store2 = stores[i + 1];
    
    // Simple Euclidean distance calculation
    const distance = Math.sqrt(
      Math.pow(store2.location.lat - store1.location.lat, 2) + 
      Math.pow(store2.location.lng - store1.location.lng, 2)
    ) * 111; // Rough conversion to kilometers
    
    totalDistance += distance;
  }
  
  // Assuming average speed of 30 km/h for city driving
  const duration = totalDistance / 30 * 60; // Convert to minutes
  
  return {
    distance: parseFloat(totalDistance.toFixed(2)),
    duration: parseFloat(duration.toFixed(0))
  };
}
