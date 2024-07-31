# Kafka Microservice

Этот микросервис на Go принимает сообщения через HTTP API, сохраняет их в PostgreSQL, а затем отправляет в Kafka для дальнейшей обработки. Обработанные сообщения помечаются, и сервис предоставляет API для получения статистики по обработанным сообщениям.

## Установка и запуск

### Требования
- Go 1.20+
- Docker и Docker Compose
- PostgreSQL

### Запуск

1. Клонируйте репозиторий:

   ```sh
   git clone github.com/K11MMILK/mess
   cd mess
Создайте файл .env с вашими переменными окружения:

2. Создайте файл .env с вашими переменными окружения:

    DB_CONN_STR=postgres://user:password@postgres:5432/messages?sslmode=disable
    KAFKA_BROKERS=kafka:9092
    KAFKA_TOPIC=messages
    KAFKA_GROUP_ID=message_group


3. Запустите Docker Compose:

    docker-compose up --build

4. Сервис будет доступен по адресу http://147.45.235.195:8080.

API
POST /messages
Принимает сообщения в формате JSON.
Пример запроса:
{
    "content": "Hello, Kafka!"
}

GET /stats
Возвращает количество обработанных сообщений.

Пример ответа:
{
    "processed_messages": 1
}
