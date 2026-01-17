import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Form, Button, Input, Alert } from 'antd';
import { CloseCircleOutlined } from '@ant-design/icons';

const RegisterPage: React.FC = () => {
  const navigate = useNavigate();
  const { register } = useAuth();
  const [form] = Form.useForm();
  
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [alertVisible, setAlertVisible] = useState(true);

  const handleSubmit = async (values: {
    username: string;
    email: string;
    firstName: string;
    lastName: string;
    password: string;
    confirmPassword: string;
  }) => {
    try {
      setLoading(true);
      setError('');
      setAlertVisible(true);

      await register(
        values.username,
        values.email,
        values.firstName,
        values.lastName,
        values.password
      );

      navigate('/');
    } catch (err: any) {
      setError(err.message || 'Registration failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[calc(100vh-500px)] flex items-center justify-center px-4 py-8">
      <div className="w-full max-w-xl p-10 bg-white rounded-lg shadow-md">
        <h1 className="text-5xl font-bold text-center mb-6">Create Account</h1>
        <p className="text-gray-600 text-center mb-8 text-lg">
          Register to book your favorite concerts
        </p>

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

        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          className="space-y-4"
          size="large"
        >
          <Form.Item
            name="username"
            label={<span className="text-base font-medium">Username</span>}
            rules={[
              { required: true, message: 'Please enter your username' },
              { min: 3, message: 'Username must be at least 3 characters' }
            ]}
          >
            <Input
              placeholder="Enter your username"
              className="py-3"
            />
          </Form.Item>

          <Form.Item
            name="email"
            label={<span className="text-base font-medium">Email</span>}
            rules={[
              { required: true, message: 'Please enter your email' },
              { type: 'email', message: 'Please enter a valid email' }
            ]}
          >
            <Input
              type="email"
              placeholder="Enter your email"
              className="py-3"
            />
          </Form.Item>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Form.Item
              name="firstName"
              label={<span className="text-base font-medium">First Name</span>}
              rules={[{ required: true, message: 'Please enter your first name' }]}
            >
              <Input
                placeholder="First name"
                className="py-3"
              />
            </Form.Item>

            <Form.Item
              name="lastName"
              label={<span className="text-base font-medium">Last Name</span>}
              rules={[{ required: true, message: 'Please enter your last name' }]}
            >
              <Input
                placeholder="Last name"
                className="py-3"
              />
            </Form.Item>
          </div>

          <Form.Item
            name="password"
            label={<span className="text-base font-medium">Password</span>}
            rules={[
              { required: true, message: 'Please enter your password' },
              { min: 6, message: 'Password must be at least 6 characters' }
            ]}
          >
            <Input.Password
              placeholder="Enter your password"
              className="py-3"
            />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            label={<span className="text-base font-medium">Confirm Password</span>}
            dependencies={['password']}
            rules={[
              { required: true, message: 'Please confirm your password' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('Password does not match!'));
                },
              }),
            ]}
          >
            <Input.Password
              placeholder="Confirm your password"
              className="py-3"
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              size="large"
              className="h-12 text-lg font-semibold"
            >
              {loading ? 'Creating Account...' : 'Register'}
            </Button>
          </Form.Item>
        </Form>

        <p className="text-center mt-8 text-gray-600 text-base">
          Already have an account?{' '}
          <Link to="/login" className="text-blue-600 hover:text-blue-800 font-semibold">
            Login here
          </Link>
        </p>
      </div>
    </div>
  );
};

export default RegisterPage;