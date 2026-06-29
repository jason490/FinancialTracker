import { Title } from "@solidjs/meta";
import {
  Show,
  createEffect,
  createResource,
  createSignal,
  onCleanup,
  onMount,
} from "solid-js";
import AppLayout from "~/layouts/AppLayout";
import PageStatusBanner, { type PageStatus } from "~/components/PageStatusBanner";
import LoadingCrossfade from "~/components/LoadingCrossfade";
import DashboardGrid from "~/components/dashboard/DashboardGrid";
import { DashboardSkeletonGrid } from "~/components/dashboard/DashboardSkeletons";
import { getDashboard } from "~/lib/dashboard";
import {
  authTransitionActive,
  endAuthTransition,
  takePrefetchedDashboard,
} from "~/lib/auth-transition";
import { useMinLoadingHold } from "~/lib/loading-transition";
import { useAuth } from "~/lib/auth-context";
import type { DashboardPayload } from "~/lib/types";
import styles from "~/styles/dashboard.module.css";

const HOLD_MS = 550;
const TRANSITION_RELEASE_MS = 320;

export default function DashboardPage() {
  const { user: profile } = useAuth();
  const prefetchedDashboard = takePrefetchedDashboard();
  const [editMode, setEditMode] = createSignal(false);
  const [dashboardData, setDashboardData] = createSignal<DashboardPayload | undefined>(
    prefetchedDashboard
  );
  const [message, setMessage] = createSignal<PageStatus | null>(null);

  const handleSyncMessage = (text: string, type: PageStatus["type"]) => {
    setMessage({ text, type });
  };

  const [dashboard] = createResource(
    profile,
    async () => {
      if (prefetchedDashboard) {
        setDashboardData(prefetchedDashboard);
        return prefetchedDashboard;
      }

      const payload = await getDashboard(false);
      setDashboardData(payload);
      return payload;
    }
  );

  const dashboardLoading = () => dashboard.loading && !prefetchedDashboard;
  const dashboardReady = () => !!(dashboardData() ?? dashboard());
  const holdingDashboard = useMinLoadingHold(dashboardLoading);

  createEffect(() => {
    if (!authTransitionActive()) {
      return;
    }

    if (!dashboardReady() || dashboardLoading() || holdingDashboard()) {
      return;
    }

    const timer = window.setTimeout(() => {
      endAuthTransition();
    }, TRANSITION_RELEASE_MS);

    onCleanup(() => {
      window.clearTimeout(timer);
    });
  });

  let holdTimer: number | undefined;

  const clearHoldTimer = () => {
    if (holdTimer !== undefined) {
      window.clearTimeout(holdTimer);
      holdTimer = undefined;
    }
  };

  const enterEditMode = async () => {
    try {
      const payload = await getDashboard(true);
      setDashboardData(payload);
      setEditMode(true);
    } catch (err) {
      console.error(err);
    }
  };

  const exitEditMode = async () => {
    setEditMode(false);
    try {
      const payload = await getDashboard(false);
      setDashboardData(payload);
    } catch (err) {
      console.error(err);
    }
  };

  const handleSaved = (payload: DashboardPayload) => {
    setDashboardData(payload);
    setEditMode(false);
  };

  onMount(() => {
    const isCoarse = window.matchMedia("(pointer: coarse)").matches;

    const handlePointerDown = (event: PointerEvent) => {
      if (!isCoarse || editMode()) return;
      if (!(event.target instanceof Element)) return;
      if (!event.target.closest("[data-widget-id]")) return;

      clearHoldTimer();
      holdTimer = window.setTimeout(() => {
        void enterEditMode();
      }, HOLD_MS);
    };

    const handlePointerUp = () => clearHoldTimer();

    document.addEventListener("pointerdown", handlePointerDown);
    document.addEventListener("pointerup", handlePointerUp);
    document.addEventListener("pointercancel", handlePointerUp);

    return () => {
      clearHoldTimer();
      document.removeEventListener("pointerdown", handlePointerDown);
      document.removeEventListener("pointerup", handlePointerUp);
      document.removeEventListener("pointercancel", handlePointerUp);
    };
  });

  return (
    <AppLayout>
      <Title>Dashboard | Financial Tracker</Title>

      <div class={styles.page}>
        <header class={styles.header}>
          <div class={styles.headerTop}>
            <div>
              <p class={styles.eyebrow}>Overview</p>
              <h1 class={styles.title}>
                Welcome back
                <Show when={profile()}>
                  {(current) => `, ${current().first_name}`}
                </Show>
              </h1>
              <p class={styles.subtitle}>
                Your linked accounts, spending trends, and tagged activity in one
                customizable view.
              </p>
            </div>

            <div class={styles.headerActions}>
              <Show when={!editMode()}>
                <button
                  type="button"
                  class={styles.secondaryButton}
                  onClick={() => void enterEditMode()}
                >
                  Customize
                </button>
              </Show>
            </div>
          </div>
        </header>

        <PageStatusBanner message={message} onDismiss={() => setMessage(null)} />

        <LoadingCrossfade
          loading={dashboardLoading}
          ready={dashboardReady}
          skeleton={<DashboardSkeletonGrid />}
        >
          <DashboardGrid
            data={() => dashboardData() ?? dashboard()}
            editMode={editMode}
            onSaved={handleSaved}
            onCancel={() => void exitEditMode()}
            onSyncMessage={handleSyncMessage}
            reveal
          />
        </LoadingCrossfade>
      </div>
    </AppLayout>
  );
}
