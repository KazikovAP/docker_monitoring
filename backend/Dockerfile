FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY backend ./backend

RUN CGO_ENABLED=0 GOOS=linux go build -o docker_monitoring ./backend/cmd/dock_mon

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN adduser -D appuser
USER appuser

WORKDIR /home/appuser

COPY --from=builder /app/docker_monitoring .

ENV DB_HOST=localhost \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=password \
    DB_NAME=pinger \
    SERVER_PORT=8080

EXPOSE 8080

CMD ["./docker_monitoring"]
