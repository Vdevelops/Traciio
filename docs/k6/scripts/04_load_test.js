// =============================================================
// LOAD TEST — Simulates realistic CRM usage (mixed read/write)
// Duration: ~16 min | Max VUs: 30 (production safe)
// + Per-endpoint tracking to identify slow endpoints
// =============================================================

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';
import { ENDPOINTS, authHeaders } from '../lib/config.js';
import { login } from '../lib/auth.js';

const apiLatency = new Trend('api_latency_ms', true);
const apiErrors = new Rate('api_error_rate');
const requestCount = new Counter('api_requests_total');

// ── Per-Endpoint Trends (to identify which endpoint is slow) ──
const ep_dashboard_overview = new Trend('ep_dashboard_overview', true);
const ep_dashboard_pipeline = new Trend('ep_dashboard_pipeline', true);
const ep_leads_list = new Trend('ep_leads_list', true);
const ep_accounts_list = new Trend('ep_accounts_list', true);
const ep_contacts_list = new Trend('ep_contacts_list', true);
const ep_deals_list = new Trend('ep_deals_list', true);
const ep_visit_reports = new Trend('ep_visit_reports', true);
const ep_tasks_list = new Trend('ep_tasks_list', true);
const ep_schedules_list = new Trend('ep_schedules_list', true);
const ep_activities_list = new Trend('ep_activities_list', true);
const ep_products_list = new Trend('ep_products_list', true);
const ep_profile = new Trend('ep_profile', true);
const ep_reports = new Trend('ep_reports', true);
const ep_notifications = new Trend('ep_notifications', true);
const ep_pipeline_summary = new Trend('ep_pipeline_summary', true);
const ep_users_list = new Trend('ep_users_list', true);
const ep_create_lead = new Trend('ep_create_lead', true);
const ep_health = new Trend('ep_health', true);

const MAX = parseInt(__ENV.MAX_VUS || '30');

export const options = {
  stages: [
    { duration: '1m', target: Math.ceil(MAX * 0.3) },
    { duration: '3m', target: Math.ceil(MAX * 0.3) },
    { duration: '1m', target: Math.ceil(MAX * 0.7) },
    { duration: '5m', target: Math.ceil(MAX * 0.7) },
    { duration: '1m', target: MAX },
    { duration: '3m', target: MAX },
    { duration: '2m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(90)<500', 'p(95)<800', 'p(99)<1500'],
    http_req_failed: ['rate<0.05'],
    api_error_rate: ['rate<0.05'],
    checks: ['rate>0.90'],
  },
};

export function setup() {
  console.log(`━━━ LOAD TEST | Max VUs: ${MAX} | ~16 min ━━━`);
  const tokens = login('admin');
  if (!tokens) return { token: null };
  const p = authHeaders(tokens.accessToken);

  // Find a valid account ID for reports
  const r = http.get(ENDPOINTS.accounts, p);
  let firstAccountId = null;
  try {
    const data = JSON.parse(r.body).data;
    if (data && data.length > 0) firstAccountId = data[0].id;
  } catch (e) { }

  return { token: tokens.accessToken, accountId: firstAccountId };
}

export default function (data) {
  if (!data.token) return;
  const p = authHeaders(data.token);

  // Weighted random endpoint selection (realistic CRM usage pattern)
  const rand = Math.random();

  if (rand < 0.10) {
    // 10% — Dashboard Overview
    group('Dashboard Overview', function () {
      track(http.get(ENDPOINTS.dashboardOverview, p), 'dashboard overview: 200', ep_dashboard_overview);
    });
  } else if (rand < 0.16) {
    // 6% — Dashboard Pipeline
    group('Dashboard Pipeline', function () {
      track(http.get(ENDPOINTS.dashboardPipeline, p), 'dashboard pipeline: 200', ep_dashboard_pipeline);
    });
  } else if (rand < 0.28) {
    // 12% — Leads list
    group('Leads List', function () {
      const page = Math.floor(Math.random() * 3) + 1;
      track(http.get(`${ENDPOINTS.leads}?page=${page}&per_page=10`, p), 'leads: 200', ep_leads_list);
    });
  } else if (rand < 0.36) {
    // 8% — Accounts list
    group('Accounts List', function () {
      track(http.get(`${ENDPOINTS.accounts}?page=1&per_page=10`, p), 'accounts: 200', ep_accounts_list);
    });
  } else if (rand < 0.44) {
    // 8% — Contacts list
    group('Contacts List', function () {
      track(http.get(`${ENDPOINTS.contacts}?page=1&per_page=10`, p), 'contacts: 200', ep_contacts_list);
    });
  } else if (rand < 0.52) {
    // 8% — Deals list
    group('Deals List', function () {
      track(http.get(`${ENDPOINTS.deals}?page=1&per_page=10`, p), 'deals: 200', ep_deals_list);
    });
  } else if (rand < 0.58) {
    // 6% — Visit Reports list
    group('Visit Reports List', function () {
      track(http.get(`${ENDPOINTS.visitReports}?page=1&per_page=10`, p), 'visit-reports: 200', ep_visit_reports);
    });
  } else if (rand < 0.64) {
    // 6% — Tasks list
    group('Tasks List', function () {
      track(http.get(`${ENDPOINTS.tasks}?page=1&per_page=10`, p), 'tasks: 200', ep_tasks_list);
    });
  } else if (rand < 0.69) {
    // 5% — Schedules list
    group('Schedules List', function () {
      track(http.get(`${ENDPOINTS.schedules}?page=1&per_page=10`, p), 'schedules: 200', ep_schedules_list);
    });
  } else if (rand < 0.74) {
    // 5% — Activities list
    group('Activities List', function () {
      track(http.get(`${ENDPOINTS.activities}?page=1&per_page=10`, p), 'activities: 200', ep_activities_list);
    });
  } else if (rand < 0.78) {
    // 4% — Products list
    group('Products List', function () {
      track(http.get(`${ENDPOINTS.products}?page=1&per_page=10`, p), 'products: 200', ep_products_list);
    });
  } else if (rand < 0.82) {
    // 4% — User profile
    group('User Profile', function () {
      track(http.get(ENDPOINTS.myProfile, p), 'profile: 200', ep_profile);
    });
  } else if (rand < 0.86) {
    // 4% — Reports
    group('Reports', function () {
      const reports = [
        ENDPOINTS.reportSalesPerformance,
        ENDPOINTS.reportPipeline,
        ENDPOINTS.reportVisitReports,
        ENDPOINTS.reportAccountActivity,
      ];
      let url = reports[Math.floor(Math.random() * reports.length)];

      if (url === ENDPOINTS.reportAccountActivity && data.accountId) {
        url += `?account_id=${data.accountId}`;
      }

      track(http.get(url, p), 'report: 200', ep_reports);
    });
  } else if (rand < 0.89) {
    // 3% — Notifications
    group('Notifications', function () {
      track(http.get(`${ENDPOINTS.notifications}?page=1&per_page=10`, p), 'notifications: 200', ep_notifications);
    });
  } else if (rand < 0.92) {
    // 3% — Pipelines Summary
    group('Pipeline Summary', function () {
      track(http.get(ENDPOINTS.pipelineSummary, p), 'pipeline summary: 200', ep_pipeline_summary);
    });
  } else if (rand < 0.95) {
    // 3% — Users list
    group('Users List', function () {
      track(http.get(`${ENDPOINTS.users}?page=1&per_page=10`, p), 'users: 200', ep_users_list);
    });
  } else if (rand < 0.98) {
    // 3% — Create lead (write operation)
    group('Create Lead', function () {
      const payload = JSON.stringify({
        first_name: 'K6',
        last_name: `Load ${__VU}-${Date.now()}`,
        email: `k6load${__VU}${Date.now()}@test.com`,
        phone: `08${Math.floor(1000000000 + Math.random() * 9000000000)}`,
        lead_source: 'website',
        lead_status: 'new',
        notes: 'K6 load test auto-generated',
      });
      const r = http.post(ENDPOINTS.leads, payload, p);
      ep_create_lead.add(r.timings.duration);
      apiLatency.add(r.timings.duration);
      requestCount.add(1);
      const ok = check(r, {
        'create lead: 200/201': (r) => r.status === 200 || r.status === 201,
      });
      apiErrors.add(!ok);
    });
  } else {
    // 2% — Health check
    group('Health', function () {
      track(http.get(ENDPOINTS.health), 'health: 200', ep_health);
    });
  }

  sleep(Math.random() * 2 + 1); // 1-3s think time
}

function track(res, checkName, epTrend) {
  const dur = res.timings.duration;
  apiLatency.add(dur);
  epTrend.add(dur);
  requestCount.add(1);
  const ok = check(res, { [checkName]: (r) => r.status === 200 });
  apiErrors.add(!ok);

  // Log slow requests (> 2 seconds) to identify bottlenecks
  if (dur > 2000) {
    console.warn(`🐌 SLOW [${dur.toFixed(0)}ms] ${checkName} | status=${res.status}`);
  }
}

export function teardown() {
  console.log('━━━ LOAD TEST COMPLETE ━━━');
}

export function handleSummary(data) {
  const m = data.metrics || {};

  const val = (name, stat) => {
    try {
      const v = m[name]?.values?.[stat];
      return v !== undefined && v !== null ? Math.round(v * 100) / 100 : null;
    } catch (_) { return null; }
  };

  const fmt = (ms) => {
    if (ms === null || ms === undefined) return 'N/A';
    return ms < 1000 ? `${Math.round(ms)}ms` : `${(ms / 1000).toFixed(2)}s`;
  };

  // ── Collect per-endpoint metrics ──
  const endpoints = [
    'ep_dashboard_overview', 'ep_dashboard_pipeline',
    'ep_leads_list', 'ep_accounts_list', 'ep_contacts_list',
    'ep_deals_list', 'ep_visit_reports', 'ep_tasks_list',
    'ep_schedules_list', 'ep_activities_list', 'ep_products_list',
    'ep_profile', 'ep_reports', 'ep_notifications',
    'ep_pipeline_summary', 'ep_users_list', 'ep_create_lead', 'ep_health',
  ];

  const epData = [];
  for (const ep of endpoints) {
    const count = val(ep, 'count');
    if (count === null || count === 0) continue;
    epData.push({
      name: ep.replace('ep_', '').replace(/_/g, ' '),
      count: count,
      avg: val(ep, 'avg'),
      med: val(ep, 'med'),
      p90: val(ep, 'p(90)'),
      p95: val(ep, 'p(95)'),
      max: val(ep, 'max'),
    });
  }

  // Sort by p95 descending (slowest first)
  epData.sort((a, b) => (b.p95 || 0) - (a.p95 || 0));

  // ── Checks ──
  let totalPassed = 0, totalFailed = 0;
  const checks = [];
  const walk = (g) => {
    if (!g) return;
    for (const c of Object.values(g.checks || {})) {
      checks.push({ name: c.name, passes: c.passes, fails: c.fails });
      totalPassed += c.passes;
      totalFailed += c.fails;
    }
    for (const sub of Object.values(g.groups || {})) walk(sub);
  };
  walk(data.root_group);

  const httpReqs = val('http_reqs', 'count') || 0;
  const failRate = val('http_req_failed', 'rate') || 0;
  const errPct = (failRate * 100).toFixed(2) + '%';
  const totalChecks = totalPassed + totalFailed;
  const checkRate = totalChecks > 0 ? ((totalPassed / totalChecks) * 100).toFixed(1) : 'N/A';
  const rps = Math.round(val('http_reqs', 'rate') || 0);
  const avgMs = val('http_req_duration', 'avg');
  const medMs = val('http_req_duration', 'med');
  const p90Ms = val('http_req_duration', 'p(90)');
  const p95Ms = val('http_req_duration', 'p(95)');
  const p99Ms = val('http_req_duration', 'p(99)');
  const maxMs = val('http_req_duration', 'max');

  // ── Build console output ──
  const L = [];
  L.push('');
  L.push('╔══════════════════════════════════════════════════════════════════════╗');
  L.push('║  LOAD TEST REPORT — CRM Healthcare                                 ║');
  L.push('╠══════════════════════════════════════════════════════════════════════╣');
  L.push(`║  Duration   : ${fmt(data.state?.testRunDurationMs).padEnd(15)} Requests: ${String(httpReqs).padEnd(20)} ║`);
  L.push(`║  Throughput : ${(rps + ' req/s').padEnd(15)} Errors  : ${errPct.padEnd(20)} ║`);
  L.push(`║  Checks     : ✓ ${totalPassed} / ✗ ${totalFailed} (${checkRate}%)${' '.repeat(Math.max(0, 33 - String(totalPassed).length - String(totalFailed).length - checkRate.length))}║`);
  L.push('╠══════════════════════════════════════════════════════════════════════╣');
  L.push('║  RESPONSE TIMES (all endpoints)                                    ║');
  L.push(`║  avg=${fmt(avgMs).padEnd(9)} med=${fmt(medMs).padEnd(9)} p90=${fmt(p90Ms).padEnd(9)} p95=${fmt(p95Ms).padEnd(9)} ║`);
  L.push(`║  p99=${fmt(p99Ms).padEnd(9)} max=${fmt(maxMs).padEnd(9)}                             ║`);
  L.push('╠══════════════════════════════════════════════════════════════════════╣');
  L.push('║  📊 PER-ENDPOINT BREAKDOWN (sorted by p95, slowest first)          ║');
  L.push('║  ───────────────────────────────────────────────────────────────    ║');
  L.push('║  Endpoint            Count   Avg     Med     p90     p95     Max    ║');
  L.push('║  ───────────────────────────────────────────────────────────────    ║');

  for (const ep of epData) {
    const name = ep.name.length > 18 ? ep.name.substring(0, 18) : ep.name.padEnd(18);
    const isSlow = ep.p95 !== null && ep.p95 > 800;
    const marker = isSlow ? '⚠️' : '  ';
    L.push(`║ ${marker}${name} ${String(ep.count).padStart(5)}  ${fmt(ep.avg).padEnd(7)} ${fmt(ep.med).padEnd(7)} ${fmt(ep.p90).padEnd(7)} ${fmt(ep.p95).padEnd(7)} ${fmt(ep.max).padEnd(7)}║`);
  }

  L.push('║  ───────────────────────────────────────────────────────────────    ║');

  // Identify the worst offenders
  const slowEps = epData.filter(ep => ep.p95 !== null && ep.p95 > 800);
  if (slowEps.length > 0) {
    L.push('║                                                                    ║');
    L.push('║  🔴 SLOW ENDPOINTS (p95 > 800ms) — Optimization needed:            ║');
    for (const ep of slowEps) {
      L.push(`║     → ${ep.name}: p95=${fmt(ep.p95)}, max=${fmt(ep.max)}, ${ep.count} reqs${' '.repeat(Math.max(0, 28 - ep.name.length - fmt(ep.p95).length - fmt(ep.max).length - String(ep.count).length))}║`);
    }
  }

  const okEps = epData.filter(ep => ep.p95 !== null && ep.p95 <= 200);
  if (okEps.length > 0) {
    L.push('║                                                                    ║');
    L.push('║  🟢 FAST ENDPOINTS (p95 < 200ms):                                  ║');
    for (const ep of okEps.slice(0, 5)) {
      L.push(`║     ✓ ${ep.name}: p95=${fmt(ep.p95)}${' '.repeat(Math.max(0, 46 - ep.name.length - fmt(ep.p95).length))}║`);
    }
  }

  L.push('╠══════════════════════════════════════════════════════════════════════╣');
  L.push('║  📋 BUMN PRODUCTION READINESS                                      ║');

  const grade = (pass, txt) => (pass ? '✓ PASS' : '✗ FAIL') + ' ' + txt;
  L.push(`║  ${grade(avgMs !== null && avgMs < 500, `avg < 500ms (${fmt(avgMs)})`).padEnd(67)}║`);
  L.push(`║  ${grade(p95Ms !== null && p95Ms < 1000, `p95 < 1000ms (${fmt(p95Ms)})`).padEnd(67)}║`);
  L.push(`║  ${grade(failRate * 100 < 1, `errors < 1% (${errPct})`).padEnd(67)}║`);
  L.push(`║  ${grade(rps > 50, `throughput > 50 req/s (${rps})`).padEnd(67)}║`);

  L.push('╠══════════════════════════════════════════════════════════════════════╣');
  const ts = new Date().toISOString().replace(/[:.]/g, '-').substring(0, 19);
  const jsonFile = `reports/load_test_${ts}.json`;
  L.push(`║  📄 ${jsonFile.padEnd(64)}║`);
  L.push('╚══════════════════════════════════════════════════════════════════════╝');
  L.push('');

  // ── JSON report ──
  const report = {
    testName: 'load_test',
    timestamp: new Date().toISOString(),
    thresholds: Object.fromEntries(
      Object.entries(m).filter(([_, v]) => v.thresholds).map(([k, v]) => [k, v.thresholds])
    ),
    http_req_duration: { avg: avgMs, min: val('http_req_duration', 'min'), med: medMs, max: maxMs, 'p(90)': p90Ms, 'p(95)': p95Ms, 'p(99)': p99Ms },
    http_reqs: { count: httpReqs, rate: rps },
    http_req_failed: { count: Math.round(failRate * httpReqs), rate: failRate },
    checks: { total: totalChecks, passed: totalPassed, failed: totalFailed, success_rate: checkRate },
    per_endpoint: epData,
    state: data.state,
  };

  return {
    [jsonFile]: JSON.stringify(report, null, 2),
    stdout: L.join('\n'),
  };
}
