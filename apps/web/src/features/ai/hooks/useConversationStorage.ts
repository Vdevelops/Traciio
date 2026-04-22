import { useEffect, useCallback, useRef } from "react";

const STORAGE_KEY = "ai_chatbot_conversation";

export interface StoredMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: string; // ISO string for serialization
}

/**
 * Hook to manage conversation persistence in sessionStorage
 * Conversation is saved automatically when messages change
 * and loaded on component mount
 */
export function useConversationStorage(
  messages: Array<{ id: string; role: "user" | "assistant"; content: string; timestamp: Date }>,
  setMessages: React.Dispatch<React.SetStateAction<Array<{ id: string; role: "user" | "assistant"; content: string; timestamp: Date }>>>
) {
  const isInitialMount = useRef(true);

  // Load conversation from sessionStorage on mount (only once)
  useEffect(() => {
    if (!isInitialMount.current) {
      return;
    }

    try {
      const stored = sessionStorage.getItem(STORAGE_KEY);
      if (stored) {
        const parsed: StoredMessage[] = JSON.parse(stored);
        if (Array.isArray(parsed) && parsed.length > 0) {
          // Convert stored messages back to Message format with Date objects
          const restoredMessages = parsed.map((msg) => ({
            id: msg.id,
            role: msg.role,
            content: msg.content,
            timestamp: new Date(msg.timestamp),
          }));
          setMessages(restoredMessages);
        }
      }
    } catch (error) {
      // Clear corrupted data
      sessionStorage.removeItem(STORAGE_KEY);
    } finally {
      isInitialMount.current = false;
    }
  }, [setMessages]);

  // Save conversation to sessionStorage whenever messages change
  // Skip saving during initial mount to avoid overwriting with default state
  useEffect(() => {
    // Don't save during initial mount (before we've checked storage)
    if (isInitialMount.current) {
      return;
    }

    try {
      // Only save if there are messages
      if (messages.length > 0) {
        // Convert messages to storable format
        const messagesToStore: StoredMessage[] = messages.map((msg) => ({
          id: msg.id,
          role: msg.role,
          content: msg.content,
          timestamp: msg.timestamp.toISOString(),
        }));
        sessionStorage.setItem(STORAGE_KEY, JSON.stringify(messagesToStore));
      } else {
        // If messages array is empty, clear storage
        sessionStorage.removeItem(STORAGE_KEY);
      }
    } catch (error) {
      // Handle quota exceeded error gracefully
      if (error instanceof DOMException && error.name === "QuotaExceededError") {
        sessionStorage.removeItem(STORAGE_KEY);
      }
    }
  }, [messages]);

  // Function to clear conversation from storage
  const clearConversation = useCallback(() => {
    try {
      sessionStorage.removeItem(STORAGE_KEY);
    } catch (error) {
    }
  }, []);

  return { clearConversation };
}


