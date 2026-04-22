"use client";

import { Drawer } from "@/components/ui/drawer";
import { NotificationList } from "./notification-list";
import { useTranslations } from "next-intl";

interface NotificationDrawerProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
}

export function NotificationDrawer({ open, onOpenChange }: NotificationDrawerProps) {
  const t = useTranslations("notifications");

  return (
    <Drawer
      open={open}
      onOpenChange={onOpenChange}
      title={t("drawerTitle")}
      description={t("drawerDescription")}
      side="right"
      defaultWidth={480}
      minWidth={320}
      maxWidth={800}
    >
      <NotificationList />
    </Drawer>
  );
}

