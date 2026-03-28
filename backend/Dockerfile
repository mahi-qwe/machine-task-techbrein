FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/api

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server ./server
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
