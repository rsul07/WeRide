FROM golang:1.23.6-alpine3.21 as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY . .
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /order ./cmd/main.go

FROM alpline:3.21

WORKDIR /app

COPY --from=builder /order .

EXPOSE 8081
EXPOSE 50051

# Run
CMD ["./order"]