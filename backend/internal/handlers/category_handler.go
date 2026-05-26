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

// CreateCategory xử lý API Thêm danh mục mới
func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	collection := database.Client.Database("gearpc_db").Collection("categories")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu danh mục"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo danh mục thành công",
		"id":      result.InsertedID,
	})
}

// GetCategories xử lý API Lấy toàn bộ danh mục
func GetCategories(c *gin.Context) {
	collection := database.Client.Database("gearpc_db").Collection("categories")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var categories []models.Category
	// Tìm tất cả (không có điều kiện lọc)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi truy vấn dữ liệu"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var category models.Category
		cursor.Decode(&category)
		categories = append(categories, category)
	}

	// Xử lý trường hợp database chưa có danh mục nào thì trả về mảng rỗng [] thay vì null
	if categories == nil {
		categories = []models.Category{}
	}

	c.JSON(http.StatusOK, categories)
}
