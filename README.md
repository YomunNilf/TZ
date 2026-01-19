# Numbers Service

Микросервис на Go с REST эндпоинтом для сохранения и получения отсортированных чисел.

## Описание

Сервис принимает числа через REST API, сохраняет их в PostgreSQL и возвращает отсортированный список всех сохраненных чисел.

### Примеры использования:

1. Отправили число 3 => Записали в БД => Вернули `[3]`
2. Отправили число 2 => Записали в БД => Вернули `[2, 3]`
3. Отправили число 1 => Записали в БД => Вернули `[1, 2, 3]`

## Быстрый старт

### Запуск одной командой (Docker Compose)

```bash
docker-compose up --build
```

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

### Требования

- Go 1.21+
- PostgreSQL 15+

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

## Технологии

- **Go 1.21** - основной язык
- **PostgreSQL** - база данных
- **Docker** - контейнеризация
- **Docker Compose** - оркестрация

