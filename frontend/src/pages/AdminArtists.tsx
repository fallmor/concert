import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { adminAPI } from "../services/adminApi";
import type { Artist } from "../types";
import type { CreateArtistInput } from "../services/adminApi";
import { Button, Card, Modal, Form, Input, Typography, Spin, Alert, Empty, Space, Tag } from "antd";
import { PlusOutlined, EditOutlined, DeleteOutlined, ArrowLeftOutlined, CloseCircleOutlined } from "@ant-design/icons";

const { Title, Text } = Typography;
const { TextArea } = Input;

const AdminArtists: React.FC = () => {
  const navigate = useNavigate();
  const { user, isAuthenticated, isLoading } = useAuth();
  const [form] = Form.useForm();

  const [artists, setArtists] = useState<Artist[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [editingArtist, setEditingArtist] = useState<Artist | null>(null);
  const [alertVisible, setAlertVisible] = useState(true);

  useEffect(() => {
    if (isLoading) return;

    if (
      !isAuthenticated ||
      (user?.role !== "admin" && user?.role !== "moderator")
    ) {
      navigate("/");
      return;
    }

    fetchArtists();
  }, [isAuthenticated, isLoading, user, navigate]);

  const fetchArtists = async () => {
    try {
      const data = await adminAPI.listArtists();
      setArtists(data);
    } catch (err: any) {
      setError(err.message || "Failed to load artists");
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingArtist(null);
    form.resetFields();
    setShowModal(true);
  };

  const handleEdit = (artist: Artist) => {
    setEditingArtist(artist);
    form.setFieldsValue({
      name: artist.name,
      genre: artist.genre,
      bio: artist.bio,
      imageUrl: artist.imageUrl,
    });
    setShowModal(true);
  };

  const handleDelete = async (id: number, name: string) => {
    Modal.confirm({
      title: "Delete Artist",
      content: `Are you sure you want to delete "${name}"?`,
      okText: "Yes, Delete",
      okType: "danger",
      cancelText: "Cancel",
      onOk: async () => {
        try {
          await adminAPI.deleteArtist(id);
          setArtists((prev) => prev.filter((a) => a.ID !== id));
        } catch (err: any) {
          Modal.error({ title: "Error", content: err.message || "Failed to delete artist" });
        }
      },
    });
  };

  const handleSave = async (values: CreateArtistInput) => {
    try {
      if (editingArtist) {
        const updated = await adminAPI.updateArtist(editingArtist.ID, values);
        setArtists((prev) => prev.map((a) => (a.ID === editingArtist.ID ? updated : a)));
      } else {
        const created = await adminAPI.createArtist(values);
        setArtists((prev) => [created, ...prev]);
      }
      setShowModal(false);
      setEditingArtist(null);
      form.resetFields();
    } catch (err: any) {
      throw err;
    }
  };

  if (isLoading || loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <Spin size="large" tip="Loading artists..." />
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

  return (
    <div className="w-full flex justify-center px-4 py-12">
      <div className="max-w-7xl w-full">
        <div className="flex justify-between items-center mb-6">
          <div>
            <Title level={1} className="mb-2">Manage Artists</Title>
            <Button icon={<ArrowLeftOutlined />} onClick={() => navigate("/admin")} className="mb-4">
              Back to Dashboard
            </Button>
          </div>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate} size="large">
            Create Artist
          </Button>
        </div>

        {artists.length === 0 ? (
          <Card>
            <Empty description="No artists found. Create your first artist!" />
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {artists.map((artist) => (
              <Card
                key={artist.ID}
                hoverable
                className="transition-all hover:shadow-lg hover:-translate-y-1"
                cover={
                  <div className="h-48 bg-gray-100 flex items-center justify-center overflow-hidden">
                    {artist.imageUrl ? (
                      <img src={artist.imageUrl} alt={artist.name} className="w-full h-full object-cover" />
                    ) : (
                      <span className="text-6xl">🎸</span>
                    )}
                  </div>
                }
                actions={[
                  <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(artist)}>
                    Edit
                  </Button>,
                  <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleDelete(artist.ID, artist.name)}>
                    Delete
                  </Button>,
                ]}
              >
                <Card.Meta
                  title={artist.name}
                  description={
                    <div>
                      <Tag color="blue" className="mb-2">{artist.genre}</Tag>
                      {artist.bio && (
                        <Text className="text-sm text-gray-600 line-clamp-3 block">{artist.bio}</Text>
                      )}
                    </div>
                  }
                />
              </Card>
            ))}
          </div>
        )}

        <Modal
          title={editingArtist ? "Edit Artist" : "Create Artist"}
          open={showModal}
          onCancel={() => {
            setShowModal(false);
            setEditingArtist(null);
            form.resetFields();
          }}
          footer={null}
          width={500}
        >
          <Form form={form} layout="vertical" onFinish={handleSave}>
            <Form.Item
              name="name"
              label="Name"
              rules={[{ required: true, message: "Artist name is required" }]}
            >
              <Input placeholder="e.g. The Beatles" />
            </Form.Item>

            <Form.Item name="genre" label="Genre">
              <Input placeholder="e.g. Rock, Jazz, Pop" />
            </Form.Item>

            <Form.Item name="bio" label="Bio">
              <TextArea rows={4} placeholder="Brief biography of the artist..." />
            </Form.Item>

            <Form.Item name="imageUrl" label="Image URL">
              <Input placeholder="https://example.com/artist.jpg" />
            </Form.Item>

            <Form.Item className="mb-0">
              <Space className="w-full justify-end">
                <Button
                  onClick={() => {
                    setShowModal(false);
                    setEditingArtist(null);
                    form.resetFields();
                  }}
                >
                  Cancel
                </Button>
                <Button type="primary" htmlType="submit">
                  {editingArtist ? "Update Artist" : "Create Artist"}
                </Button>
              </Space>
            </Form.Item>
          </Form>
        </Modal>
      </div>
    </div>
  );
};

export default AdminArtists;
