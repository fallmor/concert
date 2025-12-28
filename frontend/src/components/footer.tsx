import React from "react";

const Footer: React.FC = () => {
  return (
    <footer
      style={{
        backgroundColor: "#1a1a2e",
        color: "white",
        padding: "30px 40px",
        marginTop: "60px",
        textAlign: "center",
      }}
    >
      <p style={{ margin: 0 }}>
        © {new Date().getFullYear()} Concert Booking System. All rights
        reserved.
      </p>
      <p style={{ margin: "10px 0 0 0", fontSize: "14px", color: "#aaa" }}>
        Built with React + TypeScript + Go
      </p>
    </footer>
  );
};

export default Footer;
