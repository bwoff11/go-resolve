# Builder stage
FROM golang:latest as builder

WORKDIR /build
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./main

# Install ca-certificates
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

# Final stage
FROM scratch

# Copy the built binary and ca-certificates
WORKDIR /app
COPY --from=builder /build/main ./main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the configuration file
COPY ./config.yaml ./config.yaml

# Set the port and entrypoint
EXPOSE 1053
ENTRYPOINT ["./main"]
