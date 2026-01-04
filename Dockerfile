# Stage 1: Build Frontend
FROM node:20-alpine AS web-builder
WORKDIR /web

# Copy package files first for caching
COPY web/package.json web/package-lock.json ./
RUN npm ci

# Copy the rest of the frontend source
COPY web/ .
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.25-alpine AS go-builder
WORKDIR /app

# Install git if needed for go mod (usually not needed for public modules, but good practice)
# RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source
COPY cmd/ cmd/
COPY internal/ internal/
COPY ui.go .

# Copy built frontend assets to the location expected by go:embed
# ui.go expects "web/dist/*" relative to itself
COPY --from=web-builder /web/dist ./web/dist

# Build the binary
# CGO_ENABLED=0 for static binary (we use pure Go sqlite driver)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o synapse cmd/synapse/main.go

# Stage 3: Final Image
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS requests and tzdata for timezones
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=go-builder /app/synapse .

# Create a directory for the database
RUN mkdir -p /data

# Set default environment variables
ENV SYNAPSE_HTTP_PORT=:8080 \
    SYNAPSE_MQTT_PORT=:1883 \
    SYNAPSE_WS_PORT=:8083 \
    SYNAPSE_DB_PATH=/data/synapse.db

# Expose ports
EXPOSE 8080 1883 8083

# Volume for persistence
VOLUME ["/data"]

# Run the binary
ENTRYPOINT ["/app/synapse"]
