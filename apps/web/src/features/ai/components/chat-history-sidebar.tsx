"use client";

import { useState, useMemo, memo } from "react";
import {
  Search,
  Plus,
  Trash2,
  MessageSquare,
  Settings,
  MoreHorizontal,
  X,
  Trash,
} from "lucide-react";
import { useChatHistoryStore, type ChatConversation } from "../stores/useChatHistoryStore";
import { AISettingsDrawer } from "./ai-settings-drawer";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";

interface ChatHistorySidebarProps {
  onNewChat: () => void;
}

const ConversationItem = memo(function ConversationItem({
  conversation,
  isActive,
  onSelect,
  onDelete,
}: {
  conversation: ChatConversation;
  isActive: boolean;
  onSelect: () => void;
  onDelete: () => void;
}) {
  return (
    <div
      role="button"
      tabIndex={0}
      onClick={onSelect}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          onSelect();
        }
      }}
      className={cn(
        "group relative w-full text-left flex items-center gap-3 px-3 py-2.5 rounded-lg cursor-pointer transition-colors",
        isActive
          ? "bg-primary/10 text-foreground"
          : "hover:bg-muted/50 text-muted-foreground hover:text-foreground"
      )}
    >
      <MessageSquare className="h-4 w-4 shrink-0" />
      <span className="flex-1 truncate text-sm">{conversation.title}</span>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            className="h-6 w-6 opacity-0 group-hover:opacity-100 transition-opacity"
            onClick={(e) => e.stopPropagation()}
          >
            <MoreHorizontal className="h-3.5 w-3.5" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem
            onClick={(e) => {
              e.stopPropagation();
              onDelete();
            }}
            className="text-destructive focus:text-destructive"
          >
            <Trash2 className="h-4 w-4 mr-2" />
            Delete
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
});

const ConversationGroup = memo(function ConversationGroup({
  label,
  conversations,
  activeId,
  onSelect,
  onDelete,
}: {
  label: string;
  conversations: ChatConversation[];
  activeId: string | null;
  onSelect: (id: string) => void;
  onDelete: (id: string) => void;
}) {
  if (conversations.length === 0) return null;

  return (
    <div className="space-y-1">
      <div className="px-3 py-2">
        <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
          {label}
        </span>
      </div>
      {conversations.map((conv) => (
        <ConversationItem
          key={conv.id}
          conversation={conv}
          isActive={conv.id === activeId}
          onSelect={() => onSelect(conv.id)}
          onDelete={() => onDelete(conv.id)}
        />
      ))}
    </div>
  );
});

export const ChatHistorySidebar = memo(function ChatHistorySidebar({
  onNewChat,
}: ChatHistorySidebarProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [conversationToDelete, setConversationToDelete] = useState<string | null>(null);
  const [clearAllDialogOpen, setClearAllDialogOpen] = useState(false);
  const [settingsOpen, setSettingsOpen] = useState(false);
  const hasAISettingsPermission = useHasPermission("ai-settings.view");
  const hasAISettingsEditPermission = useHasPermission("ai-settings.edit");

  const {
    conversations,
    activeConversationId,
    isSidebarOpen,
    setActiveConversation,
    deleteConversation,
    getConversationsByDate,
    toggleSidebar,
    clearAllHistory,
    searchConversations,
  } = useChatHistoryStore();

  const filteredConversations = useMemo(() => {
    if (!searchQuery.trim()) {
      return getConversationsByDate();
    }
    const results = searchConversations(searchQuery);
    // Group filtered results by date too
    const today: ChatConversation[] = [];
    const yesterday: ChatConversation[] = [];
    const lastWeek: ChatConversation[] = [];
    const older: ChatConversation[] = [];

    const now = new Date();
    const yesterdayDate = new Date(now);
    yesterdayDate.setDate(yesterdayDate.getDate() - 1);
    const weekAgo = new Date(now);
    weekAgo.setDate(weekAgo.getDate() - 7);

    results.forEach((conv) => {
      const date = new Date(conv.updatedAt);
      const isToday =
        date.getDate() === now.getDate() &&
        date.getMonth() === now.getMonth() &&
        date.getFullYear() === now.getFullYear();
      const isYesterday =
        date.getDate() === yesterdayDate.getDate() &&
        date.getMonth() === yesterdayDate.getMonth() &&
        date.getFullYear() === yesterdayDate.getFullYear();
      const isLastWeek = date >= weekAgo && !isToday && !isYesterday;

      if (isToday) today.push(conv);
      else if (isYesterday) yesterday.push(conv);
      else if (isLastWeek) lastWeek.push(conv);
      else older.push(conv);
    });

    return { today, yesterday, lastWeek, older };
  }, [searchQuery, getConversationsByDate, searchConversations, conversations]);

  const handleDeleteConversation = (id: string) => {
    setConversationToDelete(id);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (conversationToDelete) {
      deleteConversation(conversationToDelete);
      setConversationToDelete(null);
    }
    setDeleteDialogOpen(false);
  };



  const confirmClearAll = () => {
    clearAllHistory();
    setClearAllDialogOpen(false);
  };

  if (!isSidebarOpen) {
    return (
      <div className="flex flex-col items-center py-4 px-2 border-r border-border bg-card h-full">
        <Button
          variant="ghost"
          size="icon"
          onClick={toggleSidebar}
          className="mb-4"
        >
          <MessageSquare className="h-5 w-5" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          onClick={onNewChat}
          className="mb-2"
        >
          <Plus className="h-5 w-5" />
        </Button>
      </div>
    );
  }

  return (
    <div className="flex flex-col w-64 h-full border-r border-border bg-card">
      {/* Header with search */}
      <div className="p-3 space-y-3">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search chats..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9 h-9 bg-muted/50 border-0 focus-visible:ring-1"
          />
          {searchQuery && (
            <Button
              variant="ghost"
              size="icon"
              className="absolute right-1 top-1/2 -translate-y-1/2 h-6 w-6"
              onClick={() => setSearchQuery("")}
            >
              <X className="h-3.5 w-3.5" />
            </Button>
          )}
        </div>
      </div>

      {/* Conversations list */}
      <ScrollArea className="flex-1 px-2">
        <div className="space-y-4 pb-4">
          <ConversationGroup
            label="Today"
            conversations={filteredConversations.today}
            activeId={activeConversationId}
            onSelect={setActiveConversation}
            onDelete={handleDeleteConversation}
          />
          <ConversationGroup
            label="Yesterday"
            conversations={filteredConversations.yesterday}
            activeId={activeConversationId}
            onSelect={setActiveConversation}
            onDelete={handleDeleteConversation}
          />
          <ConversationGroup
            label="7 Days Ago"
            conversations={filteredConversations.lastWeek}
            activeId={activeConversationId}
            onSelect={setActiveConversation}
            onDelete={handleDeleteConversation}
          />
          <ConversationGroup
            label="Older"
            conversations={filteredConversations.older}
            activeId={activeConversationId}
            onSelect={setActiveConversation}
            onDelete={handleDeleteConversation}
          />
          {conversations.length === 0 && (
            <div className="px-3 py-8 text-center text-muted-foreground text-sm">
              No conversations yet
            </div>
          )}
        </div>
      </ScrollArea>

      {/* Footer navigation - Settings and Clear All */}
      <div className="border-t border-border p-2 space-y-1">
        {hasAISettingsPermission && (
          <button 
            onClick={() => setSettingsOpen(true)}
            className="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-muted-foreground hover:bg-muted/50 hover:text-foreground transition-colors"
          >
            <Settings className="h-4 w-4" />
            <span>Settings</span>
          </button>
        )}
        {conversations.length > 0 && (
          <button 
            onClick={() => setClearAllDialogOpen(true)}
            className="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-destructive/80 hover:bg-destructive/10 hover:text-destructive transition-colors"
          >
            <Trash className="h-4 w-4" />
            <span>Clear All History</span>
          </button>
        )}
      </div>

      {/* New Chat Button */}
      <div className="p-3 border-t border-border">
        <Button
          onClick={onNewChat}
          className="w-full gap-2 bg-primary text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="h-4 w-4" />
          New Chat
        </Button>
      </div>

      {/* AI Settings Drawer */}
      {hasAISettingsPermission && (
        <AISettingsDrawer 
          open={settingsOpen} 
          onOpenChange={setSettingsOpen}
          readOnly={!hasAISettingsEditPermission}
        />
      )}

      {/* Delete confirmation dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Conversation</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this conversation? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteDialogOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={confirmDelete}>
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Clear all confirmation dialog */}
      <Dialog open={clearAllDialogOpen} onOpenChange={setClearAllDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Clear All History</DialogTitle>
            <DialogDescription>
              Are you sure you want to clear all chat history? This will delete all conversations and cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setClearAllDialogOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={confirmClearAll}>
              Clear All
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
});
