export type TagColorOption = {
  key: string;
  hex: string;
  name: string;
};

export const TAG_COLORS: TagColorOption[] = [
  { key: "blue", hex: "#3b82f6", name: "Blue" },
  { key: "indigo", hex: "#6366f1", name: "Indigo" },
  { key: "violet", hex: "#8b5cf6", name: "Violet" },
  { key: "emerald", hex: "#10b981", name: "Emerald" },
  { key: "teal", hex: "#14b8a6", name: "Teal" },
  { key: "cyan", hex: "#06b6d4", name: "Cyan" },
  { key: "amber", hex: "#f59e0b", name: "Amber" },
  { key: "orange", hex: "#f97316", name: "Orange" },
  { key: "rose", hex: "#f43f5e", name: "Rose" },
  { key: "red", hex: "#ef4444", name: "Red" },
  { key: "slate", hex: "#64748b", name: "Slate" },
];

const colorHexMap = Object.fromEntries(TAG_COLORS.map((c) => [c.key, c.hex]));

// tagColorHex resolves a palette key to its hex value, defaulting to slate.
export function tagColorHex(key: string): string {
  return colorHexMap[key] ?? colorHexMap.slate;
}

// defaultTagColor returns the default palette key for new tags.
export function defaultTagColor(): string {
  return "blue";
}
