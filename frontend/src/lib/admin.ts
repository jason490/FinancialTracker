import { clientApiRequest } from "./api";

// createRegistrationCode issues a single-use invite code (admin only).
export async function createRegistrationCode(): Promise<{
  code: string;
  expires_at: number;
}> {
  return clientApiRequest<{ code: string; expires_at: number }>(
    "/api/v1/admin/registration-codes",
    {
      method: "POST",
    }
  );
}
