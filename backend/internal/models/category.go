package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Category định nghĩa danh mục sản phẩm (Hỗ trợ đa cấp)
type Category struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
	// Dùng con trỏ để nếu không truyền vào, nó sẽ là null (dành cho danh mục cao nhất)
	ParentID *primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
}
