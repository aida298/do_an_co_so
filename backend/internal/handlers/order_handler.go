package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"do_an_co_so_backend/internal/database"
	"do_an_co_so_backend/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Checkout - API Đặt hàng (Hỗ trợ cả COD và MOMO)
func Checkout(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

	var input struct {
		ShippingAddress string `json:"shipping_address" binding:"required"`
		PhoneNumber     string `json:"phone_number" binding:"required"`
		PaymentMethod   string `json:"payment_method" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng điền đầy đủ thông tin"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cartCollection := database.Client.Database("gearpc_db").Collection("carts")
	productCollection := database.Client.Database("gearpc_db").Collection("products")
	orderCollection := database.Client.Database("gearpc_db").Collection("orders")

	// 1. Lấy hàng từ giỏ
	cursor, err := cartCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy giỏ hàng"})
		return
	}
	var cartItems []models.CartItem
	cursor.All(ctx, &cartItems)

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Giỏ hàng rỗng"})
		return
	}

	// 2. Tính tổng tiền
	var orderItems []models.OrderItem
	var totalAmount float64
	for _, item := range cartItems {
		var product models.Product
		err := productCollection.FindOne(ctx, bson.M{"_id": item.ProductID}).Decode(&product)
		if err != nil {
			continue
		}
		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ID,
			Name:      product.Name,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
		totalAmount += product.Price * float64(item.Quantity)
	}

	// 3. Tạo cấu trúc đơn hàng mới
	newOrder := models.Order{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		Items:           orderItems,
		TotalAmount:     totalAmount,
		ShippingAddress: input.ShippingAddress,
		PhoneNumber:     input.PhoneNumber,
		PaymentMethod:   input.PaymentMethod,
		PaymentStatus:   "Unpaid",
		Status:          "Pending",
		CreatedAt:       time.Now(),
	}

	_, err = orderCollection.InsertOne(ctx, newOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi tạo đơn hàng"})
		return
	}

	// 4. Xử lý theo phương thức thanh toán
	if input.PaymentMethod == "MOMO" {
		yourMomoPhone := "0384253830" // Thay số của bạn vào đây
		momoLink := fmt.Sprintf("https://nhantien.momo.vn/%s/%d", yourMomoPhone, int(totalAmount))
		qrCodeURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=%s", url.QueryEscape(momoLink))

		c.JSON(http.StatusCreated, gin.H{
			"message":        "Đơn hàng tạm đã mở! Vui lòng quét mã QR thanh toán để hoàn tất.",
			"order_id":       newOrder.ID,
			"payment_method": "MOMO",
			"total_amount":   totalAmount,
			"qr_code_url":    qrCodeURL,
		})
		return
	}

	// Nếu chọn COD: XÓA GIỎ HÀNG NGAY LẬP TỨC
	_, _ = cartCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	c.JSON(http.StatusCreated, gin.H{
		"message":        "Đặt hàng thành công! Bạn sẽ thanh toán khi nhận hàng (COD).",
		"order_id":       newOrder.ID,
		"payment_method": "COD",
	})
}

// MomoWebhookSimulate - Giả lập MoMo
func MomoWebhookSimulate(c *gin.Context) {
	var input struct {
		OrderID string `json:"order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu mã order_id"})
		return
	}

	orderID, err := primitive.ObjectIDFromHex(input.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã đơn hàng không hợp lệ"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderCollection := database.Client.Database("gearpc_db").Collection("orders")
	cartCollection := database.Client.Database("gearpc_db").Collection("carts")

	var order models.Order
	err = orderCollection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đơn hàng này"})
		return
	}

	_, err = orderCollection.UpdateOne(
		ctx,
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{"payment_status": "Paid", "status": "Pending"}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật trạng thái đơn"})
		return
	}

	_, err = cartCollection.DeleteMany(ctx, bson.M{"user_id": order.UserID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi dọn dẹp giỏ hàng"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "MoMo xác nhận: Đã nhận đủ tiền!"})
}

// GetMyOrders - Xem lịch sử mua hàng cá nhân
func GetMyOrders(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	orderCollection := database.Client.Database("gearpc_db").Collection("orders")

	cursor, err := orderCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy danh sách đơn"})
		return
	}
	var orders []models.Order
	cursor.All(ctx, &orders)

	if orders == nil {
		orders = []models.Order{}
	}
	c.JSON(http.StatusOK, orders)
}

// --- CÁC HÀM CỦA ADMIN ĐÃ ĐƯỢC BỔ SUNG LẠI Ở ĐÂY ---

// GetAllOrders - Xem tất cả đơn hàng (Admin)
func GetAllOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	orderCollection := database.Client.Database("gearpc_db").Collection("orders")

	cursor, err := orderCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi truy vấn danh sách đơn hàng"})
		return
	}

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi đọc dữ liệu đơn hàng"})
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateOrderStatus - Cập nhật trạng thái đơn hàng (Admin)
func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID đơn hàng không hợp lệ"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng truyền lên trạng thái (status)"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	orderCollection := database.Client.Database("gearpc_db").Collection("orders")

	result, err := orderCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": input.Status}},
	)

	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật trạng thái đơn hàng"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã cập nhật trạng thái đơn hàng thành: " + input.Status})
}
