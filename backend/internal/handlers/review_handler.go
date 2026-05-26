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

// AddReview - Thêm đánh giá cho sản phẩm
func AddReview(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

	var input struct {
		ProductID string `json:"product_id" binding:"required"`
		Rating    int    `json:"rating" binding:"required"`
		Comment   string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// Chặn đánh giá vớ vẩn (chỉ cho phép 1-5 sao)
	if input.Rating < 1 || input.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Điểm đánh giá phải từ 1 đến 5 sao"})
		return
	}

	productID, err := primitive.ObjectIDFromHex(input.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã sản phẩm không hợp lệ"})
		return
	}

	review := models.Review{
		ID:        primitive.NewObjectID(),
		ProductID: productID,
		UserID:    userID,
		Rating:    input.Rating,
		Comment:   input.Comment,
		CreatedAt: time.Now(),
	}

	collection := database.Client.Database("gearpc_db").Collection("reviews")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu đánh giá"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Cảm ơn bạn đã đánh giá sản phẩm!"})
}

// GetProductReviews - Lấy tất cả đánh giá của 1 sản phẩm
func GetProductReviews(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã sản phẩm không hợp lệ"})
		return
	}

	collection := database.Client.Database("gearpc_db").Collection("reviews")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"product_id": productID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy đánh giá"})
		return
	}

	var reviews []models.Review
	cursor.All(ctx, &reviews)

	if reviews == nil {
		reviews = []models.Review{}
	}

	c.JSON(http.StatusOK, reviews)
}
