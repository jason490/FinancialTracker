FROM golang:1.26-alpine

# Install build dependencies for CGO and SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Pre-copy go.mod and go.sum for better caching
COPY go.mod go.sum ./
RUN go mod download

# Source code will be mounted via volume for development
# CMD will be overridden or specified here to use air
CMD ["air", "-c", ".air.toml"]
