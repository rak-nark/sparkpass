import { Routes } from '@angular/router';
import { ContentListComponent } from './content-list/content-list.component';

export const CONTENT_ROUTES: Routes = [
  { 
    path: '', 
    component: ContentListComponent,
    title: 'Lista de Contenido'  // Añade título de página
  },
  { 
    path: ':id', 
    loadComponent: () => import('./content-detail/content-detail.component').then(m => m.ContentDetailComponent),
    title: 'Detalle de Contenido'
  },
  { 
    path: 'search', 
    loadComponent: () => import('./content-search/content-search.component').then(m => m.ContentSearchComponent),
    title: 'Buscar Contenido'
  }
];