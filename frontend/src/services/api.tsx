import type { Concert, Artist, User, Booking } from "../types";

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
    // Try to read the error message from the response body
    let errorMessage = response.statusText;
    try {
      const errorData = await response.text();
      if (errorData) {
        // Try to parse as JSON first
        try {
          const jsonError = JSON.parse(errorData);
          errorMessage = jsonError.error || jsonError.message || errorData;
        } catch {
          // If not JSON, use the text as is
          errorMessage = errorData;
        }
      }
    } catch {
      // If we can't read the body, use statusText
      errorMessage = response.statusText;
    }
    throw new Error(errorMessage || `API Error: ${response.statusText}`);
  }


  const contentType = response.headers.get("content-type");
  const text = await response.text();

  if (!text || text.trim() === '') {
    return '' as T;
  }

  if (contentType && contentType.includes("application/json")) {
    try {
      return JSON.parse(text) as T;
    } catch (e) {
      return text as T;
    }
  }

  return text as T;
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