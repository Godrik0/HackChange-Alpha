import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', loadChildren: () => import('./features/clients/clients.module').then(m => m.ClientsModule) },
  { path: '', loadChildren: () => import('./features/client-metrics/client-metrics.module').then(m => m.ClientMetricsModule) },
  { path: '**', redirectTo: '' }
];