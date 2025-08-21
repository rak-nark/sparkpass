import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { catchError, tap } from 'rxjs/operators';
import { of, throwError } from 'rxjs';
import { StorageService } from '../../services/storage.service'; // Usa el servicio seguro que creamos

@Injectable({ providedIn: 'root' })
export class AuthService {
  private http = inject(HttpClient);
  private router = inject(Router);
  private storage = inject(StorageService);
  private readonly apiUrl = 'http://localhost:8080/api'; // Ajusta tu URL de backend

  login(credentials: { email: string; password: string }) {
    return this.http.post<{ token: string }>(`${this.apiUrl}/login`, credentials).pipe(
      tap(response => {
        this.storage.set('auth_token', response.token);
        this.router.navigate(['/content']);
      }),
      catchError(error => {
        console.error('Login error:', error);
        return of(null);
      })
    );
  }

  register(userData: { email: string; password: string }) {
    return this.http.post(`${this.apiUrl}/register`, userData).pipe(
      catchError(error => {
        // Propaga el error tal cual para que el componente lo reciba completo
        return throwError(() => error);
      })
    );
  }

  getToken(): string | null {
    return this.storage.get('auth_token');
  }

  getUser(): any {
    return this.storage.get('user');  
  }

  isAuthenticated(): boolean {
    return !!this.getToken();
  }

  logout() {
    this.storage.remove('auth_token');
    this.router.navigate(['/auth/login']);
  }
}