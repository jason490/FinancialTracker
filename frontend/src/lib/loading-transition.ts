import { createEffect, createSignal, onCleanup, type Accessor } from "solid-js";

const DEFAULT_MIN_MS = 320;

// useMinLoadingHold keeps loading UI visible for at least minMs to avoid abrupt skeleton flashes.
export function useMinLoadingHold(loading: Accessor<boolean>, minMs = DEFAULT_MIN_MS) {
  const [holding, setHolding] = createSignal(true);
  let startedAt = Date.now();
  let timer: number | undefined;

  const clearTimer = () => {
    if (timer !== undefined) {
      window.clearTimeout(timer);
      timer = undefined;
    }
  };

  createEffect(() => {
    if (loading()) {
      clearTimer();
      startedAt = Date.now();
      setHolding(true);
      return;
    }

    const elapsed = Date.now() - startedAt;
    const remaining = Math.max(0, minMs - elapsed);
    if (remaining === 0) {
      setHolding(false);
      return;
    }

    timer = window.setTimeout(() => {
      setHolding(false);
      timer = undefined;
    }, remaining);
  });

  onCleanup(clearTimer);

  return holding;
}
