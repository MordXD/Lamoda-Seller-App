FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create a non-root user
RUN addgroup -g 1001 appuser && \
    adduser -D -s /bin/sh -u 1001 -G appuser appuser

USER appuser

EXPOSE 8080

CMD ["./main"]