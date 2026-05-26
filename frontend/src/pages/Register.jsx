import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../services/api';

function Register() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleRegister = async (e) => {
    e.preventDefault();
    try {
      await api.post('/auth/register', { name, email, password });
      alert("Tuyệt vời! Đăng ký thành công. Vui lòng đăng nhập.");
      navigate('/login'); // Chuyển hướng sang trang đăng nhập
    } catch (error) {
      alert("Lỗi: " + (error.response?.data?.error || "Đăng ký thất bại"));
    }
  };

  return (
    <div className="auth-wrapper">
      <div className="auth-box">
        <h2>Tạo tài khoản mới</h2>
        <p className="auth-subtitle">Tham gia GearPC để nhận ưu đãi</p>
        
        <form onSubmit={handleRegister} className="auth-form">
          <div className="input-group">
            <label>Họ và tên</label>
            <input type="text" placeholder="Nhập tên của bạn" value={name} onChange={(e) => setName(e.target.value)} required />
          </div>
          <div className="input-group">
            <label>Email</label>
            <input type="email" placeholder="Ví dụ: hotro@gearpc.com" value={email} onChange={(e) => setEmail(e.target.value)} required />
          </div>
          <div className="input-group">
            <label>Mật khẩu</label>
            <input type="password" placeholder="Tối thiểu 6 ký tự" value={password} onChange={(e) => setPassword(e.target.value)} required />
          </div>
          <button type="submit" className="btn-primary auth-submit">Đăng ký ngay</button>
        </form>
        
        <p className="auth-switch">
          Đã có tài khoản? <Link to="/login">Đăng nhập tại đây</Link>
        </p>
      </div>
    </div>
  );
}

export default Register;