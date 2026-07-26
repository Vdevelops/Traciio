"use client";

import { memo, useCallback, type ComponentType } from "react";
import { useRouter } from "next/navigation";
import {
  MapPin,
  TrendingUp,
  Package,
  Users,
  BarChart3,
  Settings,
  ClipboardList,
  Calendar,
  Target,
  FileText,
  User,
  Building,
  Phone,
  ExternalLink,
  Eye,
} from "lucide-react";
import type { AIActionCard, AIActionIcon, AIActionEntity } from "../types";

/** Regex pattern to extract action cards from AI responses */
const ACTION_PATTERN = /<!-- ACTION:(.*?) -->/g;

/** Icon mapping from string identifiers to Lucide components */
const iconMap: Record<AIActionIcon, ComponentType<{ className?: string }>> = {
  map: MapPin,
  "trending-up": TrendingUp,
  package: Package,
  users: Users,
  "bar-chart": BarChart3,
  settings: Settings,
  clipboard: ClipboardList,
  calendar: Calendar,
  target: Target,
  "file-text": FileText,
  user: User,
  building: Building,
  phone: Phone,
};

/**
 * Parse action card markers from an AI response string.
 * Returns the cleaned message (without markers) and parsed action cards.
 */
export function parseActionCards(message: string): {
  cleanMessage: string;
  actions: AIActionCard[];
} {
  const actions: AIActionCard[] = [];
  let match: RegExpExecArray | null;

  // Reset regex state
  ACTION_PATTERN.lastIndex = 0;

  while ((match = ACTION_PATTERN.exec(message)) !== null) {
    try {
      const parsed = JSON.parse(match[1]) as Partial<AIActionCard>;
      const isSupportedType = parsed.type === "navigate" || parsed.type === "detail";
      if (isSupportedType && parsed.label) {
        actions.push(parsed as AIActionCard);
      }
    } catch {
      // Skip malformed action cards
    }
  }

  // Remove action markers from the message
  const cleanMessage = message.replace(ACTION_PATTERN, "").trim();

  return { cleanMessage, actions };
}

interface AIActionCardsProps {
  actions: AIActionCard[];
  /** Callback when a "detail" action is clicked - opens entity detail modal */
  onEntityClick?: (entity: AIActionEntity, entityId: string) => void;
}

/**
 * Renders AI action cards as clickable cards below the AI message.
 * Supports two action types:
 * - "navigate": redirects to a CRM page
 * - "detail": opens an entity detail modal
 */
export const AIActionCards = memo(function AIActionCards({
  actions,
  onEntityClick,
}: AIActionCardsProps) {
  const router = useRouter();

  const handleClick = useCallback(
    (action: AIActionCard) => {
      if (action.type === "navigate" && action.url) {
        router.push(action.url);
      } else if (action.type === "detail" && action.entity && action.entityId) {
        onEntityClick?.(action.entity, action.entityId);
      }
    },
    [router, onEntityClick],
  );

  if (actions.length === 0) return null;

  return (
    <div className="flex flex-wrap gap-2 mt-3">
      {actions.map((action, index) => {
        const IconComponent = action.icon ? iconMap[action.icon] ?? ExternalLink : ExternalLink;
        const isDetail = action.type === "detail";

        return (
          <button
            // Ensure keys are unique even when multiple actions share the same url/entityId
            key={`${action.type}-${action.url ?? action.entityId ?? 'item'}-${index}`}
            type="button"
            onClick={() => handleClick(action)}
            className="group flex items-center gap-3 px-4 py-3 rounded-xl border border-border bg-card hover:bg-muted/50 hover:border-primary/30 transition-all cursor-pointer text-left max-w-xs"
          >
            <div className="shrink-0 w-9 h-9 rounded-lg bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
              <IconComponent className="w-4.5 h-4.5 text-primary" />
            </div>
            <div className="flex-1 min-w-0">
              <div className="text-sm font-medium text-foreground truncate">
                {action.label}
              </div>
              {action.description && (
                <div className="text-xs text-muted-foreground truncate">
                  {action.description}
                </div>
              )}
            </div>
            <div className="shrink-0 text-muted-foreground group-hover:text-primary transition-colors">
              {isDetail ? (
                <Eye className="w-4 h-4" />
              ) : (
                <ExternalLink className="w-4 h-4" />
              )}
            </div>
          </button>
        );
      })}
    </div>
  );
});
