"use client";

import { useEffect } from "react";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error("Global error:", error);
  }, [error]);

  return (
    <html>
      <body>
        <div className="flex min-h-screen items-center justify-center p-6">
          <div className="w-full max-w-md space-y-4 rounded-lg border bg-card p-6 text-card-foreground shadow-sm">
            <div className="space-y-1.5">
              <h2 className="text-2xl font-semibold leading-none tracking-tight">
                Something went wrong
              </h2>
              <p className="text-sm text-muted-foreground">
                A critical error occurred. Please refresh the page.
              </p>
            </div>
            <div className="space-y-4">
              {process.env.NODE_ENV === "development" && (
                <div className="rounded-md bg-destructive/10 p-4">
                  <p className="text-sm font-mono text-destructive">
                    {error.message}
                  </p>
                  {error.digest && (
                    <p className="mt-2 text-xs text-muted-foreground">
                      Error ID: {error.digest}
                    </p>
                  )}
                </div>
              )}
              <div className="flex gap-2">
                <button
                  onClick={reset}
                  className="inline-flex flex-1 items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground shadow transition-colors hover:bg-primary/90"
                >
                  Try again
                </button>
                <button
                  onClick={() => (window.location.href = "/")}
                  className="inline-flex flex-1 items-center justify-center rounded-md border border-input bg-background px-4 py-2 text-sm font-medium shadow-sm transition-colors hover:bg-accent hover:text-accent-foreground"
                >
                  Go to Home
                </button>
              </div>
            </div>
          </div>
        </div>
      </body>
    </html>
  );
}
