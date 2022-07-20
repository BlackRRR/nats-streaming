FROM golang:latest

WORKDIR /app

COPY . .
RUN go build ./cmd/nats-streaming

CMD ["./nats-streaming"]
