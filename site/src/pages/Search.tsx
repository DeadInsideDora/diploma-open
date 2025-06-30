import { useState, useMemo, useEffect } from "react";
import ProductCard from "@/components/ProductCard";
import { Product } from "@/types";
import { ApiCategory, ApiProduct, ProductsRequest } from "@/types/api";
import { fetchCategories, fetchProductsByCategory } from "@/lib/api";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { 
  Collapsible,
  CollapsibleContent,
} from "@/components/ui/collapsible";
import { Button } from "@/components/ui/button";
import { Filter } from "lucide-react";
import { useQuery } from "@tanstack/react-query";

// Функция для преобразования API продукта в формат Product
const transformApiProductToProduct = (apiProduct: ApiProduct): Product => {
  const brand = apiProduct.master_data.find(item => item.key === "Бренд")?.value || "Unknown";
  
  // Находим минимальную и максимальную цены из всех магазинов и делим на 100
  const allPrices: number[] = [];
  apiProduct.prices.forEach(priceInfo => {
    allPrices.push(priceInfo.price_discount / 100, priceInfo.price_regular / 100);
  });
  
  const minPrice = Math.min(...allPrices);
  const maxPrice = Math.max(...allPrices);
  
  // Проверяем, является ли товар весовым
  const isWeighed = apiProduct.master_data.some(item => 
    item.key === "Весовой" && item.value === "Да"
  );
  
  // Создаем характеристики из master_data
  const characteristics = apiProduct.master_data.map(item => ({
    name: item.key,
    value: item.value
  }));
  
  // Создаем фиктивный магазин для совместимости
  const store = {
    id: apiProduct.prices[0]?.shop_name || "unknown",
    name: apiProduct.prices[0]?.shop_name || "Unknown Store",
    address: "Unknown Address",
    location: { lat: 0, lng: 0 }
  };
  
  return {
    id: `${apiProduct.title}-${apiProduct.category}`,
    name: apiProduct.title,
    price: minPrice, // Для совместимости оставляем минимальную цену
    minPrice: minPrice,
    maxPrice: maxPrice,
    isWeighed: isWeighed,
    image: apiProduct.image,
    store: store,
    category: apiProduct.category,
    brand: brand,
    characteristics: characteristics
  };
};

export default function Search() {
  // Состояние фильтров
  const [selectedCategory, setSelectedCategory] = useState<string>("");
  const [selectedBrand, setSelectedBrand] = useState<string>("all");
  const [selectedCharacteristics, setSelectedCharacteristics] = useState<Record<string, string[]>>({});
  const [showFilters, setShowFilters] = useState(false);
  
  // Получаем категории из API
  const { data: categories = [], isLoading: categoriesLoading } = useQuery({
    queryKey: ['categories'],
    queryFn: fetchCategories,
  });
  
  // Устанавливаем первую категорию по умолчанию
  useEffect(() => {
    if (categories.length > 0 && !selectedCategory) {
      console.log('Setting default category to:', categories[0].type);
      setSelectedCategory(categories[0].type);
    }
  }, [categories, selectedCategory]);
  
  // Добавляем отладочные логи
  useEffect(() => {
    console.log('Categories data:', categories);
    console.log('Categories length:', categories.length);
    console.log('Selected category:', selectedCategory);
  }, [categories, selectedCategory]);
  
  // Получаем продукты для выбранной категории
  const { data: products = [], isLoading: productsLoading } = useQuery({
    queryKey: ['products', selectedCategory, selectedBrand, selectedCharacteristics],
    queryFn: async () => {
      console.log('Fetching products for category:', selectedCategory);
      if (!selectedCategory) return [];
      
      const selectedCategoryData = categories.find(cat => cat.type === selectedCategory);
      console.log('Found category data:', selectedCategoryData);
      if (!selectedCategoryData) return [];
      
      // Формируем фильтры для запроса
      const filters = selectedCategoryData.filters.map(filter => {
        if (filter.name === "Бренд" && selectedBrand !== "all") {
          return {
            name: filter.name,
            values: [selectedBrand]
          };
        }
        return {
          name: filter.name,
          values: selectedCharacteristics[filter.name] || filter.values
        };
      });
      
      const request: ProductsRequest = {
        type: selectedCategory,
        filters: filters
      };
      
      console.log('API request:', request);
      const apiProducts = await fetchProductsByCategory(request);
      console.log('API response:', apiProducts);
      return apiProducts.map(transformApiProductToProduct);
    },
    enabled: !!selectedCategory && categories.length > 0,
  });
  
  // Получаем список категорий для селекта (убираем "Все категории")
  const categoryOptions = useMemo(() => {
    console.log('Creating category options from:', categories);
    const options = categories.map(cat => cat.type);
    console.log('Category options:', options);
    return options;
  }, [categories]);
  
  // Получаем список ВСЕХ брендов из ВСЕХ категорий
  const brandOptions = useMemo(() => {
    console.log('Creating brand options from categories:', categories);
    const allBrands = new Set<string>();
    
    categories.forEach(category => {
      const brandFilter = category.filters.find(filter => filter.name === "Бренд");
      if (brandFilter) {
        brandFilter.values.forEach(brand => allBrands.add(brand));
      }
    });
    
    const options = ["all", ...Array.from(allBrands).sort()];
    console.log('Brand options:', options);
    return options;
  }, [categories]);
  
  // Получаем характеристики для выбранной категории (исключая бренд)
  const characteristics = useMemo(() => {
    if (!selectedCategory) return [];
    
    const selectedCategoryData = categories.find(cat => cat.type === selectedCategory);
    return selectedCategoryData?.filters.filter(filter => filter.name !== "Бренд") || [];
  }, [selectedCategory, categories]);
  
  // Группировка продуктов по категориям
  const groupedProducts = useMemo(() => {
    const grouped: Record<string, Product[]> = {};
    
    products.forEach(product => {
      if (!grouped[product.category]) {
        grouped[product.category] = [];
      }
      grouped[product.category].push(product);
    });
    
    return grouped;
  }, [products]);
  
  // Обработчик выбора характеристики
  const handleCharacteristicChange = (name: string, value: string) => {
    setSelectedCharacteristics(prev => {
      const current = prev[name] || [];
      const newValues = value === "" ? [] : [value];
      
      return {
        ...prev,
        [name]: newValues
      };
    });
  };
  
  const resetFilters = () => {
    if (categories.length > 0) {
      setSelectedCategory(categories[0].type);
    }
    setSelectedBrand("all");
    setSelectedCharacteristics({});
  };
  
  if (categoriesLoading) {
    return (
      <div className="container py-8">
        <div className="text-center">Загрузка категорий...</div>
      </div>
    );
  }
  
  return (
    <div className="container py-8">
      <h1 className="text-2xl font-bold mb-6">Поиск продуктов</h1>
      
      <div className="flex flex-col gap-4 mb-8">
        <div className="flex items-center">
          <Button 
            variant="outline"
            className="flex items-center gap-2" 
            onClick={() => setShowFilters(!showFilters)}
          >
            <Filter className="h-4 w-4" />
            {showFilters ? "Скрыть фильтры" : "Показать фильтры"}
          </Button>
          
          <Button 
            variant="ghost" 
            className="ml-auto" 
            onClick={resetFilters}
          >
            Сбросить фильтры
          </Button>
        </div>
        
        <Collapsible open={showFilters}>
          <CollapsibleContent className="space-y-4 border rounded-md p-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <Label htmlFor="category-select" className="mb-2 block">Категория</Label>
                <Select 
                  value={selectedCategory}
                  onValueChange={value => {
                    console.log('Selecting category:', value);
                    setSelectedCategory(value);
                    setSelectedCharacteristics({});
                  }}
                >
                  <SelectTrigger id="category-select" className="w-full">
                    <SelectValue placeholder="Выберите категорию" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      {categoryOptions.map((category) => (
                        <SelectItem key={category} value={category}>
                          {category}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>
              
              <div>
                <Label htmlFor="brand-select" className="mb-2 block">Бренд</Label>
                <Select 
                  value={selectedBrand}
                  onValueChange={setSelectedBrand}
                >
                  <SelectTrigger id="brand-select" className="w-full">
                    <SelectValue placeholder="Выберите бренд" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      {brandOptions.map((brand) => (
                        <SelectItem key={brand} value={brand}>
                          {brand === "all" ? "Все бренды" : brand}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>
            </div>
            
            {characteristics.length > 0 && (
              <div className="mt-4">
                <h3 className="font-medium mb-2">Характеристики</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {characteristics.map(({ name, values }) => (
                    <div key={name} className="space-y-1">
                      <Label className="mb-2 block">{name}</Label>
                      <RadioGroup 
                        onValueChange={(value) => handleCharacteristicChange(name, value)}
                        value={selectedCharacteristics[name]?.[0] || ""}
                      >
                        <div className="flex items-center space-x-2">
                          <RadioGroupItem value="" id={`${name}-all`} />
                          <Label htmlFor={`${name}-all`}>Все</Label>
                        </div>
                        
                        {values.map(value => (
                          <div key={value} className="flex items-center space-x-2">
                            <RadioGroupItem value={value} id={`${name}-${value}`} />
                            <Label htmlFor={`${name}-${value}`}>{value}</Label>
                          </div>
                        ))}
                      </RadioGroup>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </CollapsibleContent>
        </Collapsible>
      </div>
      
      {productsLoading ? (
        <div className="text-center py-12">
          <p>Загрузка продуктов...</p>
        </div>
      ) : Object.keys(groupedProducts).length === 0 ? (
        <div className="text-center py-12">
          <h3 className="text-xl font-medium mb-2">Товары не найдены</h3>
          <p className="text-muted-foreground">
            Попробуйте изменить параметры фильтрации, чтобы найти нужные товары.
          </p>
        </div>
      ) : (
        <div className="space-y-8">
          {Object.entries(groupedProducts).map(([category, products]) => (
            <Accordion 
              key={category} 
              type="single" 
              defaultValue={category} 
              collapsible 
              className="border rounded-md"
            >
              <AccordionItem value={category}>
                <AccordionTrigger className="px-4">
                  <div className="flex justify-between w-full items-center">
                    <span className="font-medium">{category}</span>
                    <span className="text-sm text-muted-foreground">
                      ({products.length} {products.length === 1 ? 'товар' : 
                         products.length > 1 && products.length < 5 ? 'товара' : 'товаров'})
                    </span>
                  </div>
                </AccordionTrigger>
                <AccordionContent className="px-4">
                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 py-2">
                    {products.map((product) => (
                      <ProductCard key={product.id} product={product} />
                    ))}
                  </div>
                </AccordionContent>
              </AccordionItem>
            </Accordion>
          ))}
        </div>
      )}
    </div>
  );
}
