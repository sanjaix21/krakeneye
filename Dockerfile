FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY . .

RUN go build -o krakeneye .

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/krakeneye .

CMD ["./krakeneye", "--web"]
