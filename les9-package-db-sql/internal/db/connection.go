package db

import (
	"context"
	"database/sql"
	"fmt"
	"gialong-dev04/internal/config"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := config.NewConfig().DNS()

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("unable to use data source name", err)
	}

	DB.SetMaxIdleConns(3)                   // Số kết nối nhàn rỗi tối đa
	DB.SetMaxOpenConns(3)                   // Số kết nối tối đa
	DB.SetConnMaxLifetime(30 * time.Minute) // Đóng kết nối sau 30m
	DB.SetConnMaxIdleTime(5 * time.Minute)  // Đóng kết nối nhàn rỗi sau 5m

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		DB.Close()
		return fmt.Errorf("DB Ping error: %w", err)
	}

	log.Println("Connected")

	return nil
}
