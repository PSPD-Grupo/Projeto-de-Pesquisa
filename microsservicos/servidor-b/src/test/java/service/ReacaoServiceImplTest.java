package service;

import io.grpc.stub.StreamObserver;

import org.junit.jupiter.api.Test;

import reacoes.*;

import java.util.concurrent.atomic.AtomicReference;

import static org.junit.jupiter.api.Assertions.*;

public class ReacaoServiceImplTest {

    @Test
    void deveReceberTresReacoes() {

        ReacaoServiceImpl service =
                new ReacaoServiceImpl();

        AtomicReference<ResumoReacoes> resposta =
                new AtomicReference<>();

        StreamObserver<ResumoReacoes> responseObserver =
                new StreamObserver<>() {

                    @Override
                    public void onNext(
                            ResumoReacoes value) {

                        resposta.set(value);
                    }

                    @Override
                    public void onError(
                            Throwable t) {
                    }

                    @Override
                    public void onCompleted() {
                    }
                };

        StreamObserver<ReacaoRequest> requestObserver =
                service.enviarReacoes(
                        responseObserver
                );

        requestObserver.onNext(
                ReacaoRequest.newBuilder()
                        .setDesabafoId("1")
                        .build()
        );

        requestObserver.onNext(
                ReacaoRequest.newBuilder()
                        .setDesabafoId("1")
                        .build()
        );

        requestObserver.onNext(
                ReacaoRequest.newBuilder()
                        .setDesabafoId("1")
                        .build()
        );

        requestObserver.onCompleted();

        assertEquals(
                3,
                resposta.get().getTotal()
        );
    }

    @Test
    void deveBuscarQuantidadeDeReacoes() {

        ReacaoServiceImpl service =
                new ReacaoServiceImpl();

        AtomicReference<ResumoReacoes> resumo =
                new AtomicReference<>();

        StreamObserver<ResumoReacoes> responseResumo =
                new StreamObserver<>() {

                    @Override
                    public void onNext(
                            ResumoReacoes value) {

                        resumo.set(value);
                    }

                    @Override
                    public void onError(
                            Throwable t) {
                    }

                    @Override
                    public void onCompleted() {
                    }
                };

        StreamObserver<ReacaoRequest> stream =
                service.enviarReacoes(
                        responseResumo
                );

        stream.onNext(
                ReacaoRequest.newBuilder()
                        .setDesabafoId("abc")
                        .build()
        );

        stream.onNext(
                ReacaoRequest.newBuilder()
                        .setDesabafoId("abc")
                        .build()
        );

        stream.onCompleted();

        AtomicReference<ReacaoResposta> resposta =
                new AtomicReference<>();

        service.buscarReacoes(

                ReacaoConsulta
                        .newBuilder()
                        .setDesabafoId("abc")
                        .build(),

                new StreamObserver<>() {

                    @Override
                    public void onNext(
                            ReacaoResposta value) {

                        resposta.set(value);
                    }

                    @Override
                    public void onError(
                            Throwable t) {
                    }

                    @Override
                    public void onCompleted() {
                    }
                }
        );

        assertEquals(
                2,
                resposta.get().getQuantidade()
        );
    }
}