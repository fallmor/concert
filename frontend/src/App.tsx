import React, { useState } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Header from "./components/header";
import Footer from "./components/footer";
import HomePage from "./pages/HomePage";
import ConcertsPage from "./pages/ConcertPages";
import ConcertDetailPage from "./pages/ConcertDetailPage";

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [userRole, setUserRole] = useState<"user" | "moderator" | "admin">(
    "user",
  );

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
        <Header isLoggedIn={isLoggedIn} userRole={userRole} />

        <main style={{ flex: 1, padding: "40px 20px" }}>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/concerts/:id" element={<ConcertDetailPage />} /> 
            <Route path="/concerts" element={<ConcertsPage />} />
          </Routes>
        </main>

        <Footer />
      </div>
    </Router>
  );
}

export default App;
