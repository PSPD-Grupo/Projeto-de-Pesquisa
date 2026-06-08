package controller;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;

import repository.ReacaoRepository;
import util.JsonUtil;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Handlers REST que espelham os dois RPCs do serviço gRPC original:
 *
 *   POST /reacoes/enviar   ←→  rpc EnviarReacoes (client-streaming)
 *   GET  /reacoes/buscar   ←→  rpc BuscarReacoes  (unary)
 *
 * Como HTTP/1.1 não tem streaming nativo igual ao gRPC, o endpoint
 * POST aceita um array JSON com múltiplas reações de uma só vez,
 * reproduzindo a semântica de "enviar várias e receber o total".
 */
public class ReacaoController {

    private final ReacaoRepository repository;

    public ReacaoController(ReacaoRepository repository) {
        this.repository = repository;
    }

    // ──────────────────────────────────────────────────────────
    // POST /reacoes/enviar
    //
    // Body esperado:
    //   {"desabafo_ids": ["id1", "id1", "id2"]}
    //
    // Resposta 200:
    //   {"total": 3}
    // ──────────────────────────────────────────────────────────
    public HttpHandler enviarReacoes() {

        return (HttpExchange exchange) -> {

            if (!exchange.getRequestMethod().equalsIgnoreCase("POST")) {
                sendResponse(exchange, 405, "{\"erro\":\"Método não permitido\"}");
                return;
            }

            String body = readBody(exchange);

            // Extrai o array de ids — formato: ["id1","id2",...]
            int arrayStart = body.indexOf('[');
            int arrayEnd   = body.lastIndexOf(']');

            if (arrayStart < 0 || arrayEnd < 0) {
                sendResponse(exchange, 400, "{\"erro\":\"Campo 'desabafo_ids' ausente ou inválido\"}");
                return;
            }

            String arrayContent = body.substring(arrayStart + 1, arrayEnd);
            int total = 0;

            if (!arrayContent.isBlank()) {

                String[] tokens = arrayContent.split(",");

                for (String token : tokens) {

                    String id = token.trim()
                                     .replace("\"", "")
                                     .replace("'", "");

                    if (!id.isBlank()) {
                        repository.adicionar(id);
                        total++;
                    }
                }
            }

            Map<String, Object> resp = new LinkedHashMap<>();
            resp.put("total", total);
            sendResponse(exchange, 200, JsonUtil.toJson(resp));
        };
    }

    // ──────────────────────────────────────────────────────────
    // GET /reacoes/buscar?desabafo_id=<id>
    //
    // Resposta 200:
    //   {"quantidade": 5}
    // ──────────────────────────────────────────────────────────
    public HttpHandler buscarReacoes() {

        return (HttpExchange exchange) -> {

            if (!exchange.getRequestMethod().equalsIgnoreCase("GET")) {
                sendResponse(exchange, 405, "{\"erro\":\"Método não permitido\"}");
                return;
            }

            String query = exchange.getRequestURI().getQuery();
            String desabafoId = extractQueryParam(query, "desabafo_id");

            if (desabafoId == null || desabafoId.isBlank()) {
                sendResponse(exchange, 400, "{\"erro\":\"Parâmetro 'desabafo_id' obrigatório\"}");
                return;
            }

            int quantidade = repository.buscar(desabafoId);

            Map<String, Object> resp = new LinkedHashMap<>();
            resp.put("quantidade", quantidade);
            sendResponse(exchange, 200, JsonUtil.toJson(resp));
        };
    }

    // ──────────────────────────────────────────────────────────
    // Helpers
    // ──────────────────────────────────────────────────────────

    private String readBody(HttpExchange exchange) throws IOException {

        try (InputStream is = exchange.getRequestBody()) {
            return new String(is.readAllBytes(), StandardCharsets.UTF_8);
        }
    }

    private void sendResponse(HttpExchange exchange, int status, String body)
            throws IOException {

        byte[] bytes = body.getBytes(StandardCharsets.UTF_8);

        exchange.getResponseHeaders().set("Content-Type", "application/json; charset=UTF-8");
        exchange.sendResponseHeaders(status, bytes.length);

        try (OutputStream os = exchange.getResponseBody()) {
            os.write(bytes);
        }
    }

    private String extractQueryParam(String query, String key) {

        if (query == null) return null;

        for (String param : query.split("&")) {
            String[] kv = param.split("=", 2);
            if (kv.length == 2 && kv[0].equals(key)) {
                return kv[1];
            }
        }

        return null;
    }
}