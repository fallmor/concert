import React from "react";
import { Link } from "react-router-dom";

const HomePage: React.FC = () => {
  return (
    <div style={{ maxWidth: "1200px", margin: "0 auto", textAlign: "center" }}>
      <div
        style={{
          backgroundColor: "#1a1a2e",
          color: "white",
          padding: "80px 40px",
          borderRadius: "15px",
          marginBottom: "60px",
        }}
      >
        <h1 style={{ fontSize: "48px", margin: "0 0 20px 0" }}>
          🎸 Welcome to Concert Booking
        </h1>
        <p style={{ fontSize: "20px", margin: "0 0 40px 0", color: "#ddd" }}>
          Discover amazing live music events and book your tickets now
        </p>
        <Link
          to="/concerts"
          style={{
            backgroundColor: "#007bff",
            color: "white",
            padding: "15px 40px",
            textDecoration: "none",
            borderRadius: "8px",
            fontSize: "18px",
            fontWeight: "bold",
            display: "inline-block",
          }}
        >
          Browse Concerts
        </Link>
      </div>

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fit, minmax(250px, 1fr))",
          gap: "30px",
          marginTop: "60px",
        }}
      >
        <FeatureCard
          icon="🎤"
          title="Top Artists"
          description="See performances from the biggest names in music"
        />
        <FeatureCard
          icon="🎫"
          title="Easy Booking"
          description="Book tickets quickly and securely online"
        />
        <FeatureCard
          icon="📍"
          title="Great Venues"
          description="Concerts at the best locations around the world"
        />
      </div>
    </div>
  );
};

interface FeatureCardProps {
  icon: string;
  title: string;
  description: string;
}

const FeatureCard: React.FC<FeatureCardProps> = ({
  icon,
  title,
  description,
}) => {
  return (
    <div
      style={{
        backgroundColor: "white",
        padding: "40px",
        borderRadius: "10px",
        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
      }}
    >
      <div style={{ fontSize: "48px", marginBottom: "20px" }}>{icon}</div>
      <h3 style={{ margin: "0 0 15px 0", color: "#333" }}>{title}</h3>
      <p style={{ margin: 0, color: "#666", lineHeight: "1.6" }}>
        {description}
      </p>
    </div>
  );
};

export default HomePage;
