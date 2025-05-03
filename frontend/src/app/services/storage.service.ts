import { Injectable, Inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';

@Injectable({ providedIn: 'root' })
export class StorageService {
  constructor(@Inject(PLATFORM_ID) private platformId: Object) {}

  get(key: string): string | null {
    if (this.isBrowser()) {
      return localStorage.getItem(key);
    }
    return null;
  }

  set(key: string, value: string): void {
    if (this.isBrowser()) {
      localStorage.setItem(key, value);
    }
  }

  remove(key: string): void {
    if (this.isBrowser()) {
      localStorage.removeItem(key);
    }
  }

  private isBrowser(): boolean {
    return isPlatformBrowser(this.platformId);
  }
}