import type { JSX } from "solid-js";
import styles from "./AppLogo.module.css";

type AppLogoProps = JSX.SvgSVGAttributes<SVGSVGElement> & {
  size?: number | string;
  title?: string;
};

// AppLogo renders the Financial Tracker mark: ascending ledger bars with a trend line.
export default function AppLogo(props: AppLogoProps) {
  const { size = 32, title = "Financial Tracker", class: className, ...rest } = props;

  return (
    <svg
      class={`${styles.logo} ${className ?? ""}`.trim()}
      width={size}
      height={size}
      viewBox="0 0 32 32"
      role="img"
      aria-label={title}
      {...rest}
    >
      <title>{title}</title>
      <rect class={styles.frame} width="32" height="32" rx="8" />
      <rect class={styles.bar} x="7" y="18" width="4.5" height="8" rx="1.3" opacity="0.3" />
      <rect class={styles.bar} x="13.75" y="15" width="4.5" height="11" rx="1.3" opacity="0.4" />
      <rect class={styles.bar} x="20.5" y="11" width="4.5" height="15" rx="1.3" opacity="0.5" />
      <path
        class={styles.trend}
        d="M6 22 12.5 16.5 17.5 19 25 9.5"
        fill="none"
        stroke-width="2.1"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
      <path
        class={styles.trend}
        d="M19.5 9.5 25 9.5 25 15"
        fill="none"
        stroke-width="2.1"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>
  );
}
