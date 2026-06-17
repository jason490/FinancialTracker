import { getPublicApiUrl, getFrontendUrl } from "~/lib/env";
import styles from "~/styles/auth.module.css";

// SSOButtons renders provider sign-in links for authentication pages.
export default function SSOButtons(props: { label: string }) {
  const returnTo = `${getFrontendUrl()}/auth/sso/complete`;
  const googleUrl = `${getPublicApiUrl()}/api/v1/auth/google?return_to=${encodeURIComponent(returnTo)}`;

  return (
    <a class={styles.buttonSecondary} href={googleUrl} rel="external">
      <img
        src="https://www.gstatic.com/firebasejs/ui/2.0.0/images/auth/google.svg"
        alt=""
        width="20"
        height="20"
      />
      <span>{props.label}</span>
    </a>
  );
}
