package models

import (
	"time"
)

type FileInfo struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FilePath       string    `gorm:"primaryKey;unique;column:file_path"`
	LastDownloadAt time.Time `gorm:"column:last_download_at"`
}
