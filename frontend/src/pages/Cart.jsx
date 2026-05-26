import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';

function Cart() {
  const [cartItems, setCartItems] = useState([]);
  const [products, setProducts] = useState([]); 
  const [loading, setLoading] = useState(true);
  const [totalPrice, setTotalPrice] = useState(0);
  const navigate = useNavigate();

  const fetchCartAndProducts = async () => {
    try {
      const prodRes = await api.get('/products/');
      const allProducts = prodRes.data;
      setProducts(allProducts);

      const cartRes = await api.get('/cart/');
      const cartData = cartRes.data;

      const detailedCart = cartData.map(item => {
        const matchedProduct = allProducts.find(p => p.id === item.product_id);
        return {
          ...item,
          product_name: matchedProduct ? matchedProduct.name : "Linh kiện PC",
          product_price: matchedProduct ? matchedProduct.price : 0,
          product_image: matchedProduct && matchedProduct.images?.length > 0 ? matchedProduct.images[0] : 'https://bizweb.dktcdn.net/100/329/122/products/cpu-amd-ryzen-5-7600-4.png'
        };
      });

      setCartItems(detailedCart);
      calculateTotal(detailedCart);
    } catch (error) {
      console.error("Lỗi cập nhật giỏ hàng:", error);
    } finally {
      setLoading(false);
    }
  };

  const calculateTotal = (items) => {
    let total = 0;
    items.forEach(item => {
      total += item.product_price * item.quantity;
    });
    setTotalPrice(total);
  };

  useEffect(() => {
    fetchCartAndProducts();
  }, []);

  const handleUpdateQuantity = async (itemId, currentQty, delta) => {
    const newQty = currentQty + delta;
    if (newQty < 1) return;

    try {
      await api.put(`/cart/${itemId}`, { quantity: newQty });
      fetchCartAndProducts(); 
      window.dispatchEvent(new Event('cartUpdate'));
    } catch (error) {
      alert("Không thể cập nhật số lượng");
    }
  };

  const handleRemoveItem = async (itemId) => {
    if (!window.confirm("Bạn có chắc muốn xóa linh kiện này khỏi giỏ?")) return;
    try {
      await api.delete(`/cart/${itemId}`);
      fetchCartAndProducts();
      window.dispatchEvent(new Event('cartUpdate'));
    } catch (error) {
      alert("Không thể xóa sản phẩm");
    }
  };

  if (loading) return <div className="main-content"><p>Đang kiểm tra giỏ hàng...</p></div>;

  return (
    <div className="main-content">
      <h2>Giỏ Hàng Của Bạn</h2>
      
      {cartItems.length === 0 ? (
        <div style={{ textAlign: 'center', marginTop: '50px' }}>
          <p style={{ color: '#666', marginBottom: '20px' }}>Giỏ hàng của bạn đang trống trơn.</p>
          <button className="btn-primary" style={{ padding: '10px 20px', borderRadius: '6px', border: 'none', cursor: 'pointer' }} onClick={() => navigate('/')}>
            Quay lại mua sắm ngay
          </button>
        </div>
      ) : (
        <div className="cart-page-container" style={{ display: 'flex', gap: '30px', marginTop: '20px' }}>
          <div style={{ flex: 2, display: 'flex', flexDirection: 'column', gap: '15px' }}>
            {cartItems.map((item) => (
              <div key={item.id} style={{ display: 'flex', alignItems: 'center', background: 'white', padding: '15px', box_shadow: '0 2px 8px rgba(0,0,0,0.05)', justifyContent: 'space-between', borderRadius: '12px' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '15px', flex: 1 }}>
                  <img src={item.product_image} alt="" style={{ width: '80px', height: '80px', objectFit: 'contain' }} />
                  <div style={{ paddingRight: '15px' }}>
                    <h4 style={{ fontSize: '15px', marginBottom: '5px', color: '#2d3436' }}>{item.product_name}</h4>
                    <p style={{ fontSize: '12px', color: '#95a5a6', marginBottom: '5px' }}>Mã: {item.product_id}</p>
                    <p style={{ color: '#e50027', fontWeight: 'bold' }}>
                      {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(item.product_price)}
                    </p>
                  </div>
                </div>

                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginRight: '20px' }}>
                  <button 
                    style={{ width: '32px', height: '32px', border: '1px solid #b2bec3', background: '#f8f9fa', cursor: 'pointer', borderRadius: '6px', fontSize: '16px', fontWeight: 'bold', color: '#2d3436', display: 'flex', alignItems: 'center', justifyContent: 'center' }} 
                    onClick={() => handleUpdateQuantity(item.id, item.quantity, -1)}
                  >
                    -
                  </button>
                  <span style={{ fontWeight: 'bold', width: '25px', textAlign: 'center', fontSize: '15px' }}>{item.quantity}</span>
                  <button 
                    style={{ width: '32px', height: '32px', border: '1px solid #b2bec3', background: '#f8f9fa', cursor: 'pointer', borderRadius: '6px', fontSize: '16px', fontWeight: 'bold', color: '#2d3436', display: 'flex', alignItems: 'center', justifyContent: 'center' }} 
                    onClick={() => handleUpdateQuantity(item.id, item.quantity, 1)}
                  >
                    +
                  </button>
                </div>

                <button style={{ background: 'none', border: 'none', color: '#666', cursor: 'pointer', fontSize: '18px', padding: '10px' }} onClick={() => handleRemoveItem(item.id)}>❌</button>
              </div>
            ))}
          </div>

          <div style={{ flex: 1, background: 'white', padding: '25px', borderRadius: '12px', boxShadow: '0 2px 12px rgba(0,0,0,0.08)', height: 'fit-content' }}>
            <h3 style={{ marginBottom: '20px', borderBottom: '1px solid #eee', paddingBottom: '10px' }}>Tóm tắt đơn hàng</h3>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px', fontSize: '16px' }}>
              <span>Thành tiền:</span>
              <span style={{ color: '#e50027', fontWeight: 'bold', fontSize: '20px' }}>
                {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(totalPrice)}
              </span>
            </div>
            <p style={{ fontSize: '13px', color: '#666', marginBottom: '20px' }}>* Miễn phí vận chuyển toàn quốc cho mọi đơn hàng linh kiện.</p>
            <button 
              className="btn-primary" 
              style={{ width: '100%', padding: '14px', border: 'none', borderRadius: '8px', fontSize: '16px', fontWeight: 'bold', cursor: 'pointer' }}
              onClick={() => navigate('/checkout')}
            >
              Tiến Hành Đặt Hàng
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default Cart;