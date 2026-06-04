# Tipos de Comunicação gRPC

O gRPC suporta **4 tipos de comunicação**, todos definidos no arquivo `.proto` e todos transportados via HTTP/2.

---

## 1. Unary Call (chamada simples)

O cliente envia **uma requisição** e recebe **uma resposta**. É o mais simples e equivale ao modelo tradicional REST.

```protobuf
// No .proto:
rpc Publicar (PublicarRequest) returns (PublicarResponse);
```

**No contexto do Diário:**  
O usuário escreve um desabafo e clica em "Publicar" → o módulo P envia para o Servidor A → recebe confirmação.

```
Cliente ──[PublicarRequest]──▶ Servidor
        ◀──[PublicarResponse]──
```

**Quando usar:**  
Operações simples de requisição/resposta: criar um recurso, autenticar, buscar um item por ID.

---

## 2. Server Streaming

O cliente envia **uma requisição** e o servidor responde com **um fluxo de mensagens**.

```protobuf
// No .proto:
rpc ListarFeed (FeedRequest) returns (stream Desabafo);
```

**No contexto do Diário:**  
O usuário abre o feed → módulo P pede ao Servidor A → Servidor A vai enviando os desabafos um a um conforme os recupera do banco.

```
Cliente ──[FeedRequest]──▶ Servidor
        ◀──[Desabafo 1]───
        ◀──[Desabafo 2]───
        ◀──[Desabafo 3]───
        ◀──[FIM]──────────
```

**Quando usar:**  
Feed de notícias, resultados de busca paginados em tempo real, exportação de dados grandes, notificações push.

---

## 3. Client Streaming

O cliente envia **um fluxo de mensagens** e o servidor responde **uma única vez** ao final.

```protobuf
// No .proto:
rpc EnviarReacoes (stream ReacaoRequest) returns (ResumoReacoes);
```

**No contexto do Diário:**  
Um cliente acumulou várias reações offline e as envia em lote → o Servidor B processa tudo e responde com um resumo.

```
Cliente ──[Reacao 1]──▶ Servidor
        ──[Reacao 2]──▶
        ──[Reacao 3]──▶
        ──[FIM]───────▶
                        ◀──[ResumoReacoes]──
```

**Quando usar:**  
Upload de arquivos em chunks, envio de logs em lote, telemetria/métricas acumuladas.

---

## 4. Bidirecional Streaming

**Ambos os lados** enviam e recebem um fluxo de mensagens simultaneamente. Full-duplex.

```protobuf
// No .proto:
rpc SalaAoVivo (stream MensagemAoVivo) returns (stream MensagemAoVivo);
```

**No contexto do Diário:**  
Uma sala de desabafos ao vivo onde os usuários veem em tempo real enquanto outros publicam — como um chat.

```
Cliente ──[Msg 1]──▶        ◀──[Msg de outro]── Servidor
        ──[Msg 2]──▶
                    ◀──[Msg de outro]──
        ──[Msg 3]──▶
```

**Quando usar:**  
Chat em tempo real, jogos multiplayer, colaboração ao vivo (tipo Google Docs), monitoramento contínuo.

---

## Resumo comparativo

| Tipo | Request | Response | Caso de uso típico |
|---|---|---|---|
| Unary | 1 | 1 | CRUD, autenticação |
| Server Streaming | 1 | N | Feed, notificações |
| Client Streaming | N | 1 | Upload, lotes |
| Bidirecional | N | N | Chat, jogos, colaboração |