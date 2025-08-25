# Multi-stage build: build the Go backend first
FROM golang:1.24-alpine AS builder

# Install git (needed for go mod download with private repos)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/. .

# Build the binary
RUN go build -o ideb-app .

# Final stage: create the minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -s /bin/sh idebuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/ideb-app .

# Copy frontend files
COPY frontend/ ../frontend/

# Copy sample data
COPY memory-bank/input.json ./memory-bank/

# Create the database directory
RUN mkdir -p ./data

# Expose port
EXPOSE 8080

# Change ownership to non-root user
RUN chown -R idebuser:idebuser /app

# Switch to non-root user
USER idebuser

# Environment variables
ENV DB_PATH=./data/ideb.db
ENV INPUT_JSON_PATH=./memory-bank/input.json
ENV SERVER_PORT=8080

# Run the binary
ENTRYPOINT ["./ideb-app"]