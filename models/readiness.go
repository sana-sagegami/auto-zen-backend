package models

import "time"

type ReadinessRecord struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Date       time.Time `gorm:"type:date;uniqueIndex;not null"`
	Score      int
	HRVBalance int
	RawJSON    string    `gorm:"type:jsonb"`
	FetchedAt  time.Time `gorm:"autoCreateTime"`
}

func (ReadinessRecord) TableName() string { return "readiness_records" }
