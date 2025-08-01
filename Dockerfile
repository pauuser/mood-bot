
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/MoodBot

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates sqlite3 musl libc6 && rm -rf /var/lib/apt/lists/*

RUN groupadd -r appuser && useradd -r -g appuser appuser

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config ./config

RUN mkdir -p logs /app/data && chown -R appuser:appuser logs /app/data

RUN chown appuser:appuser main

USER appuser

EXPOSE 8080

CMD ["./main"] 