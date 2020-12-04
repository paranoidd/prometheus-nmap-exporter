# First step - build the binary
FROM golang:1.14-buster as builder

# Create build directory
RUN mkdir -p /tmp/build
WORKDIR /tmp/build

# Copy source and run build
COPY . ./
RUN go build

# Second step - prepare runtime
FROM debian:buster-slim

# Install nmap
RUN mkdir -p /app && \
    apt-get update && \
    apt-get install -y --no-install-recommends nmap && \
    apt-get autoclean && \
    apt-get autoremove -y && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Change working directory
WORKDIR /app

# Copy built binary from first step
COPY --from=builder /tmp/build/prometheus-nmap-exporter ./

# Docker settings
EXPOSE     8080
ENTRYPOINT [ "/app/prometheus-nmap-exporter" ]
