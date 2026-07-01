import { getPublicApiUrl, getFrontendUrl } from "~/lib/env";
import styles from "~/styles/auth.module.css";

type SSOButtonsProps = {
  label: string;
  mode?: "login" | "register";
  registrationCode?: string;
  disabled?: boolean;
};

// SSOButtons renders provider sign-in links for authentication pages.
export default function SSOButtons(props: SSOButtonsProps) {
  const returnTo = `${getFrontendUrl()}/auth/sso/complete`;
  const mode = props.mode ?? "login";
  const params = new URLSearchParams({
    return_to: returnTo,
    action: mode,
  });
  if (props.registrationCode?.trim()) {
    params.set("registration_code", props.registrationCode.trim().toUpperCase());
  }
  const googleUrl = `${getPublicApiUrl()}/api/v1/auth/google?${params.toString()}`;
  const disabled = props.disabled || (mode === "register" && !props.registrationCode?.trim());

  return (
    <a
      class={styles.buttonSecondary}
      classList={{ [styles.buttonDisabled]: disabled }}
      href={disabled ? undefined : googleUrl}
      rel="external"
      aria-disabled={disabled ? "true" : undefined}
      tabIndex={disabled ? -1 : undefined}
      onClick={(event) => {
        if (disabled) {
          event.preventDefault();
        }
      }}
    >
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
