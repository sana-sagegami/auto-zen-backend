package repositories

import (
	"auto-zen-backend/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SleepRepository interface {
	Save(s *models.SleepRecord) error
	FindByDate(date time.Time) (*models.SleepRecord, error)
	FindRecent(days int) ([]models.SleepRecord, error)
}

type sleepRepository struct {
	db *gorm.DB
}

func NewSleepRepository(db *gorm.DB) SleepRepository {
	return &sleepRepository{db: db}
}

// Save は UPSERT を使い、同日データが既にあれば上書きする（べき等性確保）
func (r *sleepRepository) Save(record *models.SleepRecord) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"score", "total_minutes", "efficiency", "wake_time", "raw_json"}),
	}).Create(record).Error
}

func (r *sleepRepository) FindByDate(date time.Time) (*models.SleepRecord, error) {
	var record models.SleepRecord
	err := r.db.Where("date = ?", date.Format("2006-01-02")).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindRecent は直近 days 日分のレコードを新しい順で返す（睡眠負債計算に使用）
func (r *sleepRepository) FindRecent(days int) ([]models.SleepRecord, error) {
	var records []models.SleepRecord
	err := r.db.Order("date DESC").Limit(days).Find(&records).Error
	return records, err
}
