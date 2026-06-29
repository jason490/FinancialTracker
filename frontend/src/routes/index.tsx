import { Title } from "@solidjs/meta";
import { For } from "solid-js";
import DashboardPreview from "~/components/home/DashboardPreview";
import { BankIcon, DashboardIcon, SyncIcon, TagsIcon } from "~/components/icons";
import { RedirectIfAuth } from "~/lib/auth-context";
import styles from "~/styles/home.module.css";

const features = [
  {
    icon: SyncIcon,
    title: "Plaid sync",
    copy: "Keep transactions current with cursor-based syncing and account reconciliation.",
  },
  {
    icon: TagsIcon,
    title: "Smart tagging",
    copy: "Apply exact, regex, and amount-based rules so categories stay consistent.",
  },
  {
    icon: DashboardIcon,
    title: "Your dashboard",
    copy: "Arrange widgets on a responsive grid tuned for both mobile and desktop.",
  },
] as const;

const steps = [
  {
    number: "01",
    title: "Connect accounts",
    copy: "Link banks through Plaid or Stripe Financial Connections in a guided flow.",
  },
  {
    number: "02",
    title: "Let rules run",
    copy: "Tags apply automatically from merchant names, amounts, and patterns you define.",
  },
  {
    number: "03",
    title: "Shape your view",
    copy: "Snap widgets into place and track cash flow the way you actually think about money.",
  },
] as const;

// IntroPage is the public landing page for Financial Tracker.
export default function IntroPage() {
  return (
    <RedirectIfAuth>
      <div class={styles.page}>
        <Title>Financial Tracker — Clarity for every dollar</Title>

        <header class={styles.nav}>
          <a class={styles.brand} href="/">
            <span class={styles.brandMark} aria-hidden="true" />
            Financial Tracker
          </a>
          <nav class={styles.navActions} aria-label="Account">
            <a class={styles.navGhost} href="/login">
              Sign in
            </a>
            <a class={styles.navCta} href="/register">
              Get started
            </a>
          </nav>
        </header>

        <main class={styles.main}>
          <section class={styles.hero}>
            <div class={styles.heroCopy}>
              <p class={styles.eyebrow}>Automated expense intelligence</p>
              <h1 class={styles.title}>
                Master your finances without the spreadsheet grind.
              </h1>
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

              <ul class={styles.trustRow}>
                <li>
                  <BankIcon size={18} />
                  Bank-grade sync
                </li>
                <li>
                  <TagsIcon size={18} />
                  Rule-based tags
                </li>
                <li>
                  <DashboardIcon size={18} />
                  Custom widgets
                </li>
              </ul>
            </div>

            <div class={styles.heroVisual}>
              <DashboardPreview />
            </div>
          </section>

          <section class={styles.features} aria-labelledby="features-heading">
            <div class={styles.sectionIntro}>
              <p class={styles.sectionEyebrow}>Built for clarity</p>
              <h2 id="features-heading" class={styles.sectionTitle}>
                Everything you need to stay on top of spending
              </h2>
            </div>

            <div class={styles.grid}>
              <For each={features}>
                {(feature, index) => {
                  const Icon = feature.icon;
                  return (
                    <article
                      class={styles.card}
                      style={{ "animation-delay": `${0.12 + index() * 0.08}s` }}
                    >
                      <span class={styles.cardIcon} aria-hidden="true">
                        <Icon size={22} />
                      </span>
                      <h2>{feature.title}</h2>
                      <p>{feature.copy}</p>
                    </article>
                  );
                }}
              </For>
            </div>
          </section>

          <section class={styles.steps} aria-labelledby="steps-heading">
            <div class={styles.sectionIntro}>
              <p class={styles.sectionEyebrow}>How it works</p>
              <h2 id="steps-heading" class={styles.sectionTitle}>
                From linked accounts to living insight in minutes
              </h2>
            </div>

            <ol class={styles.stepList}>
              <For each={steps}>
                {(step, index) => (
                  <li
                    class={styles.step}
                    style={{ "animation-delay": `${0.1 + index() * 0.1}s` }}
                  >
                    <span class={styles.stepNumber}>{step.number}</span>
                    <div>
                      <h3>{step.title}</h3>
                      <p>{step.copy}</p>
                    </div>
                  </li>
                )}
              </For>
            </ol>
          </section>

          <section class={styles.ctaBand} aria-labelledby="cta-heading">
            <div class={styles.ctaInner}>
              <p class={styles.sectionEyebrow}>Ready when you are</p>
              <h2 id="cta-heading" class={styles.ctaTitle}>
                Stop reconciling. Start understanding.
              </h2>
              <p class={styles.ctaCopy}>
                Create a free account, connect your first bank, and see transactions organize
                themselves.
              </p>
              <div class={styles.actions}>
                <a class={styles.primary} href="/register">
                  Get started free
                </a>
                <a class={styles.secondary} href="/login">
                  Sign in
                </a>
              </div>
            </div>
          </section>
        </main>

        <footer class={styles.footer}>
          <p>Financial Tracker</p>
          <p class={styles.footerMuted}>Clarity for every dollar you move.</p>
        </footer>
      </div>
    </RedirectIfAuth>
  );
}
