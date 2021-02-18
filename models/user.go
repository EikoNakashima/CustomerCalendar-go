// ユーザーテーブル カラム定義

package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// type User struct {
// 	Id        uint
// 	Title     string `gorm:"size:128"`
// 	Category  int
// 	Author    string `gorm:"size:64"`
// 	CreatedAt time.Time
// }

// User モデルの宣言
type User struct {
	gorm.Model
	Username  string `form:"username" binding:"required" gorm:"unique;not null"`
	Password  string `form:"password" binding:"required"`
	CreatedAt time.Time
}
