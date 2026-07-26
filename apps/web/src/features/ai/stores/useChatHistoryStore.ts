import { create } from "zustand";
import { persist, type PersistStorage, type StorageValue } from "zustand/middleware";

export interface ChatMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: string;
}

export interface ChatConversation {
  id: string;
  title: string;
  messages: ChatMessage[];
  createdAt: string;
  updatedAt: string;
  model?: string;
}

interface ChatHistoryState {
  conversations: ChatConversation[];
  activeConversationId: string | null;
  isSidebarOpen: boolean;
  storageKey: string;
  
  // Actions
  createConversation: (model?: string) => string;
  deleteConversation: (id: string) => void;
  updateConversation: (id: string, updates: Partial<Omit<ChatConversation, "id">>) => void;
  setActiveConversation: (id: string | null) => void;
  addMessage: (conversationId: string, message: Omit<ChatMessage, "id" | "timestamp">) => void;
  getActiveConversation: () => ChatConversation | undefined;
  getConversationsByDate: () => { today: ChatConversation[]; yesterday: ChatConversation[]; lastWeek: ChatConversation[]; older: ChatConversation[] };
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  clearAllHistory: () => void;
  searchConversations: (query: string) => ChatConversation[];
}

type ChatHistoryPersistedState = Pick<
  ChatHistoryState,
  "conversations" | "activeConversationId" | "isSidebarOpen" | "storageKey"
>;

const STORAGE_KEY_PREFIX = "ai-chat-history";
const ANONYMOUS_STORAGE_KEY = `${STORAGE_KEY_PREFIX}:anonymous`;

function getAuthenticatedUserId(): string | null {
  if (typeof window === "undefined") {
    return null;
  }

  try {
    const rawAuthStorage = localStorage.getItem("auth-storage");
    if (!rawAuthStorage) {
      return null;
    }

    const parsed = JSON.parse(rawAuthStorage) as {
      state?: { user?: { id?: unknown } | null };
    };
    const userId = parsed.state?.user?.id;
    return typeof userId === "string" && userId.trim().length > 0 ? userId : null;
  } catch {
    return null;
  }
}

function getChatHistoryStorageKey(): string {
  const userId = getAuthenticatedUserId();
  return userId ? `${STORAGE_KEY_PREFIX}:user:${userId}` : ANONYMOUS_STORAGE_KEY;
}

function isPersistedChatHistoryState(value: unknown): value is ChatHistoryPersistedState {
  if (!value || typeof value !== "object") {
    return false;
  }

  const state = value as Partial<ChatHistoryPersistedState>;
  return (
    Array.isArray(state.conversations) &&
    (typeof state.activeConversationId === "string" || state.activeConversationId === null) &&
    typeof state.isSidebarOpen === "boolean"
  );
}

const chatHistoryStorage: PersistStorage<ChatHistoryPersistedState> = {
  getItem: (): StorageValue<ChatHistoryPersistedState> | null => {
    if (typeof window === "undefined") {
      return null;
    }

    const stored = localStorage.getItem(getChatHistoryStorageKey());
    if (!stored) {
      return null;
    }

    try {
      return JSON.parse(stored) as StorageValue<ChatHistoryPersistedState>;
    } catch {
      localStorage.removeItem(getChatHistoryStorageKey());
      return null;
    }
  },
  setItem: (_name: string, value: StorageValue<ChatHistoryPersistedState>): void => {
    if (typeof window === "undefined") {
      return;
    }

    localStorage.setItem(getChatHistoryStorageKey(), JSON.stringify(value));
  },
  removeItem: (): void => {
    if (typeof window === "undefined") {
      return;
    }

    localStorage.removeItem(getChatHistoryStorageKey());
  },
};

// Generate unique ID for messages and conversations
let idCounter = 1;
function generateId(): string {
  return `${Date.now()}-${idCounter++}`;
}

// Generate title from first user message
function generateTitle(messages: ChatMessage[]): string {
  const firstUserMessage = messages.find((m) => m.role === "user");
  if (firstUserMessage) {
    const content = firstUserMessage.content.trim();
    if (content.length <= 50) {
      return content;
    }
    return content.substring(0, 47) + "...";
  }
  return "New Chat";
}

// Check if date is today
function isToday(date: Date): boolean {
  const today = new Date();
  return (
    date.getDate() === today.getDate() &&
    date.getMonth() === today.getMonth() &&
    date.getFullYear() === today.getFullYear()
  );
}

// Check if date is yesterday
function isYesterday(date: Date): boolean {
  const yesterday = new Date();
  yesterday.setDate(yesterday.getDate() - 1);
  return (
    date.getDate() === yesterday.getDate() &&
    date.getMonth() === yesterday.getMonth() &&
    date.getFullYear() === yesterday.getFullYear()
  );
}

// Check if date is within last 7 days
function isLastWeek(date: Date): boolean {
  const now = new Date();
  const weekAgo = new Date();
  weekAgo.setDate(weekAgo.getDate() - 7);
  return date >= weekAgo && date < now && !isToday(date) && !isYesterday(date);
}

export const useChatHistoryStore = create<ChatHistoryState>()(
  persist(
    (set, get) => ({
      conversations: [],
      activeConversationId: null,
      isSidebarOpen: true,
      storageKey: getChatHistoryStorageKey(),

      createConversation: (model?: string) => {
        const id = generateId();
        const now = new Date().toISOString();
        const newConversation: ChatConversation = {
          id,
          title: "New Chat",
          messages: [],
          createdAt: now,
          updatedAt: now,
          model,
        };
        set((state) => ({
          conversations: [newConversation, ...state.conversations],
          activeConversationId: id,
        }));
        return id;
      },

      deleteConversation: (id: string) => {
        set((state) => {
          const newConversations = state.conversations.filter((c) => c.id !== id);
          const newActiveId =
            state.activeConversationId === id
              ? newConversations[0]?.id ?? null
              : state.activeConversationId;
          return {
            conversations: newConversations,
            activeConversationId: newActiveId,
          };
        });
      },

      updateConversation: (id: string, updates: Partial<Omit<ChatConversation, "id">>) => {
        set((state) => ({
          conversations: state.conversations.map((c) =>
            c.id === id
              ? { ...c, ...updates, updatedAt: new Date().toISOString() }
              : c
          ),
        }));
      },

      setActiveConversation: (id: string | null) => {
        set({ activeConversationId: id });
      },

      addMessage: (conversationId: string, message: Omit<ChatMessage, "id" | "timestamp">) => {
        const messageId = generateId();
        const timestamp = new Date().toISOString();
        const fullMessage: ChatMessage = {
          ...message,
          id: messageId,
          timestamp,
        };

        set((state) => {
          const updatedConversations = state.conversations.map((c) => {
            if (c.id === conversationId) {
              const newMessages = [...c.messages, fullMessage];
              return {
                ...c,
                messages: newMessages,
                title: c.messages.length === 0 ? generateTitle([fullMessage]) : c.title,
                updatedAt: timestamp,
              };
            }
            return c;
          });
          return { conversations: updatedConversations };
        });
      },

      getActiveConversation: () => {
        const state = get();
        return state.conversations.find((c) => c.id === state.activeConversationId);
      },

      getConversationsByDate: () => {
        const { conversations } = get();
        const today: ChatConversation[] = [];
        const yesterday: ChatConversation[] = [];
        const lastWeek: ChatConversation[] = [];
        const older: ChatConversation[] = [];

        conversations.forEach((conv) => {
          const date = new Date(conv.updatedAt);
          if (isToday(date)) {
            today.push(conv);
          } else if (isYesterday(date)) {
            yesterday.push(conv);
          } else if (isLastWeek(date)) {
            lastWeek.push(conv);
          } else {
            older.push(conv);
          }
        });

        return { today, yesterday, lastWeek, older };
      },

      toggleSidebar: () => {
        set((state) => ({ isSidebarOpen: !state.isSidebarOpen }));
      },

      setSidebarOpen: (open: boolean) => {
        set({ isSidebarOpen: open });
      },

      clearAllHistory: () => {
        set({ conversations: [], activeConversationId: null });
      },

      searchConversations: (query: string) => {
        const { conversations } = get();
        if (!query.trim()) {
          return conversations;
        }
        const lowerQuery = query.toLowerCase();
        return conversations.filter(
          (c) =>
            c.title.toLowerCase().includes(lowerQuery) ||
            c.messages.some((m) => m.content.toLowerCase().includes(lowerQuery))
        );
      },
    }),
    {
      name: STORAGE_KEY_PREFIX,
      storage: chatHistoryStorage,
      partialize: (state) => ({
        conversations: state.conversations,
        activeConversationId: state.activeConversationId,
        isSidebarOpen: state.isSidebarOpen,
        storageKey: getChatHistoryStorageKey(),
      }),
      merge: (persistedState, currentState) => {
        const storageKey = getChatHistoryStorageKey();
        if (!isPersistedChatHistoryState(persistedState)) {
          return {
            ...currentState,
            conversations: [],
            activeConversationId: null,
            isSidebarOpen: true,
            storageKey,
          };
        }

        const hasActiveConversation = persistedState.conversations.some(
          (conversation) => conversation.id === persistedState.activeConversationId
        );

        return {
          ...currentState,
          ...persistedState,
          activeConversationId: hasActiveConversation ? persistedState.activeConversationId : null,
          storageKey,
        };
      },
    }
  )
);
