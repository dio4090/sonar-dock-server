package routes

import (
	"github.com/dio4090/sonar-dock-server/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		controllers.RegisterAuthRoutes(authGroup)
	}
}
