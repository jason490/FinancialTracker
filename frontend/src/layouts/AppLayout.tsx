import { A, useLocation, useNavigate } from "@solidjs/router";
import { ParentProps, createEffect, Show } from "solid-js";
import { DashboardIcon, SettingsIcon, TagsIcon, TransactionsIcon } from "~/components/icons";
import AppLogo from "~/components/icons/AppLogo";
import { useAuth } from "~/lib/auth-context";
import { postAuthPath } from "~/lib/auth";
import styles from "./AppLayout.module.css";

export default function AppLayout(props: ParentProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const auth = useAuth();

  createEffect(() => {
    if (!auth.loading() && !auth.isAuthenticated()) {
      navigate("/login", { replace: true });
      return;
    }

    const user = auth.user();
    if (
      user &&
      !user.onboarding_completed &&
      location.pathname !== "/onboarding"
    ) {
      navigate("/onboarding", { replace: true });
    }
  });

  const navItems = [
    { label: "Dashboard", href: "/dashboard", icon: <DashboardIcon /> },
    { label: "Transactions", href: "/transactions", icon: <TransactionsIcon /> },
    { label: "Tags", href: "/tags", icon: <TagsIcon /> },
    { label: "Settings", href: "/settings", icon: <SettingsIcon /> },
  ];

  return (
    <Show when={auth.isAuthenticated()}>
      <div class={styles.layout}>
        <header class={styles.navbar}>
          <A href="/dashboard" class={styles.navBrand}>
            <AppLogo size={28} />
            <span>Financial Tracker</span>
          </A>
          <nav class={styles.navLinks}>
            {navItems.map((item) => (
              <A
                href={item.href}
                class={styles.navLink}
                activeClass={styles.navLinkActive}
              >
                {item.icon}
                <span>{item.label}</span>
              </A>
            ))}
          </nav>
        </header>

        <main class={styles.main}>{props.children}</main>

        <nav class={styles.tabBar}>
          {navItems.map((item) => (
            <A
              href={item.href}
              class={styles.tabItem}
              activeClass={styles.tabItemActive}
            >
              <span class={styles.tabIcon}>{item.icon}</span>
              <span>{item.label}</span>
            </A>
          ))}
        </nav>
      </div>
    </Show>
  );
}
