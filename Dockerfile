# Build stage
FROM golang:1.25.0-alpine3.21 AS builder

WORKDIR /app

# Install wire and swag
RUN go install github.com/google/wire/cmd/wire@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate wire code and swagger docs
RUN cd internal/app && wire
RUN swag init -g cmd/api/main.go

# Build the application
RUN go build -o main cmd/api/main.go

# Run stage
FROM alpine:3.21.3

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/docs ./docs

EXPOSE 3003

CMD ["./main"]
