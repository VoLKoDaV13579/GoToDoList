.PHONY: help build run test clean docker-build docker-up docker-down docker-logs migrate-up migrate-down

# Переменные
APP_NAME=todolist-api
DOCKER_IMAGE=todolist-api:latest
MAIN_PATH=./cmd/api/main.go

## help: Показать справку по доступным командам
help:
	@echo "Доступные команды:"
	@echo "  make build          - Собрать бинарник приложения"
	@echo "  make run            - Запустить приложение локально"
	@echo "  make test           - Запустить тесты"
	@echo "  make clean          - Очистить сборочные артефакты"
	@echo "  make docker-build   - Собрать Docker образ"
	@echo "  make docker-up      - Запустить приложение в Docker"
	@echo "  make docker-down    - Остановить Docker контейнеры"
	@echo "  make docker-logs    - Показать логи Docker контейнеров"
	@echo "  make lint           - Запустить линтер (golangci-lint)"
	@echo "  make fmt            - Форматировать код"
	@echo "  make deps           - Установить зависимости"

## build: Собрать бинарник приложения
build:
	@echo "Сборка приложения..."
	go build -o bin/$(APP_NAME) $(MAIN_PATH)
	@echo "Сборка завершена: bin/$(APP_NAME)"

## run: Запустить приложение локально (требуется запущенная PostgreSQL)
run:
	@echo "Запуск приложения..."
	go run $(MAIN_PATH)

## test: Запустить тесты
test:
	@echo "Запуск тестов..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Покрытие тестами сохранено в coverage.html"

## clean: Очистить сборочные артефакты
clean:
	@echo "Очистка..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Очистка завершена"

## docker-build: Собрать Docker образ
docker-build:
	@echo "Сборка Docker образа..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "Docker образ собран: $(DOCKER_IMAGE)"

## docker-up: Запустить приложение в Docker с PostgreSQL
docker-up:
	@echo "Запуск Docker контейнеров..."
	docker-compose up -d
	@echo "Приложение запущено на http://localhost:8080"
	@echo "Adminer доступен на http://localhost:8081"

## docker-down: Остановить Docker контейнеры
docker-down:
	@echo "Остановка Docker контейнеров..."
	docker-compose down
	@echo "Контейнеры остановлены"

## docker-logs: Показать логи Docker контейнеров
docker-logs:
	docker-compose logs -f

## docker-restart: Перезапустить Docker контейнеры
docker-restart: docker-down docker-up

## docker-clean: Остановить контейнеры и удалить volumes
docker-clean:
	@echo "Остановка контейнеров и удаление данных..."
	docker-compose down -v
	@echo "Контейнеры остановлены, данные удалены"

## lint: Запустить линтер
lint:
	@echo "Запуск линтера..."
	golangci-lint run ./...

## fmt: Форматировать код
fmt:
	@echo "Форматирование кода..."
	go fmt ./...
	gofmt -s -w .
	@echo "Форматирование завершено"

## deps: Установить зависимости
deps:
	@echo "Установка зависимостей..."
	go mod download
	go mod tidy
	@echo "Зависимости установлены"

## mod-update: Обновить зависимости
mod-update:
	@echo "Обновление зависимостей..."
	go get -u ./...
	go mod tidy
	@echo "Зависимости обновлены"

.DEFAULT_GOAL := help
