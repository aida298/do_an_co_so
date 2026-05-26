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
	"go.mongodb.org/mongo-driver/mongo"
)

// AddToCart - API Thêm sản phẩm vào giỏ
func AddToCart(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy thông tin người dùng"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID người dùng không hợp lệ"})
		return
	}

	var input struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	productID, err := primitive.ObjectIDFromHex(input.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sản phẩm không hợp lệ"})
		return
	}

	collection := database.Client.Database("gearpc_db").Collection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingItem models.CartItem
	err = collection.FindOne(ctx, bson.M{"user_id": userID, "product_id": productID}).Decode(&existingItem)

	switch err {
	case mongo.ErrNoDocuments:
		newItem := models.CartItem{
			UserID:    userID,
			ProductID: productID,
			Quantity:  input.Quantity,
			CreatedAt: time.Now(),
		}
		_, err := collection.InsertOne(ctx, newItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi thêm vào giỏ hàng"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Đã thêm sản phẩm vào giỏ hàng!"})

	case nil:
		newQuantity := existingItem.Quantity + input.Quantity
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": existingItem.ID},
			bson.M{"$set": bson.M{"quantity": newQuantity}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật số lượng"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Đã cộng dồn số lượng trong giỏ!"})
	}
}

// GetCart - API Xem giỏ hàng
func GetCart(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

	collection := database.Client.Database("gearpc_db").Collection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cartItems []models.CartItem
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy giỏ hàng"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item models.CartItem
		cursor.Decode(&item)
		cartItems = append(cartItems, item)
	}

	if cartItems == nil {
		cartItems = []models.CartItem{}
	}

	c.JSON(http.StatusOK, cartItems)
}

// RemoveFromCart - API Xóa hẳn 1 sản phẩm khỏi giỏ
func RemoveFromCart(c *gin.Context) {
	cartItemID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(cartItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID giỏ hàng không hợp lệ"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

	collection := database.Client.Database("gearpc_db").Collection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID, "user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa khỏi giỏ hàng"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sản phẩm trong giỏ của bạn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa sản phẩm khỏi giỏ hàng!"})
}

// UpdateCartItem - API Cập nhật số lượng (Dùng cho nút Cộng/Trừ)
func UpdateCartItem(c *gin.Context) {
	cartItemID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(cartItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID giỏ hàng không hợp lệ"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

	var input struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	collection := database.Client.Database("gearpc_db").Collection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Logic tự hủy nếu số lượng <= 0
	if input.Quantity <= 0 {
		_, err = collection.DeleteOne(ctx, bson.M{"_id": objID, "user_id": userID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa khỏi giỏ hàng"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Đã xóa sản phẩm khỏi giỏ vì số lượng bằng 0"})
		return
	}

	// Cập nhật số lượng mới
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": objID, "user_id": userID},
		bson.M{"$set": bson.M{"quantity": input.Quantity}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật số lượng"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sản phẩm trong giỏ của bạn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã cập nhật số lượng thành công!"})
}
