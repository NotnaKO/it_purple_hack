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
Нужно скачать `redis`, запустить через `systemctl` `redis.service`.

```bash
sudo apt-get install redis
sudo systemctl enable redis
sudo systemctl start redis
```

В файле `/etc/redis/redis.conf` есть строчка `port [port_num]` с портом `redis` сервиса. Его нужно добавить в `../config/price_retrieval.yaml`

Необходимо также установить библиотеки `prometheus` и `grafana` для мониторинга. В файле `/etc/prometheus/prometheus.yml` нужно добавить:
```
scrape_configs:
  ...

  - job_name: 'price_retriever'
    static_configs:
      - targets: ['localhost:7020']

  - job_name: 'price_manager'
    static_configs:
      - targets: ['localhost:8080']
```

В папке `data` лежит пример дэшборда (`dashboard.json`), где показаны графики с текущим RPS и долей cache misses.

```bash
go build
./price_retrieval -config_path=../config/price_retrieval.yaml
```

### Пример запроса
```bash
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET 'http://localhost:7020/retrieve?location_id=1&microcategory_id=1&user_id=1'
```

### API

`LocationInfo` -- класс с информацией об http-запросе.\
`Handler` -- класс с основным logger'ом и функциями обработки http-запросов\
`Retriever` -- класс с основной логикой поиска(описана в main README)
