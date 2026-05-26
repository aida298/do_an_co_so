import React, { useState, useEffect } from 'react';
import api from '../services/api';
import { useNavigate } from 'react-router-dom';

function MyOrders() {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchMyOrders = async () => {
      try {
        const response = await api.get('/orders/my-orders');
        setOrders(response.data || []);
      } catch (error) {
        console.error("Lỗi lấy danh sách đơn hàng:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchMyOrders();
  }, []);

  if (loading) return <div className="main-content"><p>Đang tải lịch sử đơn hàng...</p></div>;

  return (
    <div className="main-content">
      <h2 style={{ marginBottom: '25px' }}>Lịch sử Đơn hàng của bạn</h2>
      
      {orders.length === 0 ? (
        <div style={{ textAlign: 'center', marginTop: '50px', background: 'white', padding: '40px', borderRadius: '12px' }}>
          <p style={{ color: '#666', marginBottom: '20px', fontSize: '16px' }}>Bạn chưa có đơn hàng nào.</p>
          <button className="btn-primary" style={{ padding: '10px 20px', borderRadius: '6px', border: 'none', cursor: 'pointer' }} onClick={() => navigate('/')}>
            Quay lại mua sắm ngay
          </button>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          {/* Lặp qua từng đơn hàng và hiển thị ra */}
          {orders.map(order => (
            <div key={order.id} style={{ background: 'white', padding: '25px', borderRadius: '12px', boxShadow: '0 4px 15px rgba(0,0,0,0.05)', borderLeft: '5px solid #e50027' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '15px', borderBottom: '1px solid #f1f2f6', paddingBottom: '15px' }}>
                <div>
                  <strong style={{ fontSize: '16px', color: '#2d3436' }}>Mã đơn: #{order.id.slice(-6).toUpperCase()}</strong>
                  <p style={{ fontSize: '13px', color: '#95a5a6', marginTop: '5px' }}>
                    Ngày đặt: {new Date(order.created_at).toLocaleString('vi-VN')}
                  </p>
                </div>
                
                {/* Badge trạng thái đơn hàng có màu sắc thay đổi */}
                <span style={{ 
                  padding: '6px 12px', 
                  borderRadius: '20px', 
                  fontSize: '13px', 
                  fontWeight: 'bold',
                  backgroundColor: order.status === 'Pending' ? '#fff3cd' : (order.status === 'Shipping' ? '#cce5ff' : '#d4edda'),
                  color: order.status === 'Pending' ? '#856404' : (order.status === 'Shipping' ? '#004085' : '#155724')
                }}>
                  {order.status === 'Pending' ? 'Đang chờ duyệt' : order.status}
                </span>
              </div>
              
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '15px', fontSize: '15px' }}>
                <p><strong>Người nhận:</strong> {order.shipping_address}</p>
                <p><strong>Điện thoại:</strong> {order.phone_number}</p>
                <p><strong>Thanh toán:</strong> <span style={{ color: order.payment_method === 'MOMO' ? '#ae2070' : '#e50027', fontWeight: 'bold' }}>{order.payment_method}</span></p>
                <p><strong>Tổng tiền:</strong> <span style={{ color: '#e50027', fontWeight: 'bold', fontSize: '18px' }}>{new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(order.total_amount)}</span></p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default MyOrders;