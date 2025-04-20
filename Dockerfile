# Используем официальный образ Go
FROM golang:1.24.2-alpine as builder

# Создаем рабочую директорию
WORKDIR /app

# Копируем модули и зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/web

# Используем минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из builder
COPY --from=builder /app/main /app/main
COPY ./configs ./configs

# Открываем порт для grpc
EXPOSE 14500

# Запускаем приложение
ENTRYPOINT ["/app/main"]