package server;

import com.sun.net.httpserver.HttpServer;

import controller.ReacaoController;
import repository.ReacaoRepository;

import java.net.InetSocketAddress;
import java.util.concurrent.Executors;

/**
 * Servidor REST/JSON — espelho do ReacaoServer gRPC original.
 *
 * Porta : 8080  (era 50052 no gRPC)
 * Rotas :
 *   POST /reacoes/enviar   →  EnviarReacoes (stream → array JSON)
 *   GET  /reacoes/buscar   →  BuscarReacoes  (unary → query param)
 *
 * Usa apenas java.net.httpserver (JDK embutido) — sem frameworks
 * externos, mantendo o mesmo espírito minimalista do projeto gRPC.
 */
public class ReacaoServer {

    private static final int PORT = 8080;

    public static void main(String[] args) throws Exception {

        ReacaoRepository   repository  = new ReacaoRepository();
        ReacaoController   controller  = new ReacaoController(repository);

        HttpServer server = HttpServer.create(
                new InetSocketAddress(PORT), /*backlog*/ 0
        );

        server.createContext("/reacoes/enviar", controller.enviarReacoes());
        server.createContext("/reacoes/buscar", controller.buscarReacoes());

        // Thread pool para suportar requisições simultâneas
        server.setExecutor(Executors.newFixedThreadPool(4));

        server.start();

        System.out.println("Servidor B (REST) iniciado na porta " + PORT);
        System.out.println("  POST /reacoes/enviar");
        System.out.println("  GET  /reacoes/buscar?desabafo_id=<id>");

        // Mantém a JVM viva (mesmo comportamento do awaitTermination() gRPC)
        Thread.currentThread().join();
    }
}