# Stage 1: build app Go
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Stage 2: run app with lightweight distro
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary từ builder
COPY --from=builder /app/main .

# Copy env + firebase secret
COPY .env .env

EXPOSE 8080
CMD ["./main"]
