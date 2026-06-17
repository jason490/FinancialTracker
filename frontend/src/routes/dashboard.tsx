import { Title } from "@solidjs/meta";
import {
  Show,
  createResource,
  createSignal,
  onMount,
} from "solid-js";
import AppLayout from "~/layouts/AppLayout";
import DashboardGrid from "~/components/dashboard/DashboardGrid";
import { getDashboard } from "~/lib/dashboard";
import { useAuth } from "~/lib/auth-context";
import type { DashboardPayload } from "~/lib/types";
import styles from "~/styles/dashboard.module.css";

const HOLD_MS = 550;

export default function DashboardPage() {
  const { user: profile } = useAuth();
  const [editMode, setEditMode] = createSignal(false);
  const [dashboardData, setDashboardData] = createSignal<DashboardPayload>();

  const [dashboard] = createResource(
    profile,
    async () => {
      const payload = await getDashboard(false);
      setDashboardData(payload);
      return payload;
    }
  );

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

        <Show
          when={!dashboard.loading && (dashboardData() || dashboard())}
          fallback={<div class={styles.loadingState}>Loading dashboard...</div>}
        >
          <DashboardGrid
            data={() => dashboardData() ?? dashboard()}
            editMode={editMode}
            onSaved={handleSaved}
            onCancel={() => void exitEditMode()}
          />
        </Show>
      </div>
    </AppLayout>
  );
}
