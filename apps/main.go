package main

import (
	"apigin/apps/config"
	"apigin/apps/controller"
	"apigin/apps/pkg/token"
	"apigin/apps/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.Use(CORS())

	authController := controller.AuthController{
		Db: db,
	}

	v1 := router.Group("v1")
	router.GET("/ping", Ping)

	auth := v1.Group("auth")
	{
		auth.POST("register", authController.Register)
		auth.POST("login", authController.Login)
		auth.GET("profile", CheckAuth(), authController.Profile)
	}

	router.Run(":4000")
}

func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "w")
		ctx.Header("Access-Control-Request-Method", "GET, OPTION, POST, PUT, DELETE")
		ctx.Header("Access-Control-Request-Headers", "Authorization, Content-Type")
		ctx.Next()
	}
}

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Okee",
	})
}

func CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		bearerToken := strings.Split(header, "Bearer ")
		if len(bearerToken) != 2 {
			res := response.ResponseAPI{
				StatusCode: http.StatusUnauthorized,
				Message:    "unauthorized2",
			}
			ctx.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		payload, err := token.ValidateToken(bearerToken[1])
		if err != nil {
			res := response.ResponseAPI{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid Token",
				Payload:    err.Error(),
			}
			ctx.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		ctx.Set("authId", payload.AuthId)

		ctx.Next()
	}
}
