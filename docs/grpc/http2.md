# HTTP/2

O gRPC usa **HTTP/2** como protocolo de transporte, o que é um dos principais motivos para sua alta performance em relação a soluções REST tradicionais baseadas em HTTP/1.1.

---

## HTTP/1.1 vs HTTP/2

| Característica | HTTP/1.1 | HTTP/2 |
|---|---|---|
| Formato | Texto | Binário |
| Conexões por domínio | Múltiplas (workaround) | Uma única (multiplexada) |
| Streaming | ❌ | ✅ Nativo |
| Compressão de headers | ❌ | ✅ HPACK |
| Server push | ❌ | ✅ |
| Head-of-line blocking | ✅ Problema presente | ✅ Resolvido |

---

## Multiplexação

No HTTP/1.1, cada requisição ocupa a conexão inteira até receber resposta. No HTTP/2, **múltiplas requisições e respostas trafegam simultaneamente** pela mesma conexão TCP, usando o conceito de *streams*:

```
HTTP/1.1:
  [Req 1] ──▶ [Resp 1] ──▶ [Req 2] ──▶ [Resp 2]   (sequencial)

HTTP/2:
  [Req 1] ──▶
  [Req 2] ──▶  (simultâneo na mesma conexão)
  [Req 3] ──▶
       ◀── [Resp 2]
       ◀── [Resp 1]
       ◀── [Resp 3]
```

---

## Por que isso importa para o gRPC?

O HTTP/2 é o que torna possível os **4 tipos de streaming do gRPC**:

- **Unary** — 1 request → 1 response (funciona igual ao HTTP/1.1)
- **Server Streaming** — 1 request → múltiplas responses (HTTP/2 stream do servidor)
- **Client Streaming** — múltiplos requests → 1 response (HTTP/2 stream do cliente)
- **Bidirecional** — múltiplos requests ↔ múltiplas responses (full-duplex)

---

## Frames e Streams

O HTTP/2 divide os dados em **frames** (unidades binárias mínimas). Cada stream é identificado por um número inteiro e carrega frames independentemente:

```
Conexão TCP única
├── Stream 1: chamada gRPC para PublicarDesabafo
├── Stream 3: chamada gRPC para ListarDesabafos (streaming)
└── Stream 5: chamada gRPC para AdicionarReacao
```

---

## Compressão de headers com HPACK

Headers repetidos (como `content-type: application/grpc`) são comprimidos e referenciados por índice em chamadas subsequentes, reduzindo o overhead de rede significativamente.

!!! note "TLS"
    Na prática, o gRPC em produção usa sempre **HTTP/2 sobre TLS** (h2). Em ambientes de desenvolvimento/teste é possível usar HTTP/2 sem TLS (h2c — cleartext), que é o que usamos no minikube.