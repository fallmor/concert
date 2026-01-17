import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { adminAPI } from '../services/adminApi';
import type { Concert, Artist } from '../types';
import type { CreateShowInput } from '../services/adminApi';
import { Button, Card, Table, Modal, Form, Input, InputNumber, Select, Typography, Spin, Alert, Empty, Space } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, ArrowLeftOutlined, CloseCircleOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;
const { TextArea } = Input;

const AdminShows: React.FC = () => {
  const navigate = useNavigate();
  const { user, isAuthenticated, isLoading } = useAuth();
  const [form] = Form.useForm();
  
  const [shows, setShows] = useState<Concert[]>([]);
  const [artists, setArtists] = useState<Artist[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [editingShow, setEditingShow] = useState<Concert | null>(null);
  const [alertVisible, setAlertVisible] = useState(true);

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
    form.resetFields();
    form.setFieldsValue({ time: '20:00', price: 50, totalSeats: 100 });
    setShowModal(true);
  };

  const handleEdit = (show: Concert) => {
    setEditingShow(show);
    form.setFieldsValue({
      title: show.title,
      date: show.date ? new Date(show.date).toISOString().split('T')[0] : '',
      time: show.time || '20:00',
      artistId: show.artist?.ID || 0,
      venue: show.venue,
      price: show.price,
      totalSeats: show.totalSeats,
      description: show.description,
      imageUrl: show.imageUrl,
    });
    setShowModal(true);
  };

  const handleDelete = async (id: number, title: string) => {
    Modal.confirm({
      title: 'Delete Show',
      content: `Are you sure you want to delete "${title}"?`,
      okText: 'Yes, Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          await adminAPI.deleteShow(id);
          setShows(prev => prev.filter(s => s.ID !== id));
        } catch (err: any) {
          Modal.error({ title: 'Error', content: err.message || 'Failed to delete show' });
        }
      },
    });
  };

  const handleSave = async (values: CreateShowInput) => {
    try {
      if (editingShow) {
        const updated = await adminAPI.updateShow(editingShow.ID, values);
        setShows(prev => prev.map(s => s.ID === editingShow.ID ? updated : s));
      } else {
        const created = await adminAPI.createShow(values);
        setShows(prev => [created, ...prev]);
      }
      setShowModal(false);
      setEditingShow(null);
      form.resetFields();
    } catch (err: any) {
      throw err;
    }
  };

  if (isLoading || loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <Spin size="large" tip="Loading shows..." />
      </div>
    );
  }

  if (error) {
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
                      setError(null);
                    }}
                  />
                </div>
              }
            />
          )}
        </div>
      </div>
    );
  }

  const columns = [
    {
      title: 'Title',
      dataIndex: 'title',
      key: 'title',
      render: (text: string) => <strong>{text}</strong>,
    },
    {
      title: 'Artist',
      key: 'artist',
      render: (_: any, record: Concert) => record.artist?.name || 'N/A',
    },
    {
      title: 'Date',
      dataIndex: 'date',
      key: 'date',
      render: (date: string) => new Date(date).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      }),
    },
    {
      title: 'Venue',
      dataIndex: 'venue',
      key: 'venue',
    },
    {
      title: 'Price',
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => `€${price.toFixed(2)}`,
    },
    {
      title: 'Seats',
      key: 'seats',
      render: (_: any, record: Concert) => `${record.availableSeats}/${record.totalSeats}`,
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: Concert) => (
        <Space>
          <Button type="primary" icon={<EditOutlined />} onClick={() => handleEdit(record)} size="small">
            Edit
          </Button>
          <Button danger icon={<DeleteOutlined />} onClick={() => handleDelete(record.ID, record.title)} size="small">
            Delete
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="w-full flex justify-center px-4 py-12">
      <div className="max-w-7xl w-full">
        <div className="flex justify-between items-center mb-6">
          <div>
            <Title level={1} className="mb-2">Manage Shows</Title>
            <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/admin')} className="mb-4">
              Back to Dashboard
            </Button>
          </div>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} size="large">
            Create Show
          </Button>
        </div>

        <Card>
          <Table
            columns={columns}
            dataSource={shows}
            rowKey="ID"
            pagination={{ pageSize: 10 }}
            locale={{
              emptyText: <Empty description="No shows found. Create your first show!" />
            }}
          />
        </Card>

        <Modal
          title={editingShow ? 'Edit Show' : 'Create Show'}
          open={showModal}
          onCancel={() => {
            setShowModal(false);
            setEditingShow(null);
            form.resetFields();
          }}
          footer={null}
          width={600}
        >
          <Form
            form={form}
            layout="vertical"
            onFinish={handleSave}
            initialValues={{ time: '20:00', price: 50, totalSeats: 100 }}
          >
            <Form.Item
              name="title"
              label="Title"
              rules={[{ required: true, message: 'Please enter show title' }]}
            >
              <Input placeholder="Enter show title" />
            </Form.Item>

            <div className="grid grid-cols-2 gap-4">
              <Form.Item
                name="date"
                label="Date"
                rules={[{ required: true, message: 'Please select date' }]}
              >
                <Input type="date" />
              </Form.Item>
              <Form.Item name="time" label="Time">
                <Input type="time" />
              </Form.Item>
            </div>

            <Form.Item
              name="artistId"
              label="Artist"
              rules={[{ required: true, message: 'Please select an artist' }]}
            >
              <Select placeholder="Select an artist" disabled={artists.length === 0}>
                {artists.map(artist => (
                  <Select.Option key={artist.ID} value={artist.ID}>
                    {artist.name} ({artist.genre})
                  </Select.Option>
                ))}
              </Select>
              {artists.length === 0 && (
                <Text type="danger" className="text-xs">
                  No artists available - create artists first
                </Text>
              )}
            </Form.Item>

            <Form.Item
              name="venue"
              label="Venue"
              rules={[{ required: true, message: 'Please enter venue' }]}
            >
              <Input placeholder="Enter venue" />
            </Form.Item>

            <div className="grid grid-cols-2 gap-4">
              <Form.Item name="price" label="Price (€)">
                <InputNumber min={0} step={0.01} style={{ width: '100%' }} />
              </Form.Item>
              <Form.Item name="totalSeats" label="Total Seats">
                <InputNumber min={1} style={{ width: '100%' }} />
              </Form.Item>
            </div>

            <Form.Item name="description" label="Description">
              <TextArea rows={3} placeholder="Enter show description" />
            </Form.Item>

            <Form.Item name="imageUrl" label="Image URL">
              <Input placeholder="https://example.com/image.jpg" />
            </Form.Item>

            <Form.Item className="mb-0">
              <Space className="w-full justify-end">
                <Button onClick={() => {
                  setShowModal(false);
                  setEditingShow(null);
                  form.resetFields();
                }}>
                  Cancel
                </Button>
                <Button type="primary" htmlType="submit">
                  {editingShow ? 'Update Show' : 'Create Show'}
                </Button>
              </Space>
            </Form.Item>
          </Form>
        </Modal>
      </div>
    </div>
  );
};

export default AdminShows;
