# Servidor B - Reações

Microserviço responsável por armazenar e consultar
reações dos desabafos.

## Tecnologias

- Java 21
- Maven
- gRPC
- Protobuf
- Docker

## Porta

50052

## Serviços gRPC

### EnviarReacoes

Client Streaming

Recebe várias reações e retorna um resumo.

### BuscarReacoes

Unary

Retorna quantidade de reações de um desabafo.

## Executando localmente

mvn clean package

java -jar target/servidor-b.jar

## Docker

docker build -t servidor-b .

docker run -p 50052:50052 servidor-b