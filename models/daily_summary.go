package models

import "time"

type DailySummary struct {
	Date             time.Time  `gorm:"type:date;primaryKey"`
	ConditionScore   int
	FocusPeakStart   *time.Time
	FocusPeakEnd     *time.Time
	RecommendBedtime *time.Time
	SleepDebtMinutes int
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}

func (DailySummary) TableName() string { return "daily_summaries" }
