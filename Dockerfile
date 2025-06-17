# Use the official Go image as the build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install git (needed for go modules)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o restic-browser main.go

# Use a minimal base image for the final stage
FROM alpine:latest

# Install restic and ca-certificates
RUN apk add --no-cache restic ca-certificates tzdata

# Create a non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/restic-browser .

# Copy templates directory
COPY --from=builder /app/templates ./templates

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Environment variables with defaults
ENV RESTIC_REPO=""
ENV RESTIC_PASSWORD=""
ENV PORT="8081"

# Expose the port
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/ || exit 1

# Start the application using environment variables
CMD ["sh", "-c", "./restic-browser -repo \"${RESTIC_REPO}\" -password \"${RESTIC_PASSWORD}\" -port \"${PORT}\""]
