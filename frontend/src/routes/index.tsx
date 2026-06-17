import { Title } from "@solidjs/meta";
import { RedirectIfAuth } from "~/lib/auth-context";
import styles from "~/styles/intro.module.css";

// IntroPage is the public landing page for Financial Tracker.
export default function IntroPage() {
  return (
    <RedirectIfAuth>
      <main class={styles.page}>
        <Title>Financial Tracker</Title>

        <section class={styles.content}>
          <p class={styles.eyebrow}>Automated expense intelligence</p>
          <h1 class={styles.title}>Master your finances without the spreadsheet grind.</h1>
          <p class={styles.lede}>
            Financial Tracker connects to your banks through Plaid, organizes spending with
            smart tags, and gives you a dashboard you can shape around what matters.
          </p>

          <div class={styles.actions}>
            <a class={styles.primary} href="/register">
              Get started
            </a>
            <a class={styles.secondary} href="/login">
              Sign in
            </a>
          </div>

          <div class={styles.grid}>
            <article class={styles.card}>
              <h2>Plaid sync</h2>
              <p>Keep transactions current with cursor-based syncing and account reconciliation.</p>
            </article>
            <article class={styles.card}>
              <h2>Smart tagging</h2>
              <p>Apply exact, regex, and amount-based rules so categories stay consistent.</p>
            </article>
            <article class={styles.card}>
              <h2>Your dashboard</h2>
              <p>Arrange widgets on a responsive grid tuned for both mobile and desktop.</p>
            </article>
          </div>
        </section>
      </main>
    </RedirectIfAuth>
  );
}
