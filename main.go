package main

import (
	config "library/configs"
	_ "library/docs"
	"library/internal/database"
	"library/internal/handlers"
	"library/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Library API
// @version 1.0
// @description This is a sample library server
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
func main() {

	cfg := config.LoadConfig()

	if err := database.ConnectDatabase(); err != nil {
		panic(err)
	}

	if err := database.Migrate(); err != nil {
		panic(err)
	}

	router := gin.Default()

	router.Static("/docs", "./docs")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	router.GET("/", handlers.Welcome)
	router.GET("/getBooks", handlers.GetBooks)
	router.POST("/addBooks", middleware.RoleMiddleware("admin"), handlers.AddBook)
	router.POST("/deleteBook", middleware.RoleMiddleware("admin"), handlers.DeleteBook)
	router.GET("/getBook", middleware.JWTMiddleware(), handlers.GetBook)
	router.POST("/modifyingBook", middleware.RoleMiddleware("admin"), handlers.ModifyingBook)
	router.POST("/register", handlers.RegisterUser)
	router.POST("/login", handlers.LoginUser)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		panic(err)
	}
}
