"use client";

import { Bell, CheckCheck, Trash2, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  useNotifications,
  useMarkAllAsRead,
  useDeleteNotification,
  useMarkAsRead,
} from "../hooks/useNotifications";
import type { Notification } from "../types";
import { formatDistanceToNow } from "date-fns";
import { useTranslations } from "next-intl";
import { useState } from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export function NotificationList() {
  const t = useTranslations("notifications");
  const [filter, setFilter] = useState<"all" | "unread" | "read">("all");
  const [page, setPage] = useState(1);

  const isReadFilter = filter === "read" ? true : filter === "unread" ? false : undefined;

  const {
    data: response,
    isLoading,
    isError,
    error,
    refetch,
  } = useNotifications({
    page,
    per_page: 20,
    is_read: isReadFilter,
  });

  const markAsRead = useMarkAsRead();
  const markAllAsRead = useMarkAllAsRead();
  const deleteNotification = useDeleteNotification();

  const notifications = response?.data ?? [];
  const pagination = response?.meta?.pagination;

  const handleMarkAsRead = async (id: string) => {
    await markAsRead.mutateAsync(id);
  };

  const handleMarkAllAsRead = async () => {
    await markAllAsRead.mutateAsync();
  };

  const handleDelete = async (id: string) => {
    await deleteNotification.mutateAsync(id);
  };

  if (isLoading) {
    return (
      <div className="space-y-3">
        {[...Array(5)].map((_, i) => (
          <Card key={i}>
            <CardContent className="p-4">
              <Skeleton className="h-4 w-3/4 mb-2" />
              <Skeleton className="h-3 w-1/2" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <p className="text-sm text-destructive mb-4">
            {error instanceof Error ? error.message : t("errorLoading")}
          </p>
          <Button onClick={() => refetch()} variant="outline" size="sm">
            {t("retry")}
          </Button>
        </CardContent>
      </Card>
    );
  }

  if (notifications.length === 0) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <div className="rounded-full bg-muted p-3 mb-3 inline-block">
            <Bell className="h-5 w-5 text-muted-foreground" />
          </div>
          <p className="text-sm text-muted-foreground">{t("empty")}</p>
        </CardContent>
      </Card>
    );
  }

  const unreadCount = notifications.filter((n) => !n.is_read).length;

  return (
    <div className="space-y-4">
      {/* Filters and Actions */}
      <div className="flex items-center justify-between gap-4">
        <Select value={filter} onValueChange={(value) => {
          setFilter(value as "all" | "unread" | "read");
          setPage(1);
        }}>
          <SelectTrigger className="w-[180px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t("filterAll")}</SelectItem>
            <SelectItem value="unread">{t("filterUnread")}</SelectItem>
            <SelectItem value="read">{t("filterRead")}</SelectItem>
          </SelectContent>
        </Select>

        {unreadCount > 0 && (
          <Button
            variant="outline"
            size="sm"
            onClick={handleMarkAllAsRead}
            disabled={markAllAsRead.isPending}
            className="gap-2"
          >
            <CheckCheck className="h-4 w-4" />
            {markAllAsRead.isPending ? t("markingAllAsRead") : t("markAllAsRead")}
          </Button>
        )}
      </div>

      {/* Notification List */}
      <div className="space-y-3">
        {notifications.map((notification) => (
          <NotificationItem
            key={notification.id}
            notification={notification}
            onMarkAsRead={handleMarkAsRead}
            onDelete={handleDelete}
          />
        ))}
      </div>

      {/* Pagination */}
      {pagination && pagination.total_pages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            {t("showing", {
              from: (pagination.page - 1) * pagination.per_page + 1,
              to: Math.min(pagination.page * pagination.per_page, pagination.total),
              total: pagination.total,
            })}
          </p>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={!pagination.has_prev || page === 1}
            >
              {t("previous")}
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPage((p) => p + 1)}
              disabled={!pagination.has_next}
            >
              {t("next")}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}

interface NotificationItemProps {
  readonly notification: Notification;
  readonly onMarkAsRead: (id: string) => void;
  readonly onDelete: (id: string) => void;
}

function NotificationItem({ notification, onMarkAsRead, onDelete }: NotificationItemProps) {
  const t = useTranslations("notifications");
  const markAsRead = useMarkAsRead();
  const deleteNotification = useDeleteNotification();

  const timeAgo = formatDistanceToNow(new Date(notification.created_at), {
    addSuffix: true,
  });

  const handleMarkAsRead = () => {
    if (!notification.is_read) {
      onMarkAsRead(notification.id);
    }
  };

  const handleDelete = () => {
    onDelete(notification.id);
  };

  return (
    <Card
      className={`transition-all hover:shadow-sm ${
        !notification.is_read ? "border-primary/20 bg-primary/5" : ""
      }`}
    >
      <CardContent className="px4">
        <div className="flex items-start justify-between gap-4">
          <div className="flex-1 space-y-1 cursor-pointer" onClick={handleMarkAsRead}>
            <div className="flex items-start gap-2">
              {!notification.is_read && (
                <span className="mt-1.5 h-2 w-2 rounded-full bg-primary shrink-0" />
              )}
              <div className="flex-1">
                <h4 className={`text-sm font-medium ${!notification.is_read ? "font-medium" : ""}`}>
                  {notification.title}
                </h4>
                {notification.message && (
                  <p className="text-sm text-muted-foreground mt-1">{notification.message}</p>
                )}
                <p className="text-xs text-muted-foreground mt-2">{timeAgo}</p>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-1">
            {!notification.is_read && (
              <Button
                variant="ghost"
                size="icon-sm"
                onClick={handleMarkAsRead}
                disabled={markAsRead.isPending}
                title={t("markAsRead")}
                className="h-8 w-8"
              >
                {markAsRead.isPending ? (
                  <Loader2 className="h-3.5 w-3.5 animate-spin" />
                ) : (
                  <CheckCheck className="h-3.5 w-3.5" />
                )}
              </Button>
            )}
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={handleDelete}
              disabled={deleteNotification.isPending}
              title={t("delete")}
              className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
            >
              {deleteNotification.isPending ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin" />
              ) : (
                <Trash2 className="h-3.5 w-3.5" />
              )}
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

