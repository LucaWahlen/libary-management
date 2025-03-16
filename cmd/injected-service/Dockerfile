FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o injected-service ./cmd/injected-service/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

COPY --from=builder /app/injected-service .

EXPOSE 8080

CMD ["./injected-service"]