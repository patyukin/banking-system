FROM golang:1.21.2-alpine AS builder

COPY . /app
WORKDIR /app

RUN go mod download
# RUN go mod tidy
RUN go build -o ./bin/auth_server cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/bin/auth_server .
ADD .env .

CMD ["./auth_server", "-config-path", ".env"]
