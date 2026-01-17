import React, { useState } from "react";
import { Link } from "react-router-dom";
import { authAPI } from "../services/api";
import { Form, Button, Input, Alert } from 'antd';
import { CloseCircleOutlined } from '@ant-design/icons';

const ForgetPass: React.FC = () => {
  const [form] = Form.useForm();
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  const [alertVisible, setAlertVisible] = useState(true);

  const handleSubmit = async (values: { email: string }) => {
    try {
      setLoading(true);
      setError("");
      setSuccess(false);
      setAlertVisible(true);
      
      await authAPI.forget(values.email);
      setSuccess(true);
    } catch (err: any) {
      setError(err.message || 'Failed to send password reset email. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[calc(100vh-500px)] flex items-center justify-center px-4 py-8">
      <div className="w-full max-w-xl p-10 bg-white rounded-lg shadow-md">
        <h1 className="text-5xl font-bold text-center mb-6">Forget Password</h1>
        <p className="text-gray-600 text-center mb-8 text-lg">
          Request a new password
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

        {success && alertVisible && (
          <Alert
            type="success"
            showIcon
            className="mb-6"
            description={
              <div className="flex items-center justify-between">
                <span>Password reset email sent! Please check your email.</span>
                <CloseCircleOutlined 
                  className="cursor-pointer hover:text-green-600 ml-4"
                  onClick={() => {
                    setAlertVisible(false);
                    setSuccess(false);
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
          className="space-y-6"
          size="large"
        >
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

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              size="large"
              className="h-12 text-lg font-semibold"
            >
              {loading ? 'Sending...' : 'Request Password'}
            </Button>
          </Form.Item>
        </Form>

        <p className="text-center mt-8 text-gray-600 text-base">
          Remember your password?{' '}
          <Link to="/login" className="text-blue-600 hover:text-blue-800 font-semibold">
            Login here
          </Link>
        </p>
      </div>
    </div>
  );
};

export default ForgetPass;
