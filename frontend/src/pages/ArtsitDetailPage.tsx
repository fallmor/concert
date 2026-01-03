import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import type { Artist, Concert } from '../types';
import { artistAPI } from '../services/api';

const ArtistDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [artist, setArtist] = useState<Artist | null>(null);
  const [shows, setShows] = useState<Concert[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchArtist = async () => {
      try {
        setLoading(true);
        const data = await artistAPI.getById(Number(id));
        setArtist(data);
        setShows(data.shows || []);
      } catch (err) {
        setError('Artist not found');
        console.error('Error fetching artist:', err);
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchArtist();
    }
  }, [id]);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <h2>Loading artist...</h2>
      </div>
    );
  }

  if (error || !artist) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <h2>❌ {error || 'Artist not found'}</h2>
        <button 
          onClick={() => navigate('/artists')}
          style={{
            marginTop: '20px',
            padding: '10px 20px',
            backgroundColor: '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '5px',
            cursor: 'pointer'
          }}
        >
          Back to Artists
        </button>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '1000px', margin: '0 auto' }}>
      <button 
        onClick={() => navigate('/artists')}
        style={{
          marginBottom: '20px',
          padding: '8px 16px',
          backgroundColor: '#6c757d',
          color: 'white',
          border: 'none',
          borderRadius: '5px',
          cursor: 'pointer'
        }}
      >
        ← Back to Artists
      </button>

      <div style={{
        backgroundColor: 'white',
        borderRadius: '15px',
        overflow: 'hidden',
        boxShadow: '0 4px 20px rgba(0,0,0,0.1)'
      }}>
        {/* Hero Section */}
        <div style={{
          height: '300px',
          background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          position: 'relative'
        }}>
          <div style={{
            width: '200px',
            height: '200px',
            borderRadius: '50%',
            backgroundColor: 'rgba(255,255,255,0.95)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: '100px',
            boxShadow: '0 8px 24px rgba(0,0,0,0.3)'
          }}>
            🎤
          </div>
        </div>

        <div style={{ padding: '40px' }}>
          <div style={{ textAlign: 'center', marginBottom: '40px' }}>
            <h1 style={{ 
              margin: '0 0 15px 0', 
              fontSize: '42px', 
              color: '#333' 
            }}>
              {artist.name}
            </h1>

            {artist.genre && (
              <span style={{
                display: 'inline-block',
                padding: '8px 20px',
                backgroundColor: '#e7f3ff',
                color: '#007bff',
                borderRadius: '20px',
                fontSize: '16px',
                fontWeight: 'bold'
              }}>
                {artist.genre}
              </span>
            )}
          </div>

          {artist.bio && (
            <div style={{ 
              marginBottom: '40px',
              padding: '30px',
              backgroundColor: '#f8f9fa',
              borderRadius: '10px'
            }}>
              <h2 style={{ fontSize: '24px', marginBottom: '15px' }}>About</h2>
              <p style={{ 
                fontSize: '16px', 
                lineHeight: '1.8', 
                color: '#555',
                margin: 0
              }}>
                {artist.bio}
              </p>
            </div>
          )}

          <div>
            <h2 style={{ fontSize: '28px', marginBottom: '20px' }}>
              Upcoming Performances ({shows.length})
            </h2>

            {shows.length === 0 ? (
              <div style={{
                textAlign: 'center',
                padding: '40px',
                backgroundColor: '#f8f9fa',
                borderRadius: '10px',
                color: '#666'
              }}>
                <p>No upcoming performances scheduled</p>
              </div>
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
                {shows.map(show => (
                  <Link
                    key={show.ID}
                    to={`/concerts/${show.ID}`}
                    style={{ textDecoration: 'none', color: 'inherit' }}
                  >
                    <div style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      padding: '20px',
                      backgroundColor: '#f8f9fa',
                      borderRadius: '8px',
                      border: '2px solid transparent',
                      transition: 'all 0.3s',
                      cursor: 'pointer'
                    }}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.borderColor = '#007bff';
                      e.currentTarget.style.backgroundColor = '#e7f3ff';
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.borderColor = 'transparent';
                      e.currentTarget.style.backgroundColor = '#f8f9fa';
                    }}
                    >
                      <div>
                        <h3 style={{ margin: '0 0 8px 0', fontSize: '20px' }}>
                          {show.title}
                        </h3>
                        <p style={{ margin: '5px 0', color: '#666', fontSize: '14px' }}>
                          📍 {show.venue}
                        </p>
                        <p style={{ margin: '5px 0', color: '#666', fontSize: '14px' }}>
                          📅 {new Date(show.date).toLocaleDateString('en-US', {
                            weekday: 'long',
                            year: 'numeric',
                            month: 'long',
                            day: 'numeric'
                          })}
                        </p>
                      </div>
                      <div style={{ textAlign: 'right' }}>
                        <div style={{ 
                          fontSize: '24px', 
                          fontWeight: 'bold', 
                          color: '#007bff' 
                        }}>
                          ${show.price.toFixed(2)}
                        </div>
                        <div style={{ 
                          fontSize: '12px', 
                          color: show.availableSeats < 20 ? '#dc3545' : '#28a745',
                          marginTop: '5px'
                        }}>
                          {show.availableSeats} seats left
                        </div>
                      </div>
                    </div>
                  </Link>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ArtistDetailPage;