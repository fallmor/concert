import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { adminAPI } from '../services/adminApi';
import type { AdminStats } from '../services/adminApi';
import { Card, Statistic, Button, Typography, Spin, Alert } from 'antd';
import { EuroOutlined, CustomerServiceOutlined, TeamOutlined, ShoppingCartOutlined, UserOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

const AdminDashboard: React.FC = () => {
  const navigate = useNavigate();
  const { user, isAuthenticated, isLoading } = useAuth();
  
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isLoading) return;
    
    if (!isAuthenticated || (user?.role !== 'admin' && user?.role !== 'moderator')) {
      navigate('/');
      return;
    }

    const fetchStats = async () => {
      try {
        const data = await adminAPI.getStats();
        setStats(data);
      } catch (err: any) {
        setError(err.message || 'Failed to load stats');
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, [isAuthenticated, isLoading, user, navigate]);

  if (isLoading || loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <Spin size="large" tip="Loading dashboard..." />
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-full flex justify-center px-4 py-12">
        <div className="max-w-7xl w-full">
          <Alert
            type="error"
            message={error}
            action={
              <Button type="primary" onClick={() => navigate('/admin')}>
                Back to Dashboard
              </Button>
            }
          />
        </div>
      </div>
    );
  }

  const managementCards = [
    { title: 'Manage Shows', description: 'Create, edit, and delete concerts', path: '/admin/shows', icon: <CustomerServiceOutlined className="text-4xl" /> },
    { title: 'Manage Artists', description: 'Create, edit, and delete artists', path: '/admin/artists', icon: <TeamOutlined className="text-4xl" /> },
    { title: 'View All Bookings', description: 'See all user bookings and statistics', path: '/admin/bookings', icon: <ShoppingCartOutlined className="text-4xl" /> },
  ];

  return (
    <div className="w-full flex justify-center px-4 py-12">
      <div className="max-w-7xl w-full">
        <Title level={1} className="mb-2">Admin Dashboard</Title>
        <Text type="secondary" className="mb-8 block">
          Welcome, {user?.username}! Manage your concert booking system.
        </Text>

        <Title level={2} className="mb-6">Statistics</Title>
        {stats && (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
            <Card>
              <Statistic title="Total Shows" value={stats.totalShows} prefix={<CustomerServiceOutlined />} />
            </Card>
            <Card>
              <Statistic title="Total Artists" value={stats.totalArtists} prefix={<TeamOutlined />} />
            </Card>
            <Card>
              <Statistic title="Total Bookings" value={stats.totalBookings} prefix={<ShoppingCartOutlined />} />
            </Card>
            <Card>
              <Statistic title="Total Revenue" value={stats.totalRevenue} prefix={<EuroOutlined />} precision={2} valueStyle={{ color: '#dc3545' }} />
            </Card>
            <Card>
              <Statistic title="Total Users" value={stats.totalUsers} prefix={<UserOutlined />} />
            </Card>
          </div>
        )}

        <Title level={2} className="mb-6">Management</Title>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {managementCards.map((card) => (
            <Card
              key={card.path}
              hoverable
              className="cursor-pointer transition-all hover:shadow-lg hover:-translate-y-1 border-2 border-transparent hover:border-blue-500"
              onClick={() => navigate(card.path)}
            >
              <div className="flex flex-col items-center text-center">
                <div className="text-blue-600 mb-4">{card.icon}</div>
                <Title level={4} className="mb-2">{card.title}</Title>
                <Text type="secondary">{card.description}</Text>
              </div>
            </Card>
          ))}
        </div>
      </div>
    </div>
  );
};

export default AdminDashboard;