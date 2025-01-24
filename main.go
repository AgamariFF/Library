package main

import (
	config "library/configs"
	_ "library/docs"
	"library/internal/database"
	"library/internal/handlers"
	"time"

	// "library/internal/kafka"
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

	if err := database.ConnectWithRetry(6, time.Second); err != nil {
		panic(err)
	}

	if err := database.Migrate(); err != nil {
		panic(err)
	}

	// producer, err := kafka.NewKafkaProducer([]string{"localhost:9092"}, "library-events")
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() {
	// 	if producer != nil {
	// 		producer.Close()
	// 	}
	// }()

	// consumer, err := kafka.NewKafkaConsumer([]string{"localhost:9092"}, "library-events")
	// if err != nil {
	// 	panic(err)
	// }
	// go consumer.ConsumeMessage()

	router := gin.Default()

	router.Static("/docs", "./docs")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	router.GET("/", handlers.Welcome)
	router.GET("/getBooks", handlers.GetBooks(database.DB))
	router.POST("/addBook", middleware.RoleMiddleware("admin"), handlers.AddBook(database.DB, nil))
	router.POST("/deleteBook", middleware.RoleMiddleware("admin"), handlers.DeleteBook(database.DB))
	router.GET("/getBook", middleware.JWTMiddleware(), handlers.GetBook(database.DB))
	router.DELETE("/modifyingBook", middleware.RoleMiddleware("admin"), handlers.ModifyingBook(database.DB))
	router.POST("/register", handlers.RegisterUser(database.DB))
	router.POST("/login", handlers.LoginUser(database.DB))

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		panic(err)
	}
}
