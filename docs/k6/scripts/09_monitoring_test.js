// =============================================================
// MONITORING TEST — Health Check Endpoints Verification
// Tests new monitoring endpoints: circuit-breakers, runtime
// Duration: ~5 min | Continuous monitoring
// =============================================================

import http from "k6/http";
import { check, sleep, group } from "k6";
import { Rate, Trend } from "k6/metrics";
import { ENDPOINTS, authHeaders } from "../lib/config.js";
import { login } from "../lib/auth.js";
import { handleSummaryFor } from "../lib/reporter.js";

const monitorLatency = new Trend("monitor_latency_ms", true);
const monitorErrors = new Rate("monitor_error_rate");

export const options = {
  stages: [
    { duration: "1m", target: 10 },
    { duration: "3m", target: 10 },
    { duration: "1m", target: 0 },
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"],
    monitor_error_rate: ["rate<0.05"],
    checks: ["rate>0.95"],
  },
};

export function setup() {
  console.log("━━━ MONITORING ENDPOINTS TEST ━━━");
  const tokens = login("admin");
  if (!tokens) return { token: null };
  return { token: tokens.accessToken };
}

export default function (data) {
  if (!data.token) return;
  const p = authHeaders(data.token);

  // Test all monitoring endpoints
  group("Health Endpoints", function () {
    const res = http.get(ENDPOINTS.health);
    track(res, "health: 200");

    // Verify response structure
    check(res, {
      "health has status": (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.status !== undefined;
        } catch (e) {
          return false;
        }
      },
    });
  });

  group("Circuit Breaker Stats", function () {
    const res = http.get(ENDPOINTS.healthCircuitBreakers);
    track(res, "circuit-breakers: 200");

    // Verify circuit breaker data structure
    check(res, {
      "circuit-breakers has timestamp": (r) => {
        try {
          const body = JSON.parse(r.body);
          return (
            body.timestamp !== undefined && body.circuit_breakers !== undefined
          );
        } catch (e) {
          return false;
        }
      },
    });

    // Log circuit breaker status untuk monitoring
    try {
      const body = JSON.parse(res.body);
      if (body.circuit_breakers) {
        Object.entries(body.circuit_breakers).forEach(([name, stats]) => {
          if (stats.state === "open") {
            console.warn(`⚠️ Circuit breaker OPEN: ${name}`);
          }
        });
      }
    } catch (e) {}
  });

  group("Runtime Stats", function () {
    const res = http.get(ENDPOINTS.healthRuntime);
    track(res, "runtime: 200");

    // Verify runtime data structure
    check(res, {
      "runtime has goroutines": (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.goroutines !== undefined && body.memory !== undefined;
        } catch (e) {
          return false;
        }
      },
    });

    // Log runtime metrics
    try {
      const body = JSON.parse(res.body);
      const goroutines = body.goroutines;
      const memAlloc = body.memory?.alloc_mb;

      if (goroutines > 5000) {
        console.warn(`⚠️ High goroutine count: ${goroutines}`);
      }
      if (memAlloc > 1000) {
        console.warn(`⚠️ High memory usage: ${memAlloc}MB`);
      }
    } catch (e) {}
  });

  group("Cache Health", function () {
    const res = http.get(ENDPOINTS.healthCache);
    track(res, "cache health: 200");
  });

  group("Cache Metrics", function () {
    const res = http.get(ENDPOINTS.healthCacheMetrics);
    track(res, "cache metrics: 200");
  });

  sleep(1);
}

function track(res, checkName) {
  monitorLatency.add(res.timings.duration);
  const ok = check(res, { [checkName]: (r) => r.status === 200 });
  monitorErrors.add(!ok);
}

export function teardown() {
  console.log("━━━ MONITORING TEST COMPLETE ━━━");
  console.log("Verify circuit breaker states dan runtime metrics di reports/");
}

export function handleSummary(data) {
  return handleSummaryFor("monitoring_test", data);
}
