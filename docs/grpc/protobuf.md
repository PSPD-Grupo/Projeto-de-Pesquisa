# Protocol Buffers (Protobuf)

**Protocol Buffers** (protobuf) é o mecanismo de serialização de dados usado pelo gRPC. Desenvolvido pelo Google, ele define a estrutura das mensagens trocadas entre cliente e servidor de forma **tipada, compacta e independente de linguagem**.

---

## Como funciona?

Você define suas mensagens e serviços em um arquivo `.proto`:

```protobuf title="diario.proto"
syntax = "proto3";

package diario;

// Mensagem que representa um desabafo
message Desabafo {
  string id        = 1;
  string texto     = 2;
  string autor     = 3;  // anônimo por padrão
  int64  timestamp = 4;
}

// Requisição para publicar um desabafo
message PublicarRequest {
  string texto = 1;
}

// Resposta após publicar
message PublicarResponse {
  string id      = 1;
  bool   sucesso = 2;
}

// Serviço principal
service DiarioService {
  rpc Publicar (PublicarRequest) returns (PublicarResponse);
}
```

O protobuf gera automaticamente **código cliente e servidor** nas linguagens suportadas (Python, Go, Java, C++, etc.) a partir desse arquivo.

---

## Serialização binária vs JSON

O protobuf codifica os dados em formato **binário**, não texto. Isso traz vantagens:

=== "JSON (texto)"
    ```json
    {
      "id": "abc123",
      "texto": "Hoje foi um dia muito difícil...",
      "autor": "anônimo",
      "timestamp": 1748995200
    }
    ```
    **Tamanho:** ~85 bytes

=== "Protobuf (binário)"
    ```
    0a 06 61 62 63 31 32 33 12 1f 48 6f 6a 65 20 ...
    ```
    **Tamanho:** ~45 bytes (~47% menor)

---

## Numeração dos campos

Cada campo no `.proto` recebe um **número único** (tag). É esse número — não o nome — que é serializado no binário. Por isso:

!!! warning "Cuidado ao evoluir o schema"
    Nunca reutilize o número de um campo removido. Isso garante compatibilidade entre versões diferentes do serviço.

---

## Tipos de dados disponíveis

| Tipo proto | Equivalente |
|---|---|
| `string` | UTF-8 string |
| `int32`, `int64` | Inteiros |
| `float`, `double` | Ponto flutuante |
| `bool` | Booleano |
| `bytes` | Dados binários |
| `repeated` | Lista/array |
| `message` | Objeto aninhado |
| `enum` | Enumeração |

---

## Geração de código

```bash
# Gerar código Python a partir do .proto
python -m grpc_tools.protoc \
  -I. \
  --python_out=. \
  --grpc_python_out=. \
  diario.proto

# Gerar código Go
protoc --go_out=. --go-grpc_out=. diario.proto
```

Isso cria automaticamente os arquivos com as classes das mensagens e os stubs de cliente/servidor.