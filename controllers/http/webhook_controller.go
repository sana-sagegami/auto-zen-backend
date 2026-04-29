package http

import (
	dto "auto-zen-backend/dto/http"
	"auto-zen-backend/services"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WebhookController struct {
	analyzerService services.AnalyzerService
	verifyToken     string
}

func NewWebhookController(s services.AnalyzerService, verifyToken string) *WebhookController {
	return &WebhookController{analyzerService: s, verifyToken: verifyToken}
}

// HandleOuraEvent は Oura Ring からの Webhook イベントを受信するハンドラ。
// Oura は同一イベントを重複送信することがあるが、
// リポジトリ層の UPSERT (ON CONFLICT date) によりべき等性を担保する。
func (c *WebhookController) HandleOuraEvent(ctx *gin.Context) {
	// トークン検証
	if ctx.GetHeader("x-oura-verification-token") != c.verifyToken {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var payload dto.OuraWebhookPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	date, err := time.Parse("2006-01-02", payload.Day)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid day format"})
		return
	}

	// daily_readiness または daily_sleep 確定時に分析を非同期実行
	if payload.DataType == "daily_readiness" || payload.DataType == "daily_sleep" {
		go func() {
			if err := c.analyzerService.RunDailyAnalysis(context.Background(), date); err != nil {
				log.Printf("[webhook] RunDailyAnalysis failed (date=%s): %v", payload.Day, err)
			}
		}()
	}

	ctx.Status(http.StatusNoContent)
}
