# Многоступенчатая сборка (multi-stage build) для оптимизации размера образа

# Стадия 1: Сборка приложения
FROM golang:1.25-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go mod файлы для кеширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник с оптимизациями для production
# -ldflags="-s -w" уменьшает размер бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/bin/api \
    ./cmd/api/main.go

# Стадия 2: Финальный образ (только runtime)
FROM alpine:latest

# Устанавливаем необходимые зависимости
RUN apk --no-cache add ca-certificates tzdata

# Создаем непривилегированного пользователя для безопасности
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

# Устанавливаем рабочую директорию
WORKDIR /home/appuser

# Копируем бинарник из builder стадии
COPY --from=builder /app/bin/api .

# Меняем владельца файлов
RUN chown -R appuser:appgroup /home/appuser

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт приложения
EXPOSE 8080

# Healthcheck для проверки работоспособности контейнера
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Запускаем приложение
CMD ["./api"]
