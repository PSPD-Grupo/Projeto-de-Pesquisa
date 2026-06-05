# Exemplos e Testes

> 🚧 **Esta seção será preenchida com os testes realizados pelo grupo.**  
> Os exemplos de código abaixo são o ponto de partida — cada membro deve rodar localmente e adicionar capturas de tela e resultados reais.

---

## Arquivo `.proto` base do projeto

```protobuf title="diario.proto"
syntax = "proto3";

package diario;

// ── Mensagens ───────────────────────────────────────────────

message Desabafo {
  string id        = 1;
  string texto     = 2;
  string autor     = 3;
  int64  timestamp = 4;
}

message PublicarRequest  { string texto  = 1; }
message PublicarResponse { string id     = 1; bool sucesso = 2; }

message FeedRequest      { int32  limite = 1; }

message ReacaoRequest    { string desabafo_id = 1; string tipo = 2; }
message ResumoReacoes    { int32  total = 1; }

message MensagemAoVivo   { string texto = 1; string sala_id = 2; }

// ── Serviços ────────────────────────────────────────────────

// Servidor A — gerencia desabafos
service DesabafoService {
  // 1. Unary
  rpc Publicar (PublicarRequest) returns (PublicarResponse);

  // 2. Server Streaming
  rpc ListarFeed (FeedRequest) returns (stream Desabafo);

  // 4. Bidirecional
  rpc SalaAoVivo (stream MensagemAoVivo) returns (stream MensagemAoVivo);
}

// Servidor B — gerencia reações
service ReacaoService {
  // 3. Client Streaming
  rpc EnviarReacoes (stream ReacaoRequest) returns (ResumoReacoes);
}
```

---

## Teste 1 – Unary Call

**Descrição:** Publicar um desabafo simples e receber confirmação.

=== "Servidor (Python)"
    ```python
    import grpc
    from concurrent import futures
    import diario_pb2
    import diario_pb2_grpc
    import uuid, time

    class DesabafoServicer(diario_pb2_grpc.DesabafoServiceServicer):
        def Publicar(self, request, context):
            novo_id = str(uuid.uuid4())[:8]
            print(f"[SERVIDOR] Novo desabafo: '{request.texto[:30]}...'")
            return diario_pb2.PublicarResponse(id=novo_id, sucesso=True)

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    diario_pb2_grpc.add_DesabafoServiceServicer_to_server(DesabafoServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()
    ```

=== "Cliente (Python)"
    ```python
    import grpc
    import diario_pb2
    import diario_pb2_grpc

    with grpc.insecure_channel('localhost:50051') as channel:
        stub = diario_pb2_grpc.DesabafoServiceStub(channel)
        resp = stub.Publicar(diario_pb2.PublicarRequest(
            texto="Hoje foi um dia muito difícil, mas vou conseguir."
        ))
        print(f"ID: {resp.id} | Sucesso: {resp.sucesso}")
    ```

!!! note "Resultado esperado"
    ```
    ID: a3f1b2c4 | Sucesso: True
    ```

---

## Teste 2 – Server Streaming

**Descrição:** Solicitar o feed e receber os desabafos um a um.

=== "Servidor (Python)"
    ```python
    def ListarFeed(self, request, context):
        desabafos = [
            "Precisava desabafar isso faz tempo...",
            "Ninguém entende como me sinto.",
            "Hoje finalmente sorri de verdade.",
        ]
        for i, texto in enumerate(desabafos[:request.limite]):
            yield diario_pb2.Desabafo(
                id=str(i), texto=texto, autor="anônimo"
            )
    ```

=== "Cliente (Python)"
    ```python
    for desabafo in stub.ListarFeed(diario_pb2.FeedRequest(limite=10)):
        print(f"[{desabafo.id}] {desabafo.texto}")
    ```

---

## Teste 3 – Client Streaming

**Descrição:** Enviar várias reações em lote e receber o total.

=== "Cliente (Python)"
    ```python
    def gerar_reacoes():
        for tipo in ["❤️", "❤️", "❤️"]:
            yield diario_pb2.ReacaoRequest(desabafo_id="abc123", tipo=tipo)

    resumo = stub_reacao.EnviarReacoes(gerar_reacoes())
    print(f"Total de reações processadas: {resumo.total}")
    ```

---

## Teste 4 – Bidirecional Streaming

**Descrição:** Sala ao vivo onde cliente envia e recebe mensagens simultaneamente.

=== "Cliente (Python)"
    ```python
    import threading

    def enviar(stub):
        msgs = ["oi sala", "alguém aí?", "que dia difícil..."]
        for m in msgs:
            yield diario_pb2.MensagemAoVivo(texto=m, sala_id="sala-1")

    for resposta in stub.SalaAoVivo(enviar(stub)):
        print(f"[AO VIVO] {resposta.texto}")
    ```

---

!!! warning "TODO para o grupo"
    - [ ] Rodar cada teste localmente e adicionar print das saídas reais
    - [ ] Adicionar capturas de tela dos terminais
    - [ ] Documentar erros encontrados e como foram resolvidos