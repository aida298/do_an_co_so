package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CartItem định nghĩa một món hàng nằm trong giỏ của người dùng
type CartItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`       // Mã người dùng (Lấy từ Token)
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"` // Mã sản phẩm
	Quantity  int                `bson:"quantity" json:"quantity"`     // Số lượng
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
