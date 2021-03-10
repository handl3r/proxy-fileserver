package repository

import (
	"gorm.io/gorm"
	"proxy-fileserver/models"
	"time"
)

type FileInfoRepository struct {
	orm *gorm.DB
}

func NewFileInfoRepository(db *gorm.DB) *FileInfoRepository {
	return &FileInfoRepository{
		orm: db,
	}
}

func (r *FileInfoRepository) Create(model models.FileInfo) error {
	return r.orm.Save(&model).Error
}

func (r *FileInfoRepository) Update(model models.FileInfo) error {
	return r.orm.Model(&models.FileInfo{}).Where("file_path = ?", model.FilePath).Update("last_download_at", time.Now()).Error
}

func (r *FileInfoRepository) Delete(filePath string) error {
	return r.orm.Delete(&models.FileInfo{FilePath: filePath}).Error
}

func (r *FileInfoRepository) GetRecordOutDate(minute int) ([]models.FileInfo, error) {
	var fileInfos []models.FileInfo
	now := time.Now()
	lastTime := now.Add(time.Duration(-minute) * time.Minute)
	err := r.orm.Where("last_download_at <= ?", lastTime).Find(&fileInfos).Error
	if err != nil {
		return nil, err
	}
	return fileInfos, nil
}
