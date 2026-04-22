import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { aiService } from "../services/aiService";
import type { ChatRequest } from "../types";

export function useChat() {
  return useMutation({
    mutationFn: (data: ChatRequest) => aiService.chat(data),
    onError: (error: unknown) => {
      const message =
        error instanceof Error ? error.message : "Failed to send message";
      toast.error(message);
    },
  });
}

