import { useCallback, useEffect, useState } from "react";

interface UseDashboardCommandPaletteResult {
  readonly isOpen: boolean;
  readonly open: () => void;
  readonly close: () => void;
  readonly toggle: () => void;
}

export function useDashboardCommandPalette(): UseDashboardCommandPaletteResult {
  const [isOpen, setIsOpen] = useState(false);

  const open = useCallback(() => setIsOpen(true), []);
  const close = useCallback(() => setIsOpen(false), []);
  const toggle = useCallback(() => {
    setIsOpen((prev) => !prev);
  }, []);

  // Global keyboard shortcuts: "/" and "Ctrl+K" / "Meta+K"
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      const target = event.target as HTMLElement | null;
      const isInputLike =
        target?.tagName === "INPUT" ||
        target?.tagName === "TEXTAREA" ||
        target?.isContentEditable;

      // Ignore when typing in inputs
      if (isInputLike) {
        return;
      }

      // "/" opens palette
      if (event.key === "/" && !event.ctrlKey && !event.metaKey && !event.altKey) {
        event.preventDefault();
        setIsOpen(true);
        return;
      }

      // "Ctrl+K" or "Meta+K"
      if (
        (event.key === "k" || event.key === "K") &&
        (event.ctrlKey || event.metaKey)
      ) {
        event.preventDefault();
        setIsOpen((prev) => !prev);
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  return {
    isOpen,
    open,
    close,
    toggle,
  };
}




