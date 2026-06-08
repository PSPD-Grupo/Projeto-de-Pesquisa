FROM golang:1.25 AS builder
WORKDIR /app

# Dependências do Go
COPY go.mod go.sum ./
RUN go mod download

# Instala o buf e os plugins de geração
RUN go install github.com/bufbuild/buf/cmd/buf@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copia o projeto
COPY . .

# Gera o código a partir dos .proto
RUN buf generate

# Compila o binário
RUN go build -o servidor .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/servidor .
EXPOSE 50051
CMD ["./servidor"]