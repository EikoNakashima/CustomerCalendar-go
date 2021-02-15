package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")

	dbInit()

	//Index
	router.GET("/", func(ctx *gin.Context) {
		tweets := dbGetAll()
		ctx.HTML(200, "index.html", gin.H{
			"tweets": tweets,
		})
	})

	//Create
	router.POST("/new", func(ctx *gin.Context) {
		var form Tweet
		// ここがバリデーション部分
		if err := ctx.Bind(&form); err != nil {
			tweets := dbGetAll()
			ctx.HTML(http.StatusBadRequest, "index.html", gin.H{"tweets": tweets, "err": err})
			ctx.Abort()
		} else {
			content := ctx.PostForm("content")
			status := ctx.PostForm("status")
			dbInsert(content, status)
			ctx.Redirect(302, "/")
		}
	})

	//Detail
	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		tweet := dbGetOne(id)
		ctx.HTML(200, "detail.html", gin.H{"tweet": tweet})
	})

	//Update
	router.POST("/update/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		content := ctx.PostForm("content")
		status := ctx.PostForm("status")
		dbUpdate(id, content, status)
		ctx.Redirect(302, "/")
	})

	//削除確認
	router.GET("/delete_check/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		tweet := dbGetOne(id)
		ctx.HTML(200, "delete.html", gin.H{"tweet": tweet})
	})

	//Delete
	router.POST("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		ctx.Redirect(302, "/")

	})

	router.Run()
}

type Tweet struct {
	gorm.Model
	Content string `form:"content" binding:"required"`
	Status  string
}

func gormConnect() *gorm.DB {
	DBMS := "mysql"
	USER := "test"
	PASS := "12345678"
	DBNAME := "test"
	// MySQLだと文字コードの問題で"?parseTime=true"を末尾につける必要がある
	CONNECT := USER + ":" + PASS + "@/" + DBNAME + "?parseTime=true"
	db, err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}
	return db
}

// DBの初期化
func dbInit() {
	db := gormConnect()

	// コネクション解放解放
	db.AutoMigrate(&Tweet{}) //構造体に基づいてテーブルを作成
	defer db.Close()
}

// データインサート処理
func dbInsert(content string, status string) {
	db := gormConnect()

	// Insert処理
	db.Create(&Tweet{Content: content, Status: status})
	defer db.Close()
}

//DB更新
func dbUpdate(id int, content string, status string) {
	db := gormConnect()
	var tweet Tweet
	db.First(&tweet, id)
	tweet.Content = content
	tweet.Status = status
	db.Save(&tweet)
	db.Close()
}

// 全件取得
func dbGetAll() []Tweet {
	db := gormConnect()

	defer db.Close()
	var tweets []Tweet
	// FindでDB名を指定して取得した後、orderで登録順に並び替え
	db.Order("created_at desc").Find(&tweets)
	return tweets
}

//DB一つ取得
func dbGetOne(id int) Tweet {
	db := gormConnect()
	var tweet Tweet
	db.First(&tweet, id)
	db.Close()
	return tweet
}

//DB削除
func dbDelete(id int) {
	db := gormConnect()
	var tweet Tweet
	db.First(&tweet, id)
	db.Delete(&tweet)
	db.Close()
}
