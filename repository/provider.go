package repository

import (
	"gorm.io/gorm"
)

type ProviderRepository interface {
	GetFileInfoRepository() *FileInfoRepository
}

type providerRepositoryImpl struct {
	fileInfoRepo *FileInfoRepository
}

func NewProviderRepository(db *gorm.DB) ProviderRepository {
	return &providerRepositoryImpl{
		fileInfoRepo: NewFileInfoRepository(db),
	}
}

func (p *providerRepositoryImpl) GetFileInfoRepository() *FileInfoRepository {
	return p.fileInfoRepo
}
