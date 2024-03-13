package service

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dio4090/sonar-dock-server/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

// JWT SECRET KEY
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var tokens []string

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func redisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0, // Use default DB
	})
	return redisClient
}

func AuthResource(c *gin.Context) bool {
	token, err := c.Cookie("token")
	if err != nil {
		bearerToken := c.GetHeader("Authorization")
		if bearerToken != "" {
			token = strings.TrimPrefix(bearerToken, "Bearer ")
		} else {
			return false
		}
	}

	email, err := extractEmailFromToken(token)
	if err != nil {
		// Se não conseguir extrair o email, retorna falso
		return false
	}

	key := "token_" + email
	ctx := context.Background()

	// Validação com o REDIS
	exists, err := redisClient().Exists(ctx, key).Result()
	if err != nil || exists == 0 {
		return false
	}

	return true
}

func RedisSaveSessionToken(email string, token string, c *gin.Context) error {
	// Armazenar o token no Redis com uma chave baseada no email do usuário e definir um tempo de expiração
	err := redisClient().Set(c, "token_"+email, token, 3600*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao salvar o token da sessão no Redis"})
	}

	tokens = append(tokens, token)
	return nil
}

func GenerateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(config.TOKEN_EXPIRATION_TIME)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// SHA-256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func extractEmailFromToken(tokenStr string) (string, error) {
	claims := &Claims{}

	// Parse o token
	parsedToken, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method: " + token.Header["alg"].(string))
		}
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if !parsedToken.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Email, nil
}
