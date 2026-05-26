package handlers

import (
	"context"
	"net/http"
	"time"

	"do_an_co_so_backend/internal/database"
	"do_an_co_so_backend/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllUsers - Lấy danh sách toàn bộ khách hàng (Admin)
func GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCollection := database.Client.Database("gearpc_db").Collection("users")

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy danh sách người dùng"})
		return
	}

	var users []models.User
	cursor.All(ctx, &users)

	if users == nil {
		users = []models.User{}
	}

	c.JSON(http.StatusOK, users)
}

// ToggleBlockUser - Khóa hoặc Mở khóa tài khoản (Admin)
func ToggleBlockUser(c *gin.Context) {
	userIDStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID người dùng không hợp lệ"})
		return
	}

	// Admin gửi lên boolean: true (khóa) hoặc false (mở khóa)
	var input struct {
		IsBlocked bool `json:"is_blocked"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCollection := database.Client.Database("gearpc_db").Collection("users")

	// Không cho phép Admin tự khóa chính mình
	adminIDStr, _ := c.Get("user_id")
	if userIDStr == adminIDStr.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không thể tự khóa tài khoản của chính mình!"})
		return
	}

	result, err := userCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"is_blocked": input.IsBlocked}},
	)

	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật trạng thái người dùng"})
		return
	}

	statusMsg := "đã bị khóa"
	if !input.IsBlocked {
		statusMsg = "đã được mở khóa"
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tài khoản người dùng " + statusMsg + " thành công!"})
}
