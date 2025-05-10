// notification.service.ts
import { Injectable, inject } from '@angular/core';
import { ApiService } from '../shared/api/api.service';
import { interval, Observable } from 'rxjs';
import { switchMap, startWith } from 'rxjs/operators';

interface Notification {
    id: string;
    message: string;
    read: boolean;
    createdAt: string;
}

@Injectable({ providedIn: 'root' })
export class NotificationService {
    private api = inject(ApiService);

    getNotifications(): Observable<Notification[]> {
        return this.api.get<Notification[]>('/notifications');
    }

    markAsRead(): Observable<void> {
        return this.api.post<void>('/notifications/mark-read', {});
    }

    // Opcional: Polling cada 30 segundos
    getLiveNotifications(pollInterval: number = 30000): Observable<Notification[]> {
        return interval(pollInterval).pipe(
            startWith(0),
            switchMap(() => this.getNotifications())
        );
    }
}