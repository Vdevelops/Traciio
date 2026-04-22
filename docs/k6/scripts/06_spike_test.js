// =============================================================
// SPIKE TEST — Sudden traffic spike, test recovery
// Simulates: normal → instant spike → back to normal
// Duration: ~5 min
// =============================================================

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { ENDPOINTS, authHeaders } from '../lib/config.js';
import { login } from '../lib/auth.js';
import { handleSummaryFor } from '../lib/reporter.js';

const spikeLatency = new Trend('spike_latency_ms', true);
const spikeErrors = new Rate('spike_error_rate');

const SPIKE_VUS = parseInt(__ENV.SPIKE_VUS || '80');
const NORMAL_VUS = parseInt(__ENV.NORMAL_VUS || '5');

export const options = {
  stages: [
    { duration: '30s', target: NORMAL_VUS },     // Normal
    { duration: '1m', target: NORMAL_VUS },       // Steady
    { duration: '10s', target: SPIKE_VUS },       // 🔥 SPIKE!
    { duration: '1m', target: SPIKE_VUS },        // Hold spike
    { duration: '10s', target: NORMAL_VUS },      // Drop
    { duration: '1m30s', target: NORMAL_VUS },    // Recovery
    { duration: '20s', target: 0 },               // Done
  ],
  thresholds: {
    http_req_duration: ['p(95)<5000'],       // Very lenient
    spike_error_rate: ['rate<0.40'],          // Spikes can cause errors
  },
};

export function setup() {
  console.log(`━━━ SPIKE TEST | Normal: ${NORMAL_VUS} → Spike: ${SPIKE_VUS} VUs ━━━`);
  const tokens = login('admin');
  if (!tokens) return { token: null };
  return { token: tokens.accessToken };
}

export default function (data) {
  if (!data.token) return;
  const p = authHeaders(data.token);

  const endpoints = [
    () => http.get(ENDPOINTS.health),
    () => http.get(`${ENDPOINTS.leads}?page=1&per_page=10`, p),
    () => http.get(ENDPOINTS.myProfile, p),
    () => http.get(ENDPOINTS.dashboardOverview, p),
    () => http.get(`${ENDPOINTS.accounts}?page=1&per_page=5`, p),
    () => http.get(`${ENDPOINTS.contacts}?page=1&per_page=5`, p),
    () => http.get(`${ENDPOINTS.deals}?page=1&per_page=5`, p),
    () => http.get(`${ENDPOINTS.visitReports}?page=1&per_page=5`, p),
    () => http.get(`${ENDPOINTS.tasks}?page=1&per_page=5`, p),
  ];

  const r = endpoints[Math.floor(Math.random() * endpoints.length)]();

  spikeLatency.add(r.timings.duration);
  spikeErrors.add(!check(r, {
    'response ok': (r) => r.status === 200,
    'not timeout': (r) => r.timings.duration < 10000,
  }));

  sleep(Math.random() * 0.3 + 0.1);
}

export function teardown() {
  console.log('━━━ SPIKE TEST COMPLETE ━━━');
  console.log('Key question: How quickly did response times recover after spike dropped?');
}

export function handleSummary(data) {
  return handleSummaryFor('spike_test', data);
}
