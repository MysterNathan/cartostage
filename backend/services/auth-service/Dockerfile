# Build stage-service
FROM mysternathan.freeboxos.fr/golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy shared module first
COPY shared/ ./shared/

# Copy auth-service service
COPY services/auth-service/ ./services/auth-service/

WORKDIR /app/shared
RUN go mod download

WORKDIR /app/services/auth-service
RUN go mod download

# Build the auth-service service from cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage-service
FROM mysternathan.freeboxos.fr/alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary
COPY --from=builder /app/services/auth-service/main .

EXPOSE 80

CMD ["./main"]
