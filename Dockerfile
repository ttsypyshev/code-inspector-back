# Используем официальный образ Go
FROM golang:1.23.2-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Собираем бинарник
RUN go build -o main ./src/main.go

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["/app/main"]
