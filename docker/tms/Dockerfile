# Устанавливаем базовый образ
FROM golang:1.18-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR ./tms

# Обновляем основные зависимости среды
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY ./go.mod ./go.sum .
RUN go mod download

# Копируем все остальные файлы проекта внутрь контейнера
COPY ./cmd/tms		./cmd/tms
COPY ./config/tms	./config/tms
COPY ./internal/tms	./internal/tms
COPY ./pkg			./pkg 

EXPOSE 5001

# Делаем последнюю билд проекта
RUN go build -o tms ./cmd/tms

# Устанавливаем переменные окружения
ARG config=./config/tms/config.json
ENV CONFIG_PATH=$config

# Указываем команду, которая будет запускаться при старте контейнера
CMD ./tms -config=$CONFIG_PATH