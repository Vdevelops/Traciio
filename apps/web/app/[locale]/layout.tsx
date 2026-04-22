import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";
import { notFound } from "next/navigation";
import { routing } from "@/i18n/routing";
import type { Locale } from "@/features/dashboard/types";
import { ReactQueryProvider } from "@/lib/react-query";
import { LazyPermissionsProvider } from "@/features/auth/providers/lazy-permissions-provider";
import { ThemeProvider } from "@/components/providers/theme-provider";
import { AppLayout } from "@/components/layouts/app-layout";
import { Toaster } from "sonner";
import "../globals.css";

export default async function LocaleLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  if (!routing.locales.includes(locale as Locale)) {
    notFound();
  }

  const messages = await getMessages({ locale });

  return (
    <NextIntlClientProvider locale={locale} messages={messages} key={locale}>
      <ThemeProvider
        attribute="class"
        defaultTheme="light"
        enableSystem
        disableTransitionOnChange
      >
        <ReactQueryProvider>
          <LazyPermissionsProvider>
            <AppLayout>{children}</AppLayout>
          </LazyPermissionsProvider>
          <Toaster 
            position="bottom-right" 
            offset={8}
            toastOptions={{
              className: "toast-notification",
              classNames: {
                toast: "toast-base",
                success: "toast-success",
                error: "toast-error",
                info: "toast-info",
                warning: "toast-warning",
              },
            }}
          />
        </ReactQueryProvider>
      </ThemeProvider>
    </NextIntlClientProvider>
  );
}


