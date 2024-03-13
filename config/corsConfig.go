package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Middleware para habilitar CORS
func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"POST", "OPTIONS", "GET", "PUT"}
	config.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
		"Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"}
	return cors.New(config)
}
