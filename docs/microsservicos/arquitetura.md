# Arquitetura da Aplicação

## Visão geral

A aplicação segue a arquitetura definida pelo professor: um **frontend HTTP** (módulo P) que se comunica com dois **microserviços gRPC** (A e B) no backend, todos rodando em containers separados no Kubernetes.

<div align="center">
<font size="3"><p style="text-align: center"><b>Figura 1:</b> Arquitetura de Comunicação do projeto</p></font>

![Imagem 1](https://github.com/PSPD-Grupo/Projeto-de-Pesquisa/blob/main/docs/img/arquitetura.jpg)

<font size="3"><p style="text-align: center"><b>Autor:</b> <a href="https://github.com/gabrielfreitass1">Gabriel Freitas</a>, 2026.</p></font> 
</div>

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

## Histórico de Versão
| Versão | Data       | Descrição                                      | Autor               | Revisor               |
|--------|------------|------------------------------------------------|---------------------|-----------------------|
| 1.0    | 04/06/2026 | Primeira versão do artefato gRPC | [Milena Baruc Rodrigues Morais](https://github.com/MilenaBaruc) | [Milena Baruc Rodrigues Morais](https://github.com/MilenaBaruc) |
| 1.1    | 06/06/2026 | Atualização das imagens | [Gabriel Freitas Balbino](https://github.com/gabrielfreitass1) | [Gabriel Freitas Balbino](https://github.com/gabrielfreitass1) |
