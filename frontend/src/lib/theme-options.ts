export type ThemeOption = {
  id: string;
  label: string;
  description: string;
  preview: [string, string];
  split?: boolean;
};

// THEME_OPTIONS lists every selectable theme preference.
export const THEME_OPTIONS: ThemeOption[] = [
  {
    id: "light",
    label: "Light",
    description: "Clean daylight workspace",
    preview: ["#ffffff", "#e8f4f2"],
  },
  {
    id: "dark",
    label: "Dark",
    description: "Low-glare night mode",
    preview: ["#1a222c", "#11161d"],
  },
  {
    id: "system",
    label: "System",
    description: "Match your device",
    preview: ["#ffffff", "#11161d"],
    split: true,
  },
  {
    id: "tokyo-night",
    label: "Tokyo Night",
    description: "Cool indigo terminal glow",
    preview: ["#24283b", "#7aa2f7"],
  },
  {
    id: "coffee",
    label: "Coffee",
    description: "Warm roast and parchment",
    preview: ["#fffaf4", "#8b5e3c"],
  },
  {
    id: "forest",
    label: "Forest",
    description: "Deep moss and fern",
    preview: ["#142019", "#6fcf97"],
  },
  {
    id: "rose",
    label: "Rose",
    description: "Velvet dusk and blush",
    preview: ["#1c1218", "#e8a0b4"],
  },
  {
    id: "midnight",
    label: "Midnight",
    description: "Ink blue after hours",
    preview: ["#12182b", "#8b9cff"],
  },
  {
    id: "parchment",
    label: "Parchment",
    description: "Sunlit archive paper",
    preview: ["#fff8ee", "#a67c52"],
  },
];

// THEME_IDS is the set of valid persisted theme preference strings.
export const THEME_IDS = new Set(THEME_OPTIONS.map((theme) => theme.id));

// isThemePreference returns true when a value is a supported theme id.
export function isThemePreference(value: string): boolean {
  return THEME_IDS.has(value);
}
