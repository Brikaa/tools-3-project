FROM golang:1.21.3-alpine3.18 AS build

WORKDIR /app
COPY ["go.mod", "go.sum", "./"]
RUN go mod download
COPY . .
RUN ["go", "build", "-o", "bin/backend"]

FROM alpine:3.18

WORKDIR /app
COPY --from=build /app/bin/backend /app
CMD ["/app/backend"]
