package jobs

import (
	"proxy-fileserver/common/lock"
	"proxy-fileserver/common/log"
	"proxy-fileserver/repository"
)

func LoadLockFromDB(fileInfoRepo *repository.FileInfoRepository) {
	log.Infof("[Preload-job] start")
	fileInfos, err := fileInfoRepo.GetAll()
	if err != nil {
		log.Errorf("Can not get all file info with error: %v", err)
		return
	}
	numSuccess := 0
	numErr := 0
	for _, v := range fileInfos {
		err := lock.AddLock(v.FilePath)
		if err != nil {
			numErr += 1
			log.Errorf("[Preload-job] Can not add lock for exist fileInfo: %v", v)
			continue
		}
		numSuccess += 1
		log.Infof("[Preload-job] Add lock for exist fileInfo: %v", fileInfos)
	}
	log.Infof("[Preload-job] finish with loaded: %d lock, fail: %d lock", numSuccess, numErr)

}
