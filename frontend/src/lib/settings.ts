import { clientApiRequest } from "./api";
import { setThemePreference } from "./theme-store";
import type { SettingsProfile } from "./types";

// getSettings loads the authenticated user's settings profile.
export async function getSettings(): Promise<SettingsProfile> {
  return clientApiRequest<SettingsProfile>("/api/v1/settings");
}

// updateProfile saves the user's first and last name.
export async function updateProfile(firstName: string, lastName: string): Promise<SettingsProfile> {
  return clientApiRequest<SettingsProfile>("/api/v1/settings/profile", {
    method: "PATCH",
    body: { first_name: firstName, last_name: lastName },
  });
}

// updatePassword changes the user's password and invalidates the session.
export async function updatePassword(body: {
  current_password?: string;
  new_password: string;
  confirm_password: string;
}): Promise<{ logged_out: boolean }> {
  return clientApiRequest<{ logged_out: boolean }>("/api/v1/settings/password", {
    method: "POST",
    body,
  });
}

// updateTheme persists and applies a theme preference.
export async function updateTheme(theme: string): Promise<{ theme_preference: string }> {
  const result = await clientApiRequest<{ theme_preference: string }>("/api/v1/settings/theme", {
    method: "POST",
    body: { theme },
  });
  setThemePreference(result.theme_preference);
  return result;
}

// unlinkSSO removes a linked SSO provider from the account.
export async function unlinkSSO(provider: string): Promise<SettingsProfile> {
  return clientApiRequest<SettingsProfile>(`/api/v1/settings/unlink/${provider}`, {
    method: "POST",
  });
}

// getDeleteAccountReauthStatus checks whether the session has recent re-auth.
export async function getDeleteAccountReauthStatus(): Promise<{ reauth_verified: boolean }> {
  return clientApiRequest<{ reauth_verified: boolean }>("/api/v1/settings/delete/reauth-status");
}

// verifyDeleteAccount verifies password credentials before account deletion.
export async function verifyDeleteAccount(body: {
  email: string;
  password: string;
}): Promise<{ reauth_verified: boolean }> {
  return clientApiRequest<{ reauth_verified: boolean }>("/api/v1/settings/delete/verify", {
    method: "POST",
    body,
  });
}

// confirmDeleteAccount permanently deletes the authenticated user's account.
export async function confirmDeleteAccount(): Promise<{ deleted: boolean }> {
  return clientApiRequest<{ deleted: boolean }>("/api/v1/settings/delete/confirm", {
    method: "POST",
  });
}
