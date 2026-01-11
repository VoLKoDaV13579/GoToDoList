# Тестирование GoToDoList

Этот документ содержит информацию о запуске и работе с тестами проекта.

## Обзор покрытия

Проект имеет полное покрытие unit-тестами для всех основных слоев:

- **Service Layer**: 100% покрытие
- **Validator**: 94.1% покрытие
- **Repository**: 90.2% покрытие
- **Handler**: 79.2% покрытие

## Структура тестов

```
.
├── internal/
│   ├── handler/
│   │   ├── mocks/                    # Моки для сервисов
│   │   │   └── mock_todo_service.go
│   │   └── todo_handler_test.go      # Тесты HTTP handlers
│   ├── repository/
│   │   └── todo_repository_test.go   # Тесты работы с БД
│   └── service/
│       ├── mocks/                    # Моки для репозиториев
│       │   └── mock_todo_repository.go
│       └── todo_service_test.go      # Тесты бизнес-логики
└── pkg/
    └── validator/
        └── validator_test.go         # Тесты валидации
```

## Запуск тестов

### Запустить все тесты

```bash
go test ./...
```

### Запустить тесты с подробным выводом

```bash
go test ./... -v
```

### Запустить тесты конкретного пакета

```bash
# Repository tests
go test ./internal/repository -v

# Service tests
go test ./internal/service -v

# Handler tests
go test ./internal/handler -v

# Validator tests
go test ./pkg/validator -v
```

### Запустить конкретный тест

```bash
go test ./internal/service -v -run TestTodoService_CreateTodo
```

## Проверка покрытия

### Генерация отчета о покрытии

```bash
go test ./... -coverprofile=coverage.out
```

### Просмотр покрытия в консоли

```bash
go tool cover -func=coverage.out
```

### Генерация HTML отчета

```bash
go tool cover -html=coverage.out -o coverage.html
```

Откройте `coverage.html` в браузере для интерактивного просмотра покрытия.

### Просмотр покрытия конкретного пакета

```bash
go test ./internal/service -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Используемые библиотеки для тестирования

- **[testify](https://github.com/stretchr/testify)** - удобные assertions и моки
- **[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)** - мокирование SQL запросов
- **httptest** (стандартная библиотека) - тестирование HTTP handlers

## Особенности тестов

### Repository Layer

- Использует `sqlmock` для изоляции от реальной базы данных
- Тестирует все CRUD операции
- Проверяет корректность SQL запросов
- Тестирует обработку ошибок БД

### Service Layer

- Использует моки репозиториев
- Тестирует бизнес-логику приложения
- Проверяет валидацию и трансформацию данных
- Тестирует пагинацию и фильтрацию
- 100% покрытие всех методов

### Handler Layer

- Использует `httptest` для имитации HTTP запросов
- Тестирует все эндпоинты API
- Проверяет статус-коды и форматы ответов
- Тестирует валидацию входных данных
- Тестирует обработку ошибок

### Validator

- Тестирует все правила валидации
- Проверяет форматирование ошибок
- Тестирует граничные случаи

## Continuous Integration

Тесты могут быть легко интегрированы в CI/CD pipeline:

```yaml
# Пример для GitHub Actions
- name: Run tests
  run: go test ./... -v -coverprofile=coverage.out

- name: Check coverage
  run: |
    go tool cover -func=coverage.out
    coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    if (( $(echo "$coverage < 80" | bc -l) )); then
      echo "Coverage is below 80%"
      exit 1
    fi
```

## Лучшие практики

1. **Запускайте тесты перед коммитом**
   ```bash
   go test ./...
   ```

2. **Проверяйте покрытие**
   ```bash
   go test ./... -cover
   ```

3. **Используйте table-driven tests** для множественных тест-кейсов

4. **Изолируйте тесты** - используйте моки для внешних зависимостей

5. **Тестируйте граничные случаи** и обработку ошибок

## Добавление новых тестов

При добавлении нового функционала следуйте этим правилам:

1. Создайте тестовый файл `*_test.go` рядом с основным файлом
2. Используйте именование `Test<FunctionName>` для тестовых функций
3. Группируйте связанные тесты с помощью `t.Run()`
4. Используйте моки для изоляции от внешних зависимостей
5. Стремитесь к покрытию минимум 80%

### Пример структуры теста

```go
func TestMyFunction(t *testing.T) {
    t.Run("successful case", func(t *testing.T) {
        // Arrange
        // Act
        // Assert
    })

    t.Run("error case", func(t *testing.T) {
        // Arrange
        // Act
        // Assert
    })
}
```

## Troubleshooting

### Тесты падают с ошибкой подключения к БД

Тесты используют моки и не требуют реальной БД. Убедитесь, что вы используете правильные импорты и моки настроены корректно.

### Тесты проходят локально, но падают в CI

Проверьте версию Go и зависимости. Убедитесь, что `go.mod` и `go.sum` закоммичены.

### Низкое покрытие в новом коде

Добавьте больше тест-кейсов, особенно для:
- Граничных случаев
- Обработки ошибок
- Различных входных данных

## Полезные команды

```bash
# Очистка кеша тестов
go clean -testcache

# Запуск тестов с race detector
go test ./... -race

# Запуск тестов с таймаутом
go test ./... -timeout 30s

# Бенчмарки (если добавлены)
go test ./... -bench=.
```
