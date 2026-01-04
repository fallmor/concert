export interface Concert {
  ID: number;
  title: string;
  artist: Artist;
  venue: string;
  date: string;
  time: string;
  price: number;
  totalSeats: number;
  availableSeats: number;
  description?: string;
  imageUrl?: string;
}

export interface User {
  ID: number;
  email: string;
  username: string;
  firstName?: string;
  lastName?: string;
  role: "user" | "moderator" | "admin";
  createdAt?: string;
}

export interface Artist {
  ID: number;
  name: string;
  genre: string;
  bio?: string;
  imageUrl?: string;
}


export interface Booking {
  ID: number;
  userId: number;
  showId: number;
  show?: Concert;
  ticketCount: number;
  totalPrice: number;
  status: 'confirmed' | 'cancelled';
  createdAt: string;
  bookingDate: string;
}
export interface Artist {
  id: number;
  name: string;
  genre: string;
  bio?: string;
  shows: Concert[]
  imageUrl?: string;

}

export interface LoginResponse {
  user: User;
  token?: string;
}

export interface RegisterInput {
  username: string;
  email: string;
  password: string;
}

export interface LoginInput {
  email: string;
  password: string;
}
