package main

import (
	config "library/configs"
	"library/internal/database"
	"library/internal/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Library API
// @version 1.0
// @description This is a sample library server
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/docs/swagger.json")))

	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	router.GET("/", handlers.Welcome)
	router.GET("/getBooks", handlers.GetBooks)
	router.POST("/addBooks", handlers.AddBook)
	router.POST("/deleteBook", handlers.DeleteBook)
	router.GET("/getBook", handlers.GetBook)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		panic(err)
	}
}
