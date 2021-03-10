package models

import (
	"time"
)

type FileInfo struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FilePath       string    `gorm:"primaryKey;column:file_path"`
	LastDownloadAt time.Time `gorm:"column:last_download_at"`
}
