package controller

import (
	"github.com/gin-gonic/gin"
	"go-web-demo/config"
	"go-web-demo/middleware"
	"go-web-demo/models"
	"net/http"
)

func Login(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//  查询数据库记录
	var dbUser models.User
	result := config.DB.Where("username = ?", user.Username).First(&dbUser)
	if result.Error == nil {
		// 验证密码
		if dbUser.Password == user.Password {
			// 创建token
			token, _ := middleware.CreateToken(dbUser)

			c.Header("Authorization", "Bearer "+token)
			c.JSON(http.StatusOK, gin.H{"token": token})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
	}
}

func Register(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 检查是否存在同username用户
	var dbUser models.User
	result := config.DB.Where("username = ?", user.Username).First(&dbUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}
	config.DB.Create(&user)
	// 用jwt创建登录态并设置到header里
	token, _ := middleware.CreateToken(user)

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusCreated, user)

}

func QueryUser(c *gin.Context) {
	var user models.User

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	// 使用用户 ID 进行后续操作
	// 例如，从数据库中查询用户信息
	config.DB.First(&user, userId)

	c.JSON(http.StatusOK, user)
}
