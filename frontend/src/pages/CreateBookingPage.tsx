import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import type { Concert } from '../types';
import { concertAPI, bookingAPI } from '../services/api';
import { Button, Card, Spin, Alert, Form, InputNumber, Typography, Statistic, Row, Col } from 'antd';
import { CloseCircleOutlined, ArrowLeftOutlined, CalendarOutlined, EnvironmentOutlined, EuroOutlined, UserOutlined } from '@ant-design/icons';

const { Title } = Typography;

const CreateBookingPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated, isLoading } = useAuth();
  const [form] = Form.useForm();
  
  const [concert, setConcert] = useState<Concert | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [alertVisible, setAlertVisible] = useState<boolean>(true);
  const [ticketCount, setTicketCount] = useState<number>(1);

  useEffect(() => {
    if (isLoading) return;
    
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    const fetchConcert = async () => {
      try {
        setLoading(true);
        const data = await concertAPI.getById(Number(id));
        setConcert(data);
        
        form.setFieldsValue({ ticketCount: 1 });
        setTicketCount(1);
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
  }, [id, isAuthenticated, isLoading, navigate, form]);

  const handleSubmit = async (values: { ticketCount: number }) => {
    if (!concert) return;

    try {
      setSubmitting(true);
      setError('');
      setAlertVisible(false);

      await bookingAPI.create(concert.ID, values.ticketCount);

      navigate('/my-bookings', {
        state: { message: 'Booking created successfully!' }
      });
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create booking';
      setError(errorMessage);
      setAlertVisible(true);
    } finally {
      setSubmitting(false);
    }
  };

  if (isLoading || loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Spin size="large" tip="Loading concert details..." />
      </div>
    );
  }

  if (error && !concert) {
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
          <div className="text-center py-12">
            <Title level={2} className="mb-6">❌ {error || 'Concert not found'}</Title>
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

  if (!concert) return null;

  const isSoldOut = concert.availableSeats === 0;
  const maxTickets = concert.availableSeats;
  const totalPrice = ticketCount * concert.price;

  return (
    <div className="w-full flex justify-center px-4 py-12">
      <div className="max-w-4xl w-full">
        <Button
          type="default"
          onClick={() => navigate(`/concerts/${concert.ID}`)}
          icon={<ArrowLeftOutlined />}
          className="mb-6"
        >
          Back to Concert Details
        </Button>

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

        <Card className="shadow-lg">
          <Title level={2} className="mb-6">Book Tickets</Title>
          
          <Card className="mb-6 bg-gray-50">
            <Title level={4} className="mb-4">{concert.title}</Title>
            <Row gutter={[16, 16]}>
              <Col xs={24} sm={12} md={6}>
                <Statistic
                  title="Artist"
                  value={concert.artist.name}
                  prefix={<UserOutlined />}
                />
              </Col>
              <Col xs={24} sm={12} md={6}>
                <Statistic
                  title="Date & Time"
                  value={new Date(concert.date).toLocaleString('en-US', {
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                  })}
                  prefix={<CalendarOutlined />}
                />
              </Col>
              <Col xs={24} sm={12} md={6}>
                <Statistic
                  title="Venue"
                  value={concert.venue}
                  prefix={<EnvironmentOutlined />}
                />
              </Col>
              <Col xs={24} sm={12} md={6}>
                <Statistic
                  title="Price per Ticket"
                  value={concert.price}
                  prefix={<EuroOutlined />}
                  precision={2}
                />
              </Col>
            </Row>
          </Card>

          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
            size="large"
          >
            <Form.Item
              name="ticketCount"
              label={<span className="text-base font-medium">Number of Tickets</span>}
              rules={[
                { required: true, message: 'Please select number of tickets' },
                { type: 'number', min: 1, message: 'At least 1 ticket required' },
                { type: 'number', max: maxTickets, message: `Maximum ${maxTickets} tickets available` }
              ]}
            >
              <InputNumber
                min={1}
                max={maxTickets}
                disabled={isSoldOut}
                className="w-full"
                style={{ width: '100%' }}
                placeholder="Select number of tickets"
                value={ticketCount}
                onChange={(value) => setTicketCount(value || 1)}
              />
            </Form.Item>

            <Card className="mb-6 bg-blue-50 border-blue-200">
              <Row gutter={16}>
                <Col span={12}>
                  <Statistic
                    title="Tickets"
                    value={ticketCount}
                  />
                </Col>
                <Col span={12}>
                  <Statistic
                    title="Total Price"
                    value={totalPrice}
                    prefix="€"
                    precision={2}
                    valueStyle={{ color: '#1890ff' }}
                  />
                </Col>
              </Row>
            </Card>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={submitting}
                disabled={isSoldOut}
                block
                size="large"
                className="h-12 text-lg font-semibold"
              >
                {isSoldOut ? 'Sold Out' : submitting ? 'Processing...' : 'Confirm Booking'}
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </div>
    </div>
  );
};

export default CreateBookingPage;
