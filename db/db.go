// Migration

package db

import (
	"os"

	"github.com/EikoNakashima/CustomerCalendar-go/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var (
	db  *gorm.DB
	err error
)

// DB初期化
func Init() {
	// 実行環境取得
	env := os.Getenv("ENV")

	if "production" == env {
		env = "production"
	} else {
		env = "development"
	}

	// 環境変数取得
	godotenv.Load(".env." + env)
	godotenv.Load()

	// DB接続
	db, err = gorm.Open("mysql", os.Getenv("CONNECT"))

	if err != nil {
		panic(err)
	}

	autoMigration()
}

// DB取得
func GetDB() *gorm.DB {
	return db
}

// DB接続終了
func Close() {
	if err := db.Close(); err != nil {
		panic(err)
	}
}

// マイグレーション
func autoMigration() {
	db.AutoMigrate(&models.User{})
}
