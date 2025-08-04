# Используем версию Go, соответствующую вашему go.mod
FROM golang:1.24-alpine as builder

WORKDIR /app

# Установка зависимостей
RUN apk add --no-cache git gcc musl-dev

# Копируем сначала только go.mod и go.sum для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Установка swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Копируем остальные файлы
COPY . .

# Генерируем документацию (с явным указанием директорий)
RUN swag init 

# Собираем приложение
RUN go build -o main ./main.go

# Финальный образ
FROM alpine:3.18
WORKDIR /app

# Копируем бинарник и ресурсы
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/web ./web
COPY --from=builder /app/main.go .
COPY --from=builder /app/internal ./internal

EXPOSE 8080
CMD ["./main"]
