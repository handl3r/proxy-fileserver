package adapter

import (
	"io"
	"os"
	"proxy-fileserver/repository"
	"strings"
)

type LocalFileSystem struct {
	rootFolder   string
	fileInfoRepo *repository.FileInfoRepository
}

func NewLocalFileSystem(rootFolder string) *LocalFileSystem {
	return &LocalFileSystem{
		rootFolder: rootFolder,
	}
}

// TODO check if should remove parent folders
func (s *LocalFileSystem) Delete(filePath string) error {
	err := os.Remove(s.rootFolder + "/" + filePath)
	if err != nil {
		return err
	}
	return nil
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

func (s *LocalFileSystem) NewFile(filePath string) (io.Writer, error) {
	lastFile := ""
	directory := ""
	files := strings.Split(filePath, "/")
	if len(files) > 1 {
		lastFile = files[len(files)-1]
		directory = strings.Join(files[0:len(files)-1], "/")
		err := os.MkdirAll(s.rootFolder+"/"+directory, 0777)
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

func (s *LocalFileSystem) RenameFile(oldPath string, newPath string) error {
	return os.Rename(s.rootFolder+"/"+oldPath, s.rootFolder+"/"+newPath)
}
