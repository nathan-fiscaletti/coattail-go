# Use the official Golang image with Alpine Linux for smaller size
FROM golang:1.23.2-alpine

# Install build tools required for CGO
RUN apk add --no-cache build-base

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules files first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go binary for the "web-service" from the /cmd directory
RUN go build -o web-service ./cmd

# Expose required ports
EXPOSE 5423
EXPOSE 8082
EXPOSE 8083

# Command to run the application
CMD ["./web-service"]