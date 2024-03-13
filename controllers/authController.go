package controllers

import (
	"net/http"

	"github.com/dio4090/sonar-dock-server/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	group.GET("/verify", verifyHandler)
}

// loginHandler godoc
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
func loginHandler(c *gin.Context) {
	var loginParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// TODO: CORRIGIR ERRO DE IMPORT
	user, err := service.GetUserByEmail(loginParams.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusForbidden, gin.H{"forbidden": "invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginParams.Password))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"forbidden": "invalid email or password"})
		return
	}

	token, err := service.GenerateJWT(loginParams.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	service.RedisSaveSessionToken(loginParams.Email, token, c)

	c.Status(http.StatusOK)
}

// verifyHandler godoc
// @Summary verifyHandler
// @Description Verifica o token da sessão
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} LoginResponse "Status do token da sessão"
// @Failure 400 {object} model.ErrorResponse "Requisição inválida"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 403 {object} model.ForbiddenResponse "Token de sessão inválido"
// @Router /auth/verify [get]
func verifyHandler(c *gin.Context) {
	if service.AuthResource(c) {
		c.JSON(http.StatusOK, gin.H{"isAuthenticated": true})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"isAuthenticated": false})
	}
}
