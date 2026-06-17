import { clientApiRequest } from "./api";
import { loadUserTheme } from "./theme-store";
import type { LoginRequest, RegisterRequest, SessionProfile } from "./types";

// login authenticates a user. The backend handles cookie setting.
export async function login(body: Record<string, any>): Promise<{ status: string }> {
  const result = await clientApiRequest<{ status: string }>("/api/v1/auth/login", {
    method: "POST",
    body,
  });
  await loadUserTheme();
  return result;
}

// register creates a new account. The backend handles cookie setting.
export async function register(body: Record<string, any>): Promise<{ status: string }> {
  const result = await clientApiRequest<{ status: string }>("/api/v1/auth/register", {
    method: "POST",
    body,
  });
  await loadUserTheme();
  return result;
}

// forgotPassword requests a password reset email.
export async function forgotPassword(email: string): Promise<{ message: string }> {
  return clientApiRequest<{ message: string }>("/api/v1/auth/forgot-password", {
    method: "POST",
    body: { email },
  });
}

// logout clears the session on the backend.
export async function logout(): Promise<void> {
  return clientApiRequest<void>("/api/v1/auth/logout", {
    method: "POST",
  });
}

// getCurrentUser loads the authenticated session profile.
export async function getCurrentUser(): Promise<SessionProfile> {
  return clientApiRequest<SessionProfile>("/api/v1/auth/me");
}
