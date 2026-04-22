import type { Metadata, Viewport } from "next";
import { Geist, Geist_Mono, Sora } from "next/font/google";
import { getLocale } from "next-intl/server";
import { routing } from "@/i18n/routing";
import type { Locale } from "@/features/dashboard/types";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const headingFont = Sora({
  variable: "--font-heading",
  subsets: ["latin"],
  weight: ["500", "600", "700"],
});

export const metadata: Metadata = {
  title: "Salesview Platform",
  description: "Salesview — Sales CRM platform for field teams",
};

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  maximumScale: 5,
  userScalable: true,
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  let locale: Locale;
  try {
    const localeValue = await getLocale();
    locale = routing.locales.includes(localeValue as Locale)
      ? (localeValue as Locale)
      : routing.defaultLocale;
  } catch {
    locale = routing.defaultLocale;
  }

  return (
    <html lang={locale} suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} ${headingFont.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
