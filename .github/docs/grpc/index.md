# gRPC – Visão Geral

**gRPC** (Google Remote Procedure Call) é um framework open-source de comunicação remota de alta performance, criado pelo Google e mantido pela CNCF.

---

## O que é RPC?

RPC (*Remote Procedure Call*) é um modelo de comunicação onde um programa chama uma função que está rodando em **outro processo ou máquina**, como se fosse uma chamada local. O gRPC moderniza esse conceito usando:

- **Protobuf** como formato de serialização (substitui JSON/XML)
- **HTTP/2** como protocolo de transporte (substitui HTTP/1.1)

---

## Por que usar gRPC?

| Característica | REST/JSON | gRPC/Protobuf |
|---|---|---|
| Serialização | JSON (texto) | Binário (protobuf) |
| Protocolo | HTTP/1.1 | HTTP/2 |
| Streaming | ❌ Limitado | ✅ Nativo (4 tipos) |
| Tipagem | Fraca (JSON) | Forte (.proto) |
| Performance | Moderada | Alta |
| Contrato de API | OpenAPI/Swagger | `.proto` file |

---

## Componentes principais

```
┌─────────────┐         .proto         ┌─────────────┐
│   Cliente   │ ──── define contrato ──▶│   Servidor  │
│  (Stub)     │                         │  (Service)  │
└──────┬──────┘                         └──────┬──────┘
       │         HTTP/2 + Protobuf             │
       └─────────────────────────────────────-─┘
```

- **`.proto` file** — define os serviços e mensagens (contrato da API)
- **Stub (cliente)** — código gerado automaticamente que o cliente usa para chamar o servidor
- **Service (servidor)** — implementação dos métodos definidos no `.proto`
- **Protobuf** — mecanismo de serialização binária
- **HTTP/2** — protocolo de transporte com suporte a multiplexação e streaming

---

!!! info "Continue lendo"
    - [Protobuf em detalhes →](protobuf.md)
    - [HTTP/2 em detalhes →](http2.md)
    - [Os 4 tipos de comunicação →](tipos-comunicacao.md)
    - [Exemplos e testes realizados →](exemplos.md)