/**
 * K6 HTML Report Generator
 * Usage: import { htmlReport } from '../lib/reporter.js';
 *        export function handleSummary(data) { return htmlReport(data, 'test_name'); }
 *
 * IMPORTANT: Make sure to create the reports/ directory before running tests:
 *   mkdir -p docs/k6/reports
 */

export function htmlReport(data, testName) {
  const ts = new Date().toISOString().replace(/[:.]/g, "-").slice(0, 19);
  const fname = `${testName}_${ts}`;
  return {
    [`reports/${fname}.json`]: JSON.stringify(data, null, 2),
    stdout: textSummary(data),
  };
}

export function textSummary(data) {
  const m = data.metrics || {};
  const lines = ["\n========== K6 SUMMARY ==========\n"];

  // Show checks
  const checks = m.checks;
  if (checks && checks.values) {
    const total = checks.values.passes + checks.values.fails;
    lines.push(`  Checks:`);
    lines.push(`    ✓ ${checks.values.passes} passed`);
    lines.push(`    ✗ ${checks.values.fails} failed`);
    lines.push(`    Success Rate: ${(checks.values.rate * 100).toFixed(1)}%`);
    lines.push("");
  }

  const dur = m.http_req_duration;
  if (dur && dur.values) {
    lines.push("  HTTP Request Duration:");
    lines.push(`    avg : ${fmt(dur.values.avg)}`);
    lines.push(`    min : ${fmt(dur.values.min)}`);
    lines.push(`    med : ${fmt(dur.values.med)}`);
    lines.push(`    p90 : ${fmt(dur.values["p(90)"])}`);
    lines.push(`    p95 : ${fmt(dur.values["p(95)"])}`);
    lines.push(`    p99 : ${fmt(dur.values["p(99)"])}`);
    lines.push(`    max : ${fmt(dur.values.max)}`);
  }

  const reqs = m.http_reqs;
  if (reqs && reqs.values) {
    lines.push(`\n  Total Requests : ${reqs.values.count}`);
    lines.push(`  Throughput     : ${Math.round(reqs.values.rate)} req/s`);
  }

  const fail = m.http_req_failed;
  if (fail && fail.values) {
    lines.push(`  HTTP Errors    : ${(fail.values.rate * 100).toFixed(2)}%`);
  }

  // Show custom metrics
  for (const [key, metric] of Object.entries(m)) {
    if (
      key.startsWith("api_") ||
      key.startsWith("stress_") ||
      key.startsWith("spike_") ||
      key.startsWith("soak_") ||
      key.startsWith("login_") ||
      key.startsWith("profile_") ||
      key.startsWith("refresh_") ||
      key.startsWith("auth_") ||
      key.startsWith("monitor_")
    ) {
      if (metric.type === "trend" && metric.values) {
        lines.push(
          `  ${key}: avg=${fmt(metric.values.avg)} p95=${fmt(metric.values["p(95)"])} max=${fmt(metric.values.max)}`,
        );
      } else if (metric.type === "rate" && metric.values) {
        lines.push(`  ${key}: ${(metric.values.rate * 100).toFixed(2)}%`);
      } else if (metric.type === "counter" && metric.values) {
        lines.push(`  ${key}: ${metric.values.count}`);
      }
    }
  }

  lines.push("\n================================\n");
  return lines.join("\n");
}

function fmt(ms) {
  if (ms === undefined || ms === null) return "N/A";
  return ms < 1000 ? `${Math.round(ms)}ms` : `${(ms / 1000).toFixed(2)}s`;
}

/**
 * Standard report generator for all K6 scripts.
 * Usage:
 *   import { handleSummaryFor } from '../lib/reporter.js';
 *   export function handleSummary(data) { return handleSummaryFor('test_name', data); }
 */
export function handleSummaryFor(testName, data) {
  const m = data.metrics || {};

  const val = (name, stat) => {
    try {
      const v = m[name]?.values?.[stat];
      return v !== undefined && v !== null ? Math.round(v * 100) / 100 : "N/A";
    } catch (_) {
      return "N/A";
    }
  };

  // Collect check results
  let totalPassed = 0;
  let totalFailed = 0;
  const checks = [];

  const walk = (group) => {
    if (!group) return;
    for (const c of Object.values(group.checks || {})) {
      checks.push({ name: c.name, passes: c.passes, fails: c.fails });
      totalPassed += c.passes;
      totalFailed += c.fails;
    }
    for (const g of Object.values(group.groups || {})) {
      walk(g);
    }
  };
  walk(data.root_group);

  const httpReqs = m.http_reqs?.values?.count || 0;
  const failedReqs = m.http_req_failed?.values?.passes || 0;
  const successRate =
    httpReqs > 0
      ? (((httpReqs - failedReqs) / httpReqs) * 100).toFixed(2)
      : "N/A";

  const summaryReport = {
    testName,
    timestamp: new Date().toISOString(),
    thresholds: Object.fromEntries(
      Object.entries(m)
        .filter(([_, v]) => v.thresholds)
        .map(([k, v]) => [k, v.thresholds]),
    ),
    http_req_duration: {
      avg: val("http_req_duration", "avg"),
      min: val("http_req_duration", "min"),
      med: val("http_req_duration", "med"),
      max: val("http_req_duration", "max"),
      "p(90)": val("http_req_duration", "p(90)"),
      "p(95)": val("http_req_duration", "p(95)"),
      "p(99)": val("http_req_duration", "p(99)"),
    },
    http_reqs: { count: httpReqs },
    http_req_failed: {
      count: failedReqs,
      rate: m.http_req_failed?.values?.rate || 0,
    },
    checks: {
      total: totalPassed + totalFailed,
      passed: totalPassed,
      failed: totalFailed,
      success_rate:
        totalPassed + totalFailed > 0
          ? ((totalPassed / (totalPassed + totalFailed)) * 100).toFixed(2)
          : 0,
    },
    state: data.state,
  };

  const ts = new Date().toISOString().replace(/[:.]/g, "-").slice(0, 19);
  const fname = `${testName}_${ts}`;

  return {
    [`reports/${fname}.json`]: JSON.stringify(summaryReport, null, 2),
    stdout: textSummary(data),
  };
}
