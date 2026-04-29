package repositories

import (
	"auto-zen-backend/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReadinessRepository interface {
	Save(r *models.ReadinessRecord) error
	FindByDate(date time.Time) (*models.ReadinessRecord, error)
	FindRecent(days int) ([]models.ReadinessRecord, error)
}

type readinessRepository struct {
	db *gorm.DB
}

func NewReadinessRepository(db *gorm.DB) ReadinessRepository {
	return &readinessRepository{db: db}
}

// Save は UPSERT を使い、同日データが既にあれば上書きする（べき等性確保）
func (r *readinessRepository) Save(record *models.ReadinessRecord) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"score", "hrv_balance", "raw_json", "fetched_at"}),
	}).Create(record).Error
}

func (r *readinessRepository) FindByDate(date time.Time) (*models.ReadinessRecord, error) {
	var record models.ReadinessRecord
	err := r.db.Where("date = ?", date.Format("2006-01-02")).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindRecent は直近 days 日分のレコードを新しい順で返す（HRVスコア正規化に使用）
func (r *readinessRepository) FindRecent(days int) ([]models.ReadinessRecord, error) {
	var records []models.ReadinessRecord
	err := r.db.Order("date DESC").Limit(days).Find(&records).Error
	return records, err
}
