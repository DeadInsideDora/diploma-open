import argparse
import json
import requests
import sys

CARDS = {
    'products': ['Лента', 'Перекрёсток', 'Дикси', 'Магнит'],
    'nearby-products': [],
}

def main():
    parser = argparse.ArgumentParser(description="Вычисление эффективности построения продуктовой корзины на основе тестового кейса.")
    parser.add_argument('--json_path', required=True, help='Путь до JSON файла с описанием продуктовой корзины, скидочных кард, точки пользователя и радиуса')
    parser.add_argument('--service_url', required=True, help='Сетевой адрес сервиса optimizer, например http://localhost:8080')
    parser.add_argument('--radius', type=int, required=True, help='Радиус рассматриваемых магазинов')
    parser.add_argument('--exchange', type=int, required=True)
    parser.add_argument('--lat', type=float, required=True)
    parser.add_argument('--lon', type=float, required=True)

    args = parser.parse_args()

    try:
        with open(args.json_path, encoding='utf-8') as f:
            data = json.load(f)
    except Exception as e:
        print(f"Ошибка при загрузке JSON файла: {e}")
        sys.exit(1)

    urls = {
        path: f"{args.service_url.rstrip('/')}/{path}"
        for path in ['products', 'nearby-products']
    }

    response = dict()

    for (path, url) in urls.items():
        try:
            data['point'] = {
                'lat': args.lat,
                'lon': args.lon,
            }
            data['radius'] = args.radius
            data['exchange'] = args.exchange
            data['discount_cards'] = CARDS[path]
            response[path] = requests.post(url, json=data)
            print(f"POST {url} -> {response[path].status_code}")
            print(response[path].text)
        except Exception as e:
            print(f"Ошибка при POST запросе к {url}: {e}")

    print(f"products: {response['products']}")
    print(f"nearby-products: {response['nearby-products']}")

    print(f'Процент экономии по Cost: {(1 - (response["products"].json()["cost"] / response["nearby-products"].json()["cost"])) * 100.}')
    print(f'Процент экономии по Price: {(1 - (response["products"].json()["price"] / response["nearby-products"].json()["price"])) * 100.}')

if __name__ == '__main__':
    main()
