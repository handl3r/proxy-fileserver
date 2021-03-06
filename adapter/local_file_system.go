package adapter

import (
	"io"
	"os"
	"strings"
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
	_, err := os.Stat(s.rootFolder + "/" + filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *LocalFileSystem) GetStreamSourceByFilePath(filePath string) (io.Reader, error) {
	file, err := os.Open(s.rootFolder + "/" + filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *LocalFileSystem) NewFileAndStream(filePath string, dataSource io.Reader) error {
	w, err := os.Create(s.rootFolder + "/" + filePath)
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
	lastFile := ""
	directory := ""
	files := strings.Split(filePath, "/")
	if len(files) > 1 {
		lastFile = files[len(files)-1]
		directory = strings.Join(files[0:len(files)-1], "/")
		err := os.MkdirAll(s.rootFolder+"/"+directory, 0770)
		if err != nil {
			return nil, err
		}
		directory += "/"
	}
	w, err := os.Create(s.rootFolder + "/" + directory + lastFile)
	if err != nil {
		return nil, err
	}
	return w, nil
}
