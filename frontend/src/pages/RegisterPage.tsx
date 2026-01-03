import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const RegisterPage: React.FC = () => {
  const navigate = useNavigate();
  const { register } = useAuth();
  
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    firstName: '',
    lastName: '',
    password: '',
    confirmPassword: ''
  });
  
  const [errors, setErrors] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    general: ''
  });
  
  const [loading, setLoading] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    
    // Clear error when user types
    setErrors(prev => ({ ...prev, [name]: '', general: '' }));
  };

  const validate = (): boolean => {
    const newErrors = {
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      general: ''
    };
    
    let isValid = true;

    if (formData.username.length < 3) {
      newErrors.username = 'Username must be at least 3 characters';
      isValid = false;
    }

    if (!formData.email.includes('@')) {
      newErrors.email = 'Please enter a valid email';
      isValid = false;
    }

    if (formData.password.length < 6) {
      newErrors.password = 'Password must be at least 6 characters';
      isValid = false;
    }

    if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match';
      isValid = false;
    }

    setErrors(newErrors);
    return isValid;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validate()) return;

    try {
      setLoading(true);
      await register(formData.username, formData.email, formData.firstName, formData.lastName, formData.password);
      
      // Redirect to home page after successful registration
      navigate('/');
    } catch (error: any) {
      setErrors(prev => ({
        ...prev,
        general: error.message || 'Registration failed. Email might already be in use.'
      }));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{
      maxWidth: '450px',
      margin: '60px auto',
      backgroundColor: 'white',
      padding: '40px',
      borderRadius: '15px',
      boxShadow: '0 4px 20px rgba(0,0,0,0.1)'
    }}>
      <h1 style={{ 
        textAlign: 'center', 
        marginBottom: '10px',
        fontSize: '32px',
        color: '#333'
      }}>
        Create Account
      </h1>
      <p style={{ 
        textAlign: 'center', 
        color: '#666', 
        marginBottom: '30px' 
      }}>
        Join us to book amazing concerts
      </p>

      {errors.general && (
        <div style={{
          backgroundColor: '#f8d7da',
          color: '#721c24',
          padding: '12px',
          borderRadius: '5px',
          marginBottom: '20px',
          fontSize: '14px'
        }}>
           {errors.general}
        </div>
      )}

      <form onSubmit={handleSubmit}>

        <div style={{ marginBottom: '20px' }}>
          <label style={{ 
            display: 'block', 
            marginBottom: '8px', 
            fontWeight: 'bold',
            fontSize: '14px',
            color: '#333'
          }}>
            Username
          </label>
          <input
            type="text"
            name="username"
            value={formData.username}
            onChange={handleChange}
            required
            style={{
              width: '100%',
              padding: '12px',
              fontSize: '16px',
              border: errors.username ? '2px solid #dc3545' : '1px solid #ddd',
              borderRadius: '8px',
              outline: 'none',
              transition: 'border-color 0.3s'
            }}
            onFocus={(e) => {
              if (!errors.username) {
                e.target.style.borderColor = '#007bff';
              }
            }}
            onBlur={(e) => {
              if (!errors.username) {
                e.target.style.borderColor = '#ddd';
              }
            }}
          />
          {errors.username && (
            <p style={{ color: '#dc3545', fontSize: '13px', margin: '5px 0 0 0' }}>
              {errors.username}
            </p>
          )}
        </div>

        <div style={{ marginBottom: '20px' }}>
          <label style={{ 
            display: 'block', 
            marginBottom: '8px', 
            fontWeight: 'bold',
            fontSize: '14px',
            color: '#333'
          }}>
            Email
          </label>
          <input
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            required
            style={{
              width: '100%',
              padding: '12px',
              fontSize: '16px',
              border: errors.email ? '2px solid #dc3545' : '1px solid #ddd',
              borderRadius: '8px',
              outline: 'none'
            }}
            onFocus={(e) => {
              if (!errors.email) {
                e.target.style.borderColor = '#007bff';
              }
            }}
            onBlur={(e) => {
              if (!errors.email) {
                e.target.style.borderColor = '#ddd';
              }
            }}
          />
          {errors.email && (
            <p style={{ color: '#dc3545', fontSize: '13px', margin: '5px 0 0 0' }}>
              {errors.email}
            </p>
          )}
        </div>

        {/* FirstNamr */}
        <div style={{ marginBottom: '20px' }}>
          <label style={{ 
            display: 'block', 
            marginBottom: '8px', 
            fontWeight: 'bold',
            fontSize: '14px',
            color: '#333'
          }}>
            First Name
          </label>
          <input
            type="text"
            name="firstName"
            value={formData.firstName}
            onChange={handleChange}
            required
            style={{
              width: '100%',
              padding: '12px',
              fontSize: '16px',
              border: '1px solid #ddd',
              borderRadius: '8px',
              outline: 'none'
            }}
          />
        </div>

        <div style={{ marginBottom: '20px' }}>
          <label style={{ 
            display: 'block', 
            marginBottom: '8px', 
            fontWeight: 'bold',
            fontSize: '14px',
            color: '#333'
          }}>
            Last Name
          </label>
          <input
            type="text"
            name="lastName"
            value={formData.lastName}
            onChange={handleChange}
            required
            style={{
              width: '100%',
              padding: '12px',
              fontSize: '16px',
              border:'1px solid #ddd',
              borderRadius: '8px',
              outline: 'none'
            }}
          />
        </div>

        <div style={{ marginBottom: '20px' }}>
          <label style={{ 
            display: 'block', 
            marginBottom: '8px', 
            fontWeight: 'bold',
            fontSize: '14px',
            color: '#333'
          }}>
            Password
          </label>
          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            required
            style={{
              width: '100%',
              padding: '12px',
              fontSize: '16px',
              border: errors.password ? '2px solid #dc3545' : '1px solid #ddd',
              borderRadius: '8px',
              outline: 'none'
            }}
            onFocus={(e) => {
              if (!errors.password) {
                e.target.style.borderColor = '#007bff';
              }
            }}
            onBlur={(e) => {
              if (!errors.password) {
                e.target.style.borderColor = '#ddd';
              }
            }}
          />
          {errors.password && (
            <p style={{ color: '#dc3545', fontSize: '13px', margin: '5px 0 0 0' }}>
              {errors.password}
            </p>
          )}
        </div>

        <div style={{ marginBottom: '25px' }}>
          <label style={{ 
            display: 'block', 
            marginBottom: '8px', 
            fontWeight: 'bold',
            fontSize: '14px',
            color: '#333'
          }}>
            Confirm Password
          </label>
          <input
            type="password"
            name="confirmPassword"
            value={formData.confirmPassword}
            onChange={handleChange}
            required
            style={{
              width: '100%',
              padding: '12px',
              fontSize: '16px',
              border: errors.confirmPassword ? '2px solid #dc3545' : '1px solid #ddd',
              borderRadius: '8px',
              outline: 'none'
            }}
            onFocus={(e) => {
              if (!errors.confirmPassword) {
                e.target.style.borderColor = '#007bff';
              }
            }}
            onBlur={(e) => {
              if (!errors.confirmPassword) {
                e.target.style.borderColor = '#ddd';
              }
            }}
          />
          {errors.confirmPassword && (
            <p style={{ color: '#dc3545', fontSize: '13px', margin: '5px 0 0 0' }}>
              {errors.confirmPassword}
            </p>
          )}
        </div>
        <button
          type="submit"
          disabled={loading}
          style={{
            width: '100%',
            padding: '14px',
            backgroundColor: loading ? '#ccc' : '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '8px',
            fontSize: '18px',
            fontWeight: 'bold',
            cursor: loading ? 'not-allowed' : 'pointer',
            transition: 'background-color 0.3s'
          }}
          onMouseEnter={(e) => {
            if (!loading) {
              e.currentTarget.style.backgroundColor = '#0056b3';
            }
          }}
          onMouseLeave={(e) => {
            if (!loading) {
              e.currentTarget.style.backgroundColor = '#007bff';
            }
          }}
        >
          {loading ? 'Creating Account...' : 'Register'}
        </button>
      </form>

      <p style={{ 
        textAlign: 'center', 
        marginTop: '25px',
        color: '#666',
        fontSize: '14px'
      }}>
        Already have an account?{' '}
        <Link 
          to="/login" 
          style={{ 
            color: '#007bff', 
            textDecoration: 'none',
            fontWeight: 'bold'
          }}
        >
          Login here
        </Link>
      </p>
    </div>
  );
};

export default RegisterPage;