# Build stage
FROM golang:1.26.1-alpine AS builder

# Install build dependencies for CGO sqlite3
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files from the correct location
COPY src/go.mod src/go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . ./

# Set working directory to where the source code is
WORKDIR /app/src

# Build the application with CGO enabled for sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/src/main .

# Copy any necessary environment files
COPY .env* ./

RUN mkdir -p /var/log/leszmonitor

# Expose port 7001
EXPOSE 7001

# Command to run the executable
CMD ["./main"]