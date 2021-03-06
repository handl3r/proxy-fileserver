package repository

import (
	"database/sql"
	"proxy-fileserver/configs"
)

type ProviderRepository interface {
	GetFileInfoRepository() *FileInfoRepository
}

type providerRepositoryImpl struct {
	fileInfoRepo *FileInfoRepository
}

func NewProviderRepository(db *sql.DB, conf *configs.Config) ProviderRepository {
	return &providerRepositoryImpl{
		fileInfoRepo: NewFileInfoRepository(
			db, conf.MysqlFileInfoTable,
		),
	}
}

func (p *providerRepositoryImpl) GetFileInfoRepository() *FileInfoRepository {
	return p.fileInfoRepo
}
