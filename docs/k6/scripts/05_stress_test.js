// =============================================================
// STRESS TEST — Find the breaking point
// Gradually increase to MAX_VUS to find when API degrades.
// Duration: ~20 min
// =============================================================

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { ENDPOINTS, authHeaders } from '../lib/config.js';
import { login } from '../lib/auth.js';
import { handleSummaryFor } from '../lib/reporter.js';

const stressLatency = new Trend('stress_latency_ms', true);
const stressErrors = new Rate('stress_error_rate');

const MAX = parseInt(__ENV.MAX_VUS || '100');

export const options = {
  stages: [
    { duration: '1m', target: Math.ceil(MAX * 0.2) },   // Warm up 20%
    { duration: '2m', target: Math.ceil(MAX * 0.2) },   // Hold
    { duration: '1m', target: Math.ceil(MAX * 0.5) },   // 50%
    { duration: '3m', target: Math.ceil(MAX * 0.5) },   // Hold
    { duration: '1m', target: Math.ceil(MAX * 0.8) },   // 80%
    { duration: '3m', target: Math.ceil(MAX * 0.8) },   // Hold
    { duration: '1m', target: MAX },                     // 100%
    { duration: '3m', target: MAX },                     // Hold at peak
    { duration: '2m', target: 0 },                       // Recovery
  ],
  thresholds: {
    http_req_duration: ['p(95)<3000'],      // Lenient for stress
    http_req_failed: ['rate<0.20'],          // Up to 20% failure expected
    stress_error_rate: ['rate<0.25'],
  },
};

export function setup() {
  console.log(`━━━ STRESS TEST | Max VUs: ${MAX} | ~20 min ━━━`);
  console.log('⚠️  This test intentionally pushes the server to its limits!');
  const tokens = login('admin');
  if (!tokens) return { token: null };
  return { token: tokens.accessToken };
}

export default function (data) {
  if (!data.token) return;
  const p = authHeaders(data.token);

  const rand = Math.random();

  if (rand < 0.20) {
    // 20% — Health (lightweight baseline)
    const r = http.get(ENDPOINTS.health);
    stressLatency.add(r.timings.duration);
    stressErrors.add(!check(r, { 'health: ok': (r) => r.status === 200 }));
  } else if (rand < 0.35) {
    // 15% — Leads list (common read)
    const r = http.get(`${ENDPOINTS.leads}?page=1&per_page=10`, p);
    stressLatency.add(r.timings.duration);
    stressErrors.add(!check(r, { 'leads: ok': (r) => r.status === 200 }));
  } else if (rand < 0.48) {
    // 13% — Dashboard Overview (heavy)
    const r = http.get(ENDPOINTS.dashboardOverview, p);
    stressLatency.add(r.timings.duration);
    stressErrors.add(!check(r, { 'dashboard: ok': (r) => r.status === 200 }));
  } else if (rand < 0.58) {
    // 10% — Profile
    const r = http.get(ENDPOINTS.myProfile, p);
    stressLatency.add(r.timings.duration);
    stressErrors.add(!check(r, { 'profile: ok': (r) => r.status === 200 }));
  } else if (rand < 0.67) {
    // 9% — Accounts (was contacts)
    const r = http.get(`${ENDPOINTS.accounts}?page=1&per_page=10`, p);
    stressLatency.add(r.timings.duration);
    stressErrors.add(!check(r, { 'accounts: ok': (r) => r.status === 200 }));
  } else if (rand < 0.75) {
    // 8% — Contacts
    const r = http.get(`${ENDPOINTS.contacts}?page=1&per_page=10`, p);
    stressLatency.add(r.timings.duration);
    stressErrors.add(!check(r, { 'contacts: ok': (r) => r.status === 200 }));
  } else if (rand < 0.82) {
    // 7% — Deals
    const r = http.get(`${ENDPOINTS.deals}?page=1&per_page=10`, p);
    stressLatency.add(r.timings.duration);
    stressErrors.add(!check(r, { 'deals: ok': (r) => r.status === 200 }));
  } else if (rand < 0.88) {
    // 6% — Visit Reports
    const r = http.get(`${ENDPOINTS.visitReports}?page=1&per_page=10`, p);
    stressLatency.add(r.timings.duration);
    stressErrors.add(!check(r, { 'visit-reports: ok': (r) => r.status === 200 }));
  } else if (rand < 0.93) {
    // 5% — Tasks
    const r = http.get(`${ENDPOINTS.tasks}?page=1&per_page=10`, p);
    stressLatency.add(r.timings.duration);
    stressErrors.add(!check(r, { 'tasks: ok': (r) => r.status === 200 }));
  } else {
    // 7% — Create lead (write under stress)
    const payload = JSON.stringify({
      first_name: 'K6',
      last_name: `Stress ${__VU}-${Date.now()}`,
      email: `k6stress${__VU}${Date.now()}@test.com`,
      phone: `08${Math.floor(1000000000 + Math.random() * 9000000000)}`,
      lead_source: 'website',
      lead_status: 'new',
      notes: 'K6 stress test',
    });
    const r = http.post(ENDPOINTS.leads, payload, p);
    stressLatency.add(r.timings.duration);
    stressErrors.add(!check(r, { 'create: ok': (r) => r.status === 200 || r.status === 201 }));
  }

  sleep(Math.random() * 0.5 + 0.3); // Minimal think time
}

export function teardown() {
  console.log('━━━ STRESS TEST COMPLETE ━━━');
  console.log('Check: At what VU count did response times spike?');
  console.log('Check: At what VU count did error rate increase?');
}

export function handleSummary(data) {
  return handleSummaryFor('stress_test', data);
}
