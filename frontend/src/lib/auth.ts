import { clientApiRequest } from "./api";
import { loadUserTheme } from "./theme-store";
import type { SessionProfile } from "./types";

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

// forgotPassword requests a password reset code by email.
export async function forgotPassword(email: string): Promise<{
  message: string;
  code_expires_in_seconds: number;
}> {
  return clientApiRequest<{ message: string; code_expires_in_seconds: number }>(
    "/api/v1/auth/forgot-password",
    {
      method: "POST",
      body: { email },
    }
  );
}

// verifyResetCode checks a reset code before setting a new password.
export async function verifyResetCode(body: {
  email: string;
  code: string;
}): Promise<{ expires_at: number }> {
  return clientApiRequest<{ expires_at: number }>("/api/v1/auth/verify-reset-code", {
    method: "POST",
    body,
  });
}

// resetPassword verifies a reset code and sets a new password.
export async function resetPassword(body: {
  email: string;
  code: string;
  new_password: string;
  confirm_password: string;
}): Promise<{ status: string }> {
  return clientApiRequest<{ status: string }>("/api/v1/auth/reset-password", {
    method: "POST",
    body,
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

// completeOnboarding marks the onboarding wizard as finished.
export async function completeOnboarding(): Promise<{ onboarding_completed: boolean }> {
  return clientApiRequest<{ onboarding_completed: boolean }>("/api/v1/auth/onboarding/complete", {
    method: "POST",
  });
}

// postAuthPath returns the route new sessions should land on.
export function postAuthPath(user: SessionProfile | undefined): string {
  if (user && !user.onboarding_completed) {
    return "/onboarding";
  }
  return "/dashboard";
}
