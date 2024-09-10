# Build step
FROM golang:1.22.3-alpine AS builder

WORKDIR /bin

COPY . .

RUN go mod download

RUN go build -o server .

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /bin/server .

EXPOSE 2080

EXPOSE 4043

CMD ["./server"]