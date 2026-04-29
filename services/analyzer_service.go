package services

import (
	"auto-zen-backend/infra/oura"
	"auto-zen-backend/models"
	"auto-zen-backend/repositories"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type AnalyzerService interface {
	RunDailyAnalysis(ctx context.Context, date time.Time) error
}

type analyzerService struct {
	ouraClient    *oura.Client
	readinessRepo repositories.ReadinessRepository
	sleepRepo     repositories.SleepRepository
	summaryRepo   repositories.SummaryRepository
}

func NewAnalyzerService(
	c *oura.Client,
	rr repositories.ReadinessRepository,
	sr repositories.SleepRepository,
	smr repositories.SummaryRepository,
) AnalyzerService {
	return &analyzerService{
		ouraClient:    c,
		readinessRepo: rr,
		sleepRepo:     sr,
		summaryRepo:   smr,
	}
}

// RunDailyAnalysis は Oura API からその日のデータを取得して DB に保存する。
// Phase 1: 保存のみ。スコア計算・通知は Phase 2/3 で追加。
func (s *analyzerService) RunDailyAnalysis(ctx context.Context, date time.Time) error {
	dateStr := date.Format("2006-01-02")

	if err := s.fetchAndSaveReadiness(ctx, date, dateStr); err != nil {
		return fmt.Errorf("analyzer: %w", err)
	}
	if err := s.fetchAndSaveSleep(ctx, date, dateStr); err != nil {
		return fmt.Errorf("analyzer: %w", err)
	}
	return nil
}

func (s *analyzerService) fetchAndSaveReadiness(ctx context.Context, date time.Time, dateStr string) error {
	data, err := s.ouraClient.GetDailyReadiness(ctx, dateStr)
	if err != nil {
		return fmt.Errorf("fetch readiness: %w", err)
	}

	raw, _ := json.Marshal(data)
	record := &models.ReadinessRecord{
		Date:    date,
		RawJSON: string(raw),
	}
	if data.Score != nil {
		record.Score = *data.Score
	}
	if data.Contributors.HRVBalance != nil {
		record.HRVBalance = *data.Contributors.HRVBalance
	}

	if err := s.readinessRepo.Save(record); err != nil {
		return fmt.Errorf("save readiness: %w", err)
	}
	return nil
}

func (s *analyzerService) fetchAndSaveSleep(ctx context.Context, date time.Time, dateStr string) error {
	data, err := s.ouraClient.GetDailySleep(ctx, dateStr)
	if err != nil {
		return fmt.Errorf("fetch sleep: %w", err)
	}

	raw, _ := json.Marshal(data)
	record := &models.SleepRecord{
		Date:    date,
		RawJSON: string(raw),
	}
	if data.Score != nil {
		record.Score = *data.Score
	}
	// efficiency は contributor スコア（0-100）を暫定値として使用
	// Phase 2 で /v2/usercollection/sleep の actual_efficiency に差し替え予定
	if data.Contributors.Efficiency != nil {
		record.Efficiency = *data.Contributors.Efficiency
	}

	if err := s.sleepRepo.Save(record); err != nil {
		return fmt.Errorf("save sleep: %w", err)
	}
	return nil
}
