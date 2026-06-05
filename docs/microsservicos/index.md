# Nossa Aplicação: *Desabafo* 🗒️

## A ideia

*Desabafo* é um **diário anônimo distribuído** — um espaço simples onde qualquer pessoa pode publicar um texto curto para desabafar, sem precisar de conta ou identificação, e ver o que outros estão sentindo.

Pense num Twitter, mas:
- Sem curtidas públicas (só reações anônimas)
- Sem seguidores
- Sem histórico ligado a você
- Só o desabafo importa

---

## Por que essa aplicação?

| Critério | Por que atende |
|---|---|
| Fácil de explicar | Qualquer pessoa entende o conceito em 10 segundos |
| Demonstra todos os 4 tipos gRPC | Publicar (unary), feed (server-stream), reações em lote (client-stream), sala ao vivo (bidirecional) |
| Dados distribuídos naturalmente | Posts ficam no Servidor A, reações no Servidor B |
| Mostra valor do gRPC | Streaming de feed em tempo real é algo que REST não faz bem |

---

## Funcionalidades

### Módulo P (API Gateway + Web Server)
- Interface web simples com caixa de texto para desabafo
- Feed de desabafos recentes
- Traduz requisições HTTP → gRPC para os servidores A e B

### Servidor A – Desabafos
- Receber e armazenar novos desabafos
- Retornar o feed (stream de desabafos)
- Sala ao vivo (streaming bidirecional)

### Servidor B – Reações
- Receber reações em lote (❤️)
- Retornar contagem de reações por desabafo

---

!!! info "Continue lendo"
    - [Arquitetura técnica detalhada →](arquitetura.md)
    - [Funcionalidades planejadas e divisão entre módulos →](diario-desabafo.md)