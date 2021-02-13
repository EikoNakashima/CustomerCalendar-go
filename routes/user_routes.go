package routes

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func UserSignUp(ctx *gin.Context) {
    println("post/signup")
    username := ctx.PostForm("username")
    email := ctx.PostForm("email")
    password := ctx.PostForm("password")
    passwordConf := ctx.PostForm("passwordconfirmation")
    println("username: " + username)
    println("email: " + email)
    println("password: " + password)
    println("passwordConf: " + passwordConf)

    ctx.Redirect(http.StatusSeeOther, "//localhost:8080/")
}

func UserLogIn(ctx *gin.Context) {
    println("post/login")
    username := ctx.PostForm("username")
    password := ctx.PostForm("password")
    println("username: " + username)
    println("password: " + password)

    ctx.Redirect(http.StatusSeeOther, "//localhost:8080/")
}
