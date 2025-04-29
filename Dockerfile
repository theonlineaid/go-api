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


# FROM golang:1.24.2-alpine

# # Install necessary build tools
# RUN apk add --no-cache git

# # Install air
# RUN go install github.com/air-verse/air@latest

# # Set the working directory
# WORKDIR /app

# # Copy go.mod and go.sum first
# COPY go.mod go.sum ./

# # Download the dependencies
# RUN go mod tidy

# # Copy the rest of the app
# COPY . .

# # Expose port
# EXPOSE 8080

# # Use air for live reload
# CMD ["air"]
