"use client";

import type { ReactNode } from "react";
import Image from "next/image";

interface AuthLayoutProps {
  readonly children: ReactNode;
}

export function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="relative flex min-h-screen overflow-hidden bg-background">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(85%_45%_at_12%_12%,rgba(56,189,248,0.15),transparent_60%),radial-gradient(60%_40%_at_88%_4%,rgba(37,99,235,0.16),transparent_65%)]" />
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
          <div className="absolute inset-0 bg-gradient-to-t from-slate-950/55 via-slate-900/20 to-transparent" />
          <div className="absolute left-8 bottom-8 max-w-md rounded-2xl border border-white/25 bg-slate-950/30 p-5 text-white backdrop-blur-md">
            <p className="text-xs uppercase tracking-[0.2em] text-cyan-200/90">Tracio CRM</p>
            <h2 className="mt-2 text-2xl font-semibold leading-tight">Track customer journey with a cleaner academic showcase experience.</h2>
          </div>
        </div>
      </div>

      {/* Right Side - Form (1/3) */}
      <div className="relative z-10 flex w-full items-center justify-center p-8 lg:w-1/3">
        <div className="w-full max-w-md space-y-8 rounded-3xl border border-border/70 bg-card/85 p-7 shadow-[0_22px_44px_-30px_rgba(15,23,42,0.75)] backdrop-blur-sm">
          {/* Mobile Logo */}
          <div className="mb-8 flex items-center justify-center gap-3 lg:hidden">
            <div className="flex size-10 aspect-square items-center justify-center overflow-hidden rounded-xl border border-border/60 bg-card shadow-lg">
              <Image
                src="/tracio-logo.svg"
                alt="Tracio"
                width={40}
                height={40}
                className="object-contain"
              />
            </div>
            <div className="flex flex-col gap-0.5 leading-none">
              <span className="text-xl font-semibold text-primary">Tracio</span>
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


