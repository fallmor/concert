import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import  { adminAPI } from '../services/adminApi';
import type { AdminStats } from '../services/adminApi';

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
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <h2>Loading...</h2>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <h2>❌ {error}</h2>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
      <h1 style={{ fontSize: '36px', marginBottom: '10px' }}>
        Admin Dashboard
      </h1>
      <p style={{ color: '#666', marginBottom: '40px', fontSize: '16px' }}>
        Welcome, {user?.username}! Manage your concert booking system.
      </p>

      {stats && (
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))',
          gap: '20px',
          marginBottom: '40px'
        }}>
          <StatCard
            title="Total Shows"
            value={stats.totalShows}
            icon="🎤"
            color="#007bff"
          />
          <StatCard
            title="Total Artists"
            value={stats.totalArtists}
            icon="🎸"
            color="#28a745"
          />
          <StatCard
            title="Total Bookings"
            value={stats.totalBookings}
            icon="🎫"
            color="#ffc107"
          />
          <StatCard
            title="Total Revenue"
            value={`$${stats.totalRevenue.toFixed(2)}`}
            icon="💰"
            color="#dc3545"
          />
          <StatCard
            title="Total Users"
            value={stats.totalUsers}
            icon="👥"
            color="#6f42c1"
          />
        </div>
      )}

      {/* Management*/}
      <h2 style={{ fontSize: '24px', marginBottom: '20px' }}>
        Management
      </h2>

      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
        gap: '20px'
      }}>
        <ManagementCard
          title="Manage Shows"
          description="Create, edit, and delete concerts"
          icon="🎤"
          onClick={() => navigate('/admin/shows')}
        />
        <ManagementCard
          title="Manage Artists"
          description="Create, edit, and delete artists"
          icon="🎸"
          onClick={() => navigate('/admin/artists')}
        />
        <ManagementCard
          title="View All Bookings"
          description="See all user bookings and statistics"
          icon="🎫"
          onClick={() => navigate('/admin/bookings')}
        />
      </div>
    </div>
  );
};

interface StatCardProps {
  title: string;
  value: number | string;
  icon: string;
  color: string;
}

const StatCard: React.FC<StatCardProps> = ({ title, value, icon, color }) => {
  return (
    <div style={{
      backgroundColor: 'white',
      borderRadius: '12px',
      padding: '25px',
      boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
      borderLeft: `4px solid ${color}`
    }}>
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <div>
          <p style={{ 
            margin: '0 0 8px 0', 
            fontSize: '14px', 
            color: '#666',
            fontWeight: '500'
          }}>
            {title}
          </p>
          <p style={{ 
            margin: 0, 
            fontSize: '32px', 
            fontWeight: 'bold',
            color: '#333'
          }}>
            {value}
          </p>
        </div>
        <div style={{ fontSize: '40px' }}>
          {icon}
        </div>
      </div>
    </div>
  );
};

interface ManagementCardProps {
  title: string;
  description: string;
  icon: string;
  onClick: () => void;
}

const ManagementCard: React.FC<ManagementCardProps> = ({ title, description, icon, onClick }) => {
  return (
    <div
      onClick={onClick}
      style={{
        backgroundColor: 'white',
        borderRadius: '12px',
        padding: '30px',
        boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
        cursor: 'pointer',
        transition: 'transform 0.2s, box-shadow 0.2s',
        border: '2px solid transparent'
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = 'translateY(-5px)';
        e.currentTarget.style.boxShadow = '0 4px 20px rgba(0,0,0,0.15)';
        e.currentTarget.style.borderColor = '#007bff';
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = 'translateY(0)';
        e.currentTarget.style.boxShadow = '0 2px 10px rgba(0,0,0,0.1)';
        e.currentTarget.style.borderColor = 'transparent';
      }}
    >
      <div style={{ fontSize: '48px', marginBottom: '15px' }}>
        {icon}
      </div>
      <h3 style={{ 
        margin: '0 0 10px 0', 
        fontSize: '20px',
        color: '#333'
      }}>
        {title}
      </h3>
      <p style={{ 
        margin: 0, 
        fontSize: '14px', 
        color: '#666',
        lineHeight: '1.5'
      }}>
        {description}
      </p>
    </div>
  );
};

export default AdminDashboard;