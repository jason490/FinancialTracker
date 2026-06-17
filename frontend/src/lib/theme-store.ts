import { getCurrentUser } from "./auth";
import { applyTheme } from "./themes";

let currentPreference = "system";

// getThemePreference returns the active theme preference string.
export function getThemePreference(): string {
  return currentPreference;
}

// setThemePreference applies and remembers a theme preference globally.
export function setThemePreference(preference: string) {
  currentPreference = preference || "system";
  applyTheme(currentPreference);
}

// loadUserTheme fetches the session profile and applies the user's theme.
export async function loadUserTheme(): Promise<void> {
  try {
    const user = await getCurrentUser();
    setThemePreference(user.theme_preference);
  } catch {
    setThemePreference("system");
  }
}

// initThemeListener reapplies the theme when the OS color scheme changes.
export function initThemeListener(): () => void {
  if (typeof window === "undefined") {
    return () => undefined;
  }

  const media = window.matchMedia("(prefers-color-scheme: dark)");
  const handler = () => {
    if (currentPreference === "system") {
      applyTheme("system");
    }
  };

  media.addEventListener("change", handler);
  return () => media.removeEventListener("change", handler);
}
