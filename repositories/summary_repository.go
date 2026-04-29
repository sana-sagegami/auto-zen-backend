package repositories

import (
	"auto-zen-backend/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SummaryRepository interface {
	Save(s *models.DailySummary) error
	FindByDate(date time.Time) (*models.DailySummary, error)
}

type summaryRepository struct {
	db *gorm.DB
}

func NewSummaryRepository(db *gorm.DB) SummaryRepository {
	return &summaryRepository{db: db}
}

// Save は UPSERT を使い、同日サマリーが既にあれば全フィールドを上書きする
func (r *summaryRepository) Save(summary *models.DailySummary) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"condition_score", "focus_peak_start", "focus_peak_end",
			"recommend_bedtime", "sleep_debt_minutes",
		}),
	}).Create(summary).Error
}

func (r *summaryRepository) FindByDate(date time.Time) (*models.DailySummary, error) {
	var summary models.DailySummary
	err := r.db.Where("date = ?", date.Format("2006-01-02")).First(&summary).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}
