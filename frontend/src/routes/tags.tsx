import { Title } from "@solidjs/meta";
import { For, Show, createSignal, onMount } from "solid-js";
import { createStore, produce, reconcile } from "solid-js/store";
import PageStatusBanner, { type PageStatus } from "~/components/PageStatusBanner";
import { PlusIcon } from "~/components/icons";
import CategoryCard from "~/components/tags/CategoryCard";
import DeleteCategoryModal from "~/components/tags/DeleteCategoryModal";
import NewCategoryModal from "~/components/tags/NewCategoryModal";
import TagFormModal from "~/components/tags/TagFormModal";
import { useTagDrag } from "~/components/tags/useTagDrag";
import AppLayout from "~/layouts/AppLayout";
import { getTags, moveTag } from "~/lib/tags";
import type { CategoryWithTagsView, TagView } from "~/lib/types";
import styles from "~/styles/tags.module.css";

type TagModalState =
  | { open: false }
  | { open: true; mode: "create"; category?: CategoryWithTagsView }
  | { open: true; mode: "edit"; tag: TagView };

// TagsPage lets users organize categories, tags, and auto-tagging rules.
export default function TagsPage() {
  // ── Data store ─────────────────────────────────────────────
  // Backed by a Solid store + `reconcile` so unchanged categories keep their
  // proxy identity across refetches. <For> therefore never remounts (and never
  // replays the `cardIn` entry animation) for rows that did not actually change.
  const [data, setData] = createStore<{ categories: CategoryWithTagsView[] }>({
    categories: [],
  });
  const [loaded, setLoaded] = createSignal(false);

  const applyPayload = (payload: { categories: CategoryWithTagsView[] }) => {
    setData("categories", reconcile(payload.categories, { key: "id", merge: true }));
  };

  const fetchTags = async () => {
    try {
      applyPayload(await getTags());
    } catch (err) {
      notify(err instanceof Error ? err.message : "Failed to load tags", "error");
    } finally {
      setLoaded(true);
    }
  };

  onMount(() => void fetchTags());

  // ── UI state ───────────────────────────────────────────────
  const [message, setMessage] = createSignal<PageStatus | null>(null);
  const [showNewCategory, setShowNewCategory] = createSignal(false);
  const [tagModal, setTagModal] = createSignal<TagModalState>({ open: false });
  const [deleteCategory, setDeleteCategory] = createSignal<CategoryWithTagsView>();

  const notify = (text: string, type: PageStatus["type"]) => setMessage({ text, type });
  const notifyOk = (text: string) => notify(text, "ok");
  const notifyErr = (text: string) => notify(text, "error");

  // ── Tag actions ────────────────────────────────────────────
  const handleMoveTag = async (tag: TagView, target: CategoryWithTagsView) => {
    if (tag.category_id === target.id) return;

    // Optimistic: surgically move the tag between category proxies. Only the
    // two affected `tags` arrays are touched; other cards are untouched.
    setData(
      produce((state) => {
        for (const cat of state.categories) {
          if (cat.id === tag.category_id) {
            cat.tags = cat.tags.filter((t) => t.id !== tag.id);
          } else if (cat.id === target.id) {
            cat.tags = [...cat.tags, { ...tag, category_id: target.id }];
          }
        }
      }),
    );

    try {
      applyPayload(await moveTag(tag.id, { category_id: target.id }));
      notifyOk(`Moved "${tag.name}" to ${target.name}.`);
    } catch (err) {
      void fetchTags();
      notifyErr(err instanceof Error ? err.message : "Failed to move tag");
    }
  };

  // ── Drag controller ────────────────────────────────────────
  // Single pointer-event drag controller for the whole page. Works for mouse,
  // pen, and touch so chips can be moved between categories on mobile too.
  const drag = useTagDrag({
    onCommit: (tag, targetId) => {
      const target = data.categories.find((c) => c.id === targetId);
      if (target) void handleMoveTag(tag, target);
    },
  });

  // ── Modal handlers ─────────────────────────────────────────
  const openCreateTag = (category: CategoryWithTagsView) =>
    setTagModal({ open: true, mode: "create", category });

  const openEditTag = (tag: TagView) => setTagModal({ open: true, mode: "edit", tag });

  const closeTagModal = () => setTagModal({ open: false });

  const handleTagSaved = () => {
    const current = tagModal();
    void fetchTags();
    notifyOk(current.open && current.mode === "create" ? "Tag created." : "Tag updated.");
  };

  // ── Render ─────────────────────────────────────────────────
  return (
    <AppLayout>
      <Title>Tags | Financial Tracker</Title>

      <div class={styles.page}>
        <header class={styles.header}>
          <div class={styles.headerCopy}>
            <p class={styles.eyebrow}>Organization</p>
            <h1 class={styles.title}>Tags & Categories</h1>
            <p class={styles.subtitle}>
              Shape how transactions get labeled. Define auto-tagging rules and keep your spending
              taxonomy sharp.
            </p>
          </div>

          <div class={styles.headerActions}>
            <button
              type="button"
              class={styles.primaryBtn}
              onClick={() => setShowNewCategory(true)}
            >
              <PlusIcon size={18} />
              New Category
            </button>
          </div>
        </header>

        <PageStatusBanner message={message} onDismiss={() => setMessage(null)} />

        <Show when={loaded()} fallback={<div class={styles.loading}>Loading tags...</div>}>
          <div class={styles.grid}>
            <For each={data.categories}>
              {(category, index) => (
                <CategoryCard
                  category={category}
                  index={index()}
                  draggingTag={drag.draggingTag()}
                  hoverCategoryId={drag.hoverCategoryId()}
                  onAddTag={openCreateTag}
                  onEditTag={openEditTag}
                  onDeleteCategory={setDeleteCategory}
                  onRefresh={fetchTags}
                  onError={notifyErr}
                  onTagPointerDown={drag.begin}
                />
              )}
            </For>
          </div>
        </Show>
      </div>

      <NewCategoryModal
        open={showNewCategory()}
        onClose={() => setShowNewCategory(false)}
        onSaved={() => {
          void fetchTags();
          notifyOk("Category created.");
        }}
        onError={notifyErr}
      />

      <TagFormModal
        state={tagModal()}
        categories={data.categories}
        onClose={closeTagModal}
        onSaved={handleTagSaved}
        onError={notifyErr}
      />

      <DeleteCategoryModal
        open={deleteCategory() != null}
        category={deleteCategory()}
        categories={data.categories}
        onClose={() => setDeleteCategory(undefined)}
        onSaved={() => {
          void fetchTags();
          notifyOk("Category deleted.");
        }}
        onError={notifyErr}
      />
    </AppLayout>
  );
}
