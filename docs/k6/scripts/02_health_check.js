// =============================================================
// HEALTH CHECK TEST — Throughput baseline
// Ramp up ke MAX_VUS, ukur throughput & response time murni.
// =============================================================

import http from 'k6/http';
import { check, sleep } from 'k6';
import { ENDPOINTS } from '../lib/config.js';
import { handleSummaryFor } from '../lib/reporter.js';

const MAX = parseInt(__ENV.MAX_VUS || '20');

export const options = {
  stages: [
    { duration: '15s', target: Math.ceil(MAX * 0.3) },
    { duration: '30s', target: Math.ceil(MAX * 0.3) },
    { duration: '15s', target: MAX },
    { duration: '1m', target: MAX },
    { duration: '15s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.05'],
    checks: ['rate>0.95'],
  },
};

export default function () {
  // Alternate between /health and /ping for variety
  const useHealth = Math.random() < 0.7;
  const url = useHealth ? ENDPOINTS.health : ENDPOINTS.ping;
  const r = http.get(url);

  check(r, {
    'health/ping: 200': (r) => r.status === 200,
    'health/ping: < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(Math.random() * 0.5 + 0.5); // 0.5 - 1s
}

export function handleSummary(data) {
  return handleSummaryFor('health_check', data);
}
