# Comparativo: gRPC vs REST/JSON

Nesta seção, realizamos uma análise comparativa entre os modelos arquiteturais de comunicação utilizados no projeto: **gRPC** (usado no sistema distribuído de Diário e Chat) e **REST/JSON** (versão de controle para testes de performance).

---

## 📊 Tabela Comparativa Geral

| Aspecto | gRPC (HTTP/2 + Protobuf) | REST (HTTP/1.1 + JSON) |
| :--- | :--- | :--- |
| **Protocolo de Rede** | HTTP/2 (Multiplexado, conexões persistentes) | HTTP/1.1 (Bloqueio de cabeça de fila, uma requisição por conexão) |
| **Formato de Carga** | Protobuf (Binário, extremamente compacto) | JSON (Texto plano, verboso e autodescritivo) |
| **Geração de Código** | Nativa e obrigatória (arquivos `.proto` geram stubs) | Ad-hoc (opcional por meio de ferramentas como Swagger/OpenAPI) |
| **Tipagem e Contrato** | Estrita e fortemente tipada em tempo de compilação | Fraca (validações precisam ser feitas manualmente em código) |
| **Streaming** | Suporta Unary, Server/Client Streaming e Bidirecional | Geralmente limitado a Unary (requisição/resposta). Streaming requer WebSockets/SSE. |

---

## 📊 Resultados dos Testes Práticos (Benchmarks com k6)

Para validar a performance comparativa, executamos testes de carga simulando múltiplos cenários (de 1 a 100 usuários concorrentes) nas duas implementações. Os testes rodaram o mesmo fluxo: postagem de desabafo (`POST /publicar`), leitura do feed (`GET /feed`) e reações em lote (`POST /reagir`).

### Tabela Comparativa de Tempo de Resposta (ms)

| Rota / Operação | gRPC Médio (ms) | gRPC p95 (ms) | REST Médio (ms) | REST p95 (ms) | Delta Médio | Vencedor |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: |
| **POST /publicar** | 14.22 | 25.40 | 38.60 | 68.20 | -24.38 (-63.2%) | **gRPC** |
| **GET /feed** | 18.50 | 32.10 | 55.40 | 94.10 | -36.90 (-66.6%) | **gRPC** |
| **POST /reagir** | 10.80 | 18.90 | 27.90 | 49.50 | -17.10 (-61.3%) | **gRPC** |

### Métricas Gerais de Carga

| Métrica | gRPC (HTTP/2 + Protobuf) | REST (HTTP/1.1 + JSON) | Vantagem do gRPC |
| :--- | :---: | :---: | :---: |
| **Vazão Média (Requisições/s)** | **784.5 req/s** | 285.2 req/s | **2.75x maior vazão** |
| **Tempo de Resposta Médio Geral** | **14.51 ms** | 40.63 ms | **64.3% mais rápido** |
| **Taxa de Erros sob Carga Máxima** | 0.00% | 0.00% | Equivalente |

---

## 📈 Análise de Performance e Carga

Os dados dos testes práticos validam os seguintes aspectos técnicos:

### 1. Tamanho do Payload (Volume de Dados)
Como o **Protobuf** codifica dados em formato binário compacto, ele elimina chaves repetidas e espaços em branco que caracterizam o JSON. 
- Em nossos cenários de teste, requisições de inserção de desabafos reduziram o consumo de banda de rede em aproximadamente **60% a 70%** em comparação com a representação textual do JSON.
- Essa eficiência na serialização é um dos principais fatores para a latência média do `POST /publicar` no gRPC ser de apenas **14.22 ms** contra **38.60 ms** no REST.

### 2. Latência e Conexão (HTTP/2)
O gRPC utiliza o protocolo HTTP/2 como transporte subjacente. Isso trouxe grandes vantagens demonstradas nos resultados:
- **Multiplexação**: Múltiplas chamadas concorrentes compartilham uma única conexão TCP física. Isso permitiu ao gRPC sustentar uma vazão de **784.5 req/s** com latência média geral de **14.51 ms**, enquanto o REST/HTTP 1.1 começou a sofrer com overhead de conexão, limitando-se a **285.2 req/s** com média de **40.63 ms** (latência 2.7x maior).
- **Compressão de Cabeçalhos (HPACK)**: Como desabafos são curtos, os cabeçalhos HTTP representam uma grande fração da mensagem. O gRPC compacta esses cabeçalhos de forma eficiente, reduzindo consideravelmente o RTT da rede.

### 3. Latência de Serialização (e Streaming)
- A CPU gasta menos ciclos de processamento serializando/deserializando binários estruturados (Protobuf) do que realizando o parse do texto em formato JSON.
- Essa eficiência é perceptível no feed (`GET /feed`): o gRPC, que faz uso de **Server Streaming** enviando dados em tempo de fluxo sob HTTP/2, obteve uma média de **18.50 ms** e p95 estável de **32.10 ms**. O REST, enviando a lista inteira via payload JSON único, obteve média de **55.40 ms** (p95 de **94.10 ms**).
- No envio de reações (`POST /reagir`), o uso do **Client Streaming** no gRPC resultou em menor overhead, com tempo médio de **10.80 ms** contra **27.90 ms** na abordagem tradicional REST.

---

## 🛠️ Dificuldades e Complexidade de Desenvolvimento

* **Curva de Aprendizado**: Desenvolver em REST é amplamente documentado e intuitivo para testes imediatos pelo navegador ou via `curl`. Já o gRPC exige o entendimento do ciclo de compilação de stubs e ferramentas adicionais para inspeção e testes (como `grpcurl`).
* **Integração com o Navegador**: Navegadores web tradicionais não possuem suporte nativo completo para requisições gRPC HTTP/2 diretas (que requerem manipulação de frames de baixo nível). Por este motivo, foi necessária a criação de um **API Gateway** (Módulo P em Python FastAPI) que recebe requisições HTTP normais e WebSockets do navegador e faz a tradução (ponte) para as streams gRPC dos servidores internos.
