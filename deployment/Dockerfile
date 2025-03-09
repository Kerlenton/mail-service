# Use Go base image
FROM golang:1.24.1-alpine

# Set working directory inside container
WORKDIR /app

# Copy go mod and sum files
COPY ../go.mod ../go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY .. .

# Build the application
RUN go build -o mail-service cmd/main.go

# Expose port
EXPOSE 8080

# Run the application
CMD ["./mail-service"]
# Use Go base image