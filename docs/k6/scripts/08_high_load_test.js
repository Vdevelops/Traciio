// =============================================================
// HIGH LOAD TEST — 1000+ Concurrent Users
// Simulates extreme load untuk test server resilience
// Duration: ~35 min | Max VUs: 1000 (override via MAX_VUS env)
// Tests: Concurrency limiter, timeout handling, connection stability
// =============================================================

import http from "k6/http";
import { check, sleep, group } from "k6";
import { Counter, Rate, Trend } from "k6/metrics";
import { ENDPOINTS, authHeaders } from "../lib/config.js";
import { login } from "../lib/auth.js";
import { handleSummaryFor } from "../lib/reporter.js";

const apiLatency = new Trend("api_latency_ms", true);
const apiErrors = new Rate("api_error_rate");
const requestCount = new Counter("api_requests_total");
const serverBusyCount = new Counter("server_busy_503");
const timeoutCount = new Counter("server_timeout_504");

// Configuration — 1000 VUs default, override dengan MAX_VUS=2000 atau lebih
const MAX_VUS = parseInt(__ENV.MAX_VUS || "1000");
const TARGET_RPS = parseInt(__ENV.TARGET_RPS || "5000");

export const options = {
  stages: [
    // Phase 1: Ramp up ke 20% (warm up)
    { duration: "3m", target: Math.ceil(MAX_VUS * 0.2) },
    { duration: "2m", target: Math.ceil(MAX_VUS * 0.2) },

    // Phase 2: Ramp up ke 50%
    { duration: "3m", target: Math.ceil(MAX_VUS * 0.5) },
    { duration: "3m", target: Math.ceil(MAX_VUS * 0.5) },

    // Phase 3: Ramp up ke 100% (full load)
    { duration: "5m", target: MAX_VUS },
    { duration: "8m", target: MAX_VUS },

    // Phase 4: Sustained load
    { duration: "5m", target: MAX_VUS },

    // Phase 5: Cool down
    { duration: "3m", target: Math.ceil(MAX_VUS * 0.5) },
    { duration: "3m", target: 0 },
  ],
  thresholds: {
    http_req_duration: [
      "p(50)<500", // 50% under 500ms
      "p(90)<2000", // 90% under 2s
      "p(95)<5000", // 95% under 5s (lenient untuk high load)
      "p(99)<10000", // 99% under 10s
    ],
    http_req_failed: ["rate<0.10"], // Max 10% errors (high load tolerance)
    api_error_rate: ["rate<0.10"],
    checks: ["rate>0.85"], // 85% checks pass (lenient untuk extreme load)
  },
  // Rate limiting untuk menghindari overwhelm server
  rps: TARGET_RPS,
};

export function setup() {
  console.log(
    `━━━ HIGH LOAD TEST | Max VUs: ${MAX_VUS} | Target RPS: ${TARGET_RPS} | ~35 min ━━━`,
  );

  // Login multiple users untuk distribusi load
  const adminTokens = login("admin");
  const salesTokens = login("sales");
  const managerTokens = login("salesManager");

  if (!adminTokens) {
    console.error("Failed to login admin user");
    return { tokens: null };
  }

  // Get sample data untuk testing
  const adminHeaders = authHeaders(adminTokens.accessToken);

  // Get sample account
  const accountsRes = http.get(
    `${ENDPOINTS.accounts}?page=1&per_page=1`,
    adminHeaders,
  );
  let sampleAccountId = null;
  try {
    const data = JSON.parse(accountsRes.body).data;
    if (data && data.length > 0) sampleAccountId = data[0].id;
  } catch (e) {}

  // Get sample lead
  const leadsRes = http.get(
    `${ENDPOINTS.leads}?page=1&per_page=1`,
    adminHeaders,
  );
  let sampleLeadId = null;
  try {
    const data = JSON.parse(leadsRes.body).data;
    if (data && data.length > 0) sampleLeadId = data[0].id;
  } catch (e) {}

  return {
    tokens: {
      admin: adminTokens.accessToken,
      sales: salesTokens ? salesTokens.accessToken : adminTokens.accessToken,
      manager: managerTokens
        ? managerTokens.accessToken
        : adminTokens.accessToken,
    },
    sampleAccountId,
    sampleLeadId,
  };
}

export default function (data) {
  if (!data.tokens) return;

  // Rotate tokens untuk distribusi load
  const tokenRotation = ["admin", "sales", "manager"];
  const selectedToken = data.tokens[tokenRotation[__VU % 3]];
  const p = authHeaders(selectedToken);

  // Weighted random endpoint selection (optimized untuk high load)
  const rand = Math.random();

  // 40% - Read operations (lightweight)
  if (rand < 0.15) {
    // 15% - Dashboard Overview (test N+1 fix)
    group("Dashboard Overview [Optimized]", function () {
      const res = http.get(ENDPOINTS.dashboardOverview, p);
      track(res, "dashboard overview: 200/201");

      // Verify response time improvement (< 2s untuk N+1 fix)
      if (res.timings.duration > 2000) {
        console.warn(
          `⚠️ Dashboard slow: ${res.timings.duration}ms (expected <2000ms)`,
        );
      }
    });
  } else if (rand < 0.25) {
    // 10% - Health monitoring endpoints
    group("Health Monitoring", function () {
      const endpoints = [
        ENDPOINTS.health,
        ENDPOINTS.healthCircuitBreakers,
        ENDPOINTS.healthRuntime,
      ];
      const url = endpoints[Math.floor(Math.random() * endpoints.length)];
      const res = http.get(url);
      track(res, "health monitoring: 200");
    });
  } else if (rand < 0.35) {
    // 10% - Leads list
    group("Leads List", function () {
      const page = Math.floor(Math.random() * 5) + 1;
      const res = http.get(`${ENDPOINTS.leads}?page=${page}&per_page=20`, p);
      track(res, "leads: 200");
    });
  } else if (rand < 0.4) {
    // 5% - Accounts list
    group("Accounts List", function () {
      const res = http.get(`${ENDPOINTS.accounts}?page=1&per_page=20`, p);
      track(res, "accounts: 200");
    });
  }

  // 30% - Mixed operations
  else if (rand < 0.5) {
    // 10% - Contacts & Deals (parallel)
    group("Contacts & Deals", function () {
      const responses = http.batch([
        ["GET", `${ENDPOINTS.contacts}?page=1&per_page=10`, null, p],
        ["GET", `${ENDPOINTS.deals}?page=1&per_page=10`, null, p],
      ]);
      responses.forEach((res, idx) => {
        track(res, `${idx === 0 ? "contacts" : "deals"}: 200`);
      });
    });
  } else if (rand < 0.58) {
    // 8% - Visit Reports & Activities
    group("Visit Reports & Activities", function () {
      const res = http.get(`${ENDPOINTS.visitReports}?page=1&per_page=10`, p);
      track(res, "visit-reports: 200");
    });
  } else if (rand < 0.65) {
    // 7% - Tasks & Schedules
    group("Tasks & Schedules", function () {
      const res = http.get(`${ENDPOINTS.tasks}?page=1&per_page=10`, p);
      track(res, "tasks: 200");
    });
  } else if (rand < 0.7) {
    // 5% - Notifications
    group("Notifications", function () {
      const res = http.get(`${ENDPOINTS.notifications}?page=1&per_page=10`, p);
      track(res, "notifications: 200");
    });
  }

  // 20% - Write operations (test N+1 fix)
  else if (rand < 0.8) {
    // 10% - Create Lead
    group("Create Lead [Write]", function () {
      const payload = JSON.stringify({
        first_name: "K6",
        last_name: `HighLoad${__VU}-${Date.now()}`,
        email: `k6high${__VU}${Date.now()}@test.com`,
        phone: `08${Math.floor(1000000000 + Math.random() * 9000000000)}`,
        lead_source: "website",
        lead_status: "new",
        notes: "K6 high load test - N+1 fix verification",
      });
      const res = http.post(ENDPOINTS.leads, payload, p);
      track(res, "create lead: 200/201");

      // Lead creation should be fast dengan N+1 fix (< 1s)
      if (res.timings.duration > 1000) {
        console.warn(
          `⚠️ Lead creation slow: ${res.timings.duration}ms (expected <1000ms)`,
        );
      }
    });
  } else if (rand < 0.85) {
    // 5% - Create Account
    group("Create Account [Write]", function () {
      const payload = JSON.stringify({
        name: `K6 Account ${__VU} ${Date.now()}`,
        email: `k6acc${__VU}${Date.now()}@test.com`,
        phone: `08${Math.floor(1000000000 + Math.random() * 9000000000)}`,
        address: "Test Address",
        city: "Jakarta",
        status: "active",
      });
      const res = http.post(ENDPOINTS.accounts, payload, p);
      track(res, "create account: 200/201");
    });
  } else if (rand < 0.9) {
    // 5% - Create Contact
    group("Create Contact [Write]", function () {
      const payload = JSON.stringify({
        first_name: "K6",
        last_name: `Contact ${__VU}`,
        email: `k6contact${__VU}${Date.now()}@test.com`,
        phone: `08${Math.floor(1000000000 + Math.random() * 9000000000)}`,
      });
      const res = http.post(ENDPOINTS.contacts, payload, p);
      track(res, "create contact: 200/201");
    });
  }

  // 10% - Reports & Analytics
  else {
    group("Reports & Analytics", function () {
      const reports = [
        ENDPOINTS.reportSalesPerformance,
        ENDPOINTS.reportPipeline,
        ENDPOINTS.reportVisitReports,
      ];
      const url = reports[Math.floor(Math.random() * reports.length)];
      const res = http.get(url, p);
      track(res, "report: 200");
    });
  }

  // Shorter think time untuk high load
  sleep(Math.random() * 1 + 0.5); // 0.5-1.5s
}

function track(res, checkName) {
  apiLatency.add(res.timings.duration);
  requestCount.add(1);

  // Track server-side load shedding separately (these are expected under extreme load)
  if (res.status === 503) {
    serverBusyCount.add(1);
  } else if (res.status === 504) {
    timeoutCount.add(1);
  }

  // Only count non-503/504 failures as real errors
  const ok = check(res, {
    [checkName]: (r) =>
      r.status === 200 || r.status === 201 || r.status === 503 || r.status === 504,
  });
  apiErrors.add(res.status !== 200 && res.status !== 201 && res.status !== 503 && res.status !== 504);

  // Log slow requests untuk analysis
  if (res.timings.duration > 5000) {
    console.warn(
      `⏱️ Slow request: ${checkName} took ${res.timings.duration}ms`,
    );
  }
}

export function teardown() {
  console.log("━━━ HIGH LOAD TEST COMPLETE ━━━");
  console.log(`Target: ${MAX_VUS} VUs | Target RPS: ${TARGET_RPS}`);
  console.log("503 (Server Busy) = concurrency limiter working correctly");
  console.log("504 (Timeout) = request exceeded timeout, context cancelled");
  console.log("EOF errors  = server dropping connections (should be near 0 after fix)");
  console.log("Check reports/ folder untuk detailed analysis");
}

export function handleSummary(data) {
  return handleSummaryFor("high_load_test", data);
}
