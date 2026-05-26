import axios from 'axios';

const api = axios.create({
    baseURL: 'http://localhost:8080/api', // Link Backend Golang của bạn
    timeout: 10000,
});

// Tự động gắn Token vào thẻ bảo vệ (Header) nếu khách đã đăng nhập
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export default api;