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
			dbInsert(content)
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
		tweet := ctx.PostForm("tweet")
		dbUpdate(id, tweet)
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
	defer db.Close()
	db.AutoMigrate(&Tweet{}) //構造体に基づいてテーブルを作成
}

// データインサート処理
func dbInsert(content string) {
	db := gormConnect()

	defer db.Close()
	// Insert処理
	db.Create(&Tweet{Content: content})
}

//DB更新
func dbUpdate(id int, tweetText string) {
	db := gormConnect()
	var tweet Tweet
	db.First(&tweet, id)
	tweet.Content = tweetText
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

// type Todo struct {
// 	gorm.Model
// 	Text   string
// 	Status string
// }

// //DB初期化
// func dbInit() {
// 	db, err := gorm.Open("mysql", "root@/sample?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic("データベース開けず！（dbInit）")
// 	}
// 	db.AutoMigrate(&Todo{})
// 	defer db.Close()
// 	db.LogMode(true)
// }

// //DB追加
// func dbInsert(text string, status string) {
// 	db, err := gorm.Open("mysql", "root@/sample?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic("データベース開けず！（dbInsert)")
// 	}
// 	db.Create(&Todo{Text: text, Status: status})
// 	defer db.Close()
// }

// //DB更新
// func dbUpdate(id int, text string, status string) {
// 	db, err := gorm.Open("mysql", "root@/sample?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic("データベース開けず！（dbUpdate)")
// 	}
// 	var todo Todo
// 	db.First(&todo, id)
// 	todo.Text = text
// 	todo.Status = status
// 	db.Save(&todo)
// 	db.Close()
// }

// //DB削除
// func dbDelete(id int) {
// 	db, err := gorm.Open("mysql", "root@/sample?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic("データベース開けず！（dbDelete)")
// 	}
// 	var todo Todo
// 	db.First(&todo, id)
// 	db.Delete(&todo)
// 	db.Close()
// }

// //DB全取得
// func dbGetAll() []Todo {
// 	db, err := gorm.Open("mysql", "root@/sample?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic("データベース開けず！(dbGetAll())")
// 	}
// 	var todos []Todo
// 	db.Order("created_at desc").Find(&todos)
// 	db.Close()
// 	return todos
// }

// //DB一つ取得
// func dbGetOne(id int) Todo {
// 	db, err := gorm.Open("mysql", "root@/sample?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic("データベース開けず！(dbGetOne())")
// 	}
// 	var todo Todo
// 	db.First(&todo, id)
// 	db.Close()
// 	return todo
// }
