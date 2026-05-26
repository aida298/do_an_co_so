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

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu sản phẩm không hợp lệ"})
		return
	}
	product.CreatedAt = time.Now()

	collection := database.Client.Database("gearpc_db").Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu sản phẩm vào kho"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Thêm linh kiện thành công",
		"id":      result.InsertedID,
	})
}

func GetProducts(c *gin.Context) {
	collection := database.Client.Database("gearpc_db").Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var products []models.Product
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi truy vấn dữ liệu sản phẩm"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product models.Product
		cursor.Decode(&product)
		products = append(products, product)
	}

	if products == nil {
		products = []models.Product{}
	}
	c.JSON(http.StatusOK, products)
}

func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sản phẩm không hợp lệ"})
		return
	}

	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu cập nhật không hợp lệ"})
		return
	}

	collection := database.Client.Database("gearpc_db").Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Đã bổ sung đầy đủ các trường mới vào mảng update
	updateData := bson.M{
		"$set": bson.M{
			"name":        input.Name,
			"brand":       input.Brand,
			"short_desc":  input.ShortDesc,
			"detail_desc": input.DetailDesc,
			"price":       input.Price,
			"old_price":   input.OldPrice,
			"stock":       input.Stock,
			"category_id": input.CategoryID,
			"images":      input.Images,
			"specs":       input.Specs,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật sản phẩm"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật thông tin linh kiện thành công!"})
}

func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sản phẩm không hợp lệ"})
		return
	}

	collection := database.Client.Database("gearpc_db").Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa sản phẩm"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sản phẩm để xóa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa sản phẩm khỏi hệ thống thành công!"})
}
