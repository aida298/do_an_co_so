package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product định nghĩa thông tin linh kiện PC chi tiết (Chuẩn Web Gaming)
type Product struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Brand      string             `bson:"brand" json:"brand"`             // Thương hiệu (VD: HyperX, ASUS)
	ShortDesc  string             `bson:"short_desc" json:"short_desc"`   // Mô tả ngắn gọn (Hiển thị ở Card sản phẩm)
	DetailDesc string             `bson:"detail_desc" json:"detail_desc"` // Bài viết mô tả dài (Có thể chứa mã HTML)
	Price      float64            `bson:"price" json:"price"`             // Giá bán hiện tại (Giá đỏ)
	OldPrice   float64            `bson:"old_price" json:"old_price"`     // Giá cũ (Để gạch ngang)
	Stock      int                `bson:"stock" json:"stock"`
	CategoryID primitive.ObjectID `bson:"category_id" json:"category_id"`
	Images     []string           `bson:"images" json:"images"` // Mảng chứa nhiều link ảnh
	Specs      map[string]string  `bson:"specs" json:"specs"`   // Bảng thông số kỹ thuật động
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
