
### Запуск (Docker Compose)

```bash
docker-compose up --build
```
- Соберет Docker образ приложения
- Запустит PostgreSQL
- Запустит микросервис
- Настроит все зависимости

Сервис будет доступен на `http://localhost:8080`

### Использование API

#### Добавить число (POST)

```bash
# JSON форматом
curl -X POST http://localhost:8080/numbers \
  -H "Content-Type: application/json" \
  -d '{"number": 3}'

# Query параметром
curl -X POST "http://localhost:8080/numbers?number=2"
```

#### Получить все числа (GET)

```bash
curl http://localhost:8080/numbers
```

Ответ:
```json
{
  "numbers": [1, 2, 3]
}
```

## Локальная разработка

### Установка зависимостей

```bash
go mod download
```

### Настройка базы данных

Создайте базу данных:

```sql
CREATE DATABASE numbersdb;
```

Или используйте переменную окружения `DATABASE_URL`:
```bash
export DATABASE_URL="postgres://user:password@localhost/numbersdb?sslmode=disable"
```

### Запуск

```bash
go run main.go
```

### Тесты

```bash
# Убедитесь, что PostgreSQL запущен и доступен
# Создайте тестовую БД:
createdb numbersdb_test

# Запустите тесты
go test -v
```

## Структура проекта

```
.
├── main.go           # Основной код приложения
├── main_test.go      # Тесты
├── go.mod            # Go модули
├── go.sum            # Зависимости
├── Dockerfile        # Docker образ приложения
├── docker-compose.yml # Docker Compose конфигурация
└── README.md         # Документация
```

## API Endpoints

### POST /numbers
Добавляет число в базу данных и возвращает отсортированный список всех чисел.

**Запрос:**
- JSON: `{"number": 3}`
- Query param: `?number=3`

**Ответ:**
```json
{
  "numbers": [1, 2, 3]
}
```

### GET /numbers
Возвращает отсортированный список всех сохраненных чисел.

**Ответ:**
```json
{
  "numbers": [1, 2, 3]
}
```

## Переменные окружения

- `DATABASE_URL` - URL подключения к PostgreSQL (по умолчанию: `postgres://postgres:postgres@localhost/numbersdb?sslmode=disable`)
- `PORT` - Порт для HTTP сервера (по умолчанию: `8080`)
