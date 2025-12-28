import React from 'react';
import { Link } from 'react-router-dom';
import type { Concert } from '../types';


interface ConcertCardProps {
  concert: Concert;
}

const ConcertCard: React.FC<ConcertCardProps> = ({ concert }) => {
    const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', { 
      weekday: 'short',
      year: 'numeric', 
      month: 'short', 
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

  return (
    <Link 
      to={`/concerts/${concert.ID}`} 
      style={{ textDecoration: 'none', color: 'inherit' }}
    >
       <div style={{
        backgroundColor: 'white',
        borderRadius: '10px',
        overflow: 'hidden',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
        transition: 'transform 0.3s, box-shadow 0.3s',
        cursor: 'pointer',
        height: '100%'
      }}>
        {/* Image placeholder */}
        
        <div style={{
          height: '200px',
          background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: '64px'
        }}>
          {concert.artist?.album_url? (
              <img 
                src={concert.artist.album_url} 
                alt={concert.artist.name} 
                style={{
                  width: '100%',
                  height: '100%',
                  objectFit: 'cover' // Makes image fill the area
                }}
                loading="lazy"
              />
            ) : (
              <div className="image-placeholder">🎵</div>
            )}
        </div>

        <div style={{ padding: '20px' }}>
          {/* Title */}
          <h2 style={{ 
            margin: '0 0 15px 0', 
            fontSize: '22px', 
            color: '#333',
            fontWeight: 'bold'
          }}>
            {concert.title || 'Concert'}
          </h2>
          
          <div style={{ 
            display: 'flex', 
            alignItems: 'center', 
            gap: '8px',
            marginBottom: '10px'
          }}>
            <span style={{ fontSize: '18px' }}>🎤</span>
            <span style={{ fontSize: '16px', color: '#555', fontWeight: '500' }}>
              {concert.artist?.name || 'Artist TBA'}
            </span>
          </div>

          <div style={{ 
            display: 'flex', 
            alignItems: 'center', 
            gap: '8px',
            marginBottom: '10px'
          }}>
            <span style={{ fontSize: '16px' }}>📍</span>
            <span style={{ fontSize: '14px', color: '#666' }}>
              {concert.venue}
            </span>
          </div>
          
          <div style={{ 
            display: 'flex', 
            alignItems: 'center', 
            gap: '8px',
            marginBottom: '15px'
          }}>
            <span style={{ fontSize: '16px' }}>📅</span>
            <span style={{ fontSize: '14px', color: '#666' }}>
              {formatDate(concert.date)}
            </span>
            <span style={{ fontSize: '14px', color: '#999' }}>•</span>
            <span style={{ fontSize: '14px', color: '#666' }}>
              {formatTime(concert.date)}
            </span>
          </div>

          {concert.price && (
            <div style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              marginTop: '15px',
              paddingTop: '15px',
              borderTop: '1px solid #eee'
            }}>
              <span style={{ 
                fontSize: '24px', 
                fontWeight: 'bold', 
                color: '#007bff' 
              }}>
                ${concert.price.toFixed(2)}
              </span>
              
              {concert.available_seats !== undefined && (
                <span style={{ 
                  fontSize: '12px', 
                  color: concert.available_seats < 20 ? '#dc3545' : '#28a745',
                  fontWeight: 'bold'
                }}>
                  {concert.available_seats} seats left
                </span>
              )}
            </div>
          )}

          <button 
            onClick={(e) => {
              e.preventDefault();
              alert(`Booking for ${concert.title}`);
            }}
            style={{
              width: '100%',
              marginTop: '15px',
              backgroundColor: '#007bff',
              color: 'white',
              padding: '12px',
              border: 'none',
              borderRadius: '5px',
              fontSize: '16px',
              fontWeight: 'bold',
              cursor: 'pointer',
              transition: 'background-color 0.3s'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = '#0056b3';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = '#007bff';
            }}
          >
            Book Now
          </button>
        </div>
      </div>
    </Link>
  );
};

export default ConcertCard;