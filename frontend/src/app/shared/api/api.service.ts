import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly baseUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  get<T>(endpoint: string) {
    return this.http.get<T>(`${this.baseUrl}${endpoint}`);
  }

  post<T>(endpoint: string, body: any) {
    return this.http.post<T>(`${this.baseUrl}${endpoint}`, body);
  }

  // Añade put, delete, etc según necesidad
}