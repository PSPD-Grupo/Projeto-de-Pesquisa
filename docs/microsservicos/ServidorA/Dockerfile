FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/bufbuild/buf/cmd/buf@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
COPY . .
RUN buf generate
RUN go build -o servidor .
RUN ls /app/db   

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/servidor .
COPY --from=builder /app/db ./db
EXPOSE 50051
CMD ["./servidor"]