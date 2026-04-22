// =============================================================
// AUTH FLOW TEST — Login → Profile → Refresh → Logout lifecycle
// =============================================================

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { ENDPOINTS, jsonHeaders, authHeaders } from '../lib/config.js';
import { handleSummaryFor } from '../lib/reporter.js';

const loginDuration = new Trend('login_duration_ms', true);
const loginFailRate = new Rate('login_failures');

const MAX = parseInt(__ENV.MAX_VUS || '10');

export const options = {
  stages: [
    { duration: '15s', target: Math.ceil(MAX * 0.5) },
    { duration: '1m', target: Math.ceil(MAX * 0.5) },
    { duration: '15s', target: MAX },
    { duration: '1m', target: MAX },
    { duration: '15s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'],
    http_req_failed: ['rate<0.10'],
    login_duration_ms: ['p(95)<1500'],
    login_failures: ['rate<0.10'],
    checks: ['rate>0.80'],
  },
};

export default function () {
  let accessToken = null;
  let refreshTokenVal = null;

  // ── LOGIN ──
  group('01 Login', function () {
    const payload = JSON.stringify({
      email: 'admin@example.com',
      password: 'admin123',
    });

    const r = http.post(ENDPOINTS.login, payload, jsonHeaders());
    loginDuration.add(r.timings.duration);

    const ok = check(r, {
      'login: 200': (r) => r.status === 200,
      'login: has token': (r) => {
        try {
          return JSON.parse(r.body).data.token !== undefined;
        } catch (_) {
          return false;
        }
      },
    });

    loginFailRate.add(!ok);

    if (ok) {
      try {
        const body = JSON.parse(r.body);
        accessToken = body.data.token;
        refreshTokenVal = body.data.refresh_token;
      } catch (_) { }
    }
  });

  if (!accessToken) {
    sleep(2);
    return;
  }

  sleep(0.5);

  // ── GET PROFILE (settings-summary) ──
  group('02 Profile', function () {
    const r = http.get(ENDPOINTS.myProfile, authHeaders(accessToken));
    check(r, {
      'profile: 200': (r) => r.status === 200,
      'profile: has data': (r) => {
        try {
          return JSON.parse(r.body).data !== undefined;
        } catch (_) {
          return false;
        }
      },
    });
  });

  sleep(0.5);

  // ── REFRESH TOKEN ──
  if (refreshTokenVal) {
    group('03 Refresh Token', function () {
      const payload = JSON.stringify({ refresh_token: refreshTokenVal });
      const r = http.post(ENDPOINTS.refresh, payload, jsonHeaders());
      check(r, {
        'refresh: 200': (r) => r.status === 200,
        'refresh: has new token': (r) => {
          try {
            return JSON.parse(r.body).data.token !== undefined;
          } catch (_) {
            return false;
          }
        },
      });

      // Use new token for logout
      if (r.status === 200) {
        try {
          const body = JSON.parse(r.body);
          accessToken = body.data.token;
          refreshTokenVal = body.data.refresh_token || refreshTokenVal;
        } catch (_) { }
      }
    });
  }

  sleep(0.5);

  // ── LOGOUT ──
  group('04 Logout', function () {
    const payload = JSON.stringify({ refresh_token: refreshTokenVal });
    const r = http.post(ENDPOINTS.logout, payload, authHeaders(accessToken));
    check(r, {
      'logout: 200': (r) => r.status === 200,
    });
  });

  sleep(1);
}

export function handleSummary(data) {
  return handleSummaryFor('auth_flow', data);
}
