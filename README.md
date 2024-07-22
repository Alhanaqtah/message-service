# Message Service

## Подключение

Сервис находится по адресу: `http://94.228.112.24:8080`

### Запросы через curl

#### Для отправки сообщения:

```sh
curl --location 'http://94.228.112.24:8080/messages' \
--header 'Content-Type: application/json' \
--data '{
    "content": "some awesome message"
}'
```

#### Для получения статистики:

```sh
curl --location 'http://94.228.112.24:8080/messages/stats'
```

### Запросы через Postman

Необходимо импортировать коллекцию, файл которой находится в корне репозитория по пути `docs/postman`. После импорта коллекции можно будет делать запросы к серверу.

### Запуск проекта локально

Если вы хотите запустить проект локально, выполните следующие шаги:

1. **Клонирование репозитория:**

   ```sh
   git clone https://github.com/your-repo/message-service.git
   cd message-service
   ```

2. **Настройка переменных окружения:**

   Создайте файл `.env` в корне проекта и добавьте следующие переменные:

   ```env
    ENV=prod #local/env/prod
    
    SERVER_HOST=message-service
    SERVER_PORT=8080
    SERVER_TIMEOUT=4
    
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_HOST=database
    POSTGRES_PORT=5432
    POSTGRES_DB=postgres
    
    KAFKA_BROKERS=kafka:9092
    KAFKA_TOPIC=messages
   ```

3. **Сборка и запуск Docker контейнеров:**

   ```sh
   docker-compose up --build
   ```
   
## Документация API

### Сохранение сообщения:

```http
POST /messages
Content-Type: application/json

{
  "content": "Текст сообщения"
}
```

#### Предполагаемый ответ:

```json
{
    "status": "Ok",
    "message": "Message saved successfully"
}
```

### Получение статистики:

```http
GET /messages/stats
```

#### Предполагаемый ответ:

```json
{
    "total_messages": 102724,
    "processed_messages": 102724,
    "average_processing_time": 273
}
```
