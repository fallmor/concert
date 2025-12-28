import type { Concert, Artist, User } from "../types";

// Your Go backend URL
const API_BASE_URL = "http://localhost:8080";

// Helper function for API calls
async function fetchAPI<T>(
  endpoint: string,
  options?: RequestInit,
): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
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

// Concert/Show APIs
export const concertAPI = {
  // Get all concerts
  getAll: async (): Promise<Concert[]> => {
    return fetchAPI<Concert[]>("/shows");
  },

  // Get single concert
  getById: async (id: number): Promise<Concert> => {
    return fetchAPI<Concert>(`/shows/${id}`);
  },

  // Create concert (admin only)
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

// Artist APIs
export const artistAPI = {
  getAll: async (): Promise<Artist[]> => {
    return fetchAPI<Artist[]>("/artists");
  },

  getById: async (id: number): Promise<Artist> => {
    return fetchAPI<Artist>(`/artists/${id}`);
  },
};

// User/Auth APIs
export const authAPI = {
  register: async (
    username: string,
    email: string,
    password: string,
  ): Promise<User> => {
    return fetchAPI<User>("/users/register", {
      method: "POST",
      body: JSON.stringify({ username, email, password }),
    });
  },

  login: async (
    email: string,
    password: string,
  ): Promise<{ user: User; token: string }> => {
    return fetchAPI<{ user: User; token: string }>("/users/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },
};

// Booking APIs
export const bookingAPI = {
  create: async (
    concertId: number,
    numberOfTickets: number,
    token: string,
  ): Promise<any> => {
    return fetchAPI<any>("/bookings", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        show_id: concertId,
        tickets: numberOfTickets,
      }),
    });
  },
};
