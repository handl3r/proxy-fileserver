package api

import (
	"proxy-fileserver/adapter"
	"proxy-fileserver/common/lock"
	"proxy-fileserver/common/log"
	"proxy-fileserver/repository"
)

type Cleaner struct {
	FileInfoRepo    *repository.FileInfoRepository
	LocalFileSystem *adapter.LocalFileSystem
	SharedFolder    string
	ExpiredTime     int // minute
}

func NewCleaner(fileInfoRepo *repository.FileInfoRepository, expiredTime int, sharedFolder string, localFileSystem *adapter.LocalFileSystem) *Cleaner {
	return &Cleaner{
		FileInfoRepo:    fileInfoRepo,
		LocalFileSystem: localFileSystem,
		SharedFolder:    sharedFolder,
		ExpiredTime:     expiredTime,
	}
}

func (c *Cleaner) Run() {
	fileInfos, err := c.FileInfoRepo.GetRecordOutDate(c.ExpiredTime)
	if err != nil {
		log.Errorf("[Cleaner]Error when get record outDate, error: %v", err)
		return
	}
	if len(fileInfos) != 0 {
		for _, fileInfo := range fileInfos {
			err := lock.WLockWithKey(fileInfo.FilePath)
			if err != nil {
				log.Errorf("[Cleaner]Can not WLOCK for filepath %s with error: %v", fileInfo.FilePath, err)
				continue
			}
			err = c.LocalFileSystem.Delete(c.SharedFolder + "/" + fileInfo.FilePath)
			if err != nil {
				log.Errorf("[Cleaner]Can not delete file %s at local file system with error: %v", fileInfo.FilePath, err)
			} else {
				log.Infof("[Cleaner]Deleted file %s at local file system", fileInfo.FilePath)
			}
			err = c.FileInfoRepo.Delete(fileInfo.FilePath)
			if err != nil {
				log.Errorf("[Cleaner]Can not remove file %s, last_download_at %v in database with error: %v",
					fileInfo.FilePath, fileInfo.LastDownloadAt, err)
			} else {
				log.Infof("[Cleaner]Remove file %s, last_download_at %v in database",
					fileInfo.FilePath, fileInfo.LastDownloadAt)
			}
			err = lock.WUnLockWithKey(fileInfo.FilePath)
			if err != nil {
				log.Errorf("[Cleaner]Can not WUNLOCK for filepath %s with error: %v", fileInfo.FilePath, err)
			}
		}
	}
	log.Infof("[Cleaner]Finish clean %d file", len(fileInfos))
}
