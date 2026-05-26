import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../services/api';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      const response = await api.post('/auth/login', { email, password });
      localStorage.setItem('token', response.data.token);
      navigate('/');
      window.location.reload(); // Load lại trang để cập nhật Header
    } catch (error) {
      alert("Lỗi: " + (error.response?.data?.error || "Đăng nhập thất bại"));
    }
  };

  return (
    <div className="auth-wrapper">
      <div className="auth-box">
        <h2>Đăng Nhập</h2>
        <p className="auth-subtitle">Mừng bạn quay lại với GearPC</p>
        
        <form onSubmit={handleLogin} className="auth-form">
          <div className="input-group">
            <label>Email</label>
            <input type="email" placeholder="Nhập email của bạn" value={email} onChange={(e) => setEmail(e.target.value)} required />
          </div>
          <div className="input-group">
            <label>Mật khẩu</label>
            <input type="password" placeholder="Nhập mật khẩu" value={password} onChange={(e) => setPassword(e.target.value)} required />
          </div>
          <button type="submit" className="btn-primary auth-submit">Đăng nhập</button>
        </form>

        <p className="auth-switch">
          Chưa có tài khoản? <Link to="/register">Đăng ký tài khoản mới</Link>
        </p>
      </div>
    </div>
  );
}

export default Login;