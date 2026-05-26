package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	Name      string             `bson:"name" json:"name"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"`
}

type Order struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	Items           []OrderItem        `bson:"items" json:"items"`
	TotalAmount     float64            `bson:"total_amount" json:"total_amount"`
	ShippingAddress string             `bson:"shipping_address" json:"shipping_address"`
	PhoneNumber     string             `bson:"phone_number" json:"phone_number"`
	PaymentMethod   string             `bson:"payment_method" json:"payment_method"` // "COD" hoặc "MOMO"
	PaymentStatus   string             `bson:"payment_status" json:"payment_status"` // "Unpaid" (Chưa trả), "Paid" (Đã trả)
	Status          string             `bson:"status" json:"status"`                 // "Pending", "Shipping", "Completed", "Cancelled"
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
}
