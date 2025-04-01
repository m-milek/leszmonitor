# Build stage
FROM golang:1.24-alpine AS builder

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

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/src/main .

# Copy any necessary environment files
COPY .env* ./

# Create directories for logs if needed
RUN mkdir -p logs

# Expose port 7001
EXPOSE 7001

# Command to run the executable
CMD ["./main"]