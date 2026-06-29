import { clientApiRequest } from "./api";
import type { DashboardLayout, DashboardPayload, DashboardWidget } from "./types";

// normalizeDashboardPayload coalesces nullable API arrays to empty lists.
function normalizeDashboardPayload(payload: DashboardPayload): DashboardPayload {
  return {
    ...payload,
    groups: payload.groups ?? {},
    transactions: payload.transactions ?? [],
    spending_trend: payload.spending_trend ?? [],
    spending_by_tag: payload.spending_by_tag ?? [],
    income_by_tag: payload.income_by_tag ?? [],
    layout: {
      desktop: payload.layout?.desktop ?? [],
      mobile: payload.layout?.mobile ?? [],
    },
  };
}

// getDashboard loads the dashboard payload for the authenticated user.
export async function getDashboard(editMode = false): Promise<DashboardPayload> {
  const query = editMode ? "?edit=1" : "";
  const payload = await clientApiRequest<DashboardPayload>(`/api/v1/dashboard${query}`);
  return normalizeDashboardPayload(payload);
}

// saveDashboardLayout persists a customized widget layout for a specific device.
export async function saveDashboardLayout(
  device_type: "desktop" | "mobile",
  widgets: DashboardWidget[]
): Promise<DashboardPayload> {
  const payload = await clientApiRequest<DashboardPayload>("/api/v1/dashboard/layout", {
    method: "POST",
    body: { device_type, widgets },
  });
  return normalizeDashboardPayload(payload);
}
