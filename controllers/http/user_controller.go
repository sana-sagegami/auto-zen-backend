package http

import (
	"auto-zen-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{Service: s}
}

// POST / signup
func (ctrl *UserController) Signup(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザー名とパスワードは必須です"})
		return
	}

	err := ctrl.Service.Signup(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "ユーザーの作成に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
}

func (ctrl *UserController) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力が正しくありません"})
		return
	}

	token, err := ctrl.Service.Login(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザー名またはパスワードが違います"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
