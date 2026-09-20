# Use the official Go image to build Envoy.
FROM golang:1.27 AS builder

# Work inside /app.
WORKDIR /app

# Copy dependency files first.
COPY go.mod go.sum ./

# Download Go dependencies.
RUN go mod download

# Copy the Envoy source code.
COPY . .

# Build the Envoy CLI binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o envoy ./cmd/envoy


# Use a small image for the final application.
FROM debian:bookworm-slim

# Work inside /app.
WORKDIR /app

# Copy only the compiled Envoy binary.
COPY --from=builder /app/envoy /usr/local/bin/envoy

# Start Envoy when the container runs.
ENTRYPOINT ["envoy"]

# Default command.
CMD ["--help"]
