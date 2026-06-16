import grpc

from app.config import settings
from app.grpc_generated import reacoes_pb2, reacoes_pb2_grpc


class ReacaoClient:
    def __init__(self, target: str | None = None):
        self.target = target or settings.servidor_b_host

    def buscar_reacoes(self, desabafo_id: str) -> dict:
        with grpc.insecure_channel(self.target) as channel:
            stub = reacoes_pb2_grpc.ReacaoServiceStub(channel)
            response = stub.BuscarReacoes(
                reacoes_pb2.ReacaoConsulta(desabafo_id=desabafo_id)
            )
            return {
                "desabafo_id": desabafo_id,
                "quantidade": response.quantidade,
            }

    def reagir(self, desabafo_id: str) -> dict:
        with grpc.insecure_channel(self.target) as channel:
            stub = reacoes_pb2_grpc.ReacaoServiceStub(channel)
            
            def request_generator():
                yield reacoes_pb2.ReacaoRequest(desabafo_id=desabafo_id)
                
            response = stub.EnviarReacoes(request_generator())
            return {
                "desabafo_id": desabafo_id,
                "total": response.total,
            }
