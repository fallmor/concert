import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import type { Booking } from '../types';
import { bookingAPI } from '../services/api';

const MyBookingsPage: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated, isLoading } = useAuth();
  
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState('');

  useEffect(() => {
    if (isLoading) return;
    
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    if (location.state?.message) {
      setSuccessMessage(location.state.message);
      setTimeout(() => setSuccessMessage(''), 5000);
    }

    const fetchBookings = async () => {
      try {
        const data = await bookingAPI.getMyBookings();
        setBookings(data);
      } catch (err: any) {
        setError(err.message || 'Failed to load bookings');
      } finally {
        setLoading(false);
      }
    };

    fetchBookings();
  }, [isAuthenticated, isLoading, navigate, location.state]);

  const handleCancel = async (bookingId: number) => {
    if (!confirm('Are you sure you want to cancel this booking?')) {
      return;
    }

    try {
      await bookingAPI.cancel(bookingId);
      
      setBookings(prev => 
        prev.map(b => 
          b.id === bookingId 
            ? { ...b, status: 'cancelled' as const }
            : b
        )
      );
      
      setSuccessMessage('Booking cancelled successfully');
      setTimeout(() => setSuccessMessage(''), 5000);
    } catch (err: any) {
      alert(err.message || 'Failed to cancel booking');
    }
  };

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <h2>Loading...</h2>
      </div>
    );
  }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <h2>Loading your bookings...</h2>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '1000px', margin: '0 auto' }}>
      <h1 style={{ fontSize: '36px', marginBottom: '30px' }}>
        My Bookings
      </h1>

      {successMessage && (
        <div style={{
          backgroundColor: '#d4edda',
          color: '#155724',
          padding: '15px',
          borderRadius: '8px',
          marginBottom: '20px',
          border: '1px solid #c3e6cb'
        }}>
          {successMessage}
        </div>
      )}

      {error && (
        <div style={{
          backgroundColor: '#f8d7da',
          color: '#721c24',
          padding: '15px',
          borderRadius: '8px',
          marginBottom: '20px'
        }}>
          {error}
        </div>
      )}

      {bookings.length === 0 ? (
        <div style={{
          textAlign: 'center',
          padding: '60px',
          backgroundColor: 'white',
          borderRadius: '15px',
          boxShadow: '0 2px 10px rgba(0,0,0,0.1)'
        }}>
          <p style={{ fontSize: '18px', color: '#666', marginBottom: '20px' }}>
            You haven't booked any concerts yet.
          </p>
          <Link
            to="/concerts"
            style={{
              display: 'inline-block',
              padding: '12px 30px',
              backgroundColor: '#007bff',
              color: 'white',
              textDecoration: 'none',
              borderRadius: '8px',
              fontWeight: 'bold'
            }}
          >
            Browse Concerts
          </Link>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          {bookings.map(booking => (
            <BookingCard
              key={booking.id}
              booking={booking}
              onCancel={handleCancel}
            />
          ))}
        </div>
      )}
    </div>
  );
};

interface BookingCardProps {
  booking: Booking;
  onCancel: (bookingId: number) => void;
}

const BookingCard: React.FC<BookingCardProps> = ({ booking, onCancel }) => {
  const isCancelled = booking.status === 'cancelled';

  return (
    <div style={{
      backgroundColor: 'white',
      borderRadius: '15px',
      padding: '25px',
      boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'center',
      opacity: isCancelled ? 0.6 : 1,
      border: isCancelled ? '2px solid #dc3545' : 'none'
    }}>
      <div style={{ flex: 1 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '10px' }}>
          <h3 style={{ margin: 0, fontSize: '22px' }}>
            {booking.show?.title}
          </h3>
          {isCancelled && (
            <span style={{
              backgroundColor: '#dc3545',
              color: 'white',
              padding: '4px 12px',
              borderRadius: '12px',
              fontSize: '12px',
              fontWeight: 'bold'
            }}>
              CANCELLED
            </span>
          )}
        </div>

        <div style={{ color: '#666', fontSize: '14px' }}>
          <p style={{ margin: '5px 0' }}>
            🎟️ {booking.ticketCount} ticket{booking.ticketCount !== 1 ? 's' : ''}
          </p>
          <p style={{ margin: '5px 0' }}>
            📅 Booked on: {new Date(booking.bookingDate).toLocaleDateString('en-US', {
              year: 'numeric',
              month: 'long',
              day: 'numeric'
            })}
          </p>
        </div>
      </div>

      <div style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'flex-end',
        gap: '15px'
      }}>
        <div style={{
          fontSize: '28px',
          fontWeight: 'bold',
          color: isCancelled ? '#999' : '#007bff'
        }}>
          ${booking.totalPrice.toFixed(2)}
        </div>

        {!isCancelled && (
          <button
            onClick={() => onCancel(booking.id)}
            style={{
              padding: '8px 20px',
              backgroundColor: '#dc3545',
              color: 'white',
              border: 'none',
              borderRadius: '5px',
              cursor: 'pointer',
              fontSize: '14px',
              fontWeight: 'bold',
              transition: 'background-color 0.3s'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = '#c82333';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = '#dc3545';
            }}
          >
            Cancel Booking
          </button>
        )}
      </div>
    </div>
  );
};

export default MyBookingsPage;