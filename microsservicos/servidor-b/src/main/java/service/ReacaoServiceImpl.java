package service;

import io.grpc.stub.StreamObserver;

import repository.ReacaoRepository;

import reacoes.ReacaoConsulta;
import reacoes.ReacaoRequest;
import reacoes.ReacaoResposta;
import reacoes.ResumoReacoes;
import reacoes.ReacaoServiceGrpc;

public class ReacaoServiceImpl extends ReacaoServiceGrpc.ReacaoServiceImplBase {

    private final ReacaoRepository repository =
            new ReacaoRepository();

    @Override
    public StreamObserver<ReacaoRequest> enviarReacoes(
            StreamObserver<ResumoReacoes> responseObserver) {

        return new StreamObserver<ReacaoRequest>() {

            int total = 0;

            @Override
            public void onNext(ReacaoRequest request) {

                repository.adicionar(
                        request.getDesabafoId()
                );

                total++;
            }

            @Override
            public void onError(Throwable throwable) {

                System.err.println(
                        "Erro ao receber reações: "
                                + throwable.getMessage()
                );
            }

            @Override
            public void onCompleted() {

                responseObserver.onNext(
                        ResumoReacoes.newBuilder()
                                .setTotal(total)
                                .build()
                );

                responseObserver.onCompleted();
            }
        };
    }

    @Override
    public void buscarReacoes(
            ReacaoConsulta request,
            StreamObserver<ReacaoResposta> responseObserver) {

        int quantidade =
                repository.buscar(
                        request.getDesabafoId()
                );

        responseObserver.onNext(
                ReacaoResposta.newBuilder()
                        .setQuantidade(quantidade)
                        .build()
        );

        responseObserver.onCompleted();
    }
}