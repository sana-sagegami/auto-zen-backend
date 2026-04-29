package apns

import "auto-zen-backend/models"

// Sender は APNs プッシュ通知の送信者。Phase 3 で実装。
type Sender struct{}

func NewSender() *Sender { return &Sender{} }

// Push は Phase 3 で実装する。現在は no-op。
func (s *Sender) Push(_ *models.DailySummary) error { return nil }
