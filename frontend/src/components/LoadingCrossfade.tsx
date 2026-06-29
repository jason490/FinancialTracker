import { Show, type Accessor, type JSX } from "solid-js";
import { useMinLoadingHold } from "~/lib/loading-transition";
import styles from "~/styles/loading-crossfade.module.css";

type LoadingCrossfadeProps = {
  loading: Accessor<boolean>;
  ready: Accessor<boolean>;
  skeleton: JSX.Element;
  children: JSX.Element;
  minMs?: number;
  class?: string;
};

// LoadingCrossfade crossfades skeleton placeholders into loaded content without layout pops.
export default function LoadingCrossfade(props: LoadingCrossfadeProps) {
  const holding = useMinLoadingHold(props.loading, props.minMs);
  const showSkeleton = () => holding() || !props.ready();
  const showContent = () => props.ready() && !holding();

  return (
    <div class={`${styles.shell} ${props.class ?? ""}`}>
      <div
        class={`${styles.layer} ${styles.skeletonLayer}`}
        classList={{ [styles.skeletonHidden]: !showSkeleton() }}
        aria-busy={showSkeleton()}
        aria-hidden={!showSkeleton()}
      >
        {props.skeleton}
      </div>

      <Show when={props.ready()}>
        <div
          class={`${styles.layer} ${styles.contentLayer}`}
          classList={{ [styles.contentVisible]: showContent() }}
        >
          {props.children}
        </div>
      </Show>
    </div>
  );
}
