import { createContext, useContext, createResource, ParentProps, createEffect, Show } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { getCurrentUser, postAuthPath } from "./auth";
import { authTransitionActive } from "./auth-transition";
import type { SessionProfile } from "./types";

interface AuthContextValue {
  user: () => SessionProfile | undefined;
  loading: () => boolean;
  error: () => any;
  refetch: () => void;
  isAuthenticated: () => boolean;
}

const AuthContext = createContext<AuthContextValue>();

export function AuthProvider(props: ParentProps) {
  const [user, { refetch }] = createResource(async () => {
    try {
      return await getCurrentUser();
    } catch (err) {
      return undefined;
    }
  });

  const value: AuthContextValue = {
    user: () => user(),
    loading: () => user.loading,
    error: () => user.error,
    refetch,
    isAuthenticated: () => !!user(),
  };

  return (
    <AuthContext.Provider value={value}>
      {props.children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}

/**
 * RequireAuth protects a route by redirecting unauthenticated users to /login.
 */
export function RequireAuth(props: ParentProps) {
  const auth = useAuth();
  const navigate = useNavigate();

  createEffect(() => {
    if (!auth.loading() && !auth.isAuthenticated()) {
      navigate("/login", { replace: true });
    }
  });

  return (
    <Show when={auth.isAuthenticated()} fallback={null}>
      {props.children}
    </Show>
  );
}

/**
 * RedirectIfAuth redirects authenticated users away from a page (e.g. login).
 */
export function RedirectIfAuth(props: ParentProps) {
  const auth = useAuth();
  const navigate = useNavigate();

  createEffect(() => {
    if (authTransitionActive()) {
      return;
    }

    if (auth.isAuthenticated()) {
      navigate(postAuthPath(auth.user()), { replace: true });
    }
  });

  return <>{props.children}</>;
}
