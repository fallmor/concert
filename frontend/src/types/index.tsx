export interface Concert {
  ID: number;
  title: string;
  artist: Artist;
  venue: string;
  date: string;
  time: string;
  price: number;
  total_seats: number;
  available_seats: number;
  description?: string;
  image_url?: string;
}

export interface Artist {
  id: number;
  name: string;
  genre: string;
  bio?: string;
  image_url?: string;
  album_url?: string;
}

export interface User {
  id: number;
  username: string;
  email: string;
  role: "user" | "moderator" | "admin";
}

export interface Booking {
  id: number;
  concertId: number;
  userId: number;
  numberOfTickets: number;
  totalPrice: number;
  bookingDate: string;
}
