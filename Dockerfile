# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files for better caching
COPY frontend/package*.json ./
RUN npm install --legacy-peer-deps

# Copy source and build
COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS backend-builder

# Build arguments for cross-compilation
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    GOTOOLCHAIN=auto go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY embed.go ./

# Copy frontend build from previous stage
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Build with cache mounts and target platform
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOTOOLCHAIN=auto \
    go build -ldflags="-w -s" -trimpath -o tailtunnel ./cmd/tailtunnel

# Stage 3: Runtime
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates iptables ip6tables

WORKDIR /app

# Copy the binary from builder
COPY --from=backend-builder /app/tailtunnel .

# Create state directory
RUN mkdir -p /var/lib/tailtunnel

# Expose port
EXPOSE 8080

# Environment variables
ENV PORT=8080
ENV STATE_DIR=/var/lib/tailtunnel

# Run as non-root user
RUN adduser -D -u 1000 tailtunnel && \
    chown -R tailtunnel:tailtunnel /app /var/lib/tailtunnel
USER tailtunnel

CMD ["./tailtunnel"]
