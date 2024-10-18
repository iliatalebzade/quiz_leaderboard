# Step 1: Build the Go application
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to cache dependencies
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Step 2: Create a lightweight container to run the app
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /root/

# Copy the compiled Go binary from the builder image
COPY --from=builder /app/main .

# Expose the port the app will run on
EXPOSE 8000

# Command to run the Go application
CMD ["./main"]
