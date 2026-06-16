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
├── docs/                    # Fonte da documentação (MkDocs)
│   ├── grpc/                # Teoria e exemplos do framework gRPC
│   ├── microsservicos/      # Aplicação Desabafo e arquitetura
│   ├── especificacao.md     # O que o professor pediu + arquivo .proto
│   └── index.md             # Página inicial
├── proto/                   # Arquivo .proto compartilhado entre os módulos
├── modulo-p/                # API Gateway (Web Server + gRPC Stub)
├── servidor-a/              # gRPC Server — Desabafos
├── servidor-b/              # gRPC Server — Reações ❤️
├── mkdocs.yml               # Configuração do MkDocs Material
└── .github/workflows/       # Deploy automático para GitHub Pages
```

---

## 🚀 Sobre o projeto

O projeto envolve dois objetivos principais:

1. **Estudar o framework gRPC** — componentes (protobuf, HTTP/2) e os 4 tipos de comunicação
2. **Construir uma aplicação distribuída** — microserviços com deploy no Kubernetes (minikube)

A aplicação escolhida é o **Desabafo**, um diário anônimo distribuído onde usuários podem publicar textos curtos, ver um feed e interagir em tempo real.

--- 
## Screenshots Do Trabalho
![Imagem 1]()<br>



---

## 🛠️ Instalação de Dependências

O projeto é composto por três módulos principais rodando em diferentes linguagens:
- **Go** (Servidor A - Desabafos e Chat)
- **Java** (Servidor B - Reações)
- **Python / FastAPI** (Módulo P - API Gateway + Web Server que serve a interface estática)

### Pré-requisitos locais:
1. **Go Compiler** (Versão 1.22 ou superior)
2. **Java JDK** (Versão 21 ou superior) e **Maven** (Versão 3.9 ou superior)
3. **Python 3** (Versão 3.10 ou superior)
4. **Minikube** e **kubectl** (caso queira rodar no Kubernetes)

---

## 🚀 Execução do Projeto

### Opção 1: Executando Localmente (Recomendado para Teste Rápido)

Para rodar todo o sistema de maneira local e integrada na sua máquina:

#### 1. Iniciar o Servidor A (Go - Porta 50053):
Este servidor gerencia desabafos e disponibiliza a sala de chat bidirecional em tempo real. Ele está configurado para salvar os dados persistentemente em um banco de dados SQLite (`db/db.sqlite3`).
```bash
cd docs/microsservicos/ServidorA
go build -o servidor .
./servidor -port 50053
```

#### 2. Iniciar o Servidor B (Java - Porta 50052):
Este servidor gerencia e acumula as reações (corações ❤️) em tempo real na memória.
```bash
cd docs/microsservicos/servidor-b
mvn clean package -DskipTests
java -jar target/servidor-b.jar
```

#### 3. Iniciar o Módulo P (Python FastAPI - Porta 8000):
Este módulo atua como API Gateway (fazendo a ponte gRPC e agregando o feed com as reações) e serve a interface Web.
```bash
cd docs/microsservicos/stub
# Crie e ative um ambiente virtual
python3 -m venv .venv
source .venv/bin/activate  # No Windows use: .\.venv\Scripts\Activate.ps1
pip install -r requirements.txt

# Inicie o gateway passando os endereços gRPC locais como variáveis de ambiente
SERVIDOR_A_HOST="localhost:50053" SERVIDOR_B_HOST="localhost:50052" uvicorn app.main:app --port 8000
```

#### 4. Acessar a aplicação:
Abra o navegador no seguinte endereço:
👉 **[http://localhost:8000/](http://localhost:8000/)**

---

### Opção 2: Executando no Kubernetes (Minikube)

Para simular o ambiente de produção rodando cada microsserviço em pods isolados no K8s:

#### 1. Configurar o ambiente do Minikube:
```bash
minikube start
eval $(minikube docker-env)  # Conecta seu terminal local ao Docker interno do Minikube
```

#### 2. Buildar as imagens Docker para dentro do Minikube:
```bash
docker build -t servidora-servidor-a:latest docs/microsservicos/ServidorA/
docker build -t servidor-b-servidor-b:latest docs/microsservicos/servidor-b/
docker build -t stub-api:latest docs/microsservicos/stub/
```

#### 3. Realizar o Deploy dos manifestos:
```bash
kubectl apply -f docs/microsservicos/k8s/
```

#### 4. Obter a URL de acesso externa do Gateway:
```bash
minikube service stub-api --url
```
Abra o endereço gerado no navegador (ex: `http://192.168.49.2:32456`).

---

### 🎥 Link da apresentação
[Vídeo da Apresentação]()<br>