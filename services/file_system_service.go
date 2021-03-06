package services

import (
	"io"
	"proxy-fileserver/adapter"
	"proxy-fileserver/common/log"
	"proxy-fileserver/enums"
	"strings"
)

type FileSystemService struct {
	rootFolder            string
	GoogleDriveFileSystem adapter.GoogleDriveFileSystem
	LocalFileSystem       adapter.LocalFileSystem
}

func NewLocalStorageService(rootFolder string, googleDrive adapter.GoogleDriveFileSystem, localStorage adapter.LocalFileSystem) *FileSystemService {
	return &FileSystemService{
		GoogleDriveFileSystem: googleDrive,
		LocalFileSystem:       localStorage,
	}
}

// StreamFile Public method control all process to stream file to client and sync file from drive to local server
func (s *FileSystemService) StreamFile(outStreamHttp io.Writer, filePath string) enums.Response {
	existed, err := s.LocalFileSystem.IsExisted(filePath)
	if err != nil {
		log.Errorf("Can not check if file exist from file %s with error: %v", filePath, err)
		return enums.ErrorSystem
	}
	if existed {
		err := s.StreamFromFileSystem(outStreamHttp, filePath)
		if err != nil {
			log.Errorf("Failure when streaming file from local server to client with filepath: %s", filePath)
			return enums.ErrorSystem
		}
		log.Infof("Finish stream file %s from  local file system to client", filePath)
		return nil
	}

	id, srcStream, err := s.GoogleDriveFileSystem.GetStreamSourceByFilePath(filePath)
	if err == enums.ErrFileNotExist {
		return enums.ErrorNoContent
	}
	if err != nil {
		log.Errorf("Can not get source stream from drive with file path %s, id %s , error: %v", filePath, id, err)
		return enums.ErrorSystem
	}
	// TODO make new function to stream from reader to multi writer with condition: when 1 writer is failure, another still continue
	go func() {
		err = s.StreamFromDriveToLocalFileSystem(id, filePath)
		if err == nil {
			log.Infof("Finished getting new file from drive to local server with filepath: %s, id: %s", filePath, id)
		}

	}()
	_, err = io.Copy(outStreamHttp, srcStream)
	if err != nil {
		log.Errorf("Error when streaming file from driver to client with filePath: %s, id: %s, error: %s", filePath, id, err)
		return enums.ErrorSystem
	}
	log.Errorf("Finished stream file from drive to client with filepath: %s, id: %s", filePath, id)
	return nil
}

func (s *FileSystemService) StreamFromFileSystem(outStream io.Writer, filePath string) error {
	//existed, err := s.LocalFileSystem.IsExisted(filePath)
	//if err != nil {
	//	log.Errorf("Can not check if file exist from file %s with error: %v", filePath, err)
	//	return err
	//}
	//if existed {
	srcStream, err := s.LocalFileSystem.GetStreamSourceByFilePath(filePath)
	if err != nil {
		log.Errorf("Can not get stream source from local file %s with error: %v", filePath, err)
		return err
	}
	_, err = io.Copy(outStream, srcStream)
	if err != nil {
		log.Errorf("Stream file %s from local file system to client error: %v", filePath, err)
		return err
	}
	//}
	return nil
}

func (s *FileSystemService) StreamFromDriveToLocalFileSystem(id string, filePath string) error {
	newFileStream, err := s.LocalFileSystem.NewFile(filePath)
	if err != nil {
		log.Errorf("Can not create new file with filePath: %s, error: %v", filePath, err)
		return err
	}
	srcStreamDrive, err := s.GoogleDriveFileSystem.GetStreamBySourceByID(id)
	if err != nil {
		log.Errorf("Can not get srcStream from drive with filepath: %s, id: %s, error: %s", filePath, id, err)
		return err
	}
	_, err = io.Copy(newFileStream, srcStreamDrive)
	if err != nil {
		log.Errorf("Can not stream file from drive to local with filepath: %s, id: %s, err: %v", filePath, id, err)
		return err
	}
	return nil
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
