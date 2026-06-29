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

  onMount(() => {
    chart = new ApexCharts(container, {
      ...props.options,
      series: props.series,
      chart: {
        ...(props.options.chart ?? {}),
        height: props.height ?? 240,
      },
    });
    
    if (!container) return;
    
    chart.render().catch((e) => {
      console.warn("ApexCharts render error ignored (likely unmounted):", e);
    });
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
    if (chart) {
      chart.destroy();
      chart = undefined;
    }
  });

  return <div class={props.class} ref={container} />;
}
