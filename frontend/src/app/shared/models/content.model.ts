export interface ContentItem {
    id: string;
    title: string;
    description: string;
    thumbnailUrl: string;
    createdAt: Date;
    tags?: string[];
    creatorId: string;
    // Añade más propiedades según tu API
  }
  
  export interface CreateContentResponse {
    id: string;
    message?: string;
  }
  
  export interface ContentSearchParams {
    query?: string;
    tags?: string[];
    creatorId?: string;
    page?: number;
    limit?: number;
  }