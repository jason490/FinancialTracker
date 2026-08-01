import { For, Show, createSignal, onMount } from "solid-js";
import { createStore, produce, reconcile } from "solid-js/store";
import { PlusIcon } from "~/components/icons";
import CategoryCard from "~/components/tags/CategoryCard";
import DeleteCategoryModal from "~/components/tags/DeleteCategoryModal";
import NewCategoryModal from "~/components/tags/NewCategoryModal";
import TagFormModal from "~/components/tags/TagFormModal";
import { useTagDrag } from "~/components/tags/useTagDrag";
import { getTags, moveTag } from "~/lib/tags";
import type { CategoryWithTagsView, TagView } from "~/lib/types";
import settingsStyles from "~/styles/settings.module.css";
import tagStyles from "~/styles/tags.module.css";

type TagsPanelProps = {
  onMessage: (message: string, type: "ok" | "error" | "info") => void;
};

type TagModalState =
  | { open: false }
  | { open: true; mode: "create"; category?: CategoryWithTagsView }
  | { open: true; mode: "edit"; tag: TagView };

// TagsPanel lets users organize categories, tags, and auto-tagging rules from Settings.
export default function TagsPanel(props: TagsPanelProps) {
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

  const notifyOk = (text: string) => props.onMessage(text, "ok");
  const notifyErr = (text: string) => props.onMessage(text, "error");

  const fetchTags = async () => {
    try {
      applyPayload(await getTags());
    } catch (err) {
      notifyErr(err instanceof Error ? err.message : "Failed to load tags");
    } finally {
      setLoaded(true);
    }
  };

  onMount(() => void fetchTags());

  const [showNewCategory, setShowNewCategory] = createSignal(false);
  const [tagModal, setTagModal] = createSignal<TagModalState>({ open: false });
  const [deleteCategory, setDeleteCategory] = createSignal<CategoryWithTagsView>();

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

  // Single pointer-event drag controller for the panel. Works for mouse, pen,
  // and touch so chips can be moved between categories on mobile too.
  const drag = useTagDrag({
    onCommit: (tag, targetId) => {
      const target = data.categories.find((c) => c.id === targetId);
      if (target) void handleMoveTag(tag, target);
    },
  });

  const openCreateTag = (category: CategoryWithTagsView) =>
    setTagModal({ open: true, mode: "create", category });

  const openEditTag = (tag: TagView) => setTagModal({ open: true, mode: "edit", tag });

  const closeTagModal = () => setTagModal({ open: false });

  const handleTagSaved = () => {
    const current = tagModal();
    void fetchTags();
    notifyOk(current.open && current.mode === "create" ? "Tag created." : "Tag updated.");
  };

  return (
    <>
      <div class={settingsStyles.panelInner}>
        <section class={settingsStyles.connectionsIntro}>
          <div class={settingsStyles.toolbar}>
            <div>
              <h2 class={settingsStyles.sectionTitle}>Tags & Categories</h2>
              <p class={settingsStyles.sectionHint}>
                Shape how transactions get labeled. Define auto-tagging rules and keep your spending
                taxonomy sharp.
              </p>
            </div>
            <div class={settingsStyles.actions}>
              <button
                type="button"
                class={tagStyles.primaryBtn}
                onClick={() => setShowNewCategory(true)}
              >
                <PlusIcon size={18} />
                New Category
              </button>
            </div>
          </div>

          <Show when={loaded()} fallback={<div class={tagStyles.loading}>Loading tags...</div>}>
            <div class={tagStyles.grid}>
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
        </section>
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
    </>
  );
}
