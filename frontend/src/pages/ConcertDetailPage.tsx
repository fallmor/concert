import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import type { Concert } from '../types';
import { concertAPI } from '../services/api';

const ConcertDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [concert, setConcert] = useState<Concert | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchConcert = async () => {
      try {
        setLoading(true);
        const data = await concertAPI.getById(Number(id));
        setConcert(data);
      } catch (err) {
        setError('Concert not found');
        console.error('Error fetching concert:', err);
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchConcert();
    }
  }, [id]);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const formatTime = (dateString: string) => {
    return new Date(dateString).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: true
    });
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <h2>Loading concert details...</h2>
      </div>
    );
  }

  if (error || !concert) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <h2>❌ {error || 'Concert not found'}</h2>
        <button 
          onClick={() => navigate('/concerts')}
          style={{
            marginTop: '20px',
            padding: '10px 20px',
            backgroundColor: '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '5px',
            cursor: 'pointer',
            fontSize: '16px'
          }}
        >
          Back to Concerts
        </button>
      </div>
    );
  }

  const isAlmostSoldOut = concert.availableSeats < 20;
  const isSoldOut = concert.availableSeats === 0;

  return (
    <div style={{ maxWidth: '1000px', margin: '0 auto' }}>
      <button 
        onClick={() => navigate('/concerts')}
        style={{
          marginBottom: '20px',
          padding: '8px 16px',
          backgroundColor: '#6c757d',
          color: 'white',
          border: 'none',
          borderRadius: '5px',
          cursor: 'pointer',
          fontSize: '14px'
        }}
      >
        ← Back to Concerts
      </button>

      <div style={{
        backgroundColor: 'white',
        borderRadius: '15px',
        overflow: 'hidden',
        boxShadow: '0 4px 20px rgba(0,0,0,0.1)'
      }}>
        <div style={{
          height: '400px',
          background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: '120px'
        }}>
          🎵
        </div>

        <div style={{ padding: '40px' }}>
          <h1 style={{ 
            margin: '0 0 10px 0', 
            fontSize: '42px', 
            color: '#333' 
          }}>
            {concert.title}
          </h1>

          <div style={{ 
            display: 'flex', 
            alignItems: 'center', 
            gap: '10px',
            marginBottom: '30px'
          }}>
            <span style={{ fontSize: '24px' }}>🎤</span>
            <span style={{ fontSize: '24px', color: '#555', fontWeight: 'bold' }}>
              {concert.artist.name}
            </span>
            <span style={{
              marginLeft: '15px',
              padding: '4px 12px',
              backgroundColor: '#e7f3ff',
              color: '#007bff',
              borderRadius: '12px',
              fontSize: '14px',
              fontWeight: 'bold'
            }}>
              {concert.artist.genre}
            </span>
          </div>

          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
            gap: '20px',
            marginBottom: '30px',
            padding: '30px',
            backgroundColor: '#f8f9fa',
            borderRadius: '10px'
          }}>
            <div>
              <div style={{ fontSize: '14px', color: '#666', marginBottom: '5px' }}>
                📅 Date & Time
              </div>
              <div style={{ fontSize: '18px', fontWeight: 'bold', color: '#333' }}>
                {formatDate(concert.date)}
              </div>
              <div style={{ fontSize: '16px', color: '#555' }}>
                {formatTime(concert.date)}
              </div>
            </div>

            <div>
              <div style={{ fontSize: '14px', color: '#666', marginBottom: '5px' }}>
                📍 Venue
              </div>
              <div style={{ fontSize: '18px', fontWeight: 'bold', color: '#333' }}>
                {concert.venue}
              </div>
            </div>

            <div>
              <div style={{ fontSize: '14px', color: '#666', marginBottom: '5px' }}>
                💰 Ticket Price
              </div>
              <div style={{ fontSize: '28px', fontWeight: 'bold', color: '#007bff' }}>
                ${concert.price.toFixed(2)}
              </div>
            </div>

            <div>
              <div style={{ fontSize: '14px', color: '#666', marginBottom: '5px' }}>
                🎫 Availability
              </div>
              <div style={{ 
                fontSize: '18px', 
                fontWeight: 'bold', 
                color: isSoldOut ? '#dc3545' : isAlmostSoldOut ? '#ffc107' : '#28a745' 
              }}>
                {isSoldOut ? 'SOLD OUT' : `${concert.availableSeats} / ${concert.totalSeats} seats`}
              </div>
              {isAlmostSoldOut && !isSoldOut && (
                <div style={{ fontSize: '14px', color: '#4b4444ff', marginTop: '5px' }}>
                  ⚠️ Hurry! Almost sold out
                </div>
              )}
            </div>
          </div>

          {concert.description && (
            <div style={{ marginBottom: '30px' }}>
              <h2 style={{ fontSize: '24px', marginBottom: '15px' }}>About This Event</h2>
              <p style={{ fontSize: '16px', lineHeight: '1.6', color: '#555' }}>
                {concert.description}
              </p>
            </div>
          )}

          {concert.artist.bio && (
            <div style={{ marginBottom: '30px' }}>
              <h2 style={{ fontSize: '24px', marginBottom: '15px' }}>About {concert.artist.name}</h2>
              <p style={{ fontSize: '16px', lineHeight: '1.6', color: '#555' }}>
                {concert.artist.bio}
              </p>
            </div>
          )}

          <div style={{
            marginTop: '40px',
            padding: '30px',
            backgroundColor: '#f0f8ff',
            borderRadius: '10px',
            textAlign: 'center'
          }}>
            <h2 style={{ fontSize: '28px', marginBottom: '20px' }}>
              Ready to Book?
            </h2>
            <button
              disabled={isSoldOut}
              onClick={() => navigate(`/book/${concert.ID}`)} 
              style={{
                padding: '18px 50px',
                fontSize: '20px',
                fontWeight: 'bold',
                backgroundColor: isSoldOut ? '#ccc' : '#007bff',
                color: 'white',
                border: 'none',
                borderRadius: '8px',
                cursor: isSoldOut ? 'not-allowed' : 'pointer',
                transition: 'background-color 0.3s'
              }}
              onMouseEnter={(e) => {
                if (!isSoldOut) {
                  e.currentTarget.style.backgroundColor = '#0056b3';
                }
              }}
              onMouseLeave={(e) => {
                if (!isSoldOut) {
                  e.currentTarget.style.backgroundColor = '#007bff';
                }
              }}
            >
              {isSoldOut ? 'Sold Out' : 'Book Tickets Now'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ConcertDetailPage;