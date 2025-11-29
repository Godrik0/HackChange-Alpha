import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', loadChildren: () => import('./features/clients/clients.module').then(m => m.ClientsModule) },
  { path: '**', redirectTo: '' }
];
