package server;

import io.grpc.Server;
import io.grpc.ServerBuilder;

import service.ReacaoServiceImpl;

public class ReacaoServer {

    public static void main(String[] args)
        throws Exception {

        Server server =
            ServerBuilder
                .forPort(50052)
                .addService(
                    new ReacaoServiceImpl()
                )
                .build();

        server.start();

        System.out.println(
            "Servidor B iniciado na porta 50052"
        );

        server.awaitTermination();
    }
}