import React from "react";
import { Link } from "react-router-dom";
import { Button, Card } from "antd";

const HomePage: React.FC = () => {
  return (
    <div className="w-full">
      <div className="w-full flex justify-center bg-gradient-to-r from-blue-500 to-blue-600 text-white py-16 mb-8 shadow-lg">
        <div className="max-w-7xl mx-auto px-4">
          <div className="flex flex-col items-center justify-center text-center">
            <h1 className="text-4xl font-bold mb-4">
              🎸 Welcome to Concert Booking
            </h1>
            <p className="text-lg text-blue-100 mb-8 max-w-2xl">
              Discover amazing live music events and book your tickets now
            </p>
            <Link to="/concerts">
              <Button type="primary" size="large" className="h-12 px-8 text-lg font-semibold">
                Browse Concerts
              </Button>
            </Link>
          </div>
        </div>
      </div>

      <div className="w-full flex justify-center px-4 py-12">
        <div className="max-w-7xl w-full">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 justify-items-center">
            <div className="w-full max-w-sm">
              <FeatureCard
                icon="🎤"
                title="Top Artists"
                description="See performances from the biggest names in music"
              />
            </div>
            <div className="w-full max-w-sm">
              <FeatureCard
                icon="🎫"
                title="Easy Booking"
                description="Book tickets quickly and securely online"
              />
            </div>
            <div className="w-full max-w-sm">
              <FeatureCard
                icon="📍"
                title="Favorite Locations"
                description="Concerts at the best locations around the world"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

interface FeatureCardProps {
  icon: string;
  title: string;
  description: string;
}

const FeatureCard: React.FC<FeatureCardProps> = ({
  icon,
  title,
  description,
}) => {
  return (
    <Card
      className="text-center hover:shadow-xl transition-all duration-300 h-full"
      hoverable
      style={{ borderRadius: '12px' }}
    >
      <div className="text-6xl mb-6">{icon}</div>
      <h3 className="mb-4 text-2xl font-bold text-gray-800">{title}</h3>
      <p className="text-gray-600 leading-relaxed text-base">
        {description}
      </p>
    </Card>
  );
};

export default HomePage;
