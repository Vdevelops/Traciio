"use client";

import React, { memo, useEffect, useState } from "react";
import { HelpCircle, Copy, Check, Settings, LogOut } from "lucide-react";
import { NotificationBadge } from "@/features/notifications/components/notification-badge";
import { ThemeToggleButton as ThemeToggle } from "@/components/ui/theme-toggle";
import { Separator } from "@/components/ui/separator";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Link, usePathname } from "@/i18n/routing";
import { useLocale } from "next-intl";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import { useLogout } from "@/features/auth/hooks/useLogout";

export interface HeaderControlsProps {
  showNotifications?: boolean;
  showThemeToggle?: boolean;
  showLocaleToggle?: boolean;
  showProfile?: boolean;
  extraIcon?: React.ReactNode;
  showCopy?: boolean;
  copied?: boolean;
  onCopy?: () => void;
}

export const HeaderControls = memo(function HeaderControls({
  showNotifications = false,
  showThemeToggle = false,
  showLocaleToggle = false,
  showProfile = false,
  extraIcon,
  showCopy = false,
  copied = false,
  onCopy,
}: HeaderControlsProps) {
  const locale = useLocale();
  const pathname = usePathname();
  const { user } = useAuthStore();
  const logout = useLogout();

  const userName = user?.name?.split(" ")[0] || "User";
  const userAvatarUrl = user?.avatar_url || undefined;

  const [mounted, setMounted] = useState(false);
  useEffect(() => {
    const t = setTimeout(() => setMounted(true), 0);
    return () => clearTimeout(t);
  }, []);

  return (
    <div className="flex items-center gap-2">
      {showNotifications && <NotificationBadge />}

      {showCopy && (
        <Button variant="ghost" size="icon" onClick={onCopy} className="h-9 w-9">
          {copied ? <Check className="h-4 w-4 text-primary" /> : <Copy className="h-4 w-4" />}
        </Button>
      )}

      {showThemeToggle && <ThemeToggle />}

      {showLocaleToggle && (
        <Link
          href={pathname || "/dashboard"}
          locale={locale === "en" ? "id" : "en"}
          scroll={false}
        >
          <Button
            variant="ghost"
            size="icon"
            className="h-9 w-11 rounded-2xl bg-background/80 text-xs font-semibold shadow-sm hover:bg-accent/60"
          >
            {locale === "en" ? "ID" : "EN"}
          </Button>
        </Link>
      )}

      <Separator orientation="vertical" className="h-6" />

      {showProfile && (
        <div>
          {mounted ? (
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant="ghost"
                  className="flex h-9 w-9 items-center justify-center rounded-full p-0 hover:bg-muted transition-colors"
                >
                  <Avatar className="h-8 w-8">
                    <AvatarImage src={userAvatarUrl} alt={userName} />
                  </Avatar>
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-56 p-2" align="end">
                <div className="px-2 py-1.5 text-xs text-muted-foreground">
                  <div className="text-foreground text-sm font-medium">{userName}</div>
                </div>
                <div className="flex flex-col gap-1">
                  <Link
                    href="/profile"
                    className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left text-sm hover:bg-accent transition-colors"
                  >
                    <Settings className="h-4 w-4" />
                    <span className="text-sm">Settings</span>
                  </Link>
                  <button
                    type="button"
                    onClick={logout}
                    className="flex w-full items-center rounded-md px-2 py-1.5 text-left text-sm text-destructive hover:bg-destructive/10"
                  >
                    <LogOut className="h-4 w-4 mr-2" />
                    Logout
                  </button>
                </div>
              </PopoverContent>
            </Popover>
          ) : (
            <Button
              variant="ghost"
              className="flex h-9 w-9 items-center justify-center rounded-full p-0 hover:bg-muted transition-colors"
              disabled
            >
              <Avatar className="h-8 w-8">
                <AvatarImage src={userAvatarUrl} alt={userName} />
              </Avatar>
            </Button>
          )}
        </div>
      )}
    </div>
  );
});
