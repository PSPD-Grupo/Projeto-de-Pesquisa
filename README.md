# PSPD – Projeto de Pesquisa

**UnB/FCTE – Engenharia de Software**  
Programação para Sistemas Paralelos e Distribuídos · Prof. Fernando W. Cruz

---

## 📖 Documentação

Acesse a documentação completa do projeto em:

**👉 [https://pspd-grupo.github.io/Projeto-de-Pesquisa/](https://pspd-grupo.github.io/Projeto-de-Pesquisa/)**

---

## 👥 Grupo

| Matrícula | Nome |
|---|---|
| 211062339 | Milena Baruc Rodrigues Morais |
| 212005444 | Pedro Fonseca Cruz |
| 211030980 | Daniel dos Santos Barros de Sousa |
| 180075462 | Gabriel Freitas Balbino |

---

## 🗂️ Estrutura do repositório

```
Projeto-de-Pesquisa/
├── docs/               # Fonte da documentação (MkDocs)
│   ├── grpc/           # Teoria e exemplos do framework gRPC
│   ├── microsservicos/ # Aplicação Desabafo e arquitetura
│   └── index.md        # Página inicial
├── mkdocs.yml          # Configuração do MkDocs Material
└── .github/workflows/  # Deploy automático para GitHub Pages
```

---

## 🚀 Sobre o projeto

O projeto envolve dois objetivos principais:

1. **Estudar o framework gRPC** — componentes (protobuf, HTTP/2) e os 4 tipos de comunicação
2. **Construir uma aplicação distribuída** — microserviços com deploy no Kubernetes (minikube)

A aplicação escolhida é o **Desabafo**, um diário anônimo distribuído onde usuários podem publicar textos curtos, ver um feed e interagir em tempo real.