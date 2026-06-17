import { isThemePreference } from "./theme-options";

// applyTheme sets the global data-theme attribute from a user preference.
export function applyTheme(preference: string) {
  if (typeof document === "undefined") return;

  if (preference === "system") {
    const resolved = window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
    document.documentElement.dataset.theme = resolved;
    return;
  }

  document.documentElement.dataset.theme = isThemePreference(preference) ? preference : "light";
}

// readCssVar returns a CSS custom property value from the document root.
export function readCssVar(name: string, fallback = "#888"): string {
  if (typeof document === "undefined") return fallback;
  const value = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
  return value || fallback;
}
