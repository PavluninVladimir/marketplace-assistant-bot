# Используем официальный образ Golang
FROM golang:1.21.5-alpine3.18

ENV TOKEN_TELEGRAM_BOT=121212121212
#ENV CGO_ENABLED=1

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /go/src/app

# Копируем содержимое текущей директории внутрь контейнера
COPY . .

#RUN apk add build-base

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev
# Загружаем зависимости проекта
RUN go install ./...

# Собираем приложение
RUN GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app ./...

# Экспонируем порт, на котором работает наше приложение
EXPOSE 8080

# Определяем команду для запуска приложения при старте контейнера
CMD ["./app"]
