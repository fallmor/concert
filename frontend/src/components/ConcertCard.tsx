import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Card, Button, Tag, Typography } from 'antd';
import { CalendarOutlined, EnvironmentOutlined, UserOutlined, EuroOutlined } from '@ant-design/icons';
import type { Concert } from '../types';

const { Title, Text } = Typography;

interface ConcertCardProps {
  concert: Concert;
}

const ConcertCard: React.FC<ConcertCardProps> = ({ concert }) => {
  const navigate = useNavigate();

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

  const isSoldOut = concert.availableSeats === 0;
  const isLowStock = concert.availableSeats !== undefined && concert.availableSeats < 20 && concert.availableSeats > 0;

  const handleBookClick = (e: React.MouseEvent) => {
    e.preventDefault();
    navigate(`/book/${concert.ID}`);
  };

  return (
    <Link to={`/concerts/${concert.ID}`} className="no-underline">
      <Card
        hoverable
        className="h-full transition-all duration-300 hover:shadow-lg"
        cover={
          <div className="h-48 bg-gradient-to-r from-purple-500 to-pink-500 flex items-center justify-center overflow-hidden">
            {concert.imageUrl ? (
              <img 
                src={concert.imageUrl} 
                alt={concert.artist?.name || 'Concert'}
                className="w-full h-full object-cover"
                loading="lazy"
              />
            ) : (
                <div className="text-6xl">🎵</div>
            )}
          </div>
        }
        actions={[
          <Button
            type="primary"
            block
            size="large"
            onClick={handleBookClick}
            disabled={isSoldOut}
            className="m-2"
          >
            {isSoldOut ? 'Sold Out' : 'Book Now'}
          </Button>
        ]}
      >
        <Title level={4} className="mb-3 line-clamp-2 min-h-12">
          {concert.title || 'Concert'}
        </Title>

        <div className="space-y-2 mb-4">
          <div className="flex items-center gap-2">
            <UserOutlined className="text-blue-500" />
            <Text className="text-base font-medium">
              {concert.artist?.name || 'Artist TBA'}
            </Text>
          </div>

          <div className="flex items-center gap-2">
            <EnvironmentOutlined className="text-green-500" />
            <Text type="secondary">{concert.venue}</Text>
          </div>
          
          <div className="flex items-center gap-2">
            <CalendarOutlined className="text-orange-500" />
            <Text type="secondary">
              {formatDate(concert.date)} • {formatTime(concert.date)}
            </Text>
          </div>
        </div>

        <div className="flex items-center justify-between pt-3 border-t border-gray-200">
          {concert.price && (
            <div className="flex items-center gap-1">
              <EuroOutlined className="text-blue-500" />
              <Text strong className="text-xl text-blue-600">
                {concert.price.toFixed(2)}
              </Text>
            </div>
          )}

          {concert.availableSeats !== undefined && (
            <Tag color={isSoldOut ? 'red' : isLowStock ? 'orange' : 'green'}>
              {isSoldOut ? 'Sold Out' : `${concert.availableSeats} seats left`}
            </Tag>
          )}
        </div>
      </Card>
    </Link>
  );
};

export default ConcertCard;