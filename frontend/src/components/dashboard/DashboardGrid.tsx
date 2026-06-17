import Sortable from "sortablejs";
import {
  For,
  Show,
  createEffect,
  createSignal,
  onCleanup,
  type Accessor,
} from "solid-js";
import { DragHandle, EyeIcon, EyeOffIcon } from "~/components/icons";
import WidgetBody, { widgetLabel } from "~/components/dashboard/WidgetBody";
import { saveDashboardLayout } from "~/lib/dashboard";
import { widgetMeta, widgetsForRender } from "~/lib/dashboard-widgets";
import type { DashboardPayload, DashboardWidget } from "~/lib/types";
import styles from "~/styles/dashboard.module.css";

type DashboardGridProps = {
  data: Accessor<DashboardPayload | undefined>;
  editMode: Accessor<boolean>;
  onSaved: (payload: DashboardPayload) => void;
  onCancel: () => void;
};

// DashboardGrid renders the customizable widget grid with drag-and-drop edit mode.
export default function DashboardGrid(props: DashboardGridProps) {
  let gridRef!: HTMLDivElement;
  let sortable: Sortable | undefined;
  const [layout, setLayout] = createSignal<DashboardWidget[]>([]);
  const [saving, setSaving] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);
  const [deviceType, setDeviceType] = createSignal<"desktop" | "mobile">("desktop");

  createEffect(() => {
    if (typeof window === "undefined") return;
    const media = window.matchMedia("(max-width: 768px)");
    const handler = (e: MediaQueryListEvent | MediaQueryList) => {
      setDeviceType(e.matches ? "mobile" : "desktop");
    };
    handler(media);
    media.addEventListener("change", handler);
    onCleanup(() => media.removeEventListener("change", handler));
  });

  createEffect(() => {
    const payload = props.data();
    if (payload) {
      const widgets = deviceType() === "mobile" ? payload.layout.mobile : payload.layout.desktop;
      setLayout([...widgets].sort((a, b) => a.order - b.order));
    }
  });

  const syncOrderFromDom = () => {
    if (!gridRef) return;
    const ids = Array.from(gridRef.querySelectorAll("[data-widget-id]")).map(
      (node) => (node as HTMLElement).dataset.widgetId!
    );
    setLayout((current) => {
      const map = new Map(current.map((widget) => [widget.id, widget]));
      return ids
        .map((id, index) => {
          const widget = map.get(id);
          return widget ? { ...widget, order: index } : null;
        })
        .filter(Boolean) as DashboardWidget[];
    });
  };

  const destroySortable = () => {
    sortable?.destroy();
    sortable = undefined;
  };

  createEffect(() => {
    if (!props.editMode()) {
      destroySortable();
      return;
    }

    queueMicrotask(() => {
      if (!gridRef) return;
      destroySortable();
      sortable = Sortable.create(gridRef, {
        animation: 180,
        handle: ".dashboard-drag-handle",
        draggable: "[data-widget-id]",
        onEnd: syncOrderFromDom,
      });
    });
  });

  onCleanup(() => destroySortable());

  const toggleVisibility = (id: string) => {
    setLayout((current) =>
      current.map((widget) =>
        widget.id === id ? { ...widget, visible: !widget.visible } : widget
      )
    );
  };

  const handleSave = async () => {
    setSaving(true);
    setError(null);
    try {
      const payload = await saveDashboardLayout(deviceType(), layout());
      props.onSaved(payload);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save layout");
    } finally {
      setSaving(false);
    }
  };

  return (
    <section class={styles.gridSection}>
      <Show when={props.editMode()}>
        <div class={styles.editBar}>
          <p>
            Editing <strong>{deviceType()}</strong> layout. Drag widgets to reorder.
          </p>
          <div class={styles.editActions}>
            <button type="button" class={styles.secondaryButton} onClick={props.onCancel}>
              Cancel
            </button>
            <button
              type="button"
              class={styles.primaryButton}
              disabled={saving()}
              onClick={handleSave}
            >
              {saving() ? "Saving..." : "Save layout"}
            </button>
          </div>
        </div>
      </Show>

      <Show when={error()}>
        <div class={styles.errorBanner} role="alert">
          {error()}
        </div>
      </Show>

      <div class={styles.grid} ref={gridRef}>
        <Show when={props.data()}>
          {(payload) => (
            <For
              each={widgetsForRender(
                layout(),
                props.editMode()
              )}
            >
              {(widget) => {
                const meta = widgetMeta(widget.id);
                return (
                <div
                  class={styles.gridItem}
                  classList={{ [styles.span2]: meta?.span === 2 }}
                  data-widget-id={widget.id}
                  data-visible={widget.visible ? "true" : "false"}
                  style={{
                    "--widget-min-rows": String(meta?.minRows ?? 2),
                    "--widget-max-rows": String(meta?.maxRows ?? 3),
                  }}
                >
                  <Show when={props.editMode()}>
                    <div class={styles.widgetToolbar}>
                      <button
                        type="button"
                        class={`${styles.iconButton} dashboard-drag-handle`}
                        aria-label="Drag widget"
                        title="Drag to reorder"
                      >
                        <DragHandle />
                      </button>
                      <span class={styles.toolbarLabel}>{widgetLabel(widget.id)}</span>
                      <button
                        type="button"
                        class={styles.iconButton}
                        aria-label={widget.visible ? "Hide widget" : "Show widget"}
                        title={widget.visible ? "Hide widget" : "Show widget"}
                        onClick={() => toggleVisibility(widget.id)}
                      >
                        <Show
                          when={widget.visible}
                          fallback={<EyeOffIcon />}
                        >
                          <EyeIcon />
                        </Show>
                      </button>
                    </div>
                  </Show>

                  <div
                    class={styles.widgetPreview}
                    classList={{ [styles.widgetHiddenPreview]: props.editMode() && !widget.visible }}
                  >
                    <WidgetBody data={payload()} widgetId={widget.id} />
                  </div>
                </div>
              );
              }}
            </For>
          )}
        </Show>
      </div>
    </section>
  );
}
