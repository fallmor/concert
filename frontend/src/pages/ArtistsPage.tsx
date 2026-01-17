import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import type { Artist } from '../types';
import { artistAPI } from '../services/api';
import { Typography, Spin, Alert, Card, Tag, Empty, Button } from 'antd';
import { CloseCircleOutlined } from '@ant-design/icons';

const { Title } = Typography;

const ArtistsPage: React.FC = () => {
  const [artists, setArtists] = useState<Artist[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  const [selectedGenre, setSelectedGenre] = useState<string>('all');
  const [alertVisible, setAlertVisible] = useState<boolean>(true);

  useEffect(() => {
    const fetchArtists = async () => {
      try {
        setLoading(true);
        const data = await artistAPI.getAll();
        setArtists(data);
      } catch (err: unknown) {
        const errorMessage = err instanceof Error ? err.message : 'Failed to load artists';
        setError(errorMessage);
        setAlertVisible(true);
        console.error('Error fetching artists:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchArtists();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Spin size="large" tip="Loading artists..." />
      </div>
    );
  }

  // get unique genres
  const genres = ['all', ...new Set(artists.map(a => a.genre).filter(Boolean))];

  // filter artists by genre
  const filteredArtists = selectedGenre === 'all' 
    ? artists 
    : artists.filter(a => a.genre === selectedGenre);

  return (
    <div className="w-full flex justify-center px-4 py-12">
      <div className="max-w-7xl w-full">
        {error && alertVisible && (
          <Alert
            type="error"
            showIcon
            className="mb-6"
            description={
              <div className="flex items-center justify-between">
                <span>{error}</span>
                <CloseCircleOutlined
                  className="cursor-pointer hover:text-red-600 ml-4"
                  onClick={() => {
                    setAlertVisible(false);
                    setError('');
                  }}
                />
              </div>
            }
          />
        )}

        <Title level={1} className="text-4xl font-bold mb-2">
          Featured Artists
        </Title>
        <p className="text-gray-600 mb-8 text-lg">
          Discover the registered artists
        </p>

        {/* Filter by genre */}
        <div className="mb-8">
          <label className="block mb-3 text-base font-medium text-gray-700">
            Filter by Genre:
          </label>
          <div className="flex gap-3 flex-wrap">
            {genres.map(genre => (
              <Button
                key={genre}
                type={selectedGenre === genre ? 'primary' : 'default'}
                onClick={() => setSelectedGenre(genre)}
                className="capitalize"
              >
                {genre} ({genre === 'all' ? artists.length : artists.filter(a => a.genre === genre).length})
              </Button>
            ))}
          </div>
        </div>

        {filteredArtists.length === 0 ? (
          <Empty
            description="No artists found in this genre"
            className="py-12"
          />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 justify-items-center">
            {filteredArtists.map((artist: Artist) => (
              <div key={artist.ID} className="w-full max-w-sm">
                <ArtistCard artist={artist} />
              </div>
          ))}
        </div>
      )}
      </div>
    </div>
  );
};

interface ArtistCardProps {
  artist: Artist;
}

const ArtistCard: React.FC<ArtistCardProps> = ({ artist }) => {
  return (
    <Link to={`/artists/${artist.ID}`} className="no-underline">
      <Card
        hoverable
        className="h-full transition-all duration-300 hover:shadow-xl"
        cover={
          <div className="h-64 bg-gradient-to-r from-pink-400 to-red-500 flex items-center justify-center relative">
            {artist.imageUrl ? (
              <img
                src={artist.imageUrl}
                alt={artist.name}
                className="w-full h-full object-cover"
              />
            ) : (
              <div className="w-36 h-36 rounded-full bg-white/90 flex items-center justify-center text-6xl shadow-lg">
                🎤
              </div>
            )}
          </div>
        }
      >
        <Card.Meta
          title={<h3 className="text-xl font-bold mb-2">{artist.name}</h3>}
          description={
            <div>
              {artist.genre && (
                <Tag color="blue" className="mb-3">
                  {artist.genre}
                </Tag>
              )}
              {artist.bio && (
                <p className="text-gray-600 text-sm line-clamp-3 leading-relaxed">
                  {artist.bio}
                </p>
              )}
              <div className="mt-4 pt-4 border-t border-gray-200">
                <span className="text-blue-600 text-sm font-semibold">
                  View Profile →
                </span>
              </div>
            </div>
          }
        />
      </Card>
    </Link>
  );
};

export default ArtistsPage;