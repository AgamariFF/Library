package main

import (
	config "library/configs"
	_ "library/docs"
	"library/internal/database"
	"library/internal/handlers"
	"library/logger"
	"time"

	"library/internal/kafka"
	"library/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Library API
// @version 1.0
// @description This is a sample library server
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
func main() {

	if err := logger.InitLog(); err != nil {
		panic("Failed to initialized logger: " + err.Error())
	}

	logger.InfoLog.Println("App started")

	cfg := config.LoadConfig()

	if err := database.ConnectWithRetry(6, time.Second); err != nil {
		logger.ErrorLog.Println("Failed connect to database with retry: " + err.Error())
	}
	time.Sleep(5 * time.Second)

	if err := database.Migrate(); err != nil {
		logger.ErrorLog.Println("Failed to migrate database: " + err.Error())
	}

	producer, err := kafka.NewKafkaProducer([]string{"kafka:9092"}, "library-events")
	if err != nil {
		logger.ErrorLog.Panicln("Failed to create kafke producer: " + err.Error())
	}
	defer func() {
		if producer != nil {
			producer.Close()
		}
	}()

	consumer, err := kafka.NewKafkaConsumer([]string{"kafka:9092"}, "library-events")
	if err != nil {
		logger.ErrorLog.Panicln("Failed to create kafka consumer: " + err.Error())
	}
	go consumer.ConsumeMessage()

	router := gin.Default()

	router.Static("/docs", "./docs")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	router.GET("/", handlers.Welcome)
	router.GET("/getBooks", handlers.GetBooks(database.DB))
	router.POST("/addBook", middleware.RoleMiddleware("admin"), handlers.AddBook(database.DB, producer))
	router.DELETE("/deleteBook", middleware.RoleMiddleware("admin"), handlers.DeleteBook(database.DB))
	router.GET("/getBook", middleware.JWTMiddleware(), handlers.GetBook(database.DB))
	router.POST("/modifyingBook", middleware.RoleMiddleware("admin"), handlers.ModifyingBook(database.DB))
	router.POST("/register", handlers.RegisterUser(database.DB))
	router.POST("/login", handlers.LoginUser(database.DB))
	router.POST("/unsubMailing", handlers.UnsubscribeMailing(database.DB))
	router.POST("/subMailing", handlers.SubscribeMailing(database.DB))

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		panic(err)
	}
}
