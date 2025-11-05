package main

import (
	"fmt"
	"log/slog"
	"github.com/BahadirAhmedov/data-aggregation/internal/app"
	"github.com/BahadirAhmedov/data-aggregation/internal/config"
	"github.com/gin-gonic/gin"
	"os"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files" 
	_ "github.com/BahadirAhmedov/data-aggregation/cmd/data-aggregation/docs"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

const (
	envLocal = "local"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg.Storage)

	logger := setupLogger(envLocal)

	handlers := app.New(logger, cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User,
		cfg.Storage.Password, cfg.Storage.DbName)


	router := gin.Default()
	
	// Create
	router.POST("/subscriptions", handlers.CreateSubscription(logger))
	// Read
 	router.GET("/subscriptions/:id", handlers.ReadSubscription(logger))
	// Update
	router.PUT("/subscriptions/:id", handlers.UpdateSubscription(logger))
	// Delete
	router.DELETE("/subscriptions/:id", handlers.DeleteSubscription(logger))
	// List
	router.GET("/subscriptions", handlers.ListSubscription(logger))


	// Sum
	router.POST("/subscriptions/sum", handlers.SumSubscriptions(logger))

	
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}