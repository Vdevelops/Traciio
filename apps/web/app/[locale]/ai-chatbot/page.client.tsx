"use client";

import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { ChatbotRedesigned } from "@/features/ai/components/chatbot";

function AIChatbotPageContent() {
  return (
    <div className="h-screen w-full overflow-hidden relative">
      <ChatbotRedesigned />
    </div>
  );
}

export default function AIChatbotPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="ai-chatbot.view">
        <AIChatbotPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
