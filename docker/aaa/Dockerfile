# Устанавливаем базовый образ
FROM golang:1.18-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR ./aaa

# Обновляем основные зависимости среды
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY ./go.mod ./go.sum .
RUN go mod download

# Копируем все остальные файлы проекта внутрь контейнера
COPY ./cmd/aaa		./cmd/aaa
COPY ./config/aaa	./config/aaa
COPY ./internal/aaa	./internal/aaa
COPY ./pkg			./pkg

EXPOSE 5000

# Делаем последнюю билд проекта
RUN go build -o aaa ./cmd/aaa

# Устанавливаем переменные окружения
ARG config=./config/aaa/config.json
ENV CONFIG_PATH=$config

# Указываем команду, которая будет запускаться при старте контейнера
CMD ./aaa -config=$CONFIG_PATH