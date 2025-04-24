FROM golang:1.23.5-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
WORKDIR /app/internal/services/payment_service
RUN CGO_ENABLED=0 GOOS=linux go build -o /payment-service ./cmd/main.go

FROM alpine:latest
COPY --from=builder /payment-service /app/
COPY --from=builder /app/internal/services/payment_service/config/local.yaml /app/config/

WORKDIR /app
EXPOSE 50053
CMD ["./payment-service"]