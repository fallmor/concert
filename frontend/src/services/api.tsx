import type { Concert, Artist, User, LoginResponse, Booking } from "../types";

const API_BASE_URL = "/api";

async function fetchAPI<T>(
  endpoint: string,
  options?: RequestInit,
): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    credentials: 'include',
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });

  if (!response.ok) {
    throw new Error(`API Error: ${response.statusText}`);
  }

  return response.json();
}

export const concertAPI = {
  getAll: async (): Promise<Concert[]> => {
    return fetchAPI<Concert[]>("/public/shows");
  },

  getById: async (id: number): Promise<Concert> => {
    return fetchAPI<Concert>(`/public/shows/${id}`);
  },

  create: async (
    concert: Omit<Concert, "id">,
    token: string,
  ): Promise<Concert> => {
    return fetchAPI<Concert>("/shows", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(concert),
    });
  },
};

export const artistAPI = {
  getAll: async (): Promise<Artist[]> => {
    return fetchAPI<Artist[]>("/public/artists");
  },

  getById: async (id: number): Promise<Artist> => {
    return fetchAPI<Artist>(`/public/artists/${id}`);
  },
};


export const authAPI = {
  register: async (
    username: string,
    email: string,
    password: string,
    firstName?: string,
    lastName?: string
  ): Promise<User> => {
    return fetchAPI<User>("/public/register", {
      method: "POST",
      body: JSON.stringify({
        username, email, password, firstName,
        lastName
      }),
    });
  },

  login: async (
    email: string,
    password: string,
  ): Promise<User> => {
    return fetchAPI<User>("/public/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },

  forget: async (
    email: string
  ): Promise<string> => {
    return fetchAPI<string>("/public/forget", {
      method: "Post",
      body: JSON.stringify({ email })
    });

  },
  getCurrentUser: async (token: string): Promise<User> => {
    return fetchAPI<User>('/me', {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
  },
};


export const bookingAPI = {
  create: async (showId: number, ticketCount: number): Promise<Booking> => {
    return fetchAPI<Booking>("/bookings", {
      method: 'POST',
      body: JSON.stringify({ showId, ticketCount })
    });
  },

  getMyBookings: async (): Promise<Booking[]> => {
    return fetchAPI<Booking[]>('/bookings');
  },

  cancel: async (bookingId: number): Promise<void> => {
    return fetchAPI<void>(`/bookings/${bookingId}`, {
      method: 'DELETE',
    });
  },
};