/**
 * Teste de performance — Versão gRPC (via Módulo P)
 *
 * Módulo P expõe HTTP na porta 8000 e traduz para gRPC internamente.
 * Testamos as mesmas rotas que a versão REST para comparação justa.
 *
 * Instalar k6: https://k6.io/docs/get-started/installation/
 * Rodar:  k6 run test_grpc.js --out json=resultado_grpc.json
 */

import http from "k6/http";
import { check, sleep } from "k6";
import { Trend } from "k6/metrics";

// ── Métricas customizadas por rota ────────────────────────────
const tendenciaPublicar = new Trend("grpc_publicar_ms", true);
const tendenciaFeed     = new Trend("grpc_feed_ms",     true);
const tendenciaReagir   = new Trend("grpc_reagir_ms",   true);

// ── Cenários de carga ─────────────────────────────────────────
export const options = {
  scenarios: {

    // C1 — Baseline: 1 usuário, 100 iterações sequenciais
    c1_baseline: {
      executor:    "per-vu-iterations",
      vus:         1,
      iterations:  100,
      maxDuration: "2m",
      tags:        { cenario: "C1-baseline" },
    },

    // C2 — Carga moderada: 10 usuários por 30s
    c2_moderado: {
      executor:  "constant-vus",
      vus:       10,
      duration:  "30s",
      startTime: "2m10s",   // começa após C1 terminar
      tags:      { cenario: "C2-moderado" },
    },

    // C3 — Carga alta: 50 usuários por 30s
    c3_alto: {
      executor:  "constant-vus",
      vus:       50,
      duration:  "30s",
      startTime: "3m",
      tags:      { cenario: "C3-alto" },
    },

    // C4 — Pico: rampa 0 → 100 usuários em 30s, sustenta 30s, desce
    c4_pico: {
      executor: "ramping-vus",
      startTime: "4m",
      stages: [
        { duration: "30s", target: 100 },
        { duration: "30s", target: 100 },
        { duration: "15s", target: 0   },
      ],
      tags: { cenario: "C4-pico" },
    },
  },

  thresholds: {
    // Limites aceitáveis (não falham o teste, só registram)
    grpc_publicar_ms: ["p(95)<2000"],
    grpc_feed_ms:     ["p(95)<3000"],
    grpc_reagir_ms:   ["p(95)<2000"],
  },
};

const BASE = "http://localhost:8000";

// IDs fixos para reagir e buscar (populados no setup)
let desabafoIds = ["teste-grpc-1", "teste-grpc-2", "teste-grpc-3"];

// ── Fluxo principal ───────────────────────────────────────────
export default function () {

  // 1. POST /publicar
  const tInicio = Date.now();
  const resPublicar = http.post(
    `${BASE}/publicar`,
    JSON.stringify({ texto: `Desabafo de teste gRPC - VU ${__VU}` }),
    { headers: { "Content-Type": "application/json" } }
  );
  tendenciaPublicar.add(Date.now() - tInicio);

  check(resPublicar, {
    "publicar: status 200": (r) => r.status === 200,
    "publicar: tem id":     (r) => {
      try { return JSON.parse(r.body).id !== undefined; }
      catch { return false; }
    },
  });

  sleep(0.1);

  // 2. GET /feed
  const tFeed = Date.now();
  const resFeed = http.get(`${BASE}/feed`);
  tendenciaFeed.add(Date.now() - tFeed);

  check(resFeed, {
    "feed: status 200":  (r) => r.status === 200,
    "feed: corpo válido": (r) => r.body.length > 0,
  });

  sleep(0.1);

  // 3. POST /reagir
  const id = desabafoIds[__VU % desabafoIds.length];
  const tReagir = Date.now();
  const resReagir = http.post(
    `${BASE}/reagir`,
    JSON.stringify({ desabafo_ids: [id, id] }),
    { headers: { "Content-Type": "application/json" } }
  );
  tendenciaReagir.add(Date.now() - tReagir);

  check(resReagir, {
    "reagir: status 200": (r) => r.status === 200,
  });

  sleep(0.1);
}

export function handleSummary(data) {
  return {
    "resultado_grpc.json": JSON.stringify(data, null, 2),
  };
}