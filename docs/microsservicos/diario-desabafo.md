# Funcionalidades Planejadas

> Esta seção detalha **o que cada módulo vai fazer**, servindo como referência para apresentação ao professor e para orientar a implementação dos outros membros do grupo.

---

## Módulo P — API Gateway (Python + FastAPI)

### Endpoints HTTP expostos ao browser

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/` | Página principal com o feed |
| `POST` | `/publicar` | Publica um novo desabafo |
| `GET` | `/feed` | Retorna lista de desabafos (streaming interno) |
| `POST` | `/reagir` | Envia reações para o Servidor B |
| `WS` | `/sala/{id}` | WebSocket para sala ao vivo (bidirecional) |

### Como P usa o gRPC

P **não armazena dados** — apenas traduz requisições HTTP para chamadas gRPC:

```python
# Exemplo: POST /publicar → Unary gRPC para o Servidor A
@app.post("/publicar")
async def publicar(texto: str):
    with grpc.insecure_channel("servidor-a:50051") as ch:
        stub = DesabafoServiceStub(ch)
        resp = stub.Publicar(PublicarRequest(texto=texto))
        return {"id": resp.id, "sucesso": resp.sucesso}
```

---

## Servidor A — Desabafos (Go)

### Funcionalidades

#### 1. `Publicar` — Unary
Recebe um texto, gera um ID único, armazena em memória (ou SQLite) e retorna confirmação.

```go
func (s *server) Publicar(ctx context.Context, req *pb.PublicarRequest) (*pb.PublicarResponse, error) {
    id := uuid.New().String()[:8]
    s.armazenar(id, req.Texto)
    return &pb.PublicarResponse{Id: id, Sucesso: true}, nil
}
```

#### 2. `ListarFeed` — Server Streaming
Retorna os N desabafos mais recentes, enviando um a um via stream (simula latência de banco de dados).

```go
func (s *server) ListarFeed(req *pb.FeedRequest, stream pb.DesabafoService_ListarFeedServer) error {
    for _, d := range s.ultimos(int(req.Limite)) {
        time.Sleep(50 * time.Millisecond) // simula latência
        stream.Send(d)
    }
    return nil
}
```

#### 3. `SalaAoVivo` — Bidirecional Streaming
Mantém um canal bidirecional: quando um cliente envia uma mensagem, o servidor a retransmite para todos os outros clientes conectados na mesma sala.

---

## Servidor B — Reações (Java)

### Funcionalidades

#### `EnviarReacoes` — Client Streaming
Recebe um fluxo de ❤️ relacionados a desabafos específicos, acumula em memória e ao final retorna um resumo com o total processado.

```java
@Override
public StreamObserver<ReacaoRequest> enviarReacoes(
    StreamObserver<ResumoReacoes> responseObserver) {
  return new StreamObserver<>() {
    int total = 0;

    @Override
    public void onNext(ReacaoRequest req) {
      armazenar(req.getDesabafoId()); // sempre ❤️
      total++;
    }

    @Override
    public void onCompleted() {
      responseObserver.onNext(ResumoReacoes.newBuilder().setTotal(total).build());
      responseObserver.onCompleted();
    }
  };
}
```

#### `BuscarReacoes` — Unary (extra)
Retorna a contagem de cada tipo de reação para um desabafo específico (chamado por P quando exibe o feed).

---

## Dados armazenados

Para simplicidade, os servidores usam **armazenamento em memória** (maps/dicionários). Isso significa que os dados são perdidos ao reiniciar o container — o que é aceitável para fins de demonstração.

| Servidor | Dados |
|---|---|
| A | `map[string]Desabafo` — id → desabafo |
| B | `map[string]int` — desabafo_id → contagem de ❤️ |

---

## Cenários de demonstração para o professor

!!! example "Demo 1 – Fluxo básico"
    1. Abrir o browser em `http://<minikube-ip>:8000`
    2. Digitar um desabafo e clicar Publicar
    3. Ver o desabafo aparecer no feed
    4. Mostrar nos logs do terminal que passou por P → A via gRPC

!!! example "Demo 2 – Server Streaming"
    1. Publicar 5 desabafos
    2. Abrir o feed e mostrar no log que os itens chegam em stream (um a um)
    3. Comparar com versão REST/JSON onde tudo chega de uma vez

!!! example "Demo 3 – Bidirecional (Sala ao vivo)"
    1. Abrir dois browsers na rota `/sala/sala-1`
    2. Digitar em um e mostrar aparecendo no outro em tempo real
    3. Demonstrar que é full-duplex

!!! example "Demo 4 – Kubernetes"
    1. Mostrar `kubectl get pods` com os 3 containers rodando
    2. Mostrar `kubectl get services` com o serviço P exposto
    3. Derrubar o pod A e mostrar o Kubernetes reiniciando automaticamente