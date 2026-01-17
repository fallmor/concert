import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import type { Booking } from '../types';
import { bookingAPI } from '../services/api';
import { Typography, Spin, Alert, Card, Tag, Empty, Button, Modal } from 'antd';
import { CloseCircleOutlined, CheckCircleOutlined, CalendarOutlined, ShoppingOutlined } from '@ant-design/icons';

const { Title } = Typography;

const MyBookingsPage: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated, isLoading } = useAuth();
  
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  const [successMessage, setSuccessMessage] = useState<string>('');
  const [alertVisible, setAlertVisible] = useState<boolean>(true);

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
      } catch (err: unknown) {
        const errorMessage = err instanceof Error ? err.message : 'Failed to load bookings';
        setError(errorMessage);
        setAlertVisible(true);
      } finally {
        setLoading(false);
      }
    };

    fetchBookings();
  }, [isAuthenticated, isLoading, navigate, location.state]);

  const handleCancel = async (bookingId: number) => {
    Modal.confirm({
      title: 'Cancel Booking',
      content: 'Are you sure you want to cancel this booking?',
      okText: 'Yes, Cancel',
      okType: 'danger',
      cancelText: 'No',
      onOk: async () => {
        try {
          await bookingAPI.cancel(bookingId);

          setBookings(prev =>
            prev.map(b =>
              b.ID === bookingId
                ? { ...b, status: 'cancelled' as const }
                : b
            )
          );

          setSuccessMessage('Booking cancelled successfully');
          setAlertVisible(true);
          setTimeout(() => {
            setSuccessMessage('');
            setAlertVisible(false);
          }, 5000);
        } catch (err: unknown) {
          const errorMessage = err instanceof Error ? err.message : 'Failed to cancel booking';
          setError(errorMessage);
          setAlertVisible(true);
        }
      }
    });
  };

  if (isLoading || loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Spin size="large" tip={isLoading ? 'Loading...' : 'Loading your bookings...'} />
      </div>
    );
  }

  return (
    <div className="w-full flex justify-center px-4 py-12">
      <div className="max-w-7xl w-full">
        <Title level={1} className="mb-8">
          My Bookings
        </Title>

        {successMessage && alertVisible && (
          <Alert
            type="success"
            showIcon
            icon={<CheckCircleOutlined />}
            className="mb-6"
            description={
              <div className="flex items-center justify-between">
                <span>{successMessage}</span>
                <CloseCircleOutlined
                  className="cursor-pointer hover:text-green-600 ml-4"
                  onClick={() => {
                    setAlertVisible(false);
                    setSuccessMessage('');
                  }}
                />
              </div>
          }
          />
      )}

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

      {bookings.length === 0 ? (
          <Empty
            description={
              <div>
                <p className="text-lg text-gray-600 mb-4">
                  You haven't booked any concerts yet.
                </p>
                <Link to="/concerts">
                  <Button type="primary" size="large">
                    Browse Concerts
                  </Button>
                </Link>
              </div>
          }
            className="py-12"
          />
      ) : (
            <div className="flex flex-col gap-6">
              {bookings.map((booking: Booking) => (
            <BookingCard
                key={booking.ID}
              booking={booking}
              onCancel={handleCancel}
            />
          ))}
        </div>
      )}
      </div>
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
    <Card
      className={`transition-all ${isCancelled ? 'opacity-60 border-2 border-red-500' : ''}`}
    >
      <div className="flex justify-between items-center flex-wrap gap-4">
        <div className="flex-1">
          <div className="flex items-center gap-3 mb-4 flex-wrap">
            <Title level={3} className="mb-0">
              {booking.show?.title || 'Concert'}
            </Title>
            {isCancelled && (
              <Tag color="red" className="text-xs font-bold">
                CANCELLED
              </Tag>
            )}
          </div>

          <div className="text-gray-600 space-y-2">
            <p className="mb-0">
              <ShoppingOutlined className="mr-2" />
              {booking.ticketCount} ticket{booking.ticketCount !== 1 ? 's' : ''}
            </p>
            <p className="mb-0">
              <CalendarOutlined className="mr-2" />
              Booked on: {new Date(booking.CreatedAt).toLocaleDateString('en-US', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
              })}
            </p>
          </div>
        </div>

        <div className="flex flex-col items-end gap-4">
          <div className={`text-3xl font-bold ${isCancelled ? 'text-gray-400' : 'text-blue-600'}`}>
            €{booking.totalPrice.toFixed(2)}
          </div>

          {!isCancelled && (
            <Button
              danger
              onClick={() => onCancel(booking.ID)}
            >
              Cancel Booking
            </Button>
          )}
        </div>
      </div>
    </Card>
  );
};

export default MyBookingsPage;