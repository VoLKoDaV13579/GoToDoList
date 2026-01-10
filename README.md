# 📝 GoToDoList - Production-Ready REST API

> Современное приложение для управления задачами (To-Do List) на Go с использованием лучших практик разработки.

## 🚀 Особенности

- ✅ **Clean Architecture** - разделение на слои (handler → service → repository)
- ✅ **RESTful API** - стандартизированные эндпоинты
- ✅ **PostgreSQL** - надежная реляционная БД с GORM ORM
- ✅ **Graceful Shutdown** - корректное завершение работы
- ✅ **Structured Logging** - структурированное логирование через slog
- ✅ **Middleware** - логирование, recovery, CORS, request ID
- ✅ **Валидация данных** - проверка входящих данных
- ✅ **Docker Support** - контейнеризация приложения
- ✅ **Пагинация и фильтрация** - удобная работа со списками
- ✅ **Soft Delete** - безопасное удаление записей
- ✅ **Health Check** - эндпоинт для проверки работоспособности

## 📋 Требования

- Go 1.22+
- PostgreSQL 14+
- Docker и Docker Compose (опционально)
- Make (опционально)

## 🏗️ Архитектура проекта

```
GoToDoList/
├── cmd/
│   └── api/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── config/                  # Конфигурация приложения
│   ├── database/                # Подключение к БД
│   ├── handler/                 # HTTP обработчики (контроллеры)
│   ├── middleware/              # HTTP middleware
│   ├── model/                   # Модели данных иDTO
│   ├── repository/              # Слой работы с БД
│   └── service/                 # Бизнес-логика
├── pkg/
│   ├── logger/                  # Логгер
│   └── validator/               # Валидатор
├── docker-compose.yml           # Docker Compose конфигурация
├── Dockerfile                   # Dockerfile для сборки образа
├── Makefile                     # Удобные команды для разработки
├── .env.example                 # Пример переменных окружения
└── README.md                    # Документация
```

## 🚦 Быстрый старт

### Вариант 1: Docker Compose (рекомендуется)

1. Клонируйте репозиторий:
```bash
git clone https://github.com/VoLKoDaV13579/GoToDoList.git
cd GoToDoList
```

2. Запустите приложение:
```bash
make docker-up
```

Приложение будет доступно по адресу: `http://localhost:8080`

Adminer (веб-интерфейс БД): `http://localhost:8081`

### Вариант 2: Локальный запуск

1. Установите зависимости:
```bash
make deps
```

2. Запустите PostgreSQL (например, через Docker):
```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=todolist -p 5432:5432 -d postgres:16-alpine
```

3. Скопируйте `.env.example` в `.env` и настройте переменные окружения:
```bash
cp .env.example .env
```

4. Запустите приложение:
```bash
make run
```

## 📚 API Документация

### Base URL
```
http://localhost:8080/api/v1
```

### Эндпоинты

#### 🏥 Health Check
```http
GET /health
```

**Ответ:**
```json
{
  "status": "ok",
  "message": "API работает нормально"
}
```

---

#### ✨ Создать задачу
```http
POST /api/v1/todos
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "title": "Купить продукты",
  "description": "Молоко, хлеб, яйца",
  "status": "pending",
  "priority": 1,
  "due_date": "2024-12-31T23:59:59Z"
}
```

**Параметры:**
- `title` (обязательно) - название задачи (1-255 символов)
- `description` (опционально) - описание задачи (до 1000 символов)
- `status` (опционально) - статус: `pending`, `in_progress`, `completed` (по умолчанию: `pending`)
- `priority` (опционально) - приоритет: 0 (низкий), 1 (средний), 2 (высокий)
- `due_date` (опционально) - срок выполнения в формате ISO 8601

**Ответ (201 Created):**
```json
{
  "id": 1,
  "title": "Купить продукты",
  "description": "Молоко, хлеб, яйца",
  "status": "pending",
  "priority": 1,
  "due_date": "2024-12-31T23:59:59Z",
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T10:00:00Z"
}
```

---

#### 📋 Получить список задач
```http
GET /api/v1/todos?status=pending&priority=1&page=1&page_size=10
```

**Query параметры:**
- `status` (опционально) - фильтр по статусу
- `priority` (опционально) - фильтр по приоритету
- `page` (опционально) - номер страницы (по умолчанию: 1)
- `page_size` (опционально) - размер страницы (по умолчанию: 10, максимум: 100)

**Ответ (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Купить продукты",
      "description": "Молоко, хлеб, яйца",
      "status": "pending",
      "priority": 1,
      "due_date": "2024-12-31T23:59:59Z",
      "created_at": "2024-01-10T10:00:00Z",
      "updated_at": "2024-01-10T10:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

---

#### 🔍 Получить задачу по ID
```http
GET /api/v1/todos/:id
```

**Ответ (200 OK):**
```json
{
  "id": 1,
  "title": "Купить продукты",
  "description": "Молоко, хлеб, яйца",
  "status": "pending",
  "priority": 1,
  "due_date": "2024-12-31T23:59:59Z",
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T10:00:00Z"
}
```

---

#### ✏️ Обновить задачу
```http
PATCH /api/v1/todos/:id
Content-Type: application/json
```

**Тело запроса (все поля опциональны):**
```json
{
  "title": "Купить продукты и цветы",
  "status": "in_progress",
  "priority": 2
}
```

**Ответ (200 OK):**
```json
{
  "id": 1,
  "title": "Купить продукты и цветы",
  "description": "Молоко, хлеб, яйца",
  "status": "in_progress",
  "priority": 2,
  "due_date": "2024-12-31T23:59:59Z",
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T11:00:00Z"
}
```

---

#### 🔄 Обновить статус задачи
```http
PATCH /api/v1/todos/:id/status
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "status": "completed"
}
```

**Ответ (200 OK):**
```json
{
  "id": 1,
  "title": "Купить продукты",
  "status": "completed",
  ...
}
```

---

#### ❌ Удалить задачу
```http
DELETE /api/v1/todos/:id
```

**Ответ (200 OK):**
```json
{
  "message": "Задача успешно удалена"
}
```

---

### Коды ошибок

| Код | Описание |
|-----|----------|
| 400 | Неверный запрос (невалидные данные) |
| 404 | Задача не найдена |
| 500 | Внутренняя ошибка сервера |

**Пример ответа с ошибкой:**
```json
{
  "error": "validation_error",
  "message": "поле 'title' обязательно для заполнения"
}
```

## 🛠️ Makefile команды

```bash
make help           # Показать справку
make build          # Собрать бинарник
make run            # Запустить приложение
make test           # Запустить тесты
make docker-up      # Запустить в Docker
make docker-down    # Остановить Docker
make docker-logs    # Показать логи Docker
make fmt            # Форматировать код
make deps           # Установить зависимости
```

## 🔧 Конфигурация

Конфигурация осуществляется через переменные окружения. Смотрите `.env.example` для списка всех доступных параметров.

Основные переменные:
- `APP_PORT` - порт приложения (по умолчанию: 8080)
- `APP_DEBUG` - режим отладки (true/false)
- `DB_HOST` - хост PostgreSQL
- `DB_PORT` - порт PostgreSQL
- `DB_USER` - пользователь БД
- `DB_PASSWORD` - пароль БД
- `DB_NAME` - имя базы данных

## 🧪 Примеры использования

### Создание задачи через curl:

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Написать документацию",
    "description": "Создать подробный README",
    "priority": 2,
    "status": "in_progress"
  }'
```

### Получение всех задач:

```bash
curl http://localhost:8080/api/v1/todos?page=1&page_size=10
```

### Обновление статуса:

```bash
curl -X PATCH http://localhost:8080/api/v1/todos/1/status \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

## 🐳 Docker

### Сборка образа:
```bash
make docker-build
```

### Запуск через docker-compose:
```bash
docker-compose up -d
```

### Остановка:
```bash
docker-compose down
```

### Полная очистка (включая volumes):
```bash
docker-compose down -v
```

## 📦 Используемые библиотеки

- [Gin](https://github.com/gin-gonic/gin) - HTTP фреймворк
- [GORM](https://gorm.io/) - ORM для работы с БД
- [validator](https://github.com/go-playground/validator) - валидация данных
- [UUID](https://github.com/google/uuid) - генерация UUID

## 🏆 Best Practices

Проект следует следующим практикам:

- ✅ Clean Architecture (разделение на слои)
- ✅ Dependency Injection
- ✅ Graceful Shutdown
- ✅ Structured Logging
- ✅ Error Handling
- ✅ Input Validation
- ✅ Connection Pooling
- ✅ Health Checks
- ✅ Multi-stage Docker builds
- ✅ Non-root Docker user
- ✅ Комментарии на русском языке

## 📝 Лицензия

MIT License

## 👨‍💻 Автор

VoLKoDaV13579

---

**Примечание:** Это production-ready приложение, готовое к развертыванию. Для production окружения рекомендуется:
- Использовать HTTPS
- Настроить аутентификацию и авторизацию
- Добавить rate limiting
- Настроить мониторинг и алертинг
- Использовать managed PostgreSQL
- Настроить CI/CD
