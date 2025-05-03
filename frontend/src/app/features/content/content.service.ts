import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../shared/api/api.service';
import { ContentItem, CreateContentResponse } from '../../shared/models/content.model';

@Injectable({ providedIn: 'root' })
export class ContentService {
  private api = inject(ApiService);

  getContent(): Observable<ContentItem[]> {
    return this.api.get<ContentItem[]>('/content');
  }

  getContentById(id: string): Observable<ContentItem> {
    return this.api.get<ContentItem>(`/content/${id}`);
  }

  createContent(contentData: Omit<ContentItem, 'id'>): Observable<CreateContentResponse> {
    return this.api.post<CreateContentResponse>('/content', contentData);
  }
}