import { getPublicApiUrl, isCapacitorClient } from "./env";
import { ClientApiError } from "./api-error";
import type { APIError } from "./types";

type RequestOptions = {
  method?: string;
  body?: unknown;
};

// getCsrfToken reads the _csrf cookie value set by the backend.
function getCsrfToken(): string {
  if (typeof document === "undefined") return "";
  const match = document.cookie
    .split("; ")
    .find((row) => row.startsWith("_csrf="));
  return match ? decodeURIComponent(match.split("=")[1]) : "";
}

// Tracks whether the CSRF cookie has been primed for this page load.
let csrfPrimed = false;

// initCsrf ensures the _csrf cookie exists by hitting the GET /api/v1/csrf
// endpoint once per page load. Subsequent calls are no-ops.
export async function initCsrf(): Promise<void> {
  if (csrfPrimed || typeof window === "undefined") return;
  if (getCsrfToken()) {
    csrfPrimed = true;
    return;
  }
  const url = buildClientApiUrl("/api/v1/csrf");
  await fetch(url, { credentials: "include" });
  csrfPrimed = true;
}

function buildClientApiUrl(path: string): string {
  if (path.startsWith("http://") || path.startsWith("https://")) {
    return path;
  }

  // Capacitor must always use absolute URLs to reach the Go API.
  if (isCapacitorClient()) {
    const baseUrl = getPublicApiUrl();
    const normalizedBase = baseUrl.endsWith('/') ? baseUrl.slice(0, -1) : baseUrl;
    const normalizedPath = path.startsWith('/') ? path : `/${path}`;
    return `${normalizedBase}${normalizedPath}`;
  }

  // In the browser (Web), relative paths are preferred when using a reverse proxy
  // (like Caddy). This ensures the browser treats the API as same-origin,
  // making cookie management and CSRF protection seamless.
  if (typeof window !== "undefined") {
    return path.startsWith('/') ? path : `/${path}`;
  }

  // During SSG prerendering (Node/Nitro), relative URLs would fail, so we
  // fallback to the absolute public API URL.
  const baseUrl = getPublicApiUrl();
  const normalizedBase = baseUrl.endsWith('/') ? baseUrl.slice(0, -1) : baseUrl;
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  
  return `${normalizedBase}${normalizedPath}`;
}

// State-changing HTTP methods that require a CSRF token.
const csrfMethods = new Set(["POST", "PUT", "PATCH", "DELETE"]);

// clientApiRequest performs a request to the Go API from the browser.
// It relies on native cookie handling (credentials: "include") and
// automatically attaches the CSRF token for state-changing requests.
export async function clientApiRequest<T>(
  path: string,
  options: RequestOptions = {}
): Promise<T> {
  const method = (options.method || "GET").toUpperCase();

  // Prime the CSRF cookie before the first state-changing request.
  if (csrfMethods.has(method)) {
    await initCsrf();
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  // Attach the CSRF token for state-changing methods.
  if (csrfMethods.has(method)) {
    const token = getCsrfToken();
    if (token) {
      headers["X-CSRF-Token"] = token;
    }
  }

  const url = buildClientApiUrl(path);
  const response = await fetch(url, {
    method,
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
    credentials: "include",
  });

  // Handle non-OK responses first to avoid parsing issues if the error is HTML
  if (!response.ok) {
    if (response.status === 401 && typeof window !== "undefined" && !window.location.pathname.startsWith("/login")) {
      if (path !== "/api/v1/auth/me") {
        window.location.href = "/login";
        // We throw a specific error that can be ignored or handled if needed,
        // but the page reload will take over.
        throw new Error("Session expired. Redirecting to login...");
      }
    }

    let errorMessage = `API Request failed with status ${response.status}`;
    let errorCode = "request_failed";
    try {
      const errorData = await response.json() as APIError;
      errorMessage = errorData.message || errorMessage;
      errorCode = errorData.code || errorCode;
    } catch (e) {
      // If parsing fails, we keep the default error message
    }
    throw new ClientApiError(errorMessage, response.status, errorCode);
  }

  // Handle OK responses
  try {
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.includes("application/json")) {
      return await response.json() as T;
    }
    // If not JSON but OK, return as text or null
    const text = await response.text();
    return text as unknown as T;
  } catch (e) {
    throw new Error(`Failed to parse response from ${url}: ${e instanceof Error ? e.message : String(e)}`);
  }
}
