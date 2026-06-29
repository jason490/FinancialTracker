import { createSignal } from "solid-js";
import { getDashboard } from "./dashboard";
import type { DashboardPayload } from "./types";

type AuthTransitionCopy = {
  title: string;
  hint: string;
};

const DEFAULT_COPY: AuthTransitionCopy = {
  title: "Welcome back",
  hint: "Opening your dashboard",
};

const [active, setActive] = createSignal(false);
const [copy, setCopy] = createSignal<AuthTransitionCopy>(DEFAULT_COPY);

let prefetchedDashboard: DashboardPayload | undefined;
let prefetchPromise: Promise<DashboardPayload | undefined> | undefined;

// beginAuthTransition shows the global post-auth loading screen across route changes.
export function beginAuthTransition(nextCopy: Partial<AuthTransitionCopy> = {}) {
  setCopy({ ...DEFAULT_COPY, ...nextCopy });
  setActive(true);
}

// endAuthTransition hides the global post-auth loading screen.
export function endAuthTransition() {
  setActive(false);
  prefetchPromise = undefined;
}

// authTransitionActive reports whether the global post-auth loading screen is visible.
export function authTransitionActive() {
  return active();
}

// authTransitionCopy returns the active overlay title and hint text.
export function authTransitionCopy() {
  return copy();
}

// prefetchDashboardForAuth loads dashboard data during the auth transition.
export async function prefetchDashboardForAuth(): Promise<DashboardPayload | undefined> {
  if (prefetchedDashboard) {
    return prefetchedDashboard;
  }

  if (!prefetchPromise) {
    prefetchPromise = getDashboard(false)
      .then((payload) => {
        prefetchedDashboard = payload;
        return payload;
      })
      .catch(() => undefined);
  }

  return prefetchPromise;
}

// takePrefetchedDashboard returns cached dashboard data once for the destination page.
export function takePrefetchedDashboard(): DashboardPayload | undefined {
  const payload = prefetchedDashboard;
  prefetchedDashboard = undefined;
  prefetchPromise = undefined;
  return payload;
}
