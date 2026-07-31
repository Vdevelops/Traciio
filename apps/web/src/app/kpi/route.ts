import { NextResponse } from "next/server";
import { routing } from "@/i18n/routing";

// Redirects non-prefixed /kpi to the default locale (/en/kpi)
export function GET(request: Request) {
  try {
    const defaultLocale = routing.defaultLocale || "en";
    const url = new URL(request.url);
    url.pathname = `/${defaultLocale}/kpi`;
    return NextResponse.redirect(url);
  } catch (err) {
    // Fallback: redirect to /en/kpi
    return NextResponse.redirect(new URL("/en/kpi", request.url));
  }
}
