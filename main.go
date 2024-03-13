package main

import (
	"log"
	"os"

	"github.com/dio4090/sonar-dock-server/config"
	"github.com/dio4090/sonar-dock-server/docs"
	"github.com/dio4090/sonar-dock-server/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API de Projetos Sonar Dock
// @description Esta é a API de gerenciamento de projetos Sonar Dock.
// @version 1.0
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	loadEnv()
	setGinMode()
	r := setupRouter()
	runServer(r)
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Erro loading .env file: %v", err)
	}
}

func setGinMode() {
	if mode := os.Getenv("GIN_MODE"); mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(config.CORSMiddleware())

	// Configuração do Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Chamada para registrar as rotas
	routes.RegisterRoutes(r)
	return r
}

func runServer(r *gin.Engine) {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
