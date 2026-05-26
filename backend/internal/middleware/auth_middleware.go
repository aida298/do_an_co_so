package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"do_an_co_so_backend/internal/database"
	"do_an_co_so_backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Lưu ý: Chuỗi Secret Key này phải GIỐNG HỆT với chuỗi bạn đang dùng bên file auth_handler.go (lúc phát hành Token)
var jwtSecret = []byte("my_super_secret_key")

// RequireAuth - Chốt chặn bảo vệ: Bắt buộc Đăng nhập & Không bị khóa
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy Token từ Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Yêu cầu đăng nhập (Thiếu Token)!"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Giải mã Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ hoặc đã hết hạn!"})
			c.Abort()
			return
		}

		// 3. Lấy thông tin user_id từ Token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Lỗi phân tích Token"})
			c.Abort()
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Dữ liệu Token không đúng cấu trúc"})
			c.Abort()
			return
		}

		// --- ĐOẠN CODE NÂNG CẤP: KIỂM TRA TRẠNG THÁI KHÓA TÀI KHOẢN TỪ DATABASE ---
		objID, err := primitive.ObjectIDFromHex(userIDStr)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			userCollection := database.Client.Database("gearpc_db").Collection("users")
			var user models.User
			err := userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)

			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy tài khoản trong hệ thống!"})
				c.Abort()
				return
			}

			// NẾU BỊ KHÓA -> CHẶN NGAY LẬP TỨC TRẢ VỀ 403 FORBIDDEN
			if user.IsBlocked {
				c.JSON(http.StatusForbidden, gin.H{"error": "Tài khoản của bạn đã bị khóa! Vui lòng liên hệ Admin."})
				c.Abort()
				return
			}

			// Tiện tay nhét luôn role vào Context để lát nữa hàm RequireAdmin xài, đỡ phải truy vấn DB 2 lần
			c.Set("role", user.Role)
		}
		// -------------------------------------------------------------------------

		// 4. Mọi thứ hợp lệ -> Lưu user_id vào Context và cho đi tiếp
		c.Set("user_id", userIDStr)
		c.Next()
	}
}

// RequireAdmin - Chốt chặn phân quyền: Chỉ cho phép Admin đi qua
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Rút quyền (role) đã được hàm RequireAuth gài sẵn vào Context
		role, exists := c.Get("role")

		if !exists || role.(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền truy cập! (Tính năng dành riêng cho Admin)"})
			c.Abort()
			return
		}

		c.Next() // Đúng là Admin thì cho đi tiếp vào xử lý API
	}
}
