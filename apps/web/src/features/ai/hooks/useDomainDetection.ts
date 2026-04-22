import { useCallback, useMemo } from "react";
import type { AIDomain, AIDataPrivacySettings } from "../types";
import {
  DOMAIN_MODULES,
  SPECIFIC_DOMAINS,
  isDomainAccessible,
} from "../data/domain-registry";

/**
 * Detects the most relevant AI domain from user input using keyword matching.
 *
 * Scoring: each keyword match adds 1 point. The domain with the highest score wins.
 * If no domain scores above the threshold, falls back to "general".
 *
 * This runs entirely client-side to avoid sending unnecessary tokens to the LLM
 * for intent classification; the backend also performs its own keyword matching,
 * but sending a domain hint from the frontend improves accuracy and reduces
 * context-fetching overhead.
 */
export function useDomainDetection(privacy: AIDataPrivacySettings) {
  const accessibleDomains = useMemo(
    () => SPECIFIC_DOMAINS.filter((d) => isDomainAccessible(d, privacy)),
    [privacy],
  );

  /**
   * Detect domain from a user message string.
   * Returns the best-matching domain and its confidence score.
   */
  const detectDomain = useCallback(
    (message: string): { domain: AIDomain; confidence: number } => {
      const normalized = message.toLowerCase().trim();

      if (!normalized) return { domain: "general", confidence: 0 };

      let bestDomain: AIDomain = "general";
      let bestScore = 0;

      for (const domainId of accessibleDomains) {
        const domainDef = DOMAIN_MODULES[domainId];
        let score = 0;

        for (const keyword of domainDef.intentKeywords) {
          if (normalized.includes(keyword.toLowerCase())) {
            score += 1;

            // Bonus for exact word boundary matches using split-based check
            if (hasWordBoundary(normalized, keyword.toLowerCase())) {
              score += 0.5;
            }
          }
        }

        if (score > bestScore) {
          bestScore = score;
          bestDomain = domainId;
        }
      }

      // Require at least 1 keyword match to assign a specific domain
      const confidence = bestScore > 0 ? Math.min(bestScore / 3, 1) : 0;
      if (bestScore < 1) {
        return { domain: "general", confidence: 0 };
      }

      return { domain: bestDomain, confidence };
    },
    [accessibleDomains],
  );

  /**
   * Get the prompt fragment for a given domain.
   */
  const getDomainPrompt = useCallback((domain: AIDomain): string => {
    return DOMAIN_MODULES[domain].promptFragment;
  }, []);

  /**
   * Get all accessible domain labels for UI display.
   */
  const domainOptions = useMemo(
    () =>
      accessibleDomains.map((d) => ({
        id: d,
        label: DOMAIN_MODULES[d].label,
        description: DOMAIN_MODULES[d].description,
      })),
    [accessibleDomains],
  );

  return {
    detectDomain,
    getDomainPrompt,
    domainOptions,
    accessibleDomains,
  };
}

/**
 * Check if a keyword appears as a standalone word (not as substring of another word).
 * Uses character-based boundary detection instead of regex to avoid escaping issues.
 */
function hasWordBoundary(text: string, keyword: string): boolean {
  const idx = text.indexOf(keyword);
  if (idx === -1) return false;

  const charBefore = idx > 0 ? text[idx - 1] : " ";
  const charAfter = idx + keyword.length < text.length ? text[idx + keyword.length] : " ";

  const isWordChar = (c: string) => /\w/.test(c);
  return !isWordChar(charBefore) && !isWordChar(charAfter);
}
