"use client";

import { Bell } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useNotificationCount } from "../hooks/useNotifications";
import { useTranslations } from "next-intl";
import { useNotificationStore } from "../stores/useNotificationStore";

export function NotificationBadge() {
  const t = useTranslations("notifications");
  const { data: unreadCount = 0, isLoading } = useNotificationCount();
  const { openDrawer } = useNotificationStore();

  return (
    <Button
      variant="ghost"
      size="icon"
      className="relative"
      onClick={() => openDrawer()}
      aria-label={t("badgeAriaLabel")}
    >
      <Bell className="h-5 w-5" />
      {!isLoading && unreadCount > 0 && (
        <Badge
          variant="destructive"
          className="absolute -right-1 -top-1 flex h-5 w-5 items-center justify-center rounded-full p-0 text-xs"
        >
          {unreadCount > 99 ? "99+" : unreadCount}
        </Badge>
      )}
    </Button>
  );
}

