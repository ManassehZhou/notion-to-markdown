FROM golang:1.21-alpine AS builder

# Install git and ca-certificates for building
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o notion-to-markdown main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests to Notion API
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/notion-to-markdown .

# Copy default config file
COPY --from=builder /app/config/notion-to-hugo.yaml ./config/

# Copy entrypoint script
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]