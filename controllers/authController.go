package controllers

import (
	"github.com/gin-gonic/gin"
)

// LoginRequest representa o corpo da solicitação de login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse representa a resposta do token JWT após um login bem-sucedido.
type LoginResponse struct {
	Token string `json:"token"`
}

func RegisterAuthRoutes(group *gin.RouterGroup) {
	group.POST("/login", loginHandler)
}

// AuthLogin godoc
// @Summary Login
// @Description Autentica um usuário e retorna um token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Credenciais de Login"
// @Success 200 {object} LoginResponse "Token JWT para autenticação nas requisições subsequentes"
// @Failure 400 {object} model.ErrorResponse "Requisição inválida"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 403 {object} model.ForbiddenResponse "Usuário ou senha inválidos"
// @Router /auth/login [post]
func loginHandler(c *gin.Context) {}
