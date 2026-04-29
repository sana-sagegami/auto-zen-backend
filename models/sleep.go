package models

import "time"

type SleepRecord struct {
	ID           string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Date         time.Time  `gorm:"type:date;uniqueIndex;not null"`
	Score        int
	TotalMinutes int
	Efficiency   int
	WakeTime     *time.Time
	RawJSON      string `gorm:"type:jsonb"`
}

func (SleepRecord) TableName() string { return "sleep_records" }
