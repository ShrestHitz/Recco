import { APIResponse, EnrichedMetadata, Item, ItemType, ShelfEntry } from '../types';

const API_BASE = '';

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<APIResponse<T>> {
  const headers = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  };

  try {
    const res = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    });

    const data: APIResponse<T> = await res.json();
    if (!res.ok && !data.error) {
      data.error = `HTTP Error ${res.status}: ${res.statusText}`;
    }
    return data;
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Network error';
    return {
      success: false,
      data: null as unknown as T,
      error: `Failed to connect to backend: ${message}`,
    };
  }
}

export const api = {
  checkHealth: async () => {
    return request<{ status: string; service: string; phase: string }>('/health');
  },

  enrichURL: async (rawURL: string) => {
    return request<EnrichedMetadata>(`/api/v1/items/enrich?url=${encodeURIComponent(rawURL)}`);
  },

  createItem: async (input: {
    type: ItemType;
    title: string;
    creator: string;
    external_url: string;
    cover_url?: string;
    source_id?: string;
  }) => {
    return request<Item>('/api/v1/items', {
      method: 'POST',
      body: JSON.stringify(input),
    });
  },

  listItems: async (type?: ItemType, limit = 50) => {
    const query = new URLSearchParams();
    if (type) query.append('type', type);
    query.append('limit', limit.toString());
    return request<Item[]>(`/api/v1/items?${query.toString()}`);
  },

  fetchQueue: async (userId: string) => {
    return request<ShelfEntry[]>(`/api/v1/shelf/queue?user_id=${encodeURIComponent(userId)}&limit=100`);
  },

  fetchShelf: async (userId: string) => {
    return request<ShelfEntry[]>(`/api/v1/shelf?user_id=${encodeURIComponent(userId)}&limit=100`);
  },

  addToQueue: async (userId: string, itemId: string, recommendedBy?: string) => {
    return request<ShelfEntry>('/api/v1/shelf/queue', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        item_id: itemId,
        recommended_by: recommendedBy,
      }),
    });
  },

  finishItem: async (userId: string, itemId: string, rating: number, take: string) => {
    return request<ShelfEntry>('/api/v1/shelf/finish', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        item_id: itemId,
        rating,
        take,
      }),
    });
  },

  dropItem: async (userId: string, itemId: string) => {
    return request<ShelfEntry>('/api/v1/shelf/drop', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        item_id: itemId,
      }),
    });
  },

  deleteEntry: async (id: string) => {
    return request<{ message: string }>(`/api/v1/shelf/${encodeURIComponent(id)}`, {
      method: 'DELETE',
    });
  },
};
