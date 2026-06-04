# PSPD – Projeto de Pesquisa

**UnB/FCTE – Engenharia de Software**  
Programação para Sistemas Paralelos e Distribuídos · Prof. Fernando W. Cruz

---

## O que é este projeto?

Este site documenta o **Projeto de Pesquisa (Parte 1)** da disciplina PSPD, que envolve dois grandes objetivos:

1. **Estudar e demonstrar o framework gRPC** — seus componentes (protobuf, HTTP/2) e os quatro tipos de comunicação suportados.
2. **Construir e implantar uma aplicação distribuída** baseada em microserviços, com deploy no Kubernetes (minikube).

---

## Nossa aplicação: *Desabafo* 🗒️

A aplicação escolhida pelo grupo é um **diário de desabafos anônimos** — uma espécie de Twitter minimalista onde o usuário pode:

- Publicar um desabafo (texto curto)
- Ver um feed de desabafos recentes
- Receber reações de outros usuários (❤️)

Simples de entender, mas com comunicação distribuída real entre os módulos gRPC.

### "Por que essa aplicação?"
    A escolha foi intencional: o fluxo de publicação e leitura de posts é um caso natural para demonstrar **todos os quatro tipos de comunicação gRPC** (unary, server-streaming, client-streaming e bidirecional).

---

## Navegue pela documentação

| Seção | Conteúdo |
|---|---|
| [gRPC →](grpc/index.md) | Teoria, componentes e exemplos dos 4 tipos de comunicação |
| [Microsserviços →](microsservicos/index.md) | Nossa aplicação, arquitetura e planejamento |
| [Sobre o Grupo →](sobre.md) | Membros e divisão de tarefas |