# Этап сборки
FROM golang:1.23.5-alpine AS builder

# Устанавливаем зависимости
RUN apk add --no-cache git

# Рабочая директория
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь проект
COPY . .

# Переходим в директорию сервиса
WORKDIR /app/internal/services/room_service

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /room-service ./cmd/main.go

# Этап запуска
FROM alpine:latest

# Копируем бинарник из этапа сборки
COPY --from=builder /room-service /room-service

EXPOSE 50051
EXPOSE 8081

# Запускаем сервис
CMD ["/room-service"]