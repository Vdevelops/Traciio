// =============================================================
// SOAK TEST — Long duration to detect memory leaks & conn leaks
// Duration: ~34 min at sustained VUs
// =============================================================

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { ENDPOINTS, authHeaders } from '../lib/config.js';
import { login } from '../lib/auth.js';
import { handleSummaryFor } from '../lib/reporter.js';

const soakLatency = new Trend('soak_latency_ms', true);
const soakErrors = new Rate('soak_error_rate');

const SOAK_VUS = parseInt(__ENV.SOAK_VUS || '15');

export const options = {
  stages: [
    { duration: '2m', target: SOAK_VUS },      // Ramp up
    { duration: '30m', target: SOAK_VUS },      // Sustained
    { duration: '2m', target: 0 },              // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<800', 'p(99)<1500'],
    http_req_failed: ['rate<0.05'],
    soak_error_rate: ['rate<0.05'],
  },
};

export function setup() {
  console.log(`━━━ SOAK TEST | ${SOAK_VUS} VUs for 30 min sustained ━━━`);
  console.log('Purpose: Detect memory leaks, connection pool exhaustion, gradual degradation');
  const tokens = login('admin');
  if (!tokens) return { token: null };
  return { token: tokens.accessToken };
}

export default function (data) {
  if (!data.token) return;
  const p = authHeaders(data.token);

  const rand = Math.random();

  let r;
  if (rand < 0.15) {
    // 15% — Health (lightweight baseline)
    r = http.get(ENDPOINTS.health);
  } else if (rand < 0.27) {
    // 12% — Leads list
    r = http.get(`${ENDPOINTS.leads}?page=${Math.floor(Math.random() * 3) + 1}&per_page=10`, p);
  } else if (rand < 0.37) {
    // 10% — Dashboard Overview (heavy query)
    r = http.get(ENDPOINTS.dashboardOverview, p);
  } else if (rand < 0.45) {
    // 8% — Profile
    r = http.get(ENDPOINTS.myProfile, p);
  } else if (rand < 0.52) {
    // 7% — Accounts
    r = http.get(`${ENDPOINTS.accounts}?page=1&per_page=10`, p);
  } else if (rand < 0.59) {
    // 7% — Contacts
    r = http.get(`${ENDPOINTS.contacts}?page=1&per_page=10`, p);
  } else if (rand < 0.66) {
    // 7% — Deals
    r = http.get(`${ENDPOINTS.deals}?page=1&per_page=10`, p);
  } else if (rand < 0.73) {
    // 7% — Visit Reports
    r = http.get(`${ENDPOINTS.visitReports}?page=1&per_page=10`, p);
  } else if (rand < 0.79) {
    // 6% — Tasks
    r = http.get(`${ENDPOINTS.tasks}?page=1&per_page=10`, p);
  } else if (rand < 0.84) {
    // 5% — Schedules
    r = http.get(`${ENDPOINTS.schedules}?page=1&per_page=10`, p);
  } else if (rand < 0.89) {
    // 5% — Activities
    r = http.get(`${ENDPOINTS.activities}?page=1&per_page=10`, p);
  } else if (rand < 0.93) {
    // 4% — Pipeline Summary
    r = http.get(ENDPOINTS.pipelineSummary, p);
  } else if (rand < 0.96) {
    // 3% — Notifications
    r = http.get(`${ENDPOINTS.notifications}?page=1&per_page=10`, p);
  } else {
    // 4% — Create lead (write)
    const payload = JSON.stringify({
      first_name: 'K6',
      last_name: `Soak ${__VU}-${Date.now()}`,
      email: `k6soak${__VU}${Date.now()}@test.com`,
      phone: `08${Math.floor(1000000000 + Math.random() * 9000000000)}`,
      lead_source: 'website',
      lead_status: 'new',
      notes: 'K6 soak test',
    });
    r = http.post(ENDPOINTS.leads, payload, p);
  }

  soakLatency.add(r.timings.duration);
  const ok = check(r, {
    'status ok': (r) => r.status === 200 || r.status === 201,
  });
  soakErrors.add(!ok);

  sleep(Math.random() * 3 + 1); // 1-4s think time
}

export function teardown() {
  console.log('━━━ SOAK TEST COMPLETE ━━━');
  console.log('Key questions:');
  console.log('  1. Did response times increase over time? (memory leak)');
  console.log('  2. Did error rate increase over time? (connection pool leak)');
  console.log('  3. Was throughput consistent throughout?');
}

export function handleSummary(data) {
  return handleSummaryFor('soak_test', data);
}
