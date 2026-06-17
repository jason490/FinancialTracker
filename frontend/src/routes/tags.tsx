import { Title } from "@solidjs/meta";
import { For, Show, createResource, createSignal } from "solid-js";
import { PlusIcon } from "~/components/icons";
import CategoryCard from "~/components/tags/CategoryCard";
import DeleteCategoryModal from "~/components/tags/DeleteCategoryModal";
import NewCategoryModal from "~/components/tags/NewCategoryModal";
import TagFormModal from "~/components/tags/TagFormModal";
import AppLayout from "~/layouts/AppLayout";
import { getTags } from "~/lib/tags";
import type { CategoryWithTagsView, TagView } from "~/lib/types";
import styles from "~/styles/tags.module.css";

// TagsPage lets users organize categories, tags, and auto-tagging rules.
export default function TagsPage() {
  const [tags, { refetch }] = createResource(getTags);
  const [message, setMessage] = createSignal<{ text: string; type: "ok" | "error" } | null>(null);

  const [showNewCategory, setShowNewCategory] = createSignal(false);
  const [tagModal, setTagModal] = createSignal<{
    open: boolean;
    mode: "create" | "edit";
    category?: CategoryWithTagsView;
    tag?: TagView;
  }>({ open: false, mode: "create" });
  const [deleteCategory, setDeleteCategory] = createSignal<CategoryWithTagsView | undefined>();

  const categories = () => tags()?.categories ?? [];

  const notify = (text: string, type: "ok" | "error") => {
    setMessage({ text, type });
  };

  const refresh = () => {
    void refetch();
  };

  const openCreateTag = (category: CategoryWithTagsView) => {
    setTagModal({ open: true, mode: "create", category });
  };

  const openEditTag = (tag: TagView) => {
    setTagModal({ open: true, mode: "edit", tag });
  };

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

        <Show when={message()}>
          {(current) => (
            <div
              class={current().type === "error" ? styles.statusError : styles.statusOk}
              role="status"
            >
              {current().text}
            </div>
          )}
        </Show>

        <Show when={!tags.loading} fallback={<div class={styles.loading}>Loading tags...</div>}>
          <div class={styles.grid}>
            <For each={categories()}>
              {(category, index) => (
                <CategoryCard
                  category={category}
                  index={index()}
                  onAddTag={openCreateTag}
                  onEditTag={openEditTag}
                  onDeleteCategory={setDeleteCategory}
                  onRefresh={refresh}
                  onError={(text) => notify(text, "error")}
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
          refresh();
          notify("Category created.", "ok");
        }}
        onError={(text) => notify(text, "error")}
      />

      <TagFormModal
        open={tagModal().open}
        mode={tagModal().mode}
        category={tagModal().category}
        tag={tagModal().tag}
        categories={categories()}
        onClose={() => setTagModal({ open: false, mode: "create" })}
        onSaved={() => {
          refresh();
          notify(tagModal().mode === "create" ? "Tag created." : "Tag updated.", "ok");
        }}
        onError={(text) => notify(text, "error")}
      />

      <DeleteCategoryModal
        open={deleteCategory() != null}
        category={deleteCategory()}
        categories={categories()}
        onClose={() => setDeleteCategory(undefined)}
        onSaved={() => {
          refresh();
          notify("Category deleted.", "ok");
        }}
        onError={(text) => notify(text, "error")}
      />
    </AppLayout>
  );
}
