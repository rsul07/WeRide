FROM golang:1.23.6-alpine3.21 as builder

# Set destination for COPY
WORKDIR /room

# Download Go modules
COPY . .
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /room ./cmd/main.go

FROM alpine:3.21

WORKDIR /room

COPY --from=builder /room .

EXPOSE 8081
EXPOSE 50051

# Run
CMD ["./room"]