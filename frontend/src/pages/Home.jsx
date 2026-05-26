import React, { useState, useEffect } from 'react';
import api from '../services/api';

function Home() {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);

  // CÁC STATE DÙNG ĐỂ LỌC VÀ SẮP XẾP
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('');
  const [sortOption, setSortOption] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        // 1. Lấy danh sách sản phẩm
        const prodRes = await api.get('/products/');
        setProducts(prodRes.data || []);

        // 2. Lấy danh sách danh mục (để đưa vào ô Select)
        try {
          const catRes = await api.get('/categories/');
          setCategories(catRes.data || []);
        } catch (error) {
          console.error("Chưa có danh mục nào từ Backend");
        }
      } catch (error) {
        console.error("Lỗi lấy dữ liệu:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  // Gọi API thêm vào giỏ hàng
  const handleAddToCart = async (product) => {
    const token = localStorage.getItem('token');
    if (!token) {
      alert("Vui lòng đăng nhập để tiến hành mua hàng!");
      return;
    }
    try {
      await api.post('/cart/', { product_id: product.id, quantity: 1 });
      alert(`Đã thêm thành công ${product.name} vào giỏ hàng!`);
      window.dispatchEvent(new Event('cartUpdate'));
    } catch (error) {
      alert("Không thể thêm vào giỏ: " + (error.response?.data?.error || "Lỗi hệ thống"));
    }
  };

  // ==========================================
  // LOGIC BỘ LỌC VÀ SẮP XẾP THÔNG MINH
  // ==========================================
  let displayedProducts = [...products];

  // 1. Lọc theo tên (Tìm kiếm)
  if (searchTerm) {
    displayedProducts = displayedProducts.filter(p => 
      p.name.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }

  // 2. Lọc theo Danh mục
  if (selectedCategory) {
    displayedProducts = displayedProducts.filter(p => p.category_id === selectedCategory);
  }

  // 3. Sắp xếp (Tăng/Giảm giá, A-Z)
  if (sortOption === 'price_asc') {
    displayedProducts.sort((a, b) => a.price - b.price);
  } else if (sortOption === 'price_desc') {
    displayedProducts.sort((a, b) => b.price - a.price);
  } else if (sortOption === 'name_asc') {
    displayedProducts.sort((a, b) => a.name.localeCompare(b.name));
  } else if (sortOption === 'name_desc') {
    displayedProducts.sort((a, b) => b.name.localeCompare(a.name));
  }

  return (
    <main className="main-content">
      
      {/* THANH CÔNG CỤ TÌM KIẾM & LỌC */}
      <div className="filter-bar">
        
        {/* BÊN TRÁI: Nút Danh mục 3 gạch */}
        <select 
          value={selectedCategory} 
          onChange={(e) => setSelectedCategory(e.target.value)}
          className="filter-category"
        >
          <option value="">☰ Danh mục</option>
          {categories.map(cat => (
            <option key={cat.id} value={cat.id}>{cat.name}</option>
          ))}
        </select>

        {/* Ở GIỮA: Ô tìm kiếm to, dài, bo tròn */}
        <input 
          type="text" 
          placeholder="🔍 Tìm kiếm tên linh kiện..." 
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="filter-input-center"
        />

        {/* BÊN PHẢI: Ô Sắp xếp thu bé lại */}
        <select 
          value={sortOption} 
          onChange={(e) => setSortOption(e.target.value)}
          className="filter-sort"
        >
          <option value="">↕️ Sắp xếp</option>
          <option value="price_asc">Giá: Thấp - Cao</option>
          <option value="price_desc">Giá: Cao - Thấp</option>
          <option value="name_asc">Tên: A - Z</option>
          <option value="name_desc">Tên: Z - A</option>
        </select>
      </div>

      <h2>Sản phẩm của chúng tôi</h2>

      {loading ? (
        <p>Đang tải dữ liệu từ Backend Golang...</p>
      ) : (
        <>
          {displayedProducts.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '50px', background: 'white', borderRadius: '8px' }}>
              <h3>Không tìm thấy sản phẩm nào phù hợp 😢</h3>
            </div>
          ) : (
            <div className="product-grid">
              {displayedProducts.map((product) => (
                <div key={product.id} className="product-card">
                  <img 
                    src={product.images && product.images.length > 0 && product.images[0].startsWith('http') ? product.images[0] : 'https://bizweb.dktcdn.net/100/329/122/products/cpu-amd-ryzen-5-7600-4.png'} 
                    alt={product.name} 
                    className="product-img"
                  />
                  <h3 className="product-name">{product.name}</h3>
                  <p className="product-price">
                    {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(product.price)}
                  </p>
                  <button className="btn-add-cart" onClick={() => handleAddToCart(product)}>
                    Thêm vào giỏ
                  </button>
                </div>
              ))}
            </div>
          )}
        </>
      )}
    </main>
  );
}

export default Home;