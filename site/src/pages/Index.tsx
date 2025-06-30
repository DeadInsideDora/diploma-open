import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import { ShoppingCart, Map, Search } from "lucide-react";

export default function Index() {
  return (
    <div className="flex flex-col min-h-[calc(100vh-4rem)]">
      {/* Hero section */}
      <section className="bg-gradient-to-b from-secondary/20 to-background py-12 md:py-24">
        <div className="container px-4 md:px-6">
          <div className="grid gap-6 lg:grid-cols-[1fr_400px] lg:gap-12 xl:grid-cols-[1fr_600px]">
            <div className="flex flex-col justify-center space-y-4">
              <div className="space-y-2">
                <h1 className="text-3xl font-bold tracking-tighter sm:text-5xl xl:text-6xl/none">
                  Покупайте выгодно, экономьте время
                </h1>
                <p className="max-w-[600px] text-gray-500 md:text-xl dark:text-gray-400">
                  Оптимизируйте свои покупки, находите лучшие товары, планируйте маршрут между магазинами и экономьте время и деньги.
                </p>
              </div>
              <div className="flex flex-col gap-2 min-[400px]:flex-row">
                <Link to="/search">
                  <Button size="lg" className="gap-2">
                    <Search className="h-4 w-4" />
                    Поиск товаров
                  </Button>
                </Link>
                <Link to="/map">
                  <Button size="lg" variant="outline" className="gap-2">
                    <Map className="h-4 w-4" />
                    Карта маршрута
                  </Button>
                </Link>
              </div>
            </div>
            <div className="flex items-center justify-center">
              <img
                alt="Изображение приложения"
                className="aspect-video overflow-hidden rounded-xl object-cover object-center"
                src="/placeholder.svg"
                width={550}
                height={310}
              />
            </div>
          </div>
        </div>
      </section>

      {/* Features section */}
      <section className="container py-12 md:py-24 lg:py-32">
        <div className="mx-auto grid max-w-5xl gap-6 px-4 md:px-6 lg:grid-cols-3 lg:gap-12">
          <div className="flex flex-col justify-center space-y-4">
            <div className="rounded-lg bg-primary p-3 w-12 h-12 flex items-center justify-center">
              <ShoppingCart className="h-6 w-6 text-white" />
            </div>
            <h3 className="text-xl font-bold">Умный список покупок</h3>
            <p className="text-gray-500 dark:text-gray-400">
              Создайте свой список покупок, и мы подскажем, в каком магазине есть каждый товар.
            </p>
          </div>
          <div className="flex flex-col justify-center space-y-4">
            <div className="rounded-lg bg-primary p-3 w-12 h-12 flex items-center justify-center">
              <Map className="h-6 w-6 text-white" />
            </div>
            <h3 className="text-xl font-bold">Оптимальные маршруты</h3>
            <p className="text-gray-500 dark:text-gray-400">
              Планируйте наиболее эффективный маршрут между магазинами, чтобы сэкономить время.
            </p>
          </div>
          <div className="flex flex-col justify-center space-y-4">
            <div className="rounded-lg bg-primary p-3 w-12 h-12 flex items-center justify-center">
              <Search className="h-6 w-6 text-white" />
            </div>
            <h3 className="text-xl font-bold">Поиск товаров</h3>
            <p className="text-gray-500 dark:text-gray-400">
              Находите лучшие товары в разных магазинах в одном месте.
            </p>
          </div>
        </div>
      </section>

      {/* How it works section */}
      <section className="bg-muted py-12 md:py-24 lg:py-32">
        <div className="container px-4 md:px-6">
          <div className="mx-auto max-w-3xl space-y-4 text-center">
            <h2 className="text-3xl font-bold tracking-tighter md:text-4xl">
              Как работает Умный Шоппинг
            </h2>
            <p className="text-gray-500 dark:text-gray-400">
              Наше приложение помогает вам экономить время и деньги, планируя наиболее эффективный маршрут покупок.
            </p>
          </div>
          <div className="mx-auto grid max-w-5xl grid-cols-1 gap-6 md:grid-cols-3 md:gap-12 pt-12">
            <div className="flex flex-col items-center space-y-4 text-center">
              <div className="flex h-12 w-12 items-center justify-center rounded-full border border-gray-200 bg-white">
                <span className="text-xl font-bold">1</span>
              </div>
              <h3 className="text-xl font-bold">Поиск товаров</h3>
              <p className="text-gray-500 dark:text-gray-400">
                Ищите товары в разных магазинах и добавляйте их в корзину.
              </p>
            </div>
            <div className="flex flex-col items-center space-y-4 text-center">
              <div className="flex h-12 w-12 items-center justify-center rounded-full border border-gray-200 bg-white">
                <span className="text-xl font-bold">2</span>
              </div>
              <h3 className="text-xl font-bold">Соберите корзину</h3>
              <p className="text-gray-500 dark:text-gray-400">
                Добавьте все необходимые товары в корзину из разных магазинов.
              </p>
            </div>
            <div className="flex flex-col items-center space-y-4 text-center">
              <div className="flex h-12 w-12 items-center justify-center rounded-full border border-gray-200 bg-white">
                <span className="text-xl font-bold">3</span>
              </div>
              <h3 className="text-xl font-bold">Просмотр маршрута</h3>
              <p className="text-gray-500 dark:text-gray-400">
                Получите оптимизированную карту маршрута, показывающую наиболее эффективный путь между магазинами.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA section */}
      <section className="container py-12 md:py-24 lg:py-32">
        <div className="mx-auto max-w-5xl space-y-6 text-center">
          <h2 className="text-3xl font-bold tracking-tighter md:text-4xl">
            Готовы начать умные покупки?
          </h2>
          <p className="text-xl text-gray-500 dark:text-gray-400">
            Зарегистрируйтесь сейчас и оптимизируйте свой опыт покупок.
          </p>
          <div className="mx-auto flex max-w-sm flex-col gap-2 min-[400px]:flex-row justify-center">
            <Link to="/login">
              <Button size="lg">Начать</Button>
            </Link>
            <Link to="/search">
              <Button size="lg" variant="outline">
                Просмотр товаров
              </Button>
            </Link>
          </div>
        </div>
      </section>
      
      {/* Footer */}
      <footer className="border-t bg-muted mt-auto">
        <div className="container flex flex-col gap-2 py-10 md:h-16 md:flex-row md:items-center md:justify-between md:py-0">
          <div className="text-center md:text-left">
            <p className="text-sm text-gray-500 dark:text-gray-400">
              © 2025 Умный Шоппинг. Все права защищены.
            </p>
          </div>
          <div className="flex items-center justify-center md:justify-end gap-4">
            <Link to="/search" className="text-sm text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-50">Товары</Link>
            <Link to="/map" className="text-sm text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-50">Карта маршрута</Link>
          </div>
        </div>
      </footer>
    </div>
  );
}
