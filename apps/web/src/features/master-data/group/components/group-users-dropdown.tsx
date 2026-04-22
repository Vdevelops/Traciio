"use client";

import { User as UserIcon, Mail } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { useUsers } from "@/features/master-data/user-management/hooks/useUsers";
import { Badge } from "@/components/ui/badge";
import type { User } from "@/features/master-data/user-management/types";

interface GroupUsersDropdownProps {
  readonly groupId: string;
  readonly groupName: string;
}

export function GroupUsersDropdown({
  groupId,
  groupName,
}: GroupUsersDropdownProps) {
  // Fetch users by group_id (lazy loading - only when component is rendered)
  const { data, isLoading } = useUsers({
    per_page: 100,
    group_id: groupId,
  });

  const groupUsers = data?.data ?? [];

  return (
    <div className="mt-2 border rounded-lg bg-muted/30">
      <div className="px-3 py-2 border-b bg-background/50">
        <span className="text-sm font-medium">
          Users in {groupName} ({groupUsers.length})
        </span>
      </div>
      <div className="p-3 space-y-2 max-h-[300px] overflow-y-auto">
          {isLoading ? (
            <div className="space-y-2">
              {Array.from({ length: 3 }, (_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : groupUsers.length === 0 ? (
            <div className="text-center py-4 text-sm text-muted-foreground">
              No users in this group
            </div>
          ) : (
            groupUsers.map((user: User) => (
              <div
                key={user.id}
                className="flex items-center gap-3 p-2 rounded-md hover:bg-background border cursor-pointer"
              >
                <div className="flex-shrink-0">
                  {user.avatar_url ? (
                    <img
                      src={user.avatar_url}
                      alt={user.name}
                      className="h-8 w-8 rounded-full"
                    />
                  ) : (
                    <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center">
                      <UserIcon className="h-4 w-4 text-primary" />
                    </div>
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="font-medium text-sm truncate">
                      {user.name}
                    </span>
                    {user.status === "active" ? (
                      <Badge variant="active" className="text-xs">
                        {user.status}
                      </Badge>
                    ) : (
                      <Badge variant="inactive" className="text-xs">
                        {user.status}
                      </Badge>
                    )}
                  </div>
                  <div className="flex items-center gap-1 mt-0.5">
                    <Mail className="h-3 w-3 text-muted-foreground" />
                    <span className="text-xs text-muted-foreground truncate">
                      {user.email}
                    </span>
                  </div>
                  {user.role && (
                    <span className="text-xs text-muted-foreground mt-0.5 block">
                      {user.role.name}
                    </span>
                  )}
                </div>
              </div>
            ))
          )}
      </div>
    </div>
  );
}

