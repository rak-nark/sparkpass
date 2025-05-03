import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';

export const routes: Routes = [
  { 
    path: 'auth',
    loadChildren: () => import('./features/auth/routes').then(m => m.AUTH_ROUTES)
  },
  { 
    path: 'users',
    canActivate: [authGuard],
    loadChildren: () => import('./features/user/routes').then(m => m.USER_ROUTES)
  },
  {
    path: 'creators',
    canActivate: [authGuard],
    loadChildren: () => import('./features/creators/routes').then(m => m.CREATORS_ROUTES)
  },
  {
    path: 'content',
    canActivate: [authGuard],
    loadChildren: () => import('./features/content/routes').then(m => m.CONTENT_ROUTES)
  },
  // Cambia estas líneas:
  { path: '', redirectTo: 'auth/login', pathMatch: 'full' }, // Redirige a login por defecto
  { path: '**', redirectTo: 'auth/login' } // Manejo de rutas no encontradas
];