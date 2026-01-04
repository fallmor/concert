import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { adminAPI } from '../services/adminApi';
import type { Concert, Artist } from '../types';
import type {CreateShowInput } from '../services/adminApi';

const AdminShows: React.FC = () => {
  const navigate = useNavigate();
  const { user, isAuthenticated, isLoading } = useAuth();
  
  const [shows, setShows] = useState<Concert[]>([]);
  const [artists, setArtists] = useState<Artist[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [editingShow, setEditingShow] = useState<Concert | null>(null);

  useEffect(() => {
    if (isLoading) return;
    
    if (!isAuthenticated || (user?.role !== 'admin' && user?.role !== 'moderator')) {
      navigate('/');
      return;
    }

    fetchData();
  }, [isAuthenticated, isLoading, user, navigate]);

  const fetchData = async () => {
    try {
      const [showsData, artistsData] = await Promise.all([
        adminAPI.listShows(),
        adminAPI.listArtists(),
      ]);
       console.log('Artists loaded:', artistsData);
      setShows(showsData);
      setArtists(artistsData);
    } catch (err: any) {
      setError(err.message || 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingShow(null);
    setShowModal(true);
  };

  const handleEdit = (show: Concert) => {
    setEditingShow(show);
    setShowModal(true);
  };

  const handleDelete = async (id: number, title: string) => {
    if (!confirm(`Are you sure you want to delete "${title}"?`)) {
      return;
    }

    try {
      await adminAPI.deleteShow(id);
      setShows(prev => prev.filter(s => s.ID !== id));
      alert('Show deleted successfully!');
    } catch (err: any) {
      alert(err.message || 'Failed to delete show');
    }
  };

  const handleSave = async (show: CreateShowInput) => {
    try {
      if (editingShow) {
        const updated = await adminAPI.updateShow(editingShow.ID, show);
        setShows(prev => prev.map(s => s.ID === editingShow.ID ? updated : s));
      } else {
        const created = await adminAPI.createShow(show);
        setShows(prev => [created, ...prev]);
      }
      setShowModal(false);
      setEditingShow(null);
    } catch (err: any) {
      throw err;
    }
  };

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
        <h2> {error}</h2>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        marginBottom: '30px'
      }}>
        <div>
          <h1 style={{ fontSize: '36px', margin: '0 0 10px 0' }}>
            Manage Shows
          </h1>
          <button
            onClick={() => navigate('/admin')}
            style={{
              padding: '8px 16px',
              backgroundColor: '#6c757d',
              color: 'white',
              border: 'none',
              borderRadius: '5px',
              cursor: 'pointer',
              fontSize: '14px'
            }}
          >
            ← Back to Dashboard
          </button>
        </div>
        <button
          onClick={handleCreate}
          style={{
            padding: '12px 24px',
            backgroundColor: '#28a745',
            color: 'white',
            border: 'none',
            borderRadius: '8px',
            cursor: 'pointer',
            fontSize: '16px',
            fontWeight: 'bold'
          }}
        >
          + Create Show
        </button>
      </div>

      <div style={{ 
        backgroundColor: 'white', 
        borderRadius: '12px', 
        boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
        overflow: 'hidden'
      }}>
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr style={{ backgroundColor: '#f8f9fa', borderBottom: '2px solid #dee2e6' }}>
              <th style={tableHeaderStyle}>Title</th>
              <th style={tableHeaderStyle}>Artist</th>
              <th style={tableHeaderStyle}>Date</th>
              <th style={tableHeaderStyle}>Venue</th>
              <th style={tableHeaderStyle}>Price</th>
              <th style={tableHeaderStyle}>Seats</th>
              <th style={tableHeaderStyle}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {shows.length === 0 ? (
              <tr>
                <td colSpan={7} style={{ textAlign: 'center', padding: '40px', color: '#666' }}>
                  No shows found. Create your first show!
                </td>
              </tr>
            ) : (
              shows.map(show => (
                <tr key={show.ID} style={{ borderBottom: '1px solid #dee2e6' }}>
                  <td style={tableCellStyle}>
                    <strong>{show.title}</strong>
                  </td>
                  <td style={tableCellStyle}>{show.artist.name}</td>
                  <td style={tableCellStyle}>
                    {new Date(show.date).toLocaleDateString('en-US', {
                      year: 'numeric',
                      month: 'short',
                      day: 'numeric'
                    })}
                  </td>
                  <td style={tableCellStyle}>{show.venue}</td>
                  <td style={tableCellStyle}>${show.price.toFixed(2)}</td>
                  <td style={tableCellStyle}>
                    {show.availableSeats}/{show.totalSeats}
                  </td>
                  <td style={tableCellStyle}>
                    <button
                      onClick={() => handleEdit(show)}
                      style={{
                        padding: '6px 12px',
                        backgroundColor: '#007bff',
                        color: 'white',
                        border: 'none',
                        borderRadius: '4px',
                        cursor: 'pointer',
                        marginRight: '8px',
                        fontSize: '13px'
                      }}
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDelete(show.ID, show.title)}
                      style={{
                        padding: '6px 12px',
                        backgroundColor: '#dc3545',
                        color: 'white',
                        border: 'none',
                        borderRadius: '4px',
                        cursor: 'pointer',
                        fontSize: '13px'
                      }}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {showModal && (
        <ShowModal
          show={editingShow}
          artists={artists}
          onSave={handleSave}
          onClose={() => {
            setShowModal(false);
            setEditingShow(null);
          }}
        />
      )}
    </div>
  );
};

const tableHeaderStyle: React.CSSProperties = {
  padding: '16px',
  textAlign: 'left',
  fontWeight: 'bold',
  fontSize: '14px',
  color: '#495057'
};

const tableCellStyle: React.CSSProperties = {
  padding: '16px',
  fontSize: '14px',
  color: '#212529'
};


interface ShowModalProps {
  show: Concert | null;
  artists: Artist[];
  onSave: (show: CreateShowInput) => Promise<void>;
  onClose: () => void;
}

const ShowModal: React.FC<ShowModalProps> = ({ show, artists, onSave, onClose }) => {
  const [formData, setFormData] = useState<CreateShowInput>({
    title: show?.title || '',
    date: show?.date ? new Date(show.date).toISOString().split('T')[0] : '',
    time: show?.time || '20:00',
    artistId: show?.artist?.ID || 0,
    venue: show?.venue || '',
    price: show?.price || 50,
    totalSeats: show?.totalSeats || 100,
    description: show?.description || '',
    imageUrl: show?.imageUrl || '',
  });

  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;


  let finalValue: any = value;
    if (name === 'artistId' || name === 'totalSeats') {
    finalValue = parseInt(value, 10);
    console.log('  - converted to number:', finalValue);
  } else if (name === 'price') {
    finalValue = parseFloat(value);
  }
  setFormData(prev => {
    const newData = {
      ...prev,
      [name]: finalValue
    };
    console.log('  - New formData:', newData);
    return newData;
  });
};

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.title || !formData.date || !formData.venue || !formData.artistId || formData.artistId === 0) {
      setError('Please fill in all required fields');
      return;
    }

    try {
      setSaving(true);
      setError('');
      await onSave(formData);
    } catch (err: any) {
      setError(err.message || 'Failed to save show');
      setSaving(false);
    }
  };

  return (
    <div style={{
      position: 'fixed',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      backgroundColor: 'rgba(0,0,0,0.5)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      zIndex: 1000,
      padding: '20px'
    }}>
      <div style={{
        backgroundColor: 'white',
        borderRadius: '12px',
        maxWidth: '600px',
        width: '100%',
        maxHeight: '90vh',
        overflow: 'auto',
        padding: '30px'
      }}>
        <h2 style={{ margin: '0 0 20px 0', fontSize: '24px' }}>
          {show ? 'Edit Show' : 'Create Show'}
        </h2>

        {error && (
          <div style={{
            backgroundColor: '#f8d7da',
            color: '#721c24',
            padding: '12px',
            borderRadius: '5px',
            marginBottom: '20px'
          }}>
             {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '20px' }}>
            <label style={labelStyle}>Title *</label>
            <input
              type="text"
              name="title"
              value={formData.title}
              onChange={handleChange}
              required
              style={inputStyle}
            />
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '15px', marginBottom: '20px' }}>
            <div>
              <label style={labelStyle}>Date *</label>
              <input
                type="date"
                name="date"
                value={formData.date}
                onChange={handleChange}
                required
                style={inputStyle}
              />
            </div>
            <div>
              <label style={labelStyle}>Time</label>
              <input
                type="time"
                name="time"
                value={formData.time}
                onChange={handleChange}
                style={inputStyle}
              />
            </div>
          </div>
         <div style={{ marginBottom: '20px' }}>
  <label style={labelStyle}>
    Artist * 
    {artists.length === 0 && (
      <span style={{ color: '#dc3545', fontSize: '12px', marginLeft: '8px' }}>
        (No artists available - create artists first)
      </span>
    )}
  </label>
<select
  name="artistId"
  value={formData.artistId}
  onChange={handleChange}
  required
  style={{
    ...inputStyle,
    cursor: 'pointer',
    backgroundColor: 'white'
  }}
>
  <option value="">-- Select an artist --</option>
  {artists.map(artist => {
    console.log('Artist option:', artist.ID, artist.name); 
    return (
      <option key={artist.ID} value={artist.ID}>
        {artist.name} ({artist.genre})
      </option>
    );
  })}
</select>
  {formData.artistId > 0 && (
    <p style={{ margin: '5px 0 0 0', fontSize: '12px', color: '#28a745' }}>
      ✓ Artist selected: {artists.find(a => a.id === formData.artistId)?.name}
    </p>
  )}
</div>

          <div style={{ marginBottom: '20px' }}>
            <label style={labelStyle}>Venue *</label>
            <input
              type="text"
              name="venue"
              value={formData.venue}
              onChange={handleChange}
              required
              style={inputStyle}
            />
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '15px', marginBottom: '20px' }}>
            <div>
              <label style={labelStyle}>Price ($)</label>
              <input
                type="number"
                name="price"
                value={formData.price}
                onChange={handleChange}
                min="0"
                step="0.01"
                style={inputStyle}
              />
            </div>
            <div>
              <label style={labelStyle}>Total Seats</label>
              <input
                type="number"
                name="totalSeats"
                value={formData.totalSeats}
                onChange={handleChange}
                min="1"
                style={inputStyle}
              />
            </div>
          </div>

          <div style={{ marginBottom: '20px' }}>
            <label style={labelStyle}>Description</label>
            <textarea
              name="description"
              value={formData.description}
              onChange={handleChange}
              rows={3}
              style={{ ...inputStyle, resize: 'vertical' }}
            />
          </div>

          <div style={{ marginBottom: '25px' }}>
            <label style={labelStyle}>Image URL</label>
            <input
              type="url"
              name="imageUrl"
              value={formData.imageUrl}
              onChange={handleChange}
              placeholder="https://example.com/image.jpg"
              style={inputStyle}
            />
          </div>

          <div style={{ display: 'flex', gap: '10px', justifyContent: 'flex-end' }}>
            <button
              type="button"
              onClick={onClose}
              disabled={saving}
              style={{
                padding: '10px 20px',
                backgroundColor: '#6c757d',
                color: 'white',
                border: 'none',
                borderRadius: '6px',
                cursor: saving ? 'not-allowed' : 'pointer',
                fontSize: '14px'
              }}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={saving}
              style={{
                padding: '10px 20px',
                backgroundColor: saving ? '#ccc' : '#007bff',
                color: 'white',
                border: 'none',
                borderRadius: '6px',
                cursor: saving ? 'not-allowed' : 'pointer',
                fontSize: '14px',
                fontWeight: 'bold'
              }}
            >
              {saving ? 'Saving...' : 'Save Show'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const labelStyle: React.CSSProperties = {
  display: 'block',
  marginBottom: '6px',
  fontSize: '14px',
  fontWeight: 'bold',
  color: '#333'
};

const inputStyle: React.CSSProperties = {
  width: '100%',
  padding: '10px',
  fontSize: '14px',
  border: '1px solid #ddd',
  borderRadius: '6px',
  outline: 'none'
};

export default AdminShows;