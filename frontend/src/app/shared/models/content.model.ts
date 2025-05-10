export interface ContentItem {
  ID: string;
  title: string;
  description: string;
  thumbnailUrl?: string; // opcional, si no tienes miniatura
  createdAt: Date;
  tags?: string[];
  creatorId: string;
  isLocked: boolean;
  price: number;
  slug: string;
  contentType: 'video' | 'image' | 'podcast' | 'document';
  s3Key: string;
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