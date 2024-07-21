# Первый этап: сборка
FROM golang:1.22 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o app ./cmd/main.go

# Второй этап: создание минимального образа
FROM alpine:latest

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/app .

RUN chmod +x /usr/src/app/app

CMD ["./app"]
