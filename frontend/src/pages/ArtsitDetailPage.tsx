import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import type { Artist, Concert } from '../types';
import { artistAPI } from '../services/api';
import { Button, Card, Spin, Tag, Alert, Typography, Empty } from 'antd';
import { CloseCircleOutlined, ArrowLeftOutlined, CalendarOutlined, EnvironmentOutlined } from '@ant-design/icons';

const { Title } = Typography;

const ArtistDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [artist, setArtist] = useState<Artist | null>(null);
  const [shows, setShows] = useState<Concert[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  const [alertVisible, setAlertVisible] = useState<boolean>(true);

  useEffect(() => {
    const fetchArtist = async () => {
      try {
        setLoading(true);
        const data = await artistAPI.getById(Number(id));
        setArtist(data);
        setShows(data.shows || []);
      } catch (err: unknown) {
        const errorMessage = err instanceof Error ? err.message : 'Artist not found';
        setError(errorMessage);
        setAlertVisible(true);
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
      <div className="flex items-center justify-center min-h-[400px]">
        <Spin size="large" tip="Loading artist..." />
      </div>
    );
  }

  if (error || !artist) {
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
                  <span>{error || 'Artist not found'}</span>
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
          <div className="text-center py-12">
            <Title level={2} className="mb-6">❌ {error || 'Artist not found'}</Title>
            <Button
              type="primary"
              size="large"
              onClick={() => navigate('/artists')}
              icon={<ArrowLeftOutlined />}
            >
              Back to Artists
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full flex justify-center px-4 py-12">
      <div className="max-w-7xl w-full">
        <Button
          type="default"
        onClick={() => navigate('/artists')}
          icon={<ArrowLeftOutlined />}
          className="mb-6"
      >
          Back to Artists
        </Button>

        <Card className="shadow-lg">
          <div className="h-80 bg-gradient-to-r from-pink-400 to-red-500 flex items-center justify-center relative mb-8 rounded-t-lg">
            {artist.imageUrl ? (
              <img
                src={artist.imageUrl}
                alt={artist.name}
                className="w-full h-full object-cover"
              />
            ) : (
              <div className="w-48 h-48 rounded-full bg-white/95 flex items-center justify-center text-8xl shadow-2xl">
                🎤
              </div>
            )}
        </div>

          <div className="p-8">
            <div className="text-center mb-8">
              <Title level={1} className="mb-4">
              {artist.name}
              </Title>
            {artist.genre && (
                <Tag color="blue" className="text-base px-4 py-1">
                {artist.genre}
                </Tag>
            )}
          </div>

          {artist.bio && (
              <Card className="mb-8 bg-gray-50">
                <Title level={3} className="mb-4">About</Title>
                <p className="text-base text-gray-600 leading-relaxed mb-0">
                {artist.bio}
              </p>
              </Card>
          )}

          <div>
              <Title level={2} className="mb-6">
              Upcoming Performances ({shows.length})
              </Title>

            {shows.length === 0 ? (
                <Empty
                  description="No upcoming performances scheduled"
                  className="py-12"
                />
            ) : (
                  <div className="flex flex-col gap-4">
                    {shows.map((show: Concert) => (
                  <Link
                    key={show.ID}
                    to={`/concerts/${show.ID}`}
                      className="no-underline"
                  >
                      <Card
                        hoverable
                        className="border-2 border-transparent hover:border-blue-500 hover:bg-blue-50 transition-all"
                    >
                        <div className="flex justify-between items-center">
                          <div>
                            <h3 className="text-xl font-semibold mb-2">{show.title}</h3>
                            <p className="text-gray-600 text-sm mb-1">
                              <EnvironmentOutlined className="mr-2" />
                              {show.venue}
                            </p>
                            <p className="text-gray-600 text-sm">
                              <CalendarOutlined className="mr-2" />
                              {new Date(show.date).toLocaleDateString('en-US', {
                                weekday: 'long',
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric'
                              })}
                            </p>
                          </div>
                          <div className="text-right">
                            <div className="text-2xl font-bold text-blue-600 mb-1">
                              ${show.price.toFixed(2)}
                            </div>
                            <Tag color={show.availableSeats < 20 ? 'red' : 'green'}>
                              {show.availableSeats} seats left
                            </Tag>
                        </div>
                      </div>
                      </Card>
                  </Link>
                ))}
              </div>
            )}
          </div>
        </div>
        </Card>
      </div>
    </div>
  );
};

export default ArtistDetailPage;