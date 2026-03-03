## Build stage — compiles the binary
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Download deps first (cached unless go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 go build -o /app/server

## Runtime stage — minimal image with just the binary
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/server .

CMD ["./server"]
