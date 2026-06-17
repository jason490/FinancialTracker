import { For } from "solid-js";
import { TAG_COLORS } from "~/lib/tag-colors";
import styles from "~/styles/tags.module.css";

type TagColorPickerProps = {
  value: string;
  onChange: (color: string) => void;
};

// TagColorPicker renders the selectable tag color palette.
export default function TagColorPicker(props: TagColorPickerProps) {
  return (
    <div class={styles.field}>
      <span class={styles.label}>Color</span>
      <div class={styles.colorGrid}>
        <For each={TAG_COLORS}>
          {(color) => (
            <button
              type="button"
              class={`${styles.colorSwatch} ${props.value === color.key ? styles.colorSwatchSelected : ""}`}
              style={{ "background-color": color.hex }}
              title={color.name}
              aria-label={color.name}
              aria-pressed={props.value === color.key}
              onClick={() => props.onChange(color.key)}
            />
          )}
        </For>
      </div>
    </div>
  );
}
