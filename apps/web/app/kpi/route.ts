import { NextResponse } from "next/server";

export function GET(request: Request) {
  const url = new URL(request.url);
  const locale = process.env.NEXT_PUBLIC_DEFAULT_LOCALE || "en";

  url.pathname = `/${locale}/kpi`;
  url.search = "";

  return NextResponse.redirect(url);
}
