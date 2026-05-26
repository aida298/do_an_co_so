import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';

function Checkout() {
  const [address, setAddress] = useState('');
  const [phone, setPhone] = useState('');
  const [paymentMethod, setPaymentMethod] = useState('COD');
  
  // State để lưu thông tin MoMo hiển thị mã QR
  const [qrUrl, setQrUrl] = useState('');
  const [orderId, setOrderId] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);
  
  const navigate = useNavigate();

  const handleCheckout = async (e) => {
    e.preventDefault();
    setIsProcessing(true);
    
    try {
      const response = await api.post('/orders/', {
        shipping_address: address,
        phone_number: phone,
        payment_method: paymentMethod
      });

      if (paymentMethod === 'MOMO') {
        // Nếu là MoMo, lưu lại link QR và OrderID để hiện ra màn hình
        setQrUrl(response.data.qr_code_url);
        setOrderId(response.data.order_id);
      } else {
        // Nếu là COD, báo thành công và đá về trang chủ
        alert("🎉 " + response.data.message);
        window.dispatchEvent(new Event('cartUpdate')); // Reset số lượng giỏ hàng
        navigate('/');
      }
    } catch (error) {
      alert("Lỗi đặt hàng: " + (error.response?.data?.error || "Vui lòng thử lại"));
    } finally {
      setIsProcessing(false);
    }
  };

  // Nút giả lập khách hàng đã cầm điện thoại quét xong MoMo (Dùng để test)
  const handleSimulateMoMoPayment = async () => {
    try {
      await api.post('/payments/momo-webhook', { order_id: orderId });
      alert("🎉 MoMo báo về: Thanh toán thành công! Đơn hàng đang được chờ duyệt.");
      window.dispatchEvent(new Event('cartUpdate'));
      navigate('/');
    } catch (error) {
      alert("Lỗi xác nhận thanh toán");
    }
  };

  // NẾU ĐANG CÓ MÃ QR THÌ HIỆN MÀN HÌNH QUÉT MÃ
  if (qrUrl) {
    return (
      <div className="main-content" style={{ display: 'flex', justifyContent: 'center', marginTop: '40px' }}>
        <div style={{ background: 'white', padding: '40px', borderRadius: '16px', boxShadow: '0 10px 30px rgba(0,0,0,0.1)', textAlign: 'center', maxWidth: '450px' }}>
          <h2 style={{ color: '#ae2070', marginBottom: '10px' }}>Thanh toán qua MoMo</h2>
          <p style={{ color: '#666', marginBottom: '20px' }}>Vui lòng mở ứng dụng MoMo và quét mã QR dưới đây để hoàn tất đơn hàng.</p>
          
          <div style={{ background: '#f8f9fa', padding: '20px', borderRadius: '12px', display: 'inline-block', marginBottom: '25px' }}>
            <img src={qrUrl} alt="MoMo QR Code" style={{ width: '250px', height: '250px', borderRadius: '8px' }} />
          </div>
          
          {/* Nút giả lập (Vì chúng ta code ở localhost nên phải có nút này để báo Backend là đã quét xong) */}
          <button 
            onClick={handleSimulateMoMoPayment}
            style={{ width: '100%', padding: '15px', backgroundColor: '#ae2070', color: 'white', border: 'none', borderRadius: '8px', fontSize: '16px', fontWeight: 'bold', cursor: 'pointer' }}
          >
            Đã quét mã & Chuyển tiền xong
          </button>
        </div>
      </div>
    );
  }

  // MÀN HÌNH NHẬP THÔNG TIN GIAO HÀNG (Mặc định)
  return (
    <div className="main-content" style={{ display: 'flex', justifyContent: 'center', marginTop: '20px' }}>
      <div style={{ background: 'white', width: '100%', maxWidth: '600px', padding: '30px', borderRadius: '12px', boxShadow: '0 4px 15px rgba(0,0,0,0.05)' }}>
        <h2 style={{ marginBottom: '25px', borderBottom: '2px solid #f1f2f6', paddingBottom: '10px' }}>Thông tin Giao hàng</h2>
        
        <form onSubmit={handleCheckout} style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          <div className="input-group">
            <label>Địa chỉ nhận hàng</label>
            <input 
              type="text" 
              placeholder="Ví dụ: 123 Đường Điện Biên Phủ, Quận Bình Thạnh, TP.HCM" 
              value={address} 
              onChange={(e) => setAddress(e.target.value)} 
              required 
            />
          </div>
          
          <div className="input-group">
            <label>Số điện thoại liên hệ</label>
            <input 
              type="text" 
              placeholder="Nhập số điện thoại của bạn" 
              value={phone} 
              onChange={(e) => setPhone(e.target.value)} 
              required 
            />
          </div>

          <div style={{ marginTop: '10px' }}>
            <label style={{ fontWeight: 'bold', display: 'block', marginBottom: '15px' }}>Phương thức thanh toán:</label>
            
            <div style={{ display: 'flex', gap: '20px' }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer', padding: '15px', border: paymentMethod === 'COD' ? '2px solid #e50027' : '1px solid #ddd', borderRadius: '8px', flex: 1 }}>
                <input 
                  type="radio" 
                  name="payment" 
                  value="COD" 
                  checked={paymentMethod === 'COD'} 
                  onChange={(e) => setPaymentMethod(e.target.value)} 
                />
                Thanh toán khi nhận hàng (COD)
              </label>

              <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer', padding: '15px', border: paymentMethod === 'MOMO' ? '2px solid #ae2070' : '1px solid #ddd', borderRadius: '8px', flex: 1 }}>
                <input 
                  type="radio" 
                  name="payment" 
                  value="MOMO" 
                  checked={paymentMethod === 'MOMO'} 
                  onChange={(e) => setPaymentMethod(e.target.value)} 
                />
                Ví điện tử MoMo
              </label>
            </div>
          </div>

          <button type="submit" className="btn-primary" style={{ padding: '15px', fontSize: '16px', borderRadius: '8px', marginTop: '10px' }} disabled={isProcessing}>
            {isProcessing ? 'Đang xử lý...' : 'Xác nhận Đặt Hàng'}
          </button>
        </form>
      </div>
    </div>
  );
}

export default Checkout;