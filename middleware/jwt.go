package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-web-demo/models"
	"net/http"
	"strings"
	"time"
)

var SecretKey = []byte("my-secret-key")

func AuthenticateJWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.Request.Header.Get("Authorization")
		if tokenString == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			context.Abort()
			return
		}

		// 提取 Bearer 后面的 token
		bearerToken := strings.Split(tokenString, "Bearer ")
		if len(bearerToken) != 2 {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			context.Abort()
			return
		}
		tokenString = bearerToken[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					context.JSON(http.StatusBadRequest, gin.H{"error": "Malformed token"})
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					context.JSON(http.StatusUnauthorized, gin.H{"error": "Expired token"})
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					context.JSON(http.StatusUnauthorized, gin.H{"error": "Token not active yet"})
				} else {
					context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				}
			} else {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			context.Abort()
			return
		}
		if !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			context.Abort()
			return
		}

		// 解析 token 中的用户信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims"})
			context.Abort()
			return
		}

		// 获取用户 ID 或用户名
		userID, ok := claims["user_id"]
		if !ok {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in token"})
			context.Abort()
			return
		}
		// 将用户信息存储在 Gin 上下文中
		context.Set("userId", userID)
		context.Next()
	}
}

func CreateToken(user models.User) (string, error) {
	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"user_id":  user.ID,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(SecretKey)
	return tokenString, err
}
