package main

import (
	"auto-zen-backend/controllers/http"
	"auto-zen-backend/infra"
	"auto-zen-backend/middlewares"
	"auto-zen-backend/repositories"
	"auto-zen-backend/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// インフラ層
	db := infra.InitDB()

	// リポジトリ層
	logRepo := repositories.NewLogRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// サービス層
	logService := services.NewLogService(logRepo)
	userService := services.NewUserService(userRepo)

	// コントローラー層
	logCtrl := http.NewLogController(logService)
	userCtrl := http.NewUserController(userService)

	// Ginルーター
	r := gin.Default()

authorized := r.Group("/")
authorized.Use(middlewares.AuthMiddleware())
{
		// ルーティング
	r.GET("/logs", logCtrl.GetLogs)
	r.POST("/save", logCtrl.SaveLog)
	r.DELETE("/delete", logCtrl.DeleteLog)
}
	// ユーザー
	r.POST("/signup", userCtrl.Signup)
	r.POST("/login", userCtrl.Login)

	fmt.Printf("/signup handler type: %T\n", userCtrl.Signup)

	// サーバー起動
	fmt.Println("Auto-Zen Backend is starting on :8081...")
	r.Run(":8081")
}
