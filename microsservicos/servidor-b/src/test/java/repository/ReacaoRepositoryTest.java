package repository;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

public class ReacaoRepositoryTest {

    @Test
    void deveRetornarZeroQuandoNaoExistemReacoes() {

        ReacaoRepository repository =
                new ReacaoRepository();

        assertEquals(
                0,
                repository.buscar("desabafo-1")
        );
    }

    @Test
    void deveContarUmaReacao() {

        ReacaoRepository repository =
                new ReacaoRepository();

        repository.adicionar("desabafo-1");

        assertEquals(
                1,
                repository.buscar("desabafo-1")
        );
    }

    @Test
    void deveContarVariasReacoes() {

        ReacaoRepository repository =
                new ReacaoRepository();

        repository.adicionar("desabafo-1");
        repository.adicionar("desabafo-1");
        repository.adicionar("desabafo-1");

        assertEquals(
                3,
                repository.buscar("desabafo-1")
        );
    }
}