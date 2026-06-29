import { createSignal, onCleanup } from "solid-js";
import type { TagView } from "~/lib/types";

export type TagDragOptions = {
  onCommit: (tag: TagView, targetCategoryId: number) => void;
};

const DRAG_THRESHOLD_PX = 8;
const AUTO_SCROLL_EDGE_PX = 72;

// useTagDrag returns a single drag controller for moving tag chips between
// categories. It uses Pointer Events so the same code path works for mouse,
// pen, and touch input. Drop targets are any ancestor with
// `data-category-id`; a floating ghost element follows the pointer.
export function useTagDrag(options: TagDragOptions) {
  const [draggingTag, setDraggingTag] = createSignal<TagView | null>(null);
  const [hoverCategoryId, setHoverCategoryId] = createSignal<number | null>(null);

  let activePointerId = -1;
  let active = false;
  let pendingTag: TagView | null = null;
  let sourceChip: HTMLElement | null = null;
  let startX = 0;
  let startY = 0;
  let offsetX = 0;
  let offsetY = 0;
  let lastClientY = 0;
  let ghostEl: HTMLElement | null = null;
  let scrollFrame = 0;

  const removeGhost = () => {
    if (ghostEl) {
      ghostEl.remove();
      ghostEl = null;
    }
  };

  const stopAutoScroll = () => {
    if (scrollFrame !== 0) {
      cancelAnimationFrame(scrollFrame);
      scrollFrame = 0;
    }
  };

  const detachListeners = () => {
    document.removeEventListener("pointermove", onMove);
    document.removeEventListener("pointerup", onUp);
    document.removeEventListener("pointercancel", onUp);
  };

  const reset = () => {
    detachListeners();
    removeGhost();
    stopAutoScroll();
    document.body.style.removeProperty("user-select");
    document.body.style.removeProperty("touch-action");
    setDraggingTag(null);
    setHoverCategoryId(null);
    active = false;
    pendingTag = null;
    sourceChip = null;
    activePointerId = -1;
  };

  const buildGhost = () => {
    if (!sourceChip) return;
    const rect = sourceChip.getBoundingClientRect();
    const clone = sourceChip.cloneNode(true) as HTMLElement;
    clone.style.position = "fixed";
    clone.style.left = "0";
    clone.style.top = "0";
    clone.style.width = `${rect.width}px`;
    clone.style.margin = "0";
    clone.style.pointerEvents = "none";
    clone.style.zIndex = "1000";
    clone.style.opacity = "0.95";
    clone.style.boxShadow = "0 14px 36px rgba(0, 0, 0, 0.32)";
    clone.style.willChange = "transform";
    clone.setAttribute("aria-hidden", "true");
    document.body.appendChild(clone);
    ghostEl = clone;
  };

  const positionGhost = (x: number, y: number) => {
    if (!ghostEl) return;
    const gx = x - offsetX;
    const gy = y - offsetY;
    ghostEl.style.transform = `translate(${gx}px, ${gy}px) scale(1.04)`;
  };

  // findCategoryUnder hit-tests the DOM under the pointer and returns the
  // nearest `data-category-id`. The ghost is hidden during the hit-test so it
  // doesn't shadow the real drop target.
  const findCategoryUnder = (x: number, y: number): number | null => {
    const previousDisplay = ghostEl ? ghostEl.style.display : "";
    if (ghostEl) ghostEl.style.display = "none";
    const el = document.elementFromPoint(x, y);
    if (ghostEl) ghostEl.style.display = previousDisplay;
    if (!el) return null;
    const card = (el as Element).closest<HTMLElement>("[data-category-id]");
    if (!card) return null;
    const value = card.dataset.categoryId;
    return value ? Number(value) : null;
  };

  const tickAutoScroll = () => {
    scrollFrame = 0;
    if (!active) return;
    const vh = window.innerHeight;
    if (lastClientY < AUTO_SCROLL_EDGE_PX) {
      window.scrollBy(0, -Math.ceil((AUTO_SCROLL_EDGE_PX - lastClientY) / 4));
    } else if (lastClientY > vh - AUTO_SCROLL_EDGE_PX) {
      window.scrollBy(0, Math.ceil((lastClientY - (vh - AUTO_SCROLL_EDGE_PX)) / 4));
    }
    scrollFrame = requestAnimationFrame(tickAutoScroll);
  };

  const ensureAutoScroll = () => {
    if (scrollFrame === 0) {
      scrollFrame = requestAnimationFrame(tickAutoScroll);
    }
  };

  const activate = () => {
    if (active || !sourceChip || !pendingTag) return;
    active = true;
    const rect = sourceChip.getBoundingClientRect();
    offsetX = startX - rect.left;
    offsetY = startY - rect.top;
    buildGhost();
    positionGhost(startX, startY);
    document.body.style.userSelect = "none";
    document.body.style.touchAction = "none";
    setDraggingTag({ ...pendingTag });
    setHoverCategoryId(findCategoryUnder(startX, startY));
  };

  const onMove = (event: PointerEvent) => {
    if (event.pointerId !== activePointerId) return;

    if (!active) {
      const dx = event.clientX - startX;
      const dy = event.clientY - startY;
      if (Math.hypot(dx, dy) < DRAG_THRESHOLD_PX) return;
      activate();
    }

    event.preventDefault();
    lastClientY = event.clientY;
    positionGhost(event.clientX, event.clientY);
    setHoverCategoryId(findCategoryUnder(event.clientX, event.clientY));
    ensureAutoScroll();
  };

  const onUp = (event: PointerEvent) => {
    if (event.pointerId !== activePointerId) return;
    const wasActive = active;
    const tag = pendingTag;
    const targetId = wasActive ? findCategoryUnder(event.clientX, event.clientY) : null;
    reset();

    if (wasActive) {
      // Swallow the synthetic click browsers fire after a touch drag so the
      // chip's edit / delete buttons don't activate on drop.
      const swallow = (evt: MouseEvent) => {
        evt.stopPropagation();
        evt.preventDefault();
      };
      window.addEventListener("click", swallow, { capture: true, once: true });
      window.setTimeout(() => {
        window.removeEventListener("click", swallow, true);
      }, 250);
    }

    if (wasActive && tag && targetId != null && targetId !== tag.category_id) {
      options.onCommit(tag, targetId);
    }
  };

  // begin should be called from the chip's `onPointerDown`. It arms the drag
  // controller; the actual drag activates once movement crosses
  // DRAG_THRESHOLD_PX so taps still fire click handlers normally.
  const begin = (tag: TagView, event: PointerEvent) => {
    if (active || activePointerId !== -1) return;
    if (event.pointerType === "mouse" && event.button !== 0) return;
    const chip = event.currentTarget as HTMLElement | null;
    if (!chip) return;

    pendingTag = tag;
    sourceChip = chip;
    activePointerId = event.pointerId;
    startX = event.clientX;
    startY = event.clientY;
    lastClientY = event.clientY;

    document.addEventListener("pointermove", onMove, { passive: false });
    document.addEventListener("pointerup", onUp);
    document.addEventListener("pointercancel", onUp);
  };

  onCleanup(reset);

  return {
    draggingTag,
    hoverCategoryId,
    begin,
  };
}
