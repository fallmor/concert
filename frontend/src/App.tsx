import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Header from "./components/header";
import Footer from "./components/footer";
import HomePage from "./pages/HomePage";
import ConcertsPage from "./pages/ConcertPages";
import ConcertDetailPage from "./pages/ConcertDetailPage";
import ArtistsPage from "./pages/ArtistsPage";
import ArtistDetailPage from "./pages/ArtsitDetailPage";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import BookingPage from "./pages/BookingPage";
import MyBookingsPage from "./pages/BookingPage";
import AdminDashboard from "./pages/AdminDashboard";
import AdminShows from "./pages/AdminShows";

function App() {

  return (
    <Router>
      <div
        style={{
          minHeight: "100vh",
          display: "flex",
          flexDirection: "column",
          backgroundColor: "#f5f5f5",
        }}
      >
        <Header />

        <main style={{ flex: 1, padding: "40px 20px" }}>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/concerts" element={<ConcertsPage />} />
            <Route path="/concerts/:id" element={<ConcertDetailPage />} />
            <Route path="/artists" element={<ArtistsPage />} />
            <Route path="/artists/:id" element={<ArtistDetailPage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route path="/my-bookings" element={<MyBookingsPage />} />
            <Route path="/book/:id" element={<BookingPage />} />
            <Route path="/admin" element={<AdminDashboard />} />
            <Route path="/admin/shows" element={<AdminShows />} />
          </Routes>
        </main>

        <Footer />
      </div>
    </Router>
  );
}

export default App;