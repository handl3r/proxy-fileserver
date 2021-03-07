package api

import (
	"proxy-fileserver/common/lock"
	"proxy-fileserver/common/log"
	"proxy-fileserver/repository"
)

type Cleaner struct {
	FileInfoRepo *repository.FileInfoRepository
	ExpiredTime  int // hour
}

func NewCleaner(fileInfoRepo *repository.FileInfoRepository, expiredTime int) *Cleaner {
	return &Cleaner{
		FileInfoRepo: fileInfoRepo,
		ExpiredTime:  expiredTime,
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
			err = c.FileInfoRepo.Delete(fileInfo.ID)
			if err != nil {
				log.Errorf("[Cleaner]Can not remove file id %d, path %s, last_download_at %v with error: %v",
					fileInfo.ID, fileInfo.FilePath, fileInfo.LastDownloadAt, err)
			} else {
				log.Infof("[Cleaner]remove file id %d, path %s, last_download_at %v", fileInfo.ID, fileInfo.FilePath, fileInfo.LastDownloadAt)
			}
			err = lock.WUnLockWithKey(fileInfo.FilePath)
			if err != nil {
				log.Errorf("[Cleaner]Can not WUNLOCK for filepath %s with error: %v", fileInfo.FilePath, err)
			}
		}
	}
	log.Infof("[Cleaner]Finish clean %d file", len(fileInfos))
}
