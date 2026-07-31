import { cookies as nextCookies } from "next/headers";
import apiClient from "@/lib/api-client";

export async function getCurrentUser() {
  try {
    // next/headers cookies() returns a helper; call .get on resolved object
    const ck = await nextCookies();
    const tokenCookie = ck.get("token");
    const cookieHeader = tokenCookie ? `token=${tokenCookie.value}` : undefined;
    const resp = await apiClient.get("/auth/me", { headers: cookieHeader ? { cookie: cookieHeader } : undefined });
    return resp.data?.data || null;
  } catch (err) {
    return null;
  }
}
