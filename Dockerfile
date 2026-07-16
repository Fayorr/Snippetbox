# Stage 1: Build base & download dependencies
FROM golang:1.26-alpine AS base
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Stage 2: Build the static binary
FROM base AS build-server 
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/snippetbox ./cmd/web

# Stage 3: Minimal production run image (using scratch)
FROM scratch
WORKDIR /app

# Copy SSL certificates for external HTTPS requests
COPY --from=build-server /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled Go binary
COPY --from=build-server /app/snippetbox .

# Copy your UI templates/assets (HTML, CSS, JS)
COPY --from=build-server /app/ui ./ui

# ---> CRITICAL: Copy your local TLS certificates into the image <---
COPY --from=build-server /app/tls ./tls

# Expose both HTTP and HTTPS ports if needed (typically 4000 for Snippetbox HTTPS)
EXPOSE 4000

ENTRYPOINT [ "/app/snippetbox" ]