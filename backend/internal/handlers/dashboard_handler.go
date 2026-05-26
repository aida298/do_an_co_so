package handlers

import (
	"context"
	"net/http"
	"time"

	"do_an_co_so_backend/internal/database"
	"do_an_co_so_backend/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetDashboardStats - Thống kê tổng quan cho Admin
func GetDashboardStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := database.Client.Database("gearpc_db")
	orderCollection := db.Collection("orders")
	userCollection := db.Collection("users")

	// 1. Đếm tổng số lượng người dùng
	totalUsers, _ := userCollection.CountDocuments(ctx, bson.M{})

	// 2. Lấy toàn bộ đơn hàng để tính toán
	cursor, err := orderCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi truy xuất dữ liệu đơn hàng"})
		return
	}

	var orders []models.Order
	cursor.All(ctx, &orders)

	var totalRevenue float64 = 0
	var totalOrders int = len(orders)
	var pendingOrders int = 0

	// 3. Vòng lặp quét doanh thu và đếm đơn chưa duyệt
	for _, order := range orders {
		// Chỉ cộng tiền những đơn đã thanh toán hoặc đã hoàn thành
		if order.PaymentStatus == "Paid" || order.Status == "Completed" {
			totalRevenue += order.TotalAmount
		}

		if order.Status == "Pending" {
			pendingOrders++
		}
	}

	// 4. Trả về một cục JSON gọn gàng cho giao diện biểu đồ
	c.JSON(http.StatusOK, gin.H{
		"total_users":    totalUsers,
		"total_orders":   totalOrders,
		"pending_orders": pendingOrders, // Để báo chấm đỏ (thông báo) cho Admin duyệt
		"total_revenue":  totalRevenue,
	})
}
