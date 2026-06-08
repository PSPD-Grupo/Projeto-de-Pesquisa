import grpc

from app.config import settings
from app.grpc_generated import desabafo_pb2, desabafo_pb2_grpc


class FeedClient:
    def __init__(self, target: str | None = None):
        self.target = target or settings.servidor_a_host

    def postar_desabafo(self, texto: str) -> dict:
        with grpc.insecure_channel(self.target) as channel:
            stub = desabafo_pb2_grpc.FeedStub(channel)
            response = stub.PostDesabafo(desabafo_pb2.RascunhoDesabafo(texto=texto))
            return {
                "id": response.id,
                "texto": response.texto,
                "created_at": response.created_at,
            }

    def listar_feed(self, quant: int) -> list[dict]:
        with grpc.insecure_channel(self.target) as channel:
            stub = desabafo_pb2_grpc.FeedStub(channel)
            stream = stub.GetFeed(desabafo_pb2.FeedRequest(quant=quant))
            return [
                {
                    "id": item.id,
                    "texto": item.texto,
                    "created_at": item.created_at,
                }
                for item in stream
            ]
