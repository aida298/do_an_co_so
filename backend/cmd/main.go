package main

import (
	"do_an_co_so_backend/internal/database"
	"do_an_co_so_backend/internal/handlers"
	"do_an_co_so_backend/internal/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()
	router := gin.Default()

	// --- CẤU HÌNH CORS CHO PHÉP REACTJS GỌI API ---
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API Test
	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Backend Golang đang chạy ngon lành!",
		})
	})

	// --- 1. AUTH (XÁC THỰC) ---
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", handlers.Register)
		authRoutes.POST("/login", handlers.Login)
	}

	// --- 2. CATEGORY (DANH MỤC) ---
	categoryRoutes := router.Group("/api/categories")
	{
		categoryRoutes.GET("/", handlers.GetCategories)
		categoryRoutes.POST("/", middleware.RequireAuth(), middleware.RequireAdmin(), handlers.CreateCategory)
	}

	// --- 3. PRODUCT (SẢN PHẨM) ---
	productRoutes := router.Group("/api/products")
	{
		productRoutes.GET("/", handlers.GetProducts)
		productRoutes.POST("/", middleware.RequireAuth(), middleware.RequireAdmin(), handlers.CreateProduct)
		productRoutes.PUT("/:id", middleware.RequireAuth(), middleware.RequireAdmin(), handlers.UpdateProduct)
		productRoutes.DELETE("/:id", middleware.RequireAuth(), middleware.RequireAdmin(), handlers.DeleteProduct)
	}

	// --- 4. CART (GIỎ HÀNG) ---
	cartRoutes := router.Group("/api/cart")
	{
		cartRoutes.POST("/", middleware.RequireAuth(), handlers.AddToCart)
		cartRoutes.GET("/", middleware.RequireAuth(), handlers.GetCart)
		cartRoutes.PUT("/:id", middleware.RequireAuth(), handlers.UpdateCartItem)
		cartRoutes.DELETE("/:id", middleware.RequireAuth(), handlers.RemoveFromCart)
	}

	// --- 5. ORDER (ĐẶT HÀNG) ---
	orderRoutes := router.Group("/api/orders")
	{
		// Của khách hàng
		orderRoutes.POST("/", middleware.RequireAuth(), handlers.Checkout)
		orderRoutes.GET("/my-orders", middleware.RequireAuth(), handlers.GetMyOrders)

		// Của Admin (Quản lý duyệt đơn)
		orderRoutes.GET("/", middleware.RequireAuth(), middleware.RequireAdmin(), handlers.GetAllOrders)
		orderRoutes.PUT("/:id/status", middleware.RequireAuth(), middleware.RequireAdmin(), handlers.UpdateOrderStatus)
	}

	// Tuyến đường giả lập cổng thanh toán xác thực (Không cần RequireAuth)
	router.POST("/api/payments/momo-webhook", handlers.MomoWebhookSimulate)

	// --- 6. REVIEWS (ĐÁNH GIÁ SẢN PHẨM) ---
	reviewRoutes := router.Group("/api/reviews")
	{
		// Ai cũng xem được đánh giá
		reviewRoutes.GET("/:product_id", handlers.GetProductReviews)
		// Chỉ người dùng đăng nhập mới được đánh giá
		reviewRoutes.POST("/", middleware.RequireAuth(), handlers.AddReview)
	}

	// --- 7. ADMIN DASHBOARD (THỐNG KÊ DOANH THU) ---
	dashboardRoutes := router.Group("/api/dashboard")
	{
		// Chỉ Admin mới được xem các con số mật này
		dashboardRoutes.GET("/", middleware.RequireAuth(), middleware.RequireAdmin(), handlers.GetDashboardStats)
	}

	// --- 8. USER MANAGEMENT (QUẢN LÝ NGƯỜI DÙNG) ---
	userRoutes := router.Group("/api/users")
	{
		// Cả 2 quyền này đều bắt buộc phải là Admin
		userRoutes.GET("/", middleware.RequireAuth(), middleware.RequireAdmin(), handlers.GetAllUsers)
		userRoutes.PUT("/:id/block", middleware.RequireAuth(), middleware.RequireAdmin(), handlers.ToggleBlockUser)
	}

	router.Run(":8080")
}
