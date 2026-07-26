"use client";

import { useState, useRef, useEffect, useMemo, useCallback, memo } from "react";
import { Send, Loader2, Globe, Code, Palette, Search, PanelLeftClose, PanelLeft } from "lucide-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import type { Components } from "react-markdown";
import { useChat } from "../hooks/useChat";
import { useAISettings } from "../hooks/useAISettings";
import { useDomainDetection } from "../hooks/useDomainDetection";
import { useChatHistoryStore } from "../stores/useChatHistoryStore";
import { ChatHistorySidebar } from "./chat-history-sidebar";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import type { AIDomain, AIActionEntity } from "../types";
import { parseActionCards, AIActionCards } from "./ai-action-cards";
import { parseAICharts, AICharts } from "./ai-chart";
import { parseLocationNeeded, LocationShareCard } from "./location-share-card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import dynamic from "next/dynamic";
import { AccountDetailModal } from "@/features/sales-crm/account-management/components/account-detail-modal";
import { ContactDetailModal } from "@/features/sales-crm/account-management/components/contact-detail-modal";
import { TaskDetailModal } from "@/features/sales-crm/task-management/components/task-detail-modal";
import { VisitReportDetailModal } from "@/features/sales-crm/visit-report/components/visit-report-detail-modal";
import { DealDetailModal } from "@/features/sales-crm/pipeline-management/components/deal-detail-modal";
import { LeadDetailModal } from "@/features/sales-crm/lead-management/components/lead-detail-modal";

// Get greeting based on time of day
function getGreeting(): string {
  const hour = new Date().getHours();
  if (hour < 12) return "Good Morning";
  if (hour < 18) return "Good Afternoon";
  return "Good Evening";
}

// Extract custom links from markdown for proper href handling
function extractCustomLinks(markdown: string): Map<string, string> {
  const linkMap = new Map<string, string>();
  // Match [text](type://uuid) format where UUID can contain hyphens
  const customLinkRegex = /\[([^\]]+)\]\((\w+:\/\/[a-f0-9-]+)\)/gi;
  let match;
  
  while ((match = customLinkRegex.exec(markdown)) !== null) {
    const linkText = match[1];
    const linkHref = match[2];
    linkMap.set(linkText, linkHref);
  }
  
  return linkMap;
}

// Quick action buttons data - mapped to CRM domains
const quickActions = [
  { id: "sales", label: "Sales", icon: Globe },
  { id: "customers", label: "Customers", icon: Search },
  { id: "analytics", label: "Analytics", icon: Code },
  { id: "route", label: "Route", icon: Palette },
];

// Gradient orb component
const GradientOrb = memo(function GradientOrb() {
  return (
    <div className="relative w-32 h-32 mx-auto mb-8">
      <div className="absolute inset-0 rounded-full bg-linear-to-br from-primary/30 via-orange-300/40 to-pink-300/30 blur-xl animate-pulse" />
      <div className="absolute inset-2 rounded-full bg-linear-to-br from-primary/40 via-orange-200/50 to-pink-200/40 blur-lg" />
      <div className="absolute inset-4 rounded-full bg-linear-to-br from-primary/50 via-orange-100/60 to-pink-100/50" />
    </div>
  );
});

export function ChatbotRedesigned() {
  const [input, setInput] = useState("");
  const [copied, setCopied] = useState(false);
  const { mutate: sendMessage, isPending } = useChat();
  const { settings } = useAISettings();
  const { user } = useAuthStore();
  const { detectDomain, domainOptions } = useDomainDetection(settings.data_privacy);
  const [selectedDomain, setSelectedDomain] = useState<AIDomain>("general");
  const scrollRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLTextAreaElement>(null);
  
  // Chat history store
  const {
    activeConversationId,
    isSidebarOpen,
    createConversation,
    addMessage,
    getActiveConversation,
    toggleSidebar,
  } = useChatHistoryStore();
  
  const activeConversation = getActiveConversation();
  const messages = useMemo(() => activeConversation?.messages ?? [], [activeConversation?.messages]);
  const userId = user?.id ?? null;

  useEffect(() => {
    void useChatHistoryStore.persist.rehydrate();
  }, [userId]);

  // Load header controls dynamically to keep chatbot bundle small
  const HeaderControls = dynamic(
    () => import("@/components/ui/header-controls").then((m) => m.HeaderControls),
    { ssr: false }
  );

  // Model selection state
  const [userSelectedModel, setUserSelectedModel] = useState<string | null>(null);
  const selectedModel = userSelectedModel || "gpt-oss-120b";

  // State for detail modals
  const [viewingAccountId, setViewingAccountId] = useState<string | null>(null);
  const [viewingContactId, setViewingContactId] = useState<string | null>(null);
  const [viewingTaskId, setViewingTaskId] = useState<string | null>(null);
  const [viewingVisitReportId, setViewingVisitReportId] = useState<string | null>(null);
  const [viewingDealId, setViewingDealId] = useState<string | null>(null);
  const [viewingLeadId, setViewingLeadId] = useState<string | null>(null);

  const userName = user?.name?.split(" ")[0] || "User";
  const userAvatarUrl = user?.avatar_url || null;
  const greeting = getGreeting();

  // Shared handler for opening entity detail modals (used by links and action cards)
  const openEntityDetail = useCallback((entity: string, entityId: string) => {
    setViewingAccountId(null);
    setViewingContactId(null);
    setViewingTaskId(null);
    setViewingVisitReportId(null);
    setViewingDealId(null);
    setViewingLeadId(null);
    
    switch (entity) {
      case 'account': 
        setViewingAccountId(entityId); 
        break;
      case 'contact': 
        setViewingContactId(entityId); 
        break;
      case 'task': 
        setViewingTaskId(entityId); 
        break;
      case 'visit': 
        setViewingVisitReportId(entityId); 
        break;
      case 'deal': 
        setViewingDealId(entityId); 
        break;
      case 'lead': 
        setViewingLeadId(entityId); 
        break;
      default:
        break;
    }
  }, []);

  // Handle custom link clicks
  const handleLinkClick = useCallback((_e: React.MouseEvent<HTMLAnchorElement> | MouseEvent, href: string) => {
    const customLinkRegex = /^(\w+):\/\/([a-f0-9-]+)$/i;
    const match = customLinkRegex.exec(href);
    
    if (!match) {
      window.open(href, '_blank', 'noopener,noreferrer');
      return;
    }
    
    const [, type, id] = match;
    
    if (id.trim().length === 0) return;
    
    openEntityDetail(type as AIActionEntity, id);
  }, [openEntityDetail]);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.style.height = "auto";
      inputRef.current.style.height = `${Math.min(inputRef.current.scrollHeight, 120)}px`;
    }
  }, [input]);

  const handleNewChat = useCallback(() => {
    createConversation(selectedModel);
  }, [createConversation, selectedModel]);

  // Sends a message programmatically (e.g. from the LocationShareCard)
  const sendMessageDirectly = useCallback((text: string) => {
    if (!text.trim() || isPending || !settings.enabled) return;

    let conversationId = activeConversationId;
    if (!conversationId) {
      conversationId = createConversation(selectedModel);
    }

    addMessage(conversationId, { role: "user", content: text });

    const conversationHistory = messages.map((msg) => ({
      role: msg.role,
      content: msg.content,
    }));

    const detected = detectDomain(text);
    const effectiveDomain = selectedDomain === "general" ? detected.domain : selectedDomain;

    sendMessage(
      {
        message: text,
        conversation_history: conversationHistory.length > 0 ? conversationHistory : undefined,
        model: selectedModel,
        domain: effectiveDomain,
      },
      {
        onSuccess: (response) => {
          if (conversationId) {
            addMessage(conversationId, { role: "assistant", content: response.data.message });
          }
        },
        onError: () => {
          if (conversationId) {
            addMessage(conversationId, { role: "assistant", content: "Sorry, I encountered an error. Please try again." });
          }
        },
      }
    );
  }, [isPending, settings.enabled, activeConversationId, createConversation, selectedModel, addMessage, messages, sendMessage, detectDomain, selectedDomain]);

  const handleSend = useCallback(() => {
    if (!input.trim() || isPending || !settings.enabled) return;

    let conversationId = activeConversationId;
    
    // Create new conversation if none active
    if (!conversationId) {
      conversationId = createConversation(selectedModel);
    }

    // Add user message
    addMessage(conversationId, {
      role: "user",
      content: input,
    });

    const currentInput = input;
    setInput("");

    // Prepare conversation history for API
    const conversationHistory = messages.map((msg) => ({
      role: msg.role,
      content: msg.content,
    }));

    // Auto-detect domain from user input, or use manually selected domain
    const detected = detectDomain(currentInput);
    const effectiveDomain = selectedDomain === "general" ? detected.domain : selectedDomain;

    sendMessage(
      {
        message: currentInput,
        conversation_history: conversationHistory.length > 0 ? conversationHistory : undefined,
        model: selectedModel,
        domain: effectiveDomain,
      },
      {
        onSuccess: (response) => {
          if (conversationId) {
            addMessage(conversationId, {
              role: "assistant",
              content: response.data.message,
            });
          }
        },
        onError: () => {
          if (conversationId) {
            addMessage(conversationId, {
              role: "assistant",
              content: "Sorry, I encountered an error. Please try again.",
            });
          }
        },
      }
    );
  }, [input, isPending, settings.enabled, activeConversationId, createConversation, selectedModel, addMessage, messages, sendMessage, detectDomain, selectedDomain]);

  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInput(e.target.value);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const handleQuickAction = (actionId: string) => {
    const prompts: Record<string, string> = {
      sales: "Tampilkan ringkasan data sales: total leads, deals di pipeline, dan tasks pending",
      customers: "Tampilkan semua akun yang ada beserta jumlah kontak terkait",
      analytics: "Berikan analisis performa sales bulan ini termasuk conversion rate dan revenue",
      route: "Buatkan rute optimal untuk kunjungan sales hari ini berdasarkan jadwal yang ada",
    };
    // Auto-set domain based on quick action
    const domainMap: Record<string, AIDomain> = {
      sales: "sales",
      customers: "customers",
      analytics: "analytics",
      route: "route_optimization",
    };
    if (domainMap[actionId]) {
      setSelectedDomain(domainMap[actionId]);
    }
    setInput(prompts[actionId] || "");
    inputRef.current?.focus();
  };

  const handleCopyChat = () => {
    const chatText = messages
      .map((msg) => `${msg.role === "user" ? "User" : "Assistant"}: ${msg.content}`)
      .join("\n\n");
    navigator.clipboard.writeText(chatText).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  };

  // Custom markdown components - defined inline to properly capture handleLinkClick
  const markdownComponents: Components = {
    h1: ({ children, ...props }) => <h1 className="text-lg font-medium mb-2 mt-4 first:mt-0" {...props}>{children}</h1>,
    h2: ({ children, ...props }) => <h2 className="text-base font-medium mb-2 mt-3 first:mt-0" {...props}>{children}</h2>,
    h3: ({ children, ...props }) => <h3 className="text-sm font-medium mb-1 mt-2 first:mt-0" {...props}>{children}</h3>,
    p: ({ children, ...props }) => <p className="text-sm leading-relaxed mb-2 last:mb-0" {...props}>{children}</p>,
    ul: ({ children, ...props }) => <ul className="list-disc list-outside space-y-1.5 my-3 ml-6" {...props}>{children}</ul>,
    ol: ({ children, ...props }) => <ol className="list-decimal list-outside space-y-1.5 my-3 ml-6" {...props}>{children}</ol>,
    li: ({ children, ...props }) => <li className="text-sm leading-relaxed pl-1" {...props}>{children}</li>,
    strong: ({ children, ...props }) => <strong className="font-medium" {...props}>{children}</strong>,
    em: ({ children, ...props }) => <em className="italic" {...props}>{children}</em>,
    // Link component will be overridden per message to use extracted custom links
    a: ({ children, href, ...props }) => (
      <a href={href || "#"} className="text-primary underline hover:no-underline" {...props}>
        {children}
      </a>
    ),
    img: ({ src, alt, ...props }) => {
      const imageSrc = typeof src === "string" ? src : "";
      const isExternalChart =
        imageSrc.includes("quickchart.io") ||
        imageSrc.includes("image-charts.com") ||
        imageSrc.includes("chart.googleapis.com");

      if (isExternalChart) {
        return (
          <div className="my-3 rounded-md border border-border bg-muted/40 px-3 py-2 text-sm text-muted-foreground">
            Grafik eksternal tidak ditampilkan karena format chart dari AI tidak valid. Minta ulang grafik agar ditampilkan sebagai tabel atau bar chart teks.
          </div>
        );
      }

      return (
        // eslint-disable-next-line @next/next/no-img-element
        <img
          src={imageSrc}
          alt={alt ?? ""}
          className="my-3 max-w-full rounded-md border border-border"
          {...props}
        />
      );
    },
    code: ({ children, className, ...props }) => {
      const isInline = !className?.includes('language-');
      return isInline ? (
        <code className="text-xs bg-muted/50 px-1.5 py-0.5 rounded font-mono" {...props}>{children}</code>
      ) : (
        <code className="block text-xs bg-muted/50 p-3 rounded font-mono overflow-x-auto" {...props}>{children}</code>
      );
    },
    table: ({ children, ...props }) => (
      <div className="overflow-x-auto my-4 -mx-2">
        <table className="w-full text-sm border-collapse border border-border" {...props}>{children}</table>
      </div>
    ),
    thead: ({ children, ...props }) => <thead className="bg-muted/50" {...props}>{children}</thead>,
    tbody: ({ children, ...props }) => <tbody {...props}>{children}</tbody>,
    tr: ({ children, ...props }) => <tr className="hover:bg-muted/30 transition-colors" {...props}>{children}</tr>,
    th: ({ children, ...props }) => <th className="px-4 py-2.5 text-left font-medium text-xs bg-muted/50 border border-border" {...props}>{children}</th>,
    td: ({ children, ...props }) => <td className="px-4 py-2.5 text-sm border border-border" {...props}>{children}</td>,
    blockquote: ({ children, ...props }) => <blockquote className="border-l-2 border-muted-foreground/30 pl-3 italic my-2" {...props}>{children}</blockquote>,
  };

  if (!settings.enabled) {
    return (
      <div className="flex items-center justify-center h-full w-full">
        <p className="text-muted-foreground text-sm">
          AI Assistant is disabled. Please enable it in settings.
        </p>
      </div>
    );
  }

  const isEmptyState = messages.length === 0;

  return (
    <div className="flex h-full w-full bg-background">
      {/* Chat History Sidebar */}
      <ChatHistorySidebar onNewChat={handleNewChat} />

      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col min-w-0 relative">
        {/* Toggle sidebar button (when sidebar is closed) */}
        <div className="absolute top-4 left-4 z-10">
          <Button
            variant="ghost"
            size="icon"
            onClick={toggleSidebar}
            className="h-9 w-9"
          >
            {isSidebarOpen ? (
              <PanelLeftClose className="h-5 w-5" />
            ) : (
              <PanelLeft className="h-5 w-5" />
            )}
          </Button>
        </div>

        <div className="absolute top-4 right-4 z-10">
          <HeaderControls
            showCopy={messages.length > 0}
            copied={copied}
            onCopy={handleCopyChat}
            showThemeToggle
            showProfile
          />
        </div>

        {isEmptyState ? (
          /* Empty State - Welcome Screen */
          <div className="flex-1 flex flex-col items-center justify-center px-4 py-8">
            <div className="max-w-2xl w-full text-center">
              {/* Gradient Orb */}
              <GradientOrb />

              {/* Greeting */}
              <h1 className="text-3xl sm:text-4xl font-medium text-foreground mb-2">
                {greeting}, {userName}
              </h1>
              <p className="text-2xl sm:text-3xl mb-12">
                <span className="text-foreground">How Can I </span>
                <span className="bg-linear-to-r from-primary via-orange-400 to-pink-400 bg-clip-text text-transparent">
                  Assist You Today?
                </span>
              </p>

              {/* Input Area */}
              <div className="w-full max-w-xl mx-auto space-y-4">
                {/* Input box */}
                <div className="relative bg-card rounded-2xl border border-border shadow-lg">
                  <textarea
                    ref={inputRef}
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyDown={handleKeyDown}
                    placeholder="Ask me anything..."
                    className="w-full resize-none bg-transparent px-4 pt-4 pb-14 text-sm focus:outline-none min-h-14 max-h-40 rounded-2xl"
                    disabled={isPending || !settings.enabled}
                    rows={1}
                  />
                  
                  {/* Bottom bar with icons */}
                  <div className="absolute bottom-0 left-0 right-0 flex items-center justify-between px-3 py-2.5 ">
                    <div className="flex items-center gap-1">
                      <Select
                        value={selectedModel}
                        onValueChange={setUserSelectedModel}
                        disabled={isPending}
                      >
                        <SelectTrigger className="h-8 w-auto min-w-[120px] border-0 bg-transparent hover:bg-muted/50 px-2 text-xs gap-1.5">
                          <Globe className="h-3.5 w-3.5" />
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="gpt-oss-120b">GPT-OSS-120B</SelectItem>
                        </SelectContent>
                      </Select>
                      <Select
                        value={selectedDomain}
                        onValueChange={(v) => setSelectedDomain(v as AIDomain)}
                        disabled={isPending}
                      >
                        <SelectTrigger className="h-8 w-auto min-w-[90px] border-0 bg-transparent hover:bg-muted/50 px-2 text-xs gap-1.5">
                          <SelectValue placeholder="Domain" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="general">Auto-detect</SelectItem>
                          {domainOptions.map((d) => (
                            <SelectItem key={d.id} value={d.id}>
                              {d.label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex items-center gap-1">
                      <Button
                        onClick={handleSend}
                        disabled={!input.trim() || isPending || !settings.enabled}
                        size="icon"
                        className="h-8 w-8 rounded-full bg-foreground text-background hover:bg-foreground/90 disabled:opacity-50"
                      >
                        {isPending ? (
                          <Loader2 className="h-4 w-4 animate-spin" />
                        ) : (
                          <Send className="h-4 w-4" />
                        )}
                      </Button>
                    </div>
                  </div>
                </div>

                {/* Quick Action Buttons */}
                <div className="flex items-center justify-center gap-3 flex-wrap">
                  {quickActions.map((action) => (
                    <Button
                      key={action.id}
                      variant="outline"
                      size="sm"
                      onClick={() => handleQuickAction(action.id)}
                      className="rounded-full gap-2 px-4 border-border/50 hover:bg-muted/50"
                    >
                      <action.icon className="h-4 w-4" />
                      {action.label}
                    </Button>
                  ))}
                </div>
              </div>
            </div>
          </div>
        ) : (
          /* Chat Messages View */
          <>
            <div
              ref={scrollRef}
              className="flex-1 overflow-y-auto pb-[200px]"
              style={{ scrollBehavior: 'smooth' }}
            >
              <div className="flex flex-col w-full">
                {messages.map((message) => (
                  <div
                    key={message.id}
                    className={`group w-full ${message.role === "user" ? "bg-muted/10" : "bg-background"}`}
                  >
                    <div className="flex items-start gap-5 px-4 sm:px-6 py-6 max-w-4xl mx-auto">
                      {/* Avatar */}
                      <div className="shrink-0 w-8 h-8">
                        {message.role === "user" ? (
                          userAvatarUrl && (
                            <Avatar className="w-8 h-8 ring-2 ring-border">
                              <AvatarImage src={userAvatarUrl} alt={userName} />
                            </Avatar>
                          )
                        ) : (
                          <div className="w-8 h-8 rounded-full bg-linear-to-br from-primary/20 to-primary/5 flex items-center justify-center ring-2 ring-primary/10">
                            <svg className="w-4 h-4 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                            </svg>
                          </div>
                        )}
                      </div>
                      
                      {/* Message Content */}
                      <div className="flex-1 min-w-0">
                        {message.role === "assistant" ? (
                          <div className="prose prose-sm dark:prose-invert max-w-none">
                            {(() => {
                              // Parse action cards and location-needed marker from the raw message
                              const { cleanMessage: afterActions, actions } = parseActionCards(message.content);
                              const { cleanMessage: afterCharts, charts } = parseAICharts(afterActions);
                              const { cleanMessage, needsLocation } = parseLocationNeeded(afterCharts);
                              const customLinks = extractCustomLinks(cleanMessage);
                              
                              const componentsWithLinks: Components = {
                                ...markdownComponents,
                                a: ({ children, href, node, ...props }) => {
                                  const linkText = typeof children === 'string' 
                                    ? children 
                                    : Array.isArray(children) 
                                      ? children.map(c => typeof c === 'string' ? c : '').join('')
                                      : '';
                                  
                                  const customHref = customLinks.get(linkText);
                                  const actualHref = customHref || href || (node as { properties?: { href?: string } })?.properties?.href || "";
                                  
                                  const customLinkRegex = /^\w+:\/\/[a-f0-9-]+$/i;
                                  const isCustomLink = actualHref ? customLinkRegex.test(actualHref) : false;
                                  
                                  if (isCustomLink && actualHref) {
                                    return (
                                      <button
                                        type="button"
                                        className="text-primary underline hover:no-underline cursor-pointer bg-transparent border-none p-0 font-inherit inline text-left"
                                        onClick={(e) => {
                                          e.preventDefault();
                                          e.stopPropagation();
                                          handleLinkClick(e as unknown as React.MouseEvent<HTMLAnchorElement>, actualHref);
                                        }}
                                        onMouseDown={(e) => {
                                          e.preventDefault();
                                          e.stopPropagation();
                                        }}
                                        title={actualHref}
                                        data-custom-link={actualHref}
                                      >
                                        {children}
                                      </button>
                                    );
                                  }
                                  
                                  return (
                                    <a 
                                      href={actualHref || "#"} 
                                      className="text-primary underline hover:no-underline cursor-pointer" 
                                      onClick={(e) => {
                                        e.preventDefault();
                                        e.stopPropagation();
                                        if (actualHref && (actualHref.startsWith('http') || !actualHref.startsWith('#'))) {
                                          window.open(actualHref, '_blank', 'noopener,noreferrer');
                                        }
                                      }}
                                      onMouseDown={(e) => {
                                        e.preventDefault();
                                        e.stopPropagation();
                                      }}
                                      {...props}
                                    >
                                      {children}
                                    </a>
                                  );
                                },
                              };
                              
                              return (
                                <>
                                  <ReactMarkdown remarkPlugins={[remarkGfm]} components={componentsWithLinks}>
                                    {cleanMessage}
                                  </ReactMarkdown>
                                  <AICharts charts={charts} />
                                  <AIActionCards
                                    actions={actions}
                                    onEntityClick={openEntityDetail}
                                  />
                                  {needsLocation && (
                                    <LocationShareCard
                                      onLocationShared={sendMessageDirectly}
                                    />
                                  )}
                                </>
                              );
                            })()}
                          </div>
                        ) : (
                          <p className="text-sm leading-relaxed whitespace-pre-wrap text-foreground">
                            {message.content}
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
                
                {/* Loading indicator */}
                {isPending && (
                  <div className="group w-full bg-background">
                    <div className="flex items-start gap-5 px-4 sm:px-6 py-6 max-w-4xl mx-auto">
                      <div className="shrink-0 w-8 h-8 rounded-full bg-linear-to-br from-primary/20 to-primary/5 flex items-center justify-center ring-2 ring-primary/10">
                        <svg className="w-4 h-4 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                        </svg>
                      </div>
                      <div className="flex-1 flex items-center gap-3">
                        <Loader2 className="h-4 w-4 animate-spin text-primary" />
                        <span className="text-sm text-muted-foreground">Thinking...</span>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* Fixed Input Area at Bottom */}
            <div className="absolute bottom-0 left-0 right-0 bg-background/95 backdrop-blur-sm border-t border-border">
              <div className="max-w-4xl mx-auto px-4 sm:px-6 py-4">
                <div className="relative bg-card rounded-2xl border border-border shadow-lg">
                  <textarea
                    ref={inputRef}
                    value={input}
                    onChange={handleInputChange}
                    onKeyDown={handleKeyDown}
                    placeholder="Ask me anything..."
                    className="w-full resize-none bg-transparent px-4 pt-4 pb-14 text-sm focus:outline-none min-h-14 max-h-40 rounded-2xl"
                    disabled={isPending || !settings.enabled}
                    rows={1}
                  />
                  
                  <div className="absolute bottom-0 left-0 right-0 flex items-center justify-between px-3 py-2.5">
                    <div className="flex items-center gap-1">
                      <Select
                        value={selectedModel}
                        onValueChange={setUserSelectedModel}
                        disabled={isPending}
                      >
                        <SelectTrigger className="h-8 w-auto min-w-[120px] border-0 bg-transparent hover:bg-muted/50 px-2 text-xs gap-1.5">
                          <Globe className="h-3.5 w-3.5" />
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="gpt-oss-120b">GPT-OSS-120B</SelectItem>
                        </SelectContent>
                      </Select>
                      <Select
                        value={selectedDomain}
                        onValueChange={(v) => setSelectedDomain(v as AIDomain)}
                        disabled={isPending}
                      >
                        <SelectTrigger className="h-8 w-auto min-w-[90px] border-0 bg-transparent hover:bg-muted/50 px-2 text-xs gap-1.5">
                          <SelectValue placeholder="Domain" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="general">Auto-detect</SelectItem>
                          {domainOptions.map((d) => (
                            <SelectItem key={d.id} value={d.id}>
                              {d.label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex items-center gap-1">
                      <Button
                        onClick={handleSend}
                        disabled={!input.trim() || isPending || !settings.enabled}
                        size="icon"
                        className="h-8 w-8 rounded-full bg-foreground text-background hover:bg-foreground/90 disabled:opacity-50"
                      >
                        {isPending ? (
                          <Loader2 className="h-4 w-4 animate-spin" />
                        ) : (
                          <Send className="h-4 w-4" />
                        )}
                      </Button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </>
        )}
      </div>

      {/* Detail Modals */}
      {viewingAccountId && (
        <AccountDetailModal
          key={`account-${viewingAccountId}`}
          accountId={viewingAccountId}
          open={!!viewingAccountId}
          onOpenChange={(open) => { if (!open) setViewingAccountId(null); }}
        />
      )}
      {viewingContactId && (
        <ContactDetailModal
          key={`contact-${viewingContactId}`}
          contactId={viewingContactId}
          open={!!viewingContactId}
          onOpenChange={(open) => { if (!open) setViewingContactId(null); }}
        />
      )}
      {viewingTaskId && (
        <TaskDetailModal
          key={`task-${viewingTaskId}`}
          taskId={viewingTaskId}
          open={!!viewingTaskId}
          onOpenChange={(open) => { if (!open) setViewingTaskId(null); }}
        />
      )}
      {viewingVisitReportId && (
        <VisitReportDetailModal
          key={`visit-${viewingVisitReportId}`}
          visitReportId={viewingVisitReportId}
          open={!!viewingVisitReportId}
          onOpenChange={(open) => { if (!open) setViewingVisitReportId(null); }}
        />
      )}
      {viewingDealId && (
        <DealDetailModal
          key={`deal-${viewingDealId}`}
          dealId={viewingDealId}
          open={!!viewingDealId}
          onOpenChange={(open) => { if (!open) setViewingDealId(null); }}
        />
      )}
      {viewingLeadId && (
        <LeadDetailModal
          key={`lead-${viewingLeadId}`}
          leadId={viewingLeadId}
          open={!!viewingLeadId}
          onOpenChange={(open) => { if (!open) setViewingLeadId(null); }}
        />
      )}
    </div>
  );
}
