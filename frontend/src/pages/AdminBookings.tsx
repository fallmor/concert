import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { adminAPI } from "../services/adminApi";
import type { AdminBooking } from "../services/adminApi";
import { Card, Row, Col, Statistic, Button, Typography, Spin, Input, Table, Tag, Empty } from 'antd';
import { TrophyOutlined, UserOutlined, CalendarOutlined, EuroOutlined, ArrowLeftOutlined, SearchOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

const AdminBookings: React.FC = () => {
  const navigate = useNavigate();
  const { user, isAuthenticated, isLoading } = useAuth();

  const [bookings, setBookings] = useState<AdminBooking[]>([]);
  const [filteredBookings, setFilteredBookings] = useState<AdminBooking[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [statusFilter, setStatusFilter] = useState<"all" | "confirmed" | "cancelled">("all");
  const [searchTerm, setSearchTerm] = useState("");

  useEffect(() => {
    if (isLoading) return;

    if (!isAuthenticated || (user?.role !== "admin" && user?.role !== "moderator")) {
      navigate("/");
      return;
    }

    fetchBookings();
  }, [isAuthenticated, isLoading, user, navigate]);

  useEffect(() => {
    let result = bookings;

    if (statusFilter !== "all") {
      result = result.filter((b) => b.status === statusFilter);
    }

    if (searchTerm) {
      const term = searchTerm.toLowerCase();
      result = result.filter(
        (b) =>
          b.username.toLowerCase().includes(term) ||
          b.userEmail.toLowerCase().includes(term) ||
          b.showTitle.toLowerCase().includes(term) ||
          b.artistName.toLowerCase().includes(term),
      );
    }

    setFilteredBookings(result);
  }, [bookings, statusFilter, searchTerm]);

  const fetchBookings = async () => {
    try {
      const data = await adminAPI.listBookings();
      setBookings(data);
      setFilteredBookings(data);
    } catch (err: any) {
      setError(err.message || "Failed to load bookings");
    } finally {
      setLoading(false);
    }
  };

  if (isLoading || loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <Spin size="large" tip="Loading bookings..." />
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-full flex justify-center px-4 py-12">
        <div className="max-w-7xl w-full">
          <Title level={2} className="text-red-500">{error}</Title>
        </div>
      </div>
    );
  }

  const stats = {
    total: bookings.length,
    confirmed: bookings.filter((b) => b.status === "confirmed").length,
    cancelled: bookings.filter((b) => b.status === "cancelled").length,
    totalRevenue: bookings
      .filter((b) => b.status === "confirmed")
      .reduce((sum, b) => sum + b.totalPrice, 0),
  };

  const columns = [
    {
      title: "User",
      key: "user",
      render: (_: any, booking: AdminBooking) => (
        <div>
          <div className="font-medium mb-1">{booking.username}</div>
          <Text type="secondary" className="text-xs">{booking.userEmail}</Text>
        </div>
      ),
    },
    {
      title: "Show",
      dataIndex: "showTitle",
      key: "show",
      render: (text: string) => <strong>{text}</strong>,
    },
    {
      title: "Artist",
      dataIndex: "artistName",
      key: "artist",
    },
    {
      title: "Tickets",
      dataIndex: "ticketCount",
      key: "tickets",
      render: (count: number) => (
        <Tag color="blue">{count} {count === 1 ? "ticket" : "tickets"}</Tag>
      ),
    },
    {
      title: "Total",
      dataIndex: "totalPrice",
      key: "total",
      render: (price: number, booking: AdminBooking) => (
        <Text strong type={booking.status === "cancelled" ? "secondary" : "success"}>
          {price.toFixed(2)}€
        </Text>
      ),
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (status: string) => (
        <Tag color={status === "confirmed" ? "green" : "red"}>{status.toUpperCase()}</Tag>
      ),
    },
    {
      title: "Date",
      dataIndex: "createdAt",
      key: "date",
      render: (date: string) => (
        <div>
          <div className="text-sm">
            {new Date(date).toLocaleDateString("en-US", {
              month: "short",
              day: "numeric",
              year: "numeric",
            })}
          </div>
          <Text type="secondary" className="text-xs">
            {new Date(date).toLocaleTimeString("en-US", {
              hour: "2-digit",
              minute: "2-digit",
            })}
          </Text>
        </div>
      ),
    },
  ];

  return (
    <div className="w-full flex justify-center px-4 py-12">
      <div className="max-w-7xl w-full">
        <div className="mb-6">
          <Title level={1} className="mb-2">All Bookings</Title>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate("/admin")} type="default">
            Back to Dashboard
          </Button>
        </div>

        <Row gutter={[16, 16]} className="mb-6">
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic title="Total Bookings" value={stats.total} prefix={<CalendarOutlined />} valueStyle={{ color: "#1890ff" }} />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic title="Confirmed" value={stats.confirmed} prefix={<TrophyOutlined />} valueStyle={{ color: "#52c41a" }} />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic title="Cancelled" value={stats.cancelled} prefix={<UserOutlined />} valueStyle={{ color: "#ff4d4f" }} />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic title="Total Revenue" value={stats.totalRevenue} prefix={<EuroOutlined />} precision={2} valueStyle={{ color: "#faad14" }} />
            </Card>
          </Col>
        </Row>

        <Card className="mb-6">
          <div className="space-y-4">
            <Input
              placeholder="Search by user, email, show, or artist..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              prefix={<SearchOutlined />}
              size="large"
              allowClear
            />
            <div className="flex flex-wrap gap-2">
              <Button type={statusFilter === "all" ? "primary" : "default"} onClick={() => setStatusFilter("all")}>
                All ({bookings.length})
              </Button>
              <Button
                type={statusFilter === "confirmed" ? "primary" : "default"}
                onClick={() => setStatusFilter("confirmed")}
                style={statusFilter === "confirmed" ? { backgroundColor: "#52c41a", borderColor: "#52c41a" } : {}}
              >
                Confirmed ({stats.confirmed})
              </Button>
              <Button
                type={statusFilter === "cancelled" ? "primary" : "default"}
                onClick={() => setStatusFilter("cancelled")}
                danger={statusFilter === "cancelled"}
              >
                Cancelled ({stats.cancelled})
              </Button>
            </div>
          </div>
        </Card>

        <Card>
          <Table
            dataSource={filteredBookings}
            columns={columns}
            rowKey="id"
            pagination={{
              pageSize: 10,
              showSizeChanger: true,
              showTotal: (total) => `Showing ${filteredBookings.length} of ${total} bookings`,
            }}
            locale={{
              emptyText: (
                <Empty
                  description={
                    searchTerm || statusFilter !== "all"
                      ? "No bookings match your filters."
                      : "No bookings yet."
                  }
                />
              ),
            }}
          />
        </Card>
      </div>
    </div>
  );
};

export default AdminBookings;
