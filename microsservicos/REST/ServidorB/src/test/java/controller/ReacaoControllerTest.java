package controller;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.Headers;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import repository.ReacaoRepository;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.URI;
import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

class ReacaoControllerTest {

    private ReacaoRepository repository;
    private ReacaoController  controller;

    @BeforeEach
    void setUp() {
        repository = new ReacaoRepository();
        controller = new ReacaoController(repository);
    }

    // ── enviarReacoes ──────────────────────────────────────────

    @Test
    void enviarReacoes_deveContabilizarETotalizarCorretamente()
            throws Exception {

        String body = "{\"desabafo_ids\":[\"abc\",\"abc\",\"xyz\"]}";
        ByteArrayOutputStream responseBody = new ByteArrayOutputStream();

        HttpExchange exchange = mockPostExchange(body, responseBody);
        controller.enviarReacoes().handle(exchange);

        String response = responseBody.toString(StandardCharsets.UTF_8);
        assertTrue(response.contains("\"total\":3"),
                "Total deve ser 3, obtido: " + response);

        assertEquals(2, repository.buscar("abc"));
        assertEquals(1, repository.buscar("xyz"));
    }

    @Test
    void enviarReacoes_arrayVazioDeveRetornarTotalZero()
            throws Exception {

        String body = "{\"desabafo_ids\":[]}";
        ByteArrayOutputStream responseBody = new ByteArrayOutputStream();

        HttpExchange exchange = mockPostExchange(body, responseBody);
        controller.enviarReacoes().handle(exchange);

        String response = responseBody.toString(StandardCharsets.UTF_8);
        assertTrue(response.contains("\"total\":0"),
                "Total deve ser 0, obtido: " + response);
    }

    // ── buscarReacoes ──────────────────────────────────────────

    @Test
    void buscarReacoes_deveRetornarQuantidadeCorreta()
            throws Exception {

        repository.adicionar("desabafo-1");
        repository.adicionar("desabafo-1");

        ByteArrayOutputStream responseBody = new ByteArrayOutputStream();
        HttpExchange exchange = mockGetExchange("desabafo_id=desabafo-1", responseBody);

        controller.buscarReacoes().handle(exchange);

        String response = responseBody.toString(StandardCharsets.UTF_8);
        assertTrue(response.contains("\"quantidade\":2"),
                "Quantidade deve ser 2, obtido: " + response);
    }

    @Test
    void buscarReacoes_idInexistenteDeveRetornarZero()
            throws Exception {

        ByteArrayOutputStream responseBody = new ByteArrayOutputStream();
        HttpExchange exchange = mockGetExchange("desabafo_id=nao-existe", responseBody);

        controller.buscarReacoes().handle(exchange);

        String response = responseBody.toString(StandardCharsets.UTF_8);
        assertTrue(response.contains("\"quantidade\":0"),
                "Quantidade deve ser 0, obtido: " + response);
    }

    // ── Helpers ───────────────────────────────────────────────

    private HttpExchange mockPostExchange(String bodyContent,
                                          OutputStream out) throws Exception {

        HttpExchange exchange = mock(HttpExchange.class);
        InputStream  is       = new ByteArrayInputStream(
                bodyContent.getBytes(StandardCharsets.UTF_8));

        when(exchange.getRequestMethod()).thenReturn("POST");
        when(exchange.getRequestBody()).thenReturn(is);
        when(exchange.getResponseBody()).thenReturn(out);
        when(exchange.getResponseHeaders()).thenReturn(new Headers());

        return exchange;
    }

    private HttpExchange mockGetExchange(String query,
                                         OutputStream out) throws Exception {

        HttpExchange exchange = mock(HttpExchange.class);
        URI          uri      = new URI("/reacoes/buscar?" + query);

        when(exchange.getRequestMethod()).thenReturn("GET");
        when(exchange.getRequestURI()).thenReturn(uri);
        when(exchange.getResponseBody()).thenReturn(out);
        when(exchange.getResponseHeaders()).thenReturn(new Headers());

        return exchange;
    }
}