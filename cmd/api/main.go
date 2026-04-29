package main

import (
	"context"
	"fmt"
	"os"

	httpctrl "auto-zen-backend/controllers/http"
	"auto-zen-backend/infra"
	"auto-zen-backend/infra/oura"
	"auto-zen-backend/middlewares"
	"auto-zen-backend/repositories"
	"auto-zen-backend/services"
	"auto-zen-backend/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	// --- インフラ層 ---
	db := infra.InitDB()
	ouraClient := oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	// --- リポジトリ層 ---
	logRepo := repositories.NewLogRepository(db)
	userRepo := repositories.NewUserRepository(db)
	readinessRepo := repositories.NewReadinessRepository(db)
	sleepRepo := repositories.NewSleepRepository(db)
	summaryRepo := repositories.NewSummaryRepository(db)

	// --- サービス層 ---
	logService := services.NewLogService(logRepo)
	userService := services.NewUserService(userRepo)
	analyzerService := services.NewAnalyzerService(ouraClient, readinessRepo, sleepRepo, summaryRepo)

	// --- コントローラー層 ---
	logCtrl := httpctrl.NewLogController(logService)
	userCtrl := httpctrl.NewUserController(userService)
	webhookCtrl := httpctrl.NewWebhookController(analyzerService, os.Getenv("OURA_WEBHOOK_VERIFY_TOKEN"))

	// --- Gin ルーター ---
	r := gin.Default()

	authorized := r.Group("/")
	authorized.Use(middlewares.AuthMiddleware())
	{
		authorized.GET("/logs", logCtrl.GetLogs)
		authorized.POST("/save", logCtrl.SaveLog)
		authorized.DELETE("/delete", logCtrl.DeleteLog)
	}

	r.POST("/signup", userCtrl.Signup)
	r.POST("/login", userCtrl.Login)
	r.POST("/webhook/oura", webhookCtrl.HandleOuraEvent)

	// --- バックグラウンドワーカー ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	poller := worker.NewPoller(ouraClient, analyzerService)
	go poller.Start(ctx)

	fmt.Println("Auto-Zen Backend is starting on :8081...")
	r.Run(":8081")
}
