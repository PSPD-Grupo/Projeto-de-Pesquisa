# Contribuições e Aprendizados

Esta seção reúne o relato individual de cada membro do grupo sobre sua participação no projeto de pesquisa, o que aprendeu e as principais dificuldades superadas.

---

## 👩‍💻 Milena Baruc Rodrigues Morais

### O que fiz
- **Estudo e Escrita de Comunicação gRPC**: Responsável por documentar a teoria dos 4 modelos de comunicação gRPC (Unary, Server/Client/Bidirectional Streaming) e estruturar a documentação inicial do MkDocs.
- **Resolução de Concorrência em Go**: Implementei a segurança de concorrência e o gerenciamento de salas com `sync.Mutex` em `chat.go` (Servidor A), eliminando pânicos de referências de ponteiros nulos e corrida de dados (race conditions) entre goroutines dos usuários.
- **Ponte WebSocket/gRPC**: Construí a lógica no API Gateway (FastAPI) para criar pontes assíncronas em tempo real, permitindo a comunicação entre navegadores WebSockets e as streams gRPC bidirecionais do servidor Go.
- **Persistência SQLite**: Solucionei erros silenciosos de `SQLITE_MISUSE` causados pelo driver antigo de banco e incompatibilidades de tipos de timestamp, migrando para o driver oficial `go-sqlite3` e ajustando os tipos para `INTEGER`.
- **Testes Automatizados (Servidor A)**: Implementei e validei a suíte de testes unitários automatizados para a persistência no Servidor A (`go test`).
- **Ajustes de Portabilidade**: Modifiquei o frontend e os caminhos de persistência para que a aplicação rode localmente e no Minikube de forma automática, sem IPs fixos.

### O que aprendi
- **Conceitos de gRPC e HTTP/2**: Entendi na prática as vantagens de comunicação binária sobre canais multiplexados em comparação com APIs REST/JSON.
- **Concorrência Segura em Go**: Aprendi a gerenciar acessos a mapas concorrentes por meio de Mutexes e a evitar vazamento de goroutines em streams de streaming bidirecional.
- **Programação Assíncrona e Buffers**: Entendi como gerenciar leituras/escritas concorrentes em conexões assíncronas do Python (FastAPI/WebSockets) mapeadas para gRPC assíncrono (`grpc.aio`).
- **Gestão de Drivers de Banco de Dados**: A importância da escolha correta de drivers SQL nativos/CGO e a correlação de tipos de dados entre banco e protocolos de transferência.

### Dificuldades superadas
- Identificar e depurar o erro genérico do SQLite `bad parameter or other API misuse: not an error (21)`.
- Gerenciar o ciclo de vida de streams gRPC bidirecionais ativas que podiam ser canceladas abruptamente pelo cliente web.

---

## 👨‍💻 Pedro Fonseca Cruz

### O que fiz
*(A preencher)*

### O que aprendi
*(A preencher)*

### Dificuldades superadas
*(A preencher)*

---

## 👨‍💻 Daniel dos Santos Barros de Sousa

### O que fiz
*(A preencher)*

### O que aprendi
*(A preencher)*

### Dificuldades superadas
*(A preencher)*

---

## 👨‍💻 Gabriel Freitas Balbino

### O que fiz
*(A preencher)*

### O que aprendi
*(A preencher)*

### Dificuldades superadas
*(A preencher)*
