function dashboardCustomize() {
  return {
    editMode: false,
    layout: { widgets: [] },
    sortable: null,

    initLayout() {
      this.readLayoutFromDOM();
      document.body.addEventListener('htmx:afterSwap', (evt) => {
        if (evt.detail.target?.id === 'dashboard-grid') {
          this.readLayoutFromDOM();
          if (this.editMode) {
            this.$nextTick(() => {
              this.initSortable();
              this.bindToggleButtons();
            });
          }
        }
      });
    },

    readLayoutFromDOM() {
      const grid = document.getElementById('dashboard-grid');
      if (!grid) return;
      const widgets = [];
      grid.querySelectorAll('.dashboard-widget').forEach((el, index) => {
        widgets.push({
          id: el.dataset.widgetId,
          visible: el.dataset.visible === 'true',
          order: index,
        });
      });
      if (widgets.length > 0) {
        this.layout = { widgets };
      }
    },

    enterEditMode() {
      this.editMode = true;
      htmx.ajax('GET', '/dashboard/widgets?edit=1', {
        target: '#dashboard-grid',
        swap: 'outerHTML',
      }).then(() => {
        this.$nextTick(() => {
          this.initSortable();
          this.bindToggleButtons();
        });
      });
    },

    cancelEdit() {
      this.editMode = false;
      this.destroySortable();
      htmx.ajax('GET', '/dashboard/widgets', {
        target: '#dashboard-grid',
        swap: 'outerHTML',
      });
    },

    saveLayout() {
      this.readLayoutFromDOM();
      const self = this;
      htmx.ajax('POST', '/dashboard/layout', {
        values: { layout: JSON.stringify(this.layout) },
        target: '#dashboard-grid',
        swap: 'outerHTML',
      }).then(() => {
        self.editMode = false;
        self.destroySortable();
      });
    },

    initSortable() {
      this.destroySortable();
      const grid = document.getElementById('dashboard-grid');
      if (!grid || typeof Sortable === 'undefined') return;
      this.sortable = Sortable.create(grid, {
        handle: '.dashboard-widget-toolbar, .dashboard-widget-toolbar *',
        animation: 150,
        draggable: '.dashboard-widget',
        filter: '.dashboard-toggle-visibility',
        preventOnFilter: true,
        onEnd: () => this.readLayoutFromDOM(),
      });
    },

    destroySortable() {
      if (this.sortable) {
        this.sortable.destroy();
        this.sortable = null;
      }
    },

    bindToggleButtons() {
      const grid = document.getElementById('dashboard-grid');
      if (!grid) return;
      grid.querySelectorAll('.dashboard-toggle-visibility').forEach((btn) => {
        btn.onclick = (e) => {
          e.preventDefault();
          const id = btn.dataset.widgetId;
          let w = this.layout.widgets.find((x) => x.id === id);
          if (!w) return;
          w.visible = !w.visible;
          const el = grid.querySelector(`.dashboard-widget[data-widget-id="${id}"]`);
          if (!el) return;
          el.dataset.visible = w.visible ? 'true' : 'false';
          
          const visibleClasses = 'dashboard-toggle-visibility p-1.5 rounded-full transition-all duration-200 text-blue-600 bg-blue-50 hover:bg-blue-100 dark:text-blue-400 dark:bg-blue-900/30 dark:hover:bg-blue-900/50';
          const hiddenClasses = 'dashboard-toggle-visibility p-1.5 rounded-full transition-all duration-200 text-gray-400 bg-gray-100 hover:bg-gray-200 dark:text-gray-500 dark:bg-gray-800 dark:hover:bg-gray-700';
          
          btn.className = w.visible ? visibleClasses : hiddenClasses;
          btn.title = w.visible ? 'Hide widget' : 'Show widget';
          btn.innerHTML = w.visible 
            ? '<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.542-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l18 18"></path></svg>'
            : '<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg>';
          
          const preview = el.querySelector('[data-widget-preview]');
          if (preview) {
            preview.classList.toggle('opacity-40', !w.visible);
            preview.classList.toggle('pointer-events-none', !w.visible);
            preview.classList.toggle('select-none', !w.visible);
          }
        };
      });
    },
  };
}
