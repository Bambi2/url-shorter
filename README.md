# HTTP API По Созданию Сокращенных Ссылок

## Запуск

Хранилище Postgres
```
make run-postgres
```
Хранилище Redis
```
make run-redis
```

## Тестирование
```
make test
```

## Примеры запросов
### Сгенерировать коротрий URL
<img width="1226" alt="image" src="https://user-images.githubusercontent.com/53175260/206856365-32f5c837-8615-4e92-9b0f-643b316e49ff.png">

### Получить оригинальный URL
<img width="1222" alt="image" src="https://user-images.githubusercontent.com/53175260/206856395-88e1784f-d0c3-4823-8b63-6bdf8ca24c55.png">

### Несуществующий короткий URL
<img width="1223" alt="image" src="https://user-images.githubusercontent.com/53175260/206856419-f0f49429-3025-47aa-a652-86c2525c2fcc.png">

### Невалидная ссылка
<img width="1238" alt="image" src="https://user-images.githubusercontent.com/53175260/206856520-48e62721-3484-4c35-8b30-0144dbbe578c.png">

## Комментарий
Более правильным решением была бы реализация лимита по времени на существование ссылки, но в задании сказано, что на вход только ссылка
