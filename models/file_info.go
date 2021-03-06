package models

import "time"

type FileInfo struct {
	ID int64
	FilePath string
	LastDownloadAt time.Time
}