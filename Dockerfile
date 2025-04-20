# Use a more recent Go version, 1.24.2 or above
FROM golang:1.24.2-alpine

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["go", "run", "main.go"]
