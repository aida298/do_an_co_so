package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() *mongo.Database {
	// URI mặc định của MongoDB chạy qua Docker
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Lỗi kết nối MongoDB: ", err)
	}

	// Ping thử xem database có phản hồi không
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Không thể ping tới MongoDB: ", err)
	}

	fmt.Println("Đã kết nối thành công tới MongoDB!")
	Client = client
	
	// Khởi tạo và trả về database có tên là "gearpc_db"
	return client.Database("gearpc_db")
}