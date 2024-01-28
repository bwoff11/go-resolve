# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Enable cgo support
ENV CGO_ENABLED=1

# Build the Go application
RUN go build -o go-resolve .

# Set the entry point for the container
ENTRYPOINT ["./go-resolve"]
