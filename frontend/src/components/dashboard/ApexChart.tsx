import ApexCharts, { type ApexOptions } from "apexcharts";
import { onCleanup, onMount, createEffect, on } from "solid-js";

type ApexChartProps = {
  options: ApexOptions;
  series: ApexOptions["series"];
  height?: number | string;
  class?: string;
};

// ApexChart renders an ApexCharts instance that reacts to option updates.
export default function ApexChart(props: ApexChartProps) {
  let container!: HTMLDivElement;
  let chart: ApexCharts | undefined;
  let isDestroyed = false;
  let resizeObserver: ResizeObserver | undefined;
  let resizeFrame = 0;
  let lastWidth = 0;

  // refit recomputes chart dimensions from the current container size.
  // ApexCharts caches the width measured at render time and only reacts to
  // window resizes, so on SPA route changes the container can be measured at
  // the wrong width and the SVG overflows its bounds. Re-running updateOptions
  // forces ApexCharts to re-read the container and resize to fit.
  const refit = () => {
    if (!chart || isDestroyed) return;
    chart
      .updateOptions({ chart: { width: "100%" } }, false, false, false)
      .catch(() => {});
  };

  onMount(() => {
    if (!container) return;

    chart = new ApexCharts(container, {
      ...props.options,
      series: props.series,
      chart: {
        ...(props.options.chart ?? {}),
        width: "100%",
        height: props.height ?? 240,
      },
    });

    chart.render().catch((e) => {
      console.warn("ApexCharts render error ignored (likely unmounted):", e);
    });

    lastWidth = container.clientWidth;

    // Watch the container itself (not just the window) so the chart re-fits
    // whenever its width changes, e.g. after navigating back to the dashboard.
    resizeObserver = new ResizeObserver((entries) => {
      const width = entries[0]?.contentRect.width ?? 0;
      if (width <= 0 || Math.abs(width - lastWidth) < 1) return;
      lastWidth = width;
      cancelAnimationFrame(resizeFrame);
      resizeFrame = requestAnimationFrame(refit);
    });
    resizeObserver.observe(container);
  });

  createEffect(
    on(
      () => [props.options, props.series],
      () => {
        if (!chart || isDestroyed) return;
        chart
          .updateOptions(
            {
              ...props.options,
              series: props.series,
            },
            false,
            true
          )
          .catch((e) => {
            console.warn("ApexCharts update ignored during teardown:", e);
          });
      },
      { defer: true }
    )
  );

  onCleanup(() => {
    isDestroyed = true;
    cancelAnimationFrame(resizeFrame);
    resizeObserver?.disconnect();
    resizeObserver = undefined;
    if (chart) {
      chart.destroy();
      chart = undefined;
    }
  });

  return <div class={props.class} ref={container} />;
}
