FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o pinger ./pinger/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/pinger /usr/local/bin/pinger

RUN chmod +x /usr/local/bin/pinger

EXPOSE 8080

CMD ["pinger"]
