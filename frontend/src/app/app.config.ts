import { ApplicationConfig, PLATFORM_ID } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptors, withFetch } from '@angular/common/http';
import { routes } from './app.routes';
import { authInterceptor } from './core/interceptors/auth.interceptor';
import { provideClientHydration } from '@angular/platform-browser';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(
      withFetch(), // 👈 Añade esto para el warning de fetch
      withInterceptors([authInterceptor])
    ),
    provideClientHydration(), // 👈 Opcional pero recomendado para SSR
    // Solución para localStorage en SSR:
    {
      provide: 'PLATFORM_ID',
      useValue: PLATFORM_ID
    }
  ]
};