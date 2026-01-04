import type { Concert, Artist } from '../types';

const API_BASE_URL = '/api';

async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `API Error: ${response.statusText}`);
  }

  return response.json();
}

export interface AdminStats {
  totalShows: number;
  totalArtists: number;
  totalBookings: number;
  totalRevenue: number;
  totalUsers: number;
}

export interface AdminBooking {
  id: number;
  userId: number;
  username: string;
  userEmail: string;
  showId: number;
  showTitle: string;
  artistName: string;
  ticketCount: number;
  totalPrice: number;
  status: string;
  createdAt: string;
}

export interface CreateShowInput {
  title: string;
  date: string;
  time: string;
  artistId: number;
  venue: string;
  price: number;
  totalSeats: number;
  description?: string;
  imageUrl?: string;
}

export interface CreateArtistInput {
  name: string;
  genre: string;
  bio?: string;
  imageUrl?: string;
}

export const adminAPI = {

  getStats: async (): Promise<AdminStats> => {
    return fetchAPI<AdminStats>('/admin/stats');
  },

  listShows: async (): Promise<Concert[]> => {
    return fetchAPI<Concert[]>('/admin/shows');
  },

  createShow: async (show: CreateShowInput): Promise<Concert> => {
    return fetchAPI<Concert>('/admin/shows', {
      method: 'POST',
      body: JSON.stringify(show),
    });
  },

  updateShow: async (id: number, show: Partial<CreateShowInput>): Promise<Concert> => {
    return fetchAPI<Concert>(`/admin/shows/${id}`, {
      method: 'PUT',
      body: JSON.stringify(show),
    });
  },

  deleteShow: async (id: number): Promise<void> => {
    return fetchAPI<void>(`/admin/shows/${id}`, {
      method: 'DELETE',
    });
  },


  listArtists: async (): Promise<Artist[]> => {
    return fetchAPI<Artist[]>('/admin/artists');
  },

  createArtist: async (artist: CreateArtistInput): Promise<Artist> => {
    return fetchAPI<Artist>('/admin/artists', {
      method: 'POST',
      body: JSON.stringify(artist),
    });
  },

  updateArtist: async (id: number, artist: Partial<CreateArtistInput>): Promise<Artist> => {
    return fetchAPI<Artist>(`/admin/artists/${id}`, {
      method: 'PUT',
      body: JSON.stringify(artist),
    });
  },

  deleteArtist: async (id: number): Promise<void> => {
    return fetchAPI<void>(`/admin/artists/${id}`, {
      method: 'DELETE',
    });
  },

  listBookings: async (): Promise<AdminBooking[]> => {
    return fetchAPI<AdminBooking[]>('/admin/bookings');
  },
};