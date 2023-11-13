FROM redis:7.2.3-alpine3.18

RUN adduser \
  --no-create-home \
  --system \
  --shell /usr/sbin/nologin \
  messaging-runner

USER messaging-runner
