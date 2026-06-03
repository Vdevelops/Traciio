"use client";

import type { ReactNode } from "react";
import Image from "next/image";

interface AuthLayoutProps {
  readonly children: ReactNode;
}

export function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="page-frame relative flex min-h-screen overflow-hidden bg-background">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(80%_45%_at_12%_12%,rgba(20,184,166,0.18),transparent_62%),radial-gradient(60%_40%_at_88%_4%,rgba(217,119,6,0.14),transparent_65%)]" />
      {/* Left Side - Full Image (2/3) */}
      <div className="hidden p-6 lg:block lg:w-2/3">
        <div className="relative h-full w-full overflow-hidden rounded-[2rem] border border-border/55 shadow-[0_28px_54px_-32px_rgba(15,23,42,0.7)]">
          <Image
            src="/login.webp"
            alt="Tracio CRM Platform"
            fill
            className="object-cover"
            priority
          />
          <div className="absolute inset-0 bg-gradient-to-t from-slate-950/70 via-slate-900/18 to-transparent" />
          <div className="absolute left-8 bottom-8 max-w-md rounded-[1.5rem] border border-white/15 bg-slate-950/36 p-5 text-white backdrop-blur-md">
            <p className="page-kicker text-xs text-emerald-200/90">Tracio CRM</p>
            <h2 className="brand-text mt-2 text-2xl font-semibold leading-tight">Track customer journeys with a calm, disciplined workspace built for healthcare teams.</h2>
          </div>
        </div>
      </div>

      {/* Right Side - Form (1/3) */}
      <div className="relative z-10 flex w-full items-center justify-center p-8 lg:w-1/3">
        <div className="surface-panel-strong w-full max-w-md space-y-8 rounded-[1.75rem] p-7 backdrop-blur-sm">
          {/* Mobile Logo */}
          <div className="mb-8 flex items-center justify-center gap-3 lg:hidden">
            <div className="brand-mark flex size-10 aspect-square items-center justify-center overflow-hidden rounded-2xl border border-border/60 shadow-lg">
              <Image
                src="/tracio-logo.svg"
                alt="Tracio"
                width={40}
                height={40}
                className="object-contain"
              />
            </div>
            <div className="flex flex-col gap-0.5 leading-none">
              <span className="brand-text text-xl font-semibold text-primary">Tracio</span>
              <span className="text-xs text-muted-foreground">Track Better, Serve Smarter</span>
            </div>
          </div>

          {/* Form Content */}
          {children}
        </div>
      </div>
    </div>
  );
}


