(function () {
  function isDark() {
    return document.documentElement.classList.contains('dark');
  }

  function initSpendingChart() {
    const canvas = document.getElementById('spending-trend-chart');
    if (!canvas || typeof Chart === 'undefined') return;

    const existing = Chart.getChart('spending-trend-chart');
    if (existing) existing.destroy();

    let data = [];
    try {
      data = JSON.parse(canvas.dataset.chart || '[]');
    } catch {
      return;
    }
    if (!data.length) return;

    const labels = data.map((d) => d.month);
    const values = data.map((d) => d.total);

    new Chart(canvas, {
      type: 'bar',
      data: {
        labels,
        datasets: [{
          label: 'Spending',
          data: values,
          backgroundColor: isDark() ? 'rgba(96, 165, 250, 0.6)' : 'rgba(37, 99, 235, 0.6)',
          borderColor: isDark() ? 'rgb(96, 165, 250)' : 'rgb(37, 99, 235)',
          borderWidth: 1,
          borderRadius: 4,
        }],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: {
          y: {
            beginAtZero: true,
            ticks: { color: isDark() ? '#9ca3af' : '#6b7280' },
            grid: { color: isDark() ? 'rgba(75,85,99,0.4)' : 'rgba(229,231,235,0.8)' },
          },
          x: {
            ticks: { color: isDark() ? '#9ca3af' : '#6b7280' },
            grid: { display: false },
          },
        },
      },
    });
  }

  function initDonutChart(canvasId) {
    const canvas = document.getElementById(canvasId);
    if (!canvas || typeof Chart === 'undefined') return;

    const existing = Chart.getChart(canvasId);
    if (existing) existing.destroy();

    let data = [];
    try {
      data = JSON.parse(canvas.dataset.chart || '[]');
    } catch {
      return;
    }
    if (!data.length) return;

    new Chart(canvas, {
      type: 'doughnut',
      data: {
        labels: data.map((d) => d.label),
        datasets: [{
          data: data.map((d) => d.value),
          backgroundColor: data.map((d) => d.color),
          borderWidth: 0,
        }],
      },
      options: {
        responsive: true,
        maintainAspectRatio: true,
        cutout: '55%',
        plugins: {
          legend: { display: false },
          tooltip: {
            callbacks: {
              label(ctx) {
                const v = ctx.parsed || 0;
                return ` $${v.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
              },
            },
          },
        },
      },
    });
  }

  function initAllCharts() {
    initSpendingChart();
    initDonutChart('spending-by-tag-chart');
    initDonutChart('income-by-tag-chart');
  }

  document.addEventListener('DOMContentLoaded', initAllCharts);
  document.body.addEventListener('htmx:afterSwap', (evt) => {
    const target = evt.detail.target;
    if (!target) return;
    if (
      target.id === 'dashboard-grid' ||
      target.querySelector?.('#spending-trend-chart') ||
      target.querySelector?.('#spending-by-tag-chart') ||
      target.querySelector?.('#income-by-tag-chart')
    ) {
      // Small delay to ensure the browser has completed layout/reflow
      setTimeout(initAllCharts, 20);
    }
  });
})();
