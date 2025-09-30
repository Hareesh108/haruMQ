# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o broker ./cmd/broker
RUN go build -o producer ./cmd/producer
RUN go build -o consumer ./cmd/consumer

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/broker ./broker
COPY --from=builder /app/producer ./producer
COPY --from=builder /app/consumer ./consumer
COPY config.yaml ./config.yaml
COPY data ./data
EXPOSE 9092
CMD ["./broker"]
