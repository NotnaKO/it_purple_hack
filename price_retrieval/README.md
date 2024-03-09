# Price Retrieval service

В файле `server.go` создается http-сервер для обработки запросов.

### HTTP handlers
`/retrieve`:
по `location_id`, `microcategory_id`, `user_id` ищет цену для данного пользователя.

### Сборка и запуск
```bash
go build && go run .
```

### Пример запроса
```bash
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET 'http://localhost:8080/retrieve?location_id=123&microcategory_id=456&user_id=123'
```

### API

`LocationInfo` -- класс с информацией об http-запросе.
`Handler` -- класс с основным logger'ом и функциями обработки http-запросов
`Retriever` -- класс с основной логикой поиска (TODO)
