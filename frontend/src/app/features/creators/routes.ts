import { Routes } from '@angular/router';
import { CreatorProfileComponent } from './creator-profile.component';

export const CREATORS_ROUTES: Routes = [
  { path: '', component: CreatorProfileComponent },
  { path: 'plans', loadComponent: () => import('./creator-plans.component').then(m => m.CreatorPlansComponent) }
];