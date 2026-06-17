import ApexCharts, { type ApexOptions } from "apexcharts";
import { onCleanup, onMount, createEffect } from "solid-js";

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

  onMount(() => {
    chart = new ApexCharts(container, {
      ...props.options,
      series: props.series,
      chart: {
        ...(props.options.chart ?? {}),
        height: props.height ?? 240,
      },
    });
    chart.render();
  });

  createEffect(() => {
    if (!chart) return;
    chart.updateOptions(
      {
        ...props.options,
        series: props.series,
      },
      false,
      true
    );
  });

  onCleanup(() => {
    chart?.destroy();
  });

  return <div class={props.class} ref={container} />;
}
