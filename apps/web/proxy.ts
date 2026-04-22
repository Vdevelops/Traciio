import createMiddleware from "next-intl/middleware";
import { routing } from "./src/i18n/routing";
import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const intlMiddleware = createMiddleware(routing);

export default function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const token = request.cookies.get("token")?.value;

  // Public routes that don't require authentication (locale-agnostic patterns)
  const publicRoutePatterns = [
    /^\/[a-z]{2}\/login$/, // /en/login, /id/login, etc.
    /^\/[a-z]{2}$/, // Root locale pages like /en, /id (handled by index page)
    /^\/$/, // Root path
  ];

  const isPublicRoute = publicRoutePatterns.some((pattern) =>
    pattern.test(pathname)
  );

  // If accessing root and already authenticated, redirect to default-locale dashboard
  if (pathname === "/" && token) {
    const target = "/en/dashboard"; // default locale is "en"
    return NextResponse.redirect(new URL(target, request.url));
  }

  // If accessing login page and already authenticated, redirect to dashboard
  if (pathname.match(/^\/[a-z]{2}\/login$/) && token) {
    // Extract locale from pathname (e.g., /en/login -> en)
    const locale = pathname.split("/")[1];
    const target = `/${locale}/dashboard`;
    return NextResponse.redirect(new URL(target, request.url));
  }

  // If accessing protected routes without token, redirect to login
  // This prevents 404 loop when token is missing after server restart
  if (!isPublicRoute && !token) {
    // Extract locale from pathname or default to "en"
    const segments = pathname.split("/").filter(Boolean);
    const locale = segments[0]?.match(/^[a-z]{2}$/) ? segments[0] : "en";
    const target = `/${locale}/login`;
    return NextResponse.redirect(new URL(target, request.url));
  }

  // Run next-intl middleware for locale handling
  return intlMiddleware(request);
}

export const config = {
  // Match all pathnames except for
  // - … if they start with `/api`, `/_next` or `/_vercel`
  // - … the ones containing a dot (e.g. `favicon.ico`)
  matcher: ["/((?!api|_next|_vercel|.*\\..*).*)"],
};

