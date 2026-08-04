import { cookies as nextCookies } from "next/headers";
import apiClient from "@/lib/api-client";

export async function getCurrentUser() {
  try {
    // next/headers cookies() returns a helper; call .get on resolved object
    const ck = await nextCookies();
    const tokenCookie = ck.get("token");
    const headers: Record<string, string> = {};
    if (tokenCookie) {
      headers.Authorization = `Bearer ${tokenCookie.value}`;
    }
    const resp = await apiClient.get("/users/me", { headers });
    const user = resp.data?.data?.user;

    if (!user) {
      return null;
    }

    const normalizedRole =
      typeof user.role === "string"
        ? user.role
        : user.role?.code || user.role_code || "";

    return {
      ...user,
      role: normalizedRole,
    };
  } catch (err) {
    return null;
  }
}
