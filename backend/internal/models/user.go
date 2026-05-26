package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User định nghĩa cấu trúc dữ liệu của một người dùng trong hệ thống
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`            // Tự động giấu mật khẩu khi trả về JSON
	Role      string             `bson:"role" json:"role"`             // Quyền hạn: "user" hoặc "admin"
	IsBlocked bool               `bson:"is_blocked" json:"is_blocked"` // Trạng thái tài khoản: true (Bị khóa), false (Bình thường)
	CreatedAt time.Time          `bson:"created_at" json:"created_at"` // Thời gian tạo tài khoản
}
