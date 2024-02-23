# Builder stage
FROM golang:1.18 as builder

# Create app directory
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o dnsproxy .

# Use `setcap` to allow the binary to bind to well-known ports as a non-root user
# Note: This step might fail in environments where `setcap` is not available or has restricted functionality.
# If you encounter issues, you may need to adjust your base image or set capabilities at runtime or deployment.
RUN apt-get update && apt-get install -y libcap2-bin && setcap 'cap_net_bind_service=+ep' /app/dnsproxy

# Final stage
FROM gcr.io/distroless/base-debian11

# Copy the binary from the builder stage
COPY --from=builder /app/dnsproxy /app/dnsproxy

# Non-root user configuration
USER nonroot:nonroot

# Expose DNS port
EXPOSE 53/udp
EXPOSE 53/tcp

# Run the application
ENTRYPOINT ["/app/dnsproxy"]
