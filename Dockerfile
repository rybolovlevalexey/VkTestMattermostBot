# Используем официальный образ Go как базовый
FROM golang:1.23 AS builder

# Установка рабочей директории внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для скачивания зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем остальной исходный код
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./...

# Создаем конечный образ на базе Alpine Linux для минимизации размера
FROM alpine:latest

# Установка необходимых пакетов (например, tzdata для корректной работы времени)
RUN apk add --no-cache tzdata

# Создаем рабочую директорию
WORKDIR /app

# Копируем собранный бинарный файл из стадии сборки
COPY --from=builder /app/app .

# Указываем команду запуска
CMD ["./app"]

# Определяем порт, который будет открыт
EXPOSE 8080
