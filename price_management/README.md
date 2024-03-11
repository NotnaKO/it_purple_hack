# Price management service

В файле `server.go` создается http-сервер для обработки запросов.

## API
`Handler` -- класс с основным logger'ом и функциями обработки http-запросов

`HttpGetRequestInfo` -- класс с информацией о get http-запросе.

`HttpSetRequestInfo` -- класс с информацией о set http-запросе.

`PriceManager` -- класс с основной логикой взаимодействия с базой данных

## Создание базы данных, сборка и запуск

Создаем базу данных
```bash
sudo su - postgres
createdb price_management
exit
```

Создаем таблицу
```SQL
CREATE DATABASE price_management;

CREATE TABLE price_matrix (
    location_id BIGINT NOT NULL,
    microcategory_id BIGINT NOT NULL,
    price BIGINT NOT NULL,
    PRIMARY KEY (location_id, microcategory_id)
);
```

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
