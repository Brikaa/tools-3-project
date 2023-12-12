FROM mysql:8.2.0

RUN adduser \
  --no-create-home \
  --system \
  --shell /usr/sbin/nologin \
  db-user

COPY ./schema.sql /docker-entrypoint-initdb.d/
USER db-user
