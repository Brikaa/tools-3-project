FROM mysql:8.2.0

RUN adduser \
  --no-create-home \
  --system \
  --shell /usr/sbin/nologin \
  db-user

COPY ./initial.sql /initial.sql
USER db-user
