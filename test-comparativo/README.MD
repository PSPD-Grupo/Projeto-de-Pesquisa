# Testes de Performance — gRPC vs REST/JSON

## Pré-requisitos

```bash
# Instalar k6
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg \
  --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" \
  | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt update && sudo apt install k6

# Verificar
k6 version
```

---

## Estrutura de arquivos

```
testes-comparativo/
├── grpc/
│   └── test_grpc.js       ← testa versão gRPC (P na porta 8000)
├── rest/
│   └── test_rest.js       ← testa versão REST (P na porta 8001)
├── gerar_relatorio.py     ← gera tabela comparativa
└── README.md
```

---

## Passo a passo

### 1. Subir as duas versões no Minikube

```bash
# Terminal 1 — versão gRPC
kubectl port-forward service/servidor-p-grpc 8000:8000

# Terminal 2 — versão REST
kubectl port-forward service/servidor-p-rest 8001:8000
```

> ⚠️ Ajuste os nomes dos services conforme seu deployment.yaml

---

### 2. Confirmar que os dois estão respondendo

```bash
curl http://localhost:8000/feed   # deve retornar JSON
curl http://localhost:8001/feed   # deve retornar JSON
```

---

### 3. Rodar os testes (em sequência para não disputar recursos)

```bash
# Teste gRPC
cd grpc
k6 run test_grpc.js
cd ..

# Aguarde terminar (~5 min), depois:

# Teste REST
cd rest
k6 run test_rest.js
cd ..
```

Os arquivos `resultado_grpc.json` e `resultado_rest.json` serão gerados
automaticamente nas respectivas pastas.

---

### 4. Mover os resultados para a raiz e gerar o relatório

```bash
mv grpc/resultado_grpc.json .
mv rest/resultado_rest.json .

python3 gerar_relatorio.py
```

---

## Cenários de teste

| Cenário | Usuários | Duração | Objetivo |
|---|---|---|---|
| C1 — Baseline | 1 VU | 100 iterações | Latência mínima sem concorrência |
| C2 — Moderado | 10 VUs | 30s | Uso típico |
| C3 — Alto | 50 VUs | 30s | Carga elevada |
| C4 — Pico | 0→100 VUs | 75s | Comportamento sob pico |

---

## O que cada iteração faz

Cada VU executa em loop:
1. `POST /publicar` — publica um desabafo
2. `GET /feed` — busca o feed
3. `POST /reagir` — envia reações

Isso garante que **os mesmos fluxos** são testados nas duas versões.

---

## Ajustes necessários

| Arquivo | O que ajustar |
|---|---|
| `test_grpc.js` | `BASE` = porta do seu Módulo P versão gRPC |
| `test_rest.js` | `BASE` = porta do seu Módulo P versão REST |
| Ambos | Estrutura do body se suas rotas esperarem campos diferentes |