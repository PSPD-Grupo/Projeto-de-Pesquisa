FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o servidor .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/servidor .
EXPOSE 50051
CMD ["./servidor"]