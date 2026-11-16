export type ItemType = 'movie' | 'youtube' | 'book' | 'anime' | 'podcast';

export interface Item {
  id: string;
  type: ItemType;
  title: string;
  creator: string;
  cover_url: string;
  external_url: string;
  source_id: string;
  metadata?: Record<string, unknown>;
  created_at: string;
  updated_at?: string;
}

export type ShelfStatus = 'queued' | 'finished' | 'dropped';

export interface EnrichedMetadata {
  type: ItemType;
  title: string;
  creator: string;
  cover_url: string;
  external_url: string;
  source_id: string;
  description?: string;
  metadata?: Record<string, unknown>;
}

export interface ShelfEntry {
  id: string;
  user_id: string;
  item_id: string;
  item?: Item;
  recommended_by?: string;
  status: ShelfStatus;
  rating?: number;
  take?: string;
  finished_at?: string;
  created_at: string;
  updated_at?: string;
}

export interface APIResponse<T> {
  success: boolean;
  data: T;
  error?: string;
  meta?: {
    count?: number;
    queue_zero?: boolean;
    limit?: number;
    offset?: number;
  };
}
