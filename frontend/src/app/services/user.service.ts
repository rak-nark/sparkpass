// user.service.ts
import { Injectable, inject } from '@angular/core';
import { ApiService } from '../shared/api/api.service';
import { Observable } from 'rxjs';
import { StorageService } from './storage.service';

interface UserProfile {
    id: string;
    email: string;
    name: string;
    isCreator: boolean;
    avatar?: string;
}

@Injectable({ providedIn: 'root' })
export class UserService {
    private api = inject(ApiService);
    private storage = inject(StorageService);

    getCurrentUser(): Observable<UserProfile> {
        return this.api.get<UserProfile>('/users/me');
    }

    updateUserProfile(profileData: Partial<UserProfile>): Observable<UserProfile> {
        return this.api.put<UserProfile>('/users/me', profileData);
    }

    getCachedUser(): UserProfile | null {
        return this.storage.get('user_profile') ? JSON.parse(this.storage.get('user_profile')!) : null;
    }
}