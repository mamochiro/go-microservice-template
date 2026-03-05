# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install wire
RUN go install github.com/google/wire/cmd/wire@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate wire code
RUN cd internal/app && wire

# Build the application
RUN go build -o main cmd/api/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

EXPOSE 8080

CMD ["./main"]
