package services

import (
	"io"
	"proxy-fileserver/common/log"
	"proxy-fileserver/enums"
	"proxy-fileserver/src/adapter/services"
	"strings"
)

type FileSystemService struct {
	rootFolder string
	GoogleDrive  services.GoogleDriveService
	LocalStorage services.LocalStorageService
}

func NewLocalStorageService(rootFolder string, googleDrive services.GoogleDriveService, localStorage services.LocalStorageService) *FileSystemService {
	return &FileSystemService{
		GoogleDrive:  googleDrive,
		LocalStorage: localStorage,
	}
}

func (s *FileSystemService) GetStreamSourceFile(filePath string) (io.Reader, error) {
	if !s.validateFilePath(filePath) {
		return nil, enums.ErrFileNotExist
	}
	existed, err := s.LocalStorage.IsExisted(filePath)
	if err != nil {
		log.Errorf("Can not check if file exist on local: [%v]", err)
		return nil, err
	}
	if existed {
		streamSource, err := s.LocalStorage.GetStreamSourceByFilePath(filePath)
		if err != nil {
			log.Errorf("Can not getStreamSourceByFilePath [%s], with error: [%v]", filePath, err)

			return nil, err
		}
		return streamSource, nil
	}



}

// FilePath format: x/y/z/t.a
func (s *FileSystemService) validateFilePath(filePath string) bool {
	files := strings.Split(filePath, "/")
	if len(files) < 2 {
		return false
	}
	for _, file := range files {
		if len(strings.TrimSpace(file)) == 0 {
			return false
		}
	}
	return true
}
