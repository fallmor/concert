import React, { useState, useEffect } from 'react';
import type { Concert } from '../types';
import { concertAPI } from '../services/api';
import ConcertCard from '../components/ConcertCard';
import { Typography, Alert, Spin } from 'antd';
import { CloseCircleOutlined } from '@ant-design/icons';

const { Title } = Typography;

const ConcertsPage: React.FC = () => {
  const [concerts, setConcerts] = useState<Concert[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  const [alertVisible, setAlertVisible] = useState<boolean>(true);

  useEffect(() => {
    const fetchConcerts = async () => {
      try {
        setLoading(true);
        const data = await concertAPI.getAll();
        setConcerts(data);
        setError('');
        setAlertVisible(false);
      } catch (err: unknown) {
        const errorMessage = err instanceof Error ? err.message : 'Failed to load concerts. Is your Go backend running?';
        setError(errorMessage);
        setAlertVisible(true);
        console.error('Error fetching concerts:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchConcerts();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Spin size="large" tip="Loading concerts..." />
      </div>
    );
  }

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

        <Title level={2} className="text-3xl font-bold mb-6">
          Upcoming Concerts ({concerts.length})
        </Title>

        {concerts.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">No concerts available at the moment.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 justify-items-center">
            {concerts.map((concert: Concert) => (
              <div key={concert.ID} className="w-full max-w-sm">
                <ConcertCard concert={concert} />
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};


export default ConcertsPage;