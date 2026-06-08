# Especificação do Projeto

Resumo do que foi pedido pelo Prof. Fernando W. Cruz para o Projeto de Pesquisa – Parte 1.

---

## Objetivo

Construir uma aplicação distribuída baseada em **microserviços gRPC** e fazer o **deploy no Kubernetes (minikube)**.

---

## Arquitetura exigida

<div align="center">
<font size="3"><p style="text-align: center"><b>Figura 1:</b> Arquitetura de Comunicação do projeto</p></font>

![Imagem 1](docs/img/arquitetura.jpg)

<font size="3"><p style="text-align: center"><b>Autor:</b> <a href="https://github.com/gabrielfreitass1">Gabriel Freitas</a>, 2026.</p></font> 
</div>

- **(P)** — Web Server + gRPC Stub (API Gateway): recebe HTTP e repassa via gRPC
- **(A)** — gRPC Server: microserviço de desabafos
- **(B)** — gRPC Server: microserviço de reações ❤️

---

## Requisitos obrigatórios

- A e B devem ser implementados em linguagem **diferente** de P
- A e B devem ser **diferentes entre si**
- P deve ter uma **interface web** acessível pelo browser
- Os três módulos rodam em **containers separados** no minikube

---

## O que deve ser entregue

| Item | Descrição |
|---|---|
| Relatório | Documento com todas as seções abaixo |
| Arquivos de configuração | Tudo necessário para replicar o laboratório |
| Vídeo | ~4 minutos por aluno demonstrando participação |
| Documentação de código | Comentários e instruções de execução |

### Seções do relatório

1. Identificação do grupo
2. Introdução
3. **Framework gRPC** — protobuf, HTTP/2 e os 4 tipos de comunicação (com testes)
4. **Aplicação gRPC** — descrição, dificuldades e comparativo com REST/JSON (com medição de performance)
5. **Kubernetes** — arquitetura, comandos utilizados e configurações
6. Conclusão individual por aluno + autoavaliação
7. Apêndice (opcional)

---

## Tipos de comunicação gRPC a demonstrar

| Tipo | Descrição |
|---|---|
| Unary | 1 request → 1 response |
| Server Streaming | 1 request → N responses |
| Client Streaming | N requests → 1 response |
| Bidirecional | N requests ↔ N responses |

---

## Arquivo `.proto` do projeto

```protobuf
//colocar o .proto aqui
```

## Histórico de Versão
| Versão | Data       | Descrição                                      | Autor               | Revisor               |
|--------|------------|------------------------------------------------|---------------------|-----------------------|
| 1.0    | 04/06/2026 | Primeira versão do artefato grpc | [Milena Baruc Rodrigues Morais](https://github.com/MilenaBaruc) | [Milena Baruc Rodrigues Morais](https://github.com/MilenaBaruc) |
| 1.1    | 06/06/2026 | Atualizando Imagens | [Gabriel Freitas Balbino](https://github.com/gabrielfreitass1) | [Gabriel Freitas Balbino](https://github.com/gabrielfreitass1) |