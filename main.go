package main

import (
	"os"

	"github.com/badaccuracyid/softeng_backend/src/database"
	"github.com/badaccuracyid/softeng_backend/src/middleware"
	"github.com/badaccuracyid/softeng_backend/src/routes"
	"github.com/badaccuracyid/softeng_backend/src/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/badaccuracyid/softeng_backend/docs"
)

// @title Softeng Backend API
// @version 1.0
// @description This is a sample server for a software engineering project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api/v1

func main() {
	utils.LoadEnv()

	var err = database.MigrateTables()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	router.Use(cors.New(config))
	router.Use(middleware.UserMiddleware())

	router = routes.InitializeUserRoutes(router)
	router = routes.InitializeChatRoutes(router)

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
