import React, { useState, useEffect } from 'react';
import type { Concert } from '../types';
import { concertAPI } from '../services/api';
import ConcertCard from '../components/ConcertCard';

const ConcertsPage: React.FC = () => {
  const [concerts, setConcerts] = useState<Concert[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchConcerts = async () => {
      try {
        setLoading(true);
        const data = await concertAPI.getAll();
        setConcerts(data);
        setError(null);
      } catch (err) {
        setError('Failed to load concerts. Is your Go backend running?');
        console.error('Error fetching concerts:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchConcerts();
  }, []); // Empty array = run once when component mounts


  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <h2>Loading concerts...</h2>
      </div>
    );
  }

  // Error handling
  if (error) {
    return (
      <div style={{ 
        textAlign: 'center', 
        padding: '100px',
        color: '#dc3545'
      }}>
        <h2> {error}</h2>
        <p>Make sure your Go backend is running on http://localhost:8080</p>
      </div>
    );
  }

  if (concerts.length === 0) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <h2>No concerts available yet</h2>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
      <h1 style={{ fontSize: '36px', marginBottom: '30px' }}>
        Upcoming Concerts ({concerts.length})
      </h1>

      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fill, minmax(350px, 1fr))',
        gap: '30px'
      }}>
        {concerts.map(concert => (
          <ConcertCard key={concert.ID} concert={concert} />
        ))}
        
      </div>
    </div>
  );
};


export default ConcertsPage;