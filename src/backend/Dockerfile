FROM golang:1.21.3-alpine3.18 AS build

WORKDIR /app
COPY ["go.mod", "go.sum", "./"]
RUN go mod download
COPY . .
RUN ["go", "build", "-o", "bin/backend"]

FROM alpine:3.18

RUN adduser \
  --no-create-home \
  --system \
  --shell /usr/sbin/nologin \
  backend-runner

WORKDIR /app
COPY --from=build /app/bin/backend /app

USER backend-runner
CMD ["/app/backend"]
