# --- Stage 1: Build the Application ---
FROM golang:alpine AS builder

WORKDIR /app

# 1. Download dependencies first (caching layer)
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy source code and build
COPY . .
# CGO_ENABLED=0 creates a static binary (required for Alpine)
# -ldflags="-s -w" strips debug info to reduce file size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o portfolioTUI main.go

# --- Stage 2: Create the Runtime Image ---
FROM alpine:latest

WORKDIR /root/

# 1. Install Certificates (Required for MongoDB connection)
RUN apk --no-cache add ca-certificates

# 2. Create the .ssh directory (So we can mount the volume later)
RUN mkdir -p .ssh && chmod 700 .ssh

# 3. Copy the binary from Stage 1
COPY --from=builder /app/portfolioTUI .

# 4. Expose the SSH port
EXPOSE 23234

# 5. Run the application
CMD ["./portfolioTUI"]
