package services

import (
	"io"
	"os"
	"time"
)

type LocalFileSystem struct {
	rootFolder string
	cacheTime  time.Duration
}

func NewLocalFileSystem(rootFolder string, cacheTime time.Duration) *LocalFileSystem {
	return &LocalFileSystem{
		rootFolder: rootFolder,
		cacheTime:  cacheTime,
	}
}

func (s *LocalFileSystem) IsExisted(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *LocalFileSystem) GetStreamSourceByFilePath(filePath string) (io.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *LocalFileSystem) NewFileAndStream(filePath string, dataSource io.Reader) error {
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

func (s *LocalFileSystem) NewFile(filePath string) (io.Writer, error) {
	w, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return w, nil
}
