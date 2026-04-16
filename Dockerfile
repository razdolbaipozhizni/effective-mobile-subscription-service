FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

# Финальный минимальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/main .

# Копируем миграции и swagger-документацию
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

# Порт, на котором будет работать приложение
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]