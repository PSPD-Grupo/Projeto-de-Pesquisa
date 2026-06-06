package repository;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class ReacaoRepository {

    private final Map<String, Integer> reacoes =
            new ConcurrentHashMap<>();

    public void adicionar(String id) {

        reacoes.merge(id, 1, Integer::sum);
    }

    public int buscar(String id) {

        return reacoes.getOrDefault(id, 0);
    }
}