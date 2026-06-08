"""
gerar_relatorio.py
──────────────────
Lê resultado_grpc.json e resultado_rest.json gerados pelo k6
e imprime a tabela comparativa de tempo de resposta médio (ms).

Uso:
    python gerar_relatorio.py
"""

import json
import sys
from pathlib import Path

# ── Helpers ───────────────────────────────────────────────────

def carregar(caminho: str) -> dict:
    p = Path(caminho)
    if not p.exists():
        print(f"[ERRO] Arquivo não encontrado: {caminho}")
        sys.exit(1)
    with open(p) as f:
        return json.load(f)

def metrica(dados: dict, nome: str, campo: str = "avg") -> float:
    """Extrai um campo de uma métrica do sumário do k6."""
    try:
        return round(dados["metrics"][nome]["values"][campo], 2)
    except (KeyError, TypeError):
        return 0.0

def p95(dados: dict, nome: str) -> float:
    try:
        return round(dados["metrics"][nome]["values"]["p(95)"], 2)
    except (KeyError, TypeError):
        return 0.0

# ── Main ──────────────────────────────────────────────────────

def main():
    grpc = carregar("resultado_grpc.json")
    rest = carregar("resultado_rest.json")

    rotas = [
        ("POST /publicar", "grpc_publicar_ms", "rest_publicar_ms"),
        ("GET  /feed",     "grpc_feed_ms",     "rest_feed_ms"),
        ("POST /reagir",   "grpc_reagir_ms",   "rest_reagir_ms"),
    ]

    # ── Tabela comparativa ────────────────────────────────────
    col = 22
    print()
    print("=" * 80)
    print("  COMPARATIVO DE PERFORMANCE — gRPC vs REST/JSON")
    print("  Métrica: Tempo de resposta médio (ms)")
    print("=" * 80)
    print()

    header = (
        f"{'Rota':<{col}}"
        f"{'gRPC avg (ms)':>16}"
        f"{'gRPC p95 (ms)':>16}"
        f"{'REST avg (ms)':>16}"
        f"{'REST p95 (ms)':>16}"
        f"{'Δ avg (ms)':>14}"
        f"{'Vencedor':>12}"
    )
    sep = "-" * len(header)

    print(header)
    print(sep)

    for rota, metrica_grpc, metrica_rest in rotas:

        g_avg = metrica(grpc, metrica_grpc, "avg")
        g_p95 = p95(grpc, metrica_grpc)
        r_avg = metrica(rest, metrica_rest, "avg")
        r_p95 = p95(rest, metrica_rest)

        delta   = round(r_avg - g_avg, 2)
        vencedor = "gRPC" if g_avg <= r_avg else "REST"
        sinal    = f"+{delta}" if delta > 0 else str(delta)

        print(
            f"{rota:<{col}}"
            f"{g_avg:>16.2f}"
            f"{g_p95:>16.2f}"
            f"{r_avg:>16.2f}"
            f"{r_p95:>16.2f}"
            f"{sinal:>14}"
            f"{vencedor:>12}"
        )

    print(sep)
    print()

    # ── Métricas gerais ───────────────────────────────────────
    print("  MÉTRICAS GERAIS")
    print(sep)

    labels = [
        ("Duração total (s)",     "http_req_duration", "avg"),
        ("Requisições/s",         "http_reqs",         "rate"),
        ("Taxa de erros (%)",     "http_req_failed",   "rate"),
    ]

    for label, met, campo in labels:
        g_val = metrica(grpc, met, campo)
        r_val = metrica(rest, met, campo)

        if "erro" in label.lower():
            g_val = round(g_val * 100, 2)
            r_val = round(r_val * 100, 2)

        print(f"  {label:<30}  gRPC: {g_val:>10}   REST: {r_val:>10}")

    print()
    print("=" * 80)
    print()

    # ── Conclusão automática ───────────────────────────────────
    avgs_grpc = [metrica(grpc, m, "avg") for _, m, _ in rotas]
    avgs_rest = [metrica(rest, m, "avg") for _, _, m in rotas]
    media_grpc = round(sum(avgs_grpc) / len(avgs_grpc), 2)
    media_rest = round(sum(avgs_rest) / len(avgs_rest), 2)
    diff_pct   = round(abs(media_grpc - media_rest) / max(media_grpc, media_rest) * 100, 1)
    vencedor_geral = "gRPC" if media_grpc < media_rest else "REST/JSON"

    print("  CONCLUSÃO")
    print(sep)
    print(f"""
  Média geral gRPC : {media_grpc} ms
  Média geral REST : {media_rest} ms
  Diferença        : {diff_pct}% a favor de {vencedor_geral}

  A versão {vencedor_geral} apresentou menor tempo de resposta médio
  nos cenários testados. {'O gRPC leva vantagem especialmente no feed' if vencedor_geral == 'gRPC' else 'O REST/JSON se mostrou competitivo nos cenários testados'}
  (streaming), onde o HTTP/2 com Protobuf reduz overhead de
  serialização em relação ao JSON sobre HTTP/1.1.
    """)
    print("=" * 80)
    print()


if __name__ == "__main__":
    main()