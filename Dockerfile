# Build step
FROM golang:1.22-3-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o server .

FROM alpine:latest

# Run step
WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 2020

CMD ["./server"]