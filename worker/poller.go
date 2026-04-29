package worker

import (
	"auto-zen-backend/infra/oura"
	"auto-zen-backend/services"
	"context"
	"log"
	"time"
)

type Poller struct {
	ouraClient      *oura.Client
	analyzerService services.AnalyzerService
}

func NewPoller(c *oura.Client, a services.AnalyzerService) *Poller {
	return &Poller{ouraClient: c, analyzerService: a}
}

// Start は毎朝 6:00 に HRV データをポーリングするスケジューラー。
// Phase 2 で poll() の中身を実装する。
func (p *Poller) Start(ctx context.Context) {
	for {
		next := nextOccurrence(6, 0)
		select {
		case <-time.After(time.Until(next)):
			p.poll(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// nextOccurrence は今日または翌日の hour:min の時刻を返す
func nextOccurrence(hour, min int) time.Time {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}

// poll は Phase 2 で実装: interbeat_interval を取得して RMSSD を計算する。
func (p *Poller) poll(ctx context.Context) {
	date := time.Now().Format("2006-01-02")
	if err := p.analyzerService.RunDailyAnalysis(ctx, time.Now()); err != nil {
		log.Printf("[poller] RunDailyAnalysis failed (date=%s): %v", date, err)
	}
}
