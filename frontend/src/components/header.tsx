import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Header: React.FC = () => {
  const { user, logout, isAuthenticated } = useAuth();

  return (
    <header style={{
      backgroundColor: '#1a1a2e',
      color: 'white',
      padding: '20px 40px',
      boxShadow: '0 2px 10px rgba(0,0,0,0.1)'
    }}>
      <div style={{
        maxWidth: '1200px',
        margin: '0 auto',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>

        <Link to="/" style={{ textDecoration: 'none', color: 'white' }}>
          <h1 style={{ margin: 0, fontSize: '28px' }}>
            🎵 Concert Booking
          </h1>
        </Link>

        <nav style={{ display: 'flex', gap: '30px', alignItems: 'center' }}>
          <Link to="/concerts" style={navLinkStyle}>
            Concerts
          </Link>
          <Link to="/artists" style={navLinkStyle}>
            Artists
          </Link>

          {isAuthenticated ? (
            <>
              <span style={{ color: '#aaa', fontSize: '14px' }}>
                Welcome, {user?.username}!
              </span>
              <Link to="/my-bookings" style={navLinkStyle}>
                My Bookings
              </Link>
              {(user?.role === 'admin' || user?.role === 'moderator') && (
                <Link to="/admin" style={navLinkStyle}>
                  Admin Panel
                </Link>
              )}
              <button
                onClick={logout}
                style={{
                  backgroundColor: '#dc3545',
                  color: 'white',
                  padding: '8px 16px',
                  border: 'none',
                  borderRadius: '5px',
                  cursor: 'pointer',
                  fontSize: '16px',
                  fontWeight: '500'
                }}
              >
                Logout
              </button>
            </>
          ) : (
            <>
              <Link to="/login" style={navLinkStyle}>
                Login
              </Link>
                <Link to="/register" style={{
                  ...navLinkStyle,
                  backgroundColor: '#007bff',
                  padding: '8px 16px',
                  borderRadius: '5px'
                }}>
                Register
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
};

const navLinkStyle: React.CSSProperties = {
  color: 'white',
  textDecoration: 'none',
  fontSize: '16px',
  fontWeight: '500',
  transition: 'color 0.3s'
};

export default Header;