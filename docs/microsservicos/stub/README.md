# Stub FastAPI + gRPC

Este servico expõe endpoints HTTP com FastAPI e chama os microsservicos gRPC:

- `ServidorA`, em Go, usando `docs/microsservicos/ServidorA/proto`
- `servidor-b`, em Java, usando `docs/microsservicos/servidor-b/src/main/proto`

## Instalar dependencias

```powershell
cd docs\microsservicos\stub
python -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt
```

## Gerar stubs Python

Sempre que algum `.proto` mudar, rode:

```powershell
cd docs\microsservicos\stub
python scripts\generate_stubs.py
```

Os arquivos gerados ficam em:

```text
app/grpc_generated
```

## Configurar portas

Defina as variaveis de ambiente antes de rodar, se as portas forem diferentes:

```powershell
$env:SERVIDOR_A_HOST="localhost:50051"
$env:SERVIDOR_B_HOST="localhost:50052"
```

Ou mantenha os valores padrao(tambem presente na .env do projeto):

```text
SERVIDOR_A_HOST=localhost:50051
SERVIDOR_B_HOST=localhost:50052
```

## Rodar

```powershell
cd docs\microsservicos\stub
uvicorn app.main:app --reload --port 8000
```

Depois acesse:

```text
http://localhost:8000/docs
```
