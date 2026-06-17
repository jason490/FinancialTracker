import { ParentProps, onMount } from "solid-js";
import { initThemeListener, loadUserTheme, setThemePreference } from "~/lib/theme-store";

// ThemeProvider applies the user theme globally for every route.
export default function ThemeProvider(props: ParentProps) {
  onMount(() => {
    setThemePreference("system");
    void loadUserTheme();
    return initThemeListener();
  });

  return props.children;
}
