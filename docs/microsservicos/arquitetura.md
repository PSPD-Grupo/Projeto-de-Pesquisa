# Arquitetura da Aplicação

## Visão geral

A aplicação segue a arquitetura definida pelo professor: um **frontend HTTP** (módulo P) que se comunica com dois **microserviços gRPC** (A e B) no backend, todos rodando em containers separados no Kubernetes.

```
┌────────────┐   HTTP/REST   ┌──────────────────────────────────────────┐
│            │               │                K8S (minikube)            │
│ WEB CLIENT │◀─────────────▶│  ┌──────────┐                           │
│ (browser)  │               │  │    (P)   │──gRPC──▶ Servidor A (Go)  │
└────────────┘               │  │ Python   │                           │
                             │  │ FastAPI  │──gRPC──▶ Servidor B (Java)│
                             │  └──────────┘                           │
                             └──────────────────────────────────────────┘
```

---

## Tecnologias por módulo

| Módulo | Linguagem | Framework | Responsabilidade |
|---|---|---|---|
| **P** (API Gateway) | Python | FastAPI + gRPC stub | Recebe HTTP, chama A e B via gRPC |
| **A** (Desabafos) | Go | grpc-go | Armazena e serve desabafos |
| **B** (Reações) | Java | grpc-java | Armazena e serve reações |

!!! note "Requisito do projeto"
    O enunciado exige que A e B sejam implementados em linguagem **diferente** de P. Por isso escolhemos Python para P, Go para A e Java para B.

---

## Mapeamento gRPC × Tipos de Comunicação

| Operação | Módulos | Tipo gRPC |
|---|---|---|
| Publicar desabafo | P → A | **Unary** |
| Carregar feed | P → A | **Server Streaming** |
| Enviar reações em lote | P → B | **Client Streaming** |
| Sala ao vivo | P ↔ A | **Bidirecional Streaming** |

---

## Estrutura de containers no Kubernetes

```
minikube (host único)
├── Pod: container-P
│   ├── FastAPI (porta 8000, exposta externamente via NodePort)
│   └── gRPC Stub (conecta à rede interna do cluster)
│
├── Pod: container-A
│   └── gRPC Server Go (porta 50051, somente rede interna)
│
└── Pod: container-B
    └── gRPC Server Java (porta 50052, somente rede interna)
```

---

## Fluxo de uma requisição

1. Usuário acessa `http://<ip-minikube>:8000` no browser
2. FastAPI (P) recebe a requisição HTTP
3. P chama o Servidor A via gRPC (protobuf + HTTP/2)
4. Opcionalmente P também chama o Servidor B
5. P agrega as respostas e retorna JSON para o browser

---

## Diagrama de sequência – Publicar desabafo

```
Browser          P (FastAPI)         A (Go gRPC)
   │                  │                   │
   │── POST /publicar ▶│                   │
   │                  │── Publicar() ─────▶│
   │                  │                   │ (armazena)
   │                  │◀── {id, sucesso} ──│
   │◀── 200 OK {id} ──│                   │
```