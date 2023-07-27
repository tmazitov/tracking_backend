# Используем официальный образ PostgreSQL как базовый образ
FROM postgis/postgis:15-3.3-alpine

# Опционально, установить дополнительные переменные окружения для PostgreSQL
# ENV POSTGRES_USER myuser
# ENV POSTGRES_PASSWORD mypassword
# ENV POSTGRES_DB mydb

# RUN apt-get update \
#     && apt-get install -y --no-install-recommends \
#         postgresql-12-postgis-3 \
#         postgresql-12-postgis-3-scripts

# Опционально, скопировать файлы инициализации и другие скрипты, если это необходимо
COPY ./scripts/db/aaa/*				/docker-entrypoint-initdb.d/
COPY ./scripts/db/tms/*			/docker-entrypoint-initdb.d/

# Опционально, указать порт, на котором будет слушать PostgreSQL
EXPOSE 5432

# Опционально, указать другие инструкции или команды, если это необходимо
# ...

# Команда ENTRYPOINT и CMD наследуются от базового образа и обычно не требуют изменений.