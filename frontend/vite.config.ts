import { defineConfig } from 'vite';
import angular from '@angular-devkit/build-angular/plugins/vite';
import vitePluginImp from 'vite-plugin-imp';

export default defineConfig({
  plugins: [
    angular(),
    vitePluginImp({
      libList: [
        {
          libName: 'ng-zorro-antd',
          style: (name) => `ng-zorro-antd/esm/${name}/style/index.less`
        }
      ]
    })
  ],
  css: {
    preprocessorOptions: {
      less: {
        modifyVars: {
          // КРАСНАЯ ТЕМА
          'primary-color': '#dc2626',
          'success-color': '#16a34a',
          'warning-color': '#d97706',
          'error-color': '#dc2626',
          'info-color': '#dc2626',
          'body-background': '#ffffff',
          'component-background': '#ffffff',
          'layout-body-background': '#ffffff',
          'layout-header-background': '#dc2626',
          'text-color': '#1f2937',
          'text-color-secondary': '#6b7280',
          'heading-color': '#111827',
          'border-color-base': '#e5e7eb',
          'border-color-split': '#f3f4f6',
          'btn-primary-bg': '#dc2626',
          'btn-default-bg': '#fef2f2',
          'btn-default-border': '#fecaca',
          'card-head-color': '#111827',
          'card-head-background': '#fef2f2',
          'card-background': '#ffffff',
          'card-shadow': '0 1px 3px 0 rgba(0, 0, 0, 0.1)',
          'menu-item-active-bg': '#fef2f2',
          'menu-highlight-color': '#dc2626',
          'input-bg': '#ffffff',
          'input-border-color': '#e5e7eb',
          'input-hover-border-color': '#dc2626',
          'tabs-hover-color': '#dc2626',
          'tabs-active-color': '#dc2626',
          'tabs-ink-bar-color': '#dc2626',
          'tabs-highlight-color': '#dc2626',
          'table-header-bg': '#fef2f2',
          'table-header-color': '#111827',
          'table-row-hover-bg': '#fef2f2'
        },
        javascriptEnabled: true
      }
    }
  }
});
