import http from 'k6/http';
import { ENDPOINTS, USERS, jsonHeaders } from './config.js';

/**
 * Login dan return token.
 * Response format dari authservice: { success: true, data: { token, refresh_token, user } }
 */
export function login(role = 'admin') {
  const user = USERS[role];
  if (!user) {
    console.error(`[LOGIN] Unknown role: ${role}`);
    return null;
  }

  const payload = JSON.stringify({
    email: user.email,
    password: user.password,
  });

  const res = http.post(ENDPOINTS.login, payload, jsonHeaders());

  // Debug
  console.log(`[LOGIN] ${role} → ${res.status} | ${res.timings.duration.toFixed(0)}ms`);

  if (res.status === 0) {
    console.error('[LOGIN] ❌ Server unreachable (timeout/connection refused)');
    return null;
  }

  if (res.status === 429) {
    console.error('[LOGIN] ❌ Rate limited (429). Tunggu beberapa menit.');
    return null;
  }

  if (res.status !== 200) {
    console.error(`[LOGIN] ❌ Status ${res.status}: ${res.body ? res.body.substring(0, 200) : 'empty'}`);
    return null;
  }

  let body;
  try {
    body = JSON.parse(res.body);
  } catch (e) {
    console.error(`[LOGIN] ❌ Response bukan JSON: ${res.body ? res.body.substring(0, 200) : 'empty'}`);
    return null;
  }

  // Extract token — format: data.token & data.refresh_token
  const token = body.data?.token || body.data?.access_token || body.token || null;
  const refreshToken = body.data?.refresh_token || body.refresh_token || null;

  if (!token) {
    console.error(`[LOGIN] ❌ Token tidak ditemukan di response.`);
    console.error(`[LOGIN] Body keys: ${JSON.stringify(Object.keys(body))}`);
    if (body.data) console.error(`[LOGIN] Data keys: ${JSON.stringify(Object.keys(body.data))}`);
    return null;
  }

  console.log(`[LOGIN] ✅ ${role} authenticated`);
  return { accessToken: token, refreshToken };
}
