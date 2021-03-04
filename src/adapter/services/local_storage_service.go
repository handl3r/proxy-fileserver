package services

import (
	"io"
	"os"
	"time"
)

type LocalStorageService struct {
	rootFolder string
	cacheTime  time.Duration
}

func NewLocalStorageService(rootFolder string, cacheTime time.Duration) *LocalStorageService {
	return &LocalStorageService{
		rootFolder: rootFolder,
		cacheTime:  cacheTime,
	}
}

func (s *LocalStorageService) IsExisted(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *LocalStorageService) GetStreamSourceByFilePath(filePath string) (io.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *LocalStorageService) NewFile(filePath string, dataSource io.Reader) error {
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, dataSource)
	if err != nil {
		return err
	}
	return nil
}
