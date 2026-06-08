from fastapi import FastAPI

from app.grpc_clients.servidor_a import FeedClient
from app.grpc_clients.servidor_b import ReacaoClient

app = FastAPI(title="Stub FastAPI gRPC", version="0.1.0")


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/desabafos")
def postar_desabafo(texto: str):
    return FeedClient().postar_desabafo(texto)


@app.get("/feed")
def listar_feed(quant: int = 10):
    return {"items": FeedClient().listar_feed(quant)}


@app.get("/desabafos/{desabafo_id}/reacoes")
def buscar_reacoes(desabafo_id: str):
    return ReacaoClient().buscar_reacoes(desabafo_id)
