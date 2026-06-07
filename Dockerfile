# Steg 1 - bygg
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o task-mgr .

# Steg 2 - kjør
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/task-mgr .
CMD ["./task-mgr"]