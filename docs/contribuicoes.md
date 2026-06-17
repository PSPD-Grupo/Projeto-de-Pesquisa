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
- **Definição do Contrato gRPC**: Especifiquei o arquivo `reacoes.proto`, definindo as mensagens e os dois RPCs do serviço: `EnviarReacoes` (*client streaming*) e `BuscarReacoes` (unário), estabelecendo o contrato de comunicação entre o Microserviço B e os demais componentes do sistema.
- **Implementação do Servidor gRPC em Java**: Desenvolvi a classe `ReacaoServiceImpl`, estendendo a base gerada automaticamente pelo compilador Protobuf, com a lógica de recebimento e contagem de reações via streaming e de consulta por identificador de desabafo.
- **Persistência em Memória com Segurança de Concorrência**: Implementei a classe `ReacaoRepository` utilizando `ConcurrentHashMap` e a operação atômica `merge`, garantindo integridade dos dados sob acesso simultâneo de múltiplas threads do servidor gRPC.
- **Configuração do Build e Empacotamento**: Configurei o `pom.xml` com o `protobuf-maven-plugin` para geração automática de código a partir do `.proto`, o `maven-shade-plugin` para empacotamento em fat JAR executável e o `maven-surefire-plugin` para execução de testes.
- **Containerização**: Elaborei o `Dockerfile` do serviço utilizando a imagem `eclipse-temurin:21-jre`, tornando o servidor pronto para execução isolada e integração com o cluster Kubernetes do projeto.

### O que aprendi
- **Funcionamento Prático do gRPC e Protobuf**: Compreendi como o arquivo `.proto` atua como contrato entre serviços de linguagens diferentes e como o compilador gera automaticamente as classes Java de mensagens e a classe base do servidor, dispensando implementação manual dessas estruturas.
- **Tipos de Comunicação gRPC**: Aprendi a distinguir os quatro modelos de comunicação e a aplicar dois deles no mesmo serviço, entendendo em quais cenários o *client streaming* é mais adequado do que chamadas unárias repetidas.
- **Concorrência em Servidores gRPC**: Entendi que o framework processa requisições em múltiplas threads simultaneamente e que estruturas como `HashMap` não são seguras nesse contexto, tornando o uso de `ConcurrentHashMap` com operações atômicas uma exigência de correção, não uma preferência.
- **Empacotamento de Aplicações Java**: Compreendi o propósito do fat JAR e o papel do Maven Shade Plugin na geração de um artefato autossuficiente adequado a ambientes sem acesso a repositórios de dependências.

### Dificuldades superadas
- Configurar o `pom.xml` para integração correta entre o plugin Protobuf, o compilador `protoc` e o Maven Shade Plugin, lidando com detalhes como a necessidade do `os-maven-plugin` para detecção de plataforma e o risco de colisão de arquivos SPI do Netty durante o empacotamento.
- Compreender e aplicar corretamente os mecanismos de concorrência, partindo apenas da base teórica vista em aula e chegando a uma implementação que garante segurança de dados sob carga paralela real.


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
