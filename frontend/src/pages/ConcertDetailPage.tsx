import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import type { Concert } from '../types';
import { concertAPI } from '../services/api';
import { Button, Card, Spin, Tag, Alert } from 'antd';
import { CloseCircleOutlined, ArrowLeftOutlined, CalendarOutlined, EnvironmentOutlined, DollarOutlined, UserOutlined } from '@ant-design/icons';


const ConcertDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [concert, setConcert] = useState<Concert | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  const [alertVisible, setAlertVisible] = useState<boolean>(true);

  useEffect(() => {
    const fetchConcert = async () => {
      try {
        setLoading(true);
        const data = await concertAPI.getById(Number(id));
        setConcert(data);
      } catch (err: unknown) {
        const errorMessage = err instanceof Error ? err.message : 'Concert not found';
        setError(errorMessage);
        setAlertVisible(true);
        console.error('Error fetching concert:', err);
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchConcert();
    }
  }, [id]);

  const formatDateTimestamp = (dateString: string): string => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: true
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Spin size="large" tip="Loading concert details..." />
      </div>
    );
  }

  if (error || !concert) {
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
                  <span>{error || 'Concert not found'}</span>
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
            <h2 className="text-2xl font-bold mb-6">❌ {error || 'Concert not found'}</h2>
            <Button
              type="primary"
              size="large"
              onClick={() => navigate('/concerts')}
              icon={<ArrowLeftOutlined />}
            >
              Back to Concerts
            </Button>
          </div>
        </div>
      </div>
    );
  }

  const isAlmostSoldOut = concert.availableSeats < 20;
  const isSoldOut = concert.availableSeats === 0;

  return (
    <div className="w-full flex justify-center px-4 py-12">
      <div className="max-w-7xl w-full">
        <Button
          type="default"
        onClick={() => navigate('/concerts')}
          icon={<ArrowLeftOutlined />}
          className="mb-6"
      >
          Back to Concerts
        </Button>

        <Card className="shadow-lg">
          <div className="h-96 bg-gradient-to-r from-blue-500 to-purple-600 text-white flex items-center justify-center text-8xl mb-8 rounded-t-lg">
            {concert.imageUrl ? (
              <img
                src={concert.imageUrl}
                alt={concert.title}
                className="w-full h-full object-cover"
              />
            ) : (
              <span>🎵</span>
            )}
        </div>

          <div className="p-8">
            <h1 className="text-4xl font-bold mb-4 text-gray-800">
            {concert.title}
          </h1>

            <div className="flex items-center gap-4 mb-8 flex-wrap">
              <UserOutlined className="text-2xl text-gray-600" />
              <span className="text-2xl font-bold text-gray-600">
              {concert.artist.name}
            </span>
              <Tag color="blue" className="text-base px-3 py-1">
              {concert.artist.genre}
              </Tag>
          </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8 p-6 bg-gray-50 rounded-lg">
              <Card.Meta
                avatar={<CalendarOutlined className="text-2xl text-blue-600" />}
                title={<span className="text-sm text-gray-600">Date & Time</span>}
                description={<span className="text-lg font-bold text-gray-800">{formatDateTimestamp(concert.date)}</span>}
              />

              <Card.Meta
                avatar={<EnvironmentOutlined className="text-2xl text-green-600" />}
                title={<span className="text-sm text-gray-600">Venue</span>}
                description={<span className="text-lg font-bold text-gray-800">{concert.venue}</span>}
              />

              <Card.Meta
                avatar={<DollarOutlined className="text-2xl text-purple-600" />}
                title={<span className="text-sm text-gray-600">Ticket Price</span>}
                description={<span className="text-2xl font-bold text-blue-600">${concert.price.toFixed(2)}</span>}
              />

              <Card.Meta
                avatar={<UserOutlined className="text-2xl text-orange-600" />}
                title={<span className="text-sm text-gray-600">Availability</span>}
                description={
                  <div>
                    <span className={`text-lg font-bold ${isSoldOut ? 'text-red-500' : isAlmostSoldOut ? 'text-yellow-500' : 'text-green-500'}`}>
                      {isSoldOut ? 'SOLD OUT' : `${concert.availableSeats} / ${concert.totalSeats} seats`}
                    </span>
                    {isAlmostSoldOut && !isSoldOut && (
                      <Tag color="warning" className="mt-2">⚠️ Hurry! Almost sold out</Tag>
                    )}
                  </div>
              }
              />
          </div>

          {concert.description && (
              <Card className="mb-6">
                <h2 className="text-2xl font-bold mb-4">About This Event</h2>
                <p className="text-base text-gray-600 leading-relaxed">
                {concert.description}
              </p>
              </Card>
          )}

          {concert.artist.bio && (
              <Card className="mb-6">
                <h2 className="text-2xl font-bold mb-4">About {concert.artist.name}</h2>
                <p className="text-base text-gray-600 leading-relaxed">
                {concert.artist.bio}
              </p>
              </Card>
          )}

            <Card className="mt-8 bg-blue-50 border-blue-200">
              <div className="text-center">
                <h2 className="text-2xl font-bold mb-6">Ready to Book?</h2>
                <Button
                  type="primary"
                  size="large"
                  disabled={isSoldOut}
                  onClick={() => navigate(`/book/${concert.ID}`)}
                  className="px-10 py-4 text-lg font-semibold h-12"
                >
                  {isSoldOut ? 'Sold Out' : 'Book Tickets Now'}
                </Button>
              </div>
            </Card>
        </div>
        </Card>
      </div>
    </div>
  );
};

export default ConcertDetailPage;