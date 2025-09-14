# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy shared module first
COPY shared/ ./shared/

# Copy auth service
COPY services/auth/ ./services/auth/

# Set working directory to auth service
WORKDIR /app/services/auth

# Download dependencies
RUN go mod download

# Build the auth service from cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary
COPY --from=builder /app/services/auth/main .

EXPOSE 80

CMD ["./main"]
