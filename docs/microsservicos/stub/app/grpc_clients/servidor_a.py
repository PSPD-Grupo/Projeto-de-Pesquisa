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
            
            # Decorate with reactions from Servidor B
            from app.grpc_clients.servidor_b import ReacaoClient
            reacao_client = ReacaoClient()
            
            items = []
            for item in stream:
                try:
                    reactions_data = reacao_client.buscar_reacoes(str(item.id))
                    reactions_count = reactions_data.get("quantidade", 0)
                except Exception as e:
                    print(f"Error fetching reactions for desabafo {item.id}: {e}")
                    reactions_count = 0
                
                items.append({
                    "id": item.id,
                    "texto": item.texto,
                    "created_at": item.created_at,
                    "reactions": reactions_count,
                })
            return items
