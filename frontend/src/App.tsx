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
import MyBookingsPage from "./pages/BookingPage";
import CreateBookingPage from "./pages/CreateBookingPage";
import AdminDashboard from "./pages/AdminDashboard";
import AdminShows from "./pages/AdminShows";
import AdminBookings from "./pages/AdminBookings";
import AdminArtists from "./pages/AdminArtists";
import ForgetPass from "./pages/ForgetPass";

function App() {

  return (
    <Router>
      {/* Using Tailwind for layout instead of inline styles */}
      <div className="min-h-screen flex flex-col bg-gray-50">
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
            <Route path="/forget" element={<ForgetPass />} />
            <Route path="/my-bookings" element={<MyBookingsPage />} />
            <Route path="/book/:id" element={<CreateBookingPage />} />
            <Route path="/admin" element={<AdminDashboard />} />
            <Route path="/admin/shows" element={<AdminShows />} />
            <Route path="/admin/bookings" element={<AdminBookings />} />
            <Route path="/admin/artists" element={<AdminArtists />} />
          </Routes>
        </main>

        <Footer />
      </div>
    </Router>
  );
}

export default App;