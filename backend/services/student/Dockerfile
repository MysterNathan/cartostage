# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy shared module first
COPY shared/ ./shared/

# Copy stage service
COPY services/student/ ./services/student/

# Initialize go modules in shared
WORKDIR /app/shared
RUN go mod download

# Build stage service
WORKDIR /app/services/student
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary
COPY --from=builder /app/services/student/main .

EXPOSE 80

CMD ["./main"]
