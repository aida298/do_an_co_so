import React, { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Link, useNavigate } from 'react-router-dom';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Cart from './pages/Cart';
import Checkout from './pages/Checkout';
import MyOrders from './pages/MyOrders'; // 1. IMPORT TRANG LỊCH SỬ ĐƠN HÀNG
import api from './services/api';
import './App.css';

function Navigation() {
  const navigate = useNavigate();
  const token = localStorage.getItem('token');
  const [cartCount, setCartCount] = useState(0);

  const getCartCount = async () => {
    if (!token) return;
    try {
      const response = await api.get('/cart/');
      const totalItems = response.data.reduce((sum, item) => sum + item.quantity, 0);
      setCartCount(totalItems);
    } catch (error) {
      console.error("Không thể lấy số lượng giỏ hàng");
    }
  };

  useEffect(() => {
    getCartCount();
    window.addEventListener('cartUpdate', getCartCount);
    return () => {
      window.removeEventListener('cartUpdate', getCartCount);
    };
  }, [token]);

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/');
    window.location.reload();
  };

  return (
    <header className="header">
      <div className="header-container">
        <Link to="/" className="logo">
          <h1>GearPC</h1>
        </Link>
        
        <div className="header-actions">
          <div className="cart-icon" onClick={() => navigate('/cart')}>
            <span>🛒 Giỏ hàng</span>
            <span className="cart-badge">{cartCount}</span>
          </div>
          
          {token ? (
            <div style={{ display: 'flex', gap: '15px', alignItems: 'center' }}>
              {/* 2. THÊM NÚT XEM ĐƠN HÀNG Ở ĐÂY */}
              <Link to="/my-orders" style={{ color: '#2d3436', textDecoration: 'none', fontWeight: 'bold', borderRight: '1px solid #dfe6e9', paddingRight: '15px' }}>
                📦 Đơn hàng
              </Link>
              <button onClick={handleLogout} className="btn-auth btn-logout">Đăng xuất</button>
            </div>
          ) : (
            <div className="auth-buttons">
              <Link to="/login" className="btn-auth btn-login">Đăng nhập</Link>
              <Link to="/register" className="btn-auth btn-register">Đăng ký</Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}

function App() {
  return (
    <BrowserRouter>
      <div className="app-container">
        <Navigation />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/cart" element={<Cart />} />
          <Route path="/checkout" element={<Checkout />} />
          <Route path="/my-orders" element={<MyOrders />} /> {/* 3. THÊM ĐƯỜNG DẪN ROUTE */}
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;