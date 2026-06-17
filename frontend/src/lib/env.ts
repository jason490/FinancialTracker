// API_URL is used by SolidStart server functions to reach the Go API on the internal network.
export function getApiUrl(): string {
  return process.env.API_URL || "http://localhost:8080";
}

// API_PUBLIC_URL is used by the browser for SSO redirects and Capacitor builds.
export function getPublicApiUrl(): string {
  if (typeof window !== "undefined") {
    // Web builds are served behind the same reverse proxy as /api/* routes.
    return import.meta.env.VITE_API_PUBLIC_URL || window.location.origin;
  }
  return process.env.API_PUBLIC_URL || process.env.API_URL || "http://localhost:8080";
}

export function isCapacitorClient(): boolean {
  return typeof window !== "undefined" && window.location.protocol === "capacitor:";
}

// FRONTEND_URL is used when building SSO return links.
export function getFrontendUrl(): string {
  if (typeof window !== "undefined") {
    return import.meta.env.VITE_FRONTEND_URL || window.location.origin;
  }
  return process.env.FRONTEND_URL || "http://localhost:3000";
}
