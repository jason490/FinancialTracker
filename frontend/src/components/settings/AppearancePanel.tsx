import { createEffect, createSignal, For } from "solid-js";
import { THEME_OPTIONS } from "~/lib/theme-options";
import { updateTheme } from "~/lib/settings";
import type { SettingsProfile } from "~/lib/types";
import styles from "~/styles/settings.module.css";

type AppearancePanelProps = {
  profile: SettingsProfile;
  onUpdated: (profile: SettingsProfile) => void;
  onMessage: (message: string, type: "ok" | "error") => void;
};

// swatchStyle builds the preview gradient for a theme card.
function swatchStyle(preview: [string, string], split?: boolean) {
  if (split) {
    return { background: `linear-gradient(135deg, ${preview[0]} 50%, ${preview[1]} 50%)` };
  }
  return { background: `linear-gradient(135deg, ${preview[0]}, ${preview[1]})` };
}

// AppearancePanel lets the user choose a theme preference.
export default function AppearancePanel(props: AppearancePanelProps) {
  const [theme, setTheme] = createSignal(props.profile.theme_preference || "system");
  const [pending, setPending] = createSignal(false);

  createEffect(() => {
    setTheme(props.profile.theme_preference || "system");
  });

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault();
    setPending(true);
    try {
      const result = await updateTheme(theme());
      props.onUpdated({
        ...props.profile,
        theme_preference: result.theme_preference,
      });
      props.onMessage("Theme preference saved.", "ok");
    } catch (err) {
      props.onMessage(err instanceof Error ? err.message : "Failed to save theme", "error");
    } finally {
      setPending(false);
    }
  };

  return (
    <div class={styles.panelInner}>
      <section>
        <h2 class={styles.sectionTitle}>Appearance</h2>
        <p class={styles.sectionHint}>
          Pick a palette for your workspace. System follows your device setting.
        </p>

        <form class={styles.card} onSubmit={handleSubmit}>
          <div class={styles.themeGrid}>
            <For each={THEME_OPTIONS}>
              {(option) => (
                <label class={styles.themeOption}>
                  <input
                    class={styles.themeInput}
                    type="radio"
                    name="theme"
                    value={option.id}
                    checked={theme() === option.id}
                    onChange={() => setTheme(option.id)}
                  />
                  <div class={styles.themeCard}>
                    <div
                      class={styles.themeSwatch}
                      style={swatchStyle(option.preview, option.split)}
                    />
                    <span class={styles.themeLabel}>{option.label}</span>
                    <span class={styles.themeDescription}>{option.description}</span>
                  </div>
                </label>
              )}
            </For>
          </div>

          <div class={styles.actions}>
            <button class={styles.buttonPrimary} type="submit" disabled={pending()}>
              {pending() ? "Saving..." : "Save preference"}
            </button>
          </div>
        </form>
      </section>
    </div>
  );
}
