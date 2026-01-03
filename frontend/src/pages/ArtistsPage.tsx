import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import type { Artist } from '../types';
import { artistAPI } from '../services/api';

const ArtistsPage: React.FC = () => {
  const [artists, setArtists] = useState<Artist[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedGenre, setSelectedGenre] = useState<string>('all');

  useEffect(() => {
    const fetchArtists = async () => {
      try {
        setLoading(true);
        const data = await artistAPI.getAll();
        setArtists(data);
      } catch (err) {
        setError('Failed to load artists');
        console.error('Error fetching artists:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchArtists();
  }, []);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <h2>Loading artists...</h2>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ textAlign: 'center', padding: '100px', color: '#dc3545' }}>
        <h2> {error}</h2>
      </div>
    );
  }

  // Get unique genres
  const genres = ['all', ...new Set(artists.map(a => a.genre).filter(Boolean))];

  // Filter artists by genre
  const filteredArtists = selectedGenre === 'all' 
    ? artists 
    : artists.filter(a => a.genre === selectedGenre);

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
      <h1 style={{ fontSize: '36px', marginBottom: '10px' }}>
        Featured Artists
      </h1>
      <p style={{ color: '#666', marginBottom: '30px', fontSize: '18px' }}>
        Discover the registered artists
      </p>

      {/* filter on genre */}
      <div style={{ marginBottom: '30px' }}>
        <label style={{ 
          display: 'block', 
          marginBottom: '10px', 
          fontWeight: 'bold',
          fontSize: '16px'
        }}>
          Filter by Genre:
        </label>
        <div style={{ display: 'flex', gap: '10px', flexWrap: 'wrap' }}>
          {genres.map(genre => (
            <button
              key={genre}
              onClick={() => setSelectedGenre(genre)}
              style={{
                padding: '8px 20px',
                backgroundColor: selectedGenre === genre ? '#007bff' : '#e9ecef',
                color: selectedGenre === genre ? 'white' : '#495057',
                border: 'none',
                borderRadius: '20px',
                cursor: 'pointer',
                fontSize: '14px',
                fontWeight: selectedGenre === genre ? 'bold' : 'normal',
                textTransform: 'capitalize',
                transition: 'all 0.3s'
              }}
            >
              {genre} ({genre === 'all' ? artists.length : artists.filter(a => a.genre === genre).length})
            </button>
          ))}
        </div>
      </div>

      {filteredArtists.length === 0 ? (
        <div style={{ textAlign: 'center', padding: '60px', color: '#999' }}>
          <p>No artists found in this genre</p>
        </div>
      ) : (
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
          gap: '30px'
        }}>
          {filteredArtists.map(artist => (
            <ArtistCard key={artist.id} artist={artist} />
          ))}
        </div>
      )}
    </div>
  );
};

interface ArtistCardProps {
  artist: Artist;
}

const ArtistCard: React.FC<ArtistCardProps> = ({ artist }) => {
  return (
    <Link 
      to={`/artists/${artist.id}`}
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
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = 'translateY(-5px)';
        e.currentTarget.style.boxShadow = '0 4px 16px rgba(0,0,0,0.2)';
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = 'translateY(0)';
        e.currentTarget.style.boxShadow = '0 2px 8px rgba(0,0,0,0.1)';
      }}
      >
        {/* artist image */}
        <div style={{
          height: '280px',
          background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          position: 'relative'
        }}>
          <div style={{
            width: '150px',
            height: '150px',
            borderRadius: '50%',
            backgroundColor: 'rgba(255,255,255,0.9)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: '64px',
            boxShadow: '0 4px 12px rgba(0,0,0,0.2)'
          }}>
            🎤
          </div>
        </div>

        <div style={{ padding: '20px' }}>
          <h3 style={{ 
            margin: '0 0 10px 0', 
            fontSize: '22px', 
            color: '#333',
            fontWeight: 'bold'
          }}>
            {artist.name}
          </h3>

          {artist.genre && (
            <span style={{
              display: 'inline-block',
              padding: '4px 12px',
              backgroundColor: '#e7f3ff',
              color: '#007bff',
              borderRadius: '12px',
              fontSize: '12px',
              fontWeight: 'bold',
              marginBottom: '15px'
            }}>
              {artist.genre}
            </span>
          )}

          {artist.bio && (
            <p style={{
              fontSize: '14px',
              color: '#666',
              lineHeight: '1.5',
              margin: '15px 0 0 0',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              display: '-webkit-box',
              WebkitLineClamp: 3,
              WebkitBoxOrient: 'vertical'
            }}>
              {artist.bio}
            </p>
          )}

          <div style={{
            marginTop: '15px',
            paddingTop: '15px',
            borderTop: '1px solid #eee'
          }}>
            <span style={{
              color: '#007bff',
              fontSize: '14px',
              fontWeight: 'bold'
            }}>
              View Profile →
            </span>
          </div>
        </div>
      </div>
    </Link>
  );
};

export default ArtistsPage;