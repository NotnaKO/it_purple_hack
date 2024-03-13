# Price management service

В файле `server.go` создается http-сервер для обработки запросов.

## API
`Handler` -- класс с основным logger'ом и функциями обработки http-запросов

`HttpGetRequestInfo` -- класс с информацией о get http-запросе.

`HttpSetRequestInfo` -- класс с информацией о set http-запросе.

`PriceManager` -- класс с основной логикой взаимодействия с базой данных

## Создание базы данных, сборка и запуск

Создание таблиц происходит автоматически, нужно лишь установить правильную конфигурацию.

\
Помощь

```bash
go build
./price_management --help
```

Сборка и запуск
```bash
go build
./price_management -config_path=../config/price_management.yaml
```

## Примеры запросов
### Set price
```bash
curl -X POST "http://localhost:8080/set_price?location_id=1&microcategory_id=1&data_base_id=1&price=12.99"
```

### Get price
```bash
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET 'http://localhost:8080/get_price?location_id=1&microcategory_id=1&data_base_id=1'
```
