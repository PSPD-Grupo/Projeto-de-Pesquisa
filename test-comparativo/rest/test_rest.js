/**
 * Teste de performance — Versão REST/JSON
 *
 * Aponta para a versão REST do Módulo P (porta 8001).
 * Estrutura de cenários IDÊNTICA ao test_grpc.js para comparação justa.
 *
 * Rodar: k6 run test_rest.js --out json=resultado_rest.json
 */

import http from "k6/http";
import { check, sleep } from "k6";
import { Trend } from "k6/metrics";

// ── Métricas customizadas por rota ────────────────────────────
const tendenciaPublicar = new Trend("rest_publicar_ms", true);
const tendenciaFeed     = new Trend("rest_feed_ms",     true);
const tendenciaReagir   = new Trend("rest_reagir_ms",   true);

// ── Cenários de carga (idênticos ao gRPC) ─────────────────────
export const options = {
  scenarios: {

    c1_baseline: {
      executor:    "per-vu-iterations",
      vus:         1,
      iterations:  100,
      maxDuration: "2m",
      tags:        { cenario: "C1-baseline" },
    },

    c2_moderado: {
      executor:  "constant-vus",
      vus:       10,
      duration:  "30s",
      startTime: "2m10s",
      tags:      { cenario: "C2-moderado" },
    },

    c3_alto: {
      executor:  "constant-vus",
      vus:       50,
      duration:  "30s",
      startTime: "3m",
      tags:      { cenario: "C3-alto" },
    },

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
    rest_publicar_ms: ["p(95)<2000"],
    rest_feed_ms:     ["p(95)<3000"],
    rest_reagir_ms:   ["p(95)<2000"],
  },
};

// ⚠️ Troque pela porta da sua versão REST do Módulo P
const BASE = "http://localhost:8001";

let desabafoIds = ["teste-rest-1", "teste-rest-2", "teste-rest-3"];

// ── Fluxo principal ───────────────────────────────────────────
export default function () {

  // 1. POST /publicar
  const tInicio = Date.now();
  const resPublicar = http.post(
    `${BASE}/publicar`,
    JSON.stringify({ texto: `Desabafo de teste REST - VU ${__VU}` }),
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
    "feed: status 200":   (r) => r.status === 200,
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
    "resultado_rest.json": JSON.stringify(data, null, 2),
  };
}