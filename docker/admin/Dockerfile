# Устанавливаем базовый образ
FROM golang:1.18-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR ./admin

# Обновляем основные зависимости среды
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY ./go.mod ./go.sum .
RUN go mod download

# Копируем все остальные файлы проекта внутрь контейнера
COPY ./cmd/admin		./cmd/admin
COPY ./config/admin	    ./config/admin
COPY ./internal/admin	./internal/admin
COPY ./pkg			    ./pkg

EXPOSE 5002

# Делаем последнюю билд проекта
RUN go build -o admin ./cmd/admin

# Устанавливаем переменные окружения
ARG config=./config/admin/config.json
ENV CONFIG_PATH=$config

# Указываем команду, которая будет запускаться при старте контейнера
CMD ./admin -config=$CONFIG_PATH