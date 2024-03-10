# Price Retrieval service

В файле `server.go` создается http-сервер для обработки запросов.

### HTTP handlers
`/retrieve`:
по `location_id`, `microcategory_id`, `user_id` ищет цену для данного пользователя.

Помощь

```bash
go build
./price_retrieval --help
```

### Сборка и запуск
```bash
go build
./price_retrieval -config_path=../config/price_retrieval.yaml
```

### Пример запроса
```bash
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET 'http://localhost:8080/retrieve?location_id=123&microcategory_id=456&user_id=123'
```

### API

`LocationInfo` -- класс с информацией об http-запросе.
`Handler` -- класс с основным logger'ом и функциями обработки http-запросов
`Retriever` -- класс с основной логикой поиска (TODO)
